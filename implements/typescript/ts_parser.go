package typescript

import (
	"xCelFlow/core"
	"xCelFlow/implements/typescript/ts_type"
	"xCelFlow/parser"
)

type TSParse struct {
	*parser.FormParser
}

func init() {
	parser.RegisterParser("typescript", NewParser)
}

func NewParser(p *parser.FormParser) core.IParser {
	p.NewTypeFunc = ts_type.NewType
	return &TSParse{p}
}
