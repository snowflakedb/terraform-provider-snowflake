package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

type ObjectType string

const (
	ObjectTypeGrants ObjectType = "grants"
)

var AllObjectTypes = []ObjectType{
	ObjectTypeGrants,
}

func ToObjectType(s string) (ObjectType, error) {
	if slices.Contains(AllObjectTypes, ObjectType(strings.ToLower(s))) {
		return ObjectType(strings.ToLower(s)), nil
	}
	return "", fmt.Errorf("unsupported object type: %s", s)
}

type ImportStatementType string

const (
	ImportStatementTypeStatement ImportStatementType = "statement"
	ImportStatementTypeBlock     ImportStatementType = "block"
)

var AllImportStatementTypes = []ImportStatementType{
	ImportStatementTypeStatement,
	ImportStatementTypeBlock,
}

func ToImportStatementType(s string) (ImportStatementType, error) {
	if slices.Contains(AllImportStatementTypes, ImportStatementType(strings.ToLower(s))) {
		return ImportStatementType(strings.ToLower(s)), nil
	}

	return "", fmt.Errorf("invalid import statement type: %s", s)
}

type Config struct {
	ObjectType ObjectType
	ImportFlag ImportStatementType
}

func ParseInputArguments() (*Config, error) {
	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: migration_script [flags] [object_type]")
		_, _ = fmt.Fprintln(os.Stderr, "")

		_, _ = fmt.Fprintln(os.Stderr, "Object types:")
		for _, ot := range AllObjectTypes {
			_, _ = fmt.Fprintf(os.Stderr, "	- %s\n", ot)
		}
		_, _ = fmt.Fprintln(os.Stderr, "")

		_, _ = fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}

	// flags
	importFlagString := flag.String("import", "statement", collections.JoinStrings([]string{
		"Determines the output format for import statements.",
		"Possible values:",
		"	- \"statement\" will print appropriate terraform import statement at the end of generated content",
		"	- \"block\" will generate import block next to every generated resource",
		"", // required for default value formatting
	}, "\n"))

	flag.Parse()

	importFlagType, err := ToImportStatementType(*importFlagString)
	if err != nil {
		return nil, fmt.Errorf("error parsing import flag: %w", err)
	}

	// positional arguments
	positionalArguments := flag.Args()

	if len(positionalArguments) != 1 {
		return nil, fmt.Errorf("no object type specified, use -h for help")
	}

	parsedObjectType, err := ToObjectType(positionalArguments[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing object type: %w", err)
	}

	return &Config{
		ObjectType: parsedObjectType,
		ImportFlag: importFlagType,
	}, nil
}

func ReadCsvFromStdin() ([][]string, error) {
	inputBytes, err := io.ReadAll(bufio.NewReader(os.Stdin))
	if err != nil {
		return nil, fmt.Errorf("error reading STDIN input: %w", err)
	}

	csvReader := csv.NewReader(bytes.NewBuffer(inputBytes))
	return csvReader.ReadAll()
}
