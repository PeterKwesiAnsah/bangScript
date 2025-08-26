package resolver

import (
	"bangScript/gbs/parser"
	"bangScript/gbs/scanner"
)

type ResolvedStmt interface {
	Execute(env *parser.Stmtsenv) error
}
type ResolvedExpr interface {
	Evaluate(env *parser.Stmtsenv) (parser.Obj, error)
}

type ResolvedCall struct {
	Callee   ResolvedExpr
	Operator *scanner.Token
	Arrity   int
	Args     ResolvedExpr
}

type ResolvedList struct {
	Expressions []ResolvedExpr
}

type ResolvedAssignment struct {
	//l-value
	StoreTarget ResolvedExpr
	Operator    *scanner.Token
	//r-value
	Right ResolvedExpr
}

type ResolvedLogicalOr struct {
	Left     ResolvedExpr
	Operator *scanner.Token
	Right    ResolvedExpr
}
type ResolvedLogicalAnd struct {
	Left     ResolvedExpr
	Operator *scanner.Token
	Right    ResolvedExpr
}

type ResolvedBinary struct {
	Left     ResolvedExpr
	Operator *scanner.Token
	Right    ResolvedExpr
}

// TODO: for handling conditional operations
type ResolvedTernary struct {
	condition ResolvedExpr
	operator  *scanner.Token
	then      ResolvedExpr
	elsef     ResolvedExpr
}

type ResolvedUnary struct {
	Operator *scanner.Token
	Right    ResolvedExpr
}

type ResolvedPrimary struct {
	Node *scanner.Token
	//-1,0 and positive integer
	// -1 means scope depth does not apply to the node (it can be a constant literal or a variable)
	// 0 means the variable being referenced is undefined
	// >0 means the variable being referenced is defined in the scope depth
	ScopeDepth int
}

type ResolvedIfStmt struct {
	Condition ResolvedExpr
	Thenbody  ResolvedStmt
	Elsebody  ResolvedStmt
}
type ResolvedForStmt struct {
	Stmt ResolvedWhileStmt
}
type ResolvedWhileStmt struct {
	Condition ResolvedExpr
	Body      ResolvedBlockStmt
	Init      ResolvedStmt
	Env       *parser.Stmtsenv
}

type ResolvedBlockStmt struct {
	Stmts []ResolvedStmt
	Env   *parser.Stmtsenv
}

type ResolvedVarStmt struct {
	Exp ResolvedExpr
}
type ResolvedPrintStmt struct {
	Exp ResolvedExpr
}

type ResolvedExpStmt struct {
	Exp ResolvedExpr
}

type ResolvedFuncDef struct {
	Name   *scanner.Token
	Params []*scanner.Token
	Body   ResolvedBlockStmt
	Arrity int
}

type ResolvedReturnStmt struct {
	//currently we allow returns to exps for now. In the future we change to statment because of closures
	Exp ResolvedExpr
}
