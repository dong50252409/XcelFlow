package ts_type

import (
	"cfg_exporter/entities"
	"fmt"
)

type TSStr struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("str", newStr)
}

func newStr(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	s, err := entities.NewStr(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &TSStr{ITypeSystem: s}, nil
}

func (s *TSStr) Convert(val any) string {
	return fmt.Sprintf("%s", val)
}

func (s *TSStr) String() string {
	return "string"
}

func (*TSStr) Decorator() string {
	return "@cacheStrRes()"
}
