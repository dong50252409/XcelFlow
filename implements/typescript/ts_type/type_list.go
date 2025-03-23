package ts_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type TSList struct {
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
	return &TSList{List: list.(*types.List)}, nil
}

func (l *TSList) Convert(val any) string {
	return toString(val)
}

func (l *TSList) String() string {
	return "any[] | null"
}

func (*TSList) DecoratorStr() string {
	return "@cacheObjRes()"
}

func (*TSList) IsReferenceType() bool {
	return true
}
