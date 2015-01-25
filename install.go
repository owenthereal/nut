package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

var install = cli.Command{
	Name:   "install",
	Usage:  "install this project's dependencies",
	Action: runInstall,
}

func runInstall(c *cli.Context) {
	config := setting.Config()
	pl := &PkgLoader{
		Deps: config.Deps,
	}
	pkgs, err := pl.Load()
	check(err)

	l := pkgLister{
		Env: os.Environ(),
	}
	currentPkg, err := l.List(".")
	check(err)

	err = rewrite(pkgs, currentPkg[0].ImportPath)
	check(err)

	err = os.RemoveAll(setting.VendorDir())
	check(err)

	err = copyPkgs(pkgs)
	check(err)
}

func copyPkgs(pkgs []*Pkg) error {
	return copyDir(filepath.Join(setting.WorkDir(), "src"), setting.VendorDir())
}

func copyDir(source string, dest string) (err error) {
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()

		// ignore dir starting with . or _
		c := obj.Name()[0]
		if obj.IsDir() && (c == '.' || c == '_') {
			continue
		}

		if obj.IsDir() {
			err = copyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				return err
			}
		}

	}

	return
}

func copyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}
