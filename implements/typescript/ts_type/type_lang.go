package ts_type

import (
	"fmt"
	"xCelFlow/core"
	"xCelFlow/types"
)

type TSLang struct {
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
	return &TSLang{Lang: lang.(*types.Lang)}, nil
}

func (l *TSLang) Convert(val any) string {
	return fmt.Sprintf("%s", val)
}

func (l *TSLang) String() string {
	return "string | Uint8Array"
}

func (*TSLang) DecoratorStr() string {
	return "@cacheStrRes()"
}

func (*TSLang) IsReferenceType() bool {
	return false
}

func (l *TSLang) MethodStr() string {
	return "__string"
}
