package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func NewPkg(dir, importPath, rev string, vcs *VCS) *Pkg {
	return &Pkg{
		Dir:        dir,
		ImportPath: importPath,
		Rev:        rev,
		vcs:        vcs,
		goFiles:    make(map[string]bool),
	}
}

type Pkg struct {
	Dir        string
	ImportPath string
	Rev        string
	vcs        *VCS
	goFiles    map[string]bool
}

func (p *Pkg) addGoFiles(files []string) {
	for _, file := range files {
		if _, ok := p.goFiles[file]; !ok {
			p.goFiles[file] = true
		}
	}
}

func (p *Pkg) GoFiles() []string {
	files := make([]string, 0)
	for file, _ := range p.goFiles {
		files = append(files, file)
	}

	return files
}

type PkgLoader struct {
	Deps ConfigDeps
}

func (pl *PkgLoader) Load() ([]*Pkg, error) {
	var names []string
	for dep, rev := range pl.Deps {
		if rev == "" {
			rev = "latest"
		}

		fmt.Printf("Downloading %s@%s\n", dep, rev)
		err := goGet(dep)
		if err != nil {
			return nil, err
		}

		names = append(names, dep)
	}

	// declared dependencies
	ps, err := listPkgs(names...)
	if err != nil {
		return nil, err
	}

	// checkout revisions
	for _, p := range ps {
		rev := pl.Deps[p.ImportPath]
		if rev == "" {
			continue
		}

		vcs, _, err := VCSFromDir(p.Dir, filepath.Join(p.Root, "src"))
		if err != nil {
			return nil, err
		}

		err = vcs.Checkout(p.Dir, rev)
		if err != nil {
			return nil, err
		}
	}

	pkgs, err := pl.loadPkgs(ps)
	if err != nil {
		return nil, err
	}

	return pkgs, nil
}

func (pl *PkgLoader) loadPkgs(a []*pkg) ([]*Pkg, error) {
	pkgMap := make(map[string]*Pkg)
	seen := make(map[string]bool)

	err := pl.doLoadPkgs(a, pkgMap, seen)
	if err != nil {
		return nil, err
	}

	var pkgs []*Pkg
	for _, pkg := range pkgMap {
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

func (pl *PkgLoader) doLoadPkgs(ps []*pkg, pkgMap map[string]*Pkg, seen map[string]bool) error {
	// filter not-seen third party dependencies
	for _, p := range ps {
		if p.Standard {
			continue
		}

		if _, ok := seen[p.ImportPath]; ok {
			continue
		}

		seen[p.ImportPath] = true

		// for itself
		err := pl.doLoadAPkg(p, pkgMap)
		if err != nil {
			return err
		}

		// for its test dependencies
		testImports := append(p.TestImports, p.XTestImports...)
		testPkgs, err := listPkgs(testImports...)
		if err != nil {
			return err
		}

		err = pl.doLoadPkgs(testPkgs, pkgMap, seen)
		if err != nil {
			return err
		}

		// for its dependencies
		depPkgs, err := listPkgs(p.Deps...)
		if err != nil {
			return err
		}

		err = pl.doLoadPkgs(depPkgs, pkgMap, seen)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pl *PkgLoader) doLoadAPkg(p *pkg, pkgMap map[string]*Pkg) error {
	err := goGet(p.ImportPath)
	if err != nil {
		return err
	}

	vcsCmd, importPath, err := VCSFromDir(p.Dir, filepath.Join(p.Root, "src"))
	if err != nil {
		return err
	}

	pkg, ok := pkgMap[importPath]
	if !ok {
		rev := pl.Deps[importPath]
		pkg = NewPkg(p.Dir, importPath, rev, vcsCmd)
		pkgMap[importPath] = pkg
	}
	pkg.addGoFiles(p.AllGoFiles())

	return nil
}

type pkg struct {
	Dir        string
	Root       string
	ImportPath string
	Deps       []string
	Standard   bool

	GoFiles        []string
	CgoFiles       []string
	IgnoredGoFiles []string

	TestGoFiles  []string
	TestImports  []string
	XTestGoFiles []string
	XTestImports []string

	Error struct {
		Err string
	}
}

func (p *pkg) AllGoFiles() (a []string) {
	a = append(a, pathOf(p.Dir, p.GoFiles)...)
	a = append(a, pathOf(p.Dir, p.CgoFiles)...)
	a = append(a, pathOf(p.Dir, p.TestGoFiles)...)
	a = append(a, pathOf(p.Dir, p.XTestGoFiles)...)
	a = append(a, pathOf(p.Dir, p.IgnoredGoFiles)...)

	return a
}

func pathOf(dir string, files []string) []string {
	paths := make([]string, 0)
	for _, f := range files {
		paths = append(paths, filepath.Join(dir, f))
	}

	return paths
}

func listPkgs(name ...string) ([]*pkg, error) {
	l := pkgLister{
		Env: goCmdEnv(),
	}
	return l.List(name...)
}

type pkgLister struct {
	Env []string
}

func (l *pkgLister) List(name ...string) (a []*pkg, err error) {
	if len(name) == 0 {
		return nil, nil
	}

	args := []string{"list", "-e", "-json"}
	cmd := exec.Command("go", append(args, name...)...)
	cmd.Env = l.Env

	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(r)
	for {
		info := new(pkg)
		err = d.Decode(info)
		if err == io.EOF {
			break
		}
		if err != nil {
			info.Error.Err = err.Error()
		}
		a = append(a, info)
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return a, nil
}
