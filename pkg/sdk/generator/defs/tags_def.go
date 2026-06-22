package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var TagPropagationEnumDef = g.NewEnum(
	"TagPropagation", "TagPropagationValues",
	"NONE", "ON_DEPENDENCY", "ON_DATA_MOVEMENT", "ON_DEPENDENCY_AND_DATA_MOVEMENT",
)

func tagMaskingPolicy() *g.QueryStruct {
	return g.NewQueryStruct("TagMaskingPolicy").
		Identifier("Name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required())
}

func tagSetMaskingPolicies() *g.QueryStruct {
	return g.NewQueryStruct("TagSetMaskingPolicies").
		ListQueryStructField("MaskingPolicies", tagMaskingPolicy(), g.ListOptions()).
		OptionalSQL("FORCE")
}

func tagUnsetMaskingPolicies() *g.QueryStruct {
	return g.NewQueryStruct("TagUnsetMaskingPolicies").
		ListQueryStructField("MaskingPolicies", tagMaskingPolicy(), g.ListOptions())
}

func allowedValues() *g.QueryStruct {
	return g.NewQueryStruct("AllowedValues").
		List("Values", "StringAllowEmpty", g.ListOptions()).
		WithAdditionalValidations()
}

func tagOnConflict() *g.QueryStruct {
	return g.NewQueryStruct("TagOnConflict").
		AssignmentWithFieldName("ON_CONFLICT", "*string", g.ParameterOptions().SingleQuotes(), "CustomValue").
		OptionalSQLWithCustomFieldName("AllowedValuesSequence", "ON_CONFLICT = ALLOWED_VALUES_SEQUENCE").
		WithValidation(g.ExactlyOneValueSet, "CustomValue", "AllowedValuesSequence")
}

func tagPropagate() *g.QueryStruct {
	return g.NewQueryStruct("TagPropagate").
		OptionalAssignmentWithFieldName("PROPAGATE", "TagPropagation", g.ParameterOptions(), "PropagationMethod").
		OptionalQueryStructField("OnConflict", tagOnConflict(), g.KeywordOptions())
}

func tagAdd() *g.QueryStruct {
	return g.NewQueryStruct("TagAdd").
		OptionalQueryStructField("AllowedValues", allowedValues(), g.KeywordOptions().SQL("ALLOWED_VALUES"))
}

func tagDrop() *g.QueryStruct {
	return g.NewQueryStruct("TagDrop").
		OptionalQueryStructField("AllowedValues", allowedValues(), g.KeywordOptions().SQL("ALLOWED_VALUES"))
}

func tagRename() *g.QueryStruct {
	return g.NewQueryStruct("TagRename").
		Identifier("Name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		WithValidation(g.ValidIdentifier, "Name")
}

func tagSet() *g.QueryStruct {
	return g.NewQueryStruct("TagSet").
		OptionalQueryStructField("MaskingPolicies", tagSetMaskingPolicies(), g.KeywordOptions()).
		OptionalQueryStructField("AllowedValues", allowedValues(), g.KeywordOptions().SQL("ALLOWED_VALUES")).
		OptionalQueryStructField("Propagate", tagPropagate(), g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.AtLeastOneValueSet, "MaskingPolicies", "AllowedValues", "Propagate", "Comment").
		WithAdditionalValidations()
}

func tagUnset() *g.QueryStruct {
	return g.NewQueryStruct("TagUnset").
		OptionalQueryStructField("MaskingPolicies", tagUnsetMaskingPolicies(), g.KeywordOptions()).
		OptionalSQL("ALLOWED_VALUES").
		OptionalSQL("PROPAGATE").
		OptionalSQL("ON_CONFLICT").
		OptionalSQL("COMMENT").
		WithValidation(g.ExactlyOneValueSet, "MaskingPolicies", "AllowedValues", "Propagate", "OnConflict", "Comment").
		WithAdditionalValidations()
}

var tagsDef = g.NewInterface(
	"Tags",
	"Tag",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-tag",
	g.NewQueryStruct("CreateTag").
		Create().OrReplace().SQL("TAG").IfNotExists().Name().
		OptionalQueryStructField("AllowedValues", allowedValues(), g.KeywordOptions().SQL("ALLOWED_VALUES")).
		OptionalQueryStructField("Propagate", tagPropagate(), g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-tag",
	g.NewQueryStruct("AlterTag").
		Alter().SQL("TAG").IfExists().Name().
		OptionalQueryStructField("Add", tagAdd(), g.KeywordOptions().SQL("ADD")).
		OptionalQueryStructField("Drop", tagDrop(), g.KeywordOptions().SQL("DROP")).
		OptionalQueryStructField("Set", tagSet(), g.KeywordOptions().SQL("SET")).
		OptionalQueryStructField("Unset", tagUnset(), g.KeywordOptions().SQL("UNSET")).
		OptionalQueryStructField("Rename", tagRename(), g.KeywordOptions().SQL("RENAME TO")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Add", "Drop", "Set", "Unset", "Rename"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-tags",
	g.StructPair("tagRow", "Tag").
		Time("created_on").
		Text("name").
		Text("database_name").
		Text("schema_name").
		Text("owner").
		Text("comment").
		OptionalPlainField("allowed_values", "[]string").
		Text("owner_role_type").
		OptionalEnum("propagate", TagPropagationEnumDef).
		OptionalText("on_conflict"),
	g.NewQueryStruct("ShowTags").
		Show().SQL("TAGS").OptionalLike().OptionalExtendedIn().
		WithAdditionalValidations(),
	g.ShowByIDLikeFiltering,
	g.ShowByIDExtendedInFiltering,
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-tag",
	g.NewQueryStruct("DropTag").
		Drop().SQL("TAG").IfExists().Name().
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"Undrop",
	"https://docs.snowflake.com/en/sql-reference/sql/undrop-tag",
	g.NewQueryStruct("UndropTag").
		SQL("UNDROP").SQL("TAG").Name().
		WithValidation(g.ValidIdentifier, "name"),
).WithCustomInterfaceMethod(
	"Set", "",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*SetTagRequest")},
	"error",
).WithCustomInterfaceMethod(
	"Unset", "",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*UnsetTagRequest")},
	"error",
).WithCustomInterfaceMethod(
	"UnsetSafely", "",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*UnsetTagRequest")},
	"error",
).WithCustomInterfaceMethod(
	"SetOnCurrentAccount", "",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*SetTagOnCurrentAccountRequest")},
	"error",
).WithCustomInterfaceMethod(
	"UnsetOnCurrentAccount", "",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*UnsetTagOnCurrentAccountRequest")},
	"error",
).WithEnums(TagPropagationEnumDef)
