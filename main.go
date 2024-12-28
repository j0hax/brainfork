package main

import (
	"io"
	"os"

	"github.com/j0hax/brainfork/interpreter"
)

func main() {
	if len(os.Args) < 2 {
		e, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		interpreter := interpreter.NewInterpreter(os.Stdin, os.Stdout)
		interpreter.Run(e)
	} else {
		files := os.Args[1:]
		run(files...)
	}
}

func run(files ...string) {
	for _, f := range files {
		e, err := os.ReadFile(f)
		if err != nil {
			panic(err)
		}
		interpreter := interpreter.NewInterpreter(os.Stdin, os.Stdout)
		interpreter.Run(e)
	}
}
