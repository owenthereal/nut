package main

import (
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strconv"
	"strings"
)

func rewrite(pkgs []*Pkg, prefix string) error {
	importPaths := make([]string, 0)
	files := make([]string, 0)
	for _, pkg := range pkgs {
		importPaths = append(importPaths, pkg.ImportPath)
		files = append(files, pkg.GoFiles()...)
	}

	for _, file := range files {
		err := rewriteGoFile(file, prefix, importPaths)
		if err != nil {
			return err
		}
	}

	return nil
}

func rewriteGoFile(file, prefix string, importPaths []string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var changed bool
	for _, s := range f.Imports {
		name, err := strconv.Unquote(s.Path.Value)
		if err != nil {
			return err
		}

		q := qualify(unqualify(name), prefix, importPaths)
		if q != name {
			s.Path.Value = strconv.Quote(q)
			changed = true
		}
	}

	if !changed {
		return nil
	}

	wpath := file + ".temp"
	w, err := os.Create(wpath)
	if err != nil {
		return err
	}

	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(w, fset, f)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return os.Rename(wpath, file)
}

const sep = "/vendor/_nuts/"

func unqualify(importPath string) string {
	if i := strings.LastIndex(importPath, sep); i != -1 {
		importPath = importPath[i+len(sep):]
	}
	return importPath
}

func qualify(importPath, pkg string, paths []string) string {
	if containsPathPrefix(paths, importPath) {
		return pkg + sep + importPath
	}

	return importPath
}

func containsPathPrefix(pats []string, s string) bool {
	for _, pat := range pats {
		if pat == s || strings.HasPrefix(s, pat+"/") {
			return true
		}
	}
	return false
}
