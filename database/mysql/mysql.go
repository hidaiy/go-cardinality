package mysql

import (
	"github.com/hidai620/go-cardinality/database"
	"github.com/jinzhu/gorm"
	"log"
)

type MySQL struct {
	Logger *log.Logger
	DB     *gorm.DB
}

func (m MySQL) GetSchemaInformation(databaseName string, tableNames []string) *database.SchemaInformation {
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
