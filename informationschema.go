package dbindex

import "github.com/jinzhu/gorm"

type InformationSchema struct {
	DB *gorm.DB
}

func NewInformationSchema(db *gorm.DB) *InformationSchema {
	return &InformationSchema{
		DB: db,
	}
}

// テーブルごとのレコード件数を返す。
// テーブル名をキー、値をレコード件数のmapの形で返す。
func (inf *InformationSchema) TableRows(databaseName string) (map[string]int, error) {
	tables, err := inf.Tables(databaseName)
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
func (inf *InformationSchema) TableColumns(databaseName string) ([]Column, error) {
	var columns []Column
	result := inf.DB.Raw(`
	         select c.table_schema as database_name,
	                c.table_name,
	                c.column_name
	         from information_schema.columns c
	         join information_schema.tables t
	           on c.table_name = t.table_name
	          and t.table_type = 'BASE TABLE'
	         where c.table_schema = ?`, databaseName).Scan(&columns)

	if result.Error != nil {
		return nil, result.Error
	}

	// カラムにDBコネクションを追加
	for i := 0; i < len(columns); i++ {
		columns[i].DB = inf.DB
	}

	return columns, result.Error
}

// テーブル単位の件数の取得
func (i *InformationSchema) Tables(databaseName string) ([]Table, error) {
	var ret []Table
	result := i.DB.Raw(`
                 select table_name as name,
                        table_rows as rows
                   from information_schema.tables
                  where table_schema = ?
                    and table_rows is not null
                    and table_type = 'BASE TABLE'
                    `, databaseName).Scan(&ret)
	return ret, result.Error
}
