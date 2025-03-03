package fb_type

import (
	"xCelFlow/entities"
)

type FBStr struct {
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
	return &FBStr{s.(*entities.Str)}, nil
}

func (s *FBStr) String() string {
	return "string"
}
