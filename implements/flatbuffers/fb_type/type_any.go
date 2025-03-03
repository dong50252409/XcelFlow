package fb_type

import (
	"xCelFlow/entities"
)

type FBAny struct {
	*entities.Any
}

func init() {
	typeRegister("any", newAny)
}

func newAny(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	anyValue, err := entities.NewAny(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &FBAny{anyValue.(*entities.Any)}, nil
}

func (s *FBAny) String() string {
	return "[ubyte](flexbuffer)"
}
