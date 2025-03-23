package erl_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type ErlMap struct {
	*types.Map
}

func init() {
	typeRegister("map", newMap)
}

func newMap(typeStr string) (core.IType, error) {
	mapType, err := types.NewMap(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlMap{Map: mapType.(*types.Map)}, nil
}

func (*ErlMap) Convert(val any) string {
	return toString(val)
}

func (m *ErlMap) String() string {
	return "map()"
}
