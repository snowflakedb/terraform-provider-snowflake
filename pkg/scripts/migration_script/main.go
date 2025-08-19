package main

import (
	"fmt"
)

func main() {
	ParseInputArguments()

	// Read input from stdin
	csvInput := ReadCsvFromStdin()

	// Parse the input and transform to objects
	switch objectType {
	case ObjectTypeGrants:
		HandleGrants(csvInput)
	default:
		panic(fmt.Sprintf("unsupported object type: %s", objectType))
	}
}
