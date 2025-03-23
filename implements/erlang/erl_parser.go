package erlang

import (
	"xCelFlow/core"
	"xCelFlow/implements/erlang/erl_type"
	"xCelFlow/parser"
)

type ERLParser struct {
	*parser.FormParser
}

func init() {
	parser.RegisterParser("erlang", NewParser)
}

func NewParser(p *parser.FormParser) core.IParser {
	p.NewTypeFunc = erl_type.NewType
	return &ERLParser{p}
}
