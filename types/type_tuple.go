package types

import (
	"errors"
	"fmt"
	"reflect"
	"xCelFlow/core"
	"xCelFlow/util"
)

type Tuple struct {
	*Any
	T core.IType
}

func init() {
	TypeRegister("tuple", NewTuple)
}

func NewTuple(typeStr string) (core.IType, error) {
	if param := util.SubParam(typeStr); param == "" {
		return &Tuple{Any: &Any{}}, nil
	} else {
		t, err := NewType(param)
		if err != nil {
			if errors.Is(err, TypeErrorNotSupported) {
				return nil, NewTypeErrorTupleInvalid(typeStr)
			}
			return nil, err
		}
		return &Tuple{Any: &Any{}, T: t}, nil
	}
}

func (t *Tuple) ParseString(str string) (any, error) {
	if !(str[0] == '(' && str[len(str)-1] == ')') {
		return nil, NewTypeErrorParseFailed(t, str)
	}
	v, err := ParseString(str)
	if err != nil {
		return nil, NewTypeErrorParseFailed(t, str)
	}
	if v == nil {
		return v, nil
	}
	if t.T != nil {
		checkFunc := t.T.CheckFunc()
		for i, e := range v.(core.TupleT) {
			if e != nil {
				if !checkFunc(e) {
					return nil, NewTypeErrorNotMatch(t, i, e)
				}
			} else {
				break
			}
		}
	}
	return v, nil
}

func (t *Tuple) Convert(val any) string {
	return fmt.Sprintf("(%v)", val)
}

func (t *Tuple) String() string {
	return "tuple"
}

func (t *Tuple) DefaultValue() any {
	return core.TupleT{}
}

func (t *Tuple) Kind() reflect.Kind {
	return reflect.Array
}

func (t *Tuple) CheckFunc() func(any) bool {
	checkFunc := t.T.CheckFunc()
	return func(v any) bool {
		if v1, ok := v.(core.TupleT); !ok {
			return false
		} else {
			for _, e := range v1 {
				if e == nil {
					break
				}
				if !checkFunc(e) {
					return false
				}
			}
		}
		return true
	}
}
