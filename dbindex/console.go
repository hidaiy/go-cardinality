package dbindex

import (
	"errors"
	"fmt"
	"github.com/hidai620/go-cardinality/config"
	"github.com/hidai620/go-cardinality/consoleTable"
	iutil "github.com/hidai620/go-cardinality/intutil"
	sutil "github.com/hidai620/go-cardinality/stringutil"
	"io"
	_ "unicode/utf8"
)

var CONSOLE_HEADER = []string{
	"table_name",
	"column_name",
	"table_rows",
	"distinct_rows",
	"cardinality",
	"indexes",
	"create_index_ddl",
	"drop_index_ddl",
}

type Console struct {
	out    io.Writer
	config *config.Config
}

func NewConsole(out io.Writer, config *config.Config) *Console {
	return &Console{
		out:    out,
		config: config,
	}
}

// WriteDDL writes ddl.
func (c *Console) WriteDDL(columns []IColumn, tableRows TableRows) error {

	// get table body data
	body, err := c.getBody(columns, tableRows)
	if err != nil {
		return err
	}

	table := consoleTable.New(c.out)
	err = table.Write(CONSOLE_HEADER, body)
	if err != nil {
		return err
	}

	return nil
}

// body is table body rows
type body [][]string

// newBody returns table body instance.
func newBody(size int) body {
	return make([][]string, 0, size)
}

// getBody returns result table body.
func (c *Console) getBody(columns []IColumn, tableRows TableRows) (body, error) {
	body := newBody(len(columns))

	for _, col := range columns {
		// 対象外のカラムは処理から除外する
		if c.config.Ignore.HasConfig() {
			isIgnore, err := c.config.Ignore.IsIgnoreColumn(col.Table(), col.Column())
			if err != nil {
				return nil, err
			}
			if isIgnore {
				continue
			}
		}

		// テーブルのレコード件数
		rows, ok := tableRows.GetRows(col.Table())
		if !ok {
			return nil, errors.New(fmt.Sprintf("table count not found: %s\n", col.Table()))
		}

		// インデックスジェネレーターの作成
		indexGenerator, err := NewIndexGenerator(col, rows, c.config.Threshold)
		if err != nil {
			return nil, err
		}

		// body 1行分
		row := []string{
			col.Table(),
			col.Column(),
			iutil.ToString(rows),
			iutil.ToString(indexGenerator.DistinctTableRows),
			iutil.ToString(indexGenerator.GetColumnCardinality()),
			sutil.Cut(indexGenerator.ExistingIndexNames.CSV()),
			sutil.Cut(indexGenerator.GenerateCreateIndexDDL()),
			sutil.Cut(indexGenerator.GenerateDropIndexDDL()),
		}

		body = append(body, row)
	}
	return body, nil
}
