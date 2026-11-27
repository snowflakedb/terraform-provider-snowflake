package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// DataMovementRuleType is the movement type that a data movement rule controls.
var DataMovementRuleType = ssot.NewEnum(
	"DataMovementRuleType", "DataMovementRuleTypes",
	"COPY_INTO_EXTERNAL_STAGE", "COPY_INTO_INTERNAL_STAGE", "SNOWSIGHT_UI", "AGENT_ACCESS", "UI_DOWNLOAD", "PROGRAMMATIC_FETCH",
)

// DataMovementRuleAttrs holds reusable data movement rule fields.
//
// Inconsistencies with the docs, checked against Snowflake on 2026-11-26:
//   - SHOW and DESCRIBE return type, and the pages list movement_type. DESCRIBE returns max_rows_signature and
//     max_rows_body, and the page lists function_signature, function_body, and function_return_type.
//   - DESCRIBE returns database_name, schema_name, and owner, and the page does not list them.
//   - SHOW accepts STARTS WITH and LIMIT ... FROM, and the page lists neither. STARTS WITH and FROM are ignored, so the
//     SDK does not offer STARTS WITH.
//   - DROP IF EXISTS succeeds, and the page does not show it.
//   - SET MAX_ROWS and SET COMMENT are rejected in one ALTER SET, so max_rows is its own alter group.
var DataMovementRuleAttrs = struct {
	Type             ssot.Attribute
	MaxRows          ssot.Attribute
	Comment          ssot.Attribute
	CreatedOn        ssot.Attribute
	Name             ssot.Attribute
	DatabaseName     ssot.Attribute
	SchemaName       ssot.Attribute
	Owner            ssot.Attribute
	OwnerRoleType    ssot.Attribute
	MaxRowsSignature ssot.Attribute
}{
	Type: ssot.Attribute{
		Name:        "type",
		Type:        ssot.AttrEnum,
		Enum:        DataMovementRuleType,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the movement type that the rule controls. A rule controls exactly one movement type.",
		CreateSql:   ssot.Clause(ssot.SingleQuotes()),
	},
	MaxRows: ssot.Attribute{
		Name:         "max_rows",
		ShowName:     "max_rows_body",
		Type:         ssot.AttrString,
		Required:     true,
		Description:  "SQL expression that returns an INTEGER, which is the maximum number of rows that the movement type can move. `NULL` means no limit, `0` blocks the movement, and a positive integer is the limit. The expression can use CASE and SYS_CONTEXT. Snowflake returns the expression as it was sent, so write it the way Snowflake should store it, for example `(1000)`.",
		DiffSuppress: ssot.SuppressStatement,
		CreateSql:    ssot.StringParameter(ssot.Sql("MAX_ROWS AS () RETURNS INTEGER ->"), ssot.NoEquals(), ssot.NoQuotes(), ssot.Required(), ssot.FieldName("MaxRows")),
		SetSql:       ssot.StringParameter(ssot.Sql("MAX_ROWS AS () RETURNS INTEGER ->"), ssot.NoEquals(), ssot.NoQuotes(), ssot.FieldName("SetMaxRows")),
	},
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the data movement rule.",
	},
	CreatedOn:        ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	Name:             ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	DatabaseName:     ssot.Attribute{Name: "database_name", Type: ssot.AttrString, Readonly: true},
	SchemaName:       ssot.Attribute{Name: "schema_name", Type: ssot.AttrString, Readonly: true},
	Owner:            ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	OwnerRoleType:    ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, Readonly: true},
	MaxRowsSignature: ssot.Attribute{Name: "max_rows_signature", Type: ssot.AttrString, Readonly: true},
}

var DataMovementRule = ssot.Object{
	Name:             "DataMovementRule",
	NamePlural:       "DataMovementRules",
	DescribeFailRead: true,
	IdType:           "SchemaObjectIdentifier",
	PreviewFeature:   true,
	Generators: []ssot.GeneratorOption{
		ssot.WithObjectParametersAssert(),
		ssot.WithResourceParametersAssert(),
		ssot.WithObjectAssert(
			ssot.ObjectAssertShow(),
			ssot.ObjectAssertDescribe(),
		),
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkCreate(),
			ssot.SdkAlter(),
			ssot.SdkDrop(),
			ssot.SdkShow(
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalIn(),
				ssot.SdkShowOptionalLimit(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDLike, ssot.SdkShowByIDIn),
			),
			ssot.SdkDescribe(
				ssot.SdkDescribeDoc("https://docs.snowflake.com/en/sql-reference/sql/desc-data-movement-rule"),
			),
		),
	},
	Attributes: []ssot.Attribute{
		DataMovementRuleAttrs.Type,
		DataMovementRuleAttrs.MaxRows,
		DataMovementRuleAttrs.Comment,
	},
	Create: ssot.CreateSpec{
		Clauses: []*ssot.Attribute{
			&DataMovementRuleAttrs.Type,
			&DataMovementRuleAttrs.MaxRows,
			&DataMovementRuleAttrs.Comment,
		},
	},
	Alter: ssot.AlterSpec{
		Clauses: []ssot.AlterClause{
			ssot.AlterRenameTo(),
			ssot.AlterSet(
				ssot.AlterGroup("MaxRows", &DataMovementRuleAttrs.MaxRows),
				ssot.AlterGroup(ssot.MainAlterGroup, &DataMovementRuleAttrs.Comment),
			),
			ssot.AlterUnset(&DataMovementRuleAttrs.Comment),
		},
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&DataMovementRuleAttrs.CreatedOn,
			&DataMovementRuleAttrs.Name,
			&DataMovementRuleAttrs.DatabaseName,
			&DataMovementRuleAttrs.SchemaName,
			&DataMovementRuleAttrs.Comment,
			&DataMovementRuleAttrs.Type,
			&DataMovementRuleAttrs.Owner,
			&DataMovementRuleAttrs.OwnerRoleType,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&DataMovementRuleAttrs.CreatedOn,
			&DataMovementRuleAttrs.Name,
			&DataMovementRuleAttrs.DatabaseName,
			&DataMovementRuleAttrs.SchemaName,
			&DataMovementRuleAttrs.Owner,
			&DataMovementRuleAttrs.Comment,
			&DataMovementRuleAttrs.Type,
			&DataMovementRuleAttrs.MaxRowsSignature,
			&DataMovementRuleAttrs.MaxRows,
		},
	},
}
