package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func NewPkg(dir, importPathRoot, importPath, rev string) *Pkg {
	return &Pkg{
		Dir:            dir,
		ImportPathRoot: importPathRoot,
		ImportPath:     importPath,
		Rev:            rev,
		goFiles:        newSet(),
	}
}

type Pkg struct {
	Dir            string
	ImportPathRoot string
	ImportPath     string
	Rev            string
	goFiles        set
}

func (p Pkg) String() string {
	return fmt.Sprintf("Dir=%s, ImportPathRoot=%s, ImportPath=%s, Rev=%s, Files=%v", p.Dir, p.ImportPathRoot, p.ImportPath, p.Rev, p.GoFiles())
}

func (p *Pkg) addGoFiles(files []string) {
	for _, file := range files {
		p.goFiles.Add(file)
	}
}

func (p *Pkg) GoFiles() []string {
	return p.goFiles.ToSliceString()
}

type PkgLoader struct {
	GoPath string
}

func (pl *PkgLoader) Load(importPaths ...string) ([]*Pkg, error) {
	err := goGet(pl.GoPath, importPaths...)
	if err != nil {
		return nil, err
	}

	ps, err := listPkgs(importPaths...)
	if err != nil {
		return nil, err
	}

	pkgs, err := pl.loadPkgs(ps)
	if err != nil {
		return nil, err
	}

	return pkgs, nil
}

func (pl *PkgLoader) loadPkgs(a []*pkg) ([]*Pkg, error) {
	pkgMap := make(map[string]*Pkg)
	seen := newSet()

	err := pl.recursiveLoadPkgs(a, pkgMap, seen)
	if err != nil {
		return nil, err
	}

	var pkgs []*Pkg
	for _, pkg := range pkgMap {
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

func (pl *PkgLoader) getUnloadPkgs(ps []*pkg, seen set) []*pkg {
	unloadPkgsMap := make(map[string]*pkg)
	for _, p := range ps {
		if p.Standard {
			continue
		}

		if seen.Contains(p.ImportPath) {
			continue
		}

		if _, ok := unloadPkgsMap[p.ImportPath]; !ok {
			unloadPkgsMap[p.ImportPath] = p
		}
	}

	var unloadPkgs []*pkg
	for _, p := range unloadPkgsMap {
		unloadPkgs = append(unloadPkgs, p)
	}

	return unloadPkgs
}

func (pl *PkgLoader) getDepPkgs(ps []*pkg) ([]*pkg, error) {
	s := newSet()
	for _, p := range ps {
		// for its test dependencies
		imports := append(p.TestImports, p.XTestImports...)
		// for its dependencies
		imports = append(imports, p.Deps...)

		for _, i := range imports {
			s.Add(i)
		}
	}

	return listPkgs(s.ToSliceString()...)
}

func (pl *PkgLoader) recursiveLoadPkgs(ps []*pkg, pkgMap map[string]*Pkg, seen set) error {
	unloadPkgs := pl.getUnloadPkgs(ps, seen)
	if len(unloadPkgs) == 0 {
		return nil
	}

	err := pl.doLoadPkgs(unloadPkgs, pkgMap, seen)
	if err != nil {
		return err
	}

	depPkgs, err := pl.getDepPkgs(ps)
	if err != nil {
		return err
	}

	return pl.recursiveLoadPkgs(depPkgs, pkgMap, seen)
}

func (pl *PkgLoader) getImportPaths(ps []*pkg) []string {
	var importPaths []string
	for _, p := range ps {
		importPaths = append(importPaths, p.ImportPath)
	}

	return importPaths
}

func (pl *PkgLoader) doLoadPkgs(ps []*pkg, pkgMap map[string]*Pkg, seen set) error {
	importPaths := pl.getImportPaths(ps)

	// make sure dependencies are downloaded
	err := goGet(pl.GoPath, importPaths...)
	if err != nil {
		return err
	}

	return pl.cachePkgs(ps, pkgMap, seen)
}

func (pl *PkgLoader) cachePkgs(ps []*pkg, pkgMap map[string]*Pkg, seen set) error {
	for _, p := range ps {
		vcs, importPathRoot, err := VCSFromDir(p.Dir, filepath.Join(p.Root, "src"))
		if err != nil {
			return err
		}

		pkg, ok := pkgMap[p.ImportPath]
		if !ok {
			rev, err := vcs.Identify(p.Dir)
			if err != nil {
				return err
			}

			pkg = NewPkg(p.Dir, importPathRoot, p.ImportPath, rev)
			pkgMap[p.ImportPath] = pkg
		}

		pkg.addGoFiles(p.AllGoFiles())

		seen.Add(p.ImportPath)
	}

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
	name = addSubPaths(name)
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

func addSubPaths(name []string) []string {
	var paths []string
	for _, path := range name {
		paths = append(paths, addSubPath(path))
	}
	return paths
}

func addSubPath(path string) string {
	return path + "/..."
}
