package entities

import (
	"xCelFlow/util"
)

var (
	typeRegistry = make(map[string]func(typeStr string, field *Field) (ITypeSystem, error))
)

// TypeRegister 类型注册器
func TypeRegister(key string, cls func(typeStr string, field *Field) (ITypeSystem, error)) {
	typeRegistry[key] = cls
}

func NewType(typeStr string, field *Field) (ITypeSystem, error) {
	key, args := util.GetKey(typeStr)
	if cls, ok := typeRegistry[key]; ok {
		return cls(args, field)
	}
	return nil, NewTypeErrorNotSupported(key)
}
