package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var alertActionDef = g.NewEnum("AlertAction", "AlertActions", "RESUME", "SUSPEND")
var alertStateDef = g.NewEnum("AlertState", "AlertStates", "started", "suspended")

func alertSet() *g.QueryStruct {
	return g.NewQueryStruct("AlertSet").
		OptionalIdentifier("Warehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("WAREHOUSE")).
		OptionalTextAssignment("SCHEDULE", g.ParameterOptions().SingleQuotes()).
		OptionalComment()
}

func alertUnset() *g.QueryStruct {
	return g.NewQueryStruct("AlertUnset").
		OptionalSQL("WAREHOUSE").
		OptionalSQL("SCHEDULE").
		OptionalSQL("COMMENT")
}

var alertPairs = g.StructPair("alertDBRow", "Alert").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	Text("owner").
	Field("comment", "*string", "*string").
	Text("warehouse").
	Text("schedule").
	Enum("state", alertStateDef).
	Text("condition").
	Text("action").
	OptionalText("owner_role_type", g.WithRequiredInPlain())

var alertDetailPairs = g.StructPair("alertDetailRow", "AlertDetails").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	Text("owner").
	Field("comment", "*string", "*string").
	Text("warehouse").
	Text("schedule").
	Text("state").
	Text("condition").
	Text("action")

var alertsDef = g.NewInterface(
	"Alerts",
	"Alert",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-alert",
	g.NewQueryStruct("CreateAlert").
		Create().
		OrReplace().
		SQL("ALERT").
		IfNotExists().
		Name().
		Identifier("Warehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("WAREHOUSE").Required()).
		TextAssignment("SCHEDULE", g.ParameterOptions().SingleQuotes()).
		OptionalComment().
		WithField(g.NewField("condition", "[]AlertCondition", g.Tags().Keyword().Parentheses().SQL("IF"), g.KeywordOptions().NoComma().Required())).
		PredefinedQueryStructField("action", "string", g.ParameterOptions().NoEquals().SQL("THEN").Required()).
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-alert",
	g.NewQueryStruct("AlterAlert").
		Alter().
		SQL("ALERT").
		IfExists().
		Name().
		OptionalEnum("Action", alertActionDef, g.KeywordOptions()).
		OptionalQueryStructField("Set", alertSet(), g.KeywordOptions().SQL("SET")).
		OptionalQueryStructField("Unset", alertUnset(), g.KeywordOptions().SQL("UNSET")).
		WithField(g.NewField("ModifyCondition", "*[]string", g.Tags().Keyword().Parentheses().SQL("MODIFY CONDITION EXISTS"), g.KeywordOptions().NoComma())).
		PredefinedQueryStructField("ModifyAction", "*string", g.ParameterOptions().NoEquals().SQL("MODIFY ACTION")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Action", "Set", "Unset", "ModifyCondition", "ModifyAction"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-alert",
	g.NewQueryStruct("DropAlert").
		Drop().
		SQL("ALERT").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-alerts",
	alertPairs,
	g.NewQueryStruct("ShowAlert").
		Show().
		Terse().
		SQL("ALERTS").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		OptionalNumberAssignment("LIMIT", g.ParameterOptions().NoEquals()),
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-alert",
	alertDetailPairs,
	g.NewQueryStruct("DescribeAlert").
		Describe().
		SQL("ALERT").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).WithEnums(alertActionDef, alertStateDef)
