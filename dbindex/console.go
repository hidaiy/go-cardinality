package dbindex

import (
	"errors"
	"fmt"
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
	out       io.Writer
	threshold int
}

func NewConsole(out io.Writer, threshold int) *Console {
	return &Console{
		out:       out,
		threshold: threshold,
	}
}

// WriteDDL writes ddl.
func (c *Console) WriteDDL(columns []Column, tableRows TableRows) error {

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

const maxLength = 12

// body is table body rows
type body [][]string

func newBody(size int) body {
	return make([][]string, 0, size)
}

func (c *Console) getBody(columns []Column, tableRows TableRows) (body, error) {
	body := newBody(len(columns))

	for _, column := range columns {
		// :TODO 対象外からむの判定の追加

		// テーブルのレコード件数
		rows, ok := tableRows.GetRows(column.TableName)
		if !ok {
			return nil, errors.New(fmt.Sprintln("table count not found:", column.TableName))
		}

		// インデックスジェネレーターの作成
		indexGenerator, err := NewIndexGenerator(column, rows, c.threshold)
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
			sutil.Cut(indexGenerator.ExistingIndexNames.CSV(), maxLength),
			sutil.Cut(indexGenerator.GenerateCreateIndexDDL(), maxLength),
			sutil.Cut(indexGenerator.GenerateDropIndexDDL(), maxLength),
		}

		body = append(body, row)
	}
	return body, nil
}
