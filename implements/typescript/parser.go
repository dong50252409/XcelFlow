package typescript

import (
	"xCelFlow/implements/typescript/ts_type"
	"xCelFlow/parser"
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
