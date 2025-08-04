/*
TODO: run interpretor through interactive prompt (REPL)
*/
package main

import (
	"bangScript/gbs/parser"
	"bangScript/gbs/scanner"
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) > 2 {
		fmt.Printf("Usage: bs [path_to_script] or bs (to launch REPL)\n")
		os.Exit(1)
	} else if len(args) == 2 {
		path := args[1]
		source, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Failed to open file :%s\n", err.Error())
			os.Exit(1)
		}
		tokens, err := scanner.ScanTokens(source)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		globalEnv := parser.Stmtsenv{Local: map[string]parser.Obj{}, Encloser: nil}
		stmts, err := parser.Parser(tokens, &globalEnv)
		if err != nil {
			fmt.Printf("ParseError: %s\n", err.Error())
			os.Exit(1)
		}
		for _, stmt := range stmts {
			if stmt == nil {
				continue
			}
			var executionError error
			switch stmt.(type) {
			case parser.WhileStmt:
				executionError = stmt.Execute(nil)
			case parser.BlockStmt:
				executionError = stmt.Execute(nil)
				//case.parse.funcDef:
			//executionError = stmt.Execute(nil)
			default:
				//for statements that have their own env,statement.Env you will be executed with nil
				executionError = stmt.Execute(&globalEnv)
			}
			if executionError != nil {
				fmt.Printf("ExecutionError: %s\n", executionError.Error())
				os.Exit(1)
			}
		}
	} else {
		fmt.Print("Run: REPL\n")
	}
}
