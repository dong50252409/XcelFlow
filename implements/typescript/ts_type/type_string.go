package ts_type

import (
	"fmt"
	"xCelFlow/entities"
)

type TSStr struct {
	*entities.Str
}

func init() {
	typeRegister("str", newStr)
}

func newStr(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	s, err := entities.NewStr(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &TSStr{Str: s.(*entities.Str)}, nil
}

func (s *TSStr) Convert(val any) string {
	return fmt.Sprintf("%s", val)
}

func (s *TSStr) String() string {
	return "string | null"
}

func (*TSStr) DecoratorStr() string {
	return "@cacheStrRes()"
}

func (*TSStr) IsReferenceType() bool {
	return false
}

func (s *TSStr) MethodStr() string {
	return "__string"
}
