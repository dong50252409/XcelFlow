package erl_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type ErlTuple struct {
	*types.Tuple
}

func init() {
	typeRegister("tuple", newTuple)
}

func newTuple(typeStr string) (core.IType, error) {
	tuple, err := types.NewTuple(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlTuple{Tuple: tuple.(*types.Tuple)}, nil
}

func (*ErlTuple) Convert(val any) string {
	return toString(val)
}

func (*ErlTuple) String() string {
	return "tuple()"
}
