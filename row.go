package dbindex

import (
	"fmt"
	"strconv"
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
		r.toString(r.TableRows),
		r.toString(r.DistinctRows),
		r.toString(r.Cardinality),
		r.PrintExistingIndexNames(),
		r.CreateIndexDDL,
		r.DropIndexDDL,
	}
	return row
}

func (r *Row) PrintExistingIndexNames() string {
	return fmt.Sprintf("%q", strings.Join(r.Indexes, ", "))
}

func (r Row) toString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
