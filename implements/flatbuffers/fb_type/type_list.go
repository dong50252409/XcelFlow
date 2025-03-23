package fb_type

import (
	"fmt"
	"xCelFlow/core"
	"xCelFlow/types"
)

type FBList struct {
	*types.List
}

func init() {
	typeRegister("list", newList)
}

func newList(typeStr string) (core.IType, error) {
	list, err := types.NewList(typeStr)
	if err != nil {
		return nil, err
	}
	return &FBList{List: list.(*types.List)}, nil
}

func (l *FBList) String() string {
	t := l.T
	switch l.T.(type) {
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
