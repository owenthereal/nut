package main

import "github.com/BurntSushi/toml"

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
	_, err := toml.DecodeFile(setting.ConfigFile, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
