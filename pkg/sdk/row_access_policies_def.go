package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

var rowAccessPolicyDbRow = g.DbStruct("rowAccessPolicyDBRow").
	Text("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	Text("kind").
	Text("owner").
	OptionalText("comment").
	Text("options").
	Text("owner_role_type")

var rowAccessPolicy = g.PlainStruct("RowAccessPolicy").
	Text("CreatedOn").
	Text("Name").
	Text("DatabaseName").
	Text("SchemaName").
	Text("Kind").
	Text("Owner").
	Text("Comment").
	Text("Options").
	Text("OwnerRoleType")

var RowAccessPoliciesDef = g.NewInterface(
	"RowAccessPolicies",
	"RowAccessPolicy",
	g.KindOfT[SchemaObjectIdentifier](),
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
			OptionalIdentifier("RenameTo", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
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
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-row-access-policies",
		rowAccessPolicyDbRow,
		rowAccessPolicy,
		g.NewQueryStruct("ShowRowAccessPolicies").
			Show().
			SQL("ROW ACCESS POLICIES").
			OptionalLike().
			OptionalExtendedIn().
			OptionalLimitFrom(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDExtendedInFiltering,
		g.ShowByIDLikeFiltering,
	).
	DescribeOperation(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-row-access-policy",
		g.DbStruct("describeRowAccessPolicyDBRow").
			Field("name", "string").
			Field("signature", "string").
			Field("return_type", "string").
			Field("body", "string"),
		g.PlainStruct("RowAccessPolicyDescription").
			Field("Name", "string").
			Field("Signature", "[]TableColumnSignature").
			Field("ReturnType", "string").
			Field("Body", "string"),
		g.NewQueryStruct("DescribeRowAccessPolicy").
			Describe().
			SQL("ROW ACCESS POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
