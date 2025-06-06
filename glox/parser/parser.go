package parser

import (
	"fmt"
	"lox/glox/scanner"
	"strconv"
)

type Tokens []*scanner.Token
type obj interface{}

type exp interface {
	evaluate() obj
}

type binary struct {
	left     exp
	operator *scanner.Token
	right    exp
}

// for handling conditional operations
type tenary struct {
	condition exp
	operator  *scanner.Token
	then      exp
	elsef     exp
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

/**
 * evaluating expressions also defines what the user can do and what types or operands can perform this operations
 * binary arithemetic operations (+)(/)(*)(-)
 * (+) has operator overloading
 * (+) (string + string) string, (double + double) double (only)
 * (/) (double / double) double (only)
 * (-) (double - double) double (only)
 * unary operations (!)(-)
 * (!) (boolean) bool
 * (-) (double) double
 * binary logical operations (==)(!=)(>)(>=)(<)(<=)
 * (==)(string | double | boolean) bool
 * (!=)(string | double | boolean) bool
 * (>)(>=)(<)(<=) (double  (>)(>=)(<)(<=) double ) bool
 * TODO: add support for logical (&&) (||)
 */

// implement the exp interface
func (p primary) evaluate() obj {
	if p.node != nil {
		//for evaluting expressions at compile-time we can perform mathematical operations,logical operations and string concantenation
		// operands needs to be only following string , number and boolean
		switch p.node.Ttype {
		case scanner.NUMBER:
			op, err := strconv.ParseFloat(p.node.Lexem, 64)
			if err == nil {
				//handle error for failed type conversion
				return nil
			}
			return op
		//string concantenation and comparison
		case scanner.STRING:
			return p.node.Lexem
		//boolean algebra
		case scanner.TRUE:
			return true
		case scanner.FALSE:
			return false
		case scanner.NIL:
		}
	}
	return nil
}

func (u unary) evaluate() obj {
	exp := u.right.evaluate()
	//switch case for !(boolean expression) and -(integer expresion)
	if exp != nil && exp != float64(0) {
		return true
	}
	return false
}

func (b binary) evaluate() obj {
	//left := b.left.evaluate()
	//right := b.right.evaluate()
	return 1
}

//func (exp tenary) evaluate()  {}

var current int = 0

func Parser(tkn Tokens) (exp, error) {
	exp, err := tkn.expression()
	//For testing cases, we need to reset the current counter
	current = 0
	return exp, err
}

/*
 * // TODO:tenary expression
 // TODO:nested tenary expression
 // TODO:precedence
 func (tkn Tokens) tenary() (exp, error) {
	//may have a condition expression or not
	eexpleft, err := tkn.equality()
	if err != nil {
		return nil, err
	}

	eToken := tkn[current]
	// find the comma operator terminal
	if eToken.Ttype == scanner.QUESTION {
		//consume question (?)
		current++
		//then expression
		//call from the top ?? call multiple to handle comma seperated operations
		expthen, err := tkn.equality()
		if err != nil {
			return nil, err
		}
		if (tkn[current].Ttype) != scanner.COLON {
			return nil, fmt.Errorf("Expected the COLON token but got %d", tkn[current].Ttype)
		}
		//consume colon (:)
		current++
		//handle else expression
		expelse, err := tkn.equality()
		if err != nil {
			return nil, err
		}

		eexpleft = tenary{condition: eexpleft, operator: eToken, then: expthen, elsef: expelse}
	}

	return eexpleft, nil
 }
*/

/*
 *
 * func (tkn Tokens) multiple() (exp, error) {
	texpleft, err := tkn.tenary()
	if err != nil {
		return nil, err
	}
	for {
		tToken := tkn[current]
		// find the comma operator terminal
		if tToken.Ttype == scanner.COMMA {
			op := tToken
			//consume comma(,)
			current++
			texpright, err := tkn.tenary()
			if err != nil {
				return nil, err
			}
			texpleft = binary{left: texpleft, operator: op, right: texpright}
			continue
		}
		break
	}
	return texpleft, nil
 }
*/

// TODO: implement grammer for logical operators && and ||
// TODO: binary operators without left hand operands , report error but continue passing
func (tkn Tokens) expression() (exp, error) {
	return tkn.equality()
}

func (tkn Tokens) equality() (exp, error) {
	cexpleft, err := tkn.comparison()
	if err != nil {
		return nil, err
	}
	for {
		cToken := tkn[current]
		// find the operator terminal
		if cToken.Ttype == scanner.EQUAL_EQUAL || cToken.Ttype == scanner.BANG_EQUAL {
			//consume operator terminal(==,!=)
			current++
			op := cToken
			cexpright, err := tkn.comparison()
			if err != nil {
				return nil, err
			}
			cexpleft = binary{left: cexpleft, operator: op, right: cexpright}
		}
		break
	}
	return cexpleft, nil
}
func (tkn Tokens) comparison() (exp, error) {
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
				texpleft = binary{left: texpleft, operator: op, right: texpright}
				break Matching_Loop
			}
		}
		break
	}
	return texpleft, nil
}
func (tkn Tokens) term() (exp, error) {
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
				fexpleft = binary{left: fexpleft, operator: op, right: fexpright}
				break Matching_Loop
			}
		}
		break
	}
	return fexpleft, nil
}

func (tkn Tokens) factor() (exp, error) {
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
		//println(cToken.Ttype)
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
				uexpleft = binary{left: uexpleft, operator: op, right: fexpright}
				break Matching_Loop
			}
		}
		break
		//fmt.Println("Stuck here")
	}
	return uexpleft, nil
}
func (tkn Tokens) unary() (exp, error) {
	uToken := tkn[current]
	if uToken.Ttype == scanner.BANG || uToken.Ttype == scanner.MINUS {
		op := uToken
		//consume operator terminal
		current++
		uexp, err := tkn.unary()
		if err != nil {
			return nil, err
		}
		return unary{operator: op, right: uexp}, nil
	}
	return tkn.primary()
}
func (tkn Tokens) primary() (exp, error) {
	ttype := tkn[current].Ttype
	tnode := primary{}
	switch ttype {
	case scanner.IDENTIFIER:
	case scanner.NUMBER:
	case scanner.STRING:
	case scanner.TRUE:
	case scanner.FALSE:
	case scanner.NIL:
	case scanner.LEFT_PAREN:
		//check if next token is EOF , if not consume the LEFT_BRACE token and call expression
		if current+1 >= len(tkn) {
			return nil, fmt.Errorf("Expected an expression token but got EOF")
		}
		current++
		exp, err := tkn.expression()
		if err != nil {
			return nil, err
		}
		if tkn[current].Ttype != scanner.RIGHT_PAREN {
			return nil, fmt.Errorf("Expected a RIGHT_BRACE token but got %d", tkn[current].Ttype)
		}
		current++
		return exp, nil
	default:
		//invalid expresion token
		return nil, fmt.Errorf("Expected an expression token but got %d", tkn[current].Ttype)
	}
	tnode.node = tkn[current]
	current++
	return tnode, nil
}
