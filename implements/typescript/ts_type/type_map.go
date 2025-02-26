package ts_type

import (
	"cfg_exporter/entities"
)

type TSMap struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("map", newMap)
}

func newMap(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	mapType, err := entities.NewMap(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &TSMap{ITypeSystem: mapType}, nil
}

func (*TSMap) Convert(val any) string {
	return toString(val)
}

func (m *TSMap) String() string {
	return "any"
}

func (*TSMap) DefaultValue() string {
	return "new Map()"
}

func (*TSMap) Decorator() string {
	return "@cacheObjRes()"
}
