package parser

import (
	"fmt"
	"lox/glox/scanner"
	"strconv"
)

type Tokens []*scanner.Token
type obj interface{}

type Exp interface {
	Evaluate() (obj, error)
}

// statements are now into two Declaration Statements (are statements like varDeclartion,function,classes. These statements usually cant be used directly after constructs like if and while) and Regular Statements
// program->declarations*EOF
// declarations->varDeclar | statements
// varDeclar->"var" IDENTIFIER (=expression)?";"
// statements->printStmt | expressionStmt
// printStmt->"print" expression ";"
// expressionStmt->expression;
type Stmt interface {
	Execute() error
}

type binary struct {
	left     Exp
	operator *scanner.Token
	right    Exp
}

// for handling conditional operations
type tenary struct {
	condition Exp
	operator  *scanner.Token
	then      Exp
	elsef     Exp
}

type unary struct {
	operator *scanner.Token
	right    Exp
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

// implement the Exp interface
func (p primary) Evaluate() (obj, error) {
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
	//TODO: case scanner.IDENTIFIER (variable expression)
	default:
		return nil, fmt.Errorf("Expected a string, number,nil and boolean but got %d at line %d", p.node.Ttype, p.node.Line)
	}
}

func (u unary) Evaluate() (obj, error) {
	Exp, err := u.right.Evaluate()
	if err != nil {
		return nil, err
	}
	operator := u.operator.Ttype

	if operator == scanner.BANG {
		bol, isbol := Exp.(bool)
		if isbol {
			return !bol, nil
		}
		return nil, fmt.Errorf("Expected a boolean value but got something else at line %d", u.operator.Line)
	} else if operator == scanner.MINUS {
		num, isnum := Exp.(float64)
		if isnum {
			return -num, nil
		}
		return nil, fmt.Errorf("Expected a number value but got something else at line %d", u.operator.Line)
	}
	return nil, fmt.Errorf("Invalid expression")
}

func (b binary) Evaluate() (obj, error) {
	left, err := b.left.Evaluate()
	if err != nil {
		return nil, err
	}
	right, err := b.right.Evaluate()
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
	case scanner.STAR:
		{
			//integer division
			nLeft, okLeft := left.(float64)
			nRight, okRight := right.(float64)
			if okLeft && okRight {
				return (nLeft * nRight), nil
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

//func (Exp tenary) Evaluate()  {}

var current int = 0

func Parser(tkn Tokens) (Exp, error) {
	Exp, err := tkn.expression()
	//For testing cases, we need to reset the current counter
	current = 0
	return Exp, err
}

// TODO: grammer for tenary expressions
// TODO: grammer for grouped expression
// TODO: implement grammer for logical operators && and ||
// TODO: binary operators without left hand operands , report error but continue passing
// Rule for parsing expressions into trees
func (tkn Tokens) expression() (Exp, error) {
	return tkn.equality()
}

func (tkn Tokens) equality() (Exp, error) {
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
func (tkn Tokens) comparison() (Exp, error) {
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
func (tkn Tokens) term() (Exp, error) {
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

func (tkn Tokens) factor() (Exp, error) {
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
func (tkn Tokens) unary() (Exp, error) {
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
func (tkn Tokens) primary() (Exp, error) {
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
		Exp, err := tkn.expression()
		if err != nil {
			return nil, err
		}
		if tkn[current].Ttype != scanner.RIGHT_PAREN {
			return nil, fmt.Errorf("Expected a RIGHT_BRACE token but got %d", tkn[current].Ttype)
		}
		current++
		return Exp, nil
	default:
		//invalid expresion token
		return nil, fmt.Errorf("Expected an expression token but got %d", tkn[current].Ttype)
	}
	tnode.node = tkn[current]
	current++
	return tnode, nil
}

type VDT struct {
	//we expect scanner.Token to be an identifier
	name scanner.Token
	exp  Exp
}

// Rule for parsing variable declaration statements into trees
func (tkn Tokens) VarDeclaration() (VDT, error) {

}
