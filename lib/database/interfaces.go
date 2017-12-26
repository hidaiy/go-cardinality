package database

// Column
type Column interface {
	Column() string
	Table() string
	IndexNames() ([]string, error)
	DistinctRows() (int, error)
}

// SchemaProvider is interface for writer.
type SchemaProvider interface {
	GetSchemaInformation(string, []string) *SchemaInformation
}

// SchemaInformation
type SchemaInformation struct {
	TableRows TableRows
	Columns   []Column
}

// NewSchemaInformation is constructor.
func NewSchemaInformation(tableRows TableRows, columns []Column) *SchemaInformation {
	return &SchemaInformation{
		TableRows: tableRows,
		Columns:   columns,
	}
}

// TableRows has table name amd rows.
type TableRows map[string]int

// NewTableRows is constructor.
func NewTableRows() TableRows {
	return make(map[string]int)
}

// GetRows returns rows searched with given table name.
func (t TableRows) GetRows(tableName string) (int, bool) {
	rows, ok := t[tableName]
	return rows, ok
}
