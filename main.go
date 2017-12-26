package main

import (
	"errors"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "github.com/hidaiy/go-cardinality/lib"
	db "github.com/hidaiy/go-cardinality/lib/database"
	"github.com/hidaiy/go-cardinality/lib/database/mysql"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)

	// コマンドラインオプションのパース
	opt, err := ParseCommandLineOption()
	if err != nil {
		fmt.Printf("NOTICE: %s", err)
		return
	}

	// 設定ファイルの読み込み
	conf, err := LoadConfig(opt.ConfigPath)
	if err != nil {
		logger.Println(err)
		return
	}

	// Get in
	provider, err := getSchemaProvider(conf, logger)
	if err != nil {
		logger.Println(err)
		return
	}
	info := provider.GetSchemaInformation(conf.Database, opt.TableNames)

	writer := getWriter(opt.Out, conf)
	err = writer.WriteDDL(info)
	if err != nil {
		logger.Println(err)
		return
	}
}

// getSchemaProvider returns provider having database schema information.
func getSchemaProvider(conf *Config, logger *log.Logger) (db.SchemaProvider, error) {
	db, err := ConnectDatabase(conf)
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	switch conf.Dialect {
	case "mysql":
		return mysql.New(logger, db), nil
	default:
		return nil, errors.New(fmt.Sprintf("database not found:%#v", conf.Dialect))
	}
}

// getWriter returns Writer according to command line argument.
func getWriter(out OutputType, config *Config) Writer {
	switch out {
	case CSV:
		return NewCSVWriter(os.Stdout, config)
	default:
		return NewConsoleWriter(os.Stdout, config)
	}
}
