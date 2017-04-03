package config

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hidai620/go-mysql-study/stringutil"
)

type Config struct {
	User          string
	Password      string
	Host          string
	Port          int
	Dialect       string
	Database      string
	Threshold     int
	IgnoreColumns IgnoreColumns `toml:"ignore"`
}

const configFileName = "config.toml"

func Load(path string) (*Config, error) {
	if path == "" {
		path = configFileName
	}

	var config = &Config{}
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}
	fmt.Println("ignore columns:", config.IgnoreColumns)

	return config, nil
}

type IgnoreColumns map[string]interface{}

func (c Config) HasIgnoreConfig() bool {
	return len(c.IgnoreColumns) != 0
}
func (c Config) IsIgnoreColumn(table, column string) (bool, error) {
	return c.IgnoreColumns.Contains(table, column)
}

func (i IgnoreColumns) HasConfig() bool {
	return len(i) != 0
}

func (i IgnoreColumns) Contains(table, column string) (bool, error) {
	value, ok := i[table]
	if !ok {
		return false, nil
	}

	columns, err := stringutil.ToStrings(value.([]interface{}))
	if err != nil {
		return false, errors.New(fmt.Sprintf("Ignore config is not valid. %#v", value))
	}

	return stringutil.Contains(columns, column), nil
}
