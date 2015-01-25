package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/vcs"
)

func NewPkg(importPath string) *Pkg {
	return &Pkg{
		ImportPath: importPath,
		goFiles:    make(map[string]bool),
	}
}

type Pkg struct {
	ImportPath string
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

func loadAPkg(p *Package, pkgMap map[string]*Pkg) error {
	err := runGoCmd("get", "-d", "-t", p.ImportPath)
	if err != nil {
		return err
	}

	_, reporoot, err := vcs.FromDir(p.Dir, filepath.Join(p.Root, "src"))
	if err != nil {
		return err
	}

	pkg, ok := pkgMap[reporoot]
	if !ok {
		pkg = NewPkg(reporoot)
		pkgMap[reporoot] = pkg
	}
	pkg.addGoFiles(p.AllGoFiles())

	return nil
}

func loadPkg(ps []*Package, pkgMap map[string]*Pkg, seen map[string]bool) error {
	for _, p := range ps {
		if p.Standard {
			continue
		}

		if _, ok := seen[p.ImportPath]; ok {
			continue
		}
		seen[p.ImportPath] = true

		// for itself
		err := loadAPkg(p, pkgMap)
		if err != nil {
			return err
		}

		// for its test dependencies
		testImports := p.TestImports
		testImports = append(testImports, p.XTestImports...)

		testPkgs, err := loadPackages(testImports...)
		if err != nil {
			return err
		}

		err = loadPkg(testPkgs, pkgMap, seen)
		if err != nil {
			return err
		}

		// for its dependencies
		depPkgs, err := loadPackages(p.Deps...)
		if err != nil {
			return err
		}

		err = loadPkg(depPkgs, pkgMap, seen)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadPkgs(a []*Package) ([]*Pkg, error) {
	pkgMap := make(map[string]*Pkg)
	seen := make(map[string]bool)

	err := loadPkg(a, pkgMap, seen)
	if err != nil {
		return nil, err
	}

	var pkgs []*Pkg
	for _, pkg := range pkgMap {
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

type Package struct {
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

func (p *Package) AllGoFiles() (a []string) {
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

func loadPackages(name ...string) (a []*Package, err error) {
	if len(name) == 0 {
		return nil, nil
	}

	args := []string{"list", "-e", "-json"}
	cmd := newGoCmd(append(args, name...)...)
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
		info := new(Package)
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
