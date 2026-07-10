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

func setTagQueryStruct() *g.QueryStruct {
	return g.NewQueryStruct("SetTag").
		SQL("ALTER").
		PredefinedQueryStructField("objectType", "ObjectType", g.KeywordOptions().Required()).
		PredefinedQueryStructField("objectName", "ObjectIdentifier", g.IdentifierOptions().Required()).
		AssignmentWithFieldName("MODIFY COLUMN", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Column").
		OptionalSetTags().
		WithValidation(g.ValidIdentifier, "objectName").
		WithAdditionalValidations()
}

func unsetTagQueryStruct() *g.QueryStruct {
	return g.NewQueryStruct("UnsetTag").
		SQL("ALTER").
		PredefinedQueryStructField("objectType", "ObjectType", g.KeywordOptions().Required()).
		IfExists().
		PredefinedQueryStructField("objectName", "ObjectIdentifier", g.IdentifierOptions().Required()).
		AssignmentWithFieldName("MODIFY COLUMN", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Column").
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "objectName").
		WithAdditionalValidations()
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
		RenameTo().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ExactlyOneValueSet, "Add", "Drop", "Set", "Unset", "RenameTo"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-tags",
	g.StructPair("tagRow", "Tag").
		Time("created_on").
		Text("name").
		Text("database_name").
		Text("schema_name").
		Text("owner").
		Text("comment").
		OptionalPlainField("allowed_values", "[]string", g.WithManualConvert()).
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
).CustomOperationWithOpts(
	"Set",
	"https://docs.snowflake.com/en/sql-reference/sql/alter-tag",
	setTagQueryStruct(),
	[]g.CustomOperationOption{g.WithRequestAdjust()},
).CustomOperationWithOpts(
	"Unset",
	"https://docs.snowflake.com/en/sql-reference/sql/alter-tag",
	unsetTagQueryStruct(),
	[]g.CustomOperationOption{g.WithRequestAdjust()},
).WithCustomInterfaceMethod(
	"UnsetSafely", "",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*UnsetTagRequest")},
	"error",
).WithCustomInterfaceMethod(
	// TODO [next PRs]: change signature to use a proper Accounts request type when Accounts is migrated to the SDK generator.
	"SetOnCurrentAccount", "",
	[]*g.MethodParameter{g.NewMethodParameter("setTags", "[]TagAssociation")},
	"error",
).WithCustomInterfaceMethod(
	// TODO [next PRs]: change signature to use a proper Accounts request type when Accounts is migrated to the SDK generator.
	"UnsetOnCurrentAccount", "",
	[]*g.MethodParameter{g.NewMethodParameter("unsetTags", "[]ObjectIdentifier")},
	"error",
).WithEnums(TagPropagationEnumDef)
