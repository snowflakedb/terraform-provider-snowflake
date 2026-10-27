package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// CatalogIntegrationObjectStorageAttrs holds reusable object storage catalog integration
// fields so SHOW and DESCRIBE Columns can point at the same definitions.
var CatalogIntegrationObjectStorageAttrs = struct {
	Enabled                ssot.Attribute
	RefreshIntervalSeconds ssot.Attribute
	Comment                ssot.Attribute
	TableFormat            ssot.Attribute
	Name                   ssot.Attribute
	Type                   ssot.Attribute
	Category               ssot.Attribute
	CreatedOn              ssot.Attribute
	Id                     ssot.Attribute
	CatalogSource          ssot.Attribute
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
	TableFormat: ssot.Attribute{
		Name:        "table_format",
		Type:        ssot.AttrEnum,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the table format.",
		Enum:        CatalogIntegrationTableFormat,
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
		Enum:        CatalogIntegrationCatalogSourceType,
		Readonly:    true,
	},
}

var catalogIntegrationObjectStorageCatalogSourceStatic = ssot.Attribute{
	Name:      "catalog_source_object_store",
	SqlOnly:   true,
	CreateSql: ssot.StaticSQL("catalogSource", "CATALOG_SOURCE = OBJECT_STORE"),
}

var CatalogIntegrationObjectStorage = ssot.Object{
	Name:                   "CatalogIntegrationObjectStorage",
	NamePlural:             "CatalogIntegrations",
	VariantOf:              &CatalogIntegration,
	DescribeFailRead:       true,
	IdType:                 "AccountObjectIdentifier",
	IdentityNameForceNew:   true,
	AttributesModification: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		// ShowOutputShow omitted. SHOW rows are sdk.CatalogIntegration, so tests use CatalogIntegrationShowOutput.
		// WithSchema omitted. SHOW is the shared CatalogIntegration schema alias in the ext.
		// DESCRIBE stays on the SDK schema list. SSOT would register the same sdk.CatalogIntegrationObjectStorageDetails name.
		// WithModel omitted. catalog_source is added in the schema ext, so the schema-based model still emits it.
		// WithResourceAssert omitted. catalog_source is added in the schema ext, so the schema-based assert still emits it.
		// ObjectAssertShow omitted. There is no sdk.CatalogIntegrationObjectStorage type.
		// ObjectAssertDescribe omitted. The SDK list keeps DescribeObjectStorage on the catalog integration helper.
		// DESCRIBE CATALOG INTEGRATION is shared by every catalog integration variant.
		ssot.WithSdk(ssot.SdkDescribe(ssot.SdkDescribeSql("CATALOG INTEGRATION"), ssot.SdkDescribeDetails())),
		ssot.WithShowOutputAssert(ssot.ShowOutputDescribe()),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithHandleExternalChangesInDescribe(ssot.ExtCall()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
	},
	Attributes: []ssot.Attribute{
		CatalogIntegrationObjectStorageAttrs.Enabled,
		CatalogIntegrationObjectStorageAttrs.RefreshIntervalSeconds,
		CatalogIntegrationObjectStorageAttrs.Comment,
		CatalogIntegrationObjectStorageAttrs.TableFormat,
	},
	Create: ssot.CreateSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-object-storage",
		Clauses: []*ssot.Attribute{
			&catalogIntegrationObjectStorageCatalogSourceStatic,
			&CatalogIntegrationObjectStorageAttrs.TableFormat,
			&CatalogIntegrationObjectStorageAttrs.Enabled,
			&CatalogIntegrationObjectStorageAttrs.RefreshIntervalSeconds,
			&CatalogIntegrationObjectStorageAttrs.Comment,
		},
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CatalogIntegrationObjectStorageAttrs.Name,
			&CatalogIntegrationObjectStorageAttrs.Enabled,
			&CatalogIntegrationObjectStorageAttrs.Type,
			&CatalogIntegrationObjectStorageAttrs.Category,
			&CatalogIntegrationObjectStorageAttrs.Comment,
			&CatalogIntegrationObjectStorageAttrs.CreatedOn,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CatalogIntegrationObjectStorageAttrs.Id,
			&CatalogIntegrationObjectStorageAttrs.CatalogSource,
			&CatalogIntegrationObjectStorageAttrs.TableFormat,
			&CatalogIntegrationObjectStorageAttrs.Enabled,
			&CatalogIntegrationObjectStorageAttrs.RefreshIntervalSeconds,
			&CatalogIntegrationObjectStorageAttrs.Comment,
		},
	},
}
