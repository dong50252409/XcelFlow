package ts_type

import (
	"xCelFlow/entities"
)

type TSTuple struct {
	*entities.Tuple
}

func init() {
	typeRegister("tuple", newTuple)
}

func newTuple(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	tuple, err := entities.NewTuple(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &TSTuple{Tuple: tuple.(*entities.Tuple)}, nil
}

func (*TSTuple) Convert(val any) string {
	return toString(val)
}

func (*TSTuple) String() string {
	return "any[] | null"
}

func (*TSTuple) DefaultValueStr() string {
	return "[]"
}

func (*TSTuple) DecoratorStr() string {
	return "@cacheObjRes()"
}

func (*TSTuple) IsReferenceType() bool {
	return true
}
