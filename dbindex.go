package main

import (
	"fmt"
	"github.com/hidai620/go-mysql-study/config"
	. "github.com/hidai620/go-mysql-study/dbindex"
	"github.com/hidai620/go-mysql-study/option"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

func main() {
	// コマンドラインオプションのパース
	opt, err := option.Parse()
	if err != nil {
		printError(err)
		return
	}

	//　設定ファイルの読み込み
	conf, err := config.Load(opt.ConfigPath)
	if err != nil {
		printError(err)
		return
	}

	// DB接続
	db, err := Connect(conf)
	if err != nil {
		printError(err)
		return
	}
	defer db.Close()

	// 管理スキーマの取得
	informationSchema := NewInformationSchema(db)

	// テーブル単位の件数の取得
	tableRows, err := informationSchema.TableRows(conf.Database, opt.TableName)
	if err != nil {
		printError(err)
		return
	}

	if len(tableRows) != 0 {
		// カラムの取得
		columns, err := informationSchema.TableColumns(conf.Database, opt.TableName)
		if err != nil {
			printError(err)
			return
		}

		// 出力先の設定
		writer := getWriter(opt.Out, conf)
		err = writer.WriteDDL(columns, tableRows)
		if err != nil {
			printError(err)
			return
		}
	}
}

// getWriter returns Writer according to command line argument.
func getWriter(out option.OutputType, config *config.Config) Writer {
	switch out {
	case option.CONSOLE:
		return NewConsole(os.Stdout, config)
	case option.CSV:
		return NewCSV(os.Stdout, config)
	default:
		return NewConsole(os.Stdout, config)
	}
}

func printError(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
}
