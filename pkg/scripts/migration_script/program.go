package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
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

type ExitCode int

const (
	ExitCodeSuccess ExitCode = iota
	ExitCodeFailedInputArgumentParsing
	ExitCodeFailedCsvInputParsing
	ExitCodeFailedGeneratingTerraformOutput
)

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

func NewDefaultProgram() *Program {
	return &Program{
		Args:   os.Args,
		StdOut: os.Stdout,
		StdErr: os.Stderr,
		StdIn:  os.Stdin,
	}
}

func (p *Program) Run() ExitCode {
	config, err := p.parseInputArguments()
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}

		_, _ = fmt.Fprintf(p.StdErr, "Error parsing input arguments: %v, run -h to get more information on running the script", err)
		return ExitCodeFailedInputArgumentParsing
	}
	p.Config = config

	input, err := readAllAsCsv(p.StdIn)
	if err != nil {
		_, _ = fmt.Fprintf(p.StdErr, "Error reading CSV input: %v", err)
		return ExitCodeFailedCsvInputParsing
	}

	output, err := p.generateOutput(input)
	if err != nil {
		_, _ = fmt.Fprintf(p.StdErr, "Error generating output: %v", err)
		return ExitCodeFailedGeneratingTerraformOutput
	}
	_, _ = fmt.Fprint(p.StdOut, output)

	return ExitCodeSuccess
}

func (p *Program) parseInputArguments() (*Config, error) {
	commandLine := flag.NewFlagSet(p.Args[0], flag.ContinueOnError)
	commandLine.SetOutput(p.StdErr)
	commandLine.Usage = func() {
		_, _ = fmt.Fprint(p.StdErr, `Migration script's purpose is to generate terraform resources from existing Snowflake objects.
It operates on STDIN input and expects output from Snowflake commands in the CSV format.
The script writes generated terraform resources to STDOUT in case you want to redirect it to a file.
Any logs or errors are written to STDERR. You should separate outputs from STDOUT and STDERR when running the script (e.g. by redirecting STDOUT to a file)
to clearly see in case of any errors or skipped objects (due to, for example, incorrect or unexpected format).

usage: migration_script [-import=<statement|block>] <object_type>

import optional flag determines the output format for import statements. The possible values are:
	- "statement" will print appropriate terraform import command at the end of generated content (default) (see https://developer.hashicorp.com/terraform/cli/commands/import)
	- "block" will generate import block at the end of generated content (see https://developer.hashicorp.com/terraform/language/import)
	
object_type represents the type of Snowflake object you want to generate terraform resources for.
	It is a required positional argument and possible values are listed below.
	A given object_type corresponds to a specific Snowflake output expected as input to the script.
	Currently supported object types are:
		- "grants" which expects output from SHOW GRANTS command (see https://docs.snowflake.com/en/sql-reference/sql/show-grants) to generate new grant resources (see https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/grants_redesign_design_decisions#mapping-from-old-grant-resources-to-the-new-ones).
			The allowed SHOW GRANTS commands are:
				- 'SHOW GRANTS ON ACCOUNT'
				- 'SHOW GRANTS ON <object_type>'
				- 'SHOW GRANTS TO ROLE <role_name>'
				- 'SHOW GRANTS TO DATABASE ROLE <database_role_name>'
			Supported resources:
				- snowflake_grant_privileges_to_account_role
				- snowflake_grant_privileges_to_database_role
				- snowflake_grant_account_role
				- snowflake_grant_database_role
			Limitations:
				- grants on 'future' or on 'all' objects are not supported
				- all_privileges and always_apply fields are not supported
		
example usage:
	migration_script -import=block grants < show_grants_output.csv > generated_output.tf
`)
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
		return HandleGrants(p.Config, input)
	default:
		return "", fmt.Errorf("unsupported object type: %s, run -h to get more information on allowed object types", p.Config.ObjectType)
	}
}
