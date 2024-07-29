package zdpgo_mcrud

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// 获取更新数据的column字符串，形如：name=?,age=?
// @param columns 列名数组
// @return 拼接好的字符串
func getUpdateColumnStr(columns []string) string {
	var arr []string
	for _, v := range columns {
		tv := fmt.Sprintf("%s=?", v)
		arr = append(arr, tv)
	}
	return strings.Join(arr, ",")
}

// Update 修改数据
// @param db MySQL连接对象
// @param tableName 表格名称
// @param id 要修改的ID
// @param columns 列名
// @param columns 列名
func Update(
	db *sql.DB,
	tableName string,
	id string,
	columns []string,
	values []interface{},
) (affectedRows int64, err error) {
	if db == nil {
		err = errors.New("db is nil")
		return
	}
	if tableName == "" {
		err = errors.New("tableName is empty")
		return
	}
	if id == "" {
		err = errors.New("id is empty")
		return
	}
	if columns == nil || len(columns) == 0 {
		err = errors.New("columns is nil")
		return
	}
	if values == nil || len(values) != len(columns) {
		err = errors.New("values is nil or values length not equal to columns")
		return
	}

	columnStr := getUpdateColumnStr(columns)
	// update user set name=?,age=?
	sqlStr := fmt.Sprintf(
		"update %s set %s where id=?",
		tableName,
		columnStr,
	)

	var stmt *sql.Stmt
	stmt, err = db.Prepare(sqlStr)
	if err != nil {
		return
	}
	defer stmt.Close()

	var result sql.Result
	values = append(values, id)
	result, err = stmt.Exec(values...)
	if err != nil {
		return
	}

	affectedRows, err = result.RowsAffected()
	return
}
