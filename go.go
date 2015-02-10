package main

import (
	"os"
	"os/exec"
	"strings"
)

func goGet(dir string, importPath ...string) error {
	if len(importPath) == 0 {
		return nil
	}

	cmd := []string{"get", "-d", "-t"}
	cmd = append(cmd, importPath...)

	c := newGoCmd(cmd...)
	c.Dir = dir

	return c.Run()
}

func runGoCmd(args ...string) error {
	c := newGoCmd(args...)

	return c.Run()
}

func newGoCmd(args ...string) *exec.Cmd {
	c := exec.Command("go", args...)
	c.Env = goCmdEnv()
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c
}

func goCmdEnv() []string {
	return append(envNoGopath(), "GOPATH="+setting.WorkDir())
}

func envNoGopath() (a []string) {
	for _, s := range os.Environ() {
		if !strings.HasPrefix(s, "GOPATH=") {
			a = append(a, s)
		}
	}
	return a
}
