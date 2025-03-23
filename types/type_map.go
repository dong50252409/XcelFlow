package types

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"xCelFlow/core"
	"xCelFlow/util"
)

type Map struct {
	*Any
	KeyT   core.IType
	ValueT core.IType
}

func init() {
	TypeRegister("map", NewMap)
}

func NewMap(typeStr string) (core.IType, error) {
	if param := util.SubParam(typeStr); param == "" {
		return &Map{Any: &Any{}}, nil
	} else {
		if l := strings.SplitN(param, ",", 2); len(l) == 2 {
			kT, err := NewType(l[0])
			if err != nil {
				if errors.Is(err, TypeErrorNotSupported) {
					return nil, NewTypeErrorMapKeyInvalid(l[0])
				}
				return nil, err
			}

			vT, err := NewType(l[1])
			if err != nil {
				if errors.Is(err, TypeErrorNotSupported) {
					return nil, NewTypeErrorMapValueInvalid(l[1])
				}
				return nil, err
			}
			return &Map{Any: &Any{}, KeyT: kT, ValueT: vT}, nil
		}
	}
	return nil, NewTypeErrorMapInvalid(typeStr)
}

func (m *Map) ParseString(str string) (any, error) {
	if !(str[0] == '{' && str[len(str)-1] == '}') {
		return nil, NewTypeErrorParseFailed(m, str)
	}
	v, err := ParseString(str)
	if err != nil {
		return nil, NewTypeErrorParseFailed(m, str)
	}
	if m.KeyT != nil && m.ValueT != nil {
		keyCheckFunc := m.KeyT.CheckFunc()
		valueCheckFunc := m.ValueT.CheckFunc()
		for key, val := range v.(map[any]any) {
			if !keyCheckFunc(key) {
				return nil, NewTypeErrorMapKeyNotMatch(m, key)
			}
			if !valueCheckFunc(val) {
				return nil, NewTypeErrorMapValueNotMatch(m, val)
			}
		}
	}
	return v, nil
}

func (m *Map) Convert(val any) string {
	var strList []string
	for k, v := range val.(map[any]any) {
		strList = append(strList, fmt.Sprintf("%v:%v", k, v))
	}
	return fmt.Sprintf("{%s}", strings.Join(strList, ","))
}

func (m *Map) String() string {
	return "map"
}

func (m *Map) DefaultValue() any {
	return map[any]any{}
}

func (m *Map) Kind() reflect.Kind {
	return reflect.Map
}

func (m *Map) CheckFunc() func(any) bool {
	keyCheckFunc := m.KeyT.CheckFunc()
	valueCheckFunc := m.ValueT.CheckFunc()
	return func(v any) bool {
		if v1, ok := v.(map[any]any); !ok {
			return false
		} else {
			for key, val := range v1 {
				if !keyCheckFunc(key) || !valueCheckFunc(val) {
					return false
				}
			}
		}
		return true
	}
}
