package main

import "github.com/gophergala/nut/vendor/_nuts/github.com/BurntSushi/toml"

type Config struct {
	App  ConfigApp  `toml:"application"`
	Deps ConfigDeps `toml:"dependencies"`
}

type ConfigApp struct {
	Name    string
	Version string
	Authors []string
}

type ConfigDeps map[string]string

func loadConfig() (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(setting.ConfigFile, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
