package types

import (
	"fmt"
	"reflect"
	"xCelFlow/core"
)

type Lang struct {
	*Any
}

func init() {
	TypeRegister("lang", NewLang)
}

func NewLang(typeStr string) (core.IType, error) {
	return &Lang{Any: &Any{}}, nil
}

func (l *Lang) ParseString(str string) (any, error) {
	return str, nil
}

func (l *Lang) Convert(val any) string {
	return fmt.Sprintf(`"%s"`, val)
}

func (l *Lang) String() string {
	return "lang"
}

func (l *Lang) DefaultValue() any {
	return ""
}

func (l *Lang) Kind() reflect.Kind {
	return reflect.String
}

func (l *Lang) CheckFunc() func(any) bool {
	return func(v any) bool {
		_, ok := v.(string)
		return ok
	}
}
