package erl_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type ErlFloat struct {
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
	return &ErlFloat{Float: float.(*types.Float)}, nil
}

func (f *ErlFloat) String() string {
	return "float()"
}
