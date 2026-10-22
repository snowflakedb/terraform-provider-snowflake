package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// AlertAttrs holds reusable alert fields.
//
// Omitted from Attributes:
//   - alert_schedule and enabled: the nested schedule block and native optional bool are added in alertAttributesModificationExt.
//
// warehouse stays a string: the resource never suppressed identifier quoting on it.
// The resource has no show_output and no describe_output.
var AlertAttrs = struct {
	Comment   ssot.Attribute
	Warehouse ssot.Attribute
	Condition ssot.Attribute
	Action    ssot.Attribute
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
}

var Alert = ssot.Object{
	Name:                   "Alert",
	NamePlural:             "Alerts",
	IdType:                 "SchemaObjectIdentifier",
	AttributesModification: true,
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithIdentitySchema(ssot.ExtCall()),
		ssot.WithShowOutputSchema(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
	},
	Attributes: []ssot.Attribute{
		AlertAttrs.Comment,
		AlertAttrs.Warehouse,
		AlertAttrs.Condition,
		AlertAttrs.Action,
	},
}
