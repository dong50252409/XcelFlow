package entities

import (
	"errors"
	"fmt"
	"strings"
)

// TypeErrorNotSupported 类型不支持
var TypeErrorNotSupported error

func NewTypeErrorNotSupported(typeStr string) error {
	TypeErrorNotSupported = errors.New(fmt.Sprintf("%s 不支持的类型", typeStr))
	return TypeErrorNotSupported
}

// NewTypeErrorBaseInvalid 类型格式错误
func NewTypeErrorBaseInvalid(t ITypeSystem, lst []string) error {
	return fmt.Errorf("格式错误 %s|%s(%s)", t.String(), t.String(), strings.Join(lst, "|"))
}

// NewTypeErrorListInvalid 类型格式错误
func NewTypeErrorListInvalid(typeStr string) error {
	return fmt.Errorf("格式错误 list%s", typeStr)
}

// NewTypeErrorTupleInvalid 元组类型格式错误
func NewTypeErrorTupleInvalid(typeStr string) error {
	return fmt.Errorf("格式错误 tuple%s", typeStr)
}

// NewTypeErrorMapKeyInvalid 字典类型格式错误
func NewTypeErrorMapKeyInvalid(typeStr string) error {
	return fmt.Errorf("格式错误 map(%s, ...)", typeStr)
}

// NewTypeErrorMapValueInvalid 字典类型格式错误
func NewTypeErrorMapValueInvalid(typeStr string) error {
	return fmt.Errorf("键格式错误 map(..., %s)", typeStr)
}

// NewTypeErrorMapInvalid 字典类型格式错误
func NewTypeErrorMapInvalid(typeStr string) error {
	return fmt.Errorf("值格式错误 %s", typeStr)
}

// NewTypeErrorParseFailed 类型解析失败
func NewTypeErrorParseFailed(parentType ITypeSystem, str string) error {
	return fmt.Errorf("%s 类型无法解析 %s", parentType.String(), str)
}

// NewTypeErrorNotMatch 类型匹配
func NewTypeErrorNotMatch(t ITypeSystem, index int, element any) error {
	return fmt.Errorf("第 %d 个元素 %v 与类型 %s 不匹配", index+1, element, formatTypeSign(t))
}

// NewTypeErrorMapKeyNotMatch 键类型匹配
func NewTypeErrorMapKeyNotMatch(t ITypeSystem, key any) error {
	return fmt.Errorf("键:%v 与类型 %s 不匹配", key, formatTypeSign(t))
}

// NewTypeErrorMapValueNotMatch 值类型匹配
func NewTypeErrorMapValueNotMatch(t ITypeSystem, val any) error {
	return fmt.Errorf("值:%v 与类型 %s 不匹配", val, formatTypeSign(t))
}

func formatTypeSign(t ITypeSystem) string {
	switch t.(type) {
	case *List:
		return fmt.Sprintf("%s(%s)", t.String(), formatTypeSign(t.(*List).T))
	case *Tuple:
		return fmt.Sprintf("%s(%s)", t.String(), formatTypeSign(t.(*Tuple).T))
	case *Map:
		return fmt.Sprintf("%s(%s,%s)", t.String(), formatTypeSign(t.(*Map).KeyT), formatTypeSign(t.(*Map).ValueT))
	default:
		return t.String()
	}
}
