package sdk

import (
	"fmt"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

//go:generate go run ./poc/main.go

// DBT Project default version options
type DbtProjectDefaultVersion string

const (
	DbtProjectDefaultVersionFirst DbtProjectDefaultVersion = "FIRST"
	DbtProjectDefaultVersionLast  DbtProjectDefaultVersion = "LAST"
)

func ToDbtProjectDefaultVersion(s string) (DbtProjectDefaultVersion, error) {
	switch version := DbtProjectDefaultVersion(strings.ToUpper(s)); version {
	case DbtProjectDefaultVersionFirst, DbtProjectDefaultVersionLast:
		return version, nil
	default:
		// Handle VERSION$<num> format
		if strings.HasPrefix(strings.ToUpper(s), "VERSION$") {
			return DbtProjectDefaultVersion(strings.ToUpper(s)), nil
		}
		return "", fmt.Errorf("unknown dbt project default version: %s", s)
	}
}

var dbtProjectDbRow = g.DbStruct("dbtProjectDBRow").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	OptionalText("source_location").
	OptionalText("default_args").
	OptionalText("default_version").
	Text("owner").
	Text("owner_role_type").
	OptionalText("comment")

var dbtProject = g.PlainStruct("DbtProject").
	Time("CreatedOn").
	Text("Name").
	Text("DatabaseName").
	Text("SchemaName").
	OptionalText("SourceLocation").
	OptionalText("DefaultArgs").
	OptionalText("DefaultVersion").
	Text("Owner").
	Text("OwnerRoleType").
	OptionalText("Comment")

var DbtProjectsDef = g.NewInterface(
	"DbtProjects",
	"DbtProject",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-dbt-project",
	g.NewQueryStruct("CreateDbtProject").
		Create().
		OrReplace().
		SQL("DBT PROJECT").
		IfNotExists().
		Name().
		OptionalTextAssignment("FROM", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("DEFAULT_ARGS", g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("DEFAULT_VERSION", g.KindOfTPointer[DbtProjectDefaultVersion](), g.ParameterOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-dbt-project",
	g.NewQueryStruct("AlterDbtProject").
		Alter().
		SQL("DBT PROJECT").
		IfExists().
		Name().
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("DbtProjectSet").
				OptionalTextAssignment("DEFAULT_ARGS", g.ParameterOptions().SingleQuotes()).
				OptionalAssignment("DEFAULT_VERSION", g.KindOfTPointer[DbtProjectDefaultVersion](), g.ParameterOptions().NoQuotes()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("DbtProjectUnset").
				OptionalSQL("DEFAULT_ARGS").
				OptionalSQL("DEFAULT_VERSION").
				OptionalSQL("COMMENT"),
			g.KeywordOptions().SQL("UNSET"),
		).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-dbt-project",
	g.NewQueryStruct("DropDbtProject").
		Drop().
		SQL("DBT PROJECT").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-dbt-projects",
	dbtProjectDbRow,
	dbtProject,
	g.NewQueryStruct("ShowDbtProjects").
		Show().
		SQL("DBT PROJECTS").
		OptionalLike().
		OptionalIn(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-dbt-project",
	g.DbStruct("dbtProjectDetailsRow").
		Text("property").
		Text("value"),
	g.PlainStruct("DbtProjectDetails").
		Text("Property").
		Text("Value"),
	g.NewQueryStruct("DescribeDbtProject").
		Describe().
		SQL("DBT PROJECT").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
