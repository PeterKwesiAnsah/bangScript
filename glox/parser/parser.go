package parser

import (
	"fmt"
	"lox/glox/scanner"
	"strconv"
)

// statements that carry their own env, whileStmt,blockStmt,funDef (but through their body)
type Tokens []*scanner.Token
type Obj interface{}

type Exp interface {
	//should take in an environment or context
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

// callStack grows , so do env
type Callable interface {
	//we create new copies of ...env.Local=map[string]Obj{}
	// reference of body.env still stays the same and has access to it's parent environment
	// we are keeping body.env as reference but cleaning the slate anytime we have a function call
	// x86-64 abi, the same set of cpu registers are used to hold function arguments but are saved first before a function so the caller's state or arguments does not lost and returned back
	// when the callee returns so perhaps we replicate this behavior
	//
	call(*CallStack, *call) (Obj, error)
}

type returnStmt struct {
	//currently we allow returns to exps for now. In the future we change to statment because of closures
	exp Exp
}

// implements Callable interface and statement
// before we call .call() caller must set up the arguments in it's environment
type funcDef struct {
	name   *scanner.Token
	params []*scanner.Token
	body   BlockStmt
	arrity int
}

type Stmtsenv struct {
	Local    map[string]Obj
	Encloser *Stmtsenv
}

type call struct {
	callee   Exp
	operator *scanner.Token
	arrity   int
	args     Exp
}

type list struct {
	expressions []Exp
}

type assigment struct {
	//l-value
	storeTarget Exp
	operator    *scanner.Token
	//r-value
	right Exp
}

type logicalOr struct {
	left     Exp
	operator *scanner.Token
	right    Exp
}
type logicalAnd struct {
	left     Exp
	operator *scanner.Token
	right    Exp
}

type binary struct {
	left     Exp
	operator *scanner.Token
	right    Exp
}

// TODO: for handling conditional operations
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
type forStmt struct {
	stmt WhileStmt
}
type WhileStmt struct {
	condition Exp
	body      BlockStmt
	init      Stmt
	env       *Stmtsenv
}

type BlockStmt struct {
	stmts []Stmt
	env   *Stmtsenv
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
var cs CallStack = CallStack{}

func (t returnStmt) Execute(env *Stmtsenv) error {
	value, err := t.exp.Evaluate(env)
	if err != nil {
		return err
	}
	panic(value)
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
	return returnStmt{
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
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	if tkn[current].Ttype != scanner.RIGHT_BRACE {
		return nil, fmt.Errorf("ParseError: Expected a right brace but got EOF at line %d", tkn[current].Line)
	}
	current++
	stmt.stmts = stmts
	stmt.env = &inner

	return funcDef{
		arrity: len(params),
		name:   name,
		body:   stmt,
		params: params,
	}, nil
}

// evaluate to "<fn funcname >"
// func (t funcDef) Evaluate(env *Stmtsenv) (Obj, error)
// bind function name to it's value in env.
func (t funcDef) Execute(env *Stmtsenv) error {
	bs := t.body
	if bs.env.Encloser != env {
		return fmt.Errorf("ExecutionError: Body statement environment should encloses around the env passed to it ")
	}
	env.Local[t.name.Lexem] = t
	return nil
}

func (t funcDef) call(env *Stmtsenv, callStack *CallStack, callInfo *call) (value Obj, err error) {
	bs := t.body

	defer func() {
		if r := recover(); r != nil {
			value = r
		}
	}()
	//bs.env need to be a complete copy, since it's new any environment that once enclosed around bs.env will be updated
	newEnv := Stmtsenv{Local: map[string]Obj{}, Encloser: bs.env.Encloser}
	bs.env = &newEnv

	envWithFunctionArgsOnly := newEnv.Local
	if callInfo.args != nil {
		listArgs, isArgs := callInfo.args.(list)
		if isArgs {
			for argI, exp := range listArgs.expressions {
				value, err := exp.Evaluate(env)
				if err != nil {
					return nil, err
				}
				envWithFunctionArgsOnly[t.params[argI].Lexem] = value
			}
		} else {
			value, err := callInfo.args.Evaluate(env)
			if err != nil {
				return nil, err
			}
			envWithFunctionArgsOnly[t.params[0].Lexem] = value
		}
	}
	cs = append(cs, &CallDetails{
		args: callInfo.args,
		at:   callInfo.operator.Line,
		name: t.name.Lexem,
	})
	//the reason why we pass nil is that , we have already created the environments during parsing and each block has it
	// what we make fresh copies of environment during function calls, and update subsequent environments in the function as they too are dynamic

	err = bs.Execute(nil)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (t forStmt) Execute(env *Stmtsenv) error {
	t.stmt.env.Encloser = env
	return t.stmt.Execute(nil)
}

// Implementing a for loop using a while loop automatically , creates a block scope where the initializer sits
func (tkn Tokens) forStmt(env *Stmtsenv) (Stmt, error) {
	var initializer Stmt = nil
	var condition Exp = nil
	var sideEffect Exp = nil
	var body Stmt = nil
	var err error
	//create scope for initializer
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

	// empty body
	if tkn[current].Ttype == scanner.SEMICOLON {
		//initializer runs in the same scope as the condition
		return forStmt{stmt: WhileStmt{condition: condition, init: initializer,
			body: BlockStmt{
				stmts: []Stmt{BlockStmt{stmts: []Stmt{expStmt{exp: sideEffect}}, env: &whileScope}},
				env:   &whileScope,
			}, env: &whileScope}}, nil
	}
	body, err = tkn.declarations(&whileScope)
	if err != nil {
		return nil, err
	}
	bs, isBS := body.(BlockStmt)
	if isBS {
		stmts := bs.stmts
		stmts = append(stmts, expStmt{exp: sideEffect})
		//we should expect a scope already created
		return forStmt{stmt: WhileStmt{condition: condition, init: initializer, body: BlockStmt{
			//initializer runs in the same scope as the condition
			stmts: []Stmt{BlockStmt{stmts: stmts, env: bs.env}},
			env:   &whileScope,
		}, env: &whileScope}}, nil
	} else {
		return forStmt{stmt: WhileStmt{condition: condition, init: initializer, body: BlockStmt{
			//initializer runs in the same scope as the condition
			stmts: []Stmt{
				BlockStmt{stmts: []Stmt{body, expStmt{exp: sideEffect}},
					env: &whileScope}},
			env: &whileScope,
		}, env: &whileScope}}, nil
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
	if t.init != nil {
		//init is neither a while stmt or a block stmt
		executionErr = t.init.Execute(t.env)
		if executionErr != nil {
			return executionErr
		}
	}
	goto evaluateAndtest
executeBody:
	{
		body := t.body.stmts[0]
		if body == nil {
			goto evaluateAndtest
		}
		bs, isbs := body.(BlockStmt)
		if !isbs {
			return fmt.Errorf("Not block")
		}
		//TODO:assert bs.env==t.env, throw error if otherwise , only block statments in a while statement can have the same env as it's parent
		//bs.env and t.env are not aliases, otherwise we would create an infinite loop
		if bs.env != t.env {
			bs.env.Encloser = t.env
		}
		executionErr = t.body.Execute(nil)

		if executionErr != nil {
			return executionErr
		}
		goto evaluateAndtest
	}
evaluateAndtest:
	//infinite loop equivalent to while(true)
	if t.condition == nil {
		goto executeBody
	}
	obj, evalErr := t.condition.Evaluate(t.env)
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
	return WhileStmt{condition: cond, body: bs, env: Encloser}, nil
}

func (t ifStmt) Execute(env *Stmtsenv) error {

	obj, err := t.condition.Evaluate(env)
	if err != nil {
		return err
	}
	isTruth := isTruthy(obj)

	if isTruth {
		//these are statements that need their environment updated becuase they either carry own env or have nodes that carry their own env
		switch s := t.thenbody.(type) {
		case WhileStmt:
			if s.env == env {
				return fmt.Errorf("Child environment can not be the same as Parent")
			}
			//This means, from top to bottom we are replacing all static env during function calls
			// TODO: do this for only func call
			s.env = &Stmtsenv{Local: map[string]Obj{}, Encloser: nil}
			s.env.Encloser = env
			err = s.Execute(nil)
		case BlockStmt:
			if s.env == env {
				return fmt.Errorf("Child environment can not be the same as Parent")
			}
			s.env = &Stmtsenv{Local: map[string]Obj{}, Encloser: nil}
			s.env.Encloser = env
			err = s.Execute(nil)
		case funcDef:
			if s.body.env == env {
				return fmt.Errorf("Child environment can not be the same as Parent")
			}
			s.body.env = &Stmtsenv{Local: map[string]Obj{}, Encloser: nil}
			s.body.env.Encloser = env
			//TODO: for now we pass env, later we will  pass nil as functions are just block statements with dynamic env
			err = s.Execute(env)
		//case returnStmt:
		default:
			err = s.Execute(env)
		}
		if err != nil {
			return err
		}
	} else {
		if t.elsebody == nil {
			return nil
		}
		//these are statements that need their environment updated
		switch s := t.elsebody.(type) {
		case WhileStmt:
			if s.env == env {
				return fmt.Errorf("Child environment can not be the same as Parent")
			}
			s.env = &Stmtsenv{Local: map[string]Obj{}, Encloser: nil}
			s.env.Encloser = env
			err = s.Execute(nil)
		case BlockStmt:
			if s.env == env {
				return fmt.Errorf("Child environment can not be the same as Parent")
			}
			s.env = &Stmtsenv{Local: map[string]Obj{}, Encloser: nil}
			s.env.Encloser = env
			err = s.Execute(nil)
		case funcDef:
			if s.body.env == env {
				return fmt.Errorf("Child environment can not be the same as Parent")
			}
			s.body.env = &Stmtsenv{Local: map[string]Obj{}, Encloser: nil}
			s.body.env.Encloser = env
			//TODO: for now we pass env, later we will  pass nil as functions are just block statements with dynamic env
			err = s.Execute(env)
		//case returnStmt:
		default:
			//for statements that have their own env,statement.Env you will be executed with nil
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
	return ifStmt{condition: condexp, thenbody: stmtBody, elsebody: elseStmt}, nil
}

func (t BlockStmt) Execute(env *Stmtsenv) error {

	for _, stmt := range t.stmts {
		if stmt == nil {
			continue
		}
		//functions have dynamic environment,created when they are called as such nested environments should be updated to enclose around this new env before they are executed
		var err error
		switch s := stmt.(type) {
		case WhileStmt:
			if s.env == t.env {
				return fmt.Errorf("Child environment can not be the same as Parent")
			}
			s.env = &Stmtsenv{Local: map[string]Obj{}, Encloser: nil}
			///if s.env.Encloser != t.env , t.env is dynamic.
			s.env.Encloser = t.env
			err = s.Execute(nil)
		case BlockStmt:
			if s.env == t.env {
				return fmt.Errorf("Child environment can not be the same as Parent")
			}
			s.env = &Stmtsenv{Local: map[string]Obj{}, Encloser: nil}
			///if s.env.Encloser != t.env , t.env is dynamic.
			s.env.Encloser = t.env
			err = s.Execute(nil)
		case funcDef:
			if s.body.env == t.env {
				return fmt.Errorf("Child environment can not be the same as Parent")
			}
			s.body.env = &Stmtsenv{Local: map[string]Obj{}, Encloser: nil}
			s.body.env.Encloser = t.env
			//TODO: for now we pass env, later we will  pass nil as functions are just block statements with dynamic env
			err = s.Execute(t.env)
		//case returnStmt:
		default:
			//for statements that have their own env,statement.Env you will be executed with nil
			err = s.Execute(t.env)
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
		// if stmt is continue/break/return or any other loop related statement we throw error because loops will now handle parsing it's body
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	if tkn[current].Ttype != scanner.RIGHT_BRACE {
		return nil, fmt.Errorf("ParseError: Expected a right brace but got EOF at line %d", tkn[current].Line)
	}
	current++
	stmt.stmts = stmts
	stmt.env = &inner
	return stmt, nil
}

func (t printStmt) Execute(env *Stmtsenv) error {
	Obj, err := t.exp.Evaluate(env)
	if err != nil {
		return err
	}
	//check if Obj is a funcDef, if true we call .Evaluate on it again
	//TODO:function def don't evaluate to simple scalar values like integer,bool etc it evaluates to a struct which won't print nicely
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
		return nil, fmt.Errorf("Expected semi-colon but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
	return stmt, nil
}

func (t varStmt) Execute(env *Stmtsenv) error {
	var obj Obj = nil
	//check if variable declaration have a definition
	if t.exp != nil {
		var err error
		obj, err = t.exp.Evaluate(env)
		if err != nil {
			return err
		}
	}
	env.Local[t.name.Lexem] = obj
	return nil
}
func (tkn Tokens) varStmt() (Stmt, error) {
	stmt := varStmt{}
	//expect identifier
	if tkn[current].Ttype != scanner.IDENTIFIER {
		return nil, fmt.Errorf("Expected an identifier after var but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
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
		return nil, fmt.Errorf("Expected semi-colon but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
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
		return nil, fmt.Errorf("Expected semi-colon but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
	}
	current++
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
// statements->printStmt | blockStmt | expressionStmt | ifStmt
// "if(exp)
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
			return nil, fmt.Errorf("Variable %s is undefined at line %d", p.node.Lexem, p.node.Line)
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
	// need to handle nil expressions nil + string or string + nil or nil + 1
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
func (t logicalAnd) Evaluate(env *Stmtsenv) (Obj, error) {
	objL, err := t.left.Evaluate(env)

	if err != nil {
		return nil, err
	}
	isTrueL := isTruthy(objL)
	//short circuit
	if !isTrueL {
		return objL, nil
	}
	objR, err := t.right.Evaluate(env)
	if err != nil {
		return nil, err
	}
	isTrueR := isTruthy(objR)
	if isTrueR {
		return objR, nil
	}
	return objL, nil
}
func (t logicalOr) Evaluate(env *Stmtsenv) (Obj, error) {
	objL, err := t.left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	isTrueL := isTruthy(objL)

	//short circuit
	if isTrueL {
		return objL, nil
	}
	objR, err := t.right.Evaluate(env)

	if err != nil {
		return nil, err
	}
	return objR, nil
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
func (l list) Evaluate(env *Stmtsenv) (Obj, error) {
	var rvalue Obj
	for index, exp := range l.expressions {
		value, err := exp.Evaluate(env)
		if err != nil {
			return nil, err
		}
		if (index + 1) == len(l.expressions) {
			rvalue = value
		}
	}
	return rvalue, nil
}
func (t call) Evaluate(env *Stmtsenv) (Obj, error) {
	value, err := t.callee.Evaluate(env)
	if err != nil {
		return nil, err
	}
	function, isCallable := value.(funcDef)
	if !isCallable {
		return nil, fmt.Errorf("Can't call expression at line %d", t.operator.Line)
	}
	if t.arrity != function.arrity {
		return nil, fmt.Errorf("Expected %d arguments but got %d instead", function.arrity, t.arrity)
	}
	//env is a parent environment can be immediate env or global
	return function.call(env, nil, &t)
}

// TODO: grammer for tenary expressions
// TODO: grammer for grouped expression
// TODO: implement grammer for logical operators && and ||
// TODO: binary operators without left hand operands , report error but continue passing
// Rule for parsing expressions into trees

func (tkn Tokens) expression() (Exp, error) {
	//
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
				return list{expressions: exps}, nil
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
		//if i had a print method to my exp interface , i could have called it here. cool right. Maybe add this later
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
				expleft = logicalOr{left: expleft, operator: op, right: expright}
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
				expleft = logicalAnd{left: expleft, operator: op, right: expright}
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
				fexpleft = binary{left: fexpleft, operator: op, right: fexpright}
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
		return unary{operator: op, right: uexp}, nil
	}
	return tkn.call()
}
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
				callee = call{arrity: 0, callee: callee, operator: op, args: nil}
				current++
				continue
			}
			//non-empty parameters
			exp, err := tkn.list()
			if err != nil {
				return nil, err
			}
			if tkn[current].Ttype != scanner.RIGHT_PAREN {
				return nil, fmt.Errorf("ParseError: Expected a right paren but got %d at line %d", tkn[current].Ttype, tkn[current].Line)
			}
			args := exp
			list, isListExp := exp.(list)
			//TODO: assert max args
			//TODO:check function arrity here
			arrity := 1
			if isListExp {
				arrity = len(list.expressions)
			}
			callee = call{arrity: arrity, callee: callee, operator: op, args: args}
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

/////////////////////////////////////////////////////end of expression rules
