package json

import (
	"xCelFlow/core"
	"xCelFlow/parser"
)

type JSONParser struct {
	*parser.FormParser
}

func init() {
	parser.RegisterParser("json", NewParser)
}

func NewParser(p *parser.FormParser) core.IParser {
	return &JSONParser{p}
}
