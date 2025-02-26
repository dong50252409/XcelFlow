package fb_type

import (
	"cfg_exporter/entities"
)

type FBAny struct {
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
	return &FBAny{ITypeSystem: anyValue}, nil
}

func (s *FBAny) String() string {
	return "[ubyte](flexbuffer)"
}
