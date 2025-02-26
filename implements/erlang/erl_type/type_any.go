package erl_type

import (
	"cfg_exporter/entities"
)

type ErlAny struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("any", newAny)
}

func newAny(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	anyValue, err := entities.NewAny(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &ErlAny{ITypeSystem: anyValue}, nil
}

func (s *ErlAny) Convert(val any) string {
	v1 := toString(val)
	return v1
}

func (s *ErlAny) String() string {
	return "term()"
}

func (s *ErlAny) DefaultValue() string {
	return "undefined"
}
