package erl_type

import (
	"fmt"
	"xCelFlow/core"
	"xCelFlow/types"
)

type ErlStr struct {
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
	return &ErlStr{Str: s.(*types.Str)}, nil
}

func (s *ErlStr) Convert(val any) string {
	return fmt.Sprintf("<<\"%s\"/utf8>>", val)
}

func (s *ErlStr) String() string {
	return "binary()"
}
