package dbindex

import (
	sutil "github.com/hidai620/go-cardinality/stringutil"
	"github.com/jinzhu/gorm"
)

// InformationSchema is a struct to access MySQL information schema.
type InformationSchema struct {
	DB *gorm.DB
}

// NewInformationSshema is constructor.
func NewInformationSchema(db *gorm.DB) *InformationSchema {
	return &InformationSchema{
		DB: db,
	}
}

// TableRows has rows of each tables.
type TableRows map[string]int

// GetRows returns rows searched with given table name.
func (t TableRows) GetRows(tableName string) (int, bool) {
	rows, ok := t[tableName]
	return rows, ok
}

// TableRows returns each rows of tables searched with given database name from information schema.
func (inf *InformationSchema) TableRows(databaseName string, tableNames []string) (TableRows, error) {
	tables, err := inf.Tables(databaseName, tableNames)
	if err != nil {
		return nil, err
	}
	tableRows := make(map[string]int)
	for _, t := range tables {
		tableRows[t.Name] = t.Rows
	}
	return tableRows, nil
}

// データベース内のカラムの一覧を返す
func (inf *InformationSchema) TableColumns(databaseName string, tableNames []string) ([]IColumn, error) {
	var columns []Column
	sql := `select c.table_schema as database_name,
		       c.table_name,
		       c.column_name
		  from information_schema.columns c
		  join information_schema.tables t
		    on c.table_name = t.table_name
		   and t.table_type = 'BASE TABLE'
		 where c.table_schema = ?
		`
	params := NewParams(databaseName)
	if sutil.NotEmpty(tableNames) {
		sql = sql + ` and c.table_name in (?)`
		params.Add(tableNames)
	}
	result := inf.DB.Raw(sql, params.values...).Scan(&columns)

	if result.Error != nil {
		return nil, result.Error
	}

	// カラムにDBコネクションを追加
	ret := make([]IColumn, 0, len(columns))
	for i := 0; i < len(columns); i++ {
		c := columns[i]
		c.DB = inf.DB
		ret = append(ret, &c)
	}

	return ret, result.Error
}

// テーブル単位の件数の取得
func (i *InformationSchema) Tables(databaseName string, tableNames []string) ([]Table, error) {
	var ret []Table
	sql := `select table_name as name,
			table_rows as rows
		  from information_schema.tables
		 where table_schema = ?
		   and table_rows is not null
		   and table_type = 'BASE TABLE'
		`
	param := NewParams(databaseName)
	if sutil.NotEmpty(tableNames) {
		sql = sql + ` and table_name in (?)`
		param.Add(tableNames)
	}

	result := i.DB.Raw(sql, param.values...).Scan(&ret)
	return ret, result.Error
}

type Table struct {
	DB           *gorm.DB
	DatabaseName string
	Name         string
	Rows         int
}
