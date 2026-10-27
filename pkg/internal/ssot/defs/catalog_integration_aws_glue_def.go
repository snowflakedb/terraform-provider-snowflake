package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-glue
var catalogIntegrationAwsGlueCatalogSourceType = ssot.NewEnum(
	"CatalogIntegrationCatalogSourceType", "CatalogIntegrationCatalogSourceTypes",
	"GLUE", "OBJECT_STORE", "POLARIS", "ICEBERG_REST", "SAP_BDC",
)

var catalogIntegrationAwsGlueTableFormat = ssot.NewEnum(
	"CatalogIntegrationTableFormat", "CatalogIntegrationTableFormats",
	"ICEBERG", "DELTA",
)

// CatalogIntegrationAwsGlueAttrs holds reusable AWS Glue catalog integration fields so SHOW and
// DESCRIBE Columns can point at the same definitions instead of restating names.
var CatalogIntegrationAwsGlueAttrs = struct {
	Enabled                ssot.Attribute
	RefreshIntervalSeconds ssot.Attribute
	Comment                ssot.Attribute
	GlueAwsRoleArn         ssot.Attribute
	GlueCatalogId          ssot.Attribute
	GlueRegion             ssot.Attribute
	CatalogNamespace       ssot.Attribute
	Name                   ssot.Attribute
	Type                   ssot.Attribute
	Category               ssot.Attribute
	CreatedOn              ssot.Attribute
	Id                     ssot.Attribute
	CatalogSource          ssot.Attribute
	TableFormat            ssot.Attribute
	GlueAwsIamUserArn      ssot.Attribute
	GlueAwsExternalId      ssot.Attribute
}{
	Enabled: ssot.Attribute{
		Name:        "enabled",
		Type:        ssot.AttrBool,
		Required:    true,
		Description: "Specifies whether the catalog integration is available for use for Iceberg tables. `true` allows users to create new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration function normally. `false` prevents users from creating new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration cannot access the catalog in the table definition.",
	},
	RefreshIntervalSeconds: ssot.Attribute{
		Name:          "refresh_interval_seconds",
		Type:          ssot.AttrInt,
		Description:   "Specifies the number of seconds to wait between attempts to poll the external Iceberg catalog for metadata updates for automated refresh. For Delta-based tables, specifies the number of seconds to wait between attempts to poll your external cloud storage for new metadata.",
		Min:           new(1),
		CreateMapping: ssot.IntAttributeCreateBuilder,
		CreateSql:     ssot.IntParameter(ssot.NoQuotes()),
	},
	Comment: ssot.Attribute{
		Name:          "comment",
		Type:          ssot.AttrString,
		Description:   "Specifies a comment for the catalog integration.",
		OptionalInDb:  true,
		CreateMapping: ssot.StringAttributeCreateBuilder,
	},
	// GlueAwsRoleArn and GlueCatalogId are required arguments of NewAwsGlueParamsRequest, not of
	// NewCreateCatalogIntegrationRequest. They have no CreateMapping: HelperMapping would emit
	// request.With<Field>, and a required attribute with no mapping is a create constructor argument.
	GlueAwsRoleArn: ssot.Attribute{
		Name:        "glue_aws_role_arn",
		Type:        ssot.AttrString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the Amazon Resource Name (ARN) of the AWS Identity and Access Management (IAM) role to assume.",
	},
	GlueCatalogId: ssot.Attribute{
		Name:        "glue_catalog_id",
		Type:        ssot.AttrString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the ID of your AWS account.",
	},
	// GlueRegion and CatalogNamespace are written onto AwsGlueCatalogSourceParams, not request.With<Field>.
	GlueRegion: ssot.Attribute{
		Name:        "glue_region",
		Type:        ssot.AttrString,
		ForceNew:    true,
		Description: "Specifies the AWS region of your AWS Glue Data Catalog. You must specify a value for this attribute if your Snowflake account is not hosted on AWS. Otherwise, the default region is the Snowflake deployment region for the account.",
	},
	CatalogNamespace: ssot.Attribute{
		Name:        "catalog_namespace",
		Type:        ssot.AttrString,
		ForceNew:    true,
		Description: "Specifies the default AWS Glue Data Catalog namespace for all Iceberg tables that you associate with the catalog integration.",
	},
	Name:      ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	Type:      ssot.Attribute{Name: "type", Type: ssot.AttrString, Readonly: true},
	Category:  ssot.Attribute{Name: "category", Type: ssot.AttrString, Readonly: true},
	CreatedOn: ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	Id:        ssot.Attribute{Name: "id", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier", Readonly: true},
	CatalogSource: ssot.Attribute{
		Name:        "catalog_source",
		Type:        ssot.AttrEnum,
		Description: "Specifies the type of catalog source. This field is used to detect when the catalog source was changed outside of Terraform and to recreate the resource when that happens.",
		Enum:        catalogIntegrationAwsGlueCatalogSourceType,
		Readonly:    true,
	},
	TableFormat:       ssot.Attribute{Name: "table_format", Type: ssot.AttrEnum, Enum: catalogIntegrationAwsGlueTableFormat, Readonly: true},
	GlueAwsIamUserArn: ssot.Attribute{Name: "glue_aws_iam_user_arn", Type: ssot.AttrString, Readonly: true},
	GlueAwsExternalId: ssot.Attribute{Name: "glue_aws_external_id", Type: ssot.AttrString, Readonly: true},
}

var catalogIntegrationAwsGlueCatalogSourceStatic = ssot.Attribute{
	Name:      "catalog_source_glue",
	SqlOnly:   true,
	CreateSql: ssot.StaticSQL("catalogSource", "CATALOG_SOURCE = GLUE"),
}

var CatalogIntegrationAwsGlue = ssot.Object{
	Name:                   "CatalogIntegrationAwsGlue",
	NamePlural:             "CatalogIntegrations",
	VariantOf:              &CatalogIntegration,
	DescribeFailRead:       true,
	IdType:                 "AccountObjectIdentifier",
	IdentityNameForceNew:   true,
	AttributesModification: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		// DESCRIBE CATALOG INTEGRATION is shared by every catalog integration variant.
		ssot.WithSdk(ssot.SdkDescribe(ssot.SdkDescribeSql("CATALOG INTEGRATION"))),
		// ShowOutputShow omitted. SHOW rows are sdk.CatalogIntegration, so tests use CatalogIntegrationShowOutput.
		// WithSchema omitted. SHOW is the shared CatalogIntegration schema alias in the ext.
		// DESCRIBE stays on the SDK schema list. SSOT would repeat sdk.CatalogIntegrationAwsGlueDetails.
		// WithModel omitted. catalog_source is added in the schema ext, which the model generator does not emit.
		// WithResourceAssert omitted. catalog_source is added in the schema ext, so the schema-based assert still emits it.
		// ObjectAssertShow omitted. There is no sdk.CatalogIntegrationAwsGlue type.
		// ObjectAssertDescribe omitted. The SDK list keeps DescribeAwsGlue.
		ssot.WithShowOutputAssert(ssot.ShowOutputDescribe()),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithHandleExternalChangesInDescribe(ssot.ExtCall()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
	},
	Attributes: []ssot.Attribute{
		CatalogIntegrationAwsGlueAttrs.Enabled,
		CatalogIntegrationAwsGlueAttrs.RefreshIntervalSeconds,
		CatalogIntegrationAwsGlueAttrs.Comment,
		CatalogIntegrationAwsGlueAttrs.GlueAwsRoleArn,
		CatalogIntegrationAwsGlueAttrs.GlueCatalogId,
		CatalogIntegrationAwsGlueAttrs.GlueRegion,
		CatalogIntegrationAwsGlueAttrs.CatalogNamespace,
	},
	Create: ssot.CreateSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-glue",
		Clauses: []*ssot.Attribute{
			&catalogIntegrationAwsGlueCatalogSourceStatic,
			&catalogIntegrationTableFormatIceberg,
			&CatalogIntegrationAwsGlueAttrs.GlueAwsRoleArn,
			&CatalogIntegrationAwsGlueAttrs.GlueCatalogId,
			&CatalogIntegrationAwsGlueAttrs.GlueRegion,
			&CatalogIntegrationAwsGlueAttrs.CatalogNamespace,
			&CatalogIntegrationAwsGlueAttrs.Enabled,
			&CatalogIntegrationAwsGlueAttrs.RefreshIntervalSeconds,
			&CatalogIntegrationAwsGlueAttrs.Comment,
		},
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CatalogIntegrationAwsGlueAttrs.Name,
			&CatalogIntegrationAwsGlueAttrs.Enabled,
			&CatalogIntegrationAwsGlueAttrs.Type,
			&CatalogIntegrationAwsGlueAttrs.Category,
			&CatalogIntegrationAwsGlueAttrs.Comment,
			&CatalogIntegrationAwsGlueAttrs.CreatedOn,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CatalogIntegrationAwsGlueAttrs.Id,
			&CatalogIntegrationAwsGlueAttrs.CatalogSource,
			&CatalogIntegrationAwsGlueAttrs.TableFormat,
			&CatalogIntegrationAwsGlueAttrs.Enabled,
			&CatalogIntegrationAwsGlueAttrs.RefreshIntervalSeconds,
			&CatalogIntegrationAwsGlueAttrs.Comment,
			&CatalogIntegrationAwsGlueAttrs.GlueAwsRoleArn,
			&CatalogIntegrationAwsGlueAttrs.GlueCatalogId,
			&CatalogIntegrationAwsGlueAttrs.GlueRegion,
			&CatalogIntegrationAwsGlueAttrs.CatalogNamespace,
			&CatalogIntegrationAwsGlueAttrs.GlueAwsIamUserArn,
			&CatalogIntegrationAwsGlueAttrs.GlueAwsExternalId,
		},
	},
}
