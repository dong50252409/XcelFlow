package entities

import (
	"fmt"
	"reflect"
)

type Any struct {
	Field *Field
}

func init() {
	TypeRegister("any", NewAny)
}

func NewAny(_ string, field *Field) (ITypeSystem, error) {
	return &Any{Field: field}, nil
}

func (r *Any) ParseString(str string) (any, error) {
	v, err := ParseString(str)
	if err != nil {
		return nil, NewTypeErrorParseFailed(r, fmt.Sprintf("%v", err))
	}
	return v, nil
}
func (*Any) Convert(val any) string {
	return fmt.Sprintf("%v", val)
}

func (r *Any) String() string {
	return "any"
}

func (r *Any) DefaultValue() string {
	return "nil"
}

func (r *Any) Kind() reflect.Kind {
	return reflect.Interface
}

func (r *Any) CheckFunc() func(any) bool {
	return func(v any) bool {
		_, ok := v.(string)
		if !ok {
			_, ok = v.(AnyT)
			return ok
		}
		return ok
	}
}
