package zdpgo_mcrud

import (
	"database/sql"
	"errors"
	"fmt"
)

// Delete 根据ID数据
// @param db MySQL连接对象
// @param tableName 表格名称
// @param id 要修改的ID
func Delete(
	db *sql.DB,
	tableName string,
	id string,
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

	sqlStr := fmt.Sprintf("delete from %s where id=?", tableName)

	var stmt *sql.Stmt
	stmt, err = db.Prepare(sqlStr)
	if err != nil {
		return
	}
	defer stmt.Close()

	var result sql.Result
	result, err = stmt.Exec(id)
	if err != nil {
		return
	}

	affectedRows, err = result.RowsAffected()
	return
}
