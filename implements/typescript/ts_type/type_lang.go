package ts_type

import (
	"fmt"
	"xCelFlow/entities"
)

type TSLang struct {
	entities.ITypeSystem
}

func init() {
	typeRegister("lang", newLang)
}

func newLang(typeStr string, field *entities.Field) (entities.ITypeSystem, error) {
	lang, err := entities.NewLang(typeStr, field)
	if err != nil {
		return nil, err
	}
	return &TSLang{ITypeSystem: lang}, nil
}

func (l *TSLang) Convert(val any) string {
	return fmt.Sprintf("%s", val)
}

func (l *TSLang) String() string {
	return "string"
}

func (*TSLang) Decorator() string {
	return "@cacheStrRes()"
}
