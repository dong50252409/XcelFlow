package entities

import (
	"cfg_exporter/config"
	"cfg_exporter/util"
	"fmt"
	"os"
	"path/filepath"
)

// Resource 资源引用
type Resource struct {
	Path string
}

func init() {
	decoratorRegister("resource", newResource)
}

func newResource(_ *Table, field *Field, str string) error {
	if param := util.SubParam(str); param != "" {
		if str != "" {
			wd, _ := os.Getwd()
			path := filepath.Join(wd, param)
			if _, err := os.Stat(path); err != nil {
				return fmt.Errorf("参数路径不存在 完整路径：%s", path)
			}
			field.Decorators["resource"] = &Resource{Path: path}
			return nil
		}
	}
	return fmt.Errorf("参数格式错误 resource(路径)")
}

func (r *Resource) Name() string {
	return "resource"
}

func (r *Resource) RunFieldDecorator(tbl *Table, field *Field) error {
	for rowIndex, row := range tbl.DataSetIter() {
		if v := row[field.Column]; v != nil && v != "" {
			if _, err := os.Stat(filepath.Join(r.Path, v.(string))); err != nil {
				return fmt.Errorf("单元格：%s 资源不存在 %s", util.ToCell(rowIndex+config.Config.GetBodyStartRow(), field.Column), v)
			}
		}
	}
	return nil
}
