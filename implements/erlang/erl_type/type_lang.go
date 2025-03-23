package erl_type

import (
	"fmt"
	"xCelFlow/core"
	"xCelFlow/types"
)

type ErlLang struct {
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
	return &ErlLang{Lang: lang.(*types.Lang)}, nil
}

func (l *ErlLang) Convert(val any) string {
	return fmt.Sprintf("<<\"%s\"/utf8>>", val)
}

func (l *ErlLang) String() string {
	return "binary()"
}
