package reader

import (
	"errors"
	"fmt"
	"xCelFlow/util"
)

var (
	ErrorTableNotSupported    error
	ErrorTableReadFailed      error
	ErrorTableNotSheet        error
	ErrorTableSheetHeadRepeat error
)

func errorTableNotSupported(path string) error {
	ErrorTableNotSupported = errors.New(fmt.Sprintf("配置表不支持！%s", path))
	return ErrorTableNotSupported
}

func errorTableReadFailed(path string, err error) error {
	ErrorTableReadFailed = errors.New(fmt.Sprintf("配置表读取失败！%s 错误:%s", path, err))
	return ErrorTableReadFailed
}

func errorTableNotSheet(path string) error {
	ErrorTableNotSheet = errors.New(fmt.Sprintf("没有发现可读取的数据，请检查页签名命名是否正确！%s", path))
	return ErrorTableNotSheet
}

func errorTableSheetHeadRepeat(sheetName string, colIndex int) error {
	ErrorTableSheetHeadRepeat = errors.New(fmt.Sprintf("页签！%s，单元格！%s，存在重复表头", sheetName, util.ToCell(0, colIndex)))
	return ErrorTableSheetHeadRepeat
}
