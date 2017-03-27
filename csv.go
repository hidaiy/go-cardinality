package dbindex

import (
	"errors"
	"fmt"
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

type Creator interface {
	Write(*Row) (int, error)
	WriteStringArray([]string) (int, error)
}

type CSV struct {
	out    io.Writer
	config Config
}

func NewCsv(out io.Writer, config Config) *CSV {
	return &CSV{
		out:    out,
		config: config,
	}
}

func (c CSV) Write(row *Row) (int, error) {
	return c.out.Write([]byte(c.CsvString(row)))
}
func (c CSV) WriteStringArray(array []string) (int, error) {
	return c.out.Write([]byte(strings.Join(array, ", ")))
}
func (r *CSV) CsvString(row *Row) string {
	return fmt.Sprintf("%s\n", strings.Join(row.StringArray(), ", "))
}

func (r *CSV) WriteDDL(columns []Column, tableRows map[string]int) (int, error) {
	var row *Row
	r.WriteStringArray(CSV_HEADER)
	for _, column := range columns {

		// テーブルのレコード件数
		rows, ok := tableRows[column.TableName]
		if !ok {
			return 0, errors.New(fmt.Sprintln("table count not found:", column.TableName))
		}

		// インデックスジェネレーターの作成
		indexGenerator, err := NewIndexGenerator(column, rows, r.config.Threshold)
		if err != nil {
			return 0, err
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
		r.Write(row)
	}
	return 0, nil
}
