package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// AlertAttrs holds reusable alert fields.
//
// alert_schedule is a nested block with SkipCreate. SCHEDULE is one SQL string, and the
// ext NewCreateRequest builds it from the block. cron and interval conflict; the generator emits
// alert_schedule.0.cron and alert_schedule.0.interval.
// enabled uses SkipCreate. CREATE has no ENABLED clause. Resume and suspend stay in the ext
// after-create and after-alter hooks. ApplySetUnset stays ExtCall. The schema description
// stays in the ext: optional+computed fields are not a generator setting.
//
// warehouse stays a string: the resource never suppressed identifier quoting on it.
// The resource has no show_output and no describe_output.
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
	Schedule      ssot.Attribute
	State         ssot.Attribute
	DescribeState ssot.Attribute
	OwnerRoleType ssot.Attribute
}{
	Enabled: ssot.Attribute{Name: "enabled", Type: ssot.AttrBool, NativeBool: true, SkipCreate: true},
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
		Name:         "warehouse",
		Type:         ssot.AttrString,
		Required:     true,
		ForceNew:     true,
		OptionalInDb: true,
		Description:  "The warehouse the alert will use.",
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
		Name:        "alert_schedule",
		Type:        ssot.AttrList,
		MaxItems:    new(1),
		Description: "The schedule for periodically running an alert.",
		SkipCreate:  true,
		Elem: ssot.ObjectElem([]ssot.Attribute{
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
		}),
	},
	CreatedOn:    ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	Name:         ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	DatabaseName: ssot.Attribute{Name: "database_name", Type: ssot.AttrString, Readonly: true},
	SchemaName:   ssot.Attribute{Name: "schema_name", Type: ssot.AttrString, Readonly: true},
	Owner:        ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	Schedule: ssot.Attribute{
		Name:      "schedule",
		Type:      ssot.AttrString,
		Readonly:  true,
		CreateSql: ssot.StringParameter(ssot.Sql("SCHEDULE"), ssot.Required()),
		SetSql:    ssot.StringParameter(ssot.Sql("SCHEDULE")),
		UnsetSql:  ssot.BoolKeyword(ssot.Sql("SCHEDULE")),
	},
	State:         ssot.Attribute{Name: "state", Type: ssot.AttrEnum, Enum: alertStateEnum, Readonly: true},
	DescribeState: ssot.Attribute{Name: "state", Type: ssot.AttrString, Readonly: true},
	OwnerRoleType: ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, OptionalInDb: true, Readonly: true},
}

// alertStateEnum is the SHOW state column. RESUME and SUSPEND are ALTER keywords, not this enum.
var alertStateEnum = ssot.NewEnum("AlertState", "AlertStates", "started", "suspended")

// SQL clauses that are not resource attributes. The resource warehouse is a string, and condition and action
// are statement text. CREATE and ALTER use an identifier, IF (EXISTS (...)), and THEN.
var (
	alertWarehouseSql = ssot.Attribute{
		Name:           "sql_warehouse",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "AccountObjectIdentifier",
		Required:       true,
		CreateSql:      ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Sql("WAREHOUSE"), ssot.Equals(), ssot.Required()),
		SetSql:         ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Sql("WAREHOUSE"), ssot.Equals()),
		UnsetSql:       ssot.BoolKeyword(ssot.Sql("WAREHOUSE")),
	}
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
			CreateSql: ssot.StringListKeyword(ssot.Sql("EXISTS"), ssot.FieldName("Condition")),
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
	LegacyIdEncoding:          true,
	PreviewFeature:            true,
	IdentityNameForceNew:      true,
	IdentityHierarchyForceNew: true,
	AttributesModification:    true,
	// The importer is a passthrough (alertImporterExt). With ShowByIdInSdk omitted, no Import function is generated.
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithCreateInSdk(ssot.Generated()),
		ssot.WithAfterCreate(ssot.ExtCall()),
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithSetStateToValuesFromConfig(ssot.Omitted()),
		ssot.WithBeforeAlter(ssot.ExtCall()),
		ssot.WithAlterUnsetInSdk(ssot.Omitted()),
		ssot.WithAfterAlter(ssot.ExtCall()),
		ssot.WithShowByIdInSdk(ssot.Omitted()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
		ssot.WithSetFieldsNotSetByRead(ssot.Omitted()),
		// Schemas and output assertions are generated. The resource has no show_output or describe_output.
		ssot.WithShowOutputSchema(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
	},
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		ssot.WithModel(),
		ssot.WithResourceAssert(),
		ssot.WithObjectAssert(ssot.ObjectAssertShow(), ssot.ObjectAssertDescribe()),
		ssot.WithShowOutputAssert(ssot.ShowOutputShow(), ssot.ShowOutputDescribe()),
		ssot.WithSchema(ssot.SchemaShow(), ssot.SchemaDescribe()),
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
		&alertWarehouseSql,
		&AlertAttrs.Schedule,
		&AlertAttrs.Comment,
		&alertConditionSql,
		&alertActionSql,
	}},
	Set: []*ssot.Attribute{
		&alertWarehouseSql,
		&AlertAttrs.Schedule,
		&AlertAttrs.Comment,
	},
	Unset: []*ssot.Attribute{
		&alertWarehouseSql,
		&AlertAttrs.Schedule,
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
			&AlertAttrs.Schedule,
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
			&AlertAttrs.Schedule,
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
