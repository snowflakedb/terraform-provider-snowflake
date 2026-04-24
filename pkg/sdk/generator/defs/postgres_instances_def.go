package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	PostgresInstanceStateEnumDef = g.NewEnum(
		"PostgresInstanceState", "PostgresInstanceStates",
		"CREATING", "RESTORING", "STARTING", "REPLAYING", "FINALIZING",
		"READY", "RESTARTING", "RESUMING", "SUSPENDING", "SUSPENDED",
	)
	PostgresInstanceAuthenticationAuthorityEnumDef = g.NewEnum(
		"PostgresInstanceAuthenticationAuthority", "PostgresInstanceAuthenticationAuthorities",
		"POSTGRES", "POSTGRES_OR_SNOWFLAKE",
	)
)

var postgresInstancesDef = g.NewInterface(
	"PostgresInstances",
	"PostgresInstance",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-postgres-instance",
	g.NewQueryStruct("CreatePostgresInstance").
		Create().
		SQL("POSTGRES INSTANCE").
		Name().
		OptionalIdentifier("Fork", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("FORK")).
		OptionalQueryStructField(
			"At",
			g.NewQueryStruct("PostgresInstanceForkAt").
				OptionalTextAssignment("TIMESTAMP", g.ParameterOptions().NoQuotes().ArrowEquals()).
				OptionalTextAssignment("OFFSET", g.ParameterOptions().NoQuotes().ArrowEquals()).
				WithValidation(g.ExactlyOneValueSet, "Timestamp", "Offset"),
			g.KeywordOptions().SQL("AT").MustParentheses(),
		).
		OptionalQueryStructField(
			"Before",
			g.NewQueryStruct("PostgresInstanceForkBefore").
				OptionalTextAssignment("TIMESTAMP", g.ParameterOptions().NoQuotes().ArrowEquals()).
				OptionalTextAssignment("OFFSET", g.ParameterOptions().NoQuotes().ArrowEquals()).
				WithValidation(g.ExactlyOneValueSet, "Timestamp", "Offset"),
			g.KeywordOptions().SQL("BEFORE").MustParentheses(),
		).
		OptionalTextAssignment("COMPUTE_FAMILY", g.ParameterOptions().SingleQuotes()).
		OptionalNumberAssignment("STORAGE_SIZE_GB", g.ParameterOptions()).
		OptionalAssignment(
			"AUTHENTICATION_AUTHORITY",
			PostgresInstanceAuthenticationAuthorityEnumDef.Kind(),
			g.ParameterOptions().NoQuotes(),
		).
		OptionalNumberAssignment("POSTGRES_VERSION", g.ParameterOptions()).
		OptionalTextAssignment("NETWORK_POLICY", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("HIGH_AVAILABILITY", g.ParameterOptions()).
		OptionalTextAssignment("STORAGE_INTEGRATION", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("POSTGRES_SETTINGS", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-postgres-instance",
	g.NewQueryStruct("AlterPostgresInstance").
		Alter().
		SQL("POSTGRES INSTANCE").
		IfExists().
		Name().
		Identifier("RenameTo", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("PostgresInstanceSet").
				OptionalTextAssignment("NETWORK_POLICY", g.ParameterOptions().SingleQuotes()).
			OptionalAssignment(
				"AUTHENTICATION_AUTHORITY",
				PostgresInstanceAuthenticationAuthorityEnumDef.Kind(),
				g.ParameterOptions().NoQuotes(),
			).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				OptionalBooleanAssignment("HIGH_AVAILABILITY", g.ParameterOptions()).
				OptionalTextAssignment("COMPUTE_FAMILY", g.ParameterOptions().SingleQuotes()).
				OptionalNumberAssignment("STORAGE_SIZE_GB", g.ParameterOptions()).
				OptionalTextAssignment("STORAGE_INTEGRATION", g.ParameterOptions().SingleQuotes()).
				OptionalNumberAssignment("POSTGRES_VERSION", g.ParameterOptions()).
				OptionalNumberAssignment("MAINTENANCE_WINDOW_START", g.ParameterOptions()).
				OptionalTextAssignment("POSTGRES_SETTINGS", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.AtLeastOneValueSet, "NetworkPolicy", "AuthenticationAuthority", "Comment", "HighAvailability", "ComputeFamily", "StorageSizeGb", "StorageIntegration", "PostgresVersion", "MaintenanceWindowStart", "PostgresSettings"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("PostgresInstanceUnset").
				OptionalSQL("COMMENT").
				OptionalSQL("POSTGRES_SETTINGS").
				OptionalSQL("NETWORK_POLICY").
				OptionalSQL("MAINTENANCE_WINDOW_START").
				OptionalSQL("STORAGE_INTEGRATION").
				WithValidation(g.AtLeastOneValueSet, "Comment", "PostgresSettings", "NetworkPolicy", "MaintenanceWindowStart", "StorageIntegration"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		OptionalSQL("SUSPEND").
		OptionalSQL("RESUME").
		OptionalQueryStructField(
			"ResetAccess",
			g.NewQueryStruct("PostgresInstanceResetAccess").
				OptionalTextAssignment("FOR", g.ParameterOptions().NoEquals().SingleQuotes()),
			g.KeywordOptions().SQL("RESET ACCESS"),
		).
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "Set", "Unset", "Suspend", "Resume", "ResetAccess", "SetTags", "UnsetTags"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-postgres-instance",
	g.NewQueryStruct("DropPostgresInstance").
		Drop().
		SQL("POSTGRES INSTANCE").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-postgres-instances",
	g.DbStruct("postgresInstancesRow").
		Text("name").
		Text("owner").
		Text("owner_role_type").
		Time("created_on").
		Time("updated_on").
		Text("type").
		OptionalText("origin").
		OptionalText("host").
		OptionalText("privatelink_service_identifier").
		Text("compute_family").
		Text("authentication_authority").
		Number("storage_size").
		Text("postgres_version").
		OptionalText("postgres_settings").
		Text("is_ha").
		Number("retention_time").
		Text("state").
		OptionalText("comment"),
	g.PlainStruct("PostgresInstance").
		Text("Name").
		Text("Owner").
		Text("OwnerRoleType").
		Time("CreatedOn").
		Time("UpdatedOn").
		Text("Type").
		OptionalText("Origin").
		OptionalText("Host").
		OptionalText("PrivatelinkServiceIdentifier").
		Text("ComputeFamily").
		Text("AuthenticationAuthority").
		Number("StorageSize").
		Text("PostgresVersion").
		OptionalText("PostgresSettings").
		Bool("IsHa").
		Number("RetentionTime").
		Field("State", PostgresInstanceStateEnumDef.Kind()).
		OptionalText("Comment"),
	g.NewQueryStruct("ShowPostgresInstances").
		Show().
		SQL("POSTGRES INSTANCES").
		OptionalLike().
		OptionalStartsWith().
		OptionalLimitFrom(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-postgres-instance",
	g.DbStruct("postgresInstanceDetailsRow").
		Text("property").
		Text("value"),
	g.PlainStruct("PostgresInstanceProperty").
		Text("Property").
		Text("Value"),
	g.NewQueryStruct("DescribePostgresInstance").
		Describe().
		SQL("POSTGRES INSTANCE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).WithEnums(
	PostgresInstanceStateEnumDef,
	PostgresInstanceAuthenticationAuthorityEnumDef,
)
