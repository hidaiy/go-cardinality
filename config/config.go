package config

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hidai620/go-cardinality/stringutil"
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
	Ignore    Ignore `toml:"ignore"`
}

const (
	configFileName  = "config.toml"
	ignoreAllColumn = "*"
)

// Load returns config loaded with argument file path.
func Load(path string) (*Config, error) {
	if path == "" {
		path = configFileName
	}

	var config = &Config{}
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}
	log.Println("ignore columns:", config.Ignore)

	return config, nil
}

// Ignore has exclude tables, or columns config.
type Ignore map[string]interface{}

//
func (i Ignore) HasConfig() bool {
	return len(i) != 0
}

func (i Ignore) IsIgnoreTable(table string) bool {
	value, ok := i[table]
	if ok {
		return value == ignoreAllColumn
	}
	return false
}

func (i Ignore) IsIgnoreColumn(table, column string) (bool, error) {
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
