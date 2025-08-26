package resolver

import (
	"bangScript/gbs/parser"
	"fmt"
)

func (t *ResolvedForStmt) StaticToDynamic(parent *parser.Stmtsenv) error {
	//t.stmt.Env is environment for condition,initializer and single body statement
	if t.Stmt.Env == parent {
		return fmt.Errorf("Child environment can not be the same as Parent")
	}
	if parent.Policy == DYNAMIC {
		//condition,initializer and single body statement
		newEnv := &parser.Stmtsenv{Local: map[string]parser.Obj{}, Encloser: parent, Policy: DYNAMIC}
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
func (t *ResolvedFuncDef) StaticToDynamic(parent *parser.Stmtsenv) error {
	if t.Body.Env == parent {
		return fmt.Errorf("Child environment can not be the same as Parent")
	}
	if parent.Policy == DYNAMIC {
		t.Body.Env = &parser.Stmtsenv{Local: map[string]parser.Obj{}, Encloser: parent, Policy: DYNAMIC}
	} else {
		if t.Body.Env.Encloser != parent {
			return fmt.Errorf("ExecutionError: Body statement environment should encloses around the env passed to it ")
		}
	}
	return nil
}
func (t *ResolvedBlockStmt) StaticToDynamic(parent *parser.Stmtsenv) error {
	if t.Env == parent {
		return fmt.Errorf("Child environment can not be the same as Parent")
	}
	if parent.Policy == DYNAMIC {
		t.Env = &parser.Stmtsenv{Local: map[string]parser.Obj{}, Encloser: parent, Policy: DYNAMIC}
	} else {
		if t.Env.Encloser != parent {
			return fmt.Errorf("Body statement environment should encloses around the env passed to it ")
		}
	}
	return nil
}
func (t *ResolvedWhileStmt) StaticToDynamic(parent *parser.Stmtsenv) error {
	if t.Env == parent {
		return fmt.Errorf("Child environment can not be the same as Parent")
	}
	if parent.Policy == DYNAMIC {
		//condition,initializer and single body statement
		newEnv := &parser.Stmtsenv{Local: map[string]parser.Obj{}, Encloser: parent, Policy: DYNAMIC}
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

func isTruthy(val parser.Obj) bool {
	isTruthy := true
	var falsyVal []parser.Obj = []parser.Obj{"", nil, 0, false}
	for _, falsy := range falsyVal {
		if falsy == val {
			isTruthy = false
			break
		}
	}
	return isTruthy
}
