package dbindex

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hidaiy/go-utils/stringutil"
	"log"
)

// Config
type Config struct {
	User      string
	Password  string
	Host      string
	Port      int
	Dialect   string
	Database  string
	Threshold int
	Ignore    ignore `toml:"ignore"`
}

const (
	configFileName  = "config.toml"
	ignoreAllColumn = "*"
)

// Load returns config loaded with argument file path.
func LoadConfig(path string) (*Config, error) {
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
type ignore map[string]interface{}

// HasConfig returns true if config.toml has ignore tables config.
func (i ignore) HasConfig() bool {
	return len(i) != 0
}

// IsIgnoreTable returns true if table name is specified in config file as Ignore table,
// and has "*"  as column name.
func (i ignore) IsIgnoreTable(table string) bool {
	value, ok := i[table]
	if ok {
		return value == ignoreAllColumn
	}
	return false
}

// IsIgnoreColumn returns true, if table and column are specified in config file as Ignore Column.
func (i ignore) IsIgnoreColumn(table, column string) (bool, error) {
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
