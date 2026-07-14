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
	PostgresInstanceResetAccessRoleEnumDef = g.NewEnum(
		"PostgresInstanceResetAccessRole", "PostgresInstanceResetAccessRoles",
		"snowflake_admin", "application",
	)
	PostgresInstanceComputeFamilyEnumDef = g.NewEnum(
		"PostgresInstanceComputeFamily", "PostgresInstanceComputeFamilies",
		"STANDARD_M", "STANDARD_L", "STANDARD_XL", "STANDARD_2XL", "STANDARD_4XL", "STANDARD_8XL", "STANDARD_12XL", "STANDARD_24XL",
		"HIGHMEM_L", "HIGHMEM_XL", "HIGHMEM_2XL", "HIGHMEM_4XL", "HIGHMEM_8XL", "HIGHMEM_12XL", "HIGHMEM_16XL", "HIGHMEM_24XL", "HIGHMEM_32XL", "HIGHMEM_48XL",
		"BURST_XS", "BURST_S", "BURST_M",
	)
)

var postgresInstancesPairs = g.StructPair("postgresInstancesRow", "PostgresInstance").
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
	BoolFromText("is_ha", g.WithBoolTrueValue("true"), g.WithPlainFieldName("IsHighlyAvailable")).
	Number("retention_time").
	Enum("state", PostgresInstanceStateEnumDef).
	OptionalText("comment")

var postgresInstanceDetailPairs = g.StructPair("postgresInstanceDetailsRow", "PostgresInstanceProperty").
	Text("property").
	OptionalText("value", g.WithRequiredInPlain())

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
		TextAssignment("COMPUTE_FAMILY", g.ParameterOptions().SingleQuotes()).
		NumberAssignment("STORAGE_SIZE_GB", g.ParameterOptions().Required()).
		EnumAssignment(
			"AUTHENTICATION_AUTHORITY", PostgresInstanceAuthenticationAuthorityEnumDef,
			g.ParameterOptions().NoQuotes().Required(),
		).
		OptionalNumberAssignment("POSTGRES_VERSION", g.ParameterOptions()).
		OptionalIdentifier("NetworkPolicy", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
		OptionalBooleanAssignment("HIGH_AVAILABILITY", g.ParameterOptions()).
		OptionalIdentifier("StorageIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("STORAGE_INTEGRATION")).
		OptionalTextAssignment("POSTGRES_SETTINGS", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"Fork",
	"https://docs.snowflake.com/en/sql-reference/sql/create-postgres-instance",
	g.NewQueryStruct("ForkPostgresInstance").
		Create().
		SQL("POSTGRES INSTANCE").
		Name().
		Identifier("Fork", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("FORK").Required()).
		OptionalQueryStructField(
			"At",
			g.NewQueryStruct("PostgresInstanceForkAt").
				OptionalTextAssignment("TIMESTAMP", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				OptionalTextAssignment("OFFSET", g.ParameterOptions().NoQuotes().ArrowEquals()).
				WithValidation(g.ExactlyOneValueSet, "Timestamp", "Offset"),
			g.ListOptions().SQL("AT").Parentheses().NoComma(),
		).
		OptionalQueryStructField(
			"Before",
			g.NewQueryStruct("PostgresInstanceForkBefore").
				OptionalTextAssignment("TIMESTAMP", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				OptionalTextAssignment("OFFSET", g.ParameterOptions().NoQuotes().ArrowEquals()).
				WithValidation(g.ExactlyOneValueSet, "Timestamp", "Offset"),
			g.ListOptions().SQL("BEFORE").Parentheses().NoComma(),
		).
		OptionalTextAssignment("COMPUTE_FAMILY", g.ParameterOptions().SingleQuotes()).
		OptionalNumberAssignment("STORAGE_SIZE_GB", g.ParameterOptions()).
		OptionalBooleanAssignment("HIGH_AVAILABILITY", g.ParameterOptions()).
		OptionalTextAssignment("POSTGRES_SETTINGS", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "Fork").
		WithValidation(g.ConflictingFields, "At", "Before"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-postgres-instance",
	g.NewQueryStruct("AlterPostgresInstance").
		Alter().
		SQL("POSTGRES INSTANCE").
		IfExists().
		Name().
		RenameTo().
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("PostgresInstanceSet").
				OptionalIdentifier("NetworkPolicy", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("NETWORK_POLICY")).
				OptionalAssignment(
					"AUTHENTICATION_AUTHORITY",
					PostgresInstanceAuthenticationAuthorityEnumDef.Kind(),
					g.ParameterOptions().NoQuotes(),
				).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				OptionalBooleanAssignment("HIGH_AVAILABILITY", g.ParameterOptions()).
				OptionalTextAssignment("COMPUTE_FAMILY", g.ParameterOptions().SingleQuotes()).
				OptionalNumberAssignment("STORAGE_SIZE_GB", g.ParameterOptions()).
				OptionalIdentifier("StorageIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("STORAGE_INTEGRATION")).
				OptionalNumberAssignment("POSTGRES_VERSION", g.ParameterOptions()).
				OptionalNumberAssignment("MAINTENANCE_WINDOW_START", g.ParameterOptions()).
				OptionalTextAssignment("POSTGRES_SETTINGS", g.ParameterOptions().SingleQuotes()).
				OptionalQueryStructField(
					"Apply",
					g.NewQueryStruct("PostgresInstanceApply").
						OptionalSQL("IMMEDIATELY").
						OptionalTextAssignment("ON", g.ParameterOptions().SingleQuotes().NoEquals()).
						WithValidation(g.ExactlyOneValueSet, "Immediately", "On"),
					g.KeywordOptions().SQL("APPLY"),
				).
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
				EnumAssignmentWithFieldName("FOR", PostgresInstanceResetAccessRoleEnumDef, g.ParameterOptions().NoEquals().SingleQuotes().Required(), "ForRole"),
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
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-postgres-instances",
	postgresInstancesPairs,
	g.NewQueryStruct("ShowPostgresInstances").
		Show().
		SQL("POSTGRES INSTANCES").
		OptionalLike().
		OptionalStartsWith().
		OptionalLimitFrom(),
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-postgres-instance",
	postgresInstanceDetailPairs,
	g.NewQueryStruct("DescribePostgresInstance").
		Describe().
		SQL("POSTGRES INSTANCE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).WithEnums(
	PostgresInstanceStateEnumDef,
	PostgresInstanceAuthenticationAuthorityEnumDef,
	PostgresInstanceResetAccessRoleEnumDef,
	PostgresInstanceComputeFamilyEnumDef,
).WithCustomInterfaceMethod(
	"DescribeDetails",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
	"*PostgresInstanceDetails", "error",
).WithCustomInterfaceMethod(
	"CreateSafely",
	"CreateSafely creates the instance and polls until it reaches READY state.\nThe caller should set a deadline on ctx via context.WithTimeout.\nImplemented in postgres_instances_ext.go.",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*CreatePostgresInstanceRequest")},
	"*PostgresInstance", "error",
).WithCustomInterfaceMethod(
	"AlterSafely",
	"AlterSafely alters the instance and polls until it reaches READY state.\nThe caller should set a deadline on ctx via context.WithTimeout.\nImplemented in postgres_instances_ext.go.",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*AlterPostgresInstanceRequest")},
	"error",
)
