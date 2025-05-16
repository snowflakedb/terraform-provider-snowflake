package resources

import (
	"context"
	"encoding/csv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

type v085ExternalFunctionId struct {
	DatabaseName             string
	SchemaName               string
	ExternalFunctionName     string
	ExternalFunctionArgTypes string
}

func parseV085ExternalFunctionId(stringID string) (*v085ExternalFunctionId, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = '|'
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, sdk.NewError("not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, sdk.NewError("1 line at a time")
	}
	if len(lines[0]) != 4 {
		return nil, sdk.NewError("4 fields allowed")
	}

	return &v085ExternalFunctionId{
		DatabaseName:             lines[0][0],
		SchemaName:               lines[0][1],
		ExternalFunctionName:     lines[0][2],
		ExternalFunctionArgTypes: lines[0][3],
	}, nil
}

func v085ExternalFunctionStateUpgrader(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rawState == nil {
		return rawState, nil
	}

	oldId := rawState["id"].(string)
	parsedV085ExternalFunctionId, err := parseV085ExternalFunctionId(oldId)
	if err != nil {
		return nil, err
	}

	argDataTypes := make([]sdk.DataType, 0)
	if parsedV085ExternalFunctionId.ExternalFunctionArgTypes != "" {
		for _, argType := range strings.Split(parsedV085ExternalFunctionId.ExternalFunctionArgTypes, "-") {
			argDataType, err := datatypes.ParseDataType(argType)
			if err != nil {
				return nil, err
			}
			argDataTypes = append(argDataTypes, sdk.LegacyDataTypeFrom(argDataType))
		}
	}

	schemaObjectIdentifierWithArguments := sdk.NewSchemaObjectIdentifierWithArgumentsOld(parsedV085ExternalFunctionId.DatabaseName, parsedV085ExternalFunctionId.SchemaName, parsedV085ExternalFunctionId.ExternalFunctionName, argDataTypes)
	rawState["id"] = schemaObjectIdentifierWithArguments.FullyQualifiedName()

	oldDatabase := rawState["database"].(string)
	oldSchema := rawState["schema"].(string)

	rawState["database"] = strings.Trim(oldDatabase, "\"")
	rawState["schema"] = strings.Trim(oldSchema, "\"")

	if old, isPresent := rawState["return_null_allowed"]; !isPresent || old == nil {
		rawState["return_null_allowed"] = "true"
	}

	return rawState, nil
}
