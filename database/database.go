package database

import (
	"fmt"
	cnf "github.com/hidai620/go-cardinality/config"
	"github.com/jinzhu/gorm"
)
type Connection struct {

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
	Values []interface{}
}

// NewParams returns Params pointer with values.
func NewParams(v interface{}) *Params {
	ret := &Params{}
	ret.Values = append(ret.Values, v)
	return ret
}

func (p *Params) Add(v interface{}) error {
	switch x := v.(type) {
	case string:
		if x != "" {
			p.Values = append(p.Values, x)
		}
	case []string:
		if x != nil {
			p.Values = append(p.Values, x)
		}
	case int:
		p.Values = append(p.Values, x)
	}
	return nil
}

type IColumn interface {
	Column() string
	Table() string
	IndexNames() ([]string, error)
	DistinctRows() (int, error)
}

type TableRows map[string]int

// GetRows returns rows searched with given table name.
func (t TableRows) GetRows(tableName string) (int, bool) {
	rows, ok := t[tableName]
	return rows, ok
}

type SchemaInformation struct {
	TableRows TableRows
	Columns   []IColumn
}

func NewSchemaInformation(tableRows TableRows, columns []IColumn) *SchemaInformation {
	return &SchemaInformation{
		TableRows: tableRows,
		Columns:   columns,
	}
}

type Database interface {
	GetSchemaInformation(string, []string) (*SchemaInformation)
}
