package decorators

import (
	"fmt"
	"xCelFlow/core"
	"xCelFlow/util"
)

type DecoRefTable struct {
	TableName string
	Table     *core.Table
}

func init() {
	decoratorRegister("ref_table", newRefTable)
}

func newRefTable(_ *core.Table, field *core.Field, str string) error {
	if param := util.SubParam(str); param != "" {
		field.Decorators["ref_table"] = &DecoRefTable{TableName: param, Table: &core.Table{}}
		return nil
	}
	return fmt.Errorf("参数格式错误 ref_table(表名,字段名)")
}

func (r *DecoRefTable) RunFieldDecorator(tbl *core.Table, field *core.Field) error {
	return nil
}

func (r *DecoRefTable) Name() string {
	return "ref_table"
}
