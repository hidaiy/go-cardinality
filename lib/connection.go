package dbindex

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

// ConnectDatabase connects database, and returns database connection.
func ConnectDatabase(config *Config) (*gorm.DB, error) {
	return gorm.Open(config.Dialect, createDBConnectString(config))
}

// createDBConnectString returns string to connect database by database driver.
func createDBConnectString(c *Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.User, c.Password, c.Host, c.Port, c.Dialect)
}
