package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gophergala/nut/vendor/_nuts/github.com/codegangsta/cli"
)

var installCmd = cli.Command{
	Name:   "install",
	Usage:  "install this project's dependencies",
	Action: runInstall,
}

func runInstall(c *cli.Context) {
	config := setting.Config()
	if len(config.Deps) == 0 {
		return
	}

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

	if strings.HasSuffix(dest, ".go") {
		err = copyWithoutImportComment(destfile, sourcefile)
	} else {
		_, err = io.Copy(destfile, sourcefile)
	}

	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func copyWithoutImportComment(w io.Writer, r io.Reader) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		_, err := w.Write(append(stripImportComment(sc.Bytes()), '\n'))
		if err != nil {
			return err
		}
	}
	return nil
}

const (
	importAnnotation = `import\s+(?:"[^"]*"|` + "`[^`]*`" + `)`
	importComment    = `(?://\s*` + importAnnotation + `\s*$|/\*\s*` + importAnnotation + `\s*\*/)`
)

var (
	importCommentRE = regexp.MustCompile(`^\s*(package\s+\w+)\s+` + importComment + `(.*)`)
	pkgPrefix       = []byte("package ")
)

// stripImportComment returns line with its import comment removed.
// If s is not a package statement containing an import comment,
// it is returned unaltered.
// See also http://golang.org/s/go14customimport.
func stripImportComment(line []byte) []byte {
	if !bytes.HasPrefix(line, pkgPrefix) {
		// Fast path; this will skip all but one line in the file.
		// This assumes there is no whitespace before the keyword.
		return line
	}
	if m := importCommentRE.FindSubmatch(line); m != nil {
		return append(m[1], m[2]...)
	}
	return line
}
