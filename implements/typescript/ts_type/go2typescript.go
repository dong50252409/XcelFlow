package ts_type

import (
	"fmt"
	"strconv"
	"strings"
	"xCelFlow/entities"
)

func toString(data interface{}) string {
	switch v := data.(type) {
	case int8, int16, int32, int64, int: // 处理整数
		return strconv.FormatInt(v.(int64), 10)
	case float32, float64: // 处理浮点数
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case string: // 处理字符串
		return fmt.Sprintf(`"%s"`, v)
	case bool: // 处理布尔值
		return strconv.FormatBool(v)
	case entities.AnyT: // 处理原始值，在Erlang中就是atom
		return fmt.Sprintf("%s", v)
	case []interface{}: // 处理列表
		var elements []string
		for _, item := range v {
			elements = append(elements, toString(item))
		}
		return "[" + strings.Join(elements, ", ") + "]"
	case map[interface{}]interface{}: // 处理Map
		var elements []string
		for key, value := range v {
			elements = append(elements, fmt.Sprintf("[%s, %s]", toString(key), toString(value)))
		}
		return "new Map(" + strings.Join(elements, ", ") + ")"
	case entities.TupleT: // 处理数组
		var elements []string
		for _, item := range v {
			if item == nil {
				break
			}
			elements = append(elements, toString(item))
		}
		return "[" + strings.Join(elements, ", ") + "]"
	default:
		return "undefined"
	}
}
