package dbindex

type Writer interface {
	WriteDDL([]Column, TableRows) error
}
