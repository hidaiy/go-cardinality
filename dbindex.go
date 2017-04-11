package main

import (
	"errors"
	"fmt"
	"github.com/hidai620/go-cardinality/config"
	"github.com/hidai620/go-cardinality/database"
	"github.com/hidai620/go-cardinality/database/mysql"
	. "github.com/hidai620/go-cardinality/dbindex"
	"github.com/hidai620/go-cardinality/option"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)

	// コマンドラインオプションのパース
	opt, err := option.Parse()
	if err != nil {
		logger.Println(err)
		return
	}

	//　設定ファイルの読み込み
	conf, err := config.Load(opt.ConfigPath)
	if err != nil {
		logger.Println(err)
		return
	}

	// 管理スキーマの取得
	dataBase, err := getDatabase(conf, logger)
	if err != nil {
		logger.Println(err)
		return
	}
	info := dataBase.GetSchemaInformation(conf.Database, opt.TableNames)

	writer := getWriter(opt.Out, conf)
	err = writer.WriteDDL(info.Columns, info.TableRows)
	if err != nil {
		logger.Println(err)
		return
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

func getDatabase(conf *config.Config, logger *log.Logger) (database.Database, error) {
	db, err := database.Connect(conf)
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	switch conf.Dialect {
	case "mysql":
		return mysql.MySQL{Logger: logger, DB: db}, nil
	default:
		return nil, errors.New(fmt.Sprintf("database not found:%#v", conf.Dialect))
	}
}
