package resources

import (
	"context"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	icebergTableFromFilesParametersSchema     = make(map[string]*schema.Schema)
	icebergTableFromFilesParametersCustomDiff = ParametersCustomDiff(
		icebergTableFromFilesParametersProvider,
		parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterExternalVolume, valueTypeString, sdk.ParameterTypeTable},
		parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterCatalog, valueTypeString, sdk.ParameterTypeTable},
		parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterReplaceInvalidCharacters, valueTypeBool, sdk.ParameterTypeTable},
	)
)

func init() {
	forceNewFields := []parameterDef[sdk.IcebergTableParameter]{
		{
			Name:         sdk.IcebergTableParameterExternalVolume,
			Type:         schema.TypeString,
			Description:  "Specifies the identifier for the external volume where the Iceberg table stores its metadata files and data in Parquet format. If not specified, the account-level default is used.",
			DiffSuppress: suppressIdentifierQuoting,
		},
		{
			Name:         sdk.IcebergTableParameterCatalog,
			Type:         schema.TypeString,
			Description:  "Specifies the identifier for the catalog integration to use for the Iceberg table. If not specified, the account-level default is used.",
			DiffSuppress: suppressIdentifierQuoting,
		},
	}
	for _, field := range forceNewFields {
		fieldName := strings.ToLower(string(field.Name))
		icebergTableFromFilesParametersSchema[fieldName] = &schema.Schema{
			Type:             field.Type,
			Description:      field.Description,
			Computed:         true,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: field.DiffSuppress,
		}
	}

	fieldName := strings.ToLower(string(sdk.IcebergTableParameterReplaceInvalidCharacters))
	icebergTableFromFilesParametersSchema[fieldName] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: enrichWithReferenceToParameterDocs(sdk.IcebergTableParameterReplaceInvalidCharacters, "Specifies whether to replace invalid UTF-8 characters with the Unicode replacement character (`�`) in query results for an Iceberg table."),
		Computed:    true,
		Optional:    true,
	}
}

func icebergTableFromFilesParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	return parametersProvider(ctx, d, meta.(*provider.Context), icebergTableFromFilesParametersProviderFunc, sdk.ParseSchemaObjectIdentifier)
}

func icebergTableFromFilesParametersProviderFunc(c *sdk.Client) showParametersFunc[sdk.SchemaObjectIdentifier] {
	return c.IcebergTables.ShowParameters
}

func handleIcebergTableFromFilesParameterRead(d *schema.ResourceData, parameters []*sdk.Parameter) error {
	for _, p := range parameters {
		switch p.Key {
		case string(sdk.IcebergTableParameterExternalVolume),
			string(sdk.IcebergTableParameterCatalog):
			if err := d.Set(strings.ToLower(p.Key), p.Value); err != nil {
				return err
			}
		case string(sdk.IcebergTableParameterReplaceInvalidCharacters):
			value, err := strconv.ParseBool(p.Value)
			if err != nil {
				return err
			}
			if err := d.Set(strings.ToLower(p.Key), value); err != nil {
				return err
			}
		}
	}
	return nil
}
