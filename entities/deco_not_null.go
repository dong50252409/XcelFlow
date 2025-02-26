package entities

import (
	"cfg_exporter/config"
	"cfg_exporter/util"
	"fmt"
)

// NotNull 非空
type NotNull struct {
}

func init() {
	decoratorRegister("not_null", newNotNull)
}

func newNotNull(_ *Table, field *Field, _ string) error {
	field.Decorators["not_null"] = &NotNull{}
	return nil
}

func (*NotNull) Name() string {
	return "not_null"
}

func (*NotNull) RunFieldDecorator(tbl *Table, field *Field) error {
	for rowIndex, row := range tbl.DataSetIter() {
		v := row[field.Column]
		if v == nil || v == "" {
			return fmt.Errorf("单元格：%s 数值不能为空", util.ToCell(rowIndex+config.Config.GetBodyStartRow(), field.Column))
		}
	}
	return nil
}
