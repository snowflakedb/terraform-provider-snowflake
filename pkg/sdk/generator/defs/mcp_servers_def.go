package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var mcpServersDef = g.NewInterface(
	"McpServers",
	"McpServer",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-mcp-server",
		g.NewQueryStruct("CreateMcpServer").
			Create().
			OrReplace().
			SQL("MCP SERVER").
			IfNotExists().
			Name().
			TextAssignment("FROM SPECIFICATION", g.ParameterOptions().NoEquals().DoubleDollarQuotes()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
			WithValidation(g.NoDoubleDollarQuotes, "FromSpecification"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-mcp-server",
		g.NewQueryStruct("DropMcpServer").
			Drop().
			SQL("MCP SERVER").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-mcp-servers",
		g.StructPair("mcpServerDBRow", "McpServer").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			OptionalText("comment", g.WithRequiredInPlain()),
		g.NewQueryStruct("ShowMcpServers").
			Show().
			SQL("MCP SERVERS").
			OptionalLike().
			OptionalIn(),
		g.ShowByIDLikeFiltering,
		g.ShowByIDInFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-mcp-server",
		g.StructPair("mcpServerDetailsRow", "McpServerDetails").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			OptionalText("comment", g.WithRequiredInPlain()).
			Text("server_spec", g.WithCustomParser("NormalizeMcpServerSpecification")).
			Time("created_on"),
		g.NewQueryStruct("DescribeMcpServer").
			Describe().
			SQL("MCP SERVER").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
