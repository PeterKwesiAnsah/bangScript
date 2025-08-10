package resolver

import "bangScript/gbs/parser"

type Resolved interface {
	Varesolution(env *parser.Stmtsenv) error
}

type ResVarStmt struct {
	tree *parser.VarStmt
}

func (t ResVarStmt) Varesolution(env *parser.Stmtsenv) error

func Resolver(stmts []parser.Stmt) []Resolved
