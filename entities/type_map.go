package entities

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"xCelFlow/util"
)

type Map struct {
	Field  *Field
	KeyT   ITypeSystem
	ValueT ITypeSystem
}

func init() {
	TypeRegister("map", NewMap)
}

func NewMap(typeStr string, field *Field) (ITypeSystem, error) {
	if param := util.SubParam(typeStr); param == "" {
		return &Map{Field: field}, nil
	} else {
		if l := strings.SplitN(param, ",", 2); len(l) == 2 {
			kT, err := NewType(l[0], field)
			if err != nil {
				if errors.Is(err, TypeErrorNotSupported) {
					return nil, NewTypeErrorMapKeyInvalid(l[0])
				}
				return nil, err
			}

			vT, err := NewType(l[1], field)
			if err != nil {
				if errors.Is(err, TypeErrorNotSupported) {
					return nil, NewTypeErrorMapValueInvalid(l[1])
				}
				return nil, err
			}
			return &Map{Field: field, KeyT: kT, ValueT: vT}, nil
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

func (*Map) Convert(val any) string {
	var strList []string
	for k, v := range val.(map[any]any) {
		strList = append(strList, fmt.Sprintf("%v:%v", k, v))
	}
	return fmt.Sprintf("{%s}", strings.Join(strList, ","))
}

func (m *Map) String() string {
	return "map"
}

func (m *Map) DefaultValueStr() string {
	return "{}"
}

func (m *Map) Kind() reflect.Kind {
	return reflect.Map
}

func (m *Map) CheckFunc() func(any) bool {
	keyCF := m.KeyT.CheckFunc()
	valueCF := m.ValueT.CheckFunc()
	return func(v any) bool {
		v1, ok := v.(map[any]any)
		if !ok {
			return false
		}
		for key, val := range v1 {
			if !keyCF(key) || !valueCF(val) {
				return false
			}
		}
		return true
	}
}

func (m *Map) DecoratorStr() string {
	return ""
}

func (m *Map) IsReferenceType() bool {
	return true
}

func (m *Map) MethodStr() string {
	return ""
}
