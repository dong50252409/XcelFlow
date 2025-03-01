package entities

import (
	"fmt"
	"xCelFlow/config"
	"xCelFlow/util"
)

type Unique struct {
}

func init() {
	decoratorRegister("u_key", newUnique)
}

func newUnique(_ *Table, field *Field, _ string) error {
	field.Decorators["u_key"] = &Unique{}
	return nil
}

func (*Unique) Name() string {
	return "u_key"
}

func (*Unique) RunFieldDecorator(tbl *Table, field *Field) error {
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
