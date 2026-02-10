package resources

import (
	"context"
	"fmt"

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
		Description:  "Enables support for AWS access control lists (ACLs) to grant the bucket owner full control.",
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
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.StorageIntegrationAwsResource), TrackingCreateWrapper(resources.StorageIntegrationAws, DummyStorageIntegrationAws)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.StorageIntegrationAwsResource), TrackingReadWrapper(resources.StorageIntegrationAws, DummyStorageIntegrationAws)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.StorageIntegrationAwsResource), TrackingUpdateWrapper(resources.StorageIntegrationAws, DummyStorageIntegrationAws)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.StorageIntegrationAwsResource), TrackingDeleteWrapper(resources.StorageIntegrationAws, deleteFunc)),
		Description:   "Resource used to manage AWS storage integration objects. For more information, check [storage integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-storage-integration).",

		Schema: storageIntegrationAwsSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
		// TODO [next PR]: add CustomizeDiff logic
		// TODO [next PR]: react to external stage type change (recreate)
	}
}

func DummyStorageIntegrationAws(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
