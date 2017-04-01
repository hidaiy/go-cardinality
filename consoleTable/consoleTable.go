package consoleTable

import (
	"fmt"
	iutil "github.com/hidai620/go-mysql-study/intutil"
	sutil "github.com/hidai620/go-mysql-study/stringutil"
	"io"
	"strings"
)

type width []int

type ConsoleTable struct {
	out            io.Writer
	width          width
	headerTemplate string
	bodyTemplate   string
	separator      string
	padding        int
}

func New(out io.Writer) *ConsoleTable {
	return &ConsoleTable{
		out:       out,
		separator: "|",
		padding:   1,
	}
}

func (c *ConsoleTable) writeHeaderRow(row []string) (int, error) {
	return c.writeRow(c.headerTemplate, row)
}
func (c *ConsoleTable) writeHeaderTopLine() (int, error) {
	return c.writeSeparatorLine("┌", "┬", "┐")
}

func (c *ConsoleTable) writeHeaderBottomLine() (int, error) {
	return c.writeSeparatorLine("├", "┼", "┤")
}

func (c *ConsoleTable) writeBodyRow(row []string) (int, error) {
	return c.writeRow(c.bodyTemplate, row)
}

// footer
func (c *ConsoleTable) writeFooterLine() (int, error) {
	return c.writeSeparatorLine("└", "┴", "┘")
}

//
func (c *ConsoleTable) writeRow(template string, row []string) (int, error) {
	line := fmt.Sprintf(template, sutil.ToInterfaces(row)...)
	return c.out.Write([]byte(line))
}

func (c *ConsoleTable) writeSeparatorLine(left, middle, right string) (int, error) {
	lines := make([]string, 0, len(c.width))

	for _, w := range c.width {
		lines = append(lines, strings.Repeat("-", w))
	}

	line := left + strings.Join(lines, middle) + right + fmt.Sprintln()
	return c.out.Write([]byte(line))
}

// 値が数値の場合は右寄せ、文字列の場合左寄せにする。
func getBodyColumnTemplate(s string, width int) string {
	var ret string

	if sutil.IsNumber(s) {
		ret = "%"
	} else {
		ret = "%-"
	}
	return ret + iutil.ToString(width) + "s"
}

// ヘッダーカラム用のテンプレートを返す
func getHeaderColumnTemplate(width int) string {
	return "%-" + iutil.ToString(width) + "s"
}

// 1行分のテンプレートを作る
func (c *ConsoleTable) createRowTemplate(header []string, columnsSlice [][]string) string {

	width := c.getWidth(header, columnsSlice)

	// ヘッダー用テンプレート
	c.headerTemplate = c.createHeaderTemplate(width, header)

	// データ行用テンプレート
	c.bodyTemplate = c.createBodyTemplate(width, columnsSlice[1])
	return c.bodyTemplate
}

// ヘッダーのテンプレートを返す
func (c *ConsoleTable) createHeaderTemplate(width []int, columns []string) string {
	parts := make([]string, 0, len(columns))
	for i := 0; i < len(width); i++ {
		parts = append(parts, getHeaderColumnTemplate(width[i]))
	}
	template := c.separator + strings.Join(parts, c.separator) + c.separator + fmt.Sprintln()
	return template
}

// bodyのテンプレートを作成する。
func (c *ConsoleTable) createBodyTemplate(width []int, columns []string) string {
	parts := make([]string, 0, len(columns))
	for i := 0; i < len(width); i++ {
		parts = append(parts, getBodyColumnTemplate(columns[i], width[i]))
	}
	template := c.separator + strings.Join(parts, c.separator) + c.separator + fmt.Sprintln()
	return template
}

// 列ごとの文字数を計算する
func (c *ConsoleTable) getWidth(header []string, columnsList [][]string) width {
	ret := make([]int, len(header), len(header))

	for i := 0; i < len(header); i++ {
		ret[i] = iutil.Max(ret[i], len(header[i]))
	}

	for _, columns := range columnsList {
		for i := 0; i < len(columns); i++ {
			ret[i] = iutil.Max(ret[i], len(columns[i]))
		}
	}

	// 各列の幅にパディングをプラスする
	ret = iutil.Map(ret, c.appendPadding)
	c.width = ret

	return ret
}

func (c *ConsoleTable) appendPadding(i int) int {
	return i + c.padding
}

//
func (c *ConsoleTable) writeHeader(header []string) error {
	_, err := c.writeHeaderTopLine()
	if err != nil {
		return err
	}

	_, err = c.writeHeaderRow(header)
	if err != nil {
		return err
	}
	_, err = c.writeHeaderBottomLine()
	if err != nil {
		return err
	}
	return nil
}

func (c *ConsoleTable) writeBody(rowSlice [][]string) error {
	for _, row := range rowSlice {
		_, err := c.writeBodyRow(row)
		if err != nil {
			return err
		}
	}
	return nil
}

// Write writes table.
func (c *ConsoleTable) Write(header []string, body [][]string) error {
	// テンプレートの作成
	c.createRowTemplate(header, body)

	// ヘッダー
	err := c.writeHeader(header)
	if err != nil {
		return err
	}

	// ボディー
	err = c.writeBody(body)
	if err != nil {
		return err
	}

	// フッター
	c.writeFooterLine()
	if err != nil {
		return err
	}
	return nil
}
