package fb_type

import (
	"fmt"
	"xCelFlow/entities"
)

type FBList struct {
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
	return &FBList{list.(*entities.List)}, nil
}

func (l *FBList) String() string {
	t := l.T
	switch l.T.(type) {
	case *FBInteger:
		return fmt.Sprintf("[%s]", t.(*FBInteger).String())
	case *FBFloat:
		return fmt.Sprintf("[%s]", t.(*FBFloat).String())
	case *FBBoolean:
		return fmt.Sprintf("[%s]", t.(*FBBoolean).String())
	case *FBStr:
		return fmt.Sprintf("[%s]", t.(*FBStr).String())
	case *FBLang:
		return fmt.Sprintf("[%s]", t.(*FBLang).String())
	case *FBAny:
		return fmt.Sprintf("[%s]", t.(*FBAny).String())
	default:
		return "[ubyte](flexbuffer)"
	}
}

func (*FBList) DefaultValueStr() string {
	return "[]"
}
