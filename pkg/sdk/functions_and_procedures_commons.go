package sdk

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

type NormalizedPath struct {
	// StageLocation is a normalized (fully-quoted id or `~`) stage location
	StageLocation string
	// PathOnStage is path to the file on stage without opening `/`
	PathOnStage string
}

// NormalizedArgument does not contain default value because it is not returned in the Signature (or any other field).
type NormalizedArgument struct {
	Name     string
	DataType datatypes.DataType
}

// TODO [SNOW-1850370]: use ParseCommaSeparatedStringArray + collections.MapErr combo here and in other methods?
func parseFunctionOrProcedureImports(importsRaw *string) ([]NormalizedPath, error) {
	normalizedImports := make([]NormalizedPath, 0)
	if importsRaw == nil || *importsRaw == "" || *importsRaw == "[]" {
		return normalizedImports, nil
	}
	if !strings.HasPrefix(*importsRaw, "[") || !strings.HasSuffix(*importsRaw, "]") {
		return normalizedImports, fmt.Errorf("could not parse imports from Snowflake: %s, wrapping brackets not found", *importsRaw)
	}
	raw := (*importsRaw)[1 : len(*importsRaw)-1]
	imports := strings.Split(raw, ",")
	for _, imp := range imports {
		p, err := parseFunctionOrProcedureStageLocationPath(imp)
		if err != nil {
			return nil, fmt.Errorf("could not parse imports from Snowflake: %s, err: %w", *importsRaw, err)
		}
		normalizedImports = append(normalizedImports, *p)
	}
	return normalizedImports, nil
}

func parseFunctionOrProcedureStageLocationPath(location string) (*NormalizedPath, error) {
	log.Printf("[DEBUG] parsing stage location path part: %s", location)
	idx := strings.Index(location, "/")
	if idx < 0 {
		return nil, fmt.Errorf("part %s cannot be split into stage and path", location)
	}
	stageRaw := strings.TrimPrefix(strings.TrimSpace(location[:idx]), "@")
	if stageRaw != "~" {
		stageId, err := ParseSchemaObjectIdentifier(stageRaw)
		if err != nil {
			return nil, fmt.Errorf("part %s contains incorrect stage location: %w", location, err)
		}
		stageRaw = stageId.FullyQualifiedName()
	}
	pathRaw := strings.TrimPrefix(strings.TrimSpace(location[idx:]), "/")
	if pathRaw == "" {
		return nil, fmt.Errorf("part %s contains empty path", location)
	}
	return &NormalizedPath{stageRaw, pathRaw}, nil
}

func parseFunctionOrProcedureReturns(returns string) (datatypes.DataType, bool, error) {
	var returnNotNull bool
	trimmed := strings.TrimSpace(returns)
	if strings.HasSuffix(trimmed, " NOT NULL") {
		returnNotNull = true
		trimmed = strings.TrimSuffix(trimmed, " NOT NULL")
	}
	dt, err := datatypes.ParseDataType(trimmed)
	return dt, returnNotNull, err
}

// Format in Snowflake DB is: (argName argType, argName argType, ...).
func parseFunctionOrProcedureSignature(signature string) ([]NormalizedArgument, error) {
	normalizedArguments := make([]NormalizedArgument, 0)
	trimmed := strings.TrimSpace(signature)
	if trimmed == "" {
		return normalizedArguments, fmt.Errorf("could not parse signature from Snowflake: %s, can't be empty", signature)
	}
	if trimmed == "()" {
		return normalizedArguments, nil
	}
	if !strings.HasPrefix(trimmed, "(") || !strings.HasSuffix(trimmed, ")") {
		return normalizedArguments, fmt.Errorf("could not parse signature from Snowflake: %s, wrapping parentheses not found", trimmed)
	}
	raw := (trimmed)[1 : len(trimmed)-1]
	args := strings.Split(raw, ",")

	for _, arg := range args {
		a, err := parseFunctionOrProcedureArgument(arg)
		if err != nil {
			return nil, fmt.Errorf("could not parse signature from Snowflake: %s, err: %w", trimmed, err)
		}
		normalizedArguments = append(normalizedArguments, *a)
	}
	return normalizedArguments, nil
}

// TODO [SNOW-1850370]: test with strange arg names (first integration test)
func parseFunctionOrProcedureArgument(arg string) (*NormalizedArgument, error) {
	log.Printf("[DEBUG] parsing argument: %s", arg)
	trimmed := strings.TrimSpace(arg)
	idx := strings.Index(trimmed, " ")
	if idx < 0 {
		return nil, fmt.Errorf("arg %s cannot be split into arg name, data type, and default", arg)
	}
	argName := trimmed[:idx]
	rest := strings.TrimSpace(trimmed[idx:])
	dt, err := datatypes.ParseDataType(rest)
	if err != nil {
		return nil, fmt.Errorf("arg type %s cannot be parsed, err: %w", rest, err)
	}
	return &NormalizedArgument{argName, dt}, nil
}

// TODO [SNOW-1850370]: is this combo enough? - e.g. whitespace looks to be not trimmed
func parseFunctionOrProcedureExternalAccessIntegrations(raw string) ([]AccountObjectIdentifier, error) {
	log.Printf("[DEBUG] external access integrations: %s", raw)
	return collections.MapErr(ParseCommaSeparatedStringArray(raw, false), ParseAccountObjectIdentifier)
}

// TODO [before V1]: test
func parseFunctionOrProcedurePackages(raw string) ([]string, error) {
	log.Printf("[DEBUG] packages: %s", raw)
	return collections.Map(ParseCommaSeparatedStringArray(raw, true), strings.TrimSpace), nil
}

// TODO [before V1]: unit test
func parseFunctionOrProcedureSecrets(raw string) (map[string]SchemaObjectIdentifier, error) {
	log.Printf("[DEBUG] parsing secrets: %s", raw)
	secrets := make(map[string]string)
	err := json.Unmarshal([]byte(raw), &secrets)
	if err != nil {
		return nil, fmt.Errorf("could not parse secrets from Snowflake: %s, err: %w", raw, err)
	}
	normalizedSecrets := make(map[string]SchemaObjectIdentifier)
	for k, v := range secrets {
		id, err := ParseSchemaObjectIdentifier(v)
		if err != nil {
			return nil, fmt.Errorf("could not parse secrets from Snowflake: %s, err: %w", raw, err)
		}
		normalizedSecrets[k] = id
	}
	return normalizedSecrets, nil
}
