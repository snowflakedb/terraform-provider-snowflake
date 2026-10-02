package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var gatewaysDef = g.NewInterface(
	"Gateways",
	"Gateway",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-gateway",
		g.NewQueryStruct("CreateGateway").
			Create().
			OrReplace().
			SQL("GATEWAY").
			IfNotExists().
			Name().
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			TextAssignment("FROM SPECIFICATION", g.ParameterOptions().NoEquals().DoubleDollarQuotes()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
			WithValidation(g.NoDoubleDollarQuotes, "FromSpecification"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-gateway",
		g.NewQueryStruct("AlterGateway").
			Alter().
			SQL("GATEWAY").
			IfExists().
			Name().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("GatewaySet").
					OptionalAssignment("COMMENT", "StringAllowEmpty", g.ParameterOptions()).
					WithValidation(g.AtLeastOneValueSet, "Comment"),
				g.ListOptions().NoParentheses().SQL("SET"),
			).
			OptionalQueryStructField(
				"ModifyLiveVersionSet",
				g.NewQueryStruct("GatewayModifyLiveVersionSet").
					TextAssignment("SPECIFICATION", g.ParameterOptions().DoubleDollarQuotes()).
					WithValidation(g.NoDoubleDollarQuotes, "Specification"),
				g.KeywordOptions().SQL("MODIFY LIVE VERSION SET"),
			).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Set", "ModifyLiveVersionSet"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-gateway",
		g.NewQueryStruct("DropGateway").
			Drop().
			SQL("GATEWAY").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-gateways",
		g.StructPair("showGatewayDBRow", "Gateway").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			OptionalText("comment", g.WithRequiredInPlain()),
		g.NewQueryStruct("ShowGateways").
			Show().
			SQL("GATEWAYS").
			OptionalLike().
			OptionalExtendedIn().
			OptionalStartsWith().
			OptionalLimit(),
		g.ShowByIDLikeFiltering,
		g.ShowByIDExtendedInFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-gateway",
		g.StructPair("gatewayDetailsRow", "GatewayDetails").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			OptionalText("comment", g.WithRequiredInPlain()).
			Text("gateway_spec", g.WithCustomParser("NormalizeGatewaySpecification")).
			Time("created_on"),
		g.NewQueryStruct("DescribeGateway").
			Describe().
			SQL("GATEWAY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	WithShowObjectType("Gateway")
