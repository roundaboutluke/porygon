package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

func (config *Config) ParseConfig() error {
	if _, err := toml.DecodeFile("default.toml", &config); err != nil {
		log.Println("error decoding default config file,", err)
		return err
	}

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Println("error decoding user config file,", err)
		return err
	}

	return nil
}
