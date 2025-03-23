package decorators

import (
	"fmt"
	"xCelFlow/core"
	"xCelFlow/util"
)

type DecoDefault struct {
	DefaultValue string
}

func init() {
	decoratorRegister("default", newDefault)
}

func newDefault(_ *core.Table, field *core.Field, str string) error {
	if param := util.SubParam(str); param != "" {
		field.Decorators["default"] = &DecoDefault{param}
		return nil
	}
	return fmt.Errorf("参数格式错误 default(默认值)")
}

func (*DecoDefault) Name() string {
	return "default"
}

func (d *DecoDefault) RunFieldDecorator(_ *core.Table, field *core.Field) error {
	val, err := field.Type.ParseString(d.DefaultValue)
	if err != nil {
		return err
	}
	field.DefaultValue = val
	return nil
}
