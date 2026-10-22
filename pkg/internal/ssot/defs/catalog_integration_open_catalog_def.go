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

// CatalogIntegrationOpenCatalogAttrs holds reusable Open Catalog integration fields so SHOW
// and DESCRIBE Columns can point at the same definitions instead of restating names.
//
// rest_config and rest_authentication stay in the schema ext. Those nested blocks have no SSOT struct type.
var CatalogIntegrationOpenCatalogAttrs = struct {
	Enabled                ssot.Attribute
	RefreshIntervalSeconds ssot.Attribute
	Comment                ssot.Attribute
	CatalogNamespace       ssot.Attribute
	Id                     ssot.Attribute
	CatalogSource          ssot.Attribute
	TableFormat            ssot.Attribute
	Name                   ssot.Attribute
	Type                   ssot.Attribute
	Category               ssot.Attribute
	CreatedOn              ssot.Attribute
}{
	Enabled: ssot.Attribute{
		Name:     "enabled",
		Type:     ssot.AttrBool,
		Required: true,
		Description: "Specifies whether the catalog integration is available for use for Iceberg tables. " +
			"`true` allows users to create new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration function normally. " +
			"`false` prevents users from creating new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration cannot access the catalog in the table definition.",
	},
	RefreshIntervalSeconds: ssot.Attribute{
		Name: "refresh_interval_seconds",
		Type: ssot.AttrInt,
		Description: "Specifies the number of seconds to wait between attempts to poll the external Iceberg catalog for metadata updates for automated refresh. " +
			"For Delta-based tables, specifies the number of seconds to wait between attempts to poll your external cloud storage for new metadata.",
		Min:           new(1),
		CreateMapping: ssot.IntAttributeCreateBuilder,
	},
	Comment: ssot.Attribute{
		Name:          "comment",
		Type:          ssot.AttrString,
		Description:   "Specifies a comment for the catalog integration.",
		OptionalInDb:  true,
		CreateMapping: ssot.StringAttributeCreateBuilder,
	},
	CatalogNamespace: ssot.Attribute{
		Name:        "catalog_namespace",
		Type:        ssot.AttrString,
		ForceNew:    true,
		Description: "Specifies the default Open Catalog namespace for all Iceberg tables that you associate with the catalog integration.",
	},
	Id: ssot.Attribute{Name: "id", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier", Readonly: true},
	CatalogSource: ssot.Attribute{
		Name:        "catalog_source",
		Type:        ssot.AttrEnum,
		Description: "Specifies the type of catalog source. This field is used to detect when the catalog source was changed outside of Terraform and to recreate the resource when that happens.",
		Enum:        CatalogIntegrationCatalogSourceType,
		Readonly:    true,
	},
	TableFormat: ssot.Attribute{Name: "table_format", Type: ssot.AttrEnum, Enum: CatalogIntegrationTableFormat, Readonly: true},
	Name:        ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	Type:        ssot.Attribute{Name: "type", Type: ssot.AttrString, Readonly: true},
	Category:    ssot.Attribute{Name: "category", Type: ssot.AttrString, Readonly: true},
	CreatedOn:   ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
}

var CatalogIntegrationOpenCatalog = ssot.Object{
	Name:                   "CatalogIntegrationOpenCatalog",
	NamePlural:             "CatalogIntegrations",
	DescribeFailRead:       true,
	IdType:                 "AccountObjectIdentifier",
	IdentityNameForceNew:   true,
	AttributesModification: true,
	Generators: []ssot.GeneratorOption{
		// ShowOutputShow omitted. SHOW rows are sdk.CatalogIntegration, so tests use CatalogIntegrationShowOutput.
		// WithSchema omitted. SHOW is the shared CatalogIntegration schema alias in the ext.
		// DESCRIBE stays on the SDK schema list until ManualFields exist in SSOT.
		// WithModel omitted. rest_config and rest_authentication are nested blocks, which the model generator cannot emit.
		// WithResourceAssert omitted. catalog_source is added in the schema ext, so the schema-based assert still emits it.
		// ObjectAssertShow omitted. There is no sdk.CatalogIntegrationOpenCatalog type.
		// ObjectAssertDescribe omitted. The SDK list keeps DescribeOpenCatalog and the nested rest asserts.
		ssot.WithShowOutputAssert(ssot.ShowOutputDescribe()),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		// Snowflake's command is DESCRIBE CATALOG INTEGRATION. The generated text would use the resource name.
		ssot.WithDescribeOutputSchema(ssot.ExtCall()),
		ssot.WithHandleExternalChangesInDescribe(ssot.ExtCall()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithSetStateToValuesFromConfig(ssot.Omitted()),
		// SHOW has name, enabled, and comment only. catalog_source and catalog_namespace
		// come from DESCRIBE, so SetConfigFields stays on the describe hook.
		ssot.WithSetConfigFields(ssot.Omitted()),
	},
	Attributes: []ssot.Attribute{
		CatalogIntegrationOpenCatalogAttrs.Enabled,
		CatalogIntegrationOpenCatalogAttrs.RefreshIntervalSeconds,
		CatalogIntegrationOpenCatalogAttrs.Comment,
		CatalogIntegrationOpenCatalogAttrs.CatalogNamespace,
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CatalogIntegrationOpenCatalogAttrs.Name,
			&CatalogIntegrationOpenCatalogAttrs.Enabled,
			&CatalogIntegrationOpenCatalogAttrs.Type,
			&CatalogIntegrationOpenCatalogAttrs.Category,
			&CatalogIntegrationOpenCatalogAttrs.Comment,
			&CatalogIntegrationOpenCatalogAttrs.CreatedOn,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CatalogIntegrationOpenCatalogAttrs.Id,
			&CatalogIntegrationOpenCatalogAttrs.CatalogSource,
			&CatalogIntegrationOpenCatalogAttrs.TableFormat,
			&CatalogIntegrationOpenCatalogAttrs.Enabled,
			&CatalogIntegrationOpenCatalogAttrs.RefreshIntervalSeconds,
			&CatalogIntegrationOpenCatalogAttrs.Comment,
			&CatalogIntegrationOpenCatalogAttrs.CatalogNamespace,
		},
	},
}
