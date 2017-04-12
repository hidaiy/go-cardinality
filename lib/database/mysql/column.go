package mysql

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

type Index struct {
	Name       string
	TableName  string
	ColumnName string
}

func (c *Column) Table() string {
	return c.TableName
}

//
func (c *Column) Column() string {
	return c.ColumnName
}

// Indexes returns indexes belongs with this column.
func (c *Column) indexes() ([]Index, error) {
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

// IndexNames returns index names belongs with this column.
func (c *Column) IndexNames() ([]string, error) {
	indexes, err := c.indexes()
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
		return c.distinctRows, nil
	}
	sql := fmt.Sprintf(
		"SELECT count(distinct `%s`) as count from %s.%s",
		c.ColumnName, c.DatabaseName, c.TableName)
	err = c.DB.Raw(sql).Row().Scan(&ret)
	c.distinctRows = ret
	return
}
