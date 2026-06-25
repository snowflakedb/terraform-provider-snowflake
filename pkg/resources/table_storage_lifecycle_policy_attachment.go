package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var tableStorageLifecyclePolicyAttachmentSchema = map[string]*schema.Schema{
	"table_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedPipesFieldDescription("Fully qualified name of the table (or dynamic table) the storage lifecycle policy is attached to."),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
	"table_type": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Specifies the type of the table referenced in `table_name`. " + enumValuesDescription(sdk.StorageLifecyclePolicySupportedTableTypes),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToPolicyEntityDomain),
		ValidateFunc: validation.StringInSlice(collections.Map(sdk.StorageLifecyclePolicySupportedTableTypes, func(v sdk.PolicyEntityDomain) string {
			return string(v)
		}), true),
	},
	"storage_lifecycle_policy_name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedPipesFieldDescription("Fully qualified name of the storage lifecycle policy to attach to the table."),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
	"on": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "List of the columns the storage lifecycle policy applies to.",
	},
}

func TableStorageLifecyclePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Specifies the storage lifecycle policy to attach to a table or a dynamic table.",
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.TableStorageLifecyclePolicyAttachmentResource), TrackingCreateWrapper(resources.TableStorageLifecyclePolicyAttachment, CreateTableStorageLifecyclePolicyAttachment)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.TableStorageLifecyclePolicyAttachmentResource), TrackingReadWrapper(resources.TableStorageLifecyclePolicyAttachment, ReadTableStorageLifecyclePolicyAttachment)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.TableStorageLifecyclePolicyAttachmentResource), TrackingUpdateWrapper(resources.TableStorageLifecyclePolicyAttachment, UpdateTableStorageLifecyclePolicyAttachment)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.TableStorageLifecyclePolicyAttachmentResource), TrackingDeleteWrapper(resources.TableStorageLifecyclePolicyAttachment, DeleteTableStorageLifecyclePolicyAttachment)),

		Schema: tableStorageLifecyclePolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.TableStorageLifecyclePolicyAttachment, ImportTableStorageLifecyclePolicyAttachment),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportTableStorageLifecyclePolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	tableName, err := getTableNameForTableStorageLifecyclePolicyAttachment(d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("table_name", tableName.FullyQualifiedName()); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateTableStorageLifecyclePolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	tableName, err := sdk.ParseSchemaObjectIdentifier(d.Get("table_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	tableType, err := sdk.ToPolicyEntityDomain(d.Get("table_type").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	storageLifecyclePolicyName, err := sdk.ParseSchemaObjectIdentifier(d.Get("storage_lifecycle_policy_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	columns := expandStorageLifecyclePolicyColumns(d.Get("on").([]any))

	if err := addStorageLifecyclePolicyToTable(ctx, client, tableType, tableName, storageLifecyclePolicyName, columns); err != nil {
		return diag.FromErr(fmt.Errorf("error while attaching storage lifecycle policy to %s %s, err = %w", tableType, tableName.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(tableName.FullyQualifiedName(), storageLifecyclePolicyName.FullyQualifiedName()))

	return ReadTableStorageLifecyclePolicyAttachment(ctx, d, meta)
}

func ReadTableStorageLifecyclePolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	tableName, err := getTableNameForTableStorageLifecyclePolicyAttachment(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// We use PolicyEntityDomainTable for both table variants, as the POLICY_REFERENCES function does not accept
	// DYNAMIC_TABLE as a REF_ENTITY_DOMAIN value.
	policyReferences, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(tableName, sdk.PolicyEntityDomainTable))
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to get table policies. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Table id: %s, Err: %s", tableName.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	storageLifecyclePolicyReferences := make([]sdk.PolicyReference, 0)
	for _, policyReference := range policyReferences {
		if policyReference.PolicyKind == sdk.PolicyKindStorageLifecyclePolicy {
			storageLifecyclePolicyReferences = append(storageLifecyclePolicyReferences, policyReference)
		}
	}

	if len(storageLifecyclePolicyReferences) > 1 {
		return diag.FromErr(fmt.Errorf("internal error: multiple storage lifecycle policy references attached to a table. This should never happen"))
	}

	if len(storageLifecyclePolicyReferences) == 0 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to find table's storage lifecycle policy. Marking the resource as removed.",
				Detail:   fmt.Sprintf("Table id: %s", tableName.FullyQualifiedName()),
			},
		}
	}

	storageLifecyclePolicyReference := storageLifecyclePolicyReferences[0]
	var on []string
	if storageLifecyclePolicyReference.RefArgColumnNames != nil {
		on = sdk.ParseCommaSeparatedStringArray(*storageLifecyclePolicyReference.RefArgColumnNames, true)
	}

	var policyDb, policySchema string
	if storageLifecyclePolicyReference.PolicyDb != nil {
		policyDb = *storageLifecyclePolicyReference.PolicyDb
	}
	if storageLifecyclePolicyReference.PolicySchema != nil {
		policySchema = *storageLifecyclePolicyReference.PolicySchema
	}

	errs := errors.Join(
		d.Set("table_type", storageLifecyclePolicyReference.RefEntityDomain),
		d.Set(
			"storage_lifecycle_policy_name",
			sdk.NewSchemaObjectIdentifier(
				policyDb,
				policySchema,
				storageLifecyclePolicyReference.PolicyName,
			).FullyQualifiedName()),
		d.Set("on", on),
	)
	return diag.FromErr(errs)
}

func UpdateTableStorageLifecyclePolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if d.HasChange("storage_lifecycle_policy_name") || d.HasChange("on") {
		tableName, err := sdk.ParseSchemaObjectIdentifier(d.Get("table_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		tableType, err := sdk.ToPolicyEntityDomain(d.Get("table_type").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		storageLifecyclePolicyName, err := sdk.ParseSchemaObjectIdentifier(d.Get("storage_lifecycle_policy_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		columns := expandStorageLifecyclePolicyColumns(d.Get("on").([]any))

		if err := dropStorageLifecyclePolicyFromTable(ctx, client, tableType, tableName); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error while detaching old storage lifecycle policy from %s %s, err = %w", tableType, tableName.FullyQualifiedName(), err))
		}
		if err := addStorageLifecyclePolicyToTable(ctx, client, tableType, tableName, storageLifecyclePolicyName, columns); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error while attaching new storage lifecycle policy to %s %s, err = %w", tableType, tableName.FullyQualifiedName(), err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(tableName.FullyQualifiedName(), storageLifecyclePolicyName.FullyQualifiedName()))
	}

	return ReadTableStorageLifecyclePolicyAttachment(ctx, d, meta)
}

func DeleteTableStorageLifecyclePolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	tableName, err := sdk.ParseSchemaObjectIdentifier(d.Get("table_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	tableType, err := sdk.ToPolicyEntityDomain(d.Get("table_type").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := dropStorageLifecyclePolicyFromTable(ctx, client, tableType, tableName); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func getTableNameForTableStorageLifecyclePolicyAttachment(id string) (*sdk.SchemaObjectIdentifier, error) {
	parts := helpers.ParseResourceIdentifier(id)
	if len(parts) != 2 {
		return nil, fmt.Errorf("required id format '<table_fqn>|<storage_lifecycle_policy_fqn>', but got: '%s'", id)
	}

	tableName, err := sdk.ParseSchemaObjectIdentifier(parts[0])
	if err != nil {
		return nil, err
	}
	return &tableName, nil
}

func expandStorageLifecyclePolicyColumns(raw []any) []sdk.Column {
	columns := make([]sdk.Column, len(raw))
	for i, c := range raw {
		columns[i] = sdk.Column{Value: c.(string)}
	}
	return columns
}

func addStorageLifecyclePolicyToTable(ctx context.Context, client *sdk.Client, tableType sdk.PolicyEntityDomain, tableName sdk.SchemaObjectIdentifier, storageLifecyclePolicyName sdk.SchemaObjectIdentifier, columns []sdk.Column) error {
	switch tableType {
	case sdk.PolicyEntityDomainTable:
		return client.Tables.Alter(ctx, sdk.NewAlterTableRequest(tableName).
			WithAddStorageLifecyclePolicy(sdk.NewTableAddStorageLifecyclePolicyRequest(storageLifecyclePolicyName, columns)))
	case sdk.PolicyEntityDomainDynamicTable:
		return client.DynamicTables.Alter(ctx, sdk.NewAlterDynamicTableRequest(tableName).
			WithAddStorageLifecyclePolicy(sdk.NewDynamicTableAddStorageLifecyclePolicyRequest(storageLifecyclePolicyName, columns)))
	default:
		return fmt.Errorf("unsupported table type: %s", tableType)
	}
}

func dropStorageLifecyclePolicyFromTable(ctx context.Context, client *sdk.Client, tableType sdk.PolicyEntityDomain, tableName sdk.SchemaObjectIdentifier) error {
	switch tableType {
	case sdk.PolicyEntityDomainTable:
		return client.Tables.Alter(ctx, sdk.NewAlterTableRequest(tableName).WithDropStorageLifecyclePolicy(new(true)))
	case sdk.PolicyEntityDomainDynamicTable:
		return client.DynamicTables.Alter(ctx, sdk.NewAlterDynamicTableRequest(tableName).WithDropStorageLifecyclePolicy(new(true)))
	default:
		return fmt.Errorf("unsupported table type: %s", tableType)
	}
}
