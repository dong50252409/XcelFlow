package erl_type

import (
	"xCelFlow/entities"
)

type ErlInteger struct {
	*entities.Integer
}

func init() {
	typeRegister("int", newInteger)
}

func newInteger(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	integer, err := entities.NewInteger(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &ErlInteger{integer.(*entities.Integer)}, nil
}

func (i *ErlInteger) String() string {
	return "integer()"
}
