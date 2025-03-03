package ts_type

import (
	"xCelFlow/entities"
)

type TSList struct {
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
	return &TSList{List: list.(*entities.List)}, nil
}

func (l *TSList) Convert(val any) string {
	return toString(val)
}

func (l *TSList) String() string {
	return "any[] | null"
}

func (*TSList) DefaultValueStr() string {
	return "[]"
}

func (*TSList) DecoratorStr() string {
	return "@cacheObjRes()"
}

func (*TSList) IsReferenceType() bool {
	return true
}
