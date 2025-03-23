package fb_type

import (
	"xCelFlow/core"
	"xCelFlow/types"
)

type FBStr struct {
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
	return &FBStr{Str: s.(*types.Str)}, nil
}

func (s *FBStr) String() string {
	return "string"
}
