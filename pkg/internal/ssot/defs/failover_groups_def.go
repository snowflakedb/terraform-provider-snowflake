package defs

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// Pages:
// https://docs.snowflake.com/en/sql-reference/sql/create-failover-group
// https://docs.snowflake.com/en/sql-reference/sql/alter-failover-group
// https://docs.snowflake.com/en/sql-reference/sql/drop-failover-group
// https://docs.snowflake.com/en/sql-reference/sql/show-failover-groups
//
// Failover groups already have a handwritten SDK definition
// (pkg/sdk/generator/defs/failover_groups_def.go, pkg/sdk/failover_groups_*.go). It stays until every operation has moved:
// no SdkDefineFromSsot() here yet, and each operation is replaced with its FromObject call. This definition follows the
// pages where they disagree with the handwritten one:
//   - ALLOWED_EXTERNAL_VOLUMES, TAG and ERROR_INTEGRATION on CREATE; EXTERNAL VOLUMES in ADD, MOVE and REMOVE; the object types
//     LISTINGS and PROFILES; SET and UNSET of COMMENT, ERROR_INTEGRATION and READ_ONLY; SET TAG and UNSET TAG; SUSPEND IMMEDIATE.
//     The handwritten SDK has none of them.
//   - SET may carry the allowed lists, and UNSET several comma-separated properties. The handwritten SET has only OBJECT_TYPES,
//     ALLOWED_INTEGRATION_TYPES, REPLICATION_SCHEDULE, OPTIMIZED_REFRESH and RPO_ASSURANCE, and UNSET exactly one.
//   - The ADD accounts block ends with IGNORE EDITION CHECK (with spaces). The handwritten ADD uses IGNORE_EDITION_CHECK.
//   - SHOW has 22 columns; the handwritten row has 20 (no is_listing_auto_fulfillment_group, no read_only).
//
// Documentation gaps:
//   - The CREATE page has no COMMENT, OR REPLACE and RPO_ASSURANCE (RPO assurance is set after CREATE, with ALTER SET).
//   - There is no DESCRIBE FAILOVER GROUP, and the SHOW page publishes no ALLOWED_DATABASES / ALLOWED_SHARES /
//     ALLOWED_EXTERNAL_VOLUMES columns. The handwritten SHOW DATABASES / SHARES IN FAILOVER GROUP operations stay in the SDK ext.
//   - The SET page does not say which SET properties can be combined. The syntax blocks and examples show OBJECT_TYPES with the
//     allowed lists and REPLICATION_SCHEDULE with OPTIMIZED_REFRESH, while OPTIMIZED_REFRESH = TRUE and RPO_ASSURANCE = TRUE
//     conflict and RPO_ASSURANCE = TRUE needs no REPLICATION_SCHEDULE. SDK additional validations (ext) cover that.
//   - The SHOW page lists no ERROR_INTEGRATION column.
//
// Failover groups require Business Critical Edition or higher. DROP documents IF EXISTS. CREATE and the secondary CREATE have
// IF NOT EXISTS and no OR REPLACE. ObjectTypeFailoverGroup, its plural, and the singular-to-plural map already exist in
// pkg/sdk/object_types.go.
//
// The enums have their own names (like ReplicationGroup*) because sdk.IntegrationType, sdk.PluralObjectType and
// sdk.FailoverGroupSecondaryState are the handwritten types used by the resource. Unifying them with the replication group
// enums is a later step.
//
// Account object types other than SHARES and DATABASES can be a part of only one failover or replication group per account.

// FailoverGroupObjectType is the OBJECT_TYPES list on the create and alter pages.
// It is not PluralObjectType. That type has no ACCOUNT PARAMETERS value; the handwritten SHOW converter rewrites
// the SHOW token PARAMETERS into that text.
var FailoverGroupObjectType = ssot.NewEnum(
	"FailoverGroupObjectType", "FailoverGroupObjectTypes",
	"ACCOUNT PARAMETERS",
	"DATABASES",
	"EXTERNAL VOLUMES",
	"INTEGRATIONS",
	"LISTINGS",
	"NETWORK POLICIES",
	"PROFILES",
	"RESOURCE MONITORS",
	"ROLES",
	"SHARES",
	"USERS",
	"WAREHOUSES",
)

// FailoverGroupIntegrationType has the same five values as sdk.IntegrationType.
var FailoverGroupIntegrationType = ssot.NewEnum(
	"FailoverGroupIntegrationType", "FailoverGroupIntegrationTypes",
	"SECURITY INTEGRATIONS",
	"API INTEGRATIONS",
	"STORAGE INTEGRATIONS",
	"EXTERNAL ACCESS INTEGRATIONS",
	"NOTIFICATION INTEGRATIONS",
)

// The SHOW page says secondary_state is started, suspended, or NULL. NULL is an absent value.
// The handwritten enum is uppercase and has a NULL constant.
var FailoverGroupSecondaryState = ssot.NewEnum(
	"FailoverGroupSecondaryState", "FailoverGroupSecondaryStates",
	"started", "suspended",
)

func failoverGroupEnumList(enum *ssot.Enum, required bool) ssot.SqlClause {
	opts := []ssot.SqlClauseOption{ssot.NoParentheses()}
	if required {
		opts = append(opts, ssot.Required())
	}
	return ssot.ListParameter(enum.Name, opts...)
}

func failoverGroupIdList(idType string, required bool) ssot.SqlClause {
	opts := []ssot.SqlClauseOption{ssot.NoParentheses()}
	if required {
		opts = append(opts, ssot.Required())
	}
	return ssot.ListParameter(idType, opts...)
}

func failoverGroupMoveList(name, sql string) ssot.Attribute {
	opts := []ssot.SqlClauseOption{ssot.NoParentheses(), ssot.NoEquals()}
	if sql != "" {
		opts = append([]ssot.SqlClauseOption{ssot.Sql(sql)}, opts...)
	}
	return ssot.Attribute{
		Name:      name,
		Type:      ssot.AttrList,
		Elem:      &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"},
		CreateSql: ssot.ListParameter("AccountObjectIdentifier", opts...),
	}
}

// Reverse prints the identifiers before the SQL text: ADD "db1" TO ALLOWED_DATABASES.
// The field keeps the handwritten name (AllowedDatabases), not the one derived from the SQL text (ToAllowedDatabases).
func failoverGroupReverseIdentifierList(name, sql, idType string) ssot.Attribute {
	return ssot.Attribute{
		Name: name,
		Type: ssot.AttrList,
		Elem: &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: idType},
		CreateSql: ssot.ListParameter(
			idType,
			ssot.Sql(sql),
			ssot.NoParentheses(),
			ssot.Reverse(),
			ssot.FieldName(genhelpers.SnakeCaseToCamel(name)),
		),
	}
}

// FailoverGroupAttrs holds reusable failover group fields.
var FailoverGroupAttrs = struct {
	ObjectTypes                   ssot.Attribute
	AllowedDatabases              ssot.Attribute
	AllowedExternalVolumes        ssot.Attribute
	AllowedShares                 ssot.Attribute
	AllowedIntegrationTypes       ssot.Attribute
	AllowedAccounts               ssot.Attribute
	IgnoreEditionCheck            ssot.Attribute
	ReplicationSchedule           ssot.Attribute
	OptimizedRefresh              ssot.Attribute
	RpoAssurance                  ssot.Attribute
	ErrorIntegration              ssot.Attribute
	Comment                       ssot.Attribute
	ReadOnly                      ssot.Attribute
	Primary                       ssot.Attribute
	RegionGroup                   ssot.Attribute
	SnowflakeRegion               ssot.Attribute
	CreatedOn                     ssot.Attribute
	AccountName                   ssot.Attribute
	Name                          ssot.Attribute
	Type                          ssot.Attribute
	IsPrimary                     ssot.Attribute
	OrganizationName              ssot.Attribute
	AccountLocator                ssot.Attribute
	SecondaryState                ssot.Attribute
	NextScheduledRefresh          ssot.Attribute
	Owner                         ssot.Attribute
	IsListingAutoFulfillmentGroup ssot.Attribute
	IsOptimizedRefreshEnabled     ssot.Attribute
}{
	ObjectTypes: ssot.Attribute{
		Name:          "object_types",
		Type:          ssot.AttrSet,
		Required:      true,
		CreateMapping: ssot.ListValueCreate("toFailoverGroupObjectTypes"),
		Description:   "Types of objects for which replication and failover from the source account to the target accounts is enabled.",
		Elem:          &ssot.Attribute{Type: ssot.AttrEnum, Enum: FailoverGroupObjectType},
		// CSV. NULL or empty when a secondary group lists none. The plain field stays a slice.
		// SHOW uses PARAMETERS for account parameters. That token is mapped to ACCOUNT PARAMETERS.
		// The SHOW schema generator does not map an enum slice, so the schema ext maps the column.
		OptionalInDb: true,
		CustomParser: "ParseCommaSeparatedFailoverGroupObjectTypeArray",
		CreateSql:    failoverGroupEnumList(FailoverGroupObjectType, true),
		SetSql:       failoverGroupEnumList(FailoverGroupObjectType, false),
	},
	AllowedDatabases: ssot.Attribute{
		Name:        "allowed_databases",
		Type:        ssot.AttrSet,
		Description: "Databases enabled for replication and failover. OBJECT_TYPES must include DATABASES.",
		Elem:        &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"},
		CreateSql:   failoverGroupIdList("AccountObjectIdentifier", false),
		SetSql:      failoverGroupIdList("AccountObjectIdentifier", false),
	},
	AllowedExternalVolumes: ssot.Attribute{
		Name:        "allowed_external_volumes",
		Type:        ssot.AttrSet,
		Description: "External volumes enabled for replication and failover. OBJECT_TYPES must include EXTERNAL VOLUMES.",
		Elem:        &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"},
		CreateSql:   failoverGroupIdList("AccountObjectIdentifier", false),
		SetSql:      failoverGroupIdList("AccountObjectIdentifier", false),
	},
	AllowedShares: ssot.Attribute{
		Name:        "allowed_shares",
		Type:        ssot.AttrSet,
		Description: "Shares enabled for replication and failover. OBJECT_TYPES must include SHARES.",
		Elem:        &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"},
		CreateSql:   failoverGroupIdList("AccountObjectIdentifier", false),
		SetSql:      failoverGroupIdList("AccountObjectIdentifier", false),
	},
	AllowedIntegrationTypes: ssot.Attribute{
		Name:        "allowed_integration_types",
		Type:        ssot.AttrSet,
		Description: "Integration types enabled for replication and failover. OBJECT_TYPES must include INTEGRATIONS. SECURITY INTEGRATIONS also requires ROLES.",
		Elem:        &ssot.Attribute{Type: ssot.AttrEnum, Enum: FailoverGroupIntegrationType},
		// The column is always present. It is empty when no integration types are set.
		// SHOW prints short tokens such as API and EXTERNAL_ACCESS. Those tokens are mapped to the
		// create-page values (the handwritten converter replaces underscores and appends " INTEGRATIONS").
		// The SHOW schema generator does not map an enum slice, so the schema ext maps the column.
		OptionalInDb: true,
		CustomParser: "ParseCommaSeparatedFailoverGroupIntegrationTypeArray",
		CreateSql:    failoverGroupEnumList(FailoverGroupIntegrationType, false),
		SetSql:       failoverGroupEnumList(FailoverGroupIntegrationType, false),
	},
	AllowedAccounts: ssot.Attribute{
		Name:          "allowed_accounts",
		Type:          ssot.AttrSet,
		Required:      true,
		CreateMapping: ssot.ListValueCreate("toFailoverGroupAllowedAccounts"),
		Description:   "Target accounts, in the form `<organization_name>.<account_name>`, to which replication and failover is enabled. Snowflake adds the source account to every failover group, so do not list it.",
		Elem:          &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountIdentifier"},
		// The handwritten converter skips a value that is not org.account instead of failing; the parser keeps that.
		OptionalInDb: true,
		CustomParser: "ParseCommaSeparatedFailoverGroupAllowedAccountArray",
		CreateSql:    failoverGroupIdList("AccountIdentifier", true),
	},
	IgnoreEditionCheck: ssot.Attribute{
		Name:        "ignore_edition_check",
		Type:        ssot.AttrBool,
		NativeBool:  true,
		CreateOnly:  true,
		Description: "Allows replicating objects to accounts on lower editions.",
		CreateSql:   ssot.BoolKeyword(ssot.Sql("IGNORE EDITION CHECK")),
	},
	ReplicationSchedule: ssot.Attribute{
		// RPO_ASSURANCE needs the schedule to be unset.
		ConflictsWith:  []ssot.Attribute{{Name: "rpo_assurance"}},
		Name:           "replication_schedule",
		Type:           ssot.AttrString,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Schedule for refreshing secondary failover groups. An interval (`<num> MINUTE`) or `USING CRON <expr> <time_zone>`. Remove the attribute to unset the schedule.",
	},
	OptimizedRefresh: ssot.Attribute{
		// RPO_ASSURANCE enables Optimized Refresh itself.
		ConflictsWith: []ssot.Attribute{{Name: "rpo_assurance"}},
		Name:          "optimized_refresh",
		Type:          ssot.AttrBool,
		NativeBool:    true,
		Description:   "Specifies whether the failover group uses Optimized Refresh. Set on the primary failover group only. When true, a `replication_schedule` with an interval of 6 hours or less is required (Snowflake recommends `10 MINUTE` or less). Cannot be set together with `rpo_assurance`.",
	},
	RpoAssurance: ssot.Attribute{
		// See ReplicationSchedule and OptimizedRefresh.
		ConflictsWith: []ssot.Attribute{{Name: "replication_schedule"}, {Name: "optimized_refresh"}},
		Name:          "rpo_assurance",
		Type:          ssot.AttrBool,
		NativeBool:    true,
		Description:   "Specifies whether the failover group uses RPO Assurance. Set on the primary failover group only. Needs no replication schedule, enables Optimized Refresh automatically, and cannot be set together with `optimized_refresh`.",
	},
	// SHOW returns the integration name on a primary group and NULL on a secondary group. The documentation does not list the column.
	ErrorIntegration: ssot.Attribute{
		Name:           "error_integration",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "AccountObjectIdentifier",
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Notification integration that receives refresh error notifications.",
		CreateSql:      ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Equals()),
		SetSql:         ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Equals()),
	},
	Comment: ssot.Attribute{
		Name:           "comment",
		Type:           ssot.AttrString,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Comment for the failover group.",
	},
	// The documentation lists the SHOW column as read_only, but SHOW returns is_primary_read_only: TRUE or FALSE on a primary group, NULL on a secondary group.
	ReadOnly: ssot.Attribute{
		Name:           "read_only",
		ShowName:       "is_primary_read_only",
		Type:           ssot.AttrBool,
		NativeBool:     true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Enables the read-only mode of a primary failover group. Snowflake rejects changes to `object_types`, `allowed_databases`, `allowed_external_volumes`, `allowed_shares` and `allowed_integration_types` in this mode, so set it to false (or remove it) before changing them.",
	},
	// SHOW primary is org.account.name. It is not the AS REPLICA OF argument.
	Primary: ssot.Attribute{
		Name:           "primary",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "ExternalObjectIdentifier",
		Readonly:       true,
		CustomParser:   "ParseExternalObjectIdentifier",
		Description:    "Primary failover group, as org.account.name. On a primary group this is the group itself. On a secondary group this is the source group.",
	},
	// The column is present only for organizations that span multiple region groups.
	RegionGroup:     ssot.Attribute{Name: "region_group", Type: ssot.AttrString, Readonly: true, Description: "Region group where the account is located. Visible only to organizations that span multiple region groups."},
	SnowflakeRegion: ssot.Attribute{Name: "snowflake_region", Type: ssot.AttrString, Readonly: true},
	CreatedOn:       ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	AccountName:     ssot.Attribute{Name: "account_name", Type: ssot.AttrString, Readonly: true},
	Name:            ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	// The only valid value is FAILOVER (the handwritten row keeps a string).
	Type: ssot.Attribute{Name: "type", Type: ssot.AttrString, Readonly: true, Description: "Group type. The only valid value is FAILOVER."},
	IsPrimary: ssot.Attribute{
		Name:        "is_primary",
		Type:        ssot.AttrBool,
		NativeBool:  true,
		Readonly:    true,
		Description: "Indicates if the failover group is primary. When Terraform detects that the primary status changed outside of Terraform, the resource is recreated.",
	},
	OrganizationName: ssot.Attribute{Name: "organization_name", Type: ssot.AttrString, Readonly: true},
	AccountLocator:   ssot.Attribute{Name: "account_locator", Type: ssot.AttrString, Readonly: true},
	SecondaryState: ssot.Attribute{
		Name:           "secondary_state",
		Type:           ssot.AttrEnum,
		Enum:           FailoverGroupSecondaryState,
		Readonly:       true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "State of scheduled refresh: started or suspended. NULL when no replication schedule is set.",
	},
	// The page calls this a date and time. A nullable time is OptionalInDb only, so the plain field stays time.Time (the same as the replication group).
	NextScheduledRefresh: ssot.Attribute{Name: "next_scheduled_refresh", Type: ssot.AttrTime, Readonly: true, OptionalInDb: true},
	// Usually owner is not null, but the page says it is NULL when the group is in a different region.
	Owner: ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true, OptionalInShow: true, OptionalInDb: true, Description: "Role with OWNERSHIP on the group. NULL when the group is in a different region."},
	// The handwritten SHOW row does not include this column. The page does.
	IsListingAutoFulfillmentGroup: ssot.Attribute{Name: "is_listing_auto_fulfillment_group", Type: ssot.AttrBool, Readonly: true, Description: "TRUE when the group is used for Cross-Cloud Auto-Fulfillment."},
	// The SHOW column of the optimized_refresh property.
	IsOptimizedRefreshEnabled: ssot.Attribute{Name: "is_optimized_refresh_enabled", Type: ssot.AttrBool, NativeBool: true, Readonly: true, Description: "TRUE if Optimized Refresh is enabled. When FALSE, refreshes run as Replication Classic unless RPO Assurance is enabled."},
}

// AS REPLICA OF is the secondary CREATE argument. It is not the SHOW column primary.
// The SHOW column is read-only, so the source attribute has no create clause on the primary CREATE.
var failoverGroupAsReplicaOf = ssot.Attribute{
	Name:           "primary",
	Type:           ssot.AttrIdentifier,
	IdentifierType: "ExternalObjectIdentifier",
	Required:       true,
	ForceNew:       true,
	Description:    "Primary failover group to replicate, in the form `<organization_name>.<account_name>.<failover_group_name>`.",
	CreateSql:      ssot.IdentifierKeyword("ExternalObjectIdentifier", ssot.Sql("AS REPLICA OF"), ssot.Required(), ssot.FieldName("AsReplicaOf")),
}

// TODO [SNOW-4490849]: MOVE needs at least one of the three lists, but ChildrenRule also counts the always set `to`, so it cannot express this.
// The SDK extension validates it. Move them here when validations can be declared in a better way.
var failoverGroupMove = ssot.Attribute{
	Name:      "move",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		failoverGroupMoveList("databases", ""),
		failoverGroupMoveList("external_volumes", "EXTERNAL VOLUMES"),
		failoverGroupMoveList("shares", ""),
		{
			Name:           "to",
			Type:           ssot.AttrIdentifier,
			IdentifierType: "AccountObjectIdentifier",
			CreateSql:      ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Sql("TO FAILOVER GROUP"), ssot.Required(), ssot.FieldName("To")),
		},
	},
}

// TODO [SNOW-4490849]: ADD needs exactly one of the four lists, but ChildrenRule also counts ignore_edition_check, so it cannot express this.
// IGNORE EDITION CHECK is allowed only together with allowed_accounts, a dependency between fields that no rule expresses.
// The SDK extension validates both. Move them here when validations can be declared in a better way.
var failoverGroupAdd = ssot.Attribute{
	Name:      "add",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		failoverGroupReverseIdentifierList("allowed_databases", "TO ALLOWED_DATABASES", "AccountObjectIdentifier"),
		failoverGroupReverseIdentifierList("allowed_external_volumes", "TO ALLOWED_EXTERNAL_VOLUMES", "AccountObjectIdentifier"),
		failoverGroupReverseIdentifierList("allowed_shares", "TO ALLOWED_SHARES", "AccountObjectIdentifier"),
		failoverGroupReverseIdentifierList("allowed_accounts", "TO ALLOWED_ACCOUNTS", "AccountIdentifier"),
		{
			Name:      "ignore_edition_check",
			Type:      ssot.AttrBool,
			CreateSql: ssot.BoolKeyword(ssot.Sql("IGNORE EDITION CHECK")),
		},
	},
}

var failoverGroupRemove = ssot.Attribute{
	Name:         "remove",
	Type:         ssot.AttrObject,
	SqlOnly:      true,
	ChildrenRule: ssot.ChildrenExactlyOne,
	CreateSql:    ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		failoverGroupReverseIdentifierList("allowed_databases", "FROM ALLOWED_DATABASES", "AccountObjectIdentifier"),
		failoverGroupReverseIdentifierList("allowed_external_volumes", "FROM ALLOWED_EXTERNAL_VOLUMES", "AccountObjectIdentifier"),
		failoverGroupReverseIdentifierList("allowed_shares", "FROM ALLOWED_SHARES", "AccountObjectIdentifier"),
		failoverGroupReverseIdentifierList("allowed_accounts", "FROM ALLOWED_ACCOUNTS", "AccountIdentifier"),
	},
}

// SUSPEND [ IMMEDIATE ]. IMMEDIATE belongs to SUSPEND, so it is not its own exactly-one branch.
var failoverGroupSuspend = ssot.Attribute{
	Name:      "suspend",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	Children: []ssot.Attribute{{
		Name:      "immediate",
		Type:      ssot.AttrBool,
		CreateSql: ssot.BoolKeyword(),
	}},
}

// FailoverGroup is the source-account role and the shared SDK interface.
// The resource generator is off: the model builders, assertions and SHOW schema describe the resource the SSoT would generate.
// The two parameter assertions select nothing while the object has no parameters.
// SHOW stays here: SHOW FAILOVER GROUPS returns both roles. FailoverGroupSecondary holds the target-account CREATE and ALTER.
var FailoverGroup = ssot.Object{
	Name:           "FailoverGroup",
	NamePlural:     "FailoverGroups",
	IdType:         "AccountObjectIdentifier",
	PreviewFeature: true,
	// The description of allowed_external_volumes is extended in the schema extension.
	AttributesModification: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithAllGenerators(
			ssot.WithSchema(ssot.SchemaShow(ssot.SchemaShowManualFields(&FailoverGroupAttrs.ObjectTypes, &FailoverGroupAttrs.AllowedIntegrationTypes))),
			ssot.WithModel(
				ssot.ModelRequestType(&FailoverGroupAttrs.ObjectTypes, "[]string"),
				ssot.ModelRequestType(&FailoverGroupAttrs.AllowedAccounts, "[]string"),
			),
		),
		ssot.WithDataSource(),
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkShow(
				ssot.SdkShowOptionalInAccount(),
				ssot.SdkShowByIDNoFiltering(),
				ssot.SdkShowByIDPredicate(ssot.SdkShowByIDPredicateNameAndLocator),
			),
			// SET carries the page's mixes in one struct; its validations (the allowed lists need the matching OBJECT_TYPES value,
			// OPTIMIZED_REFRESH and RPO_ASSURANCE conflict) are in the SDK ext. The page's ADD / REMOVE rules (exactly one list), the MOVE rule
			// (at least one list; Snowflake accepts several) and IGNORE EDITION CHECK only with ALLOWED_ACCOUNTS are validated in the SDK ext as well.
			ssot.SdkAlter(
				ssot.SdkAlterSetAdditionalValidations(),
				ssot.SdkAlterAdditionalValidations(),
			),
			// The secondary CREATE (AS REPLICA OF) is emitted from FailoverGroupSecondary as CreateSecondary.
			ssot.SdkCreateForVariants(),
			// The target-account ALTER (REFRESH, PRIMARY, SUSPEND [IMMEDIATE], RESUME) is emitted from FailoverGroupSecondary as AlterSecondary.
			ssot.SdkAlterForVariants(),
			ssot.SdkCustomOperations(),
			ssot.SdkCustomInterfaceMethods(),
			// The page documents IF EXISTS and no pre-drop step. A primary group can be dropped only when no secondary group is linked to it.
			ssot.SdkDrop(),
			// The page has IF NOT EXISTS and no OR REPLACE. The failover group integration tests confirm both: a second IF NOT EXISTS
			// create leaves the first group in place, and Snowflake rejects CREATE OR REPLACE ("Failover group cannot be replaced").
			// The page's requirements (each OBJECT_TYPES value and its ALLOWED_* list go together, SECURITY INTEGRATIONS needs ROLES)
			// are validated in the SDK ext.
			ssot.SdkCreate(
				ssot.SdkCreateWithoutOrReplace(),
				ssot.SdkCreateAdditionalValidations(),
			),
		),
	},
	CustomOperations: []ssot.CustomOperationSpec{
		{
			Kind:            "ShowFailoverGroupDatabases",
			Doc:             "https://docs.snowflake.com/en/sql-reference/sql/show-databases-in-failover-group",
			QueryStructName: "ShowFailoverGroupDatabases",
			Clauses: []ssot.CustomOperationClause{
				ssot.CustomStaticSQL("SHOW"),
				ssot.CustomStaticSQL("DATABASES"),
				ssot.CustomStaticSQL("IN FAILOVER GROUP"),
				ssot.CustomName(),
			},
			Paired: &ssot.CustomPairedShow{
				DbRow: "failoverGroupDatabaseDBRow",
				Plain: "FailoverGroupDatabase",
				Columns: []*ssot.Attribute{
					{Name: "name", Type: ssot.AttrString, Readonly: true},
				},
			},
		},
		{
			Kind:            "ShowFailoverGroupShares",
			Doc:             "https://docs.snowflake.com/en/sql-reference/sql/show-shares-in-failover-group",
			QueryStructName: "ShowFailoverGroupShares",
			Clauses: []ssot.CustomOperationClause{
				ssot.CustomStaticSQL("SHOW"),
				ssot.CustomStaticSQL("SHARES"),
				ssot.CustomStaticSQL("IN FAILOVER GROUP"),
				ssot.CustomName(),
			},
			Paired: &ssot.CustomPairedShow{
				DbRow: "failoverGroupShareDBRow",
				Plain: "FailoverGroupShare",
				Columns: []*ssot.Attribute{
					{Name: "name", Type: ssot.AttrString, Readonly: true},
					{Name: "owner_account", Type: ssot.AttrString, Readonly: true},
				},
			},
		},
	},
	CustomInterfaceMethods: []ssot.CustomInterfaceMethodSpec{
		{
			Name:    "ShowDatabases",
			Doc:     "ShowDatabases returns the list of databases in the failover group as identifiers.",
			Returns: []string{"[]AccountObjectIdentifier", "error"},
		},
		{
			Name:    "ShowShares",
			Doc:     "ShowShares returns the list of shares in the failover group as identifiers.",
			Returns: []string{"[]AccountObjectIdentifier", "error"},
		},
	},
	// For now, the set/list fields have to be excluded from the generated recompute of show_output, so a change of object_types or
	// allowed_integration_types leaves show_output.object_types and show_output.allowed_integration_types stale until the next refresh.
	// TODO [SNOW-1648997]: address the above comment
	// Listed in the extension; rejects the membership changes of a group that stays in the read-only mode.
	AdditionalCustomizeDiffs: []ssot.CustomizeDiffFunc{
		"failoverGroupMembershipLockedInReadOnlyMode",
		"failoverGroupRecreateWhenNotPrimary",
	},
	// is_primary is a computed attribute, not a resource attribute, so it stays out of CREATE and ALTER.
	ComputedAttributes: []ssot.Attribute{
		FailoverGroupAttrs.IsPrimary,
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithNewCreateRequest(ssot.ExtCall()),
		ssot.WithApplyCreateOptionals(ssot.ExtCall()),
		ssot.WithCreateInSdk(ssot.Generated()),
		// CREATE has no RPO_ASSURANCE clause, so the extension sets it with ALTER after the group exists.
		ssot.WithAfterCreate(ssot.ExtCall()),
		ssot.WithShowByIdInSdk(ssot.Generated()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
		// The order of the SET and UNSET statements matters to Snowflake (RPO_ASSURANCE, REPLICATION_SCHEDULE, OPTIMIZED_REFRESH),
		// and OBJECT_TYPES with the allowed lists cannot share a SET with the other properties, so the extension builds the statements
		// and runs them in order, UNSET included.
		ssot.WithApplySetUnset(ssot.ExtCall()),
		ssot.WithAlterSetInSdk(ssot.ExtCall()),
		ssot.WithAlterUnsetInSdk(ssot.Omitted()),
		// SHOW stores identifier lists as SDK values, and the allowed databases and shares come from separate SHOW commands.
		ssot.WithSetConfigFields(ssot.ExtCall()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithHierarchyRename(ssot.Generated()),
		ssot.WithCustomizeDiff(ssot.Generated()),
		ssot.WithStateUpgraders(ssot.Generated()),
	},
	// The hand-written resource kept replication_schedule as a nested block and from_replica in the same resource.
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"v2_24_0_FailoverGroupStateUpgrader",
	},
	Attributes: []ssot.Attribute{
		FailoverGroupAttrs.ObjectTypes,
		FailoverGroupAttrs.AllowedDatabases,
		FailoverGroupAttrs.AllowedExternalVolumes,
		FailoverGroupAttrs.AllowedShares,
		FailoverGroupAttrs.AllowedIntegrationTypes,
		FailoverGroupAttrs.AllowedAccounts,
		FailoverGroupAttrs.IgnoreEditionCheck,
		FailoverGroupAttrs.ReplicationSchedule,
		FailoverGroupAttrs.OptimizedRefresh,
		FailoverGroupAttrs.RpoAssurance,
		FailoverGroupAttrs.ErrorIntegration,
		FailoverGroupAttrs.Comment,
		FailoverGroupAttrs.ReadOnly,
	},
	Create: ssot.CreateSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/create-failover-group",
		Clauses: []*ssot.Attribute{
			&FailoverGroupAttrs.ObjectTypes,
			&FailoverGroupAttrs.AllowedDatabases,
			&FailoverGroupAttrs.AllowedExternalVolumes,
			&FailoverGroupAttrs.AllowedShares,
			&FailoverGroupAttrs.AllowedIntegrationTypes,
			&FailoverGroupAttrs.AllowedAccounts,
			&FailoverGroupAttrs.IgnoreEditionCheck,
			&FailoverGroupAttrs.ReplicationSchedule,
			&FailoverGroupAttrs.OptimizedRefresh,
			ssot.Tags,
			&FailoverGroupAttrs.ErrorIntegration,
		},
	},
	Alter: ssot.AlterSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/alter-failover-group",
		// SET: allowed lists require the matching OBJECT_TYPES value. OPTIMIZED_REFRESH and RPO_ASSURANCE conflict.
		SetAdditionalValidations: true,
		Clauses: []ssot.AlterClause{
			ssot.AlterRenameTo(),
			// The page shows the SET properties as several statements: OBJECT_TYPES with the allowed lists (ALLOWED_INTEGRATION_TYPES needs
			// OBJECT_TYPES = INTEGRATIONS), then each of COMMENT, REPLICATION_SCHEDULE, OPTIMIZED_REFRESH, RPO_ASSURANCE, ERROR_INTEGRATION and
			// READ_ONLY alone (REPLICATION_SCHEDULE with OPTIMIZED_REFRESH is shown together in an example). The SDK SET struct holds all
			// of them, so the extension validates the mixes and issues the statements in order.
			ssot.AlterSet(
				&FailoverGroupAttrs.ObjectTypes,
				&FailoverGroupAttrs.AllowedDatabases,
				&FailoverGroupAttrs.AllowedExternalVolumes,
				&FailoverGroupAttrs.AllowedShares,
				&FailoverGroupAttrs.AllowedIntegrationTypes,
				&FailoverGroupAttrs.Comment,
				&FailoverGroupAttrs.ReplicationSchedule,
				&FailoverGroupAttrs.OptimizedRefresh,
				&FailoverGroupAttrs.RpoAssurance,
				&FailoverGroupAttrs.ErrorIntegration,
				&FailoverGroupAttrs.ReadOnly,
			),
			// The page allows several of these in one UNSET, separated by commas. The handwritten UNSET allowed exactly one field.
			ssot.AlterUnset(
				&FailoverGroupAttrs.Comment,
				&FailoverGroupAttrs.ReplicationSchedule,
				&FailoverGroupAttrs.OptimizedRefresh,
				&FailoverGroupAttrs.RpoAssurance,
				&FailoverGroupAttrs.ErrorIntegration,
				&FailoverGroupAttrs.ReadOnly,
			),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
			ssot.AlterKeywordBlock("Add", "ADD", &failoverGroupAdd),
			ssot.AlterKeywordBlock("Move", "MOVE", &failoverGroupMove),
			ssot.AlterKeywordBlock("Remove", "REMOVE", &failoverGroupRemove),
		},
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&FailoverGroupAttrs.RegionGroup,
			&FailoverGroupAttrs.SnowflakeRegion,
			&FailoverGroupAttrs.CreatedOn,
			&FailoverGroupAttrs.AccountName,
			&FailoverGroupAttrs.Name,
			&FailoverGroupAttrs.Type,
			&FailoverGroupAttrs.Comment,
			&FailoverGroupAttrs.IsPrimary,
			&FailoverGroupAttrs.Primary,
			&FailoverGroupAttrs.ObjectTypes,
			&FailoverGroupAttrs.AllowedIntegrationTypes,
			&FailoverGroupAttrs.AllowedAccounts,
			&FailoverGroupAttrs.OrganizationName,
			&FailoverGroupAttrs.AccountLocator,
			&FailoverGroupAttrs.ReplicationSchedule,
			&FailoverGroupAttrs.SecondaryState,
			&FailoverGroupAttrs.NextScheduledRefresh,
			&FailoverGroupAttrs.Owner,
			&FailoverGroupAttrs.IsListingAutoFulfillmentGroup,
			&FailoverGroupAttrs.IsOptimizedRefreshEnabled,
			&FailoverGroupAttrs.RpoAssurance,
			&FailoverGroupAttrs.ErrorIntegration,
			&FailoverGroupAttrs.ReadOnly,
		},
	},
}

// FailoverGroupSecondary is the target-account role.
// The name starts with FailoverGroup so the methods are CreateSecondary and AlterSecondary.
// SdkCreateRequestObjectName makes the query struct CreateSecondaryFailoverGroup. It cannot be SecondaryReplicationGroup (the
// replication group variant emits CreateSecondaryReplicationGroup in the same package); the handwritten operation was named
// CreateSecondaryReplicationGroup, with FailoverGroup appended to its types.
// Unlike a replication group, a failover group can be promoted (ALTER ... PRIMARY).
var FailoverGroupSecondary = ssot.Object{
	Name:                       "FailoverGroupSecondary",
	NamePlural:                 "FailoverGroups",
	IdType:                     "AccountObjectIdentifier",
	PreviewFeature:             true,
	IdentityNameForceNew:       true,
	SdkCreateRequestObjectName: "SecondaryFailoverGroup",
	VariantOf:                  &FailoverGroup,
	Generators: []ssot.GeneratorOption{
		ssot.WithAllGenerators(
			ssot.WithShowOutputAssert(ssot.ShowOutputShow(ssot.ShowOutputShowRowFromBase()), ssot.ShowOutputDescribe()),
			ssot.WithSchema(ssot.SchemaShow(ssot.SchemaShowFromBase()), ssot.SchemaDescribe()),
			ssot.WithObjectAssert(ssot.ObjectAssertShow(ssot.ObjectAssertShowRowFromBase()), ssot.ObjectAssertDescribe()),
		),
	},
	// is_primary is a computed attribute, not a resource attribute, so it stays out of CREATE and ALTER.
	AdditionalCustomizeDiffs: []ssot.CustomizeDiffFunc{
		"failoverGroupSecondaryRecreateWhenPrimary",
	},
	ComputedAttributes: []ssot.Attribute{
		FailoverGroupAttrs.IsPrimary,
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithNewCreateRequest(ssot.Generated()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithShowByIdInSdk(ssot.Generated()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
		ssot.WithApplySetUnset(ssot.Omitted()),
		ssot.WithAlterSetInSdk(ssot.Omitted()),
		ssot.WithAlterUnsetInSdk(ssot.Omitted()),
		ssot.WithSetConfigFields(ssot.Generated()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithCustomizeDiff(ssot.Generated()),
	},
	Attributes: []ssot.Attribute{
		failoverGroupAsReplicaOf,
	},
	// The same SHOW row as the base.
	Show: ssot.OutputSpec{Columns: FailoverGroup.Show.Columns},
	Create: ssot.CreateSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/create-failover-group",
		// The secondary CREATE page has IF NOT EXISTS and no OR REPLACE.
		WithoutOrReplace: true,
		// AsReplicaOf is a valid external identifier.
		AdditionalValidations: true,
		Clauses: []*ssot.Attribute{
			&failoverGroupAsReplicaOf,
		},
	},
	Alter: ssot.AlterSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/alter-failover-group",
		Clauses: []ssot.AlterClause{
			ssot.AlterKeyword("REFRESH"),
			ssot.AlterKeyword("PRIMARY"),
			ssot.AlterKeywordBlock("Suspend", "SUSPEND", &failoverGroupSuspend),
			ssot.AlterKeyword("RESUME"),
		},
	},
}
