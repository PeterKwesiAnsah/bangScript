package resolver

import (
	"bangScript/gbs/parser"
	"bangScript/gbs/scanner"
	"fmt"
)

func ResolveVarStmt(t parser.VarStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	switch expT := t.Exp.(type) {
	case parser.Primary:
		lexeme := expT.Node.Lexem
		_, itExists := env.Local[lexeme]
		if itExists {
			return nil, fmt.Errorf("Cannot redeclare a variable in the same scope.")
		}
		env.Local[lexeme] = true
		//-1 because it binds a variable to the current scope
		return ResolvedVarStmt{Exp: ResolvedPrimary{Node: expT.Node, ScopeDepth: -1}}, nil
	case parser.Assignment:
		lv, isStorageTarget := expT.StoreTarget.(parser.Primary)
		if !(isStorageTarget && lv.Node.Ttype == scanner.IDENTIFIER) {
			return nil, fmt.Errorf("Cannot use the l-value as storage target")
		}
		rv, isPrimaryExpression := expT.Right.(parser.Primary)
		if isPrimaryExpression && rv.Node.Lexem == lv.Node.Lexem {
			return nil, fmt.Errorf("Cannot initizialize a variable using the same variable being declared")
		}
		resolvedExpr, err := ResolveExpr(expT.Right, env)
		if err != nil {
			return nil, err
		}
		env.Local[lv.Node.Lexem] = true
		return ResolvedVarStmt{Exp: resolvedExpr}, nil
	case parser.List:
		//resolvedExpr, err := ResolveList(expT, env)
		//if err != nil {
		//	return nil, err
		//}
	default:
		return nil, fmt.Errorf("Invalid expression type")
	}
	return nil, nil
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
	return nil, nil
}

func ResolveReturnStmt(t parser.ReturnStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	return nil, nil
}

func ResolveIfStmt(t parser.IfStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	return nil, nil
}

func ResolveWhileStmt(t parser.WhileStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	return nil, nil
}

func ResolveForStmt(t parser.ForStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	return nil, nil
}

func ResolveExpStmt(t parser.ExpStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	return nil, nil
}

func ResolvePrintStmt(t parser.PrintStmt, env *parser.Stmtsenv) (ResolvedStmt, error) {
	resolvedExpr, err := ResolveExpr(t.Exp, env)
	if err != nil {
		return nil, err
	}
	return ResolvedPrintStmt{Exp: resolvedExpr}, nil
}

func ResolveAssignment(t parser.Assignment, env *parser.Stmtsenv) (ResolvedExpr, error) {
	return nil, nil
}

func ResolveList(t parser.List, env *parser.Stmtsenv) (ResolvedExpr, error) {
	return nil, nil
}

func ResolveLogicalOr(t parser.LogicalOr, env *parser.Stmtsenv) (ResolvedExpr, error) {
	return nil, nil
}

func ResolveLogicalAnd(t parser.LogicalAnd, env *parser.Stmtsenv) (ResolvedExpr, error) {
	return nil, nil
}

func ResolveCall(t parser.Call, env *parser.Stmtsenv) (ResolvedExpr, error) {
	return nil, nil
}

func ResolveUnary(t parser.Unary, env *parser.Stmtsenv) (ResolvedExpr, error) {
	return nil, nil
}

func ResolveBinary(t parser.Binary, env *parser.Stmtsenv) (ResolvedExpr, error) {
	return nil, nil
}

func ResolvePrimary(t parser.Primary, env *parser.Stmtsenv) (ResolvedExpr, error) {
	cur := env
	scopeDepth := 0
	for cur != nil {
		_, itExist := cur.Local[t.Node.Lexem]
		if itExist {
			break
		}
		cur = cur.Encloser
		scopeDepth++
	}
	if cur == nil {
		scopeDepth = -1
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
	case parser.ForStmt:
		return ResolveForStmt(t, env)
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
