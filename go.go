package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var goPath string

func init() {
	goPath = tempDir()
}

func goCmd(args ...string) error {
	c := exec.Command("go", args...)
	c.Env = append(envNoGopath(), "GOPATH="+goPath)
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

func tempDir() string {
	dir, err := ioutil.TempDir("", "nut")
	check(err)

	return dir
}
