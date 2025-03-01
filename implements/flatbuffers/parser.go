package flatbuffers

import (
	"xCelFlow/implements/flatbuffers/fb_type"
	"xCelFlow/parser"
)

type FBSParse struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("flatbuffers", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	p.NewTypeFunc = fb_type.NewType
	return &FBSParse{p}
}
