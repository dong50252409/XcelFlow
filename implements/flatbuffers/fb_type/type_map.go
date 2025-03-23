package fb_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type FBMap struct {
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
	return &FBMap{Map: mapType.(*types.Map)}, nil
}

func (m *FBMap) String() string {
	return "[ubyte](flexbuffer)"
}
