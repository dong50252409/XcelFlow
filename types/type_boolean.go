package types

import (
	"reflect"
	"strconv"
	"xCelFlow/core"
)

type Boolean struct {
	*Any
}

func init() {
	TypeRegister("bool", NewBoolean)
}

func NewBoolean(typeStr string) (core.IType, error) {
	return &Boolean{}, nil
}

func (b *Boolean) ParseString(str string) (any, error) {
	parseBool, err := strconv.ParseBool(str)
	if err != nil {
		return nil, NewTypeErrorParseFailed(b, str)
	}
	return parseBool, nil
}

func (*Boolean) Convert(val any) string {
	return strconv.FormatBool(val.(bool))
}

func (b *Boolean) String() string {
	return "bool"
}

func (b *Boolean) DefaultValue() any {
	return false
}

func (b *Boolean) Kind() reflect.Kind {
	return reflect.Bool
}

func (b *Boolean) CheckFunc() func(any) bool {
	return func(v any) bool {
		_, ok := v.(bool)
		return ok
	}
}
