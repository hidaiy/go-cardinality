package dbindex

import (
	"github.com/jinzhu/gorm"
)

type Table struct {
	DB           *gorm.DB
	DatabaseName string
	Name         string
	Rows         int
}
