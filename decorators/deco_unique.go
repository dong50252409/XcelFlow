package decorators

import (
	"fmt"
	"xCelFlow/config"
	"xCelFlow/core"
	"xCelFlow/util"
)

type DecoUnique struct {
}

func init() {
	decoratorRegister("u_key", newUnique)
}

func newUnique(_ *core.Table, field *core.Field, _ string) error {
	field.Decorators["u_key"] = &DecoUnique{}
	return nil
}

func (*DecoUnique) Name() string {
	return "u_key"
}

func (*DecoUnique) RunFieldDecorator(tbl *core.Table, field *core.Field) error {
	var set = make(map[any]struct{})
	for rowIndex, row := range tbl.DataSetIter() {
		if v := row[field.Column]; v != nil {
			if _, ok := set[v]; ok {
				return fmt.Errorf("单元格：%s 数值重复", util.ToCell(rowIndex+config.Config.GetBodyStartRow(), field.Column))
			}
			set[v] = struct{}{}
		}
	}
	return nil
}
