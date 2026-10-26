package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var StorageProvider = ssot.NewEnum(
	"StorageProvider", "StorageProviders",
	"GCS", "AZURE", "S3", "S3GOV", "S3COMPAT",
)

// ExternalVolumeAttrs holds the external volume fields that SHOW and DESCRIBE reuse.
// DESCRIBE EXTERNAL VOLUME returns a generic property list (ParseExternalVolumeDescribed in
// pkg/sdk/external_volumes_ext.go), not a paired db row: every DESCRIBE value, including
// allow_writes, is a plain string, unlike SHOW's typed bool. DescribeAllowWrites is therefore a
// separate attribute from AllowWrites. StorageLocations is ManualInSchema/skipped everywhere: its
// nested per-provider JSON shape has no ssot representation yet.
var ExternalVolumeAttrs = struct {
	Comment             ssot.Attribute
	AllowWrites         ssot.Attribute
	StorageLocation     ssot.Attribute
	Name                ssot.Attribute
	Active              ssot.Attribute
	DescribeAllowWrites ssot.Attribute
	StorageLocations    ssot.Attribute
}{
	Comment: ssot.Attribute{
		Name:          "comment",
		Type:          ssot.AttrString,
		Description:   "Specifies a comment for the external volume.",
		OptionalInDb:  true,
		CreateMapping: ssot.StringAttributeCreateBuilder,
	},
	AllowWrites: ssot.Attribute{
		Name:          "allow_writes",
		Type:          ssot.AttrBool,
		Description:   "Specifies whether write operations are allowed for the external volume; must be set to TRUE for Iceberg tables that use Snowflake as the catalog.",
		CreateMapping: ssot.BooleanStringAttributeCreateBuilder,
	},
	StorageLocation: ssot.Attribute{
		Name:     "storage_location",
		Type:     ssot.AttrList,
		Required: true,
		MinItems: new(1),
		Description: "List of named cloud storage locations in different regions and, optionally, cloud platforms. Minimum 1 required. " +
			"The order of the list is important as it impacts the active storage location, and updates will be triggered if it changes. " +
			"Note that not all parameter combinations are valid as they depend on the given storage_provider. " +
			"Consult [the docs](https://docs.snowflake.com/en/sql-reference/sql/create-external-volume#cloud-provider-parameters-cloudproviderparams) for more details on this.",
		Children: []ssot.Attribute{
			{
				Name:           "storage_location_name",
				Type:           ssot.AttrIdentifier,
				IdentifierType: "AccountObjectIdentifier",
				Required:       true,
				Description:    "Name of the storage location. Must be unique for the external volume. Do not use the name `terraform_provider_sentinel_storage_location` - this is reserved for the provider for performing update operations.",
			},
			{
				Name:        "storage_provider",
				Type:        ssot.AttrEnum,
				Required:    true,
				Description: "Specifies the cloud storage provider that stores your data files.",
				Enum:        StorageProvider,
			},
			{
				Name:        "storage_base_url",
				Type:        ssot.AttrString,
				Required:    true,
				Description: "Specifies the base URL for your cloud storage location.",
			},
			{
				Name:        "storage_aws_role_arn",
				Type:        ssot.AttrString,
				Description: "Specifies the case-sensitive Amazon Resource Name (ARN) of the AWS identity and access management (IAM) role that grants privileges on the S3 bucket containing your data files.",
			},
			{
				Name:        "storage_aws_external_id",
				Type:        ssot.AttrString,
				Description: "External ID that Snowflake uses to establish a trust relationship with AWS.",
			},
			{
				Name:        "encryption_type",
				Type:        ssot.AttrString,
				Description: "Specifies the encryption type used.",
			},
			{
				Name:        "encryption_kms_key_id",
				Type:        ssot.AttrString,
				Description: "Specifies the ID for the KMS-managed key used to encrypt files.",
			},
			{
				Name:        "storage_aws_access_point_arn",
				Type:        ssot.AttrString,
				Description: "Specifies the access point ARN for the S3 bucket containing your data files. Only applicable for S3 and S3GOV storage providers.",
			},
			{
				Name:        "use_privatelink_endpoint",
				Type:        ssot.AttrBool,
				Description: "Specifies whether to use a privatelink endpoint for the storage location. Only applicable for S3, S3GOV, and AZURE storage providers.",
			},
			{
				Name:        "azure_tenant_id",
				Type:        ssot.AttrString,
				Description: "Specifies the ID for your Office 365 tenant that the allowed and blocked storage accounts belong to.",
			},
			{
				Name:        "storage_endpoint",
				Type:        ssot.AttrString,
				Description: "Specifies the endpoint for the S3-compatible storage location. Only applicable for S3COMPAT storage provider.",
			},
			{
				Name:        "storage_aws_key_id",
				Type:        ssot.AttrString,
				Description: "Specifies the AWS key ID for the S3-compatible storage location. Only applicable for S3COMPAT storage provider.",
			},
			{
				Name:        "storage_aws_secret_key",
				Type:        ssot.AttrString,
				Description: "Specifies the AWS secret key for the S3-compatible storage location. Only applicable for S3COMPAT storage provider.",
			},
		},
	},
	Name:   ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	Active: ssot.Attribute{Name: "active", Type: ssot.AttrString, Readonly: true},
	// DescribeAllowWrites is a plain string on DESCRIBE's property list, unlike SHOW's AllowWrites
	// bool. Name differs from AllowWrites (ssot forbids a second Readonly column reusing a writable
	// CREATE/ALTER attribute's Name); ShowName restores the real "allow_writes" output key.
	DescribeAllowWrites: ssot.Attribute{Name: "describe_allow_writes", ShowName: "allow_writes", Type: ssot.AttrString, Readonly: true},
	// StorageLocations only needs to exist so the schema/assertion ManualFields/SkipFields generator
	// options can target it; its nested per-provider JSON shape is not expressible as an ssot Elem.
	StorageLocations: ssot.Attribute{
		Name: "storage_locations", Type: ssot.AttrList, Readonly: true,
		Elem: &ssot.Attribute{Type: ssot.AttrString},
	},
}

// ExternalVolumeObject is not named ExternalVolume because that identifier is the
// EXTERNAL_VOLUME parameter in the catalog.
var ExternalVolumeObject = ssot.Object{
	Name:                                "ExternalVolume",
	NamePlural:                          "ExternalVolumes",
	DescribeFailRead:                    true,
	IdType:                              "AccountObjectIdentifier",
	IdentityNameForceNew:                true,
	AttributesModification:              true,
	AdditionalDescribeOutputDescription: "Because of Terraform limitations, the changes on storage_location field do not mark this field as computed.",
	Generators: []ssot.GeneratorOption{
		ssot.WithSchema(
			ssot.SchemaShow(),
			ssot.SchemaDescribe(ssot.SchemaDescribeManualFields(&ExternalVolumeAttrs.StorageLocations)),
		),
		ssot.WithModel(
			ssot.ModelRequestType(&ExternalVolumeAttrs.StorageLocation, "[]sdk.ExternalVolumeStorageLocationRequest"),
		),
		ssot.WithResourceAssert(),
		ssot.WithShowOutputAssert(
			ssot.ShowOutputShow(),
			ssot.ShowOutputDescribe(ssot.ShowOutputDescribeSkipFields(&ExternalVolumeAttrs.StorageLocations)),
		),
		ssot.WithSdk(
			ssot.SdkShow(ssot.SdkShowOptionalLike()),
			ssot.SdkDrop(),
		),
		ssot.WithObjectAssert(
			ssot.ObjectAssertShow(),
			ssot.ObjectAssertDescribe(ssot.ObjectAssertDescribeNestedField(&ExternalVolumeAttrs.StorageLocations)),
		),
	},
	Attributes: []ssot.Attribute{
		ExternalVolumeAttrs.Comment,
		ExternalVolumeAttrs.AllowWrites,
		ExternalVolumeAttrs.StorageLocation,
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&ExternalVolumeAttrs.Name,
			&ExternalVolumeAttrs.AllowWrites,
			&ExternalVolumeAttrs.Comment,
		},
	},
	// Order matches sdk.ExternalVolumeDetails field order (Id excluded: not a Describe column).
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&ExternalVolumeAttrs.StorageLocations,
			&ExternalVolumeAttrs.Active,
			&ExternalVolumeAttrs.Comment,
			&ExternalVolumeAttrs.DescribeAllowWrites,
		},
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithDescribeOutputSet(ssot.Omitted()),
		ssot.WithSetDescribeOutput(ssot.ExtCall()),
		ssot.WithStateUpgraders(ssot.Generated()),
	},
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"v2_14_0_ExternalVolumeStateUpgrader",
	},
}
