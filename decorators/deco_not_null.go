package decorators

import (
	"fmt"
	"xCelFlow/config"
	"xCelFlow/core"
	"xCelFlow/util"
)

// DecoNotNull 非空
type DecoNotNull struct {
}

func init() {
	decoratorRegister("not_null", newNotNull)
}

func newNotNull(_ *core.Table, field *core.Field, _ string) error {
	field.Decorators["not_null"] = &DecoNotNull{}
	return nil
}

func (*DecoNotNull) Name() string {
	return "not_null"
}

func (*DecoNotNull) RunFieldDecorator(tbl *core.Table, field *core.Field) error {
	for rowIndex, row := range tbl.DataSetIter() {
		v := row[field.Column]
		if v == nil || v == "" {
			return fmt.Errorf("单元格：%s 数值不能为空", util.ToCell(rowIndex+config.Config.GetBodyStartRow(), field.Column))
		}
	}
	return nil
}
