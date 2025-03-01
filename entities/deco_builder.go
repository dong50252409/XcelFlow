package entities

import (
	"fmt"
	"xCelFlow/util"
)

var (
	decoratorRegistry = make(map[string]func(tbl *Table, field *Field, str string) error)
)

func decoratorRegister(key string, cls func(tbl *Table, field *Field, str string) error) {
	decoratorRegistry[key] = cls
}

func NewDecorator(tbl *Table, field *Field, str string) error {
	key, args := util.GetKey(str)
	cons, ok := decoratorRegistry[key]
	if !ok {
		return fmt.Errorf("装饰器不存在")
	}

	return cons(tbl, field, args)
}
