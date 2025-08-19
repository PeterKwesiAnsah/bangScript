package parser

import (
	"bangScript/gbs/scanner"
	"fmt"
	"strconv"
)

type bsReturn struct {
	value Obj
}

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
	Evaluate(env *Stmtsenv) (Obj, error)
	//print() string
}

type Stmt interface {
	Execute(env *Stmtsenv) error
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
var mode uint8 = REPL

func (t *Stmtsenv) Get(depth int) *Stmtsenv {
	for range depth - 1 {
		t = t.Encloser
	}
	return t
}
func (t *ForStmt) StaticToDynamic(parent *Stmtsenv) error {
	//t.stmt.Env is environment for condition,initializer and single body statement
	if t.Stmt.Env == parent {
		return fmt.Errorf("Child environment can not be the same as Parent")
	}
	if parent.Policy == DYNAMIC {
		//condition,initializer and single body statement
		newEnv := &Stmtsenv{Local: map[string]Obj{}, Encloser: parent, Policy: DYNAMIC}
		if t.Stmt.Env == t.Stmt.Body.Env {
			t.Stmt.Env = newEnv
			t.Stmt.Body.Env = newEnv
		} else {
			t.Stmt.Env = newEnv
			t.Stmt.Body.StaticToDynamic(newEnv)
		}
	} else {
		if t.Stmt.Env.Encloser != parent {
			return fmt.Errorf("Body statement environment should encloses around the env passed to it ")
		}
	}
	return nil
}
func (t *FuncDef) StaticToDynamic(parent *Stmtsenv) error {
	if t.Body.Env == parent {
		return fmt.Errorf("Child environment can not be the same as Parent")
	}
	if parent.Policy == DYNAMIC {
		t.Body.Env = &Stmtsenv{Local: map[string]Obj{}, Encloser: parent, Policy: DYNAMIC}
	} else {
		if t.Body.Env.Encloser != parent {
			return fmt.Errorf("ExecutionError: Body statement environment should encloses around the env passed to it ")
		}
	}
	return nil
}
func (t *BlockStmt) StaticToDynamic(parent *Stmtsenv) error {
	if t.Env == parent {
		return fmt.Errorf("Child environment can not be the same as Parent")
	}
	if parent.Policy == DYNAMIC {
		t.Env = &Stmtsenv{Local: map[string]Obj{}, Encloser: parent, Policy: DYNAMIC}
	} else {
		if t.Env.Encloser != parent {
			return fmt.Errorf("Body statement environment should encloses around the env passed to it ")
		}
	}
	return nil
}
func (t *WhileStmt) StaticToDynamic(parent *Stmtsenv) error {
	if t.Env == parent {
		return fmt.Errorf("Child environment can not be the same as Parent")
	}
	if parent.Policy == DYNAMIC {
		//condition,initializer and single body statement
		newEnv := &Stmtsenv{Local: map[string]Obj{}, Encloser: parent, Policy: DYNAMIC}
		if t.Env == t.Body.Env {
			t.Env = newEnv
			t.Body.Env = newEnv
		} else {
			t.Env = newEnv
			t.Body.StaticToDynamic(newEnv)
		}
	} else {
		if t.Env.Encloser != parent {
			return fmt.Errorf("Body statement environment should encloses around the env passed to it ")
		}
	}
	return nil
}

func (t ReturnStmt) Execute(env *Stmtsenv) error {
	value, err := t.exp.Evaluate(env)
	if err != nil {
		return err
	}
	panic(bsReturn{
		value: value,
	})
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
		exp: exp,
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

// evaluate to "<fn funcname >"
// func (t FuncDef) Evaluate(env *Stmtsenv) (Obj, error)
// bind function name to it's value in env.
func (t FuncDef) Execute(env *Stmtsenv) error {
	bs := t.Body
	bs.Env.Encloser.Local[t.Name.Lexem] = t
	return nil
}

func (t FuncDef) call(env *Stmtsenv, callStack *CallStack, callInfo *Call) (value Obj, err error) {
	bs := t.Body

	defer func() {
		if r := recover(); r != nil {
			switch s := r.(type) {
			case bsReturn:
				value = s.value
			default:
				panic(r)
			}
		}
	}()

	//bs.Env need to be a complete copy, since it's new, any environment that once enclosed around bs.Env will be updated to support closures
	newEnv := Stmtsenv{Local: map[string]Obj{}, Encloser: bs.Env.Encloser, Policy: DYNAMIC}
	bs.Env = &newEnv

	envWithFunctionArgsOnly := newEnv.Local
	if callInfo.Args != nil {
		listArgs, isArgs := callInfo.Args.(List)
		if isArgs {
			for argI, exp := range listArgs.Expressions {
				value, err := exp.Evaluate(env)
				if err != nil {
					return nil, err
				}
				envWithFunctionArgsOnly[t.Params[argI].Lexem] = value
			}
		} else {
			value, err := callInfo.Args.Evaluate(env)
			if err != nil {
				return nil, err
			}
			envWithFunctionArgsOnly[t.Params[0].Lexem] = value
		}
	}
	cs = append(cs, &CallDetails{
		args: callInfo.Args,
		at:   callInfo.Operator.Line,
		name: t.Name.Lexem,
	})
	//the reason why we pass nil is that , we have already created the environments during parsing and each block has it
	// what we make fresh copies of environment during function calls, and update subsequent environments in the function as they too are dynamic
	//dynamic env
	err = bs.Execute(nil)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (t ForStmt) Execute(parent *Stmtsenv) error {
	return t.Stmt.Execute(nil)
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

// TODO: implements loop keyword like continue,break
func (t WhileStmt) Execute(env *Stmtsenv) error {
	var executionErr error
	var evalErr error
	//env should be nil
	if env != nil {
		return fmt.Errorf("Block statements does not need the caller's env")
	}
	if t.Init != nil {
		//init is neither a while stmt or a block stmt
		executionErr = t.Init.Execute(t.Env)
		if executionErr != nil {
			return executionErr
		}
	}
	goto evaluateAndtest
executeBody:
	{
		//creating new environment per iteration??
		executionErr = t.Body.Execute(nil)
		if executionErr != nil {
			return executionErr
		}
		goto evaluateAndtest
	}
evaluateAndtest:
	//infinite loop equivalent to while(true)
	if t.Condition == nil {
		goto executeBody
	}
	obj, evalErr := t.Condition.Evaluate(t.Env)
	if evalErr != nil {
		return evalErr
	}
	isTruth := isTruthy(obj)
	if isTruth {
		goto executeBody
	}
	return nil
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

func (t IfStmt) Execute(env *Stmtsenv) error {

	obj, err := t.Condition.Evaluate(env)
	if err != nil {
		return err
	}
	isTruth := isTruthy(obj)

	if isTruth {
		switch s := t.Thenbody.(type) {
		case WhileStmt:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case BlockStmt:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case FuncDef:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case ForStmt:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		default:
			err = s.Execute(env)
		}
		if err != nil {
			return err
		}
	} else {
		if t.Elsebody == nil {
			return nil
		}
		switch s := t.Elsebody.(type) {
		case WhileStmt:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case BlockStmt:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case FuncDef:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(env)
		case ForStmt:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(env)
		default:
			err = s.Execute(env)
		}
		if err != nil {
			return err
		}
	}
	return nil
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

func (t BlockStmt) Execute(env *Stmtsenv) error {

	//isDynamic := env != nil
	for _, stmt := range t.Stmts {
		if stmt == nil {
			continue
		}
		//functions have dynamic environment,created when they are called as such nested environments should be updated to enclose around this new env before they are executed
		var err error
		switch s := stmt.(type) {
		case WhileStmt:
			err = s.StaticToDynamic(t.Env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case BlockStmt:
			err = s.StaticToDynamic(t.Env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case FuncDef:
			err = s.StaticToDynamic(t.Env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case ForStmt:
			err = s.StaticToDynamic(t.Env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		default:
			err = s.Execute(t.Env)
		}
		if err != nil {
			return err
		}
	}
	return nil
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

func (t PrintStmt) Execute(env *Stmtsenv) error {
	Obj, err := t.Exp.Evaluate(env)
	if err != nil {
		return err
	}
	//TODO:function def don't evaluate to simple scalar values like integer,bool etc it evaluates to a struct which won't print nicely
	fmt.Printf("%v\n", Obj)
	return nil
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

func (t VarStmt) Execute(env *Stmtsenv) error {
	//check if variable declaration have a definition
	if t.Exp != nil {
		switch s := t.Exp.(type) {
		case Primary:
			env.Local[s.Node.Lexem] = nil
		case Assignment:
			//for now only primary identifier expressions map to a storage location
			lv, isStorageTarget := s.StoreTarget.(Primary)
			if !(isStorageTarget && lv.Node.Ttype == scanner.IDENTIFIER) {
				return fmt.Errorf("Cannot use the l-value as storage target")
			}
			obj, err := s.Right.Evaluate(env)
			if err != nil {
				return err
			}
			env.Local[lv.Node.Lexem] = obj
		case List:
			for _, exp := range s.Expressions {
				err := VarStmt{exp}.Execute(env)
				if err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("Not Valid expression for variable declaration")
		}
	}
	return nil
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

func (t ExpStmt) Execute(env *Stmtsenv) error {
	obj, err := t.Exp.Evaluate(env)
	if err != nil {
		return err
	}
	if mode == REPL && env.Encloser == nil {
		fmt.Printf("%v\n", obj)
		return nil
	}
	return err
}
func (tkn Tokens) expStmt() (Stmt, error) {
	stmt := ExpStmt{}
	exp, err := tkn.expression()
	if err != nil {
		return nil, err
	}
	stmt.Exp = exp
	//expect a ";" terminator only for script mode
	if mode == REPL {
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
	stmts := []Stmt{}
	mode = m
	for tkn[current].Ttype != scanner.EOF {
		stmt, err := tkn.declarations(globalEnv)
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

// implement the Exp interface
func (p Primary) Evaluate(env *Stmtsenv) (Obj, error) {
	//for evaluting expressions at compile-time we can perform mathematical operations,logical operations and string concantenation
	// operands needs to be only following string , number and boolean
	switch p.Node.Ttype {
	case scanner.NUMBER:
		op, err := strconv.ParseFloat(p.Node.Lexem, 64)
		if err != nil {
			//handle error for failed type conversion
			return nil, err
		}
		return op, nil
	//string concantenation and comparison
	case scanner.STRING:
		return p.Node.Lexem, nil
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
				Obj, itExist := cur.Local[p.Node.Lexem]
				if itExist {
					return Obj, nil
				}
				cur = cur.Encloser
			}
			return nil, fmt.Errorf("Variable %s is undefined at line %d", p.Node.Lexem, p.Node.Line)
		}
	default:
		return nil, fmt.Errorf("Expected a string, number, a target location , nil and boolean but got %d at line %d", p.Node.Ttype, p.Node.Line)
	}
}

func (u Unary) Evaluate(env *Stmtsenv) (Obj, error) {
	Exp, err := u.Right.Evaluate(env)
	if err != nil {
		return nil, err
	}
	operator := u.Operator.Ttype

	if operator == scanner.BANG {
		bol, isbol := Exp.(bool)
		if isbol {
			return !bol, nil
		}
		return nil, fmt.Errorf("Expected a boolean value but got something else at line %d", u.Operator.Line)
	} else if operator == scanner.MINUS {
		num, isnum := Exp.(float64)
		if isnum {
			return -num, nil
		}
		return nil, fmt.Errorf("Expected a number value but got something else at line %d", u.Operator.Line)
	}
	return nil, fmt.Errorf("Invalid expression")
}

func (b Binary) Evaluate(env *Stmtsenv) (Obj, error) {
	left, err := b.Left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	right, err := b.Right.Evaluate(env)
	if err != nil {
		return nil, err
	}
	// TODO: need to handle nil expressions nil + string or string + nil or nil + 1
	switch b.Operator.Ttype {
	case scanner.PLUS:
		{
			//string concatenation
			// TODO: (parse and concatenate) string + number | number + string
			strLeft, okLeft := left.(string)
			if okLeft {
				strRight, okRight := right.(string)
				if !okRight {
					return nil, fmt.Errorf("Invalid Right operand, expected a string at line %d", b.Operator.Line)
				}
				return strLeft + strRight, nil
			}
			// integer addition
			floatLeft, okLeft := left.(float64)
			if okLeft {
				floatRight, okRight := right.(float64)
				if !okRight {
					return nil, fmt.Errorf("Invalid Right operand, expected a number at line %d", b.Operator.Line)
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
			return nil, fmt.Errorf("Invalid operator at line %d.", b.Operator.Line)
		}

	}
	return nil, fmt.Errorf("Invalid expression.")
}
func (t LogicalAnd) Evaluate(env *Stmtsenv) (Obj, error) {
	objL, err := t.Left.Evaluate(env)

	if err != nil {
		return nil, err
	}
	isTrueL := isTruthy(objL)
	//short circuit
	if !isTrueL {
		return objL, nil
	}
	objR, err := t.Right.Evaluate(env)
	if err != nil {
		return nil, err
	}
	isTrueR := isTruthy(objR)
	if isTrueR {
		return objR, nil
	}
	return objL, nil
}
func (t LogicalOr) Evaluate(env *Stmtsenv) (Obj, error) {
	objL, err := t.Left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	isTrueL := isTruthy(objL)

	//short circuit
	if isTrueL {
		return objL, nil
	}
	objR, err := t.Right.Evaluate(env)

	if err != nil {
		return nil, err
	}
	return objR, nil
}
func (a Assignment) Evaluate(env *Stmtsenv) (Obj, error) {
	cur := env
	lv, isStorageTarget := a.StoreTarget.(Primary)
	if !(isStorageTarget && lv.Node.Ttype == scanner.IDENTIFIER) {
		return nil, fmt.Errorf("Cannot use the l-value as storage target")
	}
	for cur != nil {
		_, itExist := cur.Local[lv.Node.Lexem]
		if itExist {
			rv, err := a.Right.Evaluate(env)
			if err != nil {
				return nil, err
			}
			cur.Local[lv.Node.Lexem] = rv
			return rv, nil
		}
		cur = cur.Encloser
	}
	return nil, fmt.Errorf("Undefined variable at line %d", a.Operator.Line)
}
func (l List) Evaluate(env *Stmtsenv) (Obj, error) {
	var rvalue Obj
	for index, exp := range l.Expressions {
		value, err := exp.Evaluate(env)
		if err != nil {
			return nil, err
		}
		if (index + 1) == len(l.Expressions) {
			rvalue = value
		}
	}
	return rvalue, nil
}
func (t Call) Evaluate(env *Stmtsenv) (Obj, error) {
	value, err := t.Callee.Evaluate(env)
	if err != nil {
		return nil, err
	}
	function, isCallable := value.(FuncDef)
	if !isCallable {
		return nil, fmt.Errorf("Can't call expression at line %d", t.Operator.Line)
	}
	if t.Arrity != function.Arrity {
		return nil, fmt.Errorf("Expected %d arguments but got %d instead", function.Arrity, t.Arrity)
	}
	//env is a parent environment can be immediate env or global
	return function.call(env, nil, &t)
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

func isTruthy(val Obj) bool {
	isTruthy := true
	var falsyVal []Obj = []Obj{"", nil, 0, false}
	for _, falsy := range falsyVal {
		if falsy == val {
			isTruthy = false
			break
		}
	}
	return isTruthy
}
