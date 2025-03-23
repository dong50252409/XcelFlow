package ts_type

import (
	"fmt"
	"xCelFlow/core"
	"xCelFlow/types"
)

type TSAny struct {
	*types.Any
}

func init() {
	typeRegister("any", newAny)
}

func newAny(typeStr string) (core.IType, error) {
	anyValue, err := types.NewAny(typeStr)
	if err != nil {
		return nil, err
	}
	return &TSAny{Any: anyValue.(*types.Any)}, nil
}

func (s *TSAny) Convert(val any) string {
	return fmt.Sprintf(`"%s"`, val)
}

func (s *TSAny) String() string {
	return "any | null"
}

func (*TSAny) DecoratorStr() string {
	return "@cacheObjRes()"
}

func (*TSAny) IsReferenceType() bool {
	return true
}
