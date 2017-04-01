package dbindex

import (
	"fmt"
	"github.com/hidai620/go-mysql-study/intutil"
	"strings"
)

type Row struct {
	TableName      string
	ColumnName     string
	TableRows      int
	DistinctRows   int
	Cardinality    int
	Indexes        []string
	CreateIndexDDL string
	DropIndexDDL   string
}

func (r *Row) StringArray() []string {
	var row = []string{
		r.TableName,
		r.ColumnName,
		intutil.ToString(r.TableRows),
		intutil.ToString(r.DistinctRows),
		intutil.ToString(r.Cardinality),
		r.PrintExistingIndexNames(),
		r.CreateIndexDDL,
		r.DropIndexDDL,
	}
	return row
}

func (r *Row) PrintExistingIndexNames() string {
	return fmt.Sprintf("%q", strings.Join(r.Indexes, ", "))
}
