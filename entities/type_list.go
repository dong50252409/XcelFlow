package entities

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"xCelFlow/util"
)

type List struct {
	Field *Field
	T     ITypeSystem
}

func init() {
	TypeRegister("list", NewList)
}

func NewList(typeStr string, field *Field) (ITypeSystem, error) {
	if param := util.SubParam(typeStr); param == "" {
		return &List{Field: field}, nil
	} else {
		t, err := NewType(param, field)
		if err != nil {
			if errors.Is(err, TypeErrorNotSupported) {
				return nil, NewTypeErrorListInvalid(typeStr)
			}
			return nil, err
		}
		return &List{Field: field, T: t}, nil
	}
}

func (l *List) ParseString(str string) (any, error) {
	if !(str[0] == '[' && str[len(str)-1] == ']') {
		return nil, NewTypeErrorParseFailed(l, str)
	}
	v, err := ParseString(str)
	if err != nil {
		return nil, NewTypeErrorParseFailed(l, str)
	}
	if l.T != nil {
		checkFunc := l.T.CheckFunc()
		for i, e := range v.([]any) {
			if !checkFunc(e) {
				return nil, NewTypeErrorNotMatch(l, i, e)
			}
		}
	}
	return v, nil
}

func (l *List) Convert(val any) string {
	var strList []string
	for _, e := range val.([]any) {
		strList = append(strList, fmt.Sprintf("%v", e))
	}
	return fmt.Sprintf("[%s]", strings.Join(strList, ","))
}

func (l *List) String() string {
	return "list"
}

func (l *List) DefaultValueStr() string {
	return "[]"
}

func (l *List) Kind() reflect.Kind {
	return reflect.Slice
}

func (l *List) CheckFunc() func(any) bool {
	cf := l.T.CheckFunc()
	return func(v any) bool {
		v1, ok := v.([]any)
		if !ok {
			return false
		}
		for _, e := range v1 {
			if !cf(e) {
				return false
			}
		}
		return true
	}
}

func (l *List) DecoratorStr() string {
	return ""
}

func (l *List) IsReferenceType() bool {
	return true
}

func (l *List) MethodStr() string {
	return ""
}
