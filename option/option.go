package option

import (
	"errors"
	"flag"
	"fmt"
	sutil "github.com/hidai620/go-cardinality/stringutil"
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
func (o OutputType) String() string {
	return o.Name()
}

// CommandLineOption has command line arguments.
type CommandLineOption struct {
	Out        OutputType
	ConfigPath string
	TableNames []string
}

// New returns CommandLineOption created with arguments
func New(out OutputType, configPath string, tableNames []string) *CommandLineOption {
	return &CommandLineOption{
		Out:        out,
		ConfigPath: configPath,
		TableNames: tableNames,
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
	var config, out, tableNames string
	flag.StringVar(&config, "config", "config.toml", "Absolute or relrative path of config file.")
	flag.StringVar(&out, "out", "console", `Output type of result. "console" or "csv"`)
	flag.StringVar(&tableNames, "table", "", `Analyze Target table name.`)
	flag.Parse()

	err := Exists(config)
	if err != nil {
		return nil, err
	}

	outputType, err := OutputTypeValueOf(out)
	if err != nil {
		return nil, err
	}

	tableNamesArray := sutil.Split(tableNames, ",")
	return New(outputType, config, tableNamesArray), nil
}
