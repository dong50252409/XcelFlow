package ts_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type TSFloat struct {
	*types.Float
}

func init() {
	typeRegister("float", newFloat)
}

func newFloat(typeStr string) (core.IType, error) {
	float, err := types.NewFloat(typeStr)
	if err != nil {
		return nil, err
	}
	return &TSFloat{Float: float.(*types.Float)}, nil
}

func (f *TSFloat) String() string {
	return "number"
}

func (*TSFloat) DecoratorStr() string { return "" }

func (*TSFloat) IsReferenceType() bool {
	return false
}

func (f *TSFloat) MethodStr() string {
	if f.BitSize == 32 {
		return "readFloat32"
	}
	return "readFloat64"
}
