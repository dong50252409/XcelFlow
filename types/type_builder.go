package types

import (
	"xCelFlow/core"
	"xCelFlow/util"
)

var (
	typeRegistry = make(map[string]func(typeStr string) (core.IType, error))
)

// TypeRegister 类型注册器
func TypeRegister(key string, cls func(typeStr string) (core.IType, error)) {
	typeRegistry[key] = cls
}

func NewType(typeStr string) (core.IType, error) {
	key, args := util.GetKey(typeStr)
	if cls, ok := typeRegistry[key]; ok {
		return cls(args)
	}
	return nil, NewTypeErrorNotSupported(key)
}
