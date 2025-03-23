package ts_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type TSTuple struct {
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
	return &TSTuple{Tuple: tuple.(*types.Tuple)}, nil
}

func (*TSTuple) Convert(val any) string {
	return toString(val)
}

func (*TSTuple) String() string {
	return "any[] | null"
}

func (*TSTuple) DecoratorStr() string {
	return "@cacheObjRes()"
}

func (*TSTuple) IsReferenceType() bool {
	return true
}
