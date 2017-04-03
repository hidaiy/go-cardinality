package option

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// アウトプットタイプ
type OutputType int

const (
	CONSOLE OutputType = iota
	CSV
)

func OutputTypeValueOf(o string) (OutputType, error) {
	switch strings.ToUpper(o) {
	case "CONSOLE":
		return CONSOLE, nil
	case "CSV":
		return CSV, nil
	default:
		return 0, errors.New(fmt.Sprintf("Pselese select type following list. console, csv"))
	}
}

func (o OutputType) Name() string {
	var ret string
	switch o {
	case CONSOLE:
		ret = "CONSOLE"
	case CSV:
		ret = "CSV"
	}
	return ret
}

// CommandLineOption has command line arguments.
type CommandLineOption struct {
	Out        OutputType
	ConfigPath string
	TableName  string
}

// New returns CommandLineOption created with arguments
func New(out OutputType, configPath, tableName string) *CommandLineOption {
	return &CommandLineOption{
		Out:        out,
		ConfigPath: configPath,
		TableName:  tableName,
	}
}

// ファイルパス
func Exists(f string) error {
	_, err := os.Stat(string(f))
	if err != nil {
		return errors.New(fmt.Sprintf("Specified file is not exists. %s", f))
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
	var config, out, tableName string
	flag.StringVar(&config, "config", "config.toml", "Absolute or relrative path of config file.")
	flag.StringVar(&out, "out", "console", `Output type of result. "console" or "csv"`)
	flag.StringVar(&tableName, "table", "", `Analyze Target table name.`)
	flag.Parse()

	err := Exists(config)
	if err != nil {
		return nil, err
	}

	outputType, err := OutputTypeValueOf(out)
	if err != nil {
		return nil, err
	}

	return New(outputType, config, tableName), nil
}
