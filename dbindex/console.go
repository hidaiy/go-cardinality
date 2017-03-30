package dbindex

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
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
	template  string
	width     width
}

func NewConsole(out io.Writer, threshold int) *Console {
	return &Console{
		out:       out,
		threshold: threshold,
	}
}

func (c *Console) WriteRow(row *Row) (int, error) {
	return c.out.Write([]byte(c.String(row)))
}
func (c *Console) writeStringArray(array []string) (int, error) {

	tmp := make([]interface{}, len(array), len(array))
	for i := 0; i < len(array); i++ {
		tmp[i] = array[i]
	}

	line := fmt.Sprintf(c.template, tmp...)
	return c.out.Write([]byte(line))
}

func (c *Console) createRowTemplate(width width) string {
	template := "%-" + ToString(width.tableName) + "s " +
		"%-" + ToString(width.columnName) + "s " +
		"%" + ToString(width.tableRows) + "s " +
		"%" + ToString(width.distinctRows) + "s " +
		"%" + ToString(width.cardinality) + "s " +
		"%-" + ToString(width.indexes) + "s " +
		"%-" + ToString(width.createIndexDDL) + "s " +
		"%-" + ToString(width.dropIndexDDL) + "s \n"
	c.template = template
	return template
}

func (r *Console) String(row *Row) string {
	//fmt.Println(r.template)
	return fmt.Sprintf(r.template,
		row.TableName,
		row.ColumnName,
		ToString(row.TableRows),
		ToString(row.DistinctRows),
		ToString(row.Cardinality),
		cut(strings.Join(row.Indexes, ","), 8),
		cut(row.CreateIndexDDL, 8),
		cut(row.DropIndexDDL, 8),
	)
}

var maxLength = 8

func cut(s string, length int) string {
	if len(s) > maxLength {
		return fmt.Sprintf("%s...", s[:length])
	}
	return s
}

type width struct {
	tableName      int
	columnName     int
	tableRows      int
	distinctRows   int
	cardinality    int
	indexes        int
	createIndexDDL int
	dropIndexDDL   int
}

type widthLength int

func ToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func length(i int) float64 {
	return float64(len(strconv.FormatInt(int64(i), 10)))
}

func (r *Console) WriteDDL(columns []Column, tableRows TableRows) (int, error) {
	var row *Row

	width := width{}
	for i := 0; i < len(columns); i++ {
		column := columns[i]
		width.tableName = int(math.Max(float64(width.tableName), float64(len(column.TableName))))
		width.columnName = int(math.Max(float64(width.columnName), float64(len(column.ColumnName))))

		//distinctRows, err := column.DistinctRows()
		//if err != nil {
		//	return 0, err
		//}
		//
		//width.distinctRows = int(math.Max(length(width.distinctRows), length(distinctRows)))
	}
	width.tableRows = 7
	width.distinctRows = 7
	width.cardinality = 5
	width.indexes = 10
	width.createIndexDDL = 10
	width.dropIndexDDL = 10
	r.createRowTemplate(width)
	//fmt.Printf("%#v", width)

	r.writeStringArray(CONSOLE_HEADER)
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
