package fb_type

import (
	"xCelFlow/entities"
)

type FBInteger struct {
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
	return &FBInteger{integer.(*entities.Integer)}, nil
}

func (i *FBInteger) String() string {
	switch i.BitSize {
	case 8:
		return "int8"
	case 16:
		return "int16"
	case 32:
		return "int32"
	case 64:
		return "float64"
	default:
		return "float64"
	}
}
