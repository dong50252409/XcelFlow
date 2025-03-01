package fb_type

import (
	"fmt"
	"xCelFlow/entities"
)

type FBTuple struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("tuple", newTuple)
}

func newTuple(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	tuple, err := entities.NewTuple(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &FBTuple{ITypeSystem: tuple}, nil
}

func (fbt *FBTuple) String() string {
	t := fbt.ITypeSystem.(*entities.Tuple).T
	switch t.(type) {
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

func (*FBTuple) DefaultValue() string {
	return "[]"
}
