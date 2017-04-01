package consoleTable

import (
	"fmt"
	iutil "github.com/hidai620/go-mysql-study/intutil"
	"strings"
)

type Header struct {
	Base
}

// writeTopLine writes header top line.
func (c *Header) writeTopLine() (int, error) {
	return c.writeSeparatorLine("┌", "┬", "┐")
}

// writeBottomLine writes header bottom line.
func (c *Header) writeBottomLine() (int, error) {
	return c.writeSeparatorLine("├", "┼", "┤")
}

func (c *Header) getColumnTemplate(width int) string {
	return "%-" + iutil.ToString(width) + "s"
}

func (c *Header) createRowTemplate(width []int, columns []string) string {
	parts := make([]string, 0, len(columns))
	for i := 0; i < len(width); i++ {
		parts = append(parts, c.getColumnTemplate(width[i]))
	}
	template := c.separator + strings.Join(parts, c.separator) + c.separator + fmt.Sprintln()
	return template
}

func (c *Header) Write(width width, header []string) error {
	// 出力テンプレートの作成
	c.width = width
	c.template = c.createRowTemplate(width, header)

	_, err := c.writeTopLine()
	if err != nil {
		return err
	}

	_, err = c.writeRow(c.template, header)
	if err != nil {
		return err
	}
	_, err = c.writeBottomLine()
	if err != nil {
		return err
	}
	return nil
}
