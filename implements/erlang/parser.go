package erlang

import (
	"cfg_exporter/implements/erlang/erl_type"
	"cfg_exporter/parser"
)

type ErlParser struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("erlang", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	p.NewTypeFunc = erl_type.NewType
	return &ErlParser{p}
}
