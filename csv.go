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

type CSV struct {
	out       io.Writer
	threshold int
}

func NewCSV(out io.Writer, threshold int) *CSV {
	return &CSV{
		out:       out,
		threshold: threshold,
	}
}

func (c CSV) WriteRow(row *Row) (int, error) {
	return c.out.Write([]byte(c.csvString(row)))
}
func (c CSV) WriteStringArray(array []string) (int, error) {
	return c.out.Write([]byte(strings.Join(array, ", ")))
}
func (r *CSV) csvString(row *Row) string {
	return fmt.Sprintf("%s\n", strings.Join(row.StringArray(), ", "))
}

func (r *CSV) WriteDDL(columns []Column, tableRows TableRows) (int, error) {
	var row *Row
	r.WriteStringArray(CSV_HEADER)
	for _, column := range columns {

		// テーブルのレコード件数
		rows, ok := tableRows.GetRows(column.TableName)
		if !ok {
			return 0, errors.New(fmt.Sprintln("table count not found:", column.TableName))
		}

		// インデックスジェネレーターの作成
		indexGenerator, err := NewIndexGenerator(column, rows, r.threshold)
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
		r.WriteRow(row)
	}
	return 0, nil
}
