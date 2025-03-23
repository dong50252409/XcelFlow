package erl_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type ErlList struct {
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
	return &ErlList{List: list.(*types.List)}, nil
}

func (l *ErlList) Convert(val any) string {
	return toString(val)
}

func (l *ErlList) String() string {
	return "list()"
}
