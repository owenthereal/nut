package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/jingweno/nut/vendor/_nuts/github.com/codegangsta/cli"
)

// dependencyFinder is an interface for finding project dependencies.
// This can later be used to migrate from godep, etc
type dependencyFinder interface {
	FindDependencies(vcsRefs bool) (map[string]string, error)
}

var initCmd = cli.Command{
	Name:   "init",
	Usage:  "init nut on an existing go project, writing all dependencies to Nut.toml",
	Action: runInit,
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "ignore, i",
			Usage: "import path prefixes to ignore",
			Value: &cli.StringSlice{},
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "verbose output",
		},
		cli.BoolFlag{
			Name:  "refs, r",
			Usage: "extract VCS refs for packages",
		},
	},
}

// depFinder finds the dependencies of an existing project by traversing go deps
type depFinder struct {
	// the seen dependencies during traversal
	deps map[string]*pkg

	// the root package. we ignore all its sub packages
	root string

	// a list of ignore import prefixes (e.g. gitlab.mydomain.com)
	ignored []string

	// the current working directory (this changes durgin traversal)
	wd string

	// output verbose messages during traversal
	verbose bool

	lister pkgLister
}

// isIgnored checks if an import path start with a prefix that has been set to be ignored
func (c *depFinder) isIgnored(path string) bool {

	for _, prefix := range c.ignored {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

// findDeps recursively finds all the dependencies of a package that are not part of the same
// project or the standard library
func (c *depFinder) findDeps(p *pkg) error {

	if c.root == "" {
		c.root = p.ImportPath
	}

	deps, err := c.lister.List(p.Deps...)
	if err != nil {
		return err
	}

	for _, dep := range deps {
		if dep.Standard {
			continue
		}
		if !strings.HasPrefix(dep.ImportPath, c.root) && !c.isIgnored(dep.ImportPath) {

			if _, found := c.deps[dep.ImportPath]; !found {
				c.deps[dep.ImportPath] = dep
				if c.verbose {
					fmt.Println(dep.ImportPath)
				}
				c.findDeps(dep)
			}
		}
	}

	return nil
}

// findDepsInDir finds all the deps for the packages under the dir in path.
// this is used to recursively traverse all packages in a project
func (c *depFinder) findDepsInDir(pth string) error {

	c.wd = pth
	os.Chdir(c.wd)

	pkgs, err := c.lister.List("")
	if err != nil {
		return err
	}
	for _, p := range pkgs {
		c.findDeps(p)
	}

	files, err := ioutil.ReadDir(pth)

	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), ".") || strings.HasPrefix(file.Name(), "_") {
			continue
		}

		err := c.findDepsInDir(path.Join(pth, file.Name()))
		if err != nil {
			fmt.Println("error traversing", file.Name(), ": ", err)
			return err
		}

	}

	return nil

}

// FindDependencies traverses the current WD and retrieves a dependency map (package => revision)
// for all 3rd party libraries.
// VCS refs are retrieved only if the vcsRefs flag is set to true
func (c *depFinder) FindDependencies(vcsRefs bool) (map[string]string, error) {
	err := c.findDepsInDir(c.wd)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]string)
	if c.verbose && vcsRefs {
		fmt.Println("Finding VCS refs for packages")
	}

	for dep, pkg := range c.deps {
		tag := ""
		if vcsRefs {
			_, vc, err := VCSForImportPath(dep)

			if err == nil {
				tag, _ = vc.Identify(pkg.Dir)
			} else {
				fmt.Printf("Error extracting tag for %s: %s\n", dep, err)
			}

		}
		ret[dep] = tag

	}

	return ret, nil
}

func runInit(c *cli.Context) {

	//the first arg can be the path of a project
	if len(c.Args()) > 0 {
		pth := c.Args().First()
		if pth != "" && pth != "." {
			err := os.Chdir(pth)
			if err != nil {
				fmt.Printf("Invalid path %s (%s)\n", pth, err)
				os.Exit(1)
			}
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := &depFinder{
		deps:    make(map[string]*pkg),
		ignored: c.StringSlice("ignored"),
		wd:      cwd,
		verbose: c.Bool("verbose"),
		lister: pkgLister{
			Env: os.Environ(),
		},
	}

	deps, err := ctx.FindDependencies(c.Bool("refs"))

	// create a manifest from the gathered data
	manifest := Manifest{
		App: ManifestApp{
			Name:    path.Base(cwd),
			Authors: []string{"Your Name <you@example.com>"},
			Version: "0.0.1",
		},
		Deps: ManifestDeps(deps),
	}

	os.Chdir(cwd)

	err = manifest.write()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("New manifest file written to", setting.ConfigFile)

}
