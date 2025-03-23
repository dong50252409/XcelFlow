package types

import (
	"fmt"
	"reflect"
	"xCelFlow/core"
)

type Any struct {
}

func init() {
	TypeRegister("any", NewAny)
}

func NewAny(typeStr string) (core.IType, error) {
	return &Any{}, nil
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

func (r *Any) DefaultValue() any {
	return nil
}

func (r *Any) Kind() reflect.Kind {
	return reflect.Interface
}

func (r *Any) CheckFunc() func(any) bool {
	return func(v any) bool {
		_, ok := v.(core.AnyT)
		return ok
	}
}

func (r *Any) DecoratorStr() string {
	return ""
}

func (r *Any) IsReferenceType() bool {
	return true
}

func (r *Any) MethodStr() string {
	return ""
}
