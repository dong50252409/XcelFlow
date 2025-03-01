package entities

import (
	"fmt"
	"strings"
	"xCelFlow/util"
)

// ForeignKey 外键引用
type ForeignKey struct {
	TableName string
	FieldName string
}

func init() {
	decoratorRegister("f_key", newForeignKey)
}

func newForeignKey(_ *Table, field *Field, str string) error {
	if param := util.SubParam(str); param != "" {
		if l := strings.Split(param, ","); len(l) == 2 {
			field.Decorators["f_key"] = &ForeignKey{TableName: l[0], FieldName: l[1]}
			return nil
		}
	}
	return fmt.Errorf("参数格式错误 f_key(表名,字段名)")
}

func (f *ForeignKey) Name() string {
	return "f_key"
}

func (f *ForeignKey) RunFieldDecorator(tbl *Table, field *Field) error {
	// TODO 实现读取外键数据
	return nil
}
