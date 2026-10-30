package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var CatalogIntegrationCatalogSourceType = ssot.NewEnum(
	"CatalogIntegrationCatalogSourceType", "CatalogIntegrationCatalogSourceTypes",
	"GLUE", "OBJECT_STORE", "POLARIS", "ICEBERG_REST", "SAP_BDC",
)

var CatalogIntegrationTableFormat = ssot.NewEnum(
	"CatalogIntegrationTableFormat", "CatalogIntegrationTableFormats",
	"ICEBERG", "DELTA",
)

var CatalogIntegrationCatalogApiType = ssot.NewEnum(
	"CatalogIntegrationCatalogApiType", "CatalogIntegrationCatalogApiTypes",
	"PUBLIC", "PRIVATE", "AWS_API_GATEWAY", "AWS_PRIVATE_API_GATEWAY", "AWS_GLUE", "AWS_PRIVATE_GLUE", "AWS_S3TABLES", "AWS_PRIVATE_S3TABLES",
)

var CatalogIntegrationAccessDelegationMode = ssot.NewEnum(
	"CatalogIntegrationAccessDelegationMode", "CatalogIntegrationAccessDelegationModes",
	"VENDED_CREDENTIALS", "EXTERNAL_VOLUME_CREDENTIALS",
)

var CatalogIntegrationRestAuthenticationType = ssot.NewEnum(
	"CatalogIntegrationRestAuthenticationType", "CatalogIntegrationRestAuthenticationTypes",
	"OAUTH", "BEARER", "SIGV4",
)

var catalogIntegrationTableFormatIceberg = ssot.Attribute{
	Name:      "table_format_iceberg",
	SqlOnly:   true,
	CreateSql: ssot.StaticSQL("tableFormat", "TABLE_FORMAT = ICEBERG"),
}

// catalogIntegrationRestAuthTypeOauth is the fixed TYPE of an OAuth REST_AUTHENTICATION struct.
var catalogIntegrationRestAuthTypeOauth = ssot.Attribute{
	Name:      "type_oauth",
	SqlOnly:   true,
	CreateSql: ssot.StaticSQL("restAuthType", "TYPE = OAUTH"),
}

// catalogIntegrationSetAttrs are the fields that SET changes on every catalog integration variant. They only carry the SQL of
// the shared ALTER, the variants describe the same fields for Terraform. A REST authentication is set by its secret only.
var catalogIntegrationSetAttrs = struct {
	Enabled                  ssot.Attribute
	RefreshIntervalSeconds   ssot.Attribute
	Comment                  ssot.Attribute
	OauthRestAuthentication  ssot.Attribute
	BearerRestAuthentication ssot.Attribute
}{
	Enabled:                ssot.Attribute{Name: "enabled", Type: ssot.AttrBool},
	RefreshIntervalSeconds: ssot.Attribute{Name: "refresh_interval_seconds", Type: ssot.AttrInt, SetSql: ssot.IntParameter(ssot.NoQuotes())},
	// TODO(SNOW-3121221): use COMMENT in unset and here use a plain string.
	Comment: ssot.Attribute{Name: "comment", Type: ssot.AttrString, SetSql: ssot.StringParameter(ssot.AllowEmpty())},
	OauthRestAuthentication: ssot.Attribute{
		Name:    "set_oauth_rest_authentication",
		Type:    ssot.AttrObject,
		SqlOnly: true,
		SetSql:  ssot.ObjectParameter(ssot.Sql("REST_AUTHENTICATION ="), ssot.Parentheses(), ssot.FieldName("SetOAuthRestAuthentication")),
		Children: []ssot.Attribute{
			{Name: "oauth_client_secret", Type: ssot.AttrString, Required: true},
		},
	},
	BearerRestAuthentication: ssot.Attribute{
		Name:    "set_bearer_rest_authentication",
		Type:    ssot.AttrObject,
		SqlOnly: true,
		SetSql:  ssot.ObjectParameter(ssot.Sql("REST_AUTHENTICATION ="), ssot.Parentheses(), ssot.FieldName("SetBearerRestAuthentication")),
		Children: []ssot.Attribute{
			{Name: "bearer_token", Type: ssot.AttrString, Required: true},
		},
	},
}

// catalogIntegrationShowAttrs are the SHOW columns, which are the same for every catalog integration variant.
var catalogIntegrationShowAttrs = struct {
	Name      ssot.Attribute
	Enabled   ssot.Attribute
	Type      ssot.Attribute
	Category  ssot.Attribute
	Comment   ssot.Attribute
	CreatedOn ssot.Attribute
}{
	Name:      ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	Enabled:   ssot.Attribute{Name: "enabled", Type: ssot.AttrBool, Readonly: true},
	Type:      ssot.Attribute{Name: "type", Type: ssot.AttrString, Readonly: true},
	Category:  ssot.Attribute{Name: "category", Type: ssot.AttrString, Readonly: true},
	Comment:   ssot.Attribute{Name: "comment", Type: ssot.AttrString, Readonly: true, OptionalInDb: true},
	CreatedOn: ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
}

// catalogIntegrationDescribeRowAttrs are the columns of a DESCRIBE row. Snowflake describes every catalog integration as a
// list of properties, which the details of each variant are parsed from.
var catalogIntegrationDescribeRowAttrs = struct {
	Property        ssot.Attribute
	PropertyType    ssot.Attribute
	PropertyValue   ssot.Attribute
	PropertyDefault ssot.Attribute
}{
	Property:        ssot.Attribute{Name: "property", Type: ssot.AttrString, Readonly: true},
	PropertyType:    ssot.Attribute{Name: "property_type", Type: ssot.AttrString, Readonly: true},
	PropertyValue:   ssot.Attribute{Name: "property_value", Type: ssot.AttrString, Readonly: true},
	PropertyDefault: ssot.Attribute{Name: "property_default", Type: ssot.AttrString, Readonly: true},
}

// readonlyColumn is the attribute as a DESCRIBE column of the base: the same field, read only. It is not required, so an
// object column is optional in the details of the base, which has the fields of every variant.
func readonlyColumn(attr ssot.Attribute) ssot.Attribute {
	attr.Readonly = true
	attr.Required = false
	return attr
}

// catalogIntegrationDescribeColumns are the DESCRIBE columns of the base: the union of the columns of the variants.
// Each is the attribute of a variant.
var catalogIntegrationDescribeColumns = struct {
	Enabled                  ssot.Attribute
	RefreshIntervalSeconds   ssot.Attribute
	Comment                  ssot.Attribute
	GlueAwsRoleArn           ssot.Attribute
	GlueCatalogId            ssot.Attribute
	GlueRegion               ssot.Attribute
	CatalogNamespace         ssot.Attribute
	RestConfig               ssot.Attribute
	OauthRestAuthentication  ssot.Attribute
	BearerRestAuthentication ssot.Attribute
	Sigv4RestAuthentication  ssot.Attribute
}{
	Enabled:                  readonlyColumn(CatalogIntegrationIcebergRestAttrs.Enabled),
	RefreshIntervalSeconds:   readonlyColumn(CatalogIntegrationIcebergRestAttrs.RefreshIntervalSeconds),
	Comment:                  readonlyColumn(CatalogIntegrationIcebergRestAttrs.Comment),
	GlueAwsRoleArn:           readonlyColumn(CatalogIntegrationAwsGlueAttrs.GlueAwsRoleArn),
	GlueCatalogId:            readonlyColumn(CatalogIntegrationAwsGlueAttrs.GlueCatalogId),
	GlueRegion:               readonlyColumn(CatalogIntegrationAwsGlueAttrs.GlueRegion),
	CatalogNamespace:         readonlyColumn(CatalogIntegrationIcebergRestAttrs.CatalogNamespace),
	RestConfig:               readonlyColumn(CatalogIntegrationIcebergRestAttrs.RestConfig),
	OauthRestAuthentication:  readonlyColumn(CatalogIntegrationIcebergRestAttrs.OauthRestAuthentication),
	BearerRestAuthentication: readonlyColumn(CatalogIntegrationIcebergRestAttrs.BearerRestAuthentication),
	Sigv4RestAuthentication:  readonlyColumn(CatalogIntegrationIcebergRestAttrs.Sigv4RestAuthentication),
}

// CatalogIntegration is the base of the catalog integration variants. All of them are the same Snowflake object,
// CATALOG INTEGRATION, and share one SDK interface, which is built from this object.
var CatalogIntegration = ssot.Object{
	Name:            "CatalogIntegration",
	NamePlural:      "CatalogIntegrations",
	IdType:          "AccountObjectIdentifier",
	VariantBaseOnly: true,
	// The enums of the variants are the SDK enums of the shared interface. The REST authentication type is only in DESCRIBE.
	AdditionalEnums: []*ssot.Enum{
		CatalogIntegrationCatalogSourceType,
		CatalogIntegrationTableFormat,
		CatalogIntegrationRestAuthenticationType,
		CatalogIntegrationAccessDelegationMode,
		CatalogIntegrationCatalogApiType,
	},
	Generators: []ssot.GeneratorOption{
		ssot.WithSchema(
			ssot.SchemaShow(),
			ssot.SchemaDescribe(ssot.SchemaDescribeManualFields(
				&catalogIntegrationDescribeColumns.RestConfig,
				&catalogIntegrationDescribeColumns.OauthRestAuthentication,
				&catalogIntegrationDescribeColumns.BearerRestAuthentication,
				&catalogIntegrationDescribeColumns.Sigv4RestAuthentication,
			)),
		),
		ssot.WithShowOutputAssert(
			ssot.ShowOutputShow(),
			ssot.ShowOutputDescribe(ssot.ShowOutputDescribeSkipFields(
				&catalogIntegrationDescribeColumns.RestConfig,
				&catalogIntegrationDescribeColumns.OauthRestAuthentication,
				&catalogIntegrationDescribeColumns.BearerRestAuthentication,
				&catalogIntegrationDescribeColumns.Sigv4RestAuthentication,
			)),
		),
		// Parameter assertions omitted. The catalog has no CATALOG_INTEGRATION entries.
		ssot.WithObjectAssert(
			ssot.ObjectAssertShow(),
			// Describe returns property rows. The parsed struct is fetched by DescribeDetails.
			ssot.ObjectAssertDescribe(ssot.ObjectAssertDescribeMethod("DescribeDetails")),
		),
		ssot.WithSdk(
			ssot.SdkCreateForVariants(),
			ssot.SdkAlter(),
			ssot.SdkShow(ssot.SdkShowOptionalLike()),
			ssot.SdkDrop(),
			ssot.SdkDescribe(ssot.SdkDescribeDetailsForVariants(), ssot.SdkDescribeDetails(ssot.SdkDescribeDetailsMethod())),
			ssot.SdkDefineFromSsot(),
		),
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&catalogIntegrationShowAttrs.Name,
			&catalogIntegrationShowAttrs.Enabled,
			&catalogIntegrationShowAttrs.Type,
			&catalogIntegrationShowAttrs.Category,
			&catalogIntegrationShowAttrs.Comment,
			&catalogIntegrationShowAttrs.CreatedOn,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CatalogIntegrationIcebergRestAttrs.Id,
			&CatalogIntegrationIcebergRestAttrs.CatalogSource,
			&CatalogIntegrationIcebergRestAttrs.TableFormat,
			&catalogIntegrationDescribeColumns.Enabled,
			&catalogIntegrationDescribeColumns.RefreshIntervalSeconds,
			&catalogIntegrationDescribeColumns.Comment,
			&catalogIntegrationDescribeColumns.GlueAwsRoleArn,
			&catalogIntegrationDescribeColumns.GlueCatalogId,
			&catalogIntegrationDescribeColumns.GlueRegion,
			&catalogIntegrationDescribeColumns.CatalogNamespace,
			&CatalogIntegrationAwsGlueAttrs.GlueAwsIamUserArn,
			&CatalogIntegrationAwsGlueAttrs.GlueAwsExternalId,
			&catalogIntegrationDescribeColumns.RestConfig,
			&catalogIntegrationDescribeColumns.OauthRestAuthentication,
			&catalogIntegrationDescribeColumns.BearerRestAuthentication,
			&catalogIntegrationDescribeColumns.Sigv4RestAuthentication,
		},
	},
	DescribeRows: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&catalogIntegrationDescribeRowAttrs.Property,
			&catalogIntegrationDescribeRowAttrs.PropertyType,
			&catalogIntegrationDescribeRowAttrs.PropertyValue,
			&catalogIntegrationDescribeRowAttrs.PropertyDefault,
		},
	},
	Set: []*ssot.Attribute{
		&catalogIntegrationSetAttrs.OauthRestAuthentication,
		&catalogIntegrationSetAttrs.BearerRestAuthentication,
		&catalogIntegrationSetAttrs.Enabled,
		&catalogIntegrationSetAttrs.RefreshIntervalSeconds,
		&catalogIntegrationSetAttrs.Comment,
	},
	Alter: ssot.AlterSpec{
		Clauses: []ssot.AlterClause{ssot.AlterSet(), ssot.AlterSetTags(), ssot.AlterUnsetTags()},
		SetConflictingFields: [][]*ssot.Attribute{{
			&catalogIntegrationSetAttrs.OauthRestAuthentication,
			&catalogIntegrationSetAttrs.BearerRestAuthentication,
		}},
	},
}
