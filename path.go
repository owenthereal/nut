package main

import (
	"os"
	"path/filepath"
	"strings"
)

func getAllDeepPaths(importPaths []string) []string {
	var importPathDirs []string
	for _, path := range importPaths {
		dirs, err := getDeepPaths(path)
		check(err)
		importPathDirs = append(importPathDirs, dirs...)
	}
	return importPathDirs
}

func getDeepPaths(importPath string) ([]string, error) {
	rootPath := renameToAbsPath(importPath)
	var dirs []string
	walkFn := func(path string, info os.FileInfo, err error) error {
		if isIgnoreName(rootPath, path) {
			return nil
		}
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}
		if !stat.IsDir() {
			return nil
		}
		path = renameToImportPath(path)
		dirs = append(dirs, path)
		return nil
	}
	err := filepath.Walk(rootPath, walkFn)
	if err != nil {
		return dirs, err
	}
	return dirs, nil
}

func isIgnoreName(root, path string) bool {
	// check first string
	path = strings.Replace(path, root+"/", "", 1)
	for _, r := range path {
		s := string(r)
		if s == "." || s == "_" {
			return true
		}
		break
	}
	return false
}

func renameToAbsPath(path string) string {
	return getSrcPath(setting.WorkDir()) + "/" + path
}

func renameToImportPath(path string) string {
	return strings.Replace(path, getSrcPath(setting.WorkDir())+"/", "", 1)
}

func getSrcPath(path string) string {
	return path + "/src"
}
