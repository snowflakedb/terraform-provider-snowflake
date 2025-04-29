package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var ImageRepositoriesDef = g.NewInterface(
	"ImageRepositories",
	"ImageRepository",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-image-repository",
	g.NewQueryStruct("CreateImageRepository").
		Create().
		OrReplace().
		SQL("IMAGE REPOSITORY").
		IfNotExists().
		Name().
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-image-repository",
	g.NewQueryStruct("AlterImageRepository").
		Alter().
		SQL("IMAGE REPOSITORY").
		IfExists().
		Name().
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("ImageRepositorySet").
				// TODO(SNOW-2070753): use COMMENT in unset and here use OptionalComment
				OptionalAssignment("COMMENT", "StringAllowEmpty", g.ParameterOptions()),
			g.KeywordOptions().SQL("SET"),
		).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Set"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-image-repository",
	g.NewQueryStruct("DropImageRepository").
		Drop().
		SQL("IMAGE REPOSITORY").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-image-repositories",
	g.DbStruct("imageRepositoriesRow").
		Time("created_on").
		Text("name").
		Text("database_name").
		Text("schema_name").
		Text("repository_url").
		Text("owner").
		Text("owner_role_type").
		Text("comment").
		Text("privatelink_repository_url"),
	g.PlainStruct("ImageRepository").
		Time("CreatedOn").
		Text("Name").
		Text("DatabaseName").
		Text("SchemaName").
		Text("RepositoryUrl").
		Text("Owner").
		Text("OwnerRoleType").
		Text("Comment").
		Text("PrivatelinkRepositoryUrl"),
	g.NewQueryStruct("ShowImageRepositories").
		Show().
		SQL("IMAGE REPOSITORIES").
		OptionalLike().
		OptionalIn(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
)
