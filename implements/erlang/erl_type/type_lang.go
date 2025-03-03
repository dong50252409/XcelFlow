package erl_type

import (
	"fmt"
	"xCelFlow/entities"
)

type ErlLang struct {
	*entities.Lang
}

func init() {
	typeRegister("lang", newLang)
}

func newLang(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	lang, err := entities.NewLang(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &ErlLang{lang.(*entities.Lang)}, nil
}

func (l *ErlLang) Convert(val any) string {
	return fmt.Sprintf("<<\"%s\"/utf8>>", val)
}

func (l *ErlLang) String() string {
	return "binary()"
}

func (l *ErlLang) DefaultValueStr() string {
	return "<<>>"
}
