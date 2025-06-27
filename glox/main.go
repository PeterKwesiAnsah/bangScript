/*
TODO: run interpretor through interactive prompt (REPL)
*/
package main

import (
	"fmt"
	"lox/glox/parser"
	"lox/glox/scanner"
	"os"
)

func main() {
	args := os.Args

	if len(args) > 2 {
		fmt.Printf("Usage: glox [path_to_script]\n")
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
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		for _, stmt := range stmts {
			err := stmt.Execute(&globalEnv)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				os.Exit(1)
			}
		}
	} else {
		fmt.Print("Run: REPL\n")
	}
}
