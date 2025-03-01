package sqlite

import (
	"xCelFlow/parser"
)

type SQLParser struct {
	*parser.Parser
}

func init() {
	parser.RegisterParser("sqlite", NewParser)
}

func NewParser(p *parser.Parser) parser.IParser {
	return &SQLParser{p}
}
