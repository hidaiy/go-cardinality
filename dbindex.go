package main

import (
	"fmt"
	cnf "github.com/hidai620/go-mysql-study/config"
	. "github.com/hidai620/go-mysql-study/dbindex"
	"github.com/hidai620/go-mysql-study/option"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

func main() {
	// コマンドラインオプションのパース
	commandLineOption, err := option.Parse()
	if err != nil {
		printError(err)
		return
	}

	//　設定ファイルの読み込み
	config, err := cnf.Load(commandLineOption.ConfigPath)
	if err != nil {
		printError(err)
		return
	}

	// DB接続
	db, err := Connect(config)
	if err != nil {
		printError(err)
		return
	}
	defer db.Close()

	// 管理スキーマの取得
	informationSchema := NewInformationSchema(db)

	// テーブル単位の件数の取得
	tableRows, err := informationSchema.TableRows(config.Database)
	if err != nil {
		printError(err)
		return
	}

	// カラムの取得
	columns, err := informationSchema.TableColumns(config.Database)
	if err != nil {
		printError(err)
		return
	}

	// 出力先の設定
	writer := getWriter(commandLineOption.Out, config)
	_, err = writer.WriteDDL(columns, tableRows)
	if err != nil {
		printError(err)
		return
	}
}

func getWriter(out option.Out, config *cnf.Config) Writer {
	switch out {
	case option.CONSOLE:
		return NewConsole(os.Stdout, config.Threshold)
	case option.CSV:
		return NewCSV(os.Stdout, config.Threshold)
	default:
		return NewConsole(os.Stdout, config.Threshold)
	}
}

func printError(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
}
