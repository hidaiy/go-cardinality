package option

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// アウトプットタイプ
type Out int

const (
	CONSOLE Out = iota
	CSV
)

func OutputTypeValueOf(o string) (Out, error) {
	switch strings.ToUpper(o) {
	case "CONSOLE":
		return CONSOLE, nil
	case "CSV":
		return CSV, nil
	default:
		return 0, errors.New(fmt.Sprintf("Pselese select type following list. console, csv"))
	}
}

func (o Out) Name() string {
	var ret string
	switch o {
	case CONSOLE:
		ret = "CONSOLE"
	case CSV:
		ret = "CSV"
	}
	return ret
}

// コマンドライン引数
type CommandLineOption struct {
	Out        Out
	ConfigPath string
}

func New(out Out, configPath string) *CommandLineOption {
	return &CommandLineOption{
		Out:        out,
		ConfigPath: configPath,
	}
}

// ファイルパス
func Exists(f string) error {
	_, err := os.Stat(string(f))
	if err != nil {
		return errors.New("config is not exists.")
	}
	return nil
}

// 同じ値を持つ場合、trueを返す
func (c *CommandLineOption) Equals(c2 *CommandLineOption) bool {
	return c.Out == c2.Out &&
		c.ConfigPath == c2.ConfigPath
}

// コマンドラインオプションをパースし、CommandLineOptionにして返す
func Parse() (*CommandLineOption, error) {
	var config, out string
	flag.StringVar(&config, "config", "", "コンフィルファイルのパス")
	flag.StringVar(&out, "out", "console", "出力方法")
	flag.Parse()

	err := Exists(config)
	if err != nil {
		return nil, err
	}

	outputType, err := OutputTypeValueOf(out)
	if err != nil {
		return nil, err
	}

	return New(outputType, config), nil
}
