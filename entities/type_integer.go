package entities

import (
	"maps"
	"math"
	"reflect"
	"strconv"
	"xCelFlow/util"
)

type Integer struct {
	Field   *Field
	BitSize int
}

var (
	// 默认int位数
	integerDefaultBitSizes = "64"

	// int位数
	intByteSizes = map[string]int{
		"8":  8,
		"16": 16,
		"32": 32,
		"64": 64,
	}
)

func init() {
	TypeRegister("int", NewInteger)
}

func NewInteger(typeStr string, field *Field) (ITypeSystem, error) {
	bit := integerDefaultBitSizes
	if param := util.SubParam(typeStr); param != "" {
		bit = param
	}
	if bytes, ok := intByteSizes[bit]; ok {
		return &Integer{Field: field, BitSize: bytes}, nil
	}

	l := make([]string, 0, len(intByteSizes))
	for k := range maps.Keys(intByteSizes) {
		l = append(l, k)
	}
	return nil, NewTypeErrorBaseInvalid(&Integer{}, l)
}
func (i *Integer) ParseString(str string) (any, error) {
	parseInt, err := strconv.ParseInt(str, 10, i.BitSize)
	if err != nil {
		return nil, NewTypeErrorParseFailed(i, str)
	}
	return parseInt, nil
}

func (i *Integer) Convert(val any) string {
	return strconv.FormatInt(val.(int64), 10)
}

func (i *Integer) String() string {
	switch i.BitSize {
	case 8:
		return "int8"
	case 16:
		return "int16"
	case 32:
		return "int32"
	case 64:
		return "int64"
	default:
		return "int64"
	}
}

func (i *Integer) DefaultValueStr() string {
	return "0"
}

func (i *Integer) Kind() reflect.Kind {
	switch i.BitSize {
	case 8:
		return reflect.Int8
	case 16:
		return reflect.Int16
	case 32:
		return reflect.Int32
	case 64:
		return reflect.Int64
	default:
		return reflect.Int64
	}
}

func (i *Integer) CheckFunc() func(any) bool {
	switch i.BitSize {
	case 8:
		return func(v any) bool {
			v1, ok := v.(int64)
			return ok && math.MinInt8 <= v1 && v1 <= math.MaxInt8
		}
	case 16:
		return func(v any) bool {
			v1, ok := v.(int64)
			return ok && math.MinInt16 <= v1 && v1 <= math.MaxInt16
		}
	case 32:
		return func(v any) bool {
			v1, ok := v.(int64)
			return ok && math.MinInt32 <= v1 && v1 <= math.MaxInt32
		}
	case 64:
		return func(v any) bool { _, ok := v.(int64); return ok }
	default:
		return func(v any) bool { _, ok := v.(int64); return ok }
	}
}

func (i *Integer) DecoratorStr() string {
	return ""
}

func (i *Integer) IsReferenceType() bool {
	return false
}

func (i *Integer) MethodStr() string {
	return ""
}
