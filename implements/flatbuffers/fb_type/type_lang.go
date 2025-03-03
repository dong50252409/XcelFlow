package fb_type

import (
	"xCelFlow/entities"
)

type FBLang struct {
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
	return &FBLang{lang.(*entities.Lang)}, nil
}

func (l *FBLang) String() string {
	return "string"
}
