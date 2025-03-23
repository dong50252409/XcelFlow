package json

import (
	"fmt"
	"strconv"
	"xCelFlow/core"
)

func convert(data any) any {
	switch v := data.(type) {
	case []any:
		var elements = make([]any, 0, len(v))
		for _, item := range v {
			if item == nil {
				break
			}
			elements = append(elements, convert(item))
		}
		return elements
	case core.TupleT:
		var elements = make([]any, 0)
		for _, item := range v {
			if item == nil {
				break
			}
			elements = append(elements, convert(item))
		}
		return elements
	case map[any]any:
		elements := make(map[string]any, len(v))
		for key, value := range v {
			elements[toString(key)] = convert(value)
		}
		return elements
	default:
		return v
	}
}

func toString(data any) string {
	switch v := data.(type) {
	case int8, int16, int32, int64, int:
		return strconv.FormatInt(v.(int64), 10)
	case float32, float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case string:
		return v
	case core.AnyT:
		return string(v)
	default:
		panic(fmt.Sprintf("%v 无法作为键来使用", v))
	}
}
