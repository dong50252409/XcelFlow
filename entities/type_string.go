package entities

import (
	"fmt"
	"reflect"
)

type Str struct {
	Field *Field
}

func init() {
	TypeRegister("str", NewStr)
}

func NewStr(_ string, field *Field) (ITypeSystem, error) {
	return &Str{Field: field}, nil
}

func (s *Str) ParseString(str string) (any, error) {
	return str, nil
}

func (*Str) Convert(val any) string {
	return fmt.Sprintf(`"%s"`, val)
}

func (s *Str) String() string {
	return "str"
}

func (s *Str) DefaultValue() string {
	return `""`
}

func (s *Str) Kind() reflect.Kind {
	return reflect.String
}

func (s *Str) CheckFunc() func(any) bool {
	return func(v any) bool {
		_, ok := v.(string)
		if !ok {
			_, ok = v.(AnyT)
			return ok
		}
		return ok
	}
}
