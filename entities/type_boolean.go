package entities

import (
	"reflect"
	"strconv"
)

type Boolean struct {
	Field *Field
}

func init() {
	TypeRegister("bool", NewBoolean)
}

func NewBoolean(_ string, field *Field) (ITypeSystem, error) {
	t := &Boolean{Field: field}
	return t, nil
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

func (b *Boolean) DefaultValueStr() string {
	return "false"
}

func (b *Boolean) Kind() reflect.Kind {
	return reflect.Bool
}

func (b *Boolean) CheckFunc() func(any) bool {
	return func(v any) bool { _, ok := v.(bool); return ok }
}

func (b *Boolean) DecoratorStr() string {
	return ""
}

func (b *Boolean) IsReferenceType() bool {
	return false
}

func (b *Boolean) MethodStr() string {
	return ""
}
