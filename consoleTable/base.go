package consoleTable

import (
	"fmt"
	sutil "github.com/hidai620/go-cardinality/stringutil"
	"io"
	"strings"
)

type Base struct {
	out       io.Writer
	width     width
	template  string
	separator string
}

func (c *Base) writeSeparatorLine(left, middle, right string) (int, error) {
	lines := make([]string, 0, len(c.width))

	for _, w := range c.width {
		lines = append(lines, strings.Repeat("-", w))
	}

	line := left + strings.Join(lines, middle) + right + fmt.Sprintln()
	return c.out.Write([]byte(line))
}

// writeRow writes string array with template.
func (c *Base) writeRow(template string, row []string) (int, error) {
	line := fmt.Sprintf(template, sutil.ToInterfaces(row)...)
	return c.out.Write([]byte(line))
}
