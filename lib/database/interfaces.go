package database

type Column interface {
	Column() string
	Table() string
	IndexNames() ([]string, error)
	DistinctRows() (int, error)
}

type SchemaProvider interface {
	GetSchemaInformation(string, []string) *SchemaInformation
}

type SchemaInformation struct {
	TableRows TableRows
	Columns   []Column
}

func NewSchemaInformation(tableRows TableRows, columns []Column) *SchemaInformation {
	return &SchemaInformation{
		TableRows: tableRows,
		Columns:   columns,
	}
}

type TableRows map[string]int

// GetRows returns rows searched with given table name.
func (t TableRows) GetRows(tableName string) (int, bool) {
	rows, ok := t[tableName]
	return rows, ok
}
