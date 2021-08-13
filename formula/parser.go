package formula

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	NULL     = iota
	Literal  //数字
	Operator //操作符
	Function //函数
	Param    //参数
	COMMA    //逗号
)

type Token struct {
	Tok    string
	Type   int
	Offset int
}

type AST struct {
	Tokens    []*Token
	source    string
	currTok   *Token
	currIndex int
	depth     int
	Err       error
}

type Parser struct {
	Str    string
	char   byte
	offset int
	err    error
}

func Parse(s string) ([]*Token, error) {
	p := &Parser{
		Str:  s,
		err:  nil,
		char: s[0],
	}
	tokens := p.parse()
	if p.err != nil {
		return nil, p.err
	}
	return tokens, nil
}

func (p *Parser) parse() []*Token {
	tokens := make([]*Token, 0)
	for {
		tok := p.nextTok()
		if tok == nil {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

func (p *Parser) nextTok() *Token {
	if p.offset >= len(p.Str) || p.err != nil {
		return nil
	}
	//过滤掉无意义的字符
	for p.isWhitespace(p.char) && p.nextChar() {
	}
	start := p.offset
	var tok *Token

	if isOperator(p.char) {
		tok = &Token{
			Tok:  string(p.char),
			Type: Operator,
		}
		tok.Offset = start
		p.nextChar()
	} else if isNumber(p.char) {
		for p.isDigitNum(p.char) && p.nextChar() {
			if p.char == '-' || p.char == '+' {
				break
			}
		}
		tok = &Token{
			Tok:  p.Str[start:p.offset],
			Type: Literal,
		}
		tok.Offset = start
	} else if p.char == ',' {
		tok = &Token{
			Tok:  string(p.char),
			Type: COMMA,
		}
		tok.Offset = start
		if !p.nextChar() {
			s := fmt.Sprintf("input str error after ',' is nothing pos [%v:]\n%s",
				start,
				ErrPos(p.Str, start))
			p.err = errors.New(s)
		}
	} else {
		if p.isChar(p.char) {
			for p.isWordChar(p.char) && p.nextChar() {
			}
			name := p.Str[start:p.offset]
			var tokenType int
			if _, ok := funcMap[name]; ok {
				tokenType = Function
			} else if _, ok := ParamContext[name]; ok {
				tokenType = Param
			} else {
				s := fmt.Sprintf("input str error: unknown '%v', pos [%v:]\n%s",
					name,
					start,
					ErrPos(p.Str, start))
				p.err = errors.New(s)
			}
			tok = &Token{
				Tok:  name,
				Type: tokenType,
			}
			tok.Offset = start
		} else if p.char != ' ' {
			s := fmt.Sprintf("input str  error: unknown '%v', pos [%v:]\n%s",
				string(p.char),
				start,
				ErrPos(p.Str, start))
			p.err = errors.New(s)
		}
	}
	return tok
}

func isOperator(b byte) bool {
	return b == '(' || b == ')' || b == '+' || b == '-' || b == '*' || b == '/'
}

func isNumber(b byte) bool {
	return unicode.IsDigit(rune(b))
}

func (p *Parser) nextChar() bool {
	p.offset++
	if p.offset < len(p.Str) {
		p.char = p.Str[p.offset]
		return true
	}
	return false
}

func (p *Parser) isWhitespace(c byte) bool {
	return c == ' ' ||
		c == '\t' ||
		c == '\n' ||
		c == '\v' ||
		c == '\f' ||
		c == '\r'
}

func (p *Parser) isDigitNum(c byte) bool {
	return '0' <= c && c <= '9' || c == '.'
}

func (p *Parser) isChar(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func (p *Parser) isWordChar(c byte) bool {
	return p.isChar(c) || '0' <= c && c <= '9'
}

func ErrPos(s string, pos int) string {
	r := strings.Repeat("-", len(s)) + "\n"
	s += "\n"
	for i := 0; i < pos; i++ {
		s += " "
	}
	s += "^\n"
	return r + s + r
}

func NewAST(tokens []*Token, s string) *AST {
	a := &AST{
		Tokens: tokens,
		source: s,
	}
	if a.Tokens == nil || len(a.Tokens) == 0 {
		a.Err = errors.New("empty token")
	} else {
		a.currIndex = 0
		a.currTok = a.Tokens[0]
	}
	return a
}

func (a *AST) ParseExpression() Expression {
	a.depth++
	lhs := a.parsePrimary()
	if lhs == nil {
		return nil
	}
	r := a.parseBinaryExpOpRHS(0, lhs)
	a.depth--
	if a.depth == 0 && a.currIndex != len(a.Tokens) && a.Err == nil {
		a.Err = errors.New(
			fmt.Sprintf("error expression, reaching the end or miss the operator\n%s",
				ErrPos(a.source, a.currTok.Offset)))
	}
	return r
}

func (a *AST) getNextToken() *Token {
	a.currIndex++
	if a.currIndex < len(a.Tokens) {
		a.currTok = a.Tokens[a.currIndex]
		return a.currTok
	}
	return nil
}

func (a *AST) getPrecedence() int {
	if p, ok := precedence[a.currTok.Tok]; ok {
		return p
	}
	return -1
}

func (a *AST) getConstantExp() *ConstantExp {
	f, err := strconv.ParseFloat(a.currTok.Tok, 64)
	a.getNextToken()
	if err != nil {
		a.Err = errors.New(
			fmt.Sprintf("%v\n should be '(' or '0-9' but get '%s'\n%s",
				err.Error(),
				a.currTok.Tok,
				ErrPos(a.source, a.currTok.Offset)))
		return &ConstantExp{}
	}
	n := ConstantExp{
		Val: f,
	}
	return &n
}

func (a *AST) getFunCallerExp() Expression {
	name := a.currTok.Tok
	a.getNextToken()

	if a.currTok.Tok != "(" {
		a.Err = errors.New(
			fmt.Sprintf("func `%s` next char is \n%s it should be '(' ",
				name,
				ErrPos(a.source, a.currTok.Offset)))
		return nil
	}

	if _, ok := funcMap[name]; !ok {
		a.Err = errors.New(
			fmt.Sprintf("func `%s` is undefined\n%s",
				name,
				ErrPos(a.source, a.currTok.Offset)))
		return nil
	}
	aes := make([]Expression, 0)
	for a.currTok.Tok != ")" && a.getNextToken() != nil {
		if a.currTok.Type == COMMA {
			continue
		}
		aes = append(aes, a.ParseExpression())
	}
	def := funcMap[name]
	//检查传入的参数个数
	if len(aes) != def.argsCount {
		a.Err = errors.New(
			fmt.Sprintf("error func `%s`, parameters should be %d but get %d\n%s",
				name,
				def.argsCount,
				len(aes),
				ErrPos(a.source, a.currTok.Offset)))
	}

	if a.currTok.Tok != ")" {
		a.Err = errors.New(
			fmt.Sprintf("func `%s` last char not is `)` \n%s ",
				name,
				ErrPos(a.source, a.currTok.Offset)))
	}
	a.getNextToken()
	return &FuncExp{
		Name: name,
		Arg:  aes,
	}
}

func (a *AST) getParamExp() Expression {
	name := a.currTok.Tok
	a.getNextToken()
	if v, ok := ParamContext[name]; ok {
		return &ParamExp{
			Val:  v.(float64),
			Name: name,
		}
	} else {
		a.Err = errors.New(
			fmt.Sprintf("param `%s` is undefined\n%s",
				name,
				ErrPos(a.source, a.currTok.Offset)))
		return &ConstantExp{}
	}
}

func (a *AST) parsePrimary() Expression {
	switch a.currTok.Type {
	case Function:
		return a.getFunCallerExp()
	case Param:
		return a.getParamExp()
	case Literal:
		return a.getConstantExp()
	case Operator:
		return a.parseOperatorExp()
	default:
		return nil
	}
}

func (a *AST) parseOperatorExp() Expression {
	if a.currTok.Tok == "(" {
		t := a.getNextToken()
		if t == nil {
			a.Err = errors.New(
				fmt.Sprintf("should be '0-9' but nothing at all\n%s",
					ErrPos(a.source, a.currTok.Offset)))
			return nil
		}
		e := a.ParseExpression()
		if e == nil {
			return nil
		}
		if a.currTok.Tok != ")" {
			a.Err = errors.New(
				fmt.Sprintf("should be ')' but get %s\n%s",
					a.currTok.Tok,
					ErrPos(a.source, a.currTok.Offset)))
			return nil
		}
		a.getNextToken()
		return e
	} else if a.currTok.Tok == "-" { //表达式起始符号为减号的处理
		if a.getNextToken() == nil {
			a.Err = errors.New(
				fmt.Sprintf("should be '0-9' but get '-'\n%s",
					ErrPos(a.source, a.currTok.Offset)))
			return nil
		}
		c := BinaryExp{
			Op:  "-",
			Lhs: &ConstantExp{},
			Rhs: a.parsePrimary(),
		}
		return &c
	} else {
		return a.getConstantExp()
	}
}

func (a *AST) parseBinaryExpOpRHS(execPrec int, lhs Expression) Expression {
	for {
		prec := a.getPrecedence()
		if prec < execPrec {
			return lhs
		}
		op := a.currTok.Tok
		if a.getNextToken() == nil {
			return lhs
		}
		rhs := a.parsePrimary()
		if rhs == nil {
			return nil
		}
		nextPrec := a.getPrecedence()
		if prec < nextPrec {
			rhs = a.parseBinaryExpOpRHS(prec+1, rhs)
			if rhs == nil {
				return nil
			}
		}
		lhs = &BinaryExp{
			Op:  op,
			Lhs: lhs,
			Rhs: rhs,
		}
	}
}
