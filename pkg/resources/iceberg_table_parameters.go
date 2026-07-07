package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// icebergTableParametersSchema returns the parameter-backed schema fields shared by all Iceberg
// table resources (external_volume, catalog, replace_invalid_characters). It is a builder for the
// same reason as icebergTableCommonSchema.
func icebergTableParametersSchema() map[string]*schema.Schema {
	parametersSchema := make(map[string]*schema.Schema)

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
		parametersSchema[fieldName] = &schema.Schema{
			Type:             field.Type,
			Description:      field.Description,
			Computed:         true,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: field.DiffSuppress,
		}
	}

	fieldName := strings.ToLower(string(sdk.IcebergTableParameterReplaceInvalidCharacters))
	parametersSchema[fieldName] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: enrichWithReferenceToParameterDocs(sdk.IcebergTableParameterReplaceInvalidCharacters, "Specifies whether to replace invalid UTF-8 characters with the Unicode replacement character (`�`) in query results for an Iceberg table."),
		Computed:    true,
		Optional:    true,
	}
	return parametersSchema
}

var icebergTableParametersCustomDiff = ParametersCustomDiff(
	icebergTableParametersProvider,
	parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterExternalVolume, valueTypeString, sdk.ParameterTypeTable},
	parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterCatalog, valueTypeString, sdk.ParameterTypeTable},
	parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterReplaceInvalidCharacters, valueTypeBool, sdk.ParameterTypeTable},
)

// icebergTableFromRestParametersCustomDiff extends the common Iceberg table parameters custom diff
// with the additional parameter-backed fields supported by the REST catalog create path.
var icebergTableFromRestParametersCustomDiff = ParametersCustomDiff(
	icebergTableParametersProvider,
	parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterExternalVolume, valueTypeString, sdk.ParameterTypeTable},
	parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterCatalog, valueTypeString, sdk.ParameterTypeTable},
	parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterReplaceInvalidCharacters, valueTypeBool, sdk.ParameterTypeTable},
	parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterTargetFileSize, valueTypeString, sdk.ParameterTypeTable},
	parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterStorageSerializationPolicy, valueTypeString, sdk.ParameterTypeTable},
	parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterEnableIcebergMergeOnRead, valueTypeBool, sdk.ParameterTypeTable},
	parameter[sdk.IcebergTableParameter]{sdk.IcebergTableParameterIcebergMergeOnReadBehavior, valueTypeString, sdk.ParameterTypeTable},
)

func icebergTableFromRestParametersSchema() map[string]*schema.Schema {
	parametersSchema := icebergTableParametersSchema()
	fields := []parameterDef[sdk.IcebergTableParameter]{
		{
			Name:         sdk.IcebergTableParameterTargetFileSize,
			Type:         schema.TypeString,
			ValidateDiag: StringInSlice(sdk.AsStringList(sdk.AllIcebergTableTargetFileSizes), true),
			Description:  enrichWithReferenceToParameterDocs(sdk.IcebergTableParameterTargetFileSize, fmt.Sprintf("Specifies the target file size (in bytes) used when writing the Iceberg table's Parquet files. Valid values are: %v.", sdk.AllIcebergTableTargetFileSizes)),
		},
		{
			Name:        sdk.IcebergTableParameterEnableIcebergMergeOnRead,
			Type:        schema.TypeBool,
			Description: enrichWithReferenceToParameterDocs(sdk.IcebergTableParameterEnableIcebergMergeOnRead, "Specifies whether merge-on-read is enabled for the Iceberg table."),
		},
	}
	for _, field := range fields {
		fieldName := strings.ToLower(string(field.Name))
		parametersSchema[fieldName] = &schema.Schema{
			Type:             field.Type,
			Description:      field.Description,
			Computed:         true,
			Optional:         true,
			DiffSuppressFunc: field.DiffSuppress,
			ValidateDiagFunc: field.ValidateDiag,
		}
	}

	fieldsWithForceNew := []parameterDef[sdk.IcebergTableParameter]{
		{
			Name:         sdk.IcebergTableParameterStorageSerializationPolicy,
			Type:         schema.TypeString,
			Description:  enrichWithReferenceToParameterDocs(sdk.IcebergTableParameterStorageSerializationPolicy, fmt.Sprintf("Specifies the storage serialization policy for the Iceberg table. Valid values are: %v. Cannot be changed after creation.", sdk.AllStorageSerializationPolicies)),
			ValidateDiag: StringInSlice(sdk.AsStringList(sdk.AllStorageSerializationPolicies), true),
		},
		// TODO (next PRs): this is now available in ALTER ... SET - add to sdk and make it non-force-new here.
		{
			Name:         sdk.IcebergTableParameterIcebergMergeOnReadBehavior,
			Type:         schema.TypeString,
			ValidateDiag: StringInSlice(sdk.AsStringList(sdk.AllIcebergTableIcebergMergeOnReadBehaviors), true),
			Description:  enrichWithReferenceToParameterDocs(sdk.IcebergTableParameterIcebergMergeOnReadBehavior, fmt.Sprintf("Specifies the merge-on-read behavior for the Iceberg table. Valid values are: %v. Cannot be changed after creation.", sdk.AllIcebergTableIcebergMergeOnReadBehaviors)),
		},
	}
	for _, field := range fieldsWithForceNew {
		fieldName := strings.ToLower(string(field.Name))
		parametersSchema[fieldName] = &schema.Schema{
			Type:             field.Type,
			Description:      field.Description,
			ForceNew:         true,
			Computed:         true,
			Optional:         true,
			ValidateDiagFunc: field.ValidateDiag,
		}
	}
	return parametersSchema
}

func icebergTableParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	return parametersProvider(ctx, d, meta.(*provider.Context), icebergTableParametersProviderFunc, sdk.ParseSchemaObjectIdentifier)
}

func icebergTableParametersProviderFunc(c *sdk.Client) showParametersFunc[sdk.SchemaObjectIdentifier] {
	return c.IcebergTables.ShowParameters
}

func handleIcebergTableParameterRead(d *schema.ResourceData, parameters []*sdk.Parameter) error {
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

// handleIcebergTableFromRestParameterRead extends handleIcebergTableParameterRead with the additional
// parameter-backed fields supported by the REST catalog create path.
func handleIcebergTableFromRestParameterRead(d *schema.ResourceData, parameters []*sdk.Parameter) error {
	if err := handleIcebergTableParameterRead(d, parameters); err != nil {
		return err
	}
	for _, p := range parameters {
		switch p.Key {
		case string(sdk.IcebergTableParameterTargetFileSize),
			string(sdk.IcebergTableParameterStorageSerializationPolicy),
			string(sdk.IcebergTableParameterIcebergMergeOnReadBehavior):
			if err := d.Set(strings.ToLower(p.Key), p.Value); err != nil {
				return err
			}
		case string(sdk.IcebergTableParameterEnableIcebergMergeOnRead):
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
