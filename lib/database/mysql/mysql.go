package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/hidaiy/go-cardinality/lib/database"
	"log"
)

// mySQL
type mySQL struct {
	Logger *log.Logger
	DB     *gorm.DB
}

// New is a constructor of mySQL
func New(logger *log.Logger, db *gorm.DB) *mySQL {
	return &mySQL{
		Logger: logger,
		DB:     db,
	}
}

// GetSchemaInformation returns SchemaInformation including rows of tables and column information.
func (m mySQL) GetSchemaInformation(databaseName string, tableNames []string) *database.SchemaInformation {
	// Get information schema
	informationSchema := NewInformationSchema(m.DB)

	// Get table information having each table name and rows pair.
	tableRows, err := informationSchema.TableRows(databaseName, tableNames)
	if err != nil {
		m.Logger.Println(err)
		return nil
	}

	// Get columns having column's distinct rows, index names, and more.
	columns, err := informationSchema.TableColumns(databaseName, tableNames)
	if err != nil {
		m.Logger.Println(err)
		return nil
	}

	return database.NewSchemaInformation(tableRows, columns)
}
