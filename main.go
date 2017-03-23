package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type Config struct {
	user     string
	password string
	host     string
	port     int
	database string
}

type table struct {
	name string
	rows int
}

func craeteDBConnectString(c Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.user, c.password, c.host, c.port, "mysql")
}

// テーブル単位の件数の取得
func listTables(db *sql.DB, databaseName string) ([]table, error) {
	stmt, err := db.Prepare(`
 select table_name, table_rows
   from information_schema.tables
  where table_schema = ?
    and table_rows is not null`)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, err
	}
	rows, err := stmt.Query(databaseName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, err
	}
	defer rows.Close()

	// レスポンスの生成
	ret := make([]table, 0, 10)
	for rows.Next() {
		var (
			tableName string
			tableRows int
		)
		rows.Scan(&tableName, &tableRows)
		ret = append(ret, table{name: tableName, rows: tableRows})
	}
	return ret, nil
}

func main() {
	config := Config{
		user:     "root",
		password: "root",
		host:     "127.0.0.1",
		port:     3306,
		database: "mydatabase",
	}

	// DB接続
	db, err := sql.Open("mysql", craeteDBConnectString(config))
	if err != nil {
		printError(err)
		return
	}
	defer db.Close()

	// テーブル単位の件数の取得
	tables, err := listTables(db, config.database)
	if err != nil {
		printError(err)
		return
	}

	// 件数の表示
	for _, t := range tables {
		fmt.Printf("tableName:%s, rows: %d\n", t.name, t.rows)
	}
}

func printError(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
}
