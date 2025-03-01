package erlang

import (
	"xCelFlow/implements/erlang/erl_type"
	"xCelFlow/parser"
)

type ERLParser struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("erlang", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	p.NewTypeFunc = erl_type.NewType
	return &ERLParser{p}
}
