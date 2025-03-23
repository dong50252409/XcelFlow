package fb_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type FBFloat struct {
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
	return &FBFloat{Float: float.(*types.Float)}, nil
}
