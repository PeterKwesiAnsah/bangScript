package parser

import (
	"fmt"
	"lox/glox/scanner"
)

type tokens []*scanner.Token

type exp interface {
	print()
}

type binary struct {
	left     exp
	operator *scanner.Token
	right    exp
}

type unary struct {
	operator *scanner.Token
	right    exp
}

type primary struct {
	node *scanner.Token
}

// implement the exp interface
func (exp *binary) print()  {}
func (exp *unary) print()   {}
func (exp *primary) print() {}

var current int = 0

func Parser(tkn tokens) (exp, error) {
	return tkn.expression()
}

func (tkn tokens) expression() (exp, error) {
	var temp exp
	return temp, nil
}

func (tkn tokens) equality() (exp, error) {
	var temp exp
	return temp, nil
}
func (tkn tokens) comparison() (exp, error) {
	var temp exp
	return temp, nil
}
func (tkn tokens) term() (exp, error) {
	var temp exp
	return temp, nil
}

func (tkn tokens) factor() (exp, error) {
	var temp exp
	return temp, nil
}
func (tkn tokens) unary() (exp, error) {
	var temp exp
	return temp, nil
}
func (tkn tokens) primary() (exp, error) {
	ttype := tkn[current].Ttype
	tnode := primary{}
	switch ttype {
	case scanner.IDENTIFIER:
	case scanner.STRING:
	case scanner.TRUE:
	case scanner.FALSE:
	case scanner.NIL:
	case scanner.LEFT_BRACE:
		//check if next token is EOF , if not consume the LEFT_BRACE token and call expression
		if current+1 >= len(tkn) {
			return nil, fmt.Errorf("Expected an expression token but got EOF")
		}
		exp, err := tkn.expression()
		if err != nil {
			return nil, err
		}
		if tkn[current].Ttype != scanner.RIGHT_BRACE {
			return nil, fmt.Errorf("Expected a RIGHT_BRACE token but got %d", tkn[current].Ttype)
		}
		return exp, nil
	default:
		//invalid expresion token
		return nil, fmt.Errorf("Expected an expression token but got %d", tkn[current].Ttype)
	}
	tnode.node = tkn[current]
	current++
	return &tnode, nil
}
