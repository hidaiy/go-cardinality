package dbindex

import (
	"fmt"
	db "github.com/hidai620/go-cardinality/lib/database"
	iutil "github.com/hidai620/go-utils/intutil"
	"io"
	"strings"
)

// CSVWriter
type CSVWriter struct {
	baseWriter
}

// NewCSV returns CSVWriter as pointer.
func NewCSVWriter(out io.Writer, config *Config) *CSVWriter {
	return &CSVWriter{
		baseWriter{
			out:    out,
			config: config,
		},
	}
}

// WriteDDL writes database index ddl for each columns as csv.
func (c *CSVWriter) WriteDDL(i *db.SchemaInformation) error {
	// get table body data
	body, err := c.createBody(i.Columns, i.TableRows, c.createRow)
	if err != nil {
		return err
	}

	// write csv header
	_, err = c.writeRow(RESULT_HEADER)
	if err != nil {
		return err
	}

	// write csv body
	for _, line := range body {
		c.writeRow(line)
		if err != nil {
			return err
		}
	}

	return nil
}

// writeRow writes csv row with string array.
func (c CSVWriter) writeRow(array []string) (int, error) {
	line := strings.Join(array, ", ") + fmt.Sprintln()
	return c.out.Write([]byte(line))
}

// createRow returns csv row as string array.
func (c *CSVWriter) createRow(i *indexGenerator) []string {
	return []string{
		i.TableName,
		i.ColumnName,
		iutil.ToString(i.TableRows),
		iutil.ToString(i.DistinctTableRows),
		iutil.ToString(i.GetColumnCardinality()),
		i.ExistingIndexNames.CSV(),
		i.GenerateCreateIndexDDL(),
		i.GenerateDropIndexDDL(),
	}
}
