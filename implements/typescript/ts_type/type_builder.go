package ts_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
	"xCelFlow/util"
)

var (
	typeRegistry = make(map[string]func(typeStr string) (core.IType, error))
)

// 类型注册器
func typeRegister(key string, cls func(typeStr string) (core.IType, error)) {
	typeRegistry[key] = cls
}

func NewType(typeStr string) (core.IType, error) {
	key, args := util.GetKey(typeStr)
	if cls, ok := typeRegistry[key]; ok {
		return cls(args)
	}
	return nil, types.NewTypeErrorNotSupported(key)
}
