package flatbuffers

import (
	"xCelFlow/core"
	"xCelFlow/implements/flatbuffers/fb_type"
	"xCelFlow/parser"
)

type FBSParse struct {
	*parser.FormParser
}

func init() {
	parser.RegisterParser("flatbuffers", NewParser)
}

func NewParser(p *parser.FormParser) core.IParser {
	p.NewTypeFunc = fb_type.NewType
	return &FBSParse{p}
}
