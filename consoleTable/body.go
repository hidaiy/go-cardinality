package consoleTable

import (
	"fmt"
	iutil "github.com/hidai620/go-mysql-study/intutil"
	sutil "github.com/hidai620/go-mysql-study/stringutil"
	"strings"
)

type Body struct {
	Base
}

// 値が数値の場合は右寄せ、文字列の場合左寄せにする。
func (c *Body) getColumnTemplate(s string, width int) string {
	var ret string

	if sutil.IsNumber(s) {
		ret = "%"
	} else {
		ret = "%-"
	}
	return ret + iutil.ToString(width) + "s"
}

// bodyのテンプレートを作成する。
func (c *Body) createRowTemplate(width []int, columns []string) string {
	parts := make([]string, 0, len(columns))
	for i := 0; i < len(width); i++ {
		parts = append(parts, c.getColumnTemplate(columns[i], width[i]))
	}
	template := c.separator + strings.Join(parts, c.separator) + c.separator + fmt.Sprintln()
	return template
}

// 出力する
func (c *Body) Write(width width, body [][]string) error {
	c.width = width
	c.template = c.createRowTemplate(width, body[0])

	for _, row := range body {
		_, err := c.writeRow(c.template, row)
		if err != nil {
			return err
		}
	}
	return nil
}
