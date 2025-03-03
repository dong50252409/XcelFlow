package erl_type

import (
	"xCelFlow/entities"
)

type ErlFloat struct {
	*entities.Float
}

func init() {
	typeRegister("float", newFloat)
}

func newFloat(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	float, err := entities.NewFloat(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &ErlFloat{float.(*entities.Float)}, nil
}

func (f *ErlFloat) String() string {
	return "float()"
}
