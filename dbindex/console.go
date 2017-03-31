package dbindex

import (
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
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
	out            io.Writer
	threshold      int
	columnTemplate string
	width          width
}

func NewConsole(out io.Writer, threshold int) *Console {
	return &Console{
		out:       out,
		threshold: threshold,
	}
}

func (c *Console) WriteRow(row []string) (int, error) {
	line := fmt.Sprintf(c.columnTemplate, toInterfaces(row)...)
	return c.out.Write([]byte(line))
}

func (c *Console) WriteHeaderTopLine() (int, error) {
	return c.writeSeparatorLine("┌", "┬", "┐")
}

func (c *Console) WriteHeaderBottomLine() (int, error) {
	return c.writeSeparatorLine("├", "┼", "┤")
}

func (c *Console) WriteFooterLine() (int, error) {
	return c.writeSeparatorLine("└", "┴", "┘")
}

func (c *Console) writeSeparatorLine(left, separator, right string) (int, error) {
	tmp := make([]string, 0, len(c.width))

	for _, w := range c.width {
		tmp = append(tmp, strings.Repeat("-", w))
	}

	line := left + strings.Join(tmp, separator) + right + fmt.Sprintln()
	return c.out.Write([]byte(line))
}

func getColumnTemplate(s string, width int) string {
	var ret string

	if isNumber(s) {
		ret = "%"
	} else {
		ret = "%-"
	}
	ret = ret + toString(width) + "s"
	return ret
}

// 1行分のテンプレートを作る
func (c *Console) createTemplate(columnsSlice [][]string) string {

	width := c.getWidth(columnsSlice)

	// データ行用テンプレート
	c.columnTemplate = c.createColumnTemplate(width, columnsSlice[1])
	return c.columnTemplate
}

func (c *Console) createColumnTemplate(width []int, columns []string) string {
	parts := make([]string, 0, len(columns))
	for i := 0; i < len(width); i++ {
		parts = append(parts, getColumnTemplate(columns[i], width[i]))
	}
	template := separator + strings.Join(parts, separator) + separator + fmt.Sprintln()
	return template
}

var maxLength = 12
var separator = "|"
var padding = 1
var numberRegex = regexp.MustCompile(`[0-9]`)

func cut(s string, length int) string {
	if len(s) > maxLength {
		return fmt.Sprintf("%s...", s[:length])
	}
	return s
}

func isNumber(s string) bool {
	return numberRegex.MatchString(s)
}

func toInterfaces(array []string) []interface{} {
	ret := make([]interface{}, len(array), len(array))
	for i := 0; i < len(array); i++ {
		ret[i] = array[i]
	}
	return ret
}

type width []int

// int
func intsMap(ints []int, fn func(int) int) []int {
	ret := make([]int, 0, len(ints))
	for _, i := range ints {
		ret = append(ret, fn(i))
	}
	return ret
}

func max(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

func toString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

// 列ごとの文字数を計算する
func (c *Console) getWidth(columnSlice [][]string) width {
	ret := make([]int, len(columnSlice[0]), len(columnSlice[0]))

	for _, columns := range columnSlice {
		for i := 0; i < len(columns); i++ {
			column := columns[i]
			ret[i] = max(ret[i], len(column))
		}
	}

	// 各列の幅にパディングをプラスする
	ret = intsMap(ret, appendPadding)
	c.width = ret

	return ret
}

func appendPadding(i int) int {
	return i + padding
}

func (c *Console) WriteDDL(columns []Column, tableRows TableRows) (int, error) {
	rowSlice := make([][]string, 0, len(columns))

	rowSlice = append(rowSlice, CONSOLE_HEADER)
	for _, column := range columns {

		// テーブルのレコード件数
		rows, ok := tableRows.GetRows(column.TableName)
		if !ok {
			return 0, errors.New(fmt.Sprintln("table count not found:", column.TableName))
		}

		// インデックスジェネレーターの作成
		indexGenerator, err := NewIndexGenerator(column, rows, c.threshold)
		if err != nil {
			return 0, err
		}

		// １行分
		row := []string{
			column.TableName,
			column.ColumnName,
			toString(rows),
			toString(indexGenerator.DistinctTableRows),
			toString(indexGenerator.GetColumnCardinality()),
			cut(strings.Join(indexGenerator.ExistingIndexNames, ","), maxLength),
			cut(indexGenerator.GenerateCreateIndexDDL(), maxLength),
			cut(indexGenerator.GenerateDropIndexDDL(), maxLength),
		}

		rowSlice = append(rowSlice, row)
	}

	// テンプレートの作成
	c.createTemplate(rowSlice)

	// 出力

	// ヘッダー上
	c.WriteHeaderTopLine()

	for i, row := range rowSlice {
		if i == 1 {
			// ヘッダー下
			c.WriteHeaderBottomLine()
		}
		_, err := c.WriteRow(row)
		if err != nil {
			return 0, err
		}
	}
	// フッター行
	c.WriteFooterLine()

	return 0, nil
}
