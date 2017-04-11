package dbindex

type Writer interface {
	WriteDDL([]IColumn, TableRows) error
}

type SchemaInformation struct {
	tableRows TableRows
	columns   []IColumn
}

func NewSchemaInformation(tableRows TableRows, columns []IColumn) *SchemaInformation {
	return &SchemaInformation{
		tableRows: tableRows,
		columns:   columns,
	}
}
