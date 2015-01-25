package main

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

func Test_rewriteGoFile(t *testing.T) {
	f, _ := ioutil.TempFile("", "rewriteGoFile.go")
	f.Close()

	c := `package main

import(
	"github.com/octokit/go-octokit"
)

func main() {
	fmt.Println("hi")
}

	`
	ioutil.WriteFile(f.Name(), []byte(c), os.ModePerm)

	err := rewriteGoFile(f.Name(), "github.com/gophergala/nut", []string{"github.com/octokit/go-octokit"})
	if err != nil {
		t.Fatal("err isn't nil for rewriteGoFile")
	}

	fset := token.NewFileSet()
	ff, _ := parser.ParseFile(fset, f.Name(), nil, parser.ParseComments)

	if len(ff.Imports) != 1 {
		t.Fatal("there should be only one import")
	}

	i, _ := strconv.Unquote(ff.Imports[0].Path.Value)
	if i != "github.com/gophergala/nut/vendor/_nuts/github.com/octokit/go-octokit" {
		t.Fatalf("import is wrong: %s", i)
	}
}
