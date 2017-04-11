package dbindex

import (
	"fmt"
	cnf "github.com/hidai620/go-cardinality/config"
	"github.com/jinzhu/gorm"
)

type Database struct {
	DB *gorm.DB
}

func Connect(config *cnf.Config) (*gorm.DB, error) {
	return gorm.Open(config.Dialect, createDBConnectString(config))
}

// 接続文字列を生成する。
func createDBConnectString(c *cnf.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.User, c.Password, c.Host, c.Port, c.Dialect)
}

// Params
type Params struct {
	values []interface{}
}

// NewParams returns Params pointer with values.
func NewParams(v interface{}) *Params {
	ret := &Params{}
	ret.values = append(ret.values, v)
	return ret
}

func (p *Params) Add(v interface{}) error {
	switch x := v.(type) {
	case string:
		if x != "" {
			p.values = append(p.values, x)
		}
	case []string:
		if x != nil {
			p.values = append(p.values, x)
		}
	case int:
		p.values = append(p.values, x)
	}
	return nil
}
