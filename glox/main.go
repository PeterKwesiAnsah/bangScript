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
		t, err := parser.Parser(tokens)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		res, err := t.Evaluate()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("%v\n", res)
	} else {
		fmt.Print("Run: REPL\n")
	}
}
