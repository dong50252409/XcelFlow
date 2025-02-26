package erl_type

import (
	"cfg_exporter/entities"
	"fmt"
)

type ErlStr struct {
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
	return &ErlStr{ITypeSystem: s}, nil
}

func (s *ErlStr) Convert(val any) string {
	return fmt.Sprintf("<<\"%s\"/utf8>>", val)
}

func (s *ErlStr) String() string {
	return "binary()"
}

func (s *ErlStr) DefaultValue() string {
	return "<<>>"
}
