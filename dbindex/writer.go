package dbindex

import . "github.com/hidai620/go-cardinality/database"

type Writer interface {
	WriteDDL([]IColumn, TableRows) error
}

