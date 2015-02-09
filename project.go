package main

import (
	"fmt"
	"os"
)

type Project struct {
	ImportPath string
	Deps       []*Dep
}

type Dep struct {
	ImportPath string
	Rev        string
	Deps       []*Dep
}

func NewProject() (*Project, error) {
	l := pkgLister{
		Env: os.Environ(),
	}

	pkgs, err := l.List(".")
	if err != nil {
		return nil, err
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("Error loading project info")
	}

	p := &Project{ImportPath: pkgs[0].ImportPath}

	//for n, _ := range setting.Config().Deps {
	//}

	return p, nil
}
