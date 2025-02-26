package fb_type

import (
	"cfg_exporter/entities"
)

type FBMap struct {
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
	return &FBMap{ITypeSystem: mapType}, nil
}

func (m *FBMap) String() string {
	return "[ubyte](flexbuffer)"
}

func (*FBMap) DefaultValue() string {
	return "[]"
}
