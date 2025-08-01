package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var storageIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"type": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "EXTERNAL_STAGE",
		ValidateFunc: validation.StringInSlice([]string{"EXTERNAL_STAGE"}, true),
		ForceNew:     true,
	},
	"enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"storage_allowed_locations": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Explicitly limits external stages that use the integration to reference one or more storage locations.",
		MinItems:    1,
	},
	"storage_blocked_locations": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Explicitly prohibits external stages that use the integration from referencing one or more storage locations.",
	},
	"storage_provider": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateDiagFunc: StringInSlice(sdk.AllStorageProviders, true),
		Description:      fmt.Sprintf("Specifies the storage provider for the integration. Valid options are: %s", possibleValuesListed(sdk.AllStorageProviders)),
	},
	"storage_aws_external_id": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("storage_aws_external_id"),
		Description:      "The external ID that Snowflake will use when assuming the AWS role.",
	},
	"storage_aws_iam_user_arn": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake user that will attempt to assume the AWS role.",
	},
	"storage_aws_object_acl": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"bucket-owner-full-control"}, false),
		Description:  "\"bucket-owner-full-control\" Enables support for AWS access control lists (ACLs) to grant the bucket owner full control.",
	},
	"storage_aws_role_arn": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"azure_tenant_id": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"azure_consent_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Sensitive:   true,
		Description: "The consent URL that is used to create an Azure Snowflake service principle inside your tenant.",
	},
	"azure_multi_tenant_app_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "This is the name of the Snowflake client application created for your account.",
	},
	"storage_gcp_service_account": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "This is the name of the Snowflake Google Service Account created for your account.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the storage integration was created.",
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE STORAGE INTEGRATION` for the given storage integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeStorageIntegrationSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func StorageIntegration() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.StorageIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.StorageIntegrationResource), TrackingCreateWrapper(resources.StorageIntegration, CreateStorageIntegration)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.StorageIntegrationResource), TrackingReadWrapper(resources.StorageIntegration, GetReadStorageIntegrationFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.StorageIntegrationResource), TrackingUpdateWrapper(resources.StorageIntegration, UpdateStorageIntegration)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.StorageIntegrationResource), TrackingDeleteWrapper(resources.StorageIntegration, deleteFunc)),

		Schema: storageIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(storageIntegrationSchema, DescribeOutputAttributeName, "storage_aws_external_id", "enabled", "comment"),
		),
	}
}

func GetReadStorageIntegrationFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		s, err := client.StorageIntegrations.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query storage integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Storage integration id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		if s.Category != "STORAGE" {
			return diag.FromErr(fmt.Errorf("expected %v to be a STORAGE integration, got %v", d.Id(), s.Category))
		}

		integrationProperties, err := client.StorageIntegrations.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not describe storage integration (%s), err = %w", d.Id(), err))
		}

		if withExternalChangesMarking {
			storageAwsExternalId, err := collections.FindFirst(integrationProperties, func(property sdk.StorageIntegrationProperty) bool {
				return property.Name == "STORAGE_AWS_EXTERNAL_ID"
			})
			if err == nil {
				if err = handleExternalChangesToObjectInDescribe(d,
					describeMapping{"storage_aws_external_id", "storage_aws_external_id", storageAwsExternalId.Value, storageAwsExternalId.Value, nil},
				); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if err = setStateToValuesFromConfig(d, storageIntegrationSchema, []string{
			"storage_aws_external_id",
		}); err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set("name", s.Name),
			d.Set("type", s.StorageType),
			d.Set("created_on", s.CreatedOn.String()),
			d.Set("enabled", s.Enabled),
			d.Set("comment", s.Comment),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		)

		for _, prop := range integrationProperties {
			switch prop.Name {
			case "STORAGE_PROVIDER":
				errs = errors.Join(errs, d.Set("storage_provider", prop.Value))
			case "STORAGE_ALLOWED_LOCATIONS":
				errs = errors.Join(errs, d.Set("storage_allowed_locations", strings.Split(prop.Value, ",")))
			case "STORAGE_BLOCKED_LOCATIONS":
				if prop.Value != "" {
					errs = errors.Join(errs, d.Set("storage_blocked_locations", strings.Split(prop.Value, ",")))
				}
			case "STORAGE_AWS_IAM_USER_ARN":
				errs = errors.Join(errs, d.Set("storage_aws_iam_user_arn", prop.Value))
			case "STORAGE_AWS_OBJECT_ACL":
				if prop.Value != "" {
					errs = errors.Join(errs, d.Set("storage_aws_object_acl", prop.Value))
				}
			case "STORAGE_AWS_ROLE_ARN":
				errs = errors.Join(errs, d.Set("storage_aws_role_arn", prop.Value))
			// STORAGE_AWS_EXTERNAL_ID is removed from here - handled by external changes detection above
			case "STORAGE_GCP_SERVICE_ACCOUNT":
				errs = errors.Join(errs, d.Set("storage_gcp_service_account", prop.Value))
			case "AZURE_CONSENT_URL":
				errs = errors.Join(errs, d.Set("azure_consent_url", prop.Value))
			case "AZURE_MULTI_TENANT_APP_NAME":
				errs = errors.Join(errs, d.Set("azure_multi_tenant_app_name", prop.Value))
			}
		}

		errs = errors.Join(errs,
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.DescribeStorageIntegrationToSchema(integrationProperties)}),
		)

		return diag.FromErr(errs)
	}
}

func CreateStorageIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("name").(string))
	enabled := d.Get("enabled").(bool)
	stringStorageAllowedLocations := expandStringList(d.Get("storage_allowed_locations").([]any))
	storageAllowedLocations := make([]sdk.StorageLocation, len(stringStorageAllowedLocations))
	for i, loc := range stringStorageAllowedLocations {
		storageAllowedLocations[i] = sdk.StorageLocation{
			Path: loc,
		}
	}

	req := sdk.NewCreateStorageIntegrationRequest(name, enabled, storageAllowedLocations)

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if _, ok := d.GetOk("storage_blocked_locations"); ok {
		stringStorageBlockedLocations := expandStringList(d.Get("storage_blocked_locations").([]any))
		storageBlockedLocations := make([]sdk.StorageLocation, len(stringStorageBlockedLocations))
		for i, loc := range stringStorageBlockedLocations {
			storageBlockedLocations[i] = sdk.StorageLocation{
				Path: loc,
			}
		}
		req.WithStorageBlockedLocations(storageBlockedLocations)
	}

	storageProvider := strings.ToUpper(d.Get("storage_provider").(string))

	switch {
	case slices.Contains(sdk.AllS3Protocols, sdk.S3Protocol(storageProvider)):
		s3Protocol, err := sdk.ToS3Protocol(storageProvider)
		if err != nil {
			return diag.FromErr(err)
		}

		v, ok := d.GetOk("storage_aws_role_arn")
		if !ok {
			return diag.FromErr(fmt.Errorf("if you use the S3 storage provider you must specify a storage_aws_role_arn"))
		}

		s3Params := sdk.NewS3StorageParamsRequest(s3Protocol, v.(string))
		if _, ok := d.GetOk("storage_aws_object_acl"); ok {
			s3Params.WithStorageAwsObjectAcl(d.Get("storage_aws_object_acl").(string))
		}
		if _, ok := d.GetOk("storage_aws_external_id"); ok {
			s3Params.WithStorageAwsExternalId(d.Get("storage_aws_external_id").(string))
		}
		req.WithS3StorageProviderParams(*s3Params)
	case storageProvider == "AZURE":
		v, ok := d.GetOk("azure_tenant_id")
		if !ok {
			return diag.FromErr(fmt.Errorf("if you use the Azure storage provider you must specify an azure_tenant_id"))
		}
		req.WithAzureStorageProviderParams(*sdk.NewAzureStorageParamsRequest(sdk.String(v.(string))))
	case storageProvider == "GCS":
		req.WithGCSStorageProviderParams(*sdk.NewGCSStorageParamsRequest())
	default:
		return diag.FromErr(fmt.Errorf("unexpected provider %v", storageProvider))
	}

	if err := client.StorageIntegrations.Create(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating storage integration: %w", err))
	}

	d.SetId(helpers.EncodeSnowflakeID(name))
	return GetReadStorageIntegrationFunc(false)(ctx, d, meta)
}

func UpdateStorageIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewStorageIntegrationSetRequest(), sdk.NewStorageIntegrationUnsetRequest()

	if d.HasChange("comment") {
		set.WithComment(d.Get("comment").(string))
	}

	if d.HasChange("enabled") {
		set.WithEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("storage_allowed_locations") {
		stringStorageAllowedLocations := expandStringList(d.Get("storage_allowed_locations").([]any))
		storageAllowedLocations := make([]sdk.StorageLocation, len(stringStorageAllowedLocations))
		for i, loc := range stringStorageAllowedLocations {
			storageAllowedLocations[i] = sdk.StorageLocation{
				Path: loc,
			}
		}
		set.WithStorageAllowedLocations(storageAllowedLocations)
	}

	// We need to UNSET this if we remove all storage blocked locations, because Snowflake won't accept an empty list
	if d.HasChange("storage_blocked_locations") {
		storageBlockedLocations := d.Get("storage_blocked_locations")
		if len(storageBlockedLocations.([]interface{})) > 0 {
			stringStorageBlockedLocations := expandStringList(storageBlockedLocations.([]any))
			storageBlockedLocations := make([]sdk.StorageLocation, len(stringStorageBlockedLocations))
			for i, loc := range stringStorageBlockedLocations {
				storageBlockedLocations[i] = sdk.StorageLocation{
					Path: loc,
				}
			}
			set.WithStorageBlockedLocations(storageBlockedLocations)
		} else {
			unset.WithStorageBlockedLocations(true)
		}
	}

	if d.HasChange("storage_aws_role_arn") || d.HasChange("storage_aws_object_acl") || d.HasChange("storage_aws_external_id") {
		s3SetParams := sdk.NewSetS3StorageParamsRequest(d.Get("storage_aws_role_arn").(string))

		if d.HasChange("storage_aws_object_acl") {
			if v, ok := d.GetOk("storage_aws_object_acl"); ok {
				s3SetParams.WithStorageAwsObjectAcl(v.(string))
			} else {
				unset.WithStorageAwsObjectAcl(true)
			}
		}

		if d.HasChange("storage_aws_external_id") {
			if v, ok := d.GetOk("storage_aws_external_id"); ok {
				s3SetParams.WithStorageAwsExternalId(v.(string))
			} else {
				unset.WithStorageAwsExternalId(true)
			}
		}

		set.WithS3Params(*s3SetParams)
	}

	if d.HasChange("azure_tenant_id") {
		set.WithAzureParams(*sdk.NewSetAzureStorageParamsRequest(d.Get("azure_tenant_id").(string)))
	}

	if !reflect.DeepEqual(*set, *sdk.NewStorageIntegrationSetRequest()) {
		req := sdk.NewAlterStorageIntegrationRequest(id).WithSet(*set)
		if err := client.StorageIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating storage integration, err = %w", err))
		}
	}

	if !reflect.DeepEqual(*unset, *sdk.NewStorageIntegrationUnsetRequest()) {
		req := sdk.NewAlterStorageIntegrationRequest(id).WithUnset(*unset)
		if err := client.StorageIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating storage integration, err = %w", err))
		}
	}

	return GetReadStorageIntegrationFunc(false)(ctx, d, meta)
}
