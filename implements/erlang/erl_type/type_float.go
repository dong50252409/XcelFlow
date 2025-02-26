package erl_type

import (
	"cfg_exporter/entities"
)

type ErlFloat struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("float", newFloat)
}

func newFloat(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	float, err := entities.NewFloat(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &ErlFloat{ITypeSystem: float}, nil
}

func (f *ErlFloat) String() string {
	return "float()"
}
