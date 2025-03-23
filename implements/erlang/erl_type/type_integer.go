package erl_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type ErlInteger struct {
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
	return &ErlInteger{Integer: integer.(*types.Integer)}, nil
}

func (i *ErlInteger) String() string {
	return "integer()"
}
