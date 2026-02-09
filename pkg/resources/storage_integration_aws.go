package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var storageIntegrationAwsSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: blocklistedCharactersFieldDescription("String that specifies the identifier (i.e. name) for the integration; must be unique in your account."),
	},
	"enabled": {
		Type:     schema.TypeBool,
		Required: true,
		Description: joinWithSpace(
			"Specifies whether this storage integration is available for usage in stages.",
			"`TRUE` allows users to create new stages that reference this integration. Existing stages that reference this integration function normally.",
			"`FALSE` prevents users from creating new stages that reference this integration. Existing stages that reference this integration cannot access the storage location in the stage definition.",
		),
	},
	"storage_provider": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateDiagFunc: StringInSlice(sdk.AllStorageProviders, true),
		Description:      fmt.Sprintf("Specifies the storage provider for the integration. Valid options are: %s", possibleValuesListed(sdk.AllStorageProviders)),
	},
	// TODO [this PR]: change to sets?
	"storage_allowed_locations": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Explicitly limits external stages that use the integration to reference one or more storage locations.",
		MinItems:    1,
	},
	// TODO [this PR]: change to sets?
	"storage_blocked_locations": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Explicitly prohibits external stages that use the integration from referencing one or more storage locations.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the storage integration.",
	},
	"use_privatelink_endpoint": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		// TODO [next PR]: verify the DiffSuppressFunc
		// DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("use_privatelink_endpoint"),
		Description: booleanStringFieldDescription("Specifies whether to use outbound private connectivity to harden the security posture."),
	},
	"storage_aws_role_arn": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the Amazon Resource Name (ARN) of the AWS identity and access management (IAM) role that grants privileges on the S3 bucket containing your data files.",
	},
	"storage_aws_external_id": {
		Type:     schema.TypeString,
		Optional: true,
		// TODO [next PR]: verify the DiffSuppressFunc
		// DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("storage_aws_external_id"),
		Description: "Optionally specifies an external ID that Snowflake uses to establish a trust relationship with AWS.",
	},
	"storage_aws_object_acl": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"bucket-owner-full-control"}, false),
		Description:  "Enables support for AWS access control lists (ACLs) to grant the bucket owner full control. `bucket-owner-full-control` is the only currently supported value.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW STORAGE INTEGRATIONS` for the given storage integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowStorageIntegrationSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE STORAGE INTEGRATION` for the given storage integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeStorageIntegrationAwsDetailsSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func StorageIntegrationAws() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.StorageIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.StorageIntegrationAwsResource), TrackingCreateWrapper(resources.StorageIntegrationAws, CreateStorageIntegrationAws)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.StorageIntegrationAwsResource), TrackingReadWrapper(resources.StorageIntegrationAws, GetReadStorageIntegrationAwsFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.StorageIntegrationAwsResource), TrackingUpdateWrapper(resources.StorageIntegrationAws, UpdateStorageIntegrationAws)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.StorageIntegrationAwsResource), TrackingDeleteWrapper(resources.StorageIntegrationAws, deleteFunc)),
		Description:   "Resource used to manage AWS storage integration objects. For more information, check [storage integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-storage-integration).",

		Schema: storageIntegrationAwsSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.StorageIntegrationAws, ImportStorageIntegrationAws),
		},
		Timeouts: defaultTimeouts,
		// TODO [next PR]: add CustomizeDiff logic
		// TODO [next PR]: react to external stage type change (recreate)
	}
}

// TODO [next PR]: errors when importing the wrong type/category/provider
func ImportStorageIntegrationAws(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	awsDetails, err := client.StorageIntegrations.DescribeAwsDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	errs := errors.Join(
		d.Set("name", awsDetails.ID().Name()),
		d.Set("use_privatelink_endpoint", booleanStringFromBool(awsDetails.UsePrivatelinkEndpoint)),
		d.Set("storage_aws_external_id", awsDetails.ExternalId),
	)
	if errs != nil {
		return nil, errs
	}
	return []*schema.ResourceData{d}, nil
}

func GetReadStorageIntegrationAwsFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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
						Summary:  "Failed to query aws storage integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Aws storage integration id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		// TODO [next PR]: replace with force?
		if s.Category != "STORAGE" {
			return diag.FromErr(fmt.Errorf("expected %v to be a STORAGE integration, got %v", d.Id(), s.Category))
		}

		awsDetails, err := client.StorageIntegrations.DescribeAwsDetails(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not describe aws storage integration (%s), err = %w", d.Id(), err))
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"external_id", "storage_aws_external_id", awsDetails.ExternalId, awsDetails.ExternalId, nil},
				outputMapping{"use_privatelink_endpoint", "use_privatelink_endpoint", awsDetails.UsePrivatelinkEndpoint, awsDetails.UsePrivatelinkEndpoint, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			// not reading name on purpose (we never update the name externally)
			d.Set("storage_provider", awsDetails.Provider),
			d.Set("enabled", s.Enabled),
			d.Set("storage_allowed_locations", awsDetails.AllowedLocations),
			d.Set("storage_blocked_locations", awsDetails.BlockedLocations),
			d.Set("comment", s.Comment),
			// not reading use_privatelink_endpoint on purpose (handled as external change to describe output)
			d.Set("storage_aws_role_arn", awsDetails.RoleArn),
			// not reading storage_aws_external_id on purpose (handled as external change to describe output)
			d.Set("storage_aws_object_acl", awsDetails.ObjectAcl),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		)

		errs = errors.Join(errs,
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.StorageIntegrationToSchema(s)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.StorageIntegrationAwsDetailsToSchema(awsDetails)}),
		)

		return diag.FromErr(errs)
	}
}

func CreateStorageIntegrationAws(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(name)
	enabled := d.Get("enabled").(bool)
	storageAllowedLocations, _ := parseLocations(d.Get("storage_allowed_locations").([]any))
	storageProvider := strings.ToUpper(d.Get("storage_provider").(string))
	s3Protocol, err := sdk.ToS3Protocol(storageProvider)
	if err != nil {
		return diag.FromErr(err)
	}
	awsRoleArn := d.Get("storage_aws_role_arn").(string)

	request := sdk.NewCreateStorageIntegrationRequest(id, enabled, storageAllowedLocations)
	awsRequest := sdk.NewS3StorageParamsRequest(s3Protocol, awsRoleArn)
	errs := errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		func() error {
			if _, ok := d.GetOk("storage_blocked_locations"); ok {
				storageBlockedLocations, _ := parseLocations(d.Get("storage_blocked_locations").([]any))
				request.WithStorageBlockedLocations(storageBlockedLocations)
			}
			return nil
		}(),
		booleanStringAttributeCreateBuilder(d, "use_privatelink_endpoint", awsRequest.WithUsePrivatelinkEndpoint),
		stringAttributeCreateBuilder(d, "storage_aws_external_id", awsRequest.WithStorageAwsExternalId),
		stringAttributeCreateBuilder(d, "storage_aws_object_acl", awsRequest.WithStorageAwsObjectAcl),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err = client.StorageIntegrations.Create(ctx, request.WithS3StorageProviderParams(*awsRequest)); err != nil {
		return diag.FromErr(fmt.Errorf("error creating storage integration aws: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return GetReadStorageIntegrationAwsFunc(false)(ctx, d, meta)
}

func UpdateStorageIntegrationAws(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewStorageIntegrationSetRequest(), sdk.NewStorageIntegrationUnsetRequest()
	s3SetParams, s3UnsetParams := sdk.NewSetS3StorageParamsRequest(), sdk.NewUnsetS3StorageParamsRequest()

	errs := errors.Join(
		booleanAttributeUpdateSetOnly(d, "enabled", &set.Enabled),
		// TODO [next PRs]: extract helpers for lists with builders
		func() error {
			if d.HasChange("storage_allowed_locations") {
				if v, ok := d.GetOk("storage_allowed_locations"); ok {
					locations, err := parseLocations(v.([]any))
					if err != nil {
						return err
					}
					set.WithStorageAllowedLocations(locations)
				}
			}
			return nil
		}(),
		func() error {
			if d.HasChange("storage_blocked_locations") {
				v := d.Get("storage_blocked_locations")
				if len(v.([]any)) > 0 {
					locations, err := parseLocations(v.([]any))
					if err != nil {
						return err
					}
					set.WithStorageBlockedLocations(locations)
				} else {
					unset.WithStorageBlockedLocations(true)
				}
			}
			return nil
		}(),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		// TODO(SNOW-2356049): implement & use booleanStringAttributeUnsetBuilder when UNSET starts working correctly
		booleanStringAttributeUnsetFallbackUpdateBuilder(d, "use_privatelink_endpoint", s3SetParams.WithUsePrivatelinkEndpoint, false),
		stringAttributeUpdateSetOnlyNotEmpty(d, "storage_aws_role_arn", &s3SetParams.StorageAwsRoleArn),
		stringAttributeUpdate(d, "storage_aws_external_id", &s3SetParams.StorageAwsExternalId, &s3UnsetParams.StorageAwsExternalId),
		stringAttributeUpdate(d, "storage_aws_object_acl", &s3SetParams.StorageAwsObjectAcl, &s3UnsetParams.StorageAwsObjectAcl),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(*s3SetParams, *sdk.NewSetS3StorageParamsRequest()) {
		set.WithS3Params(*s3SetParams)
	}
	if !reflect.DeepEqual(*set, *sdk.NewStorageIntegrationSetRequest()) {
		req := sdk.NewAlterStorageIntegrationRequest(id).WithSet(*set)
		if err := client.StorageIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating aws storage integration, err = %w", err))
		}
	}

	if !reflect.DeepEqual(*s3UnsetParams, *sdk.NewUnsetS3StorageParamsRequest()) {
		unset.WithS3Params(*s3UnsetParams)
	}
	if !reflect.DeepEqual(*unset, *sdk.NewStorageIntegrationUnsetRequest()) {
		req := sdk.NewAlterStorageIntegrationRequest(id).WithUnset(*unset)
		if err := client.StorageIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating aws storage integration, err = %w", err))
		}
	}

	return GetReadStorageIntegrationAwsFunc(false)(ctx, d, meta)
}

func parseLocations(allowedLocationsRaw []any) ([]sdk.StorageLocation, error) {
	stringStorageAllowedLocations := expandStringList(allowedLocationsRaw)
	storageAllowedLocations := make([]sdk.StorageLocation, len(stringStorageAllowedLocations))
	for i, loc := range stringStorageAllowedLocations {
		storageAllowedLocations[i] = sdk.StorageLocation{
			Path: loc,
		}
	}
	return storageAllowedLocations, nil
}
