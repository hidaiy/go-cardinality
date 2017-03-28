package dbindex

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Database struct {
	DB *gorm.DB
}

func Connect(config Config) (*gorm.DB, error) {
	return gorm.Open(config.Dialect, createDBConnectString(config))
}

// 接続文字列を生成する。
func createDBConnectString(c Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.User, c.Password, c.Host, c.Port, "mysql")
}
