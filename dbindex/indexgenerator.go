package dbindex

import (
	"fmt"
	"strings"
)

type indexGenerator struct {
	Column             Column
	TableName          string
	ColumnName         string
	TableRows          int
	DistinctTableRows  int
	Threshold          int
	IndexName          string
	ExistingIndexNames stringArray
}

type stringArray []string

func (s stringArray) CSV() string {
	return strings.Join(s, ",")
}

func NewIndexGenerator(column Column, tableRows, threshold int) (*indexGenerator, error) {
	// 重複を除いた件数
	distinctTableRows, err := column.DistinctRows()
	if err != nil {
		return nil, err
	}

	// Indexのリスト
	indexNames, err := column.ExistingIndexNames()
	if err != nil {
		return nil, err
	}

	ret := &indexGenerator{
		Column:             column,
		TableName:          column.TableName,
		ColumnName:         column.ColumnName,
		TableRows:          tableRows,
		DistinctTableRows:  distinctTableRows,
		Threshold:          threshold,
		IndexName:          indexName(column.TableName, column.ColumnName),
		ExistingIndexNames: indexNames,
	}
	return ret, nil
}

// インデックス名を返す。
func indexName(tableName, columnName string) string {
	return fmt.Sprintf("i_%s__%s", tableName, columnName)
}

// CREATEインデックス文を生成して返す。
// カーディナリティが閾値を満たさない場合は空文字を返す。
func (i *indexGenerator) GenerateCreateIndexDDL() string {
	if i.NeedToCreateIndex() {
		return fmt.Sprintf("alter table %s add index %s(%s);", i.TableName, i.IndexName, i.ColumnName)
	} else {
		return "[*1]"
	}
}

// DROPインデックス文を生成して返す。
// カーディナリティが閾値を満たさない場合は空文字を返す。
func (i *indexGenerator) GenerateDropIndexDDL() string {
	if i.NeedToCreateIndex() {
		return fmt.Sprintf("alter table %s drop index %s;", i.TableName, i.IndexName)
	} else {
		return ""
	}
}

// カラムのカーディナリティを計算して返す。
func (i *indexGenerator) GetColumnCardinality() int {
	if i.DistinctTableRows == 0 && i.TableRows == 0 {
		return 0
	}

	tmp := float64(i.DistinctTableRows) / float64(i.TableRows)
	return int(tmp * float64(100))
}

// カーディナリティが閾値以上、かつ、このカラムに対して既存のインデックスが存在しない場合trueを返す。
func (i *indexGenerator) NeedToCreateIndex() bool {
	return i.GetColumnCardinality() >= i.Threshold && len(i.ExistingIndexNames) == 0
}
