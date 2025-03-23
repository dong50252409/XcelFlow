package fb_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type FBAny struct {
	*types.Any
}

func init() {
	typeRegister("any", newAny)
}

func newAny(typeStr string) (core.IType, error) {
	anyValue, err := types.NewAny(typeStr)
	if err != nil {
		return nil, err
	}
	return &FBAny{Any: anyValue.(*types.Any)}, nil
}

func (s *FBAny) String() string {
	return "[ubyte](flexbuffer)"
}
