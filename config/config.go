package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	User      string
	Password  string
	Host      string
	Port      int
	Dialect   string
	Database  string
	Threshold int
}

const configFileName = "config.toml"

func Load(path string) (config *Config, err error) {
	if path == "" {
		path = configFileName
	}
	log.Println("path:", path)

	if _, err = toml.DecodeFile(path, &config); err != nil {
		return
	}
	return
}
