package dbindex

type Output interface {
	WriteRow(*Row) (int, error)
	WriteStringArray([]string) (int, error)
	WriteDDL([]Column, TableRows) (int, error)
}
