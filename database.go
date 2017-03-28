package dbindex

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Database struct {
	config Config
}

func NewDatabase(config Config) *Database {
	return &Database{config: config}
}

func (d Database) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(d.config.Dialect, createDBConnectString(d.config))
	if err != nil {
		printError(err)
		return db, err
	}
	return db, err
}

// 接続文字列を生成する。
func createDBConnectString(c Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.User, c.Password, c.Host, c.Port, "mysql")
}
