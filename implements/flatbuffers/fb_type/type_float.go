package fb_type

import (
	"xCelFlow/entities"
)

type FBFloat struct {
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
	return &FBFloat{float.(*entities.Float)}, nil
}
