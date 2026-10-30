package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// Pages:
// https://docs.snowflake.com/en/sql-reference/sql/create-replication-group
// https://docs.snowflake.com/en/sql-reference/sql/alter-replication-group
// https://docs.snowflake.com/en/sql-reference/sql/drop-replication-group
// https://docs.snowflake.com/en/sql-reference/sql/show-replication-groups
//
// There is no handwritten replication group SDK definition. Failover groups
// (pkg/sdk/generator/defs/failover_groups_def.go) are the related object.
// This definition follows the replication group pages where they disagree.
//
// Database and share replication is available on every edition. Other object
// types, and failover, require Business Critical Edition or higher. Integration
// tests need that account.
//
// CREATE and the secondary CREATE have IF NOT EXISTS and no OR REPLACE.
// SdkCreate must use SdkCreateWithoutOrReplace. The SSOT default adds OR REPLACE.
// DROP documents IF EXISTS.
//
// desc-replication-group and show-external-volumes-in-replication-group return
// 404. SHOW DATABASES IN REPLICATION GROUP and SHOW SHARES IN REPLICATION GROUP
// exist and publish no columns. Those commands are not operations here.
//
// SHOW is `SHOW REPLICATION GROUPS [ IN ACCOUNT <account> ]`.
// SdkShowOptionalInAccount emits that clause. SdkShowOptionalIn would also emit
// database and schema IN. ObjectTypeReplicationGroup, its plural, and the
// singular-to-plural map already exist in pkg/sdk/object_types.go.
//
// OBJECT_TYPES and ALLOWED_INTEGRATION_TYPES are unquoted keywords. The list
// item struct is the enum type, not the quoted <Enum>ListItem helper an enum
// element also causes the generator to emit.

// ReplicationGroupObjectType is the OBJECT_TYPES list on the create and alter pages.
// It is not PluralObjectType. That type has no ACCOUNT PARAMETERS value; the failover
// converter rewrites a SHOW token PARAMETERS into that text.
var ReplicationGroupObjectType = ssot.NewEnum(
	"ReplicationGroupObjectType", "ReplicationGroupObjectTypes",
	"ACCOUNT PARAMETERS",
	"DATABASES",
	"EXTERNAL VOLUMES",
	"INTEGRATIONS",
	"NETWORK POLICIES",
	"RESOURCE MONITORS",
	"ROLES",
	"SHARES",
	"USERS",
	"WAREHOUSES",
)

// ReplicationGroupIntegrationType has the same five values as sdk.IntegrationType.
// A second IntegrationType would not compile next to failover_groups_enums_gen.go.
// The failover SHOW converter rewrites underscores and appends " INTEGRATIONS".
// The replication group page does not document those tokens.
var ReplicationGroupIntegrationType = ssot.NewEnum(
	"ReplicationGroupIntegrationType", "ReplicationGroupIntegrationTypes",
	"SECURITY INTEGRATIONS",
	"API INTEGRATIONS",
	"STORAGE INTEGRATIONS",
	"EXTERNAL ACCESS INTEGRATIONS",
	"NOTIFICATION INTEGRATIONS",
)

// The SHOW page says secondary_state is started, suspended, or NULL.
// NULL is an absent value. The failover enum is uppercase and has a NULL constant.
var ReplicationGroupSecondaryState = ssot.NewEnum(
	"ReplicationGroupSecondaryState", "ReplicationGroupSecondaryStates",
	"started", "suspended",
)

// SHOW REPLICATION GROUPS returns both group types.
var ReplicationGroupType = ssot.NewEnum(
	"ReplicationGroupType", "ReplicationGroupTypes",
	"REPLICATION", "FAILOVER",
)

func replicationGroupEnumList(enum *ssot.Enum, required bool) ssot.SqlClause {
	opts := []ssot.SqlClauseOption{ssot.NoParentheses()}
	if required {
		opts = append(opts, ssot.Required())
	}
	return ssot.ListParameter(enum.Name, opts...)
}

func replicationGroupIdList(idType string, required bool) ssot.SqlClause {
	opts := []ssot.SqlClauseOption{ssot.NoParentheses()}
	if required {
		opts = append(opts, ssot.Required())
	}
	return ssot.ListParameter(idType, opts...)
}

func replicationGroupMoveList(name, sql string) ssot.Attribute {
	return ssot.Attribute{
		Name: name,
		Type: ssot.AttrList,
		Elem: &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"},
		CreateSql: ssot.ListParameter(
			"AccountObjectIdentifier",
			ssot.Sql(sql),
			ssot.NoParentheses(),
			ssot.NoEquals(),
		),
	}
}

// ReplicationGroupAttrs holds reusable replication group fields.
var ReplicationGroupAttrs = struct {
	ObjectTypes                   ssot.Attribute
	AllowedDatabases              ssot.Attribute
	AllowedExternalVolumes        ssot.Attribute
	AllowedShares                 ssot.Attribute
	AllowedIntegrationTypes       ssot.Attribute
	AllowedAccounts               ssot.Attribute
	IgnoreEditionCheck            ssot.Attribute
	ReplicationSchedule           ssot.Attribute
	ErrorIntegration              ssot.Attribute
	Comment                       ssot.Attribute
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
	ReadOnly                      ssot.Attribute
}{
	ObjectTypes: ssot.Attribute{
		Name:          "object_types",
		Type:          ssot.AttrList,
		Required:      true,
		CreateMapping: ssot.ListValueCreate("toReplicationGroupObjectTypes"),
		UpdateMapping: ssot.ListValueUpdateSetOnly("toReplicationGroupObjectTypes"),
		Description:   "Types of objects replicated from the source account to the target accounts.",
		Elem:          &ssot.Attribute{Type: ssot.AttrEnum, Enum: ReplicationGroupObjectType},
		// CSV. NULL or empty when a secondary group lists none. The plain field stays a slice.
		// SHOW uses PARAMETERS for account parameters. That token is mapped to ACCOUNT PARAMETERS.
		// Any other token is left unchanged.
		// The SHOW schema generator does not map an enum slice, so the column stays off the schema.
		OptionalInDb: true,
		SkipInSchema: true,
		CustomParser: "ParseCommaSeparatedReplicationGroupObjectTypeArray",
		CreateSql:    replicationGroupEnumList(ReplicationGroupObjectType, true),
		SetSql:       replicationGroupEnumList(ReplicationGroupObjectType, false),
	},
	AllowedDatabases: ssot.Attribute{
		Name:        "allowed_databases",
		Type:        ssot.AttrList,
		Description: "Databases enabled for replication. OBJECT_TYPES must include DATABASES.",
		Elem:        &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"},
		CreateSql:   replicationGroupIdList("AccountObjectIdentifier", false),
		SetSql:      replicationGroupIdList("AccountObjectIdentifier", false),
	},
	AllowedExternalVolumes: ssot.Attribute{
		Name:        "allowed_external_volumes",
		Type:        ssot.AttrList,
		Description: "External volumes enabled for replication. Preview feature, open to all accounts. OBJECT_TYPES must include EXTERNAL VOLUMES.",
		Elem:        &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"},
		CreateSql:   replicationGroupIdList("AccountObjectIdentifier", false),
		SetSql:      replicationGroupIdList("AccountObjectIdentifier", false),
	},
	AllowedShares: ssot.Attribute{
		Name:        "allowed_shares",
		Type:        ssot.AttrList,
		Description: "Shares enabled for replication. OBJECT_TYPES must include SHARES.",
		Elem:        &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"},
		CreateSql:   replicationGroupIdList("AccountObjectIdentifier", false),
		SetSql:      replicationGroupIdList("AccountObjectIdentifier", false),
	},
	AllowedIntegrationTypes: ssot.Attribute{
		Name:        "allowed_integration_types",
		Type:        ssot.AttrList,
		Description: "Integration types enabled for replication. Requires Business Critical Edition or higher. OBJECT_TYPES must include INTEGRATIONS. SECURITY INTEGRATIONS also requires ROLES.",
		Elem:        &ssot.Attribute{Type: ssot.AttrEnum, Enum: ReplicationGroupIntegrationType},
		// The column is always present. It is empty when no integration types are set.
		// SHOW prints short tokens such as API and EXTERNAL_ACCESS. Those tokens are mapped to the
		// create-page values. Any other token is left unchanged.
		// The SHOW schema generator does not map an enum slice, so the column stays off the schema.
		OptionalInDb:  true,
		SkipInSchema:  true,
		UpdateMapping: ssot.ListValueUpdateSetOnly("toReplicationGroupAllowedIntegrationTypes"),
		CustomParser:  "ParseCommaSeparatedReplicationGroupIntegrationTypeArray",
		CreateSql:     replicationGroupEnumList(ReplicationGroupIntegrationType, false),
		SetSql:        replicationGroupEnumList(ReplicationGroupIntegrationType, false),
	},
	AllowedAccounts: ssot.Attribute{
		Name:          "allowed_accounts",
		Type:          ssot.AttrSet,
		Required:      true,
		CreateMapping: ssot.ListValueCreate("toReplicationGroupAllowedAccounts"),
		Description:   "Target accounts, as organization and account name, to which replication is enabled. Snowflake also includes the source account.",
		Elem:          &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountIdentifier"},
		// TODO [next PRs]: show_output.allowed_accounts stays TypeList. The schema generator follows the SDK slice, so SHOW order remains on that block. The resource attribute is the set.
		// The SHOW example prints org.account.group. The column text says accounts, and CREATE takes org.account.
		// The parser accepts org.account. A three-part value fails conversion.
		OptionalInDb: true,
		CustomParser: "ParseCommaSeparatedAccountIdentifierArray",
		CreateSql:    replicationGroupIdList("AccountIdentifier", true),
	},
	IgnoreEditionCheck: ssot.Attribute{
		Name:        "ignore_edition_check",
		Type:        ssot.AttrBool,
		CreateOnly:  true,
		Description: "Allows replication to accounts on lower editions in the cases listed on the create page.",
		CreateSql:   ssot.BoolKeyword(ssot.Sql("IGNORE EDITION CHECK")),
	},
	ReplicationSchedule: ssot.Attribute{
		Name:           "replication_schedule",
		Type:           ssot.AttrString,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Schedule for refreshing secondary replication groups. An interval (`<num> MINUTE`) or `USING CRON <expr> <time_zone>`. NULL when unset.",
	},
	ErrorIntegration: ssot.Attribute{
		Name:           "error_integration",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "AccountObjectIdentifier",
		Description:    "Notification integration that receives refresh error notifications.",
		CreateSql:      ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Equals()),
		SetSql:         ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Equals()),
	},
	Comment: ssot.Attribute{
		Name:           "comment",
		Type:           ssot.AttrString,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Comment for the replication group. The alter page defaults it to NULL. CREATE has no COMMENT clause.",
	},
	// SHOW primary is org.account.name. It is not a source-account argument.
	Primary: ssot.Attribute{
		Name:           "primary",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "ExternalObjectIdentifier",
		Readonly:       true,
		CustomParser:   "ParseExternalObjectIdentifier",
		Description:    "Primary replication group, as org.account.name. On a primary group this is the group itself. On a secondary group this is the source group.",
	},
	// The column is present only for organizations that span multiple region groups.
	// The SHOW example omits it, and omits read_only. The column list includes both.
	RegionGroup:      ssot.Attribute{Name: "region_group", Type: ssot.AttrString, Readonly: true, Description: "Region group where the account is located. Visible only to organizations that span multiple region groups."},
	SnowflakeRegion:  ssot.Attribute{Name: "snowflake_region", Type: ssot.AttrString, Readonly: true},
	CreatedOn:        ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	AccountName:      ssot.Attribute{Name: "account_name", Type: ssot.AttrString, Readonly: true},
	Name:             ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	Type:             ssot.Attribute{Name: "type", Type: ssot.AttrEnum, Enum: ReplicationGroupType, Readonly: true, Description: "REPLICATION or FAILOVER. SHOW REPLICATION GROUPS returns both."},
	IsPrimary:        ssot.Attribute{Name: "is_primary", Type: ssot.AttrBool, Readonly: true},
	OrganizationName: ssot.Attribute{Name: "organization_name", Type: ssot.AttrString, Readonly: true},
	AccountLocator:   ssot.Attribute{Name: "account_locator", Type: ssot.AttrString, Readonly: true},
	SecondaryState: ssot.Attribute{
		Name:           "secondary_state",
		Type:           ssot.AttrEnum,
		Enum:           ReplicationGroupSecondaryState,
		Readonly:       true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "State of scheduled refresh: started or suspended. NULL when no replication schedule is set.",
	},
	// The page calls this a date and time. Failover stores the same SHOW column as text.
	// A nullable time is OptionalInDb only, so the plain field stays time.Time.
	NextScheduledRefresh: ssot.Attribute{Name: "next_scheduled_refresh", Type: ssot.AttrTime, Readonly: true, OptionalInDb: true},
	Owner:                ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true, OptionalInShow: true, OptionalInDb: true, Description: "Role with OWNERSHIP on the group. NULL when the group is in a different region."},
	// Failover's handwritten SHOW row does not include these two columns. The replication group page does.
	// TRUE/FALSE is a bool, matching is_primary on the same SHOW command.
	IsListingAutoFulfillmentGroup: ssot.Attribute{Name: "is_listing_auto_fulfillment_group", Type: ssot.AttrBool, Readonly: true, Description: "TRUE when the group is used for Cross-Cloud Auto-Fulfillment."},
	ReadOnly: ssot.Attribute{
		Name:           "read_only",
		Type:           ssot.AttrBool,
		Readonly:       true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Read-only mode. TRUE or FALSE on a primary failover group. NULL on a secondary failover group and on every replication group.",
	},
}

// AS REPLICA OF is the secondary CREATE argument. It is not the SHOW column primary.
// The SHOW column is read-only. A create clause on an Attributes entry has to be in that statement's sequence,
// so the source attribute has no create clause and the primary CREATE does not gain AS REPLICA OF.
var replicationGroupAsReplicaOf = ssot.Attribute{
	Name:           "primary",
	Type:           ssot.AttrIdentifier,
	IdentifierType: "ExternalObjectIdentifier",
	ForceNew:       true,
	Description:    "Primary replication group to replicate, as org.account.name. This is the AS REPLICA OF identifier.",
	CreateSql:      ssot.IdentifierKeyword("ExternalObjectIdentifier", ssot.Sql("AS REPLICA OF"), ssot.Required(), ssot.FieldName("Primary")),
}

func replicationGroupReverseIdentifierList(name, sql, idType string) ssot.Attribute {
	return ssot.Attribute{
		Name: name,
		Type: ssot.AttrList,
		Elem: &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: idType},
		CreateSql: ssot.ListParameter(
			idType,
			ssot.Sql(sql),
			ssot.NoParentheses(),
			ssot.Reverse(),
		),
	}
}

// One MOVE statement sets one of the three lists. AdditionalValidations checks that.
var replicationGroupMove = ssot.Attribute{
	Name:      "move",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		replicationGroupMoveList("databases", "DATABASES"),
		replicationGroupMoveList("external_volumes", "EXTERNAL VOLUMES"),
		replicationGroupMoveList("shares", "SHARES"),
		{
			Name:           "to",
			Type:           ssot.AttrIdentifier,
			IdentifierType: "AccountObjectIdentifier",
			CreateSql:      ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Sql("TO REPLICATION GROUP"), ssot.Required(), ssot.FieldName("To")),
		},
	},
}

// ADD and REMOVE each set one allowed list. IGNORE EDITION CHECK belongs only to ADD accounts.
// TODO [next PRs]: drop Reverse(). One struct per allowed-list target, then a static TO ALLOWED_* or FROM ALLOWED_* clause.
// See pkg/sdk/generator/README.md#create-and-alter.
var replicationGroupAdd = ssot.Attribute{
	Name:      "add",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		replicationGroupReverseIdentifierList("allowed_databases", "TO ALLOWED_DATABASES", "AccountObjectIdentifier"),
		replicationGroupReverseIdentifierList("allowed_external_volumes", "TO ALLOWED_EXTERNAL_VOLUMES", "AccountObjectIdentifier"),
		replicationGroupReverseIdentifierList("allowed_shares", "TO ALLOWED_SHARES", "AccountObjectIdentifier"),
		replicationGroupReverseIdentifierList("allowed_accounts", "TO ALLOWED_ACCOUNTS", "AccountIdentifier"),
		{
			Name:      "ignore_edition_check",
			Type:      ssot.AttrBool,
			CreateSql: ssot.BoolKeyword(ssot.Sql("IGNORE EDITION CHECK")),
		},
	},
}

var replicationGroupRemove = ssot.Attribute{
	Name:      "remove",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		replicationGroupReverseIdentifierList("allowed_databases", "FROM ALLOWED_DATABASES", "AccountObjectIdentifier"),
		replicationGroupReverseIdentifierList("allowed_external_volumes", "FROM ALLOWED_EXTERNAL_VOLUMES", "AccountObjectIdentifier"),
		replicationGroupReverseIdentifierList("allowed_shares", "FROM ALLOWED_SHARES", "AccountObjectIdentifier"),
		replicationGroupReverseIdentifierList("allowed_accounts", "FROM ALLOWED_ACCOUNTS", "AccountIdentifier"),
	},
}

// SUSPEND [ IMMEDIATE ]. IMMEDIATE belongs to SUSPEND, so it is not its own exactly-one branch.
var replicationGroupSuspend = ssot.Attribute{
	Name:      "suspend",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	Children: []ssot.Attribute{{
		Name:      "immediate",
		Type:      ssot.AttrBool,
		CreateSql: ssot.BoolKeyword(ssot.Sql("IMMEDIATE")),
	}},
}

// ReplicationGroup is the source-account role and the shared SDK interface.
// SDK create, show, drop, and alter are on. CREATE has no OR REPLACE.
// DROP documents IF EXISTS, so SdkDrop keeps it. The page shows no pre-drop step.
// SHOW stays here: SHOW REPLICATION GROUPS returns both roles.
// ReplicationGroupSecondary holds the target-account CREATE and ALTER.
// SdkCreateForVariants emits CreateSecondary. SdkAlterForVariants emits AlterSecondary.
var ReplicationGroup = ssot.Object{
	Name:       "ReplicationGroup",
	NamePlural: "ReplicationGroups",
	IdType:     "AccountObjectIdentifier",
	Generators: []ssot.GeneratorOption{
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkCreate(
				ssot.SdkCreateWithoutOrReplace(),
				ssot.SdkCreateAdditionalValidations(),
			),
			ssot.SdkCreateForVariants(),
			ssot.SdkShow(
				ssot.SdkShowOptionalInAccount(),
				ssot.SdkShowByIDNoFiltering(),
				ssot.SdkShowByIDPredicate(ssot.SdkShowByIDPredicateNameAndLocator),
			),
			ssot.SdkDrop(),
			ssot.SdkAlter(
				ssot.SdkAlterSetAdditionalValidations(),
				ssot.SdkAlterAdditionalValidations(),
			),
			ssot.SdkAlterForVariants(),
		),
		ssot.WithObjectAssert(ssot.ObjectAssertShow()),
		ssot.WithSchema(ssot.SchemaShow()),
		ssot.WithResource(),
		ssot.WithModel(
			ssot.ModelRequestType(&ReplicationGroupAttrs.ObjectTypes, "[]string"),
			ssot.ModelRequestType(&ReplicationGroupAttrs.AllowedAccounts, "[]string"),
		),
		ssot.WithResourceAssert(),
		ssot.WithShowOutputAssert(ssot.ShowOutputShow()),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithNewCreateRequest(ssot.ExtCall()),
		ssot.WithApplyCreateOptionals(ssot.ExtCall()),
		ssot.WithCreateInSdk(ssot.Generated()),
		ssot.WithShowByIdInSdk(ssot.Generated()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
		// OBJECT_TYPES and the allowed lists cannot share a SET with COMMENT, REPLICATION_SCHEDULE,
		// or ERROR_INTEGRATION. The extension builds one request per group and issues two ALTER SET calls.
		ssot.WithApplySetUnset(ssot.ExtCall()),
		ssot.WithAlterSetInSdk(ssot.ExtCall()),
		ssot.WithAlterUnsetInSdk(ssot.Generated()),
		// SHOW stores identifier lists as SDK values. Generated SetConfigFields passes that slice through,
		// and d.Set cannot write []AccountIdentifier into a list of strings.
		ssot.WithSetConfigFields(ssot.ExtCall()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithCustomizeDiff(ssot.Generated()),
	},
	Attributes: []ssot.Attribute{
		ReplicationGroupAttrs.ObjectTypes,
		ReplicationGroupAttrs.AllowedDatabases,
		ReplicationGroupAttrs.AllowedExternalVolumes,
		ReplicationGroupAttrs.AllowedShares,
		ReplicationGroupAttrs.AllowedIntegrationTypes,
		ReplicationGroupAttrs.AllowedAccounts,
		ReplicationGroupAttrs.IgnoreEditionCheck,
		ReplicationGroupAttrs.ReplicationSchedule,
		ReplicationGroupAttrs.ErrorIntegration,
		ReplicationGroupAttrs.Comment,
	},
	Create: ssot.CreateSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/create-replication-group",
		// ALLOWED_DATABASES requires DATABASES in OBJECT_TYPES, and the same for external volumes,
		// shares, and integration types. SECURITY INTEGRATIONS also requires ROLES.
		AdditionalValidations: true,
		Clauses: []*ssot.Attribute{
			&ReplicationGroupAttrs.ObjectTypes,
			&ReplicationGroupAttrs.AllowedDatabases,
			&ReplicationGroupAttrs.AllowedExternalVolumes,
			&ReplicationGroupAttrs.AllowedShares,
			&ReplicationGroupAttrs.AllowedIntegrationTypes,
			&ReplicationGroupAttrs.AllowedAccounts,
			&ReplicationGroupAttrs.IgnoreEditionCheck,
			&ReplicationGroupAttrs.ReplicationSchedule,
			ssot.Tags,
			&ReplicationGroupAttrs.ErrorIntegration,
		},
	},
	// TODO [next PRs]: these fields are two exclusive SET groups, and SSOT emits them as one struct.
	// OBJECT_TYPES and the allowed-object lists are one statement. COMMENT, REPLICATION_SCHEDULE, and
	// ERROR_INTEGRATION are the other. Snowflake rejects a statement that mixes them. Dedicated SSOT handling
	// should keep the attributes ordinary and emit two groups. Until then, ReplicationGroupSet.additionalValidations
	// rejects a mixed SET.
	Set: []*ssot.Attribute{
		&ReplicationGroupAttrs.ObjectTypes,
		&ReplicationGroupAttrs.AllowedDatabases,
		&ReplicationGroupAttrs.AllowedExternalVolumes,
		&ReplicationGroupAttrs.AllowedShares,
		&ReplicationGroupAttrs.AllowedIntegrationTypes,
		&ReplicationGroupAttrs.Comment,
		&ReplicationGroupAttrs.ReplicationSchedule,
		&ReplicationGroupAttrs.ErrorIntegration,
	},
	// The page allows several of these in one UNSET, separated by commas.
	// Failover's handwritten UNSET allows exactly one field.
	Unset: []*ssot.Attribute{
		&ReplicationGroupAttrs.Comment,
		&ReplicationGroupAttrs.ReplicationSchedule,
		&ReplicationGroupAttrs.ErrorIntegration,
	},
	Alter: ssot.AlterSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/alter-replication-group",
		// SET: allowed lists require the matching OBJECT_TYPES value.
		// OBJECT_TYPES and the allowed lists cannot share a SET with COMMENT, REPLICATION_SCHEDULE, or ERROR_INTEGRATION.
		SetAdditionalValidations: true,
		// MOVE sets exactly one of databases, external volumes, and shares, and To is a valid identifier.
		AdditionalValidations: true,
		Clauses: []ssot.AlterClause{
			ssot.AlterRenameTo(),
			ssot.AlterSet(),
			ssot.AlterUnset(),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
			ssot.AlterKeywordBlock("Add", "ADD", &replicationGroupAdd),
			ssot.AlterKeywordBlock("Move", "MOVE", &replicationGroupMove),
			ssot.AlterKeywordBlock("Remove", "REMOVE", &replicationGroupRemove),
		},
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&ReplicationGroupAttrs.RegionGroup,
			&ReplicationGroupAttrs.SnowflakeRegion,
			&ReplicationGroupAttrs.CreatedOn,
			&ReplicationGroupAttrs.AccountName,
			&ReplicationGroupAttrs.Name,
			&ReplicationGroupAttrs.Type,
			&ReplicationGroupAttrs.Comment,
			&ReplicationGroupAttrs.IsPrimary,
			&ReplicationGroupAttrs.Primary,
			&ReplicationGroupAttrs.ObjectTypes,
			&ReplicationGroupAttrs.AllowedIntegrationTypes,
			&ReplicationGroupAttrs.AllowedAccounts,
			&ReplicationGroupAttrs.OrganizationName,
			&ReplicationGroupAttrs.AccountLocator,
			&ReplicationGroupAttrs.ReplicationSchedule,
			&ReplicationGroupAttrs.SecondaryState,
			&ReplicationGroupAttrs.NextScheduledRefresh,
			&ReplicationGroupAttrs.Owner,
			&ReplicationGroupAttrs.IsListingAutoFulfillmentGroup,
			&ReplicationGroupAttrs.ReadOnly,
		},
	},
}

// ReplicationGroupSecondary is the target-account role.
// The name starts with ReplicationGroup so the methods are CreateSecondary and AlterSecondary.
// SdkCreateRequestObjectName keeps the query struct CreateSecondaryReplicationGroup.
// CreateSecondary and AlterSecondary are emitted from this object.
var ReplicationGroupSecondary = ssot.Object{
	Name:                       "ReplicationGroupSecondary",
	NamePlural:                 "ReplicationGroups",
	IdType:                     "AccountObjectIdentifier",
	IdentityNameForceNew:       true,
	SdkCreateRequestObjectName: "SecondaryReplicationGroup",
	VariantOf:                  &ReplicationGroup,
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		ssot.WithModel(),
		ssot.WithResourceAssert(),
		ssot.WithShowOutputAssert(ssot.ShowOutputShow()),
	},
	// SHOW is the base row. The schema is ShowReplicationGroupSchema. Create and alter use the variant methods.
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithNewCreateRequest(ssot.ExtCall()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithShowByIdInSdk(ssot.Generated()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
		ssot.WithShowOutputSchema(ssot.ExtCall()),
		ssot.WithShowOutputSet(ssot.ExtCall()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
		ssot.WithApplySetUnset(ssot.Omitted()),
		ssot.WithAlterSetInSdk(ssot.Omitted()),
		ssot.WithAlterUnsetInSdk(ssot.Omitted()),
		ssot.WithSetConfigFields(ssot.Generated()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithCustomizeDiff(ssot.Generated()),
	},
	Attributes: []ssot.Attribute{
		replicationGroupAsReplicaOf,
	},
	// The same SHOW row as the base. The assertion is named after this resource.
	Show: ssot.OutputSpec{Columns: ReplicationGroup.Show.Columns},
	Create: ssot.CreateSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/create-replication-group",
		// The secondary CREATE page has IF NOT EXISTS and no OR REPLACE.
		WithoutOrReplace: true,
		// AsReplicaOf is a valid external identifier.
		AdditionalValidations: true,
		Clauses: []*ssot.Attribute{
			&replicationGroupAsReplicaOf,
		},
	},
	Alter: ssot.AlterSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/alter-replication-group",
		Clauses: []ssot.AlterClause{
			ssot.AlterKeyword("REFRESH"),
			ssot.AlterKeywordBlock("Suspend", "SUSPEND", &replicationGroupSuspend),
			ssot.AlterKeyword("RESUME"),
		},
	},
}
