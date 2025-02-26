package typescript

import (
	"cfg_exporter/implements/typescript/ts_type"
	"cfg_exporter/parser"
)

type TSParse struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("typescript", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	p.NewTypeFunc = ts_type.NewType
	return &TSParse{p}
}
