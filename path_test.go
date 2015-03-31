package main

import (
	"os"
	"testing"
)

// TODO: add TestGetAllDeepPaths

func TestGetDeepPaths(t *testing.T) {
	gopath := setting.WorkDir()
	importPath := "github.com/jingweno/nut"
	absPath := gopath + "/src/" + importPath
	examplePath := absPath + "/example"

	makeDir(examplePath, "test1/test5/test6")
	makeDir(examplePath, "test1/test5/test7")
	makeDir(examplePath, "test1/test5/test8")
	makeDir(examplePath, "test2")
	makeDir(examplePath, "test3/test4")

	paths, err := getDeepPaths(importPath)
	if err != nil {
		t.Fatalf("err must be nil, actual=%s", err.Error())
	}

	pathMap := make(map[string]bool)
	for _, path := range paths {
		pathMap[path] = true
	}

	expects := []string{
		importPath,
		importPath + "/example",
		importPath + "/example/test1",
		importPath + "/example/test2",
		importPath + "/example/test3",
		importPath + "/example/test3/test4",
		importPath + "/example/test1/test5",
		importPath + "/example/test1/test5/test6",
		importPath + "/example/test1/test5/test7",
		importPath + "/example/test1/test5/test8",
	}

	if len(paths) != len(expects) {
		t.Fatalf("paths must have [%d] slice, actual=%v", len(expects), paths)
	}

	for _, expect := range expects {
		if _, ok := pathMap[expect]; !ok {
			t.Fatalf("paths must contain %s, paths=%v", expect, pathMap)
		}
	}
}

func makeDir(path, dir string) {
	os.MkdirAll(path+"/"+dir, 0775)
}

func TestIsIgnoreName(t *testing.T) {
	gopath := "/tmp/example/gopath"
	importPath := "github.com/jingweno/nut"
	absPath := gopath + "/src/" + importPath

	normalDir := absPath + "/myProject"
	if isIgnoreName(absPath, normalDir) == true {
		t.Fatalf("normalDir must not be ignored dir, normalDir=%s", normalDir)
	}

	gitDir := absPath + "/.git"
	if isIgnoreName(absPath, gitDir) != true {
		t.Fatalf("gitDir must be ignored dir, gitDir=%s", gitDir)
	}

	underScoreDir := absPath + "/_vendor"
	if isIgnoreName(absPath, underScoreDir) != true {
		t.Fatalf("underScoreDir must be ignored dir, underScoreDir=%s", underScoreDir)
	}
}

func TestRenameToAbsPath(t *testing.T) {
	importPath := "github.com/jingweno/nut"
	gopath := setting.WorkDir()
	renamedPath := renameToAbsPath(importPath)

	absPath := gopath + "/src/" + importPath
	if renamedPath != absPath {
		t.Fatalf("renamedPath must be %s, actual=%s", absPath, renamedPath)
	}
}

func TestRenameToImportPath(t *testing.T) {
	importPath := "github.com/jingweno/nut"
	gopath := setting.WorkDir()
	absPath := gopath + "/src/" + importPath
	renamedPath := renameToImportPath(absPath)
	if renamedPath != importPath {
		t.Fatalf("renamedPath must be %s, actual=%s", importPath, renamedPath)
	}
}

func TestGetSrcPath(t *testing.T) {
	gopath := "/tmp/example/gopath"
	srcPath := getSrcPath(gopath)
	if srcPath != "/tmp/example/gopath/src" {
		t.Fatalf("srcPath must be %s/src, actual=%s", gopath, srcPath)
	}
}
