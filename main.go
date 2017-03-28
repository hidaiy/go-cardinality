package dbindex

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

func main() {
	config := Config{
		User:      "root",
		Password:  "",
		Host:      "127.0.0.1",
		Port:      3306,
		Dialect:   "mysql",
		Database:  "employees",
		Threshold: 60,
	}

	// DB接続
	db, err := NewDatabase(config).Connect()
	if err != nil {
		printError(err)
		return
	}
	defer db.Close()

	// 管理スキーマ
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
	var output Output = NewCSV(os.Stdout, config)
	_, err = output.WriteDDL(columns, tableRows)
	if err != nil {
		printError(err)
		return
	}
}

func printError(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
}
