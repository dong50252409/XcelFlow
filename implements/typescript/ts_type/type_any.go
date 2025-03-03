package ts_type

import (
	"fmt"
	"xCelFlow/entities"
)

type TSAny struct {
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
	return &TSAny{Any: anyValue.(*entities.Any)}, nil
}

func (s *TSAny) Convert(val any) string {
	return fmt.Sprintf(`"%s"`, val)
}

func (s *TSAny) String() string {
	return "any | null"
}

func (*TSAny) DecoratorStr() string {
	return "@cacheObjRes()"
}

func (*TSAny) IsReferenceType() bool {
	return true
}
