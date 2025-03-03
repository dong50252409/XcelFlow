package ts_type

import (
	"xCelFlow/entities"
)

type TSMap struct {
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
	return &TSMap{Map: mapType.(*entities.Map)}, nil
}

func (*TSMap) Convert(val any) string {
	return toString(val)
}

func (m *TSMap) String() string {
	return "Map<any, any> | null"
}

func (*TSMap) DefaultValueStr() string {
	return "new Map()"
}

func (*TSMap) DecoratorStr() string {
	return "@cacheObjRes()"
}

func (*TSMap) IsReferenceType() bool {
	return true
}
