package dbindex

import (
	"errors"
	"fmt"
	"github.com/hidai620/go-mysql-study/config"
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

func (r *CSV) WriteDDL(columns []Column, tableRows TableRows) error {
	var row *Row
	r.writeRow(CSV_HEADER)

	for _, column := range columns {
		// 対象外のカラムは処理から除外する
		if r.config.HasIgnoreConfig() {
			isIgnore, err := r.config.IsIgnoreColumn(column.TableName, column.ColumnName)
			if err != nil {
				return err
			}
			if isIgnore {
				continue
			}
		}

		// テーブルのレコード件数の取得
		rows, ok := tableRows.GetRows(column.TableName)
		if !ok {
			return errors.New(fmt.Sprintln("table count not found:", column.TableName))
		}

		// インデックスジェネレーターの作成
		indexGenerator, err := NewIndexGenerator(column, rows, r.config.Threshold)
		if err != nil {
			return err
		}

		// １行分
		row = &Row{
			TableName:      column.TableName,
			ColumnName:     column.ColumnName,
			TableRows:      rows,
			DistinctRows:   indexGenerator.DistinctTableRows,
			Cardinality:    indexGenerator.GetColumnCardinality(),
			Indexes:        indexGenerator.ExistingIndexNames,
			CreateIndexDDL: indexGenerator.GenerateCreateIndexDDL(),
			DropIndexDDL:   indexGenerator.GenerateDropIndexDDL(),
		}

		// 出力
		r.writeRow(row.StringArray())
	}
	return nil
}
