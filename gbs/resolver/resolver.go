package resolver

import (
	"bangScript/gbs/parser"
	"bangScript/gbs/scanner"
	"fmt"
	"strconv"
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

// move parser.(type).evaluate,parser.(type).execute to resolver.(type).execute and resolver.(type).evaluate
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
		if lv.Node.Ttype != scanner.IDENTIFIER {
			return nil, fmt.Errorf("Expected an identifier but got something else %d", lv.Node.Ttype)
		}
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

// TODO: Resolved Expressions Implementations
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
