package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// icebergTableCommonSchema returns the schema fields shared by all Iceberg table resources
func icebergTableCommonSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"database": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			Description:      blocklistedCharactersFieldDescription("The database in which to create the Iceberg table."),
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"schema": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			Description:      blocklistedCharactersFieldDescription("The schema in which to create the Iceberg table."),
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"name": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the Iceberg table; must be unique for the schema in which the Iceberg table is created."),
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"comment": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies a comment for the Iceberg table.",
		},
		FullyQualifiedNameAttributeName: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fully qualified name of the resource. For more information, see [object name resolution](https://docs.snowflake.com/en/sql-reference/name-resolution).",
		},
		ShowOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW ICEBERG TABLES` for the given Iceberg table. Note that this value will be only recomputed whenever values of fields affecting the output change.",
			Elem:        &schema.Resource{Schema: schemas.ShowIcebergTableSchema},
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE ICEBERG TABLE` for the given Iceberg table.",
			Elem:        &schema.Resource{Schema: schemas.DescribeIcebergTableSchema},
		},
		ParametersAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW PARAMETERS IN ICEBERG TABLE` for the given Iceberg table.",
			Elem:        &schema.Resource{Schema: schemas.ShowIcebergTableParametersSchema},
		},
	}
}

func icebergTableDeleteFunc() schema.DeleteContextFunc {
	return ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.IcebergTables.DropSafely
		},
	)
}

func importIcebergTable(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	_, err = client.IcebergTables.ShowByIDSafely(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func readIcebergTable(ctx context.Context, d *schema.ResourceData, meta any, setExtra func(d *schema.ResourceData, table *sdk.IcebergTable) error) diag.Diagnostics {
	return readIcebergTableWithParameterHandler(ctx, d, meta, handleIcebergTableParameterRead, schemas.IcebergTableParametersToSchema, setExtra)
}

func readIcebergTableWithParameterHandler(ctx context.Context, d *schema.ResourceData, meta any, handleParameterRead func(d *schema.ResourceData, parameters []*sdk.Parameter) error, parametersToSchema func([]*sdk.Parameter, *provider.Context) map[string]any, setExtra func(d *schema.ResourceData, table *sdk.IcebergTable) error) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	table, err := client.IcebergTables.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query Iceberg table. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Iceberg table id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	details, err := client.IcebergTables.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe Iceberg table (%s), err = %w", id.FullyQualifiedName(), err))
	}

	parameters, err := client.IcebergTables.ShowParameters(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not show parameters for Iceberg table (%s), err = %w", id.FullyQualifiedName(), err))
	}

	var comment string
	if table.Comment != nil {
		comment = *table.Comment
	}

	providerCtx := meta.(*provider.Context)
	if setExtra != nil {
		if err := setExtra(d, table); err != nil {
			return diag.FromErr(err)
		}
	}
	errs := errors.Join(
		d.Set("database", table.DatabaseName),
		d.Set("schema", table.SchemaName),
		d.Set("name", table.Name),
		d.Set("comment", comment),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.IcebergTableToSchema(table)}),
		d.Set(DescribeOutputAttributeName, schemas.IcebergTableDetailsToSchema(details)),
		d.Set(ParametersAttributeName, []map[string]any{parametersToSchema(parameters, providerCtx)}),
		handleParameterRead(d, parameters),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	return nil
}

// handleIcebergTableParametersCreate populates the parameter-backed fields shared by all Iceberg
// table create requests (external_volume, catalog, replace_invalid_characters).
func handleIcebergTableParametersCreate(d *schema.ResourceData, externalVolume, catalog **sdk.AccountObjectIdentifier, replaceInvalidCharacters **bool) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterExternalVolume, externalVolume, sdk.ParseAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, sdk.IcebergTableParameterCatalog, catalog, sdk.ParseAccountObjectIdentifier),
		handleParameterCreate(d, sdk.IcebergTableParameterReplaceInvalidCharacters, replaceInvalidCharacters),
	)
}

func handleIcebergTableCommonUpdate(d *schema.ResourceData, set *sdk.IcebergTableSetPropertiesRequest, unset *sdk.IcebergTableUnsetPropertiesRequest) diag.Diagnostics {
	if errs := errors.Join(
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	); errs != nil {
		return diag.FromErr(errs)
	}
	if diags := handleParameterUpdate(d, sdk.IcebergTableParameterReplaceInvalidCharacters, &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters); len(diags) > 0 {
		return diags
	}
	return nil
}

func applyIcebergTableAlter(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, set *sdk.IcebergTableSetPropertiesRequest, unset *sdk.IcebergTableUnsetPropertiesRequest) error {
	if !reflect.DeepEqual(*set, *sdk.NewIcebergTableSetPropertiesRequest()) {
		if err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).WithSet(*set)); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(*unset, *sdk.NewIcebergTableUnsetPropertiesRequest()) {
		if err := client.IcebergTables.Alter(ctx, sdk.NewAlterIcebergTableRequest(id).WithUnset(*unset)); err != nil {
			return err
		}
	}
	return nil
}
