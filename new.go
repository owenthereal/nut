package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gophergala/nut/vendor/_nuts/github.com/codegangsta/cli"
)

const nutTomlTmpl = `[application]

name = "{{.Name}}"
version = "0.0.1"
authors = ["Your Name <you@example.com>"]
`

const readmeMd = `# {{.Name}}

This is an awesome Go projeect.
`

const gitIgnore = `{{.Name}}
`

const mainGo = `package main

import "fmt"

func main() {
	fmt.Println("Hello Gophers!")
}
`

var newCmd = cli.Command{
	Name:   "new",
	Usage:  "create a nut project",
	Action: runNew,
}

func runNew(c *cli.Context) {
	if len(c.Args()) == 0 {
		check(fmt.Errorf("Usage: \n    nut new <path>"))
	}

	name := c.Args()[0]

	err := os.Mkdir(name, os.ModePerm)
	check(err)

	err = createNutToml(name)
	check(err)

	err = createReadmeMd(name)
	check(err)

	err = createGitIgnore(name)
	check(err)

	err = createMainGo(name)
	check(err)
}

func createNutToml(name string) error {
	return createFile(name, nutTomlTmpl, "Nut.toml")
}

func createReadmeMd(name string) error {
	return createFile(name, readmeMd, "README.md")
}

func createGitIgnore(name string) error {
	return createFile(name, gitIgnore, ".gitignore")
}

func createMainGo(name string) error {
	return ioutil.WriteFile(filepath.Join(name, "main.go"), []byte(mainGo), 0644)
}

func createFile(name, templ, fileName string) error {
	t, err := template.New(name).Parse(templ)
	if err != nil {
		return err
	}

	n := struct {
		Name string
	}{
		name,
	}
	var b bytes.Buffer
	err = t.Execute(&b, n)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(name, fileName), b.Bytes(), 0644)
}
