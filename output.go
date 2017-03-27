package dbindex

type Output interface {
	Write(*Row) (int, error)
	WriteStringArray([]string) (int, error)
	WriteDDL([]Column, map[string]int) (int, error)
}
