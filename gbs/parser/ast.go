package parser

import "bangScript/gbs/scanner"

type Call struct {
	Callee   Exp
	Operator *scanner.Token
	Arrity   int
	Args     Exp
}

type List struct {
	Expressions []Exp
}

type Assignment struct {
	//l-value
	StoreTarget Exp
	Operator    *scanner.Token
	//r-value
	Right Exp
}

type LogicalOr struct {
	Left     Exp
	Operator *scanner.Token
	Right    Exp
}
type LogicalAnd struct {
	Left     Exp
	Operator *scanner.Token
	Right    Exp
}

type Binary struct {
	Left     Exp
	Operator *scanner.Token
	Right    Exp
}

// TODO: for handling conditional operations
type tenary struct {
	condition Exp
	operator  *scanner.Token
	then      Exp
	elsef     Exp
}

type Unary struct {
	Operator *scanner.Token
	Right    Exp
}

type Primary struct {
	Node *scanner.Token
}

type IfStmt struct {
	Condition Exp
	Thenbody  Stmt
	Elsebody  Stmt
}
type ForStmt struct {
	Stmt WhileStmt
}
type WhileStmt struct {
	Condition Exp
	Body      BlockStmt
	Init      Stmt
	Env       *Stmtsenv
}

type BlockStmt struct {
	Stmts []Stmt
	Env   *Stmtsenv
}

type VarStmt struct {
	Exp Exp
}
type PrintStmt struct {
	Exp Exp
}

type ExpStmt struct {
	Exp Exp
}

type FuncDef struct {
	Name   *scanner.Token
	Params []*scanner.Token
	Body   BlockStmt
	Arrity int
}

type ContinueStmt struct {
	Token *scanner.Token
}
type BreakStmt struct {
	Token *scanner.Token
}

type ReturnStmt struct {
	//currently we allow returns to exps for now. In the future we change to statment because of closures
	Exp Exp
}

// TODO: implementation of Stmt interface

func (t VarStmt) print() string {
	return ""
}
func (t PrintStmt) print() string {
	return ""
}
func (t ExpStmt) print() string {
	return ""
}
func (t FuncDef) print() string {
	return ""
}
func (t ReturnStmt) print() string {
	return ""
}
func (t IfStmt) print() string {
	return ""
}
func (t ForStmt) print() string {
	return ""
}
func (t WhileStmt) print() string {
	return ""
}
func (t BlockStmt) print() string {
	return ""
}

func (t ContinueStmt) print() string {
	return ""
}
func (t BreakStmt) print() string {
	return ""
}

// TODO: implementation of Exp interface
func (t Call) print() string {
	return ""
}
func (t List) print() string {
	return ""
}
func (t Assignment) print() string {
	return ""
}
func (t LogicalOr) print() string {
	return ""
}
func (t LogicalAnd) print() string {
	return ""
}
func (t Unary) print() string {
	return ""
}
func (t Primary) print() string {
	return ""
}
func (t Binary) print() string {
	return ""
}
