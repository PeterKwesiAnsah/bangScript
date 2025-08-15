package resolver

import (
	"bangScript/gbs/parser"
	"bangScript/gbs/scanner"
)

type ResolvedStmt interface {
	execute() error
}
type ResolvedExpr interface {
	evaluate() error
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
	condition ResolvedExpr
	thenbody  ResolvedStmt
	elsebody  ResolvedStmt
}
type ResolvedForStmt struct {
	stmt ResolvedWhileStmt
}
type ResolvedWhileStmt struct {
	condition ResolvedExpr
	body      ResolvedBlockStmt
	init      ResolvedStmt
	env       *parser.Stmtsenv
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
	exp ResolvedExpr
}

type ResolvedFuncDef struct {
	name   *scanner.Token
	params []*scanner.Token
	body   ResolvedBlockStmt
	arrity int
}

type ResolvedReturnStmt struct {
	//currently we allow returns to exps for now. In the future we change to statment because of closures
	Exp ResolvedExpr
}

func (t ResolvedVarStmt) execute() error {
	return nil
}
func (t ResolvedPrintStmt) execute() error {
	return nil
}
func (t ResolvedBlockStmt) execute() error {
	return nil
}
func (t ResolvedList) evaluate() error {
	return nil
}
func (t ResolvedAssignment) evaluate() error {
	return nil
}
func (t ResolvedLogicalOr) evaluate() error {
	return nil
}
func (t ResolvedLogicalAnd) evaluate() error {
	return nil
}
func (t ResolvedUnary) evaluate() error {
	return nil
}
func (t ResolvedPrimary) evaluate() error {
	return nil
}
func (t ResolvedCall) evaluate() error {
	return nil
}
func (t ResolvedBinary) evaluate() error {
	return nil
}
