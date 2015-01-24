package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	App  App  `toml:"application"`
	Deps Deps `toml:"dependencies"`
}

type App struct {
	Name    string
	Version string
	Authors []string
}

type Deps map[string]string

func loadConfig() (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(configFile(), &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func configFile() string {
	dir, err := os.Getwd()
	check(err)

	return filepath.Join(dir, "Nut.toml")
}
