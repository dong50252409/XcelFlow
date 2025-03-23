package types

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"xCelFlow/core"
	"xCelFlow/util"
)

type List struct {
	*Any
	T core.IType
}

func init() {
	TypeRegister("list", NewList)
}

func NewList(typeStr string) (core.IType, error) {
	if param := util.SubParam(typeStr); param == "" {
		return &List{Any: &Any{}}, nil
	} else {
		t, err := NewType(param)
		if err != nil {
			if errors.Is(err, TypeErrorNotSupported) {
				return nil, NewTypeErrorListInvalid(typeStr)
			}
			return nil, err
		}
		return &List{Any: &Any{}, T: t}, nil
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

func (l *List) DefaultValue() any {
	return []any{}
}

func (l *List) Kind() reflect.Kind {
	return reflect.Slice
}

func (l *List) CheckFunc() func(any) bool {
	checkFunc := l.T.CheckFunc()
	return func(v any) bool {
		if v1, ok := v.([]any); !ok {
			return false
		} else {
			for _, e := range v1 {
				if !checkFunc(e) {
					return false
				}
			}
		}
		return true
	}
}
