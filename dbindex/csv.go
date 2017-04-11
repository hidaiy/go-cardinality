package dbindex

import (
	"errors"
	"fmt"
	"github.com/hidai620/go-cardinality/config"
	"io"
	"strings"
)

var CSV_HEADER = []string{
	"table_name",
	"column_name",
	"table_rows",
	"distinct_rows",
	"cardinality",
	"indexes",
	"create_index_ddl",
	"drop_index_ddl",
}

type CSV struct {
	out    io.Writer
	config *config.Config
}

func NewCSV(out io.Writer, config *config.Config) *CSV {
	return &CSV{
		out:    out,
		config: config,
	}
}

func (c CSV) writeRow(array []string) (int, error) {
	return c.write(strings.Join(array, ", "))
}

func (c CSV) write(s string) (int, error) {
	return c.out.Write([]byte(s + fmt.Sprintln()))
}

func (c *CSV) WriteDDL(columns []IColumn, tableRows TableRows) error {
	var row *Row
	c.writeRow(CSV_HEADER)

	for _, col := range columns {
		// 対象外のカラムは処理から除外する
		if c.config.Ignore.HasConfig() {
			isIgnore, err := c.config.Ignore.IsIgnoreColumn(col.Table(), col.Column())
			if err != nil {
				return err
			}
			if isIgnore {
				continue
			}
		}

		// テーブルのレコード件数の取得
		rows, ok := tableRows.GetRows(col.Table())
		if !ok {
			return errors.New(fmt.Sprintln("table count not found:", col.Table()))
		}

		// インデックスジェネレーターの作成
		indexGenerator, err := NewIndexGenerator(col, rows, c.config.Threshold)
		if err != nil {
			return err
		}

		// １行分
		row = &Row{
			TableName:      col.Table(),
			ColumnName:     col.Column(),
			TableRows:      rows,
			DistinctRows:   indexGenerator.DistinctTableRows,
			Cardinality:    indexGenerator.GetColumnCardinality(),
			Indexes:        indexGenerator.ExistingIndexNames,
			CreateIndexDDL: indexGenerator.GenerateCreateIndexDDL(),
			DropIndexDDL:   indexGenerator.GenerateDropIndexDDL(),
		}

		// 出力
		c.writeRow(row.StringArray())
	}
	return nil
}
