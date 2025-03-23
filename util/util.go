package util

import (
	"fmt"
	"strconv"
	"strings"
)

// SubTableName 获取表名
func SubTableName(filename string) string {
	index := strings.Index(filename, "(")
	lastIndex := strings.LastIndex(filename, ")")
	if index == -1 || lastIndex == -1 {
		return ""
	}
	return filename[index+1 : lastIndex]
}

// GetKey 获取key以及未解析的参数
func GetKey(typeStr string) (string, string) {
	index := strings.Index(typeStr, "(")
	if index != -1 {
		return typeStr[:index], typeStr[index:]
	} else {
		return typeStr, ""
	}
}

// SubParam 获取未解析的参数
func SubParam(str string) string {
	index := strings.Index(str, "(")
	if index != -1 {
		end := strings.LastIndex(str, ")")
		if end == -1 {
			return ""
		}
		return str[index+1 : end]
	}
	return ""
}

// ToCell 将行索引和列索引转换为excel的单元格坐标,索引从0开始
func ToCell(row int, col int) string {
	return fmt.Sprintf("%c:%d", 'A'+col, row+1)
}

// Quoted 转义特殊字符
func Quoted(str string) string {
	quoted := strconv.Quote(str)
	return quoted[1 : len(quoted)-1]
}
