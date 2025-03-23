package decorators

import (
	"fmt"
	"xCelFlow/config"
	"xCelFlow/core"
	"xCelFlow/util"
)

// DecoPrimaryKey 主键
type DecoPrimaryKey struct {
	Fields []*core.Field
}

func init() {
	decoratorRegister("p_key", newPrimaryKey)
}

func newPrimaryKey(tbl *core.Table, field *core.Field, _ string) error {
	var pk *DecoPrimaryKey
	for _, d := range tbl.Decorators {
		if d1, ok := d.(*DecoPrimaryKey); ok {
			pk = d1
			break
		}
	}
	if pk == nil {
		pk = &DecoPrimaryKey{}
		tbl.Decorators = append(tbl.Decorators, pk)
	}
	pk.Fields = append(pk.Fields, field)
	return nil
}

func (*DecoPrimaryKey) Name() string {
	return "p_key"
}

func (pk *DecoPrimaryKey) RunTableDecorator(tbl *core.Table) error {
	var set = make(map[core.TupleT]struct{})
	for rowIndex, row := range tbl.DataSetIter() {
		var items core.TupleT
		for index, field := range pk.Fields {
			item := row[field.Column]
			if item == nil {
				return fmt.Errorf("单元格：%s 主键不能为空", util.ToCell(rowIndex+config.Config.GetBodyStartRow()-1, field.Column))
			}
			items[index] = item
		}
		if _, ok := set[items]; ok {
			return fmt.Errorf("第 %d 行 主键重复 %v", rowIndex+config.Config.GetBodyStartRow()-1, items)
		} else {
			set[items] = struct{}{}
		}
	}
	return nil
}

func (pk *DecoPrimaryKey) GetFields() []*core.Field {
	return pk.Fields
}
