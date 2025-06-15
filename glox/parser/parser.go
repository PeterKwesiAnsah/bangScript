package parser

import (
	"fmt"
	"lox/glox/scanner"
	"strconv"
)

type Tokens []*scanner.Token
type obj interface{}

type exp interface {
	evaluate() (obj, error)
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
func (p primary) evaluate() (obj, error) {
	//for evaluting expressions at compile-time we can perform mathematical operations,logical operations and string concantenation
	// operands needs to be only following string , number and boolean
	switch p.node.Ttype {
	case scanner.NUMBER:
		op, err := strconv.ParseFloat(p.node.Lexem, 64)
		if err != nil {
			//handle error for failed type conversion
			return nil, err
		}
		return op, nil
	//string concantenation and comparison
	case scanner.STRING:
		return p.node.Lexem, nil
	//boolean algebra
	case scanner.TRUE:
		return true, nil
	case scanner.FALSE:
		return false, nil
	case scanner.NIL:
		return nil, nil
	default:
		return nil, fmt.Errorf("Expected a string, number,nil and boolean but got %d at line %d", p.node.Ttype, p.node.Line)
	}
}

func (u unary) evaluate() (obj, error) {
	exp, err := u.right.evaluate()
	if err != nil {
		return nil, err
	}
	operator := u.operator.Ttype

	if operator == scanner.BANG {
		bol, isbol := exp.(bool)
		if isbol {
			return !bol, nil
		}
		return nil, fmt.Errorf("Expected a boolean value but got something else at line %d", u.operator.Line)
	} else if operator == scanner.MINUS {
		num, isnum := exp.(float64)
		if isnum {
			return -num, nil
		}
		return nil, fmt.Errorf("Expected a number value but got something else at line %d", u.operator.Line)
	}
	return nil, fmt.Errorf("Invalid expression")
}

func (b binary) evaluate() (obj, error) {
	left, err := b.left.evaluate()
	if err != nil {
		return nil, err
	}
	right, err := b.right.evaluate()
	if err != nil {
		return nil, err
	}
	switch b.operator.Ttype {
	case scanner.PLUS:
		{
			//string concatenation
			// TODO: (parse and concatenate) string + number | number + string
			strLeft, okLeft := left.(string)
			if okLeft {
				strRight, okRight := right.(string)
				if !okRight {
					return nil, fmt.Errorf("Invalid Right operand, expected a string at line %d", b.operator.Line)
				}
				return strLeft + strRight, nil
			}
			// integer addition
			floatLeft, okLeft := left.(float64)
			if okLeft {
				floatRight, okRight := right.(float64)
				if !okRight {
					return nil, fmt.Errorf("Invalid Right operand, expected a number at line %d", b.operator.Line)
				}
				return floatLeft + floatRight, nil
			}
		}
	case scanner.SLASH:
		{
			//integer division
			nLeft, okLeft := left.(float64)
			nRight, okRight := right.(float64)
			if okLeft && okRight {
				return (nLeft / nRight), nil
			}
			return nil, fmt.Errorf("Invalid expression.")
		}
	case scanner.MINUS:
		{
			//integer division
			nLeft, okLeft := left.(float64)
			nRight, okRight := right.(float64)
			if okLeft && okRight {
				return (nLeft - nRight), nil
			}
			return nil, fmt.Errorf("Invalid expression.")
		}
	case scanner.GREATER:
		{
			//integer comparison
			nLeft, okLeft := left.(float64)
			nRight, okRight := right.(float64)
			if okLeft && okRight {
				return (nLeft > nRight), nil
			}
			return nil, fmt.Errorf("Invalid expression.")
		}
	case scanner.GREATER_EQUAL:
		{
			//integer comparison
			nLeft, okLeft := left.(float64)
			nRight, okRight := right.(float64)
			if okLeft && okRight {
				return (nLeft >= nRight), nil
			}
			return nil, fmt.Errorf("Invalid expression.")
		}
	case scanner.LESS:
		{
			//integer comparison
			nLeft, okLeft := left.(float64)
			nRight, okRight := right.(float64)
			if okLeft && okRight {
				return (nLeft < nRight), nil
			}
			return nil, fmt.Errorf("Invalid expression.")
		}
	case scanner.LESS_EQUAL:
		{
			//integer comparison
			nLeft, okLeft := left.(float64)
			nRight, okRight := right.(float64)
			if okLeft && okRight {
				return (nLeft <= nRight), nil
			}
			return nil, fmt.Errorf("Invalid expression.")
		}
	case scanner.EQUAL_EQUAL:
		{
			switch left.(type) {
			case string:
				{
					//string comparison
					str, isStr := right.(string)
					if isStr {
						return left == str, nil
					}
				}
			case float64:
				{
					//integer equality check
					num, isNum := right.(float64)
					if isNum {
						return left == num, nil
					}
				}
			case bool:
				{
					//boolean arithemetic
					bool, isBool := right.(bool)
					if isBool {
						return left == bool, nil
					}
				}
			default:
				// no match;
				return nil, fmt.Errorf("Invalid expression.")
			}
			return nil, fmt.Errorf("Invalid expression.")
		}
	case scanner.BANG_EQUAL:
		{
			switch left.(type) {
			case string:
				{
					strRight, isStr := right.(string)
					if isStr {
						return left != strRight, nil
					}
				}
			case float64:
				{
					numRight, isNum := right.(float64)
					if isNum {
						return left != numRight, nil
					}
				}
			case bool:
				{
					boolRight, isBool := right.(bool)
					if isBool {
						return left != boolRight, nil
					}
				}
			default:
				// no match;
				return nil, fmt.Errorf("Invalid expression.")
			}
			return nil, fmt.Errorf("Invalid expression.")
		}
	default:
		{
			return nil, fmt.Errorf("Invalid operator at line %d.", b.operator.Line)
		}

	}
	return nil, fmt.Errorf("Invalid expression.")
}

//func (exp tenary) evaluate()  {}

var current int = 0

func Parser(tkn Tokens) (exp, error) {
	exp, err := tkn.expression()
	//For testing cases, we need to reset the current counter
	current = 0
	return exp, err
}

// TODO: grammer for tenary expressions
// TODO: grammer for grouped expression
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
