package main

import (
	"fmt"
	"runtime"

	"github.com/gophergala/nut/vendor/_nuts/github.com/codegangsta/cli"
	"github.com/gophergala/nut/vendor/_nuts/golang.org/x/tools/go/vcs"
)

var listCmd = cli.Command{
	Name:   "list",
	Usage:  "list vendored dependencies",
	Action: runList,
}

func runList(c *cli.Context) {
	p, err := NewProject()
	check(err)

	r, err := vcs.RepoRootForImportPath("github.com/codegangsta/cli", false)
	check(err)

	err = r.VCS.Create("../test", r.Repo)
	check(err)

	fmt.Println(p.ImportPath)
	fmt.Println(r)
	fmt.Println(runtime.Version())
}
