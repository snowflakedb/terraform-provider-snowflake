package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var CloudTypeEnumDef = g.NewEnum("CloudType", "CloudTypes", "aws", "azure", "gcp")

var replicationAccountPairs = g.StructPair("replicationAccountRow", "ReplicationAccount").
	Text("snowflake_region").
	Time("created_on").
	Text("account_name").
	Text("account_locator").
	OptionalText("comment", g.WithRequiredInPlain()).
	Text("organization_name").
	Bool("is_org_admin")

var replicationDatabasePairs = g.StructPair("replicationDatabaseRow", "ReplicationDatabase").
	OptionalText("region_group", g.WithRequiredInPlain()).
	Text("snowflake_region").
	Time("created_on").
	Text("account_name").
	Text("name").
	OptionalText("comment", g.WithRequiredInPlain()).
	Bool("is_primary").
	OptionalExternalObjectIdentifier("primary", g.WithPlainFieldName("PrimaryDatabase")).
	OptionalText("replication_allowed_to_accounts", g.WithRequiredInPlain()).
	OptionalText("failover_allowed_to_accounts", g.WithRequiredInPlain()).
	Text("organization_name").
	Text("account_locator")

var regionPairs = g.StructPair("regionRow", "Region").
	OptionalText("region_group", g.WithRequiredInPlain()).
	Text("snowflake_region").
	Enum("cloud", CloudTypeEnumDef, g.WithPlainFieldName("CloudType")).
	Text("region").
	Text("display_name")

// NameSingular is empty on purpose. A non-empty singular ("ReplicationFunction") would produce
// ShowReplicationDatabasesReplicationFunctionOptions etc.
// Identifier kind is a required placeholder — replication functions are not a named Snowflake object.
var replicationFunctionsDef = g.NewInterface(
	"ReplicationFunctions",
	"",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CustomShowOperationWithPairedStructs(
	"ShowReplicationAccounts",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/show-replication-accounts",
	replicationAccountPairs,
	g.NewQueryStruct("ShowReplicationAccounts").
		Show().
		SQL("REPLICATION ACCOUNTS"),
).CustomShowOperationWithPairedStructs(
	"ShowReplicationDatabases",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/show-replication-databases",
	replicationDatabasePairs,
	g.NewQueryStruct("ShowReplicationDatabases").
		Show().
		SQL("REPLICATION DATABASES").
		OptionalLike().
		OptionalIdentifier("WithPrimary", g.KindOfT[sdkcommons.ExternalObjectIdentifier](), g.IdentifierOptions().SQL("WITH PRIMARY")).
		WithValidation(g.ValidIdentifierIfSet, "WithPrimary"),
).CustomShowOperationWithPairedStructs(
	"ShowRegions",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/show-regions",
	regionPairs,
	g.NewQueryStruct("ShowRegions").
		Show().
		SQL("REGIONS").
		OptionalLike(),
).WithEnums(CloudTypeEnumDef).
	// TODO(next-pr): drop this allow-list once SDK unit tests move to generated + *_ext_test.go
	WithAllowedGenerationParts(
		g.PartDefault,
		g.PartDto,
		g.PartDtoBuilders,
		g.PartImpl,
		g.PartValidations,
		g.PartEnums,
	)
