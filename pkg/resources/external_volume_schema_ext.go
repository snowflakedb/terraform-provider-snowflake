package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// A list is used as the order of storage locations matter. Storage location position in the list is used to select
// the active storage location - https://docs.snowflake.com/en/user-guide/tables-iceberg-storage#active-storage-location
// This is also why it has been left as one list with optional cloud dependent parameters, rather than splitting into
// one list per cloud provider.
var externalVolumeStorageLocationSchema = &schema.Schema{
	Type:        schema.TypeList,
	Required:    true,
	MinItems:    1,
	Description: "List of named cloud storage locations in different regions and, optionally, cloud platforms. Minimum 1 required. The order of the list is important as it impacts the active storage location, and updates will be triggered if it changes. Note that not all parameter combinations are valid as they depend on the given storage_provider. Consult [the docs](https://docs.snowflake.com/en/sql-reference/sql/create-external-volume#cloud-provider-parameters-cloudproviderparams) for more details on this.",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"storage_location_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      blocklistedCharactersFieldDescription("Name of the storage location. Must be unique for the external volume. Do not use the name `terraform_provider_sentinel_storage_location` - this is reserved for the provider for performing update operations."),
				DiffSuppressFunc: suppressIdentifierQuoting,
			},
			"storage_provider": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: sdkValidation(sdk.ToStorageProvider),
				DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToStorageProvider)),
				Description:      fmt.Sprintf("Specifies the cloud storage provider that stores your data files. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllStorageProviders)),
			},
			"storage_base_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the base URL for your cloud storage location.",
			},
			"storage_aws_role_arn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the case-sensitive Amazon Resource Name (ARN) of the AWS identity and access management (IAM) role that grants privileges on the S3 bucket containing your data files.",
			},
			"storage_aws_external_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "External ID that Snowflake uses to establish a trust relationship with AWS.",
			},
			"encryption_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the encryption type used.",
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == "NONE" && newValue == ""
				},
			},
			"encryption_kms_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID for the KMS-managed key used to encrypt files.",
			},
			"storage_aws_access_point_arn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the access point ARN for the S3 bucket containing your data files. Only applicable for S3 and S3GOV storage providers.",
			},
			"use_privatelink_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          BooleanDefault,
				ValidateDiagFunc: validateBooleanString,
				Description:      booleanStringFieldDescription("Specifies whether to use a privatelink endpoint for the storage location. Only applicable for S3, S3GOV, and AZURE storage providers."),
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == "" && (newValue == BooleanDefault || newValue == BooleanFalse)
				},
			},
			"azure_tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID for your Office 365 tenant that the allowed and blocked storage accounts belong to.",
			},
			"storage_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the endpoint for the S3-compatible storage location. Only applicable for S3COMPAT storage provider.",
			},
			"storage_aws_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the AWS key ID for the S3-compatible storage location. Only applicable for S3COMPAT storage provider.",
			},
			"storage_aws_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Specifies the AWS secret key for the S3-compatible storage location. Only applicable for S3COMPAT storage provider.",
			},
		},
	},
}

func externalVolumeAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	// allow_writes is an optional boolean string, but the existing schema does not validate it.
	s["allow_writes"].ValidateDiagFunc = nil
	s["storage_location"] = externalVolumeStorageLocationSchema
	return s
}
