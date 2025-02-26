package json

import (
	"cfg_exporter/parser"
)

type JParser struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("json", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	return &JParser{p}
}
