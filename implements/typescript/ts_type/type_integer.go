package ts_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type TSInteger struct {
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
	return &TSInteger{Integer: integer.(*types.Integer)}, nil
}

func (i *TSInteger) String() string {
	return "number"
}

func (*TSInteger) DecoratorStr() string { return "" }

func (*TSInteger) IsReferenceType() bool {
	return false
}

func (i *TSInteger) MethodStr() string {
	switch i.BitSize {
	case 8:
		return "readInt8"
	case 16:
		return "readInt16"
	case 32:
		return "readInt32"
	case 64:
		return "readFloat64"
	}
	return "readFloat64"
}
