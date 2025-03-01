package fb_type

import (
	"xCelFlow/entities"
)

type FBBoolean struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("bool", newBoolean)
}

func newBoolean(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	boolean, err := entities.NewBoolean(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &FBBoolean{ITypeSystem: boolean}, nil
}

func (b *FBBoolean) String() string {
	return "bool"
}
