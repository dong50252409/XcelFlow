package ts_type

import (
	"cfg_exporter/entities"
)

type TSInteger struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("int", newInteger)
}

func newInteger(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	integer, err := entities.NewInteger(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &TSInteger{ITypeSystem: integer}, nil
}

func (i *TSInteger) String() string {
	return "number"
}

func (*TSInteger) Decorator() string { return "" }
