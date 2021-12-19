/*
问：是否应该 Wrap 这个 error?
答: 尽量只在业务层使用wrap，否则会重复打印堆栈信息

可以考虑使用错误码的形式，这样比较方便，只是在某些特殊的错误时需要特殊处理
*/
package main

import (
	"database/sql"
	"fmt"
	"strings"
)

var notFoundCode = 4001
var systemErrorCode = 5001

func isNoRow(err error) bool {
	return strings.HasPrefix(err.Error(), fmt.Sprintf("%d", notFoundCode))
}

func business_query(sqlString string) error {
	result, err := Dao(sqlString)
	if isNoRow(err) {
		// 数据查不到的处理，可以在这里加日志
	} else if err != nil {
		// 正常逻辑
	} else {
		// 异常处理
		return err
	}
	return nil
}

func mock_myqsl_query(sqlString string) (string, error) {
	sqlString += "111"
	return "", sql.ErrNoRows
}

func Dao(sqlString string) (string, error) {
	result, err := mock_myqsl_query(sqlString)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("%d, No row Error", notFoundCode)
	} else if err != nil {
		return "nil", fmt.Errorf("%d, system Error", systemErrorCode)
	} else {
		return result, nil
	}
}
