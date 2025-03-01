package fb_type

import (
	"xCelFlow/entities"
	"xCelFlow/util"
)

var (
	typeRegistry = make(map[string]func(typeStr string, field *entities.Field) (entities.ITypeSystem, error))
)

// 类型注册器
func typeRegister(key string, cls func(typeStr string, field *entities.Field) (entities.ITypeSystem, error)) {
	typeRegistry[key] = cls
}

func NewType(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	key, args := util.GetKey(typeStr)
	if cls, ok := typeRegistry[key]; ok {
		return cls(args, field)
	}
	return nil, entities.NewTypeErrorNotSupported(key)
}
