package decorators

import (
	"fmt"
	"os"
	"path/filepath"
	"xCelFlow/config"
	"xCelFlow/core"
	"xCelFlow/util"
)

// DecoResource 资源引用
type DecoResource struct {
	Path string
}

func init() {
	decoratorRegister("resource", newResource)
}

func newResource(_ *core.Table, field *core.Field, str string) error {
	if param := util.SubParam(str); param != "" {
		if str != "" {
			wd, _ := os.Getwd()
			path := filepath.Join(wd, param)
			if _, err := os.Stat(path); err != nil {
				return fmt.Errorf("参数路径不存在 完整路径：%s", path)
			}
			field.Decorators["resource"] = &DecoResource{Path: path}
			return nil
		}
	}
	return fmt.Errorf("参数格式错误 resource(路径)")
}

func (r *DecoResource) Name() string {
	return "resource"
}

func (r *DecoResource) RunFieldDecorator(tbl *core.Table, field *core.Field) error {
	for rowIndex, row := range tbl.DataSetIter() {
		if v := row[field.Column]; v != nil && v != "" {
			if _, err := os.Stat(filepath.Join(r.Path, v.(string))); err != nil {
				return fmt.Errorf("单元格：%s 资源不存在 %s", util.ToCell(rowIndex+config.Config.GetBodyStartRow(), field.Column), v)
			}
		}
	}
	return nil
}
