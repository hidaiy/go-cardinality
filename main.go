package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	//"strings"
	"math"
)

type Config struct {
	user     string
	password string
	host     string
	port     int
	database string
}

type DB interface {
	Raw(string, ...interface{}) *DB
	Scan([]interface{}) *DB
}

type Table struct {
	Name string
	Rows int
}

type Column struct {
	TableName  string
	ColumnName string
}

type Index struct {
	Name string
}

//func NewTable(db *gorm.DB) *Tables {
//	return &Tables{db: db}
//}
//
//type Table struct {
//	db *gorm.DB
//}

// テーブル単位の件数の取得
//func (t Tables) listTables(databaseName string) ([]Table, error) {
//	var ret []Table
//	result := t.db.Raw(`
//                 select table_name as name,
//                        table_rows as rows
//                   from information_schema.tables
//                  where table_schema = ?
//                    and table_rows is not null`, databaseName).Scan(&ret)
//	return ret, result.Error
//}

var CSV_HEADER = []string{
	"table_name",
	"column_name",
	"table_count",
	"distinct_count",
	"cardinarity",
	"indexes",
	"create_index_ddl",
	"drop_index_ddl",
}

func createDBConnectString(c Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.user, c.password, c.host, c.port, "mysql")
}

// テーブル単位の件数の取得
func listTables(db *gorm.DB, databaseName string) ([]Table, error) {
	var ret []Table
	result := db.Raw(`
                 select table_name as name,
                        table_rows as rows
                   from information_schema.tables
                  where table_schema = ?
                    and table_rows is not null
                    and table_type = 'BASE TABLE'
                    `, databaseName).Scan(&ret)
	return ret, result.Error
}

// テーブル単位の件数の取得
func listIndexes(db *gorm.DB, databaseName, tableName, columnName string) ([]Index, error) {
	var ret []Index
	result := db.Raw(`
                 select index_name as name
                   from information_schema.statistics
                  where table_schema = ?
                    and table_name = ?
                    and column_name = ?`,
		databaseName, tableName, columnName).Scan(&ret)
	return ret, result.Error
}
func listTableColumns(db *gorm.DB, databaseName string) ([]Column, error) {
	var ret []Column
	result := db.Raw(`
	         select c.table_name,
	                c.column_name
	         from information_schema.columns c
	         join information_schema.tables t
	           on c.table_name = t.table_name
	          and t.table_type = 'BASE TABLE'
	         where c.table_schema = ?`, databaseName).Scan(&ret)
	return ret, result.Error
}
func count(db *gorm.DB, databaseName, tableName string) (int, error) {
	var ret int
	sql := fmt.Sprintf("SELECT count(*) as count from %s.%s", databaseName, tableName)
	error := db.Raw(sql).Row().Scan(&ret)
	return ret, error
}

func countDistinctColumn(db *gorm.DB, databaseName, tableName, columnName string) (ret int, err error) {
	sql := fmt.Sprintf("SELECT count(distinct `%s`) as count from %s.%s", columnName, databaseName, tableName)
	err = db.Raw(sql).Row().Scan(&ret)
	return
}

func main() {
	config := Config{
		user:     "root",
		password: "",
		host:     "127.0.0.1",
		port:     3306,
		database: "employees",
	}

	// DB接続
	db, err := gorm.Open("mysql", createDBConnectString(config))
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
	//fmt.Println(len(tables))

	tableCount := make(map[string]int)
	for _, table := range tables {
		tableCount[table.Name] = table.Rows
	}
	fmt.Println(tableCount)

	//for _, t := range tables {
	//	fmt.Printf("tableName:%-25s rows: %5d\n", t.Name, t.Rows)
	//}

	// カラムの取得
	columns, err := listTableColumns(db, config.database)
	if err != nil {
		printError(err)
		return
	}
	fmt.Println("columns", len(columns))
	//for _, column := range columns {
	//	fmt.Printf("tableName:%-25s column: %5s\n", column.TableName, column.ColumnName)
	//}

	// インデックスの取得
	var row []string
	for _, column := range columns {
		fmt.Println(column.TableName, column.ColumnName)
		indexes, err := listIndexes(db, config.database, column.TableName, column.ColumnName)
		if err != nil {
			printError(err)
			return
		}

		//for _, index := range indexes {
		//	fmt.Print(index)
		//}

		// レコード件数
		cnt, ok := tableCount[column.TableName]
		if !ok {
			fmt.Fprintln(os.Stderr, "table count not found:", column.TableName)
			return
		}

		count, err := countDistinctColumn(db, config.database, column.TableName, column.ColumnName)
		if err != nil {
			printError(err)
			return
		}
		//fmt.Println(cnt, count)

		println(indexes)

		row = []string{
			column.TableName,
			column.ColumnName,
			ToString(cnt),
			ToString(count),
			ToString(int(math.Ceil(math.Mod(float64(count), float64(cnt))))),
			"indexes",
			"create_ddl",
			"drop_ddl",
		}
		fmt.Println(row)

		//fmt.Println(strings.Join(indexes, ","))
		//fmt.Printf("tableName:%-25s column:%-15s  index: %5s\n", column.TableName, column.ColumnName, strings.Join(indexes, ","))
	}
}
func ToString(base int) string {
	return fmt.Sprintf("%d", base)
}

func printError(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
}
