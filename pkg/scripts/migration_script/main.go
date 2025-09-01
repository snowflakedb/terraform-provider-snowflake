package main

import (
	"os"
)

func main() {
	program := NewDefaultProgram()
	exitCode := program.Run()
	os.Exit(int(exitCode))
}
