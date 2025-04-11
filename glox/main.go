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
		fmt.Printf("Usage: glox [path_to_script]\n")
		os.Exit(1)
	} else if len(args) == 2 {
		path := args[1]
		fp, err := os.Open(path)
		if err != nil {
			fmt.Errorf("Failed to open file :%w\n", err)
			os.Exit(1)
		}
		fileStat, err := fp.Stat()
		if err != nil {
			fmt.Errorf("%w\n", err)
			os.Exit(1)
		}
		sizeofFile := fileStat.Size()
		buf := make([]byte, sizeofFile)

	} else {
		fmt.Print("Run: REPL\n")
	}
}
