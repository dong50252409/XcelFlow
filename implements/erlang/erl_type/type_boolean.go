package erl_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type ErlBoolean struct {
	*types.Boolean
}

func init() {
	typeRegister("bool", newBoolean)
}

func newBoolean(typeStr string) (core.IType, error) {
	boolean, err := types.NewBoolean(typeStr)
	if err != nil {
		return nil, err
	}
	return &ErlBoolean{Boolean: boolean.(*types.Boolean)}, nil
}

func (b *ErlBoolean) String() string {
	return "boolean()"
}
