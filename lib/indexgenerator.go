package dbindex

import (
	"fmt"
	db "github.com/hidaiy/go-cardinality/lib/database"
	"strings"
)

// indexGenerator
type indexGenerator struct {
	Column             db.Column
	TableName          string
	ColumnName         string
	TableRows          int
	DistinctTableRows  int
	Threshold          int
	IndexName          string
	ExistingIndexNames stringArray
}

// stringArray is alias of string array type.
type stringArray []string

// CSV returns stringArray as csv string.
func (s stringArray) CSV() string {
	return strings.Join(s, ",")
}

// newIndexGenerator is constructor.
func newIndexGenerator(column db.Column, tableRows, threshold int) (*indexGenerator, error) {
	// Getting distinct rows of column.
	distinctTableRows, err := column.DistinctRows()
	if err != nil {
		return nil, err
	}

	// Getting index names as string array.
	indexNames, err := column.IndexNames()
	if err != nil {
		return nil, err
	}

	ret := &indexGenerator{
		Column:             column,
		TableName:          column.Table(),
		ColumnName:         column.Column(),
		TableRows:          tableRows,
		DistinctTableRows:  distinctTableRows,
		Threshold:          threshold,
		IndexName:          indexName(column.Table(), column.Column()),
		ExistingIndexNames: indexNames,
	}
	return ret, nil
}

// indexName returns index name created from table and column names.
func indexName(tableName, columnName string) string {
	return fmt.Sprintf("i_%s__%s", tableName, columnName)
}

// GenerateCreateIndexDDL returns create index ddl, created from table and column names.
// If cardinality of column is under threshold, returns empty string.
func (i *indexGenerator) GenerateCreateIndexDDL() string {
	if i.NeedToCreateIndex() {
		return fmt.Sprintf("alter table %s add index %s(%s);", i.TableName, i.IndexName, i.ColumnName)
	} else {
		return "[*1]"
	}
}

// GenerateDropIndexDDL returns drop index ddl created from table and column names.
// If cardinality of column is under threshold, returns empty string.
func (i *indexGenerator) GenerateDropIndexDDL() string {
	if i.NeedToCreateIndex() {
		return fmt.Sprintf("alter table %s drop index %s;", i.TableName, i.IndexName)
	} else {
		return ""
	}
}

// GetColumnCardinality returns cardinality.
func (i *indexGenerator) GetColumnCardinality() int {
	if i.DistinctTableRows == 0 && i.TableRows == 0 {
		return 0
	}

	tmp := float64(i.DistinctTableRows) / float64(i.TableRows)
	return int(tmp * float64(100))
}

// NeedToCreateIndex returns true if column's cardinality is higher threshold,
// and columns does not have indexes.
func (i *indexGenerator) NeedToCreateIndex() bool {
	return i.GetColumnCardinality() >= i.Threshold && len(i.ExistingIndexNames) == 0
}
