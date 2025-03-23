package fb_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type FBLang struct {
	*types.Lang
}

func init() {
	typeRegister("lang", newLang)
}

func newLang(typeStr string) (core.IType, error) {
	lang, err := types.NewLang(typeStr)
	if err != nil {
		return nil, err
	}
	return &FBLang{Lang: lang.(*types.Lang)}, nil
}

func (l *FBLang) String() string {
	return "string"
}
