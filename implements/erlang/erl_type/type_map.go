package erl_type

import (
	"xCelFlow/entities"
)

type ErlMap struct {
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
	return &ErlMap{mapType.(*entities.Map)}, nil
}

func (*ErlMap) Convert(val any) string {
	return toString(val)
}

func (m *ErlMap) String() string {
	return "map()"
}

func (*ErlMap) DefaultValueStr() string {
	return "#{}"
}
