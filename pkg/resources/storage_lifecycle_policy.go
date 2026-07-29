package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var storageLifecyclePolicySchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the storage lifecycle policy; must be unique for the database and schema in which the storage lifecycle policy is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the storage lifecycle policy."),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the storage lifecycle policy."),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"argument": {
		Type:     schema.TypeList,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The argument name.",
					ForceNew:    true,
				},
				"type": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      dataTypeFieldDescription("The argument type."),
					DiffSuppressFunc: DiffSuppressDataTypes,
					ValidateDiagFunc: IsDataTypeValid,
					StateFunc:        DataTypeStateFunc,
					ForceNew:         true,
				},
			},
		},
		Required:    true,
		Description: "List of the arguments for the storage lifecycle policy. A signature specifies a set of attributes that must be considered to determine whether the row is ready for expiration.",
		ForceNew:    true,
	},
	"body": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      diffSuppressStatementFieldDescription("Specifies the SQL expression. The expression can be any boolean-valued SQL expression."),
		DiffSuppressFunc: DiffSuppressStatement,
		ValidateFunc:     validation.StringIsNotEmpty,
	},
	"archive_tier": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToStorageLifecyclePolicyArchiveTier),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToStorageLifecyclePolicyArchiveTier),
		Description: joinWithSpace("Specifies the type of storage tier to use for archiving rows.",
			"After you set the ARCHIVE_TIER for a policy, you can’t modify it.",
			"If you don’t specify this parameter, the policy is an expiration policy that deletes rows without archiving them.",
			enumValuesDescription(sdk.AllStorageLifecyclePolicyArchiveTiers)),
	},
	"archive_for_days": {
		Type:     schema.TypeInt,
		Optional: true,
		Description: joinWithSpace("Specifies the number of days to keep rows that match the policy expression in archive storage.",
			"If set, Snowflake moves the data into archive storage according to the value you select for archive_tier.",
			"If unset, Snowflake expires the rows from the table without archiving the data."),
		ValidateFunc: validation.IntAtLeast(1),
		RequiredWith: []string{"archive_tier"},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the storage lifecycle policy.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW STORAGE LIFECYCLE POLICIES` for the given storage lifecycle policy.",
		Elem: &schema.Resource{
			Schema: schemas.ShowStorageLifecyclePolicySchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE STORAGE LIFECYCLE POLICY` for the given storage lifecycle policy.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeStorageLifecyclePolicyDetailsSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func StorageLifecyclePolicy() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.StorageLifecyclePolicies.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.StorageLifecyclePolicyResource), TrackingCreateWrapper(resources.StorageLifecyclePolicy, CreateStorageLifecyclePolicy)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.StorageLifecyclePolicyResource), TrackingReadWrapper(resources.StorageLifecyclePolicy, ReadStorageLifecyclePolicy)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.StorageLifecyclePolicyResource), TrackingUpdateWrapper(resources.StorageLifecyclePolicy, UpdateStorageLifecyclePolicy)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.StorageLifecyclePolicyResource), TrackingDeleteWrapper(resources.StorageLifecyclePolicy, deleteFunc)),
		Description:   "Resource used to manage storage lifecycle policy objects. For more information, check [storage lifecycle policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-storage-lifecycle-policy).",

		Schema: storageLifecyclePolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.StorageLifecyclePolicy, ImportName[sdk.SchemaObjectIdentifier]),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.StorageLifecyclePolicy, customdiff.All(
			ForceNewIfChangedFromNonEmptyString("archive_tier"),
			ComputedIfAnyAttributeChanged(storageLifecyclePolicySchema, ShowOutputAttributeName, "name", "schema", "database", "comment"),
			ComputedIfAnyAttributeChanged(storageLifecyclePolicySchema, DescribeOutputAttributeName, "name", "argument", "body", "archive_tier", "archive_for_days"),
			ComputedIfAnyAttributeChanged(storageLifecyclePolicySchema, FullyQualifiedNameAttributeName, "name", "schema", "database"),
		)),
		Timeouts: defaultTimeouts,
	}
}

func CreateStorageLifecyclePolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	body := d.Get("body").(string)

	args, err := HandleNestedDataTypeCreate(d, "argument", "type", func(v map[string]any, dataType datatypes.DataType) (sdk.CreateStorageLifecyclePolicyArgsRequest, error) {
		return *sdk.NewCreateStorageLifecyclePolicyArgsRequest(v["name"].(string), dataType), nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateStorageLifecyclePolicyRequest(id, args, body)

	if errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "archive_tier", request.WithArchiveTier, sdk.ToStorageLifecyclePolicyArchiveTier),
		intAttributeCreateBuilder(d, "archive_for_days", request.WithArchiveForDays),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	); errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.StorageLifecyclePolicies.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating storage lifecycle policy %v, err = %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadStorageLifecyclePolicy(ctx, d, meta)
}

func ReadStorageLifecyclePolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	storageLifecyclePolicy, err := client.StorageLifecyclePolicies.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query storage lifecycle policy. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Storage lifecycle policy id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	storageLifecyclePolicyDescription, err := client.StorageLifecyclePolicies.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	var archiveForDays int
	if storageLifecyclePolicyDescription.ArchiveForDays != nil {
		archiveForDays = *storageLifecyclePolicyDescription.ArchiveForDays
	}

	errs := errors.Join(
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		HandleNestedDataTypeSet(
			d, "argument", "type", storageLifecyclePolicyDescription.Signature,
			func(signature sdk.TableColumnSignature) datatypes.DataType { return signature.Type },
			func(signature sdk.TableColumnSignature, arg map[string]any, _ map[string]any) {
				arg["name"] = signature.Name
			},
		),
		d.Set("body", storageLifecyclePolicyDescription.Body),
		d.Set("archive_tier", storageLifecyclePolicyDescription.ArchiveTier),
		d.Set("archive_for_days", archiveForDays),
		d.Set("comment", storageLifecyclePolicy.Comment),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.StorageLifecyclePolicyToSchema(storageLifecyclePolicy)}),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.StorageLifecyclePolicyDetailsToSchema(*storageLifecyclePolicyDescription)}),
	)
	return diag.FromErr(errs)
}

func UpdateStorageLifecyclePolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		if err := client.StorageLifecyclePolicies.Alter(ctx, sdk.NewAlterStorageLifecyclePolicyRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming storage lifecycle policy from %v to %v, err = %w", id.FullyQualifiedName(), newId.FullyQualifiedName(), err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	if d.HasChange("body") {
		if err := client.StorageLifecyclePolicies.Alter(ctx, sdk.NewAlterStorageLifecyclePolicyRequest(id).WithSetBody(d.Get("body").(string))); err != nil {
			return diag.FromErr(fmt.Errorf("error updating storage lifecycle policy %v body, err = %w", id.FullyQualifiedName(), err))
		}
	}

	setRequest := sdk.NewStorageLifecyclePolicySetRequest()
	unsetRequest := sdk.NewStorageLifecyclePolicyUnsetRequest()

	if errs := errors.Join(
		attributeMappedValueUpdateSetOnly(d, "archive_tier", &setRequest.ArchiveTier, sdk.ToStorageLifecyclePolicyArchiveTier),
		intAttributeUpdate(d, "archive_for_days", &setRequest.ArchiveForDays, &unsetRequest.ArchiveForDays),
		stringAttributeUpdate(d, "comment", &setRequest.Comment, &unsetRequest.Comment),
	); errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(*setRequest, *sdk.NewStorageLifecyclePolicySetRequest()) {
		if err := client.StorageLifecyclePolicies.Alter(ctx, sdk.NewAlterStorageLifecyclePolicyRequest(id).WithSet(*setRequest)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting properties for storage lifecycle policy %v, err = %w", id.FullyQualifiedName(), err))
		}
	}

	if !reflect.DeepEqual(*unsetRequest, *sdk.NewStorageLifecyclePolicyUnsetRequest()) {
		if err := client.StorageLifecyclePolicies.Alter(ctx, sdk.NewAlterStorageLifecyclePolicyRequest(id).WithUnset(*unsetRequest)); err != nil {
			return diag.FromErr(fmt.Errorf("error unsetting properties for storage lifecycle policy %v, err = %w", id.FullyQualifiedName(), err))
		}
	}

	return ReadStorageLifecyclePolicy(ctx, d, meta)
}
