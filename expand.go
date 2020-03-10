package seed

import (
	"strconv"
	"time"
)

// OtherToString 其他数据类型转string
func OtherToString(data interface{}) string {
	if data == nil {
		return ""
	}
	switch data.(type) {
	case int:
		data = strconv.Itoa(data.(int))
	}

	return data.(string)
}

// SetDefaultValue 设置默认值
func SetDefaultValue(data map[string]string, colName string, colType string, colDefault string) {
	if colType == "datetime" && colDefault == "" {
		data[colName] = time.Now().Format("2006-01-02 15:04:05")
		return
	}

	data[colName] = colDefault
}
