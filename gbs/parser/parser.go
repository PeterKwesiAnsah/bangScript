package parser

import (
	"bangScript/gbs/scanner"
	"fmt"
)

const (
	REPL uint8 = iota
	SCRIPT
)

const (
	STATIC uint8 = iota
	DYNAMIC
)

type StmtWithEnv interface {
	StaticToDynamic(*Stmtsenv) error
}

type Tokens []*scanner.Token
type Obj interface{}

type Exp interface {
	//Evaluate(env *Stmtsenv) (Obj, error)
	print() string
}

type Stmt interface {
	//Execute(env *Stmtsenv) error
	print() string
}
type CallDetails struct {
	args Exp
	name string
	at   int
}
type CallStack []*CallDetails

// TODO: have a map of scanner token type to string (12 becomes colon)
type Callable interface {
	call(*CallStack, *Call) (Obj, error)
}

type Stmtsenv struct {
	Local    map[string]Obj
	Encloser *Stmtsenv
	Policy   uint8
}

var current int = 0
var cs CallStack = CallStack{}
var Mode uint8 = REPL

func (t *Stmtsenv) Get(depth int) *Stmtsenv {
	for range depth - 1 {
		t = t.Encloser
	}
	return t
}

func (tkn Tokens) returnStmt() (Stmt, error) {
	exp, err := tkn.expression()
	if err != nil {
		return nil, err
	}
	//expect a ";" terminator
	if tkn[current].Ttype != scanner.SEMICOLON {
		return nil, fmt.Errorf("Expected semi-colon but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
	return ReturnStmt{
		Exp: exp,
	}, nil
}
func (tkn Tokens) funcDef(env *Stmtsenv) (Stmt, error) {
	name := tkn[current]
	params := []*scanner.Token{}

	if tkn[current].Ttype != scanner.IDENTIFIER {
		return nil, fmt.Errorf("ParseError: Expected an identifier but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
	if tkn[current].Ttype != scanner.LEFT_PAREN {
		return nil, fmt.Errorf("ParseError: Expected a left paren but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
	if tkn[current].Ttype != scanner.RIGHT_PAREN {
		//non-empty parameters
		// tkn[current].tokenType should be Identifier. TODO: assert tokentype against type Identifier
		params = append(params, tkn[current])
		current++
		for {
			if tkn[current].Ttype == scanner.COMMA {
				current++
				params = append(params, tkn[current])
				current++
				continue
			}
			break
		}
	}
	if tkn[current].Ttype != scanner.RIGHT_PAREN {
		return nil, fmt.Errorf("ParseError: Expected a right paren but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
	if tkn[current].Ttype != scanner.LEFT_BRACE {
		return nil, fmt.Errorf("ParseError: Expected a left brace but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}

	inner := Stmtsenv{Local: map[string]Obj{}, Encloser: env}
	stmts := []Stmt{}
	stmt := BlockStmt{}
	current++
	for tkn[current].Ttype != scanner.EOF && tkn[current].Ttype != scanner.RIGHT_BRACE {
		stmt, err := tkn.declarations(&inner)
		//TODO:if stmt is continue/break or any other loop related statement we throw error because loops will now handle parsing it's body
		//TODO: handle the above with the resolver
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	if tkn[current].Ttype != scanner.RIGHT_BRACE {
		return nil, fmt.Errorf("ParseError: Expected a right brace but got EOF at line %d", tkn[current].Line)
	}
	current++
	stmt.Stmts = stmts
	stmt.Env = &inner

	return FuncDef{
		Arrity: len(params),
		Name:   name,
		Body:   stmt,
		Params: params,
	}, nil
}

func (tkn Tokens) forStmt(env *Stmtsenv) (Stmt, error) {
	var initializer Stmt = nil
	var condition Exp = nil
	var sideEffect Exp = nil
	var body Stmt = nil
	var err error
	//create scope for initializer and condition
	whileScope := Stmtsenv{Local: map[string]Obj{}, Encloser: env}
	if tkn[current].Ttype != scanner.LEFT_PAREN {
		return nil, fmt.Errorf("Expected left paren after for but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
	if tkn[current].Ttype != scanner.SEMICOLON {
		initializer, err = tkn.declarations(nil)
		if err != nil {
			return nil, err
		}
		//automatically consumes the semi-colon token because we are parsing the initializer as a statement (expression statement for assignment/variable decl)
	} else {
		//empty initializer (variable declaration/Assigment)
		current++
	}
	if tkn[current].Ttype != scanner.SEMICOLON {
		condition, err = tkn.expression()
		if err != nil {
			return nil, err
		}
		//expect semicolon
		if tkn[current].Ttype != scanner.SEMICOLON {
			return nil, fmt.Errorf("Expected semi-colon after condition but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
		}
		current++
	} else {
		//empty condition expression
		current++
	}
	if tkn[current].Ttype == scanner.RIGHT_PAREN {
		current++
	} else {
		sideEffect, err = tkn.expression()
		if err != nil {
			return nil, err
		}
		//expect right paren
		if tkn[current].Ttype != scanner.RIGHT_PAREN {
			return nil, fmt.Errorf("Expected right paren after side effect expression but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
		}
		current++
	}
	// empty body statement
	if tkn[current].Ttype == scanner.SEMICOLON {
		return ForStmt{Stmt: WhileStmt{Condition: condition, Init: initializer, Env: &whileScope,
			Body: BlockStmt{Stmts: []Stmt{ExpStmt{Exp: sideEffect}}, Env: &whileScope}}}, nil
	}
	body, err = tkn.declarations(&whileScope)
	if err != nil {
		return nil, err
	}
	//block body statement
	bs, isBS := body.(BlockStmt)
	if isBS {
		stmts := bs.Stmts
		stmts = append(stmts, ExpStmt{Exp: sideEffect})
		return ForStmt{Stmt: WhileStmt{Condition: condition, Init: initializer, Env: &whileScope, Body: BlockStmt{Stmts: stmts, Env: bs.Env}}}, nil
	} else {
		//single statement body
		return ForStmt{Stmt: WhileStmt{Condition: condition, Init: initializer, Env: &whileScope,
			Body: BlockStmt{Stmts: []Stmt{body, ExpStmt{Exp: sideEffect}},
				Env: &whileScope}}}, nil
	}
}
func (tkn Tokens) blockStmt(Encloser *Stmtsenv) (Stmt, error) {
	inner := Stmtsenv{Local: map[string]Obj{}, Encloser: Encloser}
	stmts := []Stmt{}
	stmt := BlockStmt{}
	for tkn[current].Ttype != scanner.EOF && tkn[current].Ttype != scanner.RIGHT_BRACE {
		stmt, err := tkn.declarations(&inner)
		// TODO:if stmt is continue/break/return or any other loop related statement we throw error because loops will now handle parsing it's body
		// TODO: handle the above in the resolver
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	if tkn[current].Ttype != scanner.RIGHT_BRACE {
		return nil, fmt.Errorf("ParseError: Expected a right brace but got EOF at line %d", tkn[current].Line)
	}
	current++
	stmt.Stmts = stmts
	stmt.Env = &inner
	return stmt, nil
}
func (tkn Tokens) whileStmt(Encloser *Stmtsenv) (Stmt, error) {
	if tkn[current].Ttype != scanner.LEFT_PAREN {
		return nil, fmt.Errorf("Expected left paren after while but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
	cond, err := tkn.expression()
	if err != nil {
		return nil, err
	}
	if tkn[current].Ttype != scanner.RIGHT_PAREN {
		return nil, fmt.Errorf("Expected right paren after condition but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
	//TODO: use tkn.statement()
	bodyStmt, err := tkn.declarations(Encloser)
	bs, isBs := bodyStmt.(BlockStmt)
	if !isBs {
		return nil, fmt.Errorf("Expected a block statement")
	}
	if err != nil {
		return nil, err
	}
	return WhileStmt{Condition: cond, Body: bs, Env: Encloser}, nil
}

func (tkn Tokens) ifStmt(Encloser *Stmtsenv) (Stmt, error) {
	if tkn[current].Ttype != scanner.LEFT_PAREN {
		return nil, fmt.Errorf("Expected left paren after if but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
	condexp, err := tkn.expression()
	if err != nil {
		return nil, err
	}
	if tkn[current].Ttype != scanner.RIGHT_PAREN {
		return nil, fmt.Errorf("Expected right paren after condition but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
	//TODO: use tkn.statement()
	stmtBody, err := tkn.declarations(Encloser)
	if err != nil {
		return nil, err
	}
	var elseStmt Stmt = nil
	if tkn[current].Ttype == scanner.ELSE {
		current++
		//TODO: use tkn.statement()
		elseStmt, err = tkn.declarations(Encloser)
		if err != nil {
			return nil, err
		}
	}
	return IfStmt{Condition: condexp, Thenbody: stmtBody, Elsebody: elseStmt}, nil
}
func (tkn Tokens) printStmt() (Stmt, error) {
	stmt := PrintStmt{}
	exp, err := tkn.expression()
	if err != nil {
		return nil, err
	}
	stmt.Exp = exp
	//expect a ";" terminator
	if tkn[current].Ttype != scanner.SEMICOLON {
		return nil, fmt.Errorf("Expected semi-colon but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
	return stmt, nil
}
func (tkn Tokens) varStmt() (Stmt, error) {
	stmt := VarStmt{}
	exp, err := tkn.expression()
	if err != nil {
		return nil, err
	}
	stmt.Exp = exp
	if tkn[current].Ttype != scanner.SEMICOLON {
		return nil, fmt.Errorf("Expected semi-colon but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
	return stmt, nil
}
func (tkn Tokens) expStmt() (Stmt, error) {
	stmt := ExpStmt{}
	exp, err := tkn.expression()
	if err != nil {
		return nil, err
	}
	stmt.Exp = exp
	//expect a ";" terminator only for script mode
	if Mode == REPL {
		if tkn[current].Ttype == scanner.SEMICOLON {
			current++
		}
	} else {
		if tkn[current].Ttype != scanner.SEMICOLON {
			return nil, fmt.Errorf("Expected semi-colon but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
		}
		current++
	}

	return stmt, nil
}

// TODO: have a statement rule , regular statements are different from variable declarations statements
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
	} else if curT.Ttype == scanner.IF {
		current++
		return tkn.ifStmt(Encloser)
	} else if curT.Ttype == scanner.WHILE {
		current++
		return tkn.whileStmt(Encloser)
	} else if curT.Ttype == scanner.FOR {
		current++
		return tkn.forStmt(Encloser)
	} else if curT.Ttype == scanner.FUN {
		current++
		return tkn.funcDef(Encloser)
	} else if curT.Ttype == scanner.RETURN {
		current++
		return tkn.returnStmt()
	}
	//expression statement
	return tkn.expStmt()
}
func Parser(tkn Tokens, globalEnv *Stmtsenv, m uint8) ([]Stmt, error) {
	current = 0
	stmts := []Stmt{}
	Mode = m
	for tkn[current].Ttype != scanner.EOF {
		stmt, err := tkn.declarations(globalEnv)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

// TODO: grammer for tenary expressions
// TODO: binary operators without left hand operands , report error but continue passing
// Rule for parsing expressions into trees

func (tkn Tokens) expression() (Exp, error) {
	return tkn.list()
}
func (tkn Tokens) list() (Exp, error) {
	exp, err := tkn.asignment()
	if err != nil {
		return nil, err
	}
	if tkn[current].Ttype == scanner.COMMA {
		exps := []Exp{}
		exps = append(exps, exp)
		for {
			if tkn[current].Ttype == scanner.COMMA {
				current++
				exp, err := tkn.asignment()
				if err != nil {
					return nil, err
				}
				exps = append(exps, exp)
			} else {
				//end of list
				return List{Expressions: exps}, nil
			}
		}
	}
	return exp, nil
}

// TODO: add right associative parsing
func (tkn Tokens) asignment() (Exp, error) {
	exp, err := tkn.logicOr()
	if err != nil {
		return nil, err
	}
	//optionally expect "="
	if tkn[current].Ttype == scanner.EQUAL {
		//assigment
		op := tkn[current]
		//currently lv are single nodes
		lv, isStorageTarget := exp.(Primary)
		if isStorageTarget && lv.Node.Ttype == scanner.IDENTIFIER {
			//consume "="
			current++
			rv, err := tkn.asignment()
			if err != nil {
				return nil, err
			}
			ass := Assignment{lv, op, rv}
			return ass, nil
		}
		//TODO: add print method to exp interface , i could have called it here. cool right. Maybe add this later
		return nil, fmt.Errorf("Cannot use the l-value as storage target ")
	}
	return exp, nil
}

func (tkn Tokens) logicOr() (Exp, error) {
	expleft, err := tkn.logicAnd()
	opsToMatch := []scanner.Tokentype{
		scanner.OR,
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
				expright, err := tkn.logicAnd()
				if err != nil {
					return nil, err
				}
				expleft = LogicalOr{Left: expleft, Operator: op, Right: expright}
				continue Matching_Loop
			}
		}
		break
	}
	return expleft, nil
}
func (tkn Tokens) logicAnd() (Exp, error) {
	expleft, err := tkn.equality()
	opsToMatch := []scanner.Tokentype{
		scanner.AND,
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
				expright, err := tkn.equality()
				if err != nil {
					return nil, err
				}
				expleft = LogicalAnd{Left: expleft, Operator: op, Right: expright}
				continue Matching_Loop
			}
		}
		break
	}
	return expleft, nil
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
			cexpleft = Binary{Left: cexpleft, Operator: op, Right: cexpright}
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
				texpleft = Binary{Left: texpleft, Operator: op, Right: texpright}
				continue Matching_Loop
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
				fexpleft = Binary{Left: fexpleft, Operator: op, Right: fexpright}
				continue Matching_Loop
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
				uexpleft = Binary{Left: uexpleft, Operator: op, Right: fexpright}
				continue Matching_Loop
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
		return Unary{Operator: op, Right: uexp}, nil
	}
	return tkn.call()
}

// TODO: support anonymous functions
func (tkn Tokens) call() (Exp, error) {
	callee, err := tkn.primary()
	opsToMatch := scanner.LEFT_PAREN
	if err != nil {
		return nil, err
	}
Matching_Loop:
	for {
		cToken := tkn[current]
		// find the operator terminal
		if cToken.Ttype == opsToMatch {
			//consume operator terminal
			current++
			op := cToken
			if tkn[current].Ttype == scanner.RIGHT_PAREN {
				//empty
				callee = Call{Arrity: 0, Callee: callee, Operator: op, Args: nil}
				current++
				continue
			}
			//non-empty parameters
			// TODO:parse args as statements and expect funcDef,expression statements
			// if funcDef, add just add the name else add the expression
			exp, err := tkn.list()
			if err != nil {
				return nil, err
			}
			if tkn[current].Ttype != scanner.RIGHT_PAREN {
				return nil, fmt.Errorf("ParseError: Expected a right paren but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
			}
			args := exp
			list, isListExp := exp.(List)
			//TODO: assert max args
			//TODO:check function arrity here
			arrity := 1
			if isListExp {
				arrity = len(list.Expressions)
			}
			callee = Call{Arrity: arrity, Callee: callee, Operator: op, Args: args}
			current++
			continue Matching_Loop
		}
		break
	}
	return callee, nil
}

// rule for producing operands
func (tkn Tokens) primary() (Exp, error) {
	ttype := tkn[current].Ttype
	tnode := Primary{}
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
	tnode.Node = tkn[current]
	current++
	return tnode, nil
}
