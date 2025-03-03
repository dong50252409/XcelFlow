package fb_type

import (
	"xCelFlow/entities"
)

type FBMap struct {
	*entities.Map
}

func init() {
	typeRegister("map", newMap)
}

func newMap(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	mapType, err := entities.NewMap(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &FBMap{mapType.(*entities.Map)}, nil
}

func (m *FBMap) String() string {
	return "[ubyte](flexbuffer)"
}

func (*FBMap) DefaultValueStr() string {
	return "[]"
}
