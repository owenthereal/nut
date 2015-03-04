package main

import (
	"fmt"
	"os"

	"github.com/jingweno/nut/vendor/_nuts/github.com/BurntSushi/toml"
)

type Manifest struct {
	App  ManifestApp  `toml:"application"`
	Deps ManifestDeps `toml:"dependencies"`
}

type ManifestApp struct {
	Name    string
	Version string
	Authors []string
}

type ManifestDeps map[string]string

func loadManifest() (*Manifest, error) {
	var m Manifest
	_, err := toml.DecodeFile(setting.ConfigFile, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// write writes the manifest data to disk
func (m Manifest) write() error {
	fp, err := os.Create(setting.ConfigFile)
	if err != nil {
		return fmt.Errorf("Error writing manifest: %s", err)
	}
	defer fp.Close()

	err = toml.NewEncoder(fp).Encode(m)
	if err != nil {
		return fmt.Errorf("Error writing manifest: %s", err)
	}

	return nil
}
