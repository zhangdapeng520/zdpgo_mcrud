package zdpgo_mcrud

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// getPositionStr 获取问号占位符
// @param length 问号个数
func getPositionStr(length int) string {
	var arr []string
	for i := 0; i < length; i++ {
		arr = append(arr, "?")
	}
	return strings.Join(arr, ",")
}

// Add 新增数据
// @param db MySQL连接对象
// @param tableName 表格名称
// @param columns 列名
// @param columns 列名
func Add(
	db *sql.DB,
	tableName string,
	columns []string,
	values []interface{},
) (id int64, err error) {
	if db == nil {
		err = errors.New("db is nil")
		return
	}
	if tableName == "" {
		err = errors.New("tableName is empty")
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

	columnStr := strings.Join(columns, ",")
	positionStr := getPositionStr(len(columns))
	sqlStr := fmt.Sprintf(
		"insert into %s(%s) values (%s)",
		tableName,
		columnStr,
		positionStr,
	)

	var stmt *sql.Stmt
	stmt, err = db.Prepare(sqlStr)
	if err != nil {
		return
	}
	defer stmt.Close()

	var result sql.Result
	result, err = stmt.Exec(values...)
	if err != nil {
		return
	}

	id, err = result.LastInsertId()
	return
}
