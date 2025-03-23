package ts_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type TSBoolean struct {
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
	return &TSBoolean{Boolean: boolean.(*types.Boolean)}, nil
}

func (b *TSBoolean) String() string {
	return "boolean"
}

func (*TSBoolean) DecoratorStr() string { return "" }

func (*TSBoolean) IsReferenceType() bool {
	return false
}

func (*TSBoolean) MethodStr() string {
	return "readInt8"
}
