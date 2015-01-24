package main

import (
	"os"
	"os/exec"
	"strings"
)

func goCmd(args ...string) error {
	c := exec.Command("go", args...)
	c.Env = append(envNoGopath(), "GOPATH="+setting.WorkDir())
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}

func envNoGopath() (a []string) {
	for _, s := range os.Environ() {
		if !strings.HasPrefix(s, "GOPATH=") {
			a = append(a, s)
		}
	}
	return a
}
