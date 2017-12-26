package dbindex

import (
	db "github.com/hidaiy/go-cardinality/lib/database"
	iutil "github.com/hidaiy/go-utils/intutil"
	sutil "github.com/hidaiy/go-utils/stringutil"
	"github.com/hidaiy/go-utils/table"
	"io"
)

// ConsoleWriter
type ConsoleWriter struct {
	baseWriter
}

// NewConsole returns ConsoleWriter constructed of io.Writer and Config.
func NewConsoleWriter(out io.Writer, config *Config) *ConsoleWriter {
	return &ConsoleWriter{
		baseWriter{
			out:    out,
			config: config,
		},
	}
}

// WriteDDL writes ddl.
func (c *ConsoleWriter) WriteDDL(i *db.SchemaInformation) error {

	// get table body data
	body, err := c.createBody(i.Columns, i.TableRows, c.createRow)
	if err != nil {
		return err
	}

	table := table.New(c.out)
	err = table.Write(RESULT_HEADER, body)
	if err != nil {
		return err
	}

	return nil
}

// createRow writes
func (c *ConsoleWriter) createRow(i *indexGenerator) []string {
	return []string{
		i.TableName,
		i.ColumnName,
		iutil.ToString(i.TableRows),
		iutil.ToString(i.DistinctTableRows),
		iutil.ToString(i.GetColumnCardinality()),
		sutil.Cut(i.ExistingIndexNames.CSV()),
		sutil.Cut(i.GenerateCreateIndexDDL()),
		sutil.Cut(i.GenerateDropIndexDDL()),
	}
}
