package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// AlertAttrs holds reusable alert fields.
//
// alert_schedule is a block in Terraform and one SCHEDULE string in SQL, SHOW, and DESCRIBE (SdkString).
// The resource converts between the two in alert_ext.go. cron and interval conflict.
// An alert without a schedule runs when new rows arrive. A probe on 2026-11-02 showed that Snowflake converts
// neither way: UNSET SCHEDULE fails with 392026 and SET SCHEDULE on such an alert fails with 392025.
// Adding or removing the block recreates the alert.
//
// enabled uses SkipCreate. CREATE has no ENABLED clause. Resume and suspend stay in the ext
// after-create and after-alter hooks.
//
// Without warehouse the alert is serverless. SHOW returns null for it, and SET and UNSET WAREHOUSE switch
// between the two. condition and action come back from SHOW as written.
var AlertAttrs = struct {
	Comment       ssot.Attribute
	Warehouse     ssot.Attribute
	Condition     ssot.Attribute
	Action        ssot.Attribute
	AlertSchedule ssot.Attribute
	Enabled       ssot.Attribute
	CreatedOn     ssot.Attribute
	Name          ssot.Attribute
	DatabaseName  ssot.Attribute
	SchemaName    ssot.Attribute
	Owner         ssot.Attribute
	State         ssot.Attribute
	DescribeState ssot.Attribute
	OwnerRoleType ssot.Attribute
}{
	Enabled: ssot.Attribute{
		Name:        "enabled",
		Type:        ssot.AttrBool,
		NativeBool:  true,
		SkipCreate:  true,
		Description: "Specifies if an alert should be 'started' (enabled) after creation or should remain 'suspended' (default).",
	},
	Comment: ssot.Attribute{
		Name:           "comment",
		OptionalInShow: true,
		OptionalInDb:   true,
		Type:           ssot.AttrString,
		Description:    "Specifies a comment for the alert.",
		CreateMapping:  ssot.StringAttributeCreateBuilder,
		CreateSql:      ssot.StringParameter(),
		SetSql:         ssot.StringParameter(),
		UnsetSql:       ssot.BoolKeyword(),
	},
	Warehouse: ssot.Attribute{
		Name:           "warehouse",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "AccountObjectIdentifier",
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "The warehouse the alert will use. When not set, the alert uses serverless compute resources. Creating, resuming, or switching to a serverless alert requires the EXECUTE MANAGED ALERT privilege; using a named warehouse requires its USAGE privilege.",
		CreateSql:      ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Sql("WAREHOUSE"), ssot.Equals()),
		SetSql:         ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Sql("WAREHOUSE"), ssot.Equals()),
		UnsetSql:       ssot.BoolKeyword(ssot.Sql("WAREHOUSE")),
	},
	Condition: ssot.Attribute{
		Name:         "condition",
		Type:         ssot.AttrString,
		Required:     true,
		Description:  "The SQL statement that represents the condition for the alert. (SELECT, SHOW, CALL)",
		DiffSuppress: ssot.SuppressStatement,
	},
	Action: ssot.Attribute{
		Name:         "action",
		Type:         ssot.AttrString,
		Required:     true,
		Description:  "The SQL statement that should be executed if the condition returns one or more rows.",
		DiffSuppress: ssot.SuppressStatement,
	},
	AlertSchedule: ssot.Attribute{
		Name:           "alert_schedule",
		ShowName:       "schedule",
		Type:           ssot.AttrObject,
		SdkString:      true,
		SkipCreate:     true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "The schedule for periodically running an alert. Without it, the alert runs when new rows arrive in the tables of the condition. Adding or removing the schedule recreates the alert.",
		CreateSql:      ssot.StringParameter(ssot.Sql("SCHEDULE"), ssot.FieldName("Schedule")),
		SetSql:         ssot.StringParameter(ssot.Sql("SCHEDULE"), ssot.FieldName("Schedule")),
		UnsetSql:       ssot.BoolKeyword(ssot.Sql("SCHEDULE"), ssot.FieldName("Schedule")),
		Children: []ssot.Attribute{
			{
				Name:          "cron",
				Type:          ssot.AttrList,
				MaxItems:      new(1),
				ConflictsWith: []ssot.Attribute{{Name: "interval"}},
				Description: "Specifies the cron expression for the alert. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". " +
					"The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
				Elem: ssot.ObjectElem([]ssot.Attribute{
					{
						Name:        "expression",
						Type:        ssot.AttrString,
						Required:    true,
						Description: "Specifies the cron expression for the alert. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
					},
					{
						Name:        "time_zone",
						Type:        ssot.AttrString,
						Required:    true,
						Description: "Specifies the time zone for alert refresh.",
					},
				}),
			},
			{
				Name:          "interval",
				Type:          ssot.AttrInt,
				ConflictsWith: []ssot.Attribute{{Name: "cron"}},
				Description:   "Specifies the interval in minutes for the alert schedule. The interval must be greater than 0 and less than 1440 (24 hours).",
			},
		},
	},
	CreatedOn:     ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	Name:          ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	DatabaseName:  ssot.Attribute{Name: "database_name", Type: ssot.AttrString, Readonly: true},
	SchemaName:    ssot.Attribute{Name: "schema_name", Type: ssot.AttrString, Readonly: true},
	Owner:         ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	State:         ssot.Attribute{Name: "state", Type: ssot.AttrEnum, Enum: alertStateEnum, Readonly: true},
	DescribeState: ssot.Attribute{Name: "state", Type: ssot.AttrString, Readonly: true},
	OwnerRoleType: ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, OptionalInDb: true, Readonly: true},
}

// alertStateEnum is the SHOW state column. RESUME and SUSPEND are ALTER keywords, not this enum.
var alertStateEnum = ssot.NewEnum("AlertState", "AlertStates", "started", "suspended")

// SQL clauses that are not resource attributes. condition and action are statement text.
// CREATE uses IF (EXISTS (...)) and THEN.
var (
	// One IF block, not a list. ObjectParameter is a struct tagged list, so the SQL word is printed
	// before the parentheses: IF (EXISTS (...)).
	alertConditionSql = ssot.Attribute{
		Name:      "condition",
		Type:      ssot.AttrObject,
		Required:  true,
		CreateSql: ssot.ObjectParameter(ssot.Sql("IF"), ssot.Parentheses(), ssot.Required(), ssot.FieldName("condition")),
		Children: []ssot.Attribute{{
			Name:      "condition",
			Type:      ssot.AttrList,
			Elem:      &ssot.Attribute{Type: ssot.AttrString},
			CreateSql: ssot.StringListKeyword(ssot.Parentheses(), ssot.NoComma(), ssot.Sql("EXISTS"), ssot.FieldName("Condition")),
		}},
	}
	alertActionSql = ssot.Attribute{
		Name:      "then_action",
		Type:      ssot.AttrString,
		Required:  true,
		CreateSql: ssot.PlainParameter(ssot.Sql("THEN"), ssot.NoEquals(), ssot.Required(), ssot.FieldName("action")),
	}
)

var Alert = ssot.Object{
	Name:                      "Alert",
	NamePlural:                "Alerts",
	IdType:                    "SchemaObjectIdentifier",
	PreviewFeature:            true,
	IdentityNameForceNew:      true,
	IdentityHierarchyForceNew: true,
	// The CREATE request takes the condition and action SQL clauses, not resource attributes, and the alert_schedule
	// string. Set and unset of alert_schedule, MODIFY CONDITION, and MODIFY ACTION are ext code too.
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithNewCreateRequest(ssot.ExtCall()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithCreateInSdk(ssot.Generated()),
		ssot.WithAfterCreate(ssot.ExtCall()),
		ssot.WithShowByIdInSdk(ssot.Generated()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
		ssot.WithDescribeInSdk(ssot.Generated()),
		ssot.WithSetConfigFields(ssot.ExtCall()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithBeforeAlter(ssot.ExtCall()),
		ssot.WithApplySetUnset(ssot.ExtCall()),
		ssot.WithAlterSetInSdk(ssot.ExtCall()),
		ssot.WithAlterUnsetInSdk(ssot.Generated()),
		ssot.WithAfterAlter(ssot.ExtCall()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
		ssot.WithStateUpgraders(ssot.Generated()),
		ssot.WithCustomizeDiff(ssot.Generated()),
	},
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName",
	},
	AdditionalCustomizeDiffs: []ssot.CustomizeDiffFunc{
		"alertForceNewIfScheduleAddedOrRemoved",
	},
	Generators: []ssot.GeneratorOption{
		ssot.WithAllGenerators(),
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkCreate(),
			ssot.SdkAlter(),
			ssot.SdkDrop(),
			ssot.SdkShow(
				ssot.SdkShowTerse(),
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalIn(),
				ssot.SdkShowOptionalStartsWith(),
				ssot.SdkShowOptionalLimitFrom(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDLike, ssot.SdkShowByIDIn),
			),
			ssot.SdkDescribe(),
		),
	},
	Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{
		&AlertAttrs.Warehouse,
		&AlertAttrs.AlertSchedule,
		&AlertAttrs.Comment,
		&alertConditionSql,
		&alertActionSql,
	}},
	Set: []*ssot.Attribute{
		&AlertAttrs.Warehouse,
		&AlertAttrs.AlertSchedule,
		&AlertAttrs.Comment,
	},
	Unset: []*ssot.Attribute{
		&AlertAttrs.Warehouse,
		&AlertAttrs.AlertSchedule,
		&AlertAttrs.Comment,
	},
	Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
		ssot.AlterKeyword("RESUME"),
		ssot.AlterKeyword("SUSPEND"),
		ssot.AlterSet(),
		ssot.AlterUnset(),
		ssot.AlterKeywordList("ModifyCondition", "MODIFY CONDITION EXISTS", "string"),
		ssot.AlterPlainText("ModifyAction", "MODIFY ACTION"),
	}},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&AlertAttrs.CreatedOn,
			&AlertAttrs.Name,
			&AlertAttrs.DatabaseName,
			&AlertAttrs.SchemaName,
			&AlertAttrs.Owner,
			&AlertAttrs.Comment,
			&AlertAttrs.Warehouse,
			&AlertAttrs.AlertSchedule,
			&AlertAttrs.State,
			&AlertAttrs.Condition,
			&AlertAttrs.Action,
			&AlertAttrs.OwnerRoleType,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&AlertAttrs.CreatedOn,
			&AlertAttrs.Name,
			&AlertAttrs.DatabaseName,
			&AlertAttrs.SchemaName,
			&AlertAttrs.Owner,
			&AlertAttrs.Comment,
			&AlertAttrs.Warehouse,
			&AlertAttrs.AlertSchedule,
			&AlertAttrs.DescribeState,
			&AlertAttrs.Condition,
			&AlertAttrs.Action,
		},
	},
	Attributes: []ssot.Attribute{
		AlertAttrs.Comment,
		AlertAttrs.Warehouse,
		AlertAttrs.Condition,
		AlertAttrs.Action,
		AlertAttrs.AlertSchedule,
		AlertAttrs.Enabled,
	},
}
