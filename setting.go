package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

var setting *Setting

func init() {
	pwd, err := os.Getwd()
	check(err)

	setting = &Setting{
		ProjectDir:     pwd,
		ConfigFile:     filepath.Join(pwd, "Nut.toml"),
		ConfigLockFile: filepath.Join(pwd, "Nut.lock"),
	}
}

type Setting struct {
	ProjectDir     string
	ConfigFile     string
	ConfigLockFile string
	goPath         string
	manifest       *Manifest
}

func (s *Setting) WorkDir() string {
	if s.goPath == "" {
		temp, err := ioutil.TempDir("", "nut")
		check(err)

		s.goPath = temp
	}

	return s.goPath
}

func (s *Setting) VendorDir() string {
	return filepath.Join(setting.ProjectDir, "internal")
}

func (s *Setting) Manifest() *Manifest {
	if s.manifest == nil {
		var err error
		s.manifest, err = loadManifest()
		check(err)
	}

	return s.manifest
}
