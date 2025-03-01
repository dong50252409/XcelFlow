package erl_type

import (
	"xCelFlow/entities"
)

type ErlBoolean struct {
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
	return &ErlBoolean{ITypeSystem: boolean}, nil
}

func (b *ErlBoolean) String() string {
	return "boolean()"
}
