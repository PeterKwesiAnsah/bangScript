package resolver

import (
	"bangScript/gbs/parser"
	"bangScript/gbs/scanner"
	"fmt"
	"strconv"
)

type BangScriptReturn struct {
	value parser.Obj
}

const (
	REPL uint8 = iota
	SCRIPT
)

const (
	STATIC uint8 = iota
	DYNAMIC
)

type VariableMetaData struct {
	isUsed     bool
	isResolved bool
}

const (
	//Define ResolverFunctions
	NONE uint32 = iota
	VAR_STMT
	FUNC_DEF_STMT
	FOR_STMT
	WHILE_STMT
	PRINT_STMT
	RETURN_STMT
	EXP_STMT
	DESTINATION_ASSIGNMENT
)

var TopOfCallStack uint32 = NONE

// TODO: Handle loop constructs/function specific keywords like return,break,continue
func ResolveVarStmt(t parser.VarStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	caller := TopOfCallStack
	TopOfCallStack = VAR_STMT
	switch expT := t.Exp.(type) {
	case parser.Primary:
		if expT.Node.Ttype != scanner.IDENTIFIER {
			return nil, fmt.Errorf("Expected an identifier but got something else %d", expT.Node.Ttype)
		}
		lexeme := expT.Node.Lexem
		_, itExists := env.Local[lexeme]
		if itExists {
			return nil, fmt.Errorf("Cannot redeclare a variable in the same scope.")
		}
		//is initializer present??
		env.Local[lexeme] = VariableMetaData{isUsed: false, isResolved: true}
		TopOfCallStack = caller
		return ResolvedVarStmt{Exp: ResolvedPrimary{Node: expT.Node, ScopeDepth: -1}}, nil
	case parser.Assignment:
		lv := expT.StoreTarget.(parser.Primary)
		resolvedExpr, err := ResolveExpr(expT.Right, env)
		if err != nil {
			return nil, err
		}
		env.Local[lv.Node.Lexem] = VariableMetaData{isUsed: false, isResolved: true}
		TopOfCallStack = caller
		return ResolvedVarStmt{Exp: ResolvedAssignment{StoreTarget: ResolvedPrimary{Node: lv.Node, ScopeDepth: -1}, Right: resolvedExpr}}, nil
	case parser.List:
		resolvedExpr, err := ResolveList(expT, env)
		if err != nil {
			return nil, err
		}
		TopOfCallStack = caller
		return ResolvedVarStmt{Exp: resolvedExpr}, nil
	default:
		return nil, fmt.Errorf("Invalid expression type")
	}
}

func ResolveBlockStmt(t parser.BlockStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	resolvedStmts := []ResolvedStmt{}
	for _, stmt := range t.Stmts {
		resolvedStmt, err := ResolveStmt(stmt, t.Env)
		if err != nil {
			return nil, err
		}
		resolvedStmts = append(resolvedStmts, resolvedStmt)
	}
	return ResolvedBlockStmt{Stmts: resolvedStmts, Env: t.Env}, nil
}

func ResolveFuncDef(t parser.FuncDef, env *parser.Stmtsenv) (ResolvedStmt, error) {
	caller := TopOfCallStack
	TopOfCallStack = FUNC_DEF_STMT

	for _, param := range t.Params {
		if param.Ttype != scanner.IDENTIFIER {
			return nil, fmt.Errorf("Function parameters need to be identifiers")
		}
		//references to the function params
		t.Body.Env.Local[param.Lexem] = VariableMetaData{isUsed: false, isResolved: true}
	}
	resolvedStmt, err := ResolveBlockStmt(t.Body, nil)
	if err != nil {
		return nil, err
	}
	resolvedBs := resolvedStmt.(ResolvedBlockStmt)
	TopOfCallStack = caller
	t.Body.Env.Encloser.Local[t.Name.Lexem] = VariableMetaData{isUsed: false, isResolved: true}
	return ResolvedFuncDef{Name: t.Name, Params: t.Params, Body: resolvedBs, Arrity: t.Arrity}, nil
}

func ResolveReturnStmt(t parser.ReturnStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	if TopOfCallStack != FUNC_DEF_STMT {
		return nil, fmt.Errorf("Return statement can only be used inside a function body")
	}
	resolvedExpr, err := ResolveExpr(t.Exp, env)
	if err != nil {
		return nil, err
	}
	return ResolvedReturnStmt{Exp: resolvedExpr}, nil
}

func ResolveIfStmt(t parser.IfStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	var resolveElse ResolvedStmt
	var err error
	resolveCondition, err := ResolveExpr(t.Condition, env)
	if err != nil {
		return nil, err
	}
	resolveThen, err := ResolveStmt(t.Thenbody, env)
	if err != nil {
		return nil, err
	}
	if t.Elsebody != nil {
		resolveElse, err = ResolveStmt(t.Elsebody, env)
		if err != nil {
			return nil, err
		}
	}
	return ResolvedIfStmt{Condition: resolveCondition, Thenbody: resolveThen, Elsebody: resolveElse}, nil
}

func ResolveWhileStmt(t parser.WhileStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	var resolvedInit ResolvedStmt
	var resolvedCondition ResolvedExpr
	var resolvedBody ResolvedStmt
	var err error
	if t.Init != nil {
		resolvedInit, err = ResolveStmt(t.Init, t.Env)
	}
	if t.Condition != nil {
		resolvedCondition, err = ResolveExpr(t.Condition, t.Env)
	}
	resolvedBody, err = ResolveStmt(t.Body, nil)
	if err != nil {
		return nil, err
	}
	//runtime type checking...definitely incur some cycles
	resolvedBs := resolvedBody.(ResolvedBlockStmt)
	return ResolvedWhileStmt{Condition: resolvedCondition, Body: resolvedBs, Env: t.Env, Init: resolvedInit}, nil
}

func ResolveForStmt(t parser.ForStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	resolvedStmt, err := ResolveWhileStmt(t.Stmt, nil)
	if err != nil {
		return nil, err
	}
	resolvedWs := resolvedStmt.(ResolvedWhileStmt)
	return ResolvedForStmt{Stmt: resolvedWs}, nil
}

func ResolveExpStmt(t parser.ExpStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	resolvedExp, err := ResolveExpr(t.Exp, env)
	if err != nil {
		return nil, err
	}
	return ResolvedExpStmt{Exp: resolvedExp}, nil
}

func ResolvePrintStmt(t parser.PrintStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	resolvedExpr, err := ResolveExpr(t.Exp, env)
	if err != nil {
		return nil, err
	}
	return ResolvedPrintStmt{Exp: resolvedExpr}, nil
}

func ResolveAssignment(t parser.Assignment, env *parser.Stmtsenv) (ResolvedExpr, error) {
	caller := TopOfCallStack
	TopOfCallStack = DESTINATION_ASSIGNMENT
	resolveStorageTarget, err := ResolveExpr(t.StoreTarget, env)
	if err != nil {
		return nil, err
	}
	resolveValue, err := ResolveExpr(t.Right, env)
	if err != nil {
		return nil, err
	}
	TopOfCallStack = caller
	return ResolvedAssignment{StoreTarget: resolveStorageTarget, Right: resolveValue, Operator: t.Operator}, nil
}

func ResolveList(t parser.List, env *parser.Stmtsenv) (ResolvedExpr, error) {
	caller := TopOfCallStack
	resolvedExps := []ResolvedExpr{}
	for _, expr := range t.Expressions {
		if caller == VAR_STMT {
			varStmt := parser.VarStmt{Exp: expr}
			resolvedStmt, err := ResolveVarStmt(varStmt, env)
			if err != nil {
				return nil, err
			}
			resolvedVarStmt := resolvedStmt.(ResolvedVarStmt)
			resolvedExps = append(resolvedExps, resolvedVarStmt.Exp)
			continue
		}
		resolvedExpr, err := ResolveExpr(expr, env)
		if err != nil {
			return nil, err
		}
		resolvedExps = append(resolvedExps, resolvedExpr)
	}
	TopOfCallStack = caller
	return ResolvedList{Expressions: resolvedExps}, nil
}

func ResolveLogicalOr(t parser.LogicalOr, env *parser.Stmtsenv) (ResolvedExpr, error) {
	resolvedExpLeft, err := ResolveExpr(t.Left, env)
	if err != nil {
		return nil, err
	}
	resolvedExpRight, err := ResolveExpr(t.Right, env)
	if err != nil {
		return nil, err
	}
	return ResolvedLogicalOr{Left: resolvedExpLeft, Right: resolvedExpRight, Operator: t.Operator}, nil
}

func ResolveLogicalAnd(t parser.LogicalAnd, env *parser.Stmtsenv) (ResolvedExpr, error) {
	resolvedExpLeft, err := ResolveExpr(t.Left, env)
	if err != nil {
		return nil, err
	}
	resolvedExpRight, err := ResolveExpr(t.Right, env)
	if err != nil {
		return nil, err
	}
	return ResolvedLogicalAnd{Left: resolvedExpLeft, Right: resolvedExpRight, Operator: t.Operator}, nil
}

func ResolveUnary(t parser.Unary, env *parser.Stmtsenv) (ResolvedExpr, error) {
	resolvedExp, err := ResolveExpr(t.Right, env)
	if err != nil {
		return nil, err
	}
	return ResolvedUnary{Operator: t.Operator, Right: resolvedExp}, nil
}

func ResolveBinary(t parser.Binary, env *parser.Stmtsenv) (ResolvedExpr, error) {
	resolvedExpLeft, err := ResolveExpr(t.Left, env)
	if err != nil {
		return nil, err
	}
	resolvedExpRight, err := ResolveExpr(t.Right, env)
	if err != nil {
		return nil, err
	}
	return ResolvedBinary{Left: resolvedExpLeft, Right: resolvedExpRight, Operator: t.Operator}, nil
}
func ResolveCall(t parser.Call, env *parser.Stmtsenv) (ResolvedExpr, error) {
	resolvedCalleeExp, err := ResolveExpr(t.Callee, env)
	var resolvedArgsExp ResolvedExpr
	if err != nil {
		return nil, err
	}
	if t.Args != nil {
		resolvedArgsExp, err = ResolveExpr(t.Args, env)
		if err != nil {
			return nil, err
		}
	}
	return ResolvedCall{Callee: resolvedCalleeExp, Args: resolvedArgsExp, Arrity: t.Arrity, Operator: t.Operator}, nil
}

func ResolvePrimary(t parser.Primary, env *parser.Stmtsenv) (ResolvedExpr, error) {
	cur := env
	scopeDepth := 1
	if t.Node.Ttype == scanner.IDENTIFIER {
		for cur != nil {
			_, itExist := cur.Local[t.Node.Lexem]
			if itExist {
				//for a variable to be used in a scope , it needs to appear at the left hand side of an assignment/either sides of a binary expression/one side of a unary expression
				// // basically if it appears to be storage target/destination is not being used other than that is being used.
				if TopOfCallStack != DESTINATION_ASSIGNMENT {
					cur.Local[t.Node.Lexem] = VariableMetaData{isUsed: true, isResolved: false}
				}
				return ResolvedPrimary{Node: t.Node, ScopeDepth: scopeDepth}, nil
			}
			cur = cur.Encloser
			scopeDepth++
		}
	} else {
		//scope Depth does not apply to literals/variables declared
		scopeDepth = -1
	}
	if cur == nil {
		//not defined
		scopeDepth = 0
	}
	return ResolvedPrimary{Node: t.Node, ScopeDepth: scopeDepth}, nil
}

func ResolveExpr(t parser.Exp, env *parser.Stmtsenv) (ResolvedExpr, error) {
	switch t := t.(type) {
	case parser.List:
		return ResolveList(t, env)
	case parser.Assignment:
		return ResolveAssignment(t, env)
	case parser.LogicalOr:
		return ResolveLogicalOr(t, env)
	case parser.LogicalAnd:
		return ResolveLogicalAnd(t, env)
	case parser.Call:
		return ResolveCall(t, env)
	case parser.Unary:
		return ResolveUnary(t, env)
	case parser.Binary:
		return ResolveBinary(t, env)
	case parser.Primary:
		return ResolvePrimary(t, env)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", t)
	}
}

func ResolveStmt(t parser.Stmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	switch t := t.(type) {
	case parser.VarStmt:
		return ResolveVarStmt(t, env)
	case parser.BlockStmt:
		return ResolveBlockStmt(t, nil)
	case parser.WhileStmt:
		return ResolveWhileStmt(t, nil)
	case parser.FuncDef:
		return ResolveFuncDef(t, nil)
	case parser.ForStmt:
		return ResolveForStmt(t, nil)
	case parser.ExpStmt:
		return ResolveExpStmt(t, env)
	case parser.PrintStmt:
		return ResolvePrintStmt(t, env)
	case parser.ReturnStmt:
		return ResolveReturnStmt(t, env)
	default:
		return nil, fmt.Errorf("unknown statement type: %T", t)
	}
}

func Resolver(stmts []parser.Stmt, env *parser.Stmtsenv) ([]ResolvedStmt, error) {
	var resolvedStmts []ResolvedStmt
	for _, stmt := range stmts {
		resolvedStmt, err := ResolveStmt(stmt, env)
		if err != nil {
			return nil, err
		}
		resolvedStmts = append(resolvedStmts, resolvedStmt)
	}
	return resolvedStmts, nil
}

func (t ResolvedReturnStmt) Execute(env *parser.Stmtsenv) error {
	value, err := t.Exp.Evaluate(env)
	if err != nil {
		return err
	}
	panic(BangScriptReturn{
		value: value,
	})
}
func (t ResolvedFuncDef) Execute(env *parser.Stmtsenv) error {
	bs := t.Body
	bs.Env.Encloser.Local[t.Name.Lexem] = t
	return nil
}
func (t ResolvedForStmt) Execute(parent *parser.Stmtsenv) error {
	return t.Stmt.Execute(nil)
}
func (t ResolvedWhileStmt) Execute(env *parser.Stmtsenv) error {
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
func (t ResolvedIfStmt) Execute(env *parser.Stmtsenv) error {

	obj, err := t.Condition.Evaluate(env)
	if err != nil {
		return err
	}
	isTruth := isTruthy(obj)

	if isTruth {
		switch s := t.Thenbody.(type) {
		case ResolvedWhileStmt:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case ResolvedBlockStmt:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case ResolvedFuncDef:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case ResolvedForStmt:
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
		case ResolvedWhileStmt:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case ResolvedBlockStmt:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case ResolvedFuncDef:
			err = s.StaticToDynamic(env)
			if err != nil {
				return err
			}
			err = s.Execute(env)
		case ResolvedForStmt:
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
func (t ResolvedBlockStmt) Execute(env *parser.Stmtsenv) error {
	for _, stmt := range t.Stmts {
		if stmt == nil {
			continue
		}
		//functions have dynamic environment,created when they are called as such nested environments should be updated to enclose around this new env before they are executed
		var err error
		switch s := stmt.(type) {
		case ResolvedWhileStmt:
			err = s.StaticToDynamic(t.Env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case ResolvedBlockStmt:
			err = s.StaticToDynamic(t.Env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case ResolvedFuncDef:
			err = s.StaticToDynamic(t.Env)
			if err != nil {
				return err
			}
			err = s.Execute(nil)
		case ResolvedForStmt:
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
func (t ResolvedPrintStmt) Execute(env *parser.Stmtsenv) error {
	Obj, err := t.Exp.Evaluate(env)
	if err != nil {
		return err
	}
	//TODO:function def don't evaluate to simple scalar values like integer,bool etc it evaluates to a struct which won't print nicely
	fmt.Printf("%v\n", Obj)
	return nil
}
func (t ResolvedVarStmt) Execute(env *parser.Stmtsenv) error {
	//check if variable declaration have a definition
	if t.Exp != nil {
		switch s := t.Exp.(type) {
		case ResolvedPrimary:
			env.Local[s.Node.Lexem] = nil
		case ResolvedAssignment:
			//for now only primary identifier expressions map to a storage location
			lv := s.StoreTarget.(ResolvedPrimary)
			obj, err := s.Right.Evaluate(env)
			if err != nil {
				return err
			}
			env.Local[lv.Node.Lexem] = obj
		case ResolvedList:
			for _, exp := range s.Expressions {
				err := ResolvedVarStmt{exp}.Execute(env)
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
func (t ResolvedExpStmt) Execute(env *parser.Stmtsenv) error {
	obj, err := t.Exp.Evaluate(env)
	if err != nil {
		return err
	}
	if parser.Mode == REPL && env.Encloser == nil {
		fmt.Printf("%v\n", obj)
		return nil
	}
	return err
}

// Resolved Expressions Implementations
func (u ResolvedUnary) Evaluate(env *parser.Stmtsenv) (parser.Obj, error) {
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
func (b ResolvedBinary) Evaluate(env *parser.Stmtsenv) (parser.Obj, error) {
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
func (t ResolvedLogicalAnd) Evaluate(env *parser.Stmtsenv) (parser.Obj, error) {
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
func (t ResolvedLogicalOr) Evaluate(env *parser.Stmtsenv) (parser.Obj, error) {
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
func (a ResolvedAssignment) Evaluate(env *parser.Stmtsenv) (parser.Obj, error) {
	cur := env
	lv := a.StoreTarget.(ResolvedPrimary)
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
func (l ResolvedList) Evaluate(env *parser.Stmtsenv) (parser.Obj, error) {
	var rvalue parser.Obj
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

func (t ResolvedFuncDef) call(env *parser.Stmtsenv, callInfo *ResolvedCall) (value parser.Obj, err error) {
	bs := t.Body

	defer func() {
		if r := recover(); r != nil {
			switch s := r.(type) {
			case BangScriptReturn:
				value = s.value
			default:
				panic(r)
			}
		}
	}()

	//bs.Env need to be a complete copy, since it's new, any environment that once enclosed around bs.Env will be updated to support closures
	newEnv := parser.Stmtsenv{Local: map[string]parser.Obj{}, Encloser: bs.Env.Encloser, Policy: parser.DYNAMIC}
	bs.Env = &newEnv

	envWithFunctionArgsOnly := newEnv.Local
	if callInfo.Args != nil {
		listArgs, isArgs := callInfo.Args.(ResolvedList)
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

	//the reason why we pass nil is that , we have already created the environments during parsing and each block has it
	// what we make fresh copies of environment during function calls, and update subsequent environments in the function as they too are dynamic
	//dynamic env
	err = bs.Execute(nil)
	if err != nil {
		return nil, err
	}
	return value, nil
}
func (t ResolvedCall) Evaluate(env *parser.Stmtsenv) (parser.Obj, error) {
	value, err := t.Callee.Evaluate(env)
	if err != nil {
		return nil, err
	}
	function, isCallable := value.(ResolvedFuncDef)
	if !isCallable {
		return nil, fmt.Errorf("Can't call expression at line %d", t.Operator.Line)
	}
	if t.Arrity != function.Arrity {
		return nil, fmt.Errorf("Expected %d arguments but got %d instead", function.Arrity, t.Arrity)
	}
	//env is a parent environment can be immediate env or global
	return function.call(env, &t)
}
func (t ResolvedPrimary) Evaluate(env *parser.Stmtsenv) (parser.Obj, error) {
	if t.ScopeDepth == 0 {
		return nil, fmt.Errorf("Variable %s is not defined", t.Node.Lexem)
	}
	// variable operand
	if t.ScopeDepth > 0 {
		env = env.Get(t.ScopeDepth)
		return env.Local[t.Node.Lexem], nil
	}
	//if scope depth is -1 then it is a constant literal
	// operands needs to be only following string , number and boolean
	switch t.Node.Ttype {
	case scanner.NUMBER:
		op, err := strconv.ParseFloat(t.Node.Lexem, 64)
		if err != nil {
			//handle error for failed type conversion
			return nil, err
		}
		return op, nil
	//string concantenation and comparison
	case scanner.STRING:
		return t.Node.Lexem, nil
	//boolean algebra
	case scanner.TRUE:
		return true, nil
	case scanner.FALSE:
		return false, nil
	case scanner.NIL:
		return nil, nil
	default:
		return nil, fmt.Errorf("Expected a string, number, a target location , nil and boolean but got %d at line %d", t.Node.Ttype, t.Node.Line)
	}
}
