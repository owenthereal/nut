package main

import (
	"os"
	"os/exec"
	"strings"
)

func runGoCmd(args ...string) error {
	c := newGoCmd(args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}

func newGoCmd(args ...string) *exec.Cmd {
	c := exec.Command("go", args...)
	c.Env = goCmdEnv()

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
