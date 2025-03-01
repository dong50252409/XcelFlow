package ts_type

import (
	"xCelFlow/entities"
)

type TSTuple struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("tuple", newTuple)
}

func newTuple(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	tuple, err := entities.NewTuple(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &TSTuple{ITypeSystem: tuple}, nil
}

func (*TSTuple) Convert(val any) string {
	return toString(val)
}

func (*TSTuple) String() string {
	return "any"
}

func (*TSTuple) DefaultValue() string {
	return "[]"
}

func (*TSTuple) Decorator() string {
	return "@cacheObjRes()"
}
