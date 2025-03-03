package ts_type

import (
	"xCelFlow/entities"
)

type TSBoolean struct {
	*entities.Boolean
}

func init() {
	typeRegister("bool", newBoolean)
}

func newBoolean(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	boolean, err := entities.NewBoolean(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &TSBoolean{Boolean: boolean.(*entities.Boolean)}, nil
}

func (b *TSBoolean) String() string {
	return "boolean"
}

func (*TSBoolean) DecoratorStr() string { return "" }

func (*TSBoolean) IsReferenceType() bool {
	return false
}

func (*TSBoolean) MethodStr() string {
	return "readInt8"
}
