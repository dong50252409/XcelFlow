package json

import (
	"xCelFlow/parser"
)

type JSONParser struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("json", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	return &JSONParser{p}
}
