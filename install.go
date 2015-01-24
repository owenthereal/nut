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
	config, err := loadConfig()
	check(err)

	for d, c := range config.Deps {
		fmt.Printf("%s@%s\n", d, c)

		err := goCmd("get", "-d", d)
		check(err)

		err = copyDep(d)
		check(err)
	}
}

func copyDep(dep string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	return copyDir(filepath.Join(goPath, "src"), filepath.Join(dir, "vendor", "_nuts"))
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
