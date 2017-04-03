package dbindex

import (
	"errors"
	"fmt"
	"github.com/hidai620/go-mysql-study/config"
	"github.com/hidai620/go-mysql-study/consoleTable"
	iutil "github.com/hidai620/go-mysql-study/intutil"
	sutil "github.com/hidai620/go-mysql-study/stringutil"
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
func (c *Console) WriteDDL(columns []Column, tableRows TableRows) error {

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

func newBody(size int) body {
	return make([][]string, 0, size)
}

func (c *Console) getBody(columns []Column, tableRows TableRows) (body, error) {
	body := newBody(len(columns))

	for _, column := range columns {
		// 対象外のカラムは処理から除外する
		if c.config.HasIgnoreConfig() {
			isIgnore, err := c.config.IsIgnoreColumn(column.TableName, column.ColumnName)
			if err != nil {
				return nil, err
			}
			if isIgnore {
				continue
			}
		}

		// テーブルのレコード件数
		rows, ok := tableRows.GetRows(column.TableName)
		if !ok {
			return nil, errors.New(fmt.Sprintln("table count not found:", column.TableName))
		}

		// インデックスジェネレーターの作成
		indexGenerator, err := NewIndexGenerator(column, rows, c.config.Threshold)
		if err != nil {
			return nil, err
		}

		// body 1行分
		row := []string{
			column.TableName,
			column.ColumnName,
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
