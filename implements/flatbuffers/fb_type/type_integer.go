package fb_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type FBInteger struct {
	*types.Integer
}

func init() {
	typeRegister("int", newInteger)
}

func newInteger(typeStr string) (core.IType, error) {
	integer, err := types.NewInteger(typeStr)
	if err != nil {
		return nil, err
	}
	return &FBInteger{Integer: integer.(*types.Integer)}, nil
}

func (i *FBInteger) String() string {
	switch i.BitSize {
	case 8:
		return "int8"
	case 16:
		return "int16"
	case 32:
		return "int32"
	case 64:
		return "float64" // TODO typescript类型问题导致不是int64，如果可以修改，可以改成int64
	default:
		return "float64"
	}
}
