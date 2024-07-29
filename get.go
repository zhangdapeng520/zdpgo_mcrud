package zdpgo_mcrud

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func GetBy(
	db *sql.DB,
	tableName string,
	columns []string,
	conditions map[string]interface{},
) (data []map[string]interface{}, err error) {
	if db == nil {
		err = errors.New("db is nil")
		return
	}
	if tableName == "" {
		err = errors.New("tableName is empty")
		return
	}

	var columnStr string
	if columns == nil || len(columns) == 0 {
		columnStr = "*"
	} else {
		columnStr = strings.Join(columns, ",")
	}
	sqlStr := fmt.Sprintf(
		"select %s from %s",
		columnStr,
		tableName,
	)

	// 构造查询条件
	whereValues := []interface{}{}
	whereKeys := []string{}
	if conditions != nil && len(conditions) > 0 {
		// select * from user
		// select * from user where k=v
		for k, v := range conditions {
			whereKeys = append(whereKeys, fmt.Sprintf("%s=?", k))
			whereValues = append(whereValues, v)
		}
		whereStr := strings.Join(whereKeys, ",")
		sqlStr += " where " + whereStr
	}

	// 准备执行查询
	var stmt *sql.Stmt
	stmt, err = db.Prepare(sqlStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	// 执行查询
	var rows *sql.Rows
	rows, err = stmt.Query(whereValues...)
	defer rows.Close()
	if rows == nil {
		err = errors.New("rows is nil")
		return
	}

	count := len(columns)                   // 列的个数
	values := make([]interface{}, count)    // 一组数据的值
	valuePtrs := make([]interface{}, count) // 一组数据的值的对应地址
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i] // 将列的数据的值的地址取出来，赋值给地址值
		}
		err = rows.Scan(valuePtrs...) // 获取各列的值，放在对应地址中
		if err != nil {
			return
		}
		item := make(map[string]interface{}) // 构建列名和值的对应关系 {name:张三,age:22}
		for i, col := range columns {
			var v interface{}     // 临时值
			val := values[i]      // 对应的值
			b, ok := val.([]byte) // 判断能不能转换为字节数组，实际上就是判断是不是字符串
			if ok {
				v = string(b) // 转换为字符串
			} else {
				v = val
			}
			item[col] = v
		}
		data = append(data, item)
	}
	return
}
