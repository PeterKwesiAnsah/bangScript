/*
TODO: run interpretor via path to source code file
*/
/*
TODO: run interpretor through interactive prompt (REPL)
*/
package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Print("Usage: glox [path_to_script]\n")
		os.Exit(1)
	} else if len(args) == 2 {
		fmt.Print("Run: Script Source Code\n")
	} else {
		fmt.Print("Run: REPL\n")
	}
}
