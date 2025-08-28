/*
TODO: run interpretor through interactive prompt (REPL)
*/
package main

import (
	"bangScript/gbs/parser"
	"bangScript/gbs/resolver"
	"bangScript/gbs/scanner"
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	REPL uint8 = iota
	SCRIPT
)

type source []byte

var globalEnv = parser.Stmtsenv{Local: map[string]parser.Obj{}, Encloser: nil}

func (t source) RunCode(mode uint8) error {
	tokens, err := scanner.ScanTokens(t)
	if err != nil {
		return err
	}
	stmts, err := parser.Parser(tokens, &globalEnv, mode)
	if err != nil {
		return err
	}
	resolvedStmts, err := resolver.Resolver(stmts, &globalEnv)

	if err != nil {
		return err
	}
	for _, rstmt := range resolvedStmts {
		if rstmt == nil {
			continue
		}
		var executionError error
		switch rstmt.(type) {
		case resolver.ResolvedWhileStmt:
			executionError = rstmt.Execute(nil)
		case resolver.ResolvedBlockStmt:
			executionError = rstmt.Execute(nil)
		case resolver.ResolvedForStmt:
			executionError = rstmt.Execute(nil)
		case resolver.ResolvedFuncDef:
			executionError = rstmt.Execute(nil)
		default:
			executionError = rstmt.Execute(&globalEnv)
		}
		if executionError != nil {
			return executionError
		}
	}
	return nil
}

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Printf("Usage: bs [path_to_script] or bs (to launch REPL)\n")
		os.Exit(1)
	} else if len(args) == 2 {
		path := args[1]
		contents, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Failed to open file :%s\n", err.Error())
			os.Exit(1)
		}
		var src source = contents
		err = src.RunCode(SCRIPT)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}
	} else {
		scannerIO := bufio.NewScanner(os.Stdin)
		fmt.Println("Welcome to bangScript interactive REPL")
		fmt.Println("Press 'exit' or 'quit' to leave the REPL")
		for {
			fmt.Print("> ")
			if !scannerIO.Scan() {
				break
			}
			input := strings.TrimSpace(scannerIO.Text())
			if input == "exit" || input == "quit" {
				break
			}
			var src source = []byte(input)
			err := src.RunCode(REPL)
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}
			if err := scannerIO.Err(); err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				os.Exit(1)
			}
		}
	}
}
