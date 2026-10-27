package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var StorageProvider = ssot.NewEnum(
	"StorageProvider", "StorageProviders",
	"GCS", "AZURE", "S3", "S3GOV", "S3COMPAT",
)

var (
	S3StorageProvider = ssot.NewEnum(
		"S3StorageProvider", "S3StorageProviders",
		"S3", "S3GOV",
	)
	S3EncryptionType = ssot.NewEnum(
		"S3EncryptionType", "S3EncryptionTypes",
		"AWS_SSE_S3", "AWS_SSE_KMS", "NONE",
	)
	GCSEncryptionType = ssot.NewEnum(
		"GCSEncryptionType", "GCSEncryptionTypes",
		"GCS_SSE_KMS", "NONE",
	)
	AzureEncryptionType = ssot.NewEnum(
		"AzureEncryptionType", "AzureEncryptionTypes",
		"NONE",
	)
	S3CompatEncryptionType = ssot.NewEnum(
		"S3CompatEncryptionType", "S3CompatEncryptionTypes",
		"NONE",
	)
)

// ExternalVolumeStorageLocationAttrs holds the fields of the storage_location element. Terraform lists them as the Children
// of the element, and each one carries the CREATE clause it has as a field of the SQL. The variants of the element refer to
// them, and to the SQL-only attributes below.
var ExternalVolumeStorageLocationAttrs = struct {
	Name                   ssot.Attribute
	StorageProvider        ssot.Attribute
	StorageBaseUrl         ssot.Attribute
	AwsRoleArn             ssot.Attribute
	AwsExternalId          ssot.Attribute
	EncryptionType         ssot.Attribute
	EncryptionKmsKeyId     ssot.Attribute
	AwsAccessPointArn      ssot.Attribute
	UsePrivatelinkEndpoint ssot.Attribute
	AzureTenantId          ssot.Attribute
	StorageEndpoint        ssot.Attribute
	AwsKeyId               ssot.Attribute
	AwsSecretKey           ssot.Attribute
}{
	Name: ssot.Attribute{
		Name:           "storage_location_name",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "AccountObjectIdentifier",
		Required:       true,
		Description:    "Name of the storage location. Must be unique for the external volume. Do not use the name `terraform_provider_sentinel_storage_location` - this is reserved for the provider for performing update operations.",
		CreateSql:      ssot.StringParameter(ssot.Sql("NAME"), ssot.Required()),
	},
	StorageProvider: ssot.Attribute{
		Name:        "storage_provider",
		Type:        ssot.AttrEnum,
		Required:    true,
		Description: "Specifies the cloud storage provider that stores your data files.",
		Enum:        StorageProvider,
		CreateSql:   ssot.EnumParameter(ssot.EnumOf(S3StorageProvider), ssot.SingleQuotes(), ssot.Required()),
	},
	StorageBaseUrl: ssot.Attribute{
		Name:        "storage_base_url",
		Type:        ssot.AttrString,
		Required:    true,
		Description: "Specifies the base URL for your cloud storage location.",
		CreateSql:   ssot.Clause(),
	},
	AwsRoleArn: ssot.Attribute{
		Name:        "storage_aws_role_arn",
		Type:        ssot.AttrString,
		Description: "Specifies the case-sensitive Amazon Resource Name (ARN) of the AWS identity and access management (IAM) role that grants privileges on the S3 bucket containing your data files.",
		CreateSql:   ssot.Clause(ssot.Required()),
	},
	AwsExternalId: ssot.Attribute{
		Name:        "storage_aws_external_id",
		Type:        ssot.AttrString,
		Description: "External ID that Snowflake uses to establish a trust relationship with AWS.",
		CreateSql:   ssot.Clause(),
	},
	EncryptionType: ssot.Attribute{
		Name:        "encryption_type",
		Type:        ssot.AttrString,
		Description: "Specifies the encryption type used.",
		CreateSql:   ssot.Clause(ssot.Sql("TYPE"), ssot.FieldName("EncryptionType"), ssot.EnumOf(S3EncryptionType), ssot.SingleQuotes(), ssot.Required()),
	},
	EncryptionKmsKeyId: ssot.Attribute{
		Name:        "encryption_kms_key_id",
		Type:        ssot.AttrString,
		Description: "Specifies the ID for the KMS-managed key used to encrypt files.",
		CreateSql:   ssot.Clause(ssot.Sql("KMS_KEY_ID")),
	},
	AwsAccessPointArn: ssot.Attribute{
		Name:        "storage_aws_access_point_arn",
		Type:        ssot.AttrString,
		Description: "Specifies the access point ARN for the S3 bucket containing your data files. Only applicable for S3 and S3GOV storage providers.",
		CreateSql:   ssot.Clause(),
	},
	UsePrivatelinkEndpoint: ssot.Attribute{
		Name:        "use_privatelink_endpoint",
		Type:        ssot.AttrBool,
		Description: "Specifies whether to use a privatelink endpoint for the storage location. Only applicable for S3, S3GOV, and AZURE storage providers.",
		CreateSql:   ssot.Clause(),
	},
	AzureTenantId: ssot.Attribute{
		Name:        "azure_tenant_id",
		Type:        ssot.AttrString,
		Description: "Specifies the ID for your Office 365 tenant that the allowed and blocked storage accounts belong to.",
		CreateSql:   ssot.Clause(ssot.Required()),
	},
	StorageEndpoint: ssot.Attribute{
		Name:        "storage_endpoint",
		Type:        ssot.AttrString,
		Description: "Specifies the endpoint for the S3-compatible storage location. Only applicable for S3COMPAT storage provider.",
		CreateSql:   ssot.Clause(ssot.Required()),
	},
	AwsKeyId: ssot.Attribute{
		Name:        "storage_aws_key_id",
		Type:        ssot.AttrString,
		Description: "Specifies the AWS key ID for the S3-compatible storage location. Only applicable for S3COMPAT storage provider.",
		CreateSql:   ssot.Clause(ssot.Sql("AWS_KEY_ID"), ssot.Required()),
	},
	AwsSecretKey: ssot.Attribute{
		Name:        "storage_aws_secret_key",
		Type:        ssot.AttrString,
		Description: "Specifies the AWS secret key for the S3-compatible storage location. Only applicable for S3COMPAT storage provider.",
		CreateSql:   ssot.Clause(ssot.Sql("AWS_SECRET_KEY"), ssot.Required()),
	},
}

// The parts of the storage_location SQL that are not Terraform fields. Each is an attribute with its own CREATE clause.
var (
	externalVolumeProviderGcs      = ssot.Attribute{Name: "storage_provider_gcs", SqlOnly: true, CreateSql: ssot.StaticSQL("storageProviderGcs", "STORAGE_PROVIDER = 'GCS'")}
	externalVolumeProviderAzure    = ssot.Attribute{Name: "storage_provider_azure", SqlOnly: true, CreateSql: ssot.StaticSQL("storageProviderAzure", "STORAGE_PROVIDER = 'AZURE'")}
	externalVolumeProviderS3Compat = ssot.Attribute{Name: "storage_provider_s3_compat", SqlOnly: true, CreateSql: ssot.StaticSQL("storageProviderS3Compat", "STORAGE_PROVIDER = 'S3COMPAT'")}

	externalVolumeS3Encryption = ssot.Attribute{
		Name:      "s3_encryption",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(ssot.Sql("ENCRYPTION ="), ssot.Parentheses(), ssot.FieldName("Encryption")),
		Children:  []ssot.Attribute{ExternalVolumeStorageLocationAttrs.EncryptionType, ExternalVolumeStorageLocationAttrs.EncryptionKmsKeyId},
	}
	externalVolumeGcsEncryption = ssot.Attribute{
		Name:      "gcs_encryption",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(ssot.Sql("ENCRYPTION ="), ssot.Parentheses(), ssot.FieldName("Encryption")),
		Children: []ssot.Attribute{
			ssot.WithCreateSql(&ExternalVolumeStorageLocationAttrs.EncryptionType, ssot.Clause(ssot.Sql("TYPE"), ssot.FieldName("EncryptionType"), ssot.EnumOf(GCSEncryptionType), ssot.SingleQuotes(), ssot.Required())),
			ExternalVolumeStorageLocationAttrs.EncryptionKmsKeyId,
		},
	}
	// UPDATE STORAGE_LOCATION '<name>' CREDENTIALS = (...) of ALTER: the name and the credentials of an existing location.
	externalVolumeUpdateCredentials = ssot.Attribute{
		Name:      "update_credentials",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(ssot.Sql("CREDENTIALS ="), ssot.Parentheses(), ssot.Required(), ssot.FieldName("Credentials")),
		Children:  []ssot.Attribute{ExternalVolumeStorageLocationAttrs.AwsKeyId, ExternalVolumeStorageLocationAttrs.AwsSecretKey},
	}
	externalVolumeUpdateStorageLocation = ssot.Attribute{
		Name:      "update_storage_location",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(),
		Children: []ssot.Attribute{
			ssot.WithCreateSql(&ExternalVolumeStorageLocationAttrs.Name, ssot.StringParameter(ssot.Sql("STORAGE_LOCATION"), ssot.NoEquals(), ssot.Required())),
			externalVolumeUpdateCredentials,
		},
	}
	externalVolumeS3CompatCredentials = ssot.Attribute{
		Name:      "s3_compat_credentials",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(ssot.Sql("CREDENTIALS ="), ssot.Parentheses(), ssot.Required(), ssot.FieldName("Credentials")),
		Children:  []ssot.Attribute{ExternalVolumeStorageLocationAttrs.AwsKeyId, ExternalVolumeStorageLocationAttrs.AwsSecretKey},
	}
)

// ExternalVolumeAttrs holds the external volume fields that SHOW and DESCRIBE reuse.
// DESCRIBE EXTERNAL VOLUME returns a generic property list (DescribeRows), which the generated
// ParseExternalVolumeDescribed parses into the details that Describe.Columns models. Every DESCRIBE value, including allow_writes, is a plain string, unlike
// SHOW's typed bool, so DescribeAllowWrites is a separate attribute from AllowWrites.
// StorageLocations is skipped or manual in the generators: its nested per-provider JSON shape has
// no ssot representation yet.
var ExternalVolumeAttrs = struct {
	Comment             ssot.Attribute
	AllowWrites         ssot.Attribute
	StorageLocation     ssot.Attribute
	Name                ssot.Attribute
	Active              ssot.Attribute
	DescribeAllowWrites ssot.Attribute
	StorageLocations    ssot.Attribute
	// The DESCRIBE rows are a generic property list, which the SDK parses into the details the Describe columns model.
	ParentProperty  ssot.Attribute
	Property        ssot.Attribute
	PropertyType    ssot.Attribute
	PropertyValue   ssot.Attribute
	PropertyDefault ssot.Attribute
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
		Name:      "storage_location",
		Type:      ssot.AttrList,
		Required:  true,
		MinItems:  new(1),
		CreateSql: ssot.ListOfObjects(ssot.Sql("STORAGE_LOCATIONS"), ssot.Parentheses()),
		Description: "List of named cloud storage locations in different regions and, optionally, cloud platforms. Minimum 1 required. " +
			"The order of the list is important as it impacts the active storage location, and updates will be triggered if it changes. " +
			"Note that not all parameter combinations are valid as they depend on the given storage_provider. " +
			"Consult [the docs](https://docs.snowflake.com/en/sql-reference/sql/create-external-volume#cloud-provider-parameters-cloudproviderparams) for more details on this.",
		Elem: &ssot.Attribute{
			Type: ssot.AttrObject,
			CreateSql: ssot.ObjectParameter(
				ssot.Parentheses(),
				ssot.WithVariants(
					ssot.Variant("S3", &ExternalVolumeStorageLocationAttrs.StorageProvider, &ExternalVolumeStorageLocationAttrs.AwsRoleArn, &ExternalVolumeStorageLocationAttrs.StorageBaseUrl, &ExternalVolumeStorageLocationAttrs.AwsExternalId, &ExternalVolumeStorageLocationAttrs.AwsAccessPointArn, &ExternalVolumeStorageLocationAttrs.UsePrivatelinkEndpoint, &externalVolumeS3Encryption),
					ssot.Variant("GCS", &externalVolumeProviderGcs, &ExternalVolumeStorageLocationAttrs.StorageBaseUrl, &externalVolumeGcsEncryption),
					ssot.Variant("Azure", &externalVolumeProviderAzure, &ExternalVolumeStorageLocationAttrs.AzureTenantId, &ExternalVolumeStorageLocationAttrs.StorageBaseUrl, &ExternalVolumeStorageLocationAttrs.UsePrivatelinkEndpoint),
					ssot.Variant("S3Compat", &externalVolumeProviderS3Compat, &ExternalVolumeStorageLocationAttrs.StorageBaseUrl, &ExternalVolumeStorageLocationAttrs.StorageEndpoint, &externalVolumeS3CompatCredentials),
				),
			),
			Children: []ssot.Attribute{
				ExternalVolumeStorageLocationAttrs.Name,
				ExternalVolumeStorageLocationAttrs.StorageProvider,
				ExternalVolumeStorageLocationAttrs.StorageBaseUrl,
				ExternalVolumeStorageLocationAttrs.AwsRoleArn,
				ExternalVolumeStorageLocationAttrs.AwsExternalId,
				ExternalVolumeStorageLocationAttrs.EncryptionType,
				ExternalVolumeStorageLocationAttrs.EncryptionKmsKeyId,
				ExternalVolumeStorageLocationAttrs.AwsAccessPointArn,
				ExternalVolumeStorageLocationAttrs.UsePrivatelinkEndpoint,
				ExternalVolumeStorageLocationAttrs.AzureTenantId,
				ExternalVolumeStorageLocationAttrs.StorageEndpoint,
				ExternalVolumeStorageLocationAttrs.AwsKeyId,
				ExternalVolumeStorageLocationAttrs.AwsSecretKey,
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
	ParentProperty:  ssot.Attribute{Name: "parent_property", Type: ssot.AttrString, Readonly: true},
	Property:        ssot.Attribute{Name: "property", Type: ssot.AttrString, Readonly: true},
	PropertyType:    ssot.Attribute{Name: "property_type", Type: ssot.AttrString, Readonly: true},
	PropertyValue:   ssot.Attribute{Name: "property_value", Type: ssot.AttrString, Readonly: true},
	PropertyDefault: ssot.Attribute{Name: "property_default", Type: ssot.AttrString, Readonly: true},
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
		ssot.WithResource(),
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
			ssot.SdkDescribe(
				ssot.SdkDescribeDetails(
					ssot.SdkDescribeDetailsNestedField(&ExternalVolumeAttrs.StorageLocations),
					ssot.SdkDescribeDetailsMethod(),
					ssot.SdkDescribeDetailsParse(
						ssot.DescribeProperty(&ExternalVolumeAttrs.Comment, "COMMENT"),
						ssot.DescribeProperty(&ExternalVolumeAttrs.Active, "ACTIVE"),
						ssot.DescribeProperty(&ExternalVolumeAttrs.DescribeAllowWrites, "ALLOW_WRITES"),
						ssot.DescribePropertyPrefix(&ExternalVolumeAttrs.StorageLocations, "STORAGE_LOCATION_"),
					),
				),
			),
			ssot.SdkDefineFromSsot(),
			ssot.SdkCreate(),
			ssot.SdkAlter(),
		),
		ssot.WithObjectAssert(
			ssot.ObjectAssertShow(),
			ssot.ObjectAssertDescribe(ssot.ObjectAssertDescribeNestedField(&ExternalVolumeAttrs.StorageLocations)),
		),
	},
	// The Azure and S3-compatible encryption types are validated by the resource, and the SQL does not send them.
	AdditionalEnums: []*ssot.Enum{AzureEncryptionType, S3CompatEncryptionType},
	Attributes: []ssot.Attribute{
		ExternalVolumeAttrs.Comment,
		ExternalVolumeAttrs.AllowWrites,
		ExternalVolumeAttrs.StorageLocation,
	},
	Create: ssot.CreateSpec{
		Clauses: []*ssot.Attribute{
			&ExternalVolumeAttrs.StorageLocation,
			&ExternalVolumeAttrs.AllowWrites,
			&ExternalVolumeAttrs.Comment,
		},
	},
	Set: []*ssot.Attribute{
		&ExternalVolumeAttrs.AllowWrites,
		&ExternalVolumeAttrs.Comment,
	},
	Alter: ssot.AlterSpec{
		Clauses: []ssot.AlterClause{
			ssot.AlterText("REMOVE STORAGE_LOCATION"),
			ssot.AlterSet(),
			ssot.AlterParameterBlock("AddStorageLocation", "ADD STORAGE_LOCATION", &ExternalVolumeAttrs.StorageLocation),
			ssot.AlterKeywordBlock("UpdateStorageLocation", "UPDATE", &externalVolumeUpdateStorageLocation),
		},
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
	DescribeRows: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&ExternalVolumeAttrs.ParentProperty,
			&ExternalVolumeAttrs.Property,
			&ExternalVolumeAttrs.PropertyType,
			&ExternalVolumeAttrs.PropertyValue,
			&ExternalVolumeAttrs.PropertyDefault,
		},
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithSetDescribeOutput(ssot.ExtCall()),
		ssot.WithStateUpgraders(ssot.Generated()),
		ssot.WithSetConfigFields(ssot.Generated()),
	},
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"v2_14_0_ExternalVolumeStateUpgrader",
	},
}
