package erl_type

import (
	"xCelFlow/entities"
)

type ErlTuple struct {
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
	return &ErlTuple{tuple.(*entities.Tuple)}, nil
}

func (*ErlTuple) Convert(val any) string {
	return toString(val)
}

func (*ErlTuple) String() string {
	return "tuple()"
}

func (*ErlTuple) DefaultValueStr() string {
	return "{}"
}
