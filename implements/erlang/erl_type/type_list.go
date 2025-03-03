package erl_type

import (
	"xCelFlow/entities"
)

type ErlList struct {
	*entities.List
}

func init() {
	typeRegister("list", newList)
}

func newList(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	list, err := entities.NewList(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &ErlList{list.(*entities.List)}, nil
}

func (l *ErlList) Convert(val any) string {
	return toString(val)
}

func (l *ErlList) String() string {
	return "list()"
}
