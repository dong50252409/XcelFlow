package erl_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type ErlAny struct {
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
	return &ErlAny{Any: anyValue.(*types.Any)}, nil
}

func (s *ErlAny) Convert(val any) string {
	v1 := toString(val)
	return v1
}

func (s *ErlAny) String() string {
	return "term()"
}
