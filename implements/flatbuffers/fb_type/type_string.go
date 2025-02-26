package fb_type

import (
	"cfg_exporter/entities"
)

type FBStr struct {
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
	return &FBStr{ITypeSystem: s}, nil
}

func (s *FBStr) String() string {
	return "string"
}
