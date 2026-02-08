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
)

var storageIntegrationAzureSchema = map[string]*schema.Schema{
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
	"azure_tenant_id": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the ID for your Office 365 tenant that the allowed and blocked storage accounts belong to.",
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
			Schema: schemas.DescribeStorageIntegrationAzureDetailsSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func StorageIntegrationAzure() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.StorageIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.StorageIntegrationAzureResource), TrackingCreateWrapper(resources.StorageIntegrationAzure, DummyStorageIntegrationAzure)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.StorageIntegrationAzureResource), TrackingReadWrapper(resources.StorageIntegrationAzure, DummyStorageIntegrationAzure)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.StorageIntegrationAzureResource), TrackingUpdateWrapper(resources.StorageIntegrationAzure, DummyStorageIntegrationAzure)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.StorageIntegrationAzureResource), TrackingDeleteWrapper(resources.StorageIntegrationAzure, deleteFunc)),

		Schema: storageIntegrationAzureSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
		// TODO [next PR]: add CustomizeDiff logic
		// TODO [next PR]: react to external stage type change (recreate)
	}
}

func DummyStorageIntegrationAzure(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}
