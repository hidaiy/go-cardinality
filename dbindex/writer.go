package dbindex

type Writer interface {
	//WriteRow(*Row) (int, error)
	WriteDDL([]Column, TableRows) (int, error)
}
