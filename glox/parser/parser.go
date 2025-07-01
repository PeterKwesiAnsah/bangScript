package parser

import (
	"fmt"
	"lox/glox/scanner"
	"strconv"
)

type Tokens []*scanner.Token
type Obj interface{}

type Exp interface {
	//should take in an environment or context
	Evaluate(env *Stmtsenv) (Obj, error)
}

type Stmt interface {
	Execute(env *Stmtsenv) error
}

type Stmtsenv struct {
	Local    map[string]Obj
	Encloser *Stmtsenv
}

type assigment struct {
	//l-value
	storeTarget Exp
	operator    *scanner.Token
	//r-value
	right Exp
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

type ifStmt struct {
	condition Exp
	thenbody  Stmt
	elsebody  Stmt
}

type blockStmt struct {
	stmts []Stmt
	env   Stmtsenv
}

type varStmt struct {
	//we expect scanner.Token to be an identifier
	name *scanner.Token
	exp  Exp
}
type printStmt struct {
	exp Exp
}

type expStmt struct {
	exp Exp
}

var current int = 0

func (t ifStmt) Execute(env *Stmtsenv) error
func (tkn Tokens) ifStmt(Encloser *Stmtsenv) (Stmt, error)

func (t blockStmt) Execute(env *Stmtsenv) error {
	for _, stmt := range t.stmts {
		err := stmt.Execute(&t.env)
		if err != nil {
			return err
		}
	}
	return nil
}

// first blockStmt in call stack , will be called with the global env and subsequent ones will be wrapped recursively in a linked list
func (tkn Tokens) blockStmt(Encloser *Stmtsenv) (Stmt, error) {
	inner := Stmtsenv{Local: map[string]Obj{}, Encloser: nil}
	stmts := []Stmt{}
	stmt := blockStmt{}
	for tkn[current].Ttype != scanner.EOF && tkn[current].Ttype != scanner.RIGHT_BRACE {
		stmt, err := tkn.declarations(&inner)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
		//current++
	}
	if tkn[current].Ttype != scanner.RIGHT_BRACE {
		return nil, fmt.Errorf("Expected a right brace but got EOF")
	}
	current++
	stmt.stmts = stmts
	inner.Encloser = Encloser
	stmt.env = inner
	return stmt, nil
}

func (t printStmt) Execute(env *Stmtsenv) error {
	Obj, err := t.exp.Evaluate(env)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", Obj)
	return nil
}
func (tkn Tokens) printStmt() (Stmt, error) {
	stmt := printStmt{}
	exp, err := tkn.expression()
	if err != nil {
		return nil, err
	}
	stmt.exp = exp
	//expect a ";" terminator
	if tkn[current].Ttype != scanner.SEMICOLON {
		return nil, fmt.Errorf("Expected semi-colon but got %d", tkn[current].Ttype)
	}
	current++
	return stmt, nil
}

func (t varStmt) Execute(env *Stmtsenv) error {
	Obj, err := t.exp.Evaluate(env)
	if err != nil {
		return err
	}
	env.Local[t.name.Lexem] = Obj
	return nil
}
func (tkn Tokens) varStmt() (Stmt, error) {
	stmt := varStmt{}
	//expect identifier
	if tkn[current].Ttype != scanner.IDENTIFIER {
		return nil, fmt.Errorf("Expected an identifier after var but got %d", tkn[current].Ttype)
	}
	stmt.name = tkn[current]
	stmt.exp = nil

	//consume identifier
	current++
	//optionally expect an initializer
	if tkn[current].Ttype == scanner.EQUAL {
		current++
		//we expect an initializer expresion or what some may call a variable expression
		exp, err := tkn.expression()
		if err != nil {
			return nil, err
		}
		stmt.exp = exp
	}
	//expect a ";" terminator
	if tkn[current].Ttype != scanner.SEMICOLON {
		return nil, fmt.Errorf("Expected semi-colon but got %d", tkn[current].Ttype)
	}
	current++
	return stmt, nil
}

func (t expStmt) Execute(env *Stmtsenv) error {
	_, err := t.exp.Evaluate(env)
	return err
}
func (tkn Tokens) expStmt() (Stmt, error) {
	stmt := expStmt{}
	exp, err := tkn.expression()
	if err != nil {
		return nil, err
	}
	stmt.exp = exp
	//expect a ";" terminator
	if tkn[current].Ttype != scanner.SEMICOLON {
		return nil, fmt.Errorf("Expected semi-colon but got %d", tkn[current].Ttype)
	}
	return stmt, nil
}

func (tkn Tokens) declarations(Encloser *Stmtsenv) (Stmt, error) {
	curT := tkn[current]
	if curT.Ttype == scanner.VAR {
		current++
		stmt, err := tkn.varStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil
		//Encloser.stmts = append(Encloser.stmts, stmt)
	} else if curT.Ttype == scanner.PRINT {
		current++
		stmt, err := tkn.printStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	} else if curT.Ttype == scanner.LEFT_BRACE {
		current++
		//block statement Only type of statement with multiple statements bounded by an env, all other statements take their context
		return tkn.blockStmt(Encloser)
	}
	//expression statement
	return tkn.expStmt()
}
func Parser(tkn Tokens, globalEnv *Stmtsenv) ([]Stmt, error) {
	stmts := []Stmt{}
	for tkn[current].Ttype != scanner.EOF {
		stmt, err := tkn.declarations(globalEnv)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

// statements are now into two Declaration Statements (are statements like varDeclartion,function,classes. These statements usually cant be used directly after constructs like if and while) and Regular Statements
// program->declarations*EOF
// declarations->varDeclar | statements
// varDeclar->"var" IDENTIFIER (=expression)?";"
// statements->printStmt | block | expressionStmt
// block->"{" declarations* "}"
// printStmt->"print" expression ";"
// expressionStmt->expression;
// expression->assignment
// assigment->equality ("=" assignment)*

////statement rules and their interface implementaion

// end

///////////////////////expression rules and their interface implementation
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
func (p primary) Evaluate(env *Stmtsenv) (Obj, error) {
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
	case scanner.IDENTIFIER:
		{
			cur := env
			for cur != nil {
				Obj, itExist := cur.Local[p.node.Lexem]
				if itExist {
					return Obj, nil
				}
				cur = cur.Encloser
			}
			return nil, fmt.Errorf("Undefined varible at line %d", p.node.Line)
		}
	default:
		return nil, fmt.Errorf("Expected a string, number, a target location , nil and boolean but got %d at line %d", p.node.Ttype, p.node.Line)
	}
}

func (u unary) Evaluate(env *Stmtsenv) (Obj, error) {
	Exp, err := u.right.Evaluate(env)
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

func (b binary) Evaluate(env *Stmtsenv) (Obj, error) {
	left, err := b.left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	right, err := b.right.Evaluate(env)
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

func (a assigment) Evaluate(env *Stmtsenv) (Obj, error) {
	cur := env
	lv, isStorageTarget := a.storeTarget.(primary)
	if !(isStorageTarget && lv.node.Ttype == scanner.IDENTIFIER) {
		return nil, fmt.Errorf("Cannot use the l-value as storage target")
	}
	for cur != nil {
		_, itExist := cur.Local[lv.node.Lexem]
		if itExist {
			rv, err := a.right.Evaluate(env)
			if err != nil {
				return nil, err
			}
			cur.Local[lv.node.Lexem] = rv
			return rv, nil
		}
		cur = cur.Encloser
	}
	return nil, fmt.Errorf("Undefined variable at line %d", a.operator.Line)
}

// TODO: grammer for tenary expressions
// TODO: grammer for grouped expression
// TODO: implement grammer for logical operators && and ||
// TODO: binary operators without left hand operands , report error but continue passing
// Rule for parsing expressions into trees
func (tkn Tokens) expression() (Exp, error) {
	return tkn.asignment()
}
func (tkn Tokens) asignment() (Exp, error) {
	exp, err := tkn.equality()
	if err != nil {
		return nil, err
	}
	//optionally expect "="
	if tkn[current].Ttype == scanner.EQUAL {
		//assigment
		op := tkn[current]
		//currently lv are single nodes, obviosuly identifiers
		lv, isStorageTarget := exp.(primary)
		if isStorageTarget && lv.node.Ttype == scanner.IDENTIFIER {
			//consume "="
			current++
			rv, err := tkn.asignment()
			if err != nil {
				return nil, err
			}
			ass := assigment{lv, op, rv}
			return ass, nil
		}
		//if i had a print method to my exp interface , i could have called it here cool right. Maybe add this later
		return nil, fmt.Errorf("Cannot use the l-value as storage target ")
	}
	return exp, nil
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

// rule for producing operands
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
		{
			//check if next token is EOF , if not consume the LEFT_BRACE token and call expression
			if current+1 >= len(tkn) {
				return nil, fmt.Errorf("Expected an expression token but got EOF")
			}
			current++
			Exp, err := tkn.asignment()
			if err != nil {
				return nil, err
			}
			if tkn[current].Ttype != scanner.RIGHT_PAREN {
				return nil, fmt.Errorf("Expected a RIGHT_BRACE token but got %d", tkn[current].Ttype)
			}
			current++
			return Exp, nil
		}
	default:
		//invalid expresion token
		return nil, fmt.Errorf("Expected an expression token but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	tnode.node = tkn[current]
	current++
	return tnode, nil
}

/////////////////////////////////////////////////////end of expression rules
