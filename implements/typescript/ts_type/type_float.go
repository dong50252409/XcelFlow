package ts_type

import (
	"cfg_exporter/entities"
)

type TSFloat struct {
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
	return &TSFloat{ITypeSystem: float}, nil
}

func (f *TSFloat) String() string {
	return "number"
}

func (*TSFloat) Decorator() string { return "" }
