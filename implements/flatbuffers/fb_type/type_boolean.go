package fb_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type FBBoolean struct {
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
	return &FBBoolean{Boolean: boolean.(*types.Boolean)}, nil
}

func (b *FBBoolean) String() string {
	return "bool"
}
