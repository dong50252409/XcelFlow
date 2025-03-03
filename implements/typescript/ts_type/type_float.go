package ts_type

import (
	"xCelFlow/entities"
)

type TSFloat struct {
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
	return &TSFloat{Float: float.(*entities.Float)}, nil
}

func (f *TSFloat) String() string {
	return "number"
}

func (*TSFloat) DecoratorStr() string { return "" }

func (*TSFloat) IsReferenceType() bool {
	return false
}

func (f *TSFloat) MethodStr() string {
	if f.BitSize == 32 {
		return "readFloat32"
	}
	return "readFloat64"
}
