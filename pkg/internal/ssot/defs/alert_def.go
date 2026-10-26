package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// AlertAttrs holds reusable alert fields.
//
// alert_schedule is a nested block with SkipCreateBlock. SCHEDULE is one SQL string, and the
// ext NewCreateRequest builds it from the block. ConflictsWith paths
// stay in the schema ext because the generator emits the sibling name, and Terraform needs
// the resource path.
// enabled stays in the ext: an optional native bool. Optional bools become boolean strings.
//
// warehouse stays a string: the resource never suppressed identifier quoting on it.
// The resource has no show_output and no describe_output.
var AlertAttrs = struct {
	Comment       ssot.Attribute
	Warehouse     ssot.Attribute
	Condition     ssot.Attribute
	Action        ssot.Attribute
	AlertSchedule ssot.Attribute
}{
	Comment: ssot.Attribute{
		Name:          "comment",
		Type:          ssot.AttrString,
		Description:   "Specifies a comment for the alert.",
		CreateMapping: ssot.StringAttributeCreateBuilder,
	},
	Warehouse: ssot.Attribute{
		Name:        "warehouse",
		Type:        ssot.AttrString,
		Required:    true,
		ForceNew:    true,
		Description: "The warehouse the alert will use.",
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
		Name:            "alert_schedule",
		Type:            ssot.AttrList,
		MaxItems:        new(1),
		Description:     "The schedule for periodically running an alert.",
		SkipCreateBlock: true,
		Children: []ssot.Attribute{
			{
				Name:     "cron",
				Type:     ssot.AttrList,
				MaxItems: new(1),
				Description: "Specifies the cron expression for the alert. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". " +
					"The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
				Children: []ssot.Attribute{
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
				},
			},
			{
				Name:        "interval",
				Type:        ssot.AttrInt,
				Description: "Specifies the interval in minutes for the alert schedule. The interval must be greater than 0 and less than 1440 (24 hours).",
			},
		},
	},
}

var Alert = ssot.Object{
	Name:                   "Alert",
	NamePlural:             "Alerts",
	IdType:                 "SchemaObjectIdentifier",
	LegacyIdEncoding:       true,
	AttributesModification: true,
	// The importer is a passthrough (alertImporterExt). With ShowByIdInSdk omitted, no Import function is generated.
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithAfterCreate(ssot.ExtCall()),
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithSetStateToValuesFromConfig(ssot.Omitted()),
		ssot.WithBeforeAlter(ssot.ExtCall()),
		ssot.WithAlterUnsetInSdk(ssot.Omitted()),
		ssot.WithAfterAlter(ssot.ExtCall()),
		ssot.WithShowByIdInSdk(ssot.Omitted()),
		ssot.WithSetFieldsNotSetByRead(ssot.Omitted()),
		ssot.WithIdentitySchema(ssot.ExtCall()),
		ssot.WithShowOutputSchema(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
	},
	Attributes: []ssot.Attribute{
		AlertAttrs.Comment,
		AlertAttrs.Warehouse,
		AlertAttrs.Condition,
		AlertAttrs.Action,
		AlertAttrs.AlertSchedule,
	},
}
