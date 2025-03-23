package ts_type

import (
	"fmt"
	"xCelFlow/core"
	"xCelFlow/types"
)

type TSStr struct {
	*types.Str
}

func init() {
	typeRegister("str", newStr)
}

func newStr(typeStr string) (core.IType, error) {
	s, err := types.NewStr(typeStr)
	if err != nil {
		return nil, err
	}
	return &TSStr{Str: s.(*types.Str)}, nil
}

func (s *TSStr) Convert(val any) string {
	return fmt.Sprintf("%s", val)
}

func (s *TSStr) String() string {
	return "string | Uint8Array"
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
