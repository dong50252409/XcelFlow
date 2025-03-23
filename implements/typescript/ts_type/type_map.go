package ts_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type TSMap struct {
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
	return &TSMap{Map: mapType.(*types.Map)}, nil
}

func (*TSMap) Convert(val any) string {
	return toString(val)
}

func (m *TSMap) String() string {
	return "Map<any, any> | null"
}

func (*TSMap) DecoratorStr() string {
	return "@cacheObjRes()"
}

func (*TSMap) IsReferenceType() bool {
	return true
}
