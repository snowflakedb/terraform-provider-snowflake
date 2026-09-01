package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var rowAccessPoliciesDef = g.NewInterface(
	"RowAccessPolicies",
	"RowAccessPolicy",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-row-access-policy",
		g.NewQueryStruct("CreateRowAccessPolicy").
			Create().
			OrReplace().
			SQL("ROW ACCESS POLICY").
			IfNotExists().
			Name().
			SQL("AS").
			ListQueryStructField(
				"args",
				g.NewQueryStruct("CreateRowAccessPolicyArgs").
					Text("Name", g.KeywordOptions().DoubleQuotes().Required()).
					PredefinedQueryStructField("DataType", "datatypes.DataType", g.ParameterOptions().NoEquals().Required()),
				g.ParameterOptions().Parentheses().NoEquals().Required(),
			).
			SQL("RETURNS BOOLEAN").
			BodyWithPrecedingArrow().
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidateValueSet, "args").
			WithValidation(g.ValidateValueSet, "body").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-row-access-policy",
		g.NewQueryStruct("AlterRowAccessPolicy").
			Alter().
			SQL("ROW ACCESS POLICY").
			Name().
			RenameTo().
			OptionalSetBodyWithPrecedingArrow().
			OptionalSetTags().
			OptionalUnsetTags().
			OptionalTextAssignment("SET COMMENT", g.ParameterOptions().SingleQuotes()).
			OptionalSQL("UNSET COMMENT").
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "SetBody", "SetTags", "UnsetTags", "SetComment", "UnsetComment"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-row-access-policy",
		g.NewQueryStruct("DropRowAccessPolicy").
			Drop().
			SQL("ROW ACCESS POLICY").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-row-access-policies",
		g.StructPair("rowAccessPolicyDBRow", "RowAccessPolicy").
			Text("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("kind").
			Text("owner").
			OptionalText("comment", g.WithRequiredInPlain()).
			Text("options").
			Text("owner_role_type"),
		g.NewQueryStruct("ShowRowAccessPolicies").
			Show().
			SQL("ROW ACCESS POLICIES").
			OptionalLike().
			OptionalExtendedIn().
			OptionalLimitFrom(),
		g.ShowByIDExtendedInFiltering,
		g.ShowByIDLikeFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-row-access-policy",
		g.StructPair("describeRowAccessPolicyDBRow", "RowAccessPolicyDescription").
			Text("name").
			Field("signature", "string", "[]TableColumnSignature", g.WithCustomParser("ParseTableColumnSignature")).
			Text("return_type").
			Text("body"),
		g.NewQueryStruct("DescribeRowAccessPolicy").
			Describe().
			SQL("ROW ACCESS POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
