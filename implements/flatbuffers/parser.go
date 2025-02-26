package flatbuffers

import (
	"cfg_exporter/implements/flatbuffers/fb_type"
	"cfg_exporter/parser"
)

type FBParse struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("flatbuffers", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	p.NewTypeFunc = fb_type.NewType
	return &FBParse{p}
}
