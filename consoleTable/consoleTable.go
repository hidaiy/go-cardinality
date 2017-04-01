package consoleTable

import (
	iutil "github.com/hidai620/go-mysql-study/intutil"
	"io"
)

type width []int

// ConsoleTable
type ConsoleTable struct {
	out     io.Writer
	Padding int
	Header
	Body
	Footer
}

func New(out io.Writer) *ConsoleTable {
	separator := "|"
	padding := 1

	base := Base{
		out:       out,
		separator: separator,
	}
	return &ConsoleTable{
		out:     out,
		Padding: padding,
		Header: Header{
			base,
		},
		Body: Body{
			base,
		},
		Footer: Footer{
			base,
		},
	}
}

// 列ごとの文字数を計算する
func (c *ConsoleTable) getWidth(header []string, body [][]string) width {
	ret := make([]int, len(header), len(header))

	for i := 0; i < len(header); i++ {
		ret[i] = len(header[i])
	}

	for _, columns := range body {
		for i := 0; i < len(columns); i++ {
			ret[i] = iutil.Max(ret[i], len(columns[i]))
		}
	}

	// 各列の幅にパディングをプラスする
	ret = iutil.Map(ret, c.appendPadding)
	return ret
}

// appendPadding adds padding to argument.
func (c *ConsoleTable) appendPadding(i int) int {
	return i + c.Padding
}

// Write writes table.
func (c *ConsoleTable) Write(header []string, body [][]string) error {
	// 幅の取得
	width := c.getWidth(header, body)

	// ヘッダー
	err := c.Header.Write(width, header)
	if err != nil {
		return err
	}

	// ボディー
	err = c.Body.Write(width, body)
	if err != nil {
		return err
	}

	// フッター
	c.Footer.Write(width)
	if err != nil {
		return err
	}
	return nil
}
