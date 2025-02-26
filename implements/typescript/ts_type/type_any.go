package ts_type

import (
	"cfg_exporter/entities"
	"fmt"
)

type TSAny struct {
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
	return &TSAny{ITypeSystem: anyValue}, nil
}

func (s *TSAny) Convert(val any) string {
	return fmt.Sprintf(`"%s"`, val)
}

func (s *TSAny) String() string {
	return "any"
}

func (*TSAny) Decorator() string {
	return "@cacheObjRes()"
}
