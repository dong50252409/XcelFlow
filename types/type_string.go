package types

import (
	"fmt"
	"reflect"
	"xCelFlow/core"
)

type Str struct {
	*Any
}

func init() {
	TypeRegister("str", NewStr)
}

func NewStr(typeStr string) (core.IType, error) {
	return &Str{Any: &Any{}}, nil
}

func (s *Str) ParseString(str string) (any, error) {
	return str, nil
}

func (s *Str) Convert(val any) string {
	return fmt.Sprintf(`"%s"`, val)
}

func (s *Str) String() string {
	return "str"
}

func (s *Str) DefaultValue() any {
	return ""
}

func (s *Str) Kind() reflect.Kind {
	return reflect.String
}

func (s *Str) CheckFunc() func(any) bool {
	return func(v any) bool {
		_, ok := v.(string)
		return ok
	}
}
