package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var PostgresMirrorTargetTableType = ssot.NewEnum(
	"PostgresMirrorTargetTableType", "PostgresMirrorTargetTableTypes",
	"SNOWFLAKE", "ICEBERG",
)

// PostgresMirrorAttrs holds the Postgres mirror fields. A mirror has no SHOW and no DDL: every operation calls a
// procedure of the SNOWFLAKE application, so the SDK stays handwritten and the attributes carry no SQL clauses.
// The Describe columns follow postgresMirrorDetailsPairs. The columns that are not resource attributes are Readonly.
// SDK int64 columns are ints here, because ssot has no int64 type and no enabled generator reads them.
// refresh_interval is a block for Terraform, and DESCRIBE_MIRROR returns it as one string, so that column is a separate
// Readonly attribute whose ShowName is refresh_interval.
// drop_target_database_on_destroy is read only by Delete, so it uses SkipCreate and has no update mapping.
// The required identifiers are strings so that the generator adds no validation for them. The ext schema modification
// sets their diff suppression and clears the validation the generator adds to the optional identifiers.
var PostgresMirrorAttrs = struct {
	PostgresInstance            ssot.Attribute
	PostgresDatabase            ssot.Attribute
	TargetDatabase              ssot.Attribute
	PostgresTables              ssot.Attribute
	PostgresSchemas             ssot.Attribute
	RefreshInterval             ssot.Attribute
	TargetTableType             ssot.Attribute
	ExternalVolume              ssot.Attribute
	Warehouse                   ssot.Attribute
	DropTargetDatabaseOnDestroy ssot.Attribute
	Name                        ssot.Attribute
	Status                      ssot.Attribute
	DescribeRefreshInterval     ssot.Attribute
	FeedOnly                    ssot.Attribute
	LastOperationId             ssot.Attribute
	LastApplyTime               ssot.Attribute
	LastOperationTime           ssot.Attribute
	ErrorMessage                ssot.Attribute
	CatalogIntegrationName      ssot.Attribute
	RecentQueryDurations        ssot.Attribute
	ApplyLagSeconds             ssot.Attribute
	ApplyLagOperations          ssot.Attribute
	ApplyLagBytes               ssot.Attribute
	LastAppliedCommitLsn        ssot.Attribute
	LastAppliedCommitTime       ssot.Attribute
	PostgresStatus              ssot.Attribute
}{
	PostgresInstance: ssot.Attribute{
		Name:         "postgres_instance",
		Type:         ssot.AttrString,
		Required:     true,
		ForceNew:     true,
		OptionalInDb: true,
		Description:  "Specifies the source Postgres instance (case-sensitive). The role used by the provider must own the instance.",
	},
	PostgresDatabase: ssot.Attribute{
		Name:         "postgres_database",
		Type:         ssot.AttrString,
		Required:     true,
		ForceNew:     true,
		OptionalInDb: true,
		Description:  "Specifies the source database on the Postgres instance.",
	},
	TargetDatabase: ssot.Attribute{
		Name:         "target_database",
		Type:         ssot.AttrString,
		Required:     true,
		ForceNew:     true,
		OptionalInDb: true,
		Description:  "Specifies the Snowflake database the mirror creates and keeps up to date. The database must not exist yet, and Snowflake uppercases its name. The mirror protects the database, so do not manage it with `snowflake_database`.",
	},
	PostgresTables: ssot.Attribute{
		Name: "postgres_tables",
		Type: ssot.AttrSet,
		Elem: &ssot.Attribute{Type: ssot.AttrString},
		Description: "Specifies the source tables to mirror, as schema-qualified names (`schema.table`) written the way `DESCRIBE_MIRROR` returns them: unquoted and lowercase, " +
			"except for parts Postgres requires to be quoted, such as `public.\"MixedCase\"`. Other spellings of the same table show as a change on every plan. " +
			"Tables are added and removed in place, without copying the other tables again.",
	},
	PostgresSchemas: ssot.Attribute{
		Name: "postgres_schemas",
		Type: ssot.AttrSet,
		Elem: &ssot.Attribute{Type: ssot.AttrString},
		Description: "Specifies the source schemas to mirror. Every current and future table in them is included. Schema names are literal and case-sensitive, written without Postgres quoting, such as `MixedCase` rather than `\"MixedCase\"`. " +
			"Schemas are added and removed in place.",
	},
	RefreshInterval: ssot.Attribute{
		Name:        "refresh_interval",
		Type:        ssot.AttrObject,
		SkipCreate:  true,
		Description: "Specifies how often changes are merged into the target tables. Set exactly one of `hours`, `minutes`, or `seconds`. When not set, Snowflake uses its default interval. Removing the block keeps the current interval, because `ALTER_MIRROR` cannot reset it.",
		Children: []ssot.Attribute{
			{
				Name:          "hours",
				Type:          ssot.AttrInt,
				ConflictsWith: []ssot.Attribute{{Name: "minutes"}, {Name: "seconds"}},
				AtLeastOneOf:  []ssot.Attribute{{Name: "hours"}, {Name: "minutes"}, {Name: "seconds"}},
				Description:   "Specifies the refresh interval in hours. (conflicts with `minutes` and `seconds`)",
			},
			{
				Name:          "minutes",
				Type:          ssot.AttrInt,
				ConflictsWith: []ssot.Attribute{{Name: "hours"}, {Name: "seconds"}},
				AtLeastOneOf:  []ssot.Attribute{{Name: "hours"}, {Name: "minutes"}, {Name: "seconds"}},
				Description:   "Specifies the refresh interval in minutes. (conflicts with `hours` and `seconds`)",
			},
			{
				Name:          "seconds",
				Type:          ssot.AttrInt,
				ConflictsWith: []ssot.Attribute{{Name: "hours"}, {Name: "minutes"}},
				AtLeastOneOf:  []ssot.Attribute{{Name: "hours"}, {Name: "minutes"}, {Name: "seconds"}},
				Description:   "Specifies the refresh interval in seconds. (conflicts with `hours` and `minutes`)",
			},
		},
	},
	TargetTableType: ssot.Attribute{
		Name:         "target_table_type",
		Type:         ssot.AttrEnum,
		Enum:         PostgresMirrorTargetTableType,
		ForceNew:     true,
		OptionalInDb: true,
		DiffSuppress: ssot.SuppressDescribe,
		Description:  "Specifies the format of the target tables. When not set, Snowflake uses `SNOWFLAKE`.",
	},
	ExternalVolume: ssot.Attribute{
		Name:           "external_volume",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "AccountObjectIdentifier",
		ForceNew:       true,
		OptionalInShow: true,
		OptionalInDb:   true,
		CreateMapping:  ssot.AttributeMappedValueCreate("quotedAccountObjectName"),
		Description:    "Specifies the external volume for Iceberg target tables. Valid only with `target_table_type = \"ICEBERG\"`. When not set, the tables use Snowflake-managed storage. Grant `USAGE` on the volume to `APPLICATION SNOWFLAKE` first.",
	},
	Warehouse: ssot.Attribute{
		Name:           "warehouse",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "AccountObjectIdentifier",
		OptionalInShow: true,
		OptionalInDb:   true,
		DiffSuppress:   ssot.SuppressDescribe,
		CreateMapping:  ssot.AttributeMappedValueCreate("bareAccountObjectName"),
		Description:    "Specifies the warehouse the mirror's apply task runs on. When not set, the task runs on serverless compute. Snowflake uppercases the name, so the warehouse itself must have an uppercase name; the case used in the configuration does not matter. Grant `USAGE` on the warehouse to `APPLICATION SNOWFLAKE` first.",
	},
	DropTargetDatabaseOnDestroy: ssot.Attribute{
		Name:       "drop_target_database_on_destroy",
		Type:       ssot.AttrBool,
		NativeBool: true,
		SkipCreate: true,
		Description: "Specifies whether destroying the mirror also drops its target database. A kept target database can never be the target of a mirror again, so with `false`, a replacement that keeps the same `target_database` fails. " +
			"A kept target database stays owned by the `SNOWFLAKE` application. To manage or drop it, first transfer its ownership with `GRANT OWNERSHIP ON DATABASE`, which needs the `MANAGE GRANTS` privilege.",
	},
	Name:                    ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	Status:                  ssot.Attribute{Name: "status", Type: ssot.AttrString, Readonly: true, OptionalInDb: true},
	DescribeRefreshInterval: ssot.Attribute{Name: "describe_refresh_interval", ShowName: "refresh_interval", Type: ssot.AttrString, Readonly: true, OptionalInDb: true},
	FeedOnly:                ssot.Attribute{Name: "feed_only", Type: ssot.AttrBool, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	LastOperationId:         ssot.Attribute{Name: "last_operation_id", Type: ssot.AttrInt, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	LastApplyTime:           ssot.Attribute{Name: "last_apply_time", Type: ssot.AttrTime, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	LastOperationTime:       ssot.Attribute{Name: "last_operation_time", Type: ssot.AttrTime, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	ErrorMessage:            ssot.Attribute{Name: "error_message", Type: ssot.AttrString, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	CatalogIntegrationName:  ssot.Attribute{Name: "catalog_integration_name", Type: ssot.AttrString, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	RecentQueryDurations:    ssot.Attribute{Name: "recent_query_durations", Type: ssot.AttrString, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	ApplyLagSeconds:         ssot.Attribute{Name: "apply_lag_seconds", Type: ssot.AttrInt, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	ApplyLagOperations:      ssot.Attribute{Name: "apply_lag_operations", Type: ssot.AttrInt, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	ApplyLagBytes:           ssot.Attribute{Name: "apply_lag_bytes", Type: ssot.AttrInt, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	LastAppliedCommitLsn:    ssot.Attribute{Name: "last_applied_commit_lsn", Type: ssot.AttrString, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	LastAppliedCommitTime:   ssot.Attribute{Name: "last_applied_commit_time", Type: ssot.AttrTime, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	PostgresStatus:          ssot.Attribute{Name: "postgres_status", Type: ssot.AttrString, Readonly: true, OptionalInShow: true, OptionalInDb: true},
}

var PostgresMirror = ssot.Object{
	Name:                   "PostgresMirror",
	NamePlural:             "PostgresMirrors",
	IdType:                 "AccountObjectIdentifier",
	IdentityNameForceNew:   true,
	AttributesModification: true,
	PreviewFeature:         true,
	// CREATE_MIRROR takes its arguments in a nested request, which has the With* methods the create helpers call.
	SdkCreateRequestObjectName: "PostgresMirrorArguments",
	Generators: []ssot.GeneratorOption{
		ssot.WithObjectParametersAssert(),
		ssot.WithResourceParametersAssert(),
		ssot.WithResource(),
		ssot.WithModel(),
		ssot.WithResourceAssert(),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseId(ssot.Generated()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithShowByIdInSdk(ssot.Omitted()),
		ssot.WithSetFieldsNotSetByRead(ssot.Omitted()),
		ssot.WithShowOutputSchema(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.ExtCall()),
		ssot.WithIdentitySchema(ssot.ExtCall()),
		ssot.WithSetStateToValuesFromConfig(ssot.ExtCall()),
	},
	Attributes: []ssot.Attribute{
		PostgresMirrorAttrs.PostgresInstance,
		PostgresMirrorAttrs.PostgresDatabase,
		PostgresMirrorAttrs.TargetDatabase,
		PostgresMirrorAttrs.PostgresTables,
		PostgresMirrorAttrs.PostgresSchemas,
		PostgresMirrorAttrs.RefreshInterval,
		PostgresMirrorAttrs.TargetTableType,
		PostgresMirrorAttrs.ExternalVolume,
		PostgresMirrorAttrs.Warehouse,
		PostgresMirrorAttrs.DropTargetDatabaseOnDestroy,
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&PostgresMirrorAttrs.Name,
			&PostgresMirrorAttrs.PostgresInstance,
			&PostgresMirrorAttrs.PostgresDatabase,
			&PostgresMirrorAttrs.TargetDatabase,
			&PostgresMirrorAttrs.Status,
			&PostgresMirrorAttrs.DescribeRefreshInterval,
			&PostgresMirrorAttrs.TargetTableType,
			&PostgresMirrorAttrs.ExternalVolume,
			&PostgresMirrorAttrs.Warehouse,
			&PostgresMirrorAttrs.PostgresTables,
			&PostgresMirrorAttrs.PostgresSchemas,
			&PostgresMirrorAttrs.FeedOnly,
			&PostgresMirrorAttrs.LastOperationId,
			&PostgresMirrorAttrs.LastApplyTime,
			&PostgresMirrorAttrs.LastOperationTime,
			&PostgresMirrorAttrs.ErrorMessage,
			&PostgresMirrorAttrs.CatalogIntegrationName,
			&PostgresMirrorAttrs.RecentQueryDurations,
			&PostgresMirrorAttrs.ApplyLagSeconds,
			&PostgresMirrorAttrs.ApplyLagOperations,
			&PostgresMirrorAttrs.ApplyLagBytes,
			&PostgresMirrorAttrs.LastAppliedCommitLsn,
			&PostgresMirrorAttrs.LastAppliedCommitTime,
			&PostgresMirrorAttrs.PostgresStatus,
		},
	},
}
