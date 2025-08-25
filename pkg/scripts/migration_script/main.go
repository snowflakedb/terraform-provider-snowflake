package main

import (
	"log"
)

func main() {
	config, err := ParseInputArguments()
	if err != nil {
		log.Fatalf("Error parsing input arguments: %v, run -h to get more information on running the script.", err)
	}

	csvInput, err := ReadCsvFromStdin()
	if err != nil {
		log.Fatalf("Error reading CSV input: %v", err)
	}

	switch config.ObjectType {
	case ObjectTypeGrants:
		if err := HandleGrants(csvInput); err != nil {
			log.Fatalf("Error handling grants generation: %v", err)
		}
	default:
		log.Fatalf("Unsupported object type: %s. Run -h to get more information on allowed object types.", config.ObjectType)
	}
}
