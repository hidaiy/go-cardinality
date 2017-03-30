package dbindex

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Column struct {
	DB           *gorm.DB
	DatabaseName string
	TableName    string
	ColumnName   string
	distinctRows int
}

// テーブル単位の件数の取得
func (c *Column) ExistingIndexes() ([]Index, error) {
	var ret []Index
	result := c.DB.Raw(`
                 select index_name as name
                        ,table_name
                        ,column_name
                   from information_schema.statistics
                  where table_schema = ?
                    and table_name = ?
                    and column_name = ?`,
		c.DatabaseName, c.TableName, c.ColumnName).Scan(&ret)
	return ret, result.Error
}

// 既存のIndex名のリストを返す
func (c *Column) ExistingIndexNames() ([]string, error) {
	indexes, err := c.ExistingIndexes()
	if err != nil {
		return nil, err
	}
	var ret = make([]string, 0, len(indexes))
	for _, index := range indexes {
		ret = append(ret, index.Name)
	}
	return ret, nil
}

// 重複を省いた件数を返す。
func (c *Column) DistinctRows() (ret int, err error) {
	if c.distinctRows != 0 {
		//fmt.Println("distinct rows returned from property")
		return c.distinctRows, nil
	}
	sql := fmt.Sprintf(
		"SELECT count(distinct `%s`) as count from %s.%s",
		c.ColumnName, c.DatabaseName, c.TableName)
	err = c.DB.Raw(sql).Row().Scan(&ret)
	c.distinctRows = ret
	//fmt.Println("distinct rows returned from SQL")
	return
}
