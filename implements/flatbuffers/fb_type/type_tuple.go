package fb_type

import (
	"fmt"
	"xCelFlow/core"
	"xCelFlow/types"
)

type FBTuple struct {
	*types.Tuple
}

func init() {
	typeRegister("tuple", newTuple)
}

func newTuple(typeStr string) (core.IType, error) {
	tuple, err := types.NewTuple(typeStr)
	if err != nil {
		return nil, err
	}
	return &FBTuple{Tuple: tuple.(*types.Tuple)}, nil
}

func (fbt *FBTuple) String() string {
	t := fbt.T
	switch t.(type) {
	case *FBInteger:
		return fmt.Sprintf("[%s]", t.String())
	case *FBFloat:
		return fmt.Sprintf("[%s]", t.String())
	case *FBBoolean:
		return fmt.Sprintf("[%s]", t.String())
	case *FBStr:
		return fmt.Sprintf("[%s]", t.String())
	case *FBLang:
		return fmt.Sprintf("[%s]", t.String())
	case *FBAny:
		return fmt.Sprintf("[%s]", t.String())
	default:
		return "[ubyte](flexbuffer)"
	}
}
