package ts_type

import (
	"xCelFlow/entities"
)

type TSBoolean struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("bool", newBoolean)
}

func newBoolean(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	boolean, err := entities.NewBoolean(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &TSBoolean{ITypeSystem: boolean}, nil
}

func (b *TSBoolean) String() string {
	return "boolean"
}

func (*TSBoolean) Decorator() string { return "" }
