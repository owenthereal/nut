package main

import (
	"fmt"
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

	deps := make([]string, 0)
	for d, c := range config.Deps {
		cc := c
		if cc == "" {
			cc = "latest"
		}
		fmt.Printf("%s@%s\n", d, cc)

		err := runGoCmd("get", "-d", "-t", d)
		check(err)

		deps = append(deps, d)
	}

	p, err := loadPackages(deps...)
	check(err)

	pkgs, err := loadPkgs(p)
	check(err)

	err = rewrite(pkgs, "github.com/gophergala/nut")
	check(err)

	err = copyDeps()
	check(err)
}

func copyDeps() error {
	return copyDir(filepath.Join(setting.WorkDir(), "src"), filepath.Join(setting.ProjectDir, "vendor", "_nuts"))
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
