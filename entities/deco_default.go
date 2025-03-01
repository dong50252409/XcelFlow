package entities

import (
	"fmt"
	"xCelFlow/util"
)

type Default struct {
	DefaultValue string
}

func init() {
	decoratorRegister("default", newDefault)
}

func newDefault(_ *Table, field *Field, str string) error {
	if param := util.SubParam(str); param != "" {
		field.Decorators["default"] = &Default{param}
		return nil
	}
	return fmt.Errorf("参数格式错误 default(默认值)")
}

func (*Default) Name() string {
	return "default"
}

func (d *Default) RunFieldDecorator(_ *Table, field *Field) error {
	parseString, err := field.Type.ParseString(d.DefaultValue)
	if err != nil {
		return err
	}
	field.DefaultValue = field.Type.Convert(parseString)
	return nil
}
