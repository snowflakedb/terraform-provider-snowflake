package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var parametersSchema = map[string]*schema.Schema{
	"parameter_type": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "ACCOUNT",
		Description:  "The type of parameter to filter by. Valid values are: \"ACCOUNT\", \"SESSION\", \"OBJECT\".",
		ValidateFunc: validation.StringInSlice([]string{"ACCOUNT", "SESSION", "OBJECT"}, true),
	},
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Allows limiting the list of parameters by name using LIKE clause. Refer to [Limiting the List of Parameters by Name](https://docs.snowflake.com/en/sql-reference/parameters.html#limiting-the-list-of-parameters-by-name)",
	},
	"user": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "If parameter_type is set to \"SESSION\" then user is the name of the user to display session parameters for.",
	},
	"object_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "If parameter_type is set to \"OBJECT\" then object_type is the type of object to display object parameters for. Valid values are any object supported by the IN clause of the [SHOW PARAMETERS](https://docs.snowflake.com/en/sql-reference/sql/show-parameters.html#parameters) statement, including: WAREHOUSE | DATABASE | SCHEMA | TASK | TABLE",
	},
	"object_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "If parameter_type is set to \"OBJECT\" then object_name is the name of the object to display object parameters for.",
	},
	"parameters": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The pipes in the schema",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The name of the parameter",
				},
				"value": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The value of the parameter",
				},
				"default": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The default value of the parameter",
				},
				"level": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The level of the parameter",
				},
				"description": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The description of the parameter",
				},
			},
		},
	},
}

func Parameters() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.ParametersDatasource), TrackingReadWrapper(datasources.Parameters, ReadParameters)),
		Schema:      parametersSchema,
	}
}

func ReadParameters(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	p, ok := d.GetOk("pattern")
	pattern := ""
	if ok {
		pattern = p.(string)
	}
	var parameters []*sdk.Parameter
	var err error
	opts := sdk.ShowParametersOptions{
		In: &sdk.ParametersIn{},
	}
	if pattern != "" {
		opts.Like = &sdk.Like{Pattern: sdk.String(pattern)}
	}
	parameterType := strings.ToUpper(d.Get("parameter_type").(string))
	switch parameterType {
	case "ACCOUNT":
		opts.In.Account = sdk.Bool(true)
	case "SESSION":
		user := d.Get("user").(string)
		if user == "" {
			return diag.FromErr(fmt.Errorf("user is required when parameter_type is set to SESSION"))
		}
		opts.In.User = sdk.NewAccountObjectIdentifier(user)
	case "OBJECT":
		objectType := sdk.ObjectType(d.Get("object_type").(string))
		objectName := d.Get("object_name").(string)
		switch objectType {
		case sdk.ObjectTypeWarehouse:
			opts.In.Warehouse = sdk.NewAccountObjectIdentifier(objectName)
		case sdk.ObjectTypeDatabase:
			opts.In.Database = sdk.NewAccountObjectIdentifier(objectName)
		case sdk.ObjectTypeSchema:
			opts.In.Schema = sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(objectName)
		case sdk.ObjectTypeTask:
			opts.In.Task = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(objectName)
		case sdk.ObjectTypeTable:
			opts.In.Table = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(objectName)
		default:
			return diag.FromErr(fmt.Errorf("object_type %s is not supported", objectType))
		}
	}
	parameters, err = client.Parameters.ShowParameters(ctx, &opts)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing parameters: %w", err))
	}
	d.SetId("parameters")

	params := []map[string]interface{}{}
	for _, param := range parameters {
		paramMap := map[string]interface{}{}

		paramMap["key"] = param.Key
		paramMap["value"] = param.Value
		paramMap["default"] = param.Default
		paramMap["level"] = param.Level
		paramMap["description"] = param.Description

		params = append(params, paramMap)
	}
	return diag.FromErr(d.Set("parameters", params))
}
