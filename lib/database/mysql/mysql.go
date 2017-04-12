package mysql

import (
	"github.com/hidai620/go-cardinality/lib/database"
	"github.com/jinzhu/gorm"
	"log"
)

type mySQL struct {
	Logger *log.Logger
	DB     *gorm.DB
}

func New(logger *log.Logger, db *gorm.DB) *mySQL {
	return &mySQL{
		Logger: logger,
		DB:     db,
	}
}

func (m mySQL) GetSchemaInformation(databaseName string, tableNames []string) *database.SchemaInformation {
	// 管理スキーマの取得
	informationSchema := NewInformationSchema(m.DB)

	// テーブル単位の件数の取得
	tableRows, err := informationSchema.TableRows(databaseName, tableNames)
	if err != nil {
		m.Logger.Println(err)
		return nil
	}

	// カラムの取得
	columns, err := informationSchema.TableColumns(databaseName, tableNames)
	if err != nil {
		m.Logger.Println(err)
		return nil
	}

	return database.NewSchemaInformation(tableRows, columns)
}
