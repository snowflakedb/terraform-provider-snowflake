package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// openflowConnectorDefinitionsDef exposes only SHOW. Connector definitions are global, immutable, and
// managed entirely by Snowflake, so there is no CREATE/ALTER/DROP/DESCRIBE for them. They are what
// CREATE OPENFLOW CONNECTOR ... FROM DEFINITION refers to.
//
// A definition is keyed by name *and* version, so ShowByID - which could only filter by name - is
// suppressed rather than silently returning an arbitrary version.
var openflowConnectorDefinitionsDef = g.NewInterface(
	"OpenflowConnectorDefinitions",
	"OpenflowConnectorDefinition",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-connector#show-openflow-connector-definitions",
	g.StructPair("openflowConnectorDefinitionRow", "OpenflowConnectorDefinition").
		Text("name").
		Text("provider").
		Text("version").
		Text("description").
		Text("display_name").
		// categories arrives as a JSON array (["Databases","Analytics"]). The shared
		// ParseCommaSeparatedStringArray handles that shape, including commas inside quoted elements, so the
		// custom parser here only adapts it to the signature the generator wants and maps the three spellings
		// of "no values" onto an empty slice.
		//
		// Nullable: scanning a SQL NULL into a plain string fails the entire SHOW, not just this column.
		Field("categories", "sql.NullString", "[]string", g.WithCustomParser("ParseOpenflowConnectorDefinitionCategories")).
		OptionalText("min_runtime_node_type").
		OptionalNumber("max_node_count"),
	g.NewQueryStruct("ShowOpenflowConnectorDefinitions").
		Show().
		SQL("OPENFLOW CONNECTOR DEFINITIONS").
		OptionalLike().
		OptionalStartsWith().
		OptionalLimitFrom(),
	g.ShowByIDSuppressed,
).WithEnabledGenerationParts(g.PartUnitTests)
