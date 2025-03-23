package types

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"xCelFlow/core"
)

// TokenType 表示不同类型的标记。
type TokenType int

const (
	EOF      TokenType = iota // 表示输入结束
	INTEGER                   // 表示整数
	FLOAT                     // 表示数字
	STRING                    // 表示字符串
	BOOLEAN                   // 表示布尔值
	ANY                       // 表示任意数据
	LBRACKET                  // 表示左方括号
	RBRACKET                  // 表示右方括号
	LPAREN                    // 表示左圆括号
	RPAREN                    // 表示右圆括号
	LBRACE                    // 表示左大括号
	RBRACE                    // 表示右大括号
	EQUAL                     // 表示等号
	COMMA                     // 表示逗号
)

// token 表示一个标记及其类型和值。
type token struct {
	Type  TokenType
	Value string
}

// lexer 负责将输入字符串转换为标记。
type lexer struct {
	text        string
	pos         int
	currentChar byte
}

// parser 负责将标记序列解析为数据结构。
type parser struct {
	lexer        *lexer
	currentToken token
}

// ParseString 将输入字符串转换为数据结构。
func ParseString(text string) (v any, err error) {
	defer func() {
		if r := recover(); r != nil {
			// 捕获 panic 的值
			err = fmt.Errorf("%v", r)
		}
	}()
	if text == "" {
		return nil, nil
	}
	lexer := &lexer{text: text, pos: 0, currentChar: text[0]}
	p := &parser{lexer: lexer, currentToken: lexer.getNextToken()}
	return p.parse(), nil
}

// parse 解析输入并返回结果数据结构。
func (p *parser) parse() any {
	switch p.currentToken.Type {
	case LBRACKET:
		return p.parseList()
	case LBRACE:
		return p.parseMap()
	case LPAREN:
		return p.parseTuple()
	default:
		return p.parseValue()
	}
}

// advance 移动 'pos' 指针并设置 'currentChar'。
func (l *lexer) advance() {
	l.pos++
	if l.pos >= len(l.text) {
		l.currentChar = 0 // 表示输入结束
	} else {
		l.currentChar = l.text[l.pos]
	}
}

// number 返回一个数字标记。
func (l *lexer) number() token {
	var result strings.Builder
	var isFloat bool

	// 检查负号
	if l.currentChar == '-' {
		result.WriteByte(l.currentChar)
		l.advance()
	}

loop:
	for {
		switch l.currentChar {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			result.WriteByte(l.currentChar)
			l.advance()
		case '.':
			result.WriteByte(l.currentChar)
			l.advance()
			isFloat = true
		default:
			break loop
		}
	}

	if isFloat {
		return token{Type: FLOAT, Value: result.String()}
	}
	return token{Type: INTEGER, Value: result.String()}
}

// string 返回一个字符串标记。
func (l *lexer) string() token {
	l.advance() // 跳过开头的引号
	var result strings.Builder
loop:
	for {
		switch l.currentChar {
		case 0:
			break loop
		case '"':
			l.advance() // 跳过结尾的引号
			break loop
		default:
			result.WriteByte(l.currentChar)
			l.advance()
		}
	}

	return token{Type: STRING, Value: result.String()}
}

// boolean 返回一个bu布尔标记
func (l *lexer) boolean() token {
	curPos, curChar := l.pos, l.currentChar
	if strings.HasPrefix(l.text[l.pos:], "true") || strings.HasPrefix(l.text[l.pos:], "TRUE") {
		l.pos = l.pos + 3
		l.advance()
		if l.currentChar == 0 || !(unicode.IsLetter(rune(l.currentChar)) || unicode.IsDigit(rune(l.currentChar))) {
			return token{Type: BOOLEAN, Value: "true"}
		}
		l.pos, l.currentChar = curPos, curChar
	}

	if strings.HasPrefix(l.text[l.pos:], "false") || strings.HasPrefix(l.text[l.pos:], "FALSE") {
		l.pos = l.pos + 4
		l.advance()
		if l.currentChar == 0 || !(unicode.IsLetter(rune(l.currentChar)) || unicode.IsDigit(rune(l.currentChar))) {
			return token{Type: BOOLEAN, Value: "false"}
		}
		l.pos, l.currentChar = curPos, curChar
	}
	return l.any()
}

func (l *lexer) any() token {
	var result strings.Builder
loop:
	for {
		switch l.currentChar {
		case 0:
			break loop
		case ']', ')', '}', ',', '=', '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0: // 遇到这些字符则Raw类型已经结束
			break loop
		default:
			result.WriteByte(l.currentChar)
			l.advance()
		}
	}

	return token{Type: ANY, Value: result.String()}
}

// getNextToken 将输入拆分为标记。
func (l *lexer) getNextToken() token {
	for {
		switch l.currentChar {
		case 0:
			return token{Type: EOF, Value: ""}
		case '[':
			l.advance()
			return token{Type: LBRACKET, Value: "["}
		case ']':
			l.advance()
			return token{Type: RBRACKET, Value: "]"}
		case '(':
			l.advance()
			return token{Type: LPAREN, Value: "("}
		case ')':
			l.advance()
			return token{Type: RPAREN, Value: ")"}
		case '{':
			l.advance()
			return token{Type: LBRACE, Value: "{"}
		case '}':
			l.advance()
			return token{Type: RBRACE, Value: "}"}
		case '=':
			l.advance()
			return token{Type: EQUAL, Value: "="}
		case ',':
			l.advance()
			return token{Type: COMMA, Value: ","}
		case '"':
			return l.string()
		case '-', '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return l.number()
		case 't', 'T', 'f', 'F':
			return l.boolean()
		case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
			l.advance()
			continue
		default:
			if unicode.IsLetter(rune(l.currentChar)) {
				return l.any()
			}

			l.panic()
		}
	}
}

func (l *lexer) panic() {
	panic(fmt.Sprintf("解析值时出错 %s 位置:%d 字符:%c", l.text, l.pos+1, l.currentChar))
}

// eat 消耗当前标记并移动到下一个标记。
func (p *parser) eat(tokenType TokenType) {
	if p.currentToken.Type == tokenType {
		p.currentToken = p.lexer.getNextToken()
	} else {
		p.panic()
	}
}

// parseList 解析列表结构。
func (p *parser) parseList() []any {
	p.eat(LBRACKET)
	var list []any
	for p.currentToken.Type != RBRACKET {
		list = append(list, p.parse())
		if p.currentToken.Type == COMMA {
			p.eat(COMMA)
		}
	}
	p.eat(RBRACKET)
	return list
}

// parseTuple 解析元组结构。
func (p *parser) parseTuple() core.TupleT {
	p.eat(LPAREN)
	var tuple core.TupleT
	for i := 0; p.currentToken.Type != RPAREN; i++ {
		if i < len(tuple) {
			tuple[i] = p.parse()
			if p.currentToken.Type == COMMA {
				p.eat(COMMA)
			}
		} else {
			p.panic()
		}
	}
	p.eat(RPAREN)
	return tuple
}

// parseMap 解析Map结构。
func (p *parser) parseMap() map[any]any {
	p.eat(LBRACE)
	m := make(map[any]any)
	for p.currentToken.Type != RBRACE {
		key := p.parse()
		p.eat(EQUAL)
		value := p.parse()
		m[key] = value
		if p.currentToken.Type == COMMA {
			p.eat(COMMA)
		}
	}
	p.eat(RBRACE)
	return m
}

// parseValue 解析单个值。
func (p *parser) parseValue() any {
	switch p.currentToken.Type {
	case INTEGER:
		value, _ := strconv.ParseInt(p.currentToken.Value, 10, 64)
		p.eat(INTEGER)
		return value
	case FLOAT:
		value, _ := strconv.ParseFloat(p.currentToken.Value, 64)
		p.eat(FLOAT)
		return value
	case STRING:
		value := p.currentToken.Value
		p.eat(STRING)
		return value
	case BOOLEAN:
		value := p.currentToken.Value == "true"
		p.eat(BOOLEAN)
		return value
	case ANY:
		value := core.AnyT(p.currentToken.Value)
		p.eat(ANY)
		return value
	default:
		p.panic()
		return nil
	}
}

func (p *parser) panic() {
	panic(fmt.Sprintf("解析值时出错 %s 位置:%d 字符:%c", p.lexer.text, p.lexer.pos+1, p.lexer.currentChar))
}
