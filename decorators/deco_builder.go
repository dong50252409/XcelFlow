package decorators

import (
	"fmt"
	"xCelFlow/core"
	"xCelFlow/util"
)

var (
	decoratorRegistry = make(map[string]func(tbl *core.Table, field *core.Field, str string) error)
)

func decoratorRegister(key string, cls func(tbl *core.Table, field *core.Field, str string) error) {
	decoratorRegistry[key] = cls
}

func NewDecorator(tbl *core.Table, field *core.Field, str string) error {
	key, args := util.GetKey(str)
	cons, ok := decoratorRegistry[key]
	if !ok {
		return fmt.Errorf("%s 装饰器不存在", key)
	}

	return cons(tbl, field, args)
}
