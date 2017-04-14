package dbindex

import (
	"errors"
	"fmt"
	db "github.com/june-twenty/go-cardinality/lib/database"
	"io"
)

// RESULT_HEADER is used to out put table's header.
var RESULT_HEADER = []string{
	"table_name",
	"column_name",
	"table_rows",
	"distinct_rows",
	"cardinality",
	"indexes",
	"create_index_ddl",
	"drop_index_ddl",
}

// Writer
type Writer interface {
	WriteDDL(*db.SchemaInformation) error
}

// baseWriter
type baseWriter struct {
	out    io.Writer
	config *Config
}

// createRow is a function to create row from indexGenerator.
type createRow func(i *indexGenerator) []string

// body is table body rows
type body [][]string

// newBody returns table body instance.
func newBody(size int) body {
	return make([][]string, 0, size)
}

// getBody returns result table body.
func (c *baseWriter) createBody(columns []db.Column, tableRows db.TableRows, createRowFunction createRow) (body, error) {
	body := newBody(len(columns))

	for _, col := range columns {
		// Exclude column if it is specified as ignore column.
		if c.config.Ignore.HasConfig() {
			isIgnore, err := c.config.Ignore.IsIgnoreColumn(col.Table(), col.Column())
			if err != nil {
				return nil, err
			}
			if isIgnore {
				continue
			}
		}

		// Get table rows
		rows, ok := tableRows.GetRows(col.Table())
		if !ok {
			return nil, errors.New(fmt.Sprintf("table count not found: %s\n", col.Table()))
		}

		// Get new IndexGenerator.
		indexGenerator, err := newIndexGenerator(col, rows, c.config.Threshold)
		if err != nil {
			return nil, err
		}

		// Create row for result column list.
		row := createRowFunction(indexGenerator)

		body = append(body, row)
	}
	return body, nil
}
