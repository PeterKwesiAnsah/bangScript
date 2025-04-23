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

// TODO: handle multiple expression rule
// TODO: handle tenary expression rule
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
	return tkn.equality()
}

func (tkn tokens) equality() (exp, error) {
	cexpleft, err := tkn.comparison()
	if err != nil {
		return nil, err
	}
	for {
		cToken := tkn[current]
		// find the operator terminal
		if cToken.Ttype == scanner.EQUAL_EQUAL || cToken.Ttype == scanner.BANG_EQUAL {
			//consume operator terminal
			current++
			op := cToken
			cexpright, err := tkn.comparison()
			if err != nil {
				return nil, err
			}
			cexpleft = &binary{left: cexpleft, operator: op, right: cexpright}
		}
		break
	}
	return cexpleft, nil
}
func (tkn tokens) comparison() (exp, error) {
	texpleft, err := tkn.term()
	opsToMatch := []scanner.Tokentype{
		scanner.GREATER,
		scanner.GREATER_EQUAL,
		scanner.LESS,
		scanner.LESS_EQUAL,
	}
	if err != nil {
		return nil, err
	}
Matching_Loop:
	for {
		cToken := tkn[current]
		// find the operator terminal
		for _, op := range opsToMatch {
			if cToken.Ttype == op {
				//consume operator terminal
				current++
				op := cToken
				texpright, err := tkn.term()
				if err != nil {
					return nil, err
				}
				texpleft = &binary{left: texpleft, operator: op, right: texpright}
				break Matching_Loop
			}
		}
	}
	return texpleft, nil
}
func (tkn tokens) term() (exp, error) {
	fexpleft, err := tkn.factor()
	opsToMatch := []scanner.Tokentype{
		scanner.PLUS,
		scanner.MINUS,
	}
	if err != nil {
		return nil, err
	}
Matching_Loop:
	for {
		cToken := tkn[current]
		// find the operator terminal
		for _, op := range opsToMatch {
			if cToken.Ttype == op {
				//consume operator terminal
				current++
				op := cToken
				fexpright, err := tkn.factor()
				if err != nil {
					return nil, err
				}
				fexpleft = &binary{left: fexpleft, operator: op, right: fexpright}
				break Matching_Loop
			}
		}
	}
	return fexpleft, nil
}

func (tkn tokens) factor() (exp, error) {
	uexpleft, err := tkn.unary()
	opsToMatch := []scanner.Tokentype{
		scanner.STAR,
		scanner.SLASH,
	}
	if err != nil {
		return nil, err
	}
Matching_Loop:
	for {
		cToken := tkn[current]
		// find the operator terminal
		for _, op := range opsToMatch {
			if cToken.Ttype == op {
				//consume operator terminal
				current++
				op := cToken
				fexpright, err := tkn.unary()
				if err != nil {
					return nil, err
				}
				uexpleft = &binary{left: uexpleft, operator: op, right: fexpright}
				break Matching_Loop
			}
		}
	}
	return uexpleft, nil
}
func (tkn tokens) unary() (exp, error) {
	uToken := tkn[current]
	if uToken.Ttype == scanner.BANG || uToken.Ttype == scanner.MINUS {
		op := uToken
		//consume operator terminal
		current++
		uexp, err := tkn.unary()
		if err != nil {
			return nil, err
		}
		return &unary{operator: op, right: uexp}, nil
	}
	return tkn.primary()
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
		current++
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
