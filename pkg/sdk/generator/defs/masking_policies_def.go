package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var createMaskingPolicySignatureDef = g.NewQueryStruct("CreateMaskingPolicySignature").
	Text("Name", g.KeywordOptions().DoubleQuotes().Required()).
	PredefinedQueryStructField("DataType", "datatypes.DataType", g.ParameterOptions().NoEquals().Required())

var maskingPoliciesDef = g.NewInterface(
	"MaskingPolicies",
	"MaskingPolicy",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-masking-policy",
		g.NewQueryStruct("CreateMaskingPolicy").
			Create().
			OrReplace().
			SQL("MASKING POLICY").
			IfNotExists().
			Name().
			SQL("AS").
			ListQueryStructField(
				"signature",
				createMaskingPolicySignatureDef,
				g.ParameterOptions().Parentheses().NoEquals().Required(),
			).
			PredefinedQueryStructField("returns", "datatypes.DataType", g.ParameterOptions().NoEquals().SQL("RETURNS").Required()).
			BodyWithPrecedingArrow().
			OptionalComment().
			OptionalBooleanAssignment("EXEMPT_OTHER_POLICIES", nil).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidateValueSet, "signature").
			WithValidation(g.ValidateValueSet, "returns").
			WithValidation(g.ValidateValueSet, "body").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-masking-policy",
		g.NewQueryStruct("AlterMaskingPolicy").
			Alter().
			SQL("MASKING POLICY").
			IfExists().
			Name().
			OptionalIdentifier("NewName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("MaskingPolicySet").
					PredefinedQueryStructField("Body", "*string", g.ParameterOptions().NoEquals().NoQuotes().SQL("BODY ->")).
					OptionalComment().
					WithValidation(g.ExactlyOneValueSet, "Body", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("MaskingPolicyUnset").
					OptionalSQL("COMMENT").
					WithValidation(g.ExactlyOneValueSet, "Comment"),
				g.KeywordOptions().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "NewName", "Set", "Unset", "SetTags", "UnsetTags").
			WithAdditionalValidations(),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-masking-policy",
		g.NewQueryStruct("DropMaskingPolicy").
			Drop().
			SQL("MASKING POLICY").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-masking-policies",
		g.StructPair("maskingPolicyDBRow", "MaskingPolicy").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("kind").
			Text("owner").
			OptionalText("comment", g.WithRequiredInPlain()).
			JsonField("options", "MaskingPolicyOptions").
			Text("owner_role_type").
			PlainOnlyField("ExemptOtherPolicies", "bool"),
		g.NewQueryStruct("ShowMaskingPolicies").
			Show().
			SQL("MASKING POLICIES").
			OptionalLike().
			OptionalExtendedIn().
			OptionalLimitFrom(),
		g.ShowByIDExtendedInFiltering,
		g.ShowByIDLikeFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-masking-policy",
		g.StructPair("describeMaskingPolicyDBRow", "MaskingPolicyDetails").
			Text("name").
			Field("signature", "string", "[]TableColumnSignature", g.WithCustomParser("ParseTableColumnSignature")).
			DataType("return_type").
			Text("body"),
		g.NewQueryStruct("DescribeMaskingPolicy").
			Describe().
			SQL("MASKING POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
