package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"io"
	"os"
	"slices"
	"strings"
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

type Program struct {
	Args           []string
	StdOut, StdErr io.Writer
	StdIn          io.Reader
	Config         *Config
}

func NewProgram() *Program {
	return &Program{
		Args:   os.Args,
		StdOut: os.Stdout,
		StdErr: os.Stderr,
		StdIn:  os.Stdin,
	}
}

func (p *Program) Run() int {
	config, err := p.parseInputArguments()
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}

		_, _ = fmt.Fprintf(p.StdErr, "Error parsing input arguments: %v, run -h to get more information on running the script", err)
		return 1
	}
	p.Config = config

	input, err := readAllAsCsv(p.StdIn)
	if err != nil {
		_, _ = fmt.Fprintf(p.StdErr, "Error reading CSV input: %v", err)
		return 2
	}

	output, err := p.generateOutput(input)
	if err != nil {
		_, _ = fmt.Fprintf(p.StdErr, "Error generating output: %v", err)
		return 3
	}
	_, _ = fmt.Fprint(p.StdOut, output)

	return 0
}

func (p *Program) parseInputArguments() (*Config, error) {
	commandLine := flag.NewFlagSet(p.Args[0], flag.ContinueOnError)
	commandLine.SetOutput(p.StdErr)
	commandLine.Usage = func() {
		_, _ = fmt.Fprintln(p.StdErr, "Usage: migration_script [flags] [object_type]")
		_, _ = fmt.Fprintln(p.StdErr, "")

		_, _ = fmt.Fprintln(p.StdErr, "Object types:")
		for _, ot := range AllObjectTypes {
			_, _ = fmt.Fprintf(p.StdErr, "	- %s\n", ot)
		}
		_, _ = fmt.Fprintln(p.StdErr, "")

		_, _ = fmt.Fprintln(p.StdErr, "Flags:")
		commandLine.PrintDefaults()
	}

	// flags
	importFlagString := commandLine.String("import", "statement", collections.JoinStrings([]string{
		"Determines the output format for import statements.",
		"Possible values:",
		"	- \"statement\" will print appropriate terraform import statement at the end of generated content",
		"	- \"block\" will generate import block next to every generated resource",
		"", // required for default value formatting
	}, "\n"))

	if err := commandLine.Parse(p.Args[1:]); err != nil {
		return nil, err
	}

	// positional arguments
	args := commandLine.Args()
	if len(args) != 1 {
		return nil, fmt.Errorf("no object type specified, use -h for help")
	}
	parsedObjectType, err := ToObjectType(args[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing object type: %w", err)
	}

	importFlagType, err := ToImportStatementType(*importFlagString)
	if err != nil {
		return nil, fmt.Errorf("error parsing import flag: %w", err)
	}

	return &Config{
		ObjectType: parsedObjectType,
		ImportFlag: importFlagType,
	}, nil
}

func readAllAsCsv(reader io.Reader) ([][]string, error) {
	inputBytes, err := io.ReadAll(bufio.NewReader(reader))
	if err != nil {
		return nil, fmt.Errorf("error reading CSV input: %w", err)
	}

	csvReader := csv.NewReader(bytes.NewBuffer(inputBytes))
	return csvReader.ReadAll()
}

func (p *Program) generateOutput(input [][]string) (string, error) {
	switch p.Config.ObjectType {
	case ObjectTypeGrants:
		return HandleGrants(input)
	default:
		return "", fmt.Errorf("unsupported object type: %s, run -h to get more information on allowed object types", p.Config.ObjectType)
	}
}
