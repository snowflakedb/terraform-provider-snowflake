package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// pipeRefresh is the ALTER REFRESH block. Its children are the PREFIX and MODIFIED_AFTER parameters.
var pipeRefresh = ssot.Attribute{
	Name:    "refresh",
	Type:    ssot.AttrObject,
	SqlOnly: true,
	Children: []ssot.Attribute{
		{Name: "prefix", Type: ssot.AttrString},
		{Name: "modified_after", Type: ssot.AttrString},
	},
}

// PipeAttrs holds reusable pipe fields.
//
// Omitted from Attributes and added in pipeAttributesModificationExt:
//   - owner and notification_channel: computed only. A Readonly attribute must not sit on Attributes.
//
// pipeAttributesModificationExt also sets the Default of auto_ingest and the custom DiffSuppressFunc of copy_statement.
// integration and error_integration stay strings: the resource never suppressed identifier quoting on them.
// SHOW and DESCRIBE schemas and assertions are generated. The resource hooks omit show_output and
// describe_output, so the resource schema does not expose them.
// ALTER SET and UNSET spread the parameter catalog for pipes: LOG_EVENT_LEVEL and PIPE_EXECUTION_PAUSED.
// Those stay off Attributes, so the resource schema does not gain parameter fields.
var PipeAttrs = struct {
	Comment             ssot.Attribute
	CopyStatement       ssot.Attribute
	AutoIngest          ssot.Attribute
	AwsSnsTopicArn      ssot.Attribute
	Integration         ssot.Attribute
	ErrorIntegration    ssot.Attribute
	CreatedOn           ssot.Attribute
	Name                ssot.Attribute
	DatabaseName        ssot.Attribute
	SchemaName          ssot.Attribute
	Owner               ssot.Attribute
	NotificationChannel ssot.Attribute
	Pattern             ssot.Attribute
	OwnerRoleType       ssot.Attribute
	InvalidReason       ssot.Attribute
}{
	Comment: ssot.Attribute{
		Name:         "comment",
		Type:         ssot.AttrString,
		OptionalInDb: true,
		Description:  "Specifies a comment for the pipe.",
	},
	CopyStatement: ssot.Attribute{
		Name:        "copy_statement",
		Type:        ssot.AttrString,
		ShowName:    "definition",
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the copy statement for the pipe.",
		CreateSql:   ssot.StringParameter(ssot.Sql("AS"), ssot.NoEquals(), ssot.NoQuotes(), ssot.Required(), ssot.FieldName("copyStatement")),
	},
	AutoIngest: ssot.Attribute{
		Name:        "auto_ingest",
		Type:        ssot.AttrBool,
		ForceNew:    true,
		NativeBool:  true,
		Description: "Specifies a auto_ingest param for the pipe.",
	},
	AwsSnsTopicArn: ssot.Attribute{
		Name:        "aws_sns_topic_arn",
		Type:        ssot.AttrString,
		ForceNew:    true,
		Description: "Specifies the Amazon Resource Name (ARN) for the SNS topic for your S3 bucket.",
		CreateSql:   ssot.StringParameter(ssot.Sql("AWS_SNS_TOPIC"), ssot.FieldName("AwsSnsTopicArn")),
	},
	Integration: ssot.Attribute{
		Name:         "integration",
		Type:         ssot.AttrString,
		ForceNew:     true,
		OptionalInDb: true,
		Description:  "Specifies an integration for the pipe.",
	},
	ErrorIntegration: ssot.Attribute{
		Name:         "error_integration",
		Type:         ssot.AttrString,
		OptionalInDb: true,
		Description:  "Specifies the name of the notification integration used for error notifications.",
		CreateSql:    ssot.StringParameter(ssot.NoQuotes()),
		SetSql:       ssot.StringParameter(ssot.NoQuotes()),
	},
	CreatedOn:           ssot.Attribute{Name: "created_on", Type: ssot.AttrString, Readonly: true},
	Name:                ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	DatabaseName:        ssot.Attribute{Name: "database_name", Type: ssot.AttrString, Readonly: true},
	SchemaName:          ssot.Attribute{Name: "schema_name", Type: ssot.AttrString, Readonly: true},
	Owner:               ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true, Description: "Name of the role that owns the pipe."},
	NotificationChannel: ssot.Attribute{Name: "notification_channel", Type: ssot.AttrString, OptionalInDb: true, Readonly: true, Description: "Amazon Resource Name of the Amazon SQS queue for the stage named in the DEFINITION column."},
	Pattern:             ssot.Attribute{Name: "pattern", Type: ssot.AttrString, OptionalInDb: true, Readonly: true},
	OwnerRoleType:       ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, OptionalInDb: true, Readonly: true},
	InvalidReason:       ssot.Attribute{Name: "invalid_reason", Type: ssot.AttrString, OptionalInDb: true, Readonly: true},
}

var Pipe = ssot.Object{
	Name:                      "Pipe",
	NamePlural:                "Pipes",
	IdType:                    "SchemaObjectIdentifier",
	DescribeFailRead:          true,
	PreviewFeature:            true,
	IdentityNameForceNew:      true,
	IdentityHierarchyForceNew: true,
	AttributesModification:    true,
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithNewCreateRequest(ssot.Generated()),
		ssot.WithCreateInSdk(ssot.Generated()),
		ssot.WithDescribeInSdk(ssot.Generated()),
		ssot.WithHandleExternalChanges(ssot.Generated()),
		ssot.WithSetStateToValuesFromConfig(ssot.Omitted()),
		ssot.WithShowByIdInSdk(ssot.Generated()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
		ssot.WithSetConfigFields(ssot.Generated()),
		ssot.WithSetFieldsNotSetByRead(ssot.ExtCall()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
		ssot.WithSetDescribeOutput(ssot.ExtCall()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithApplySetUnset(ssot.Generated()),
		ssot.WithBeforeAlter(ssot.ExtCall()),
		ssot.WithAfterCreate(ssot.ExtCall()),
		ssot.WithAlterSetInSdk(ssot.Generated()),
		ssot.WithAlterUnsetInSdk(ssot.Generated()),
		ssot.WithShowParametersInSdk(ssot.Generated()),
		ssot.WithParametersDetailsFromRaw(ssot.Generated()),
		ssot.WithSetParametersFields(ssot.Generated()),
		ssot.WithParametersOutputSet(ssot.Generated()),
		ssot.WithApplyParametersChanges(ssot.Generated()),
		ssot.WithParametersAttributesSchema(ssot.Generated()),
		ssot.WithParametersOutputSchema(ssot.Generated()),
		ssot.WithParametersFieldNames(ssot.Generated()),
		ssot.WithParametersCustomDiff(ssot.Generated()),
		ssot.WithStateUpgraders(ssot.Generated()),
		ssot.WithCustomizeDiff(ssot.Generated()),
	},
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName",
	},
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		ssot.WithSchema(ssot.SchemaShow(), ssot.SchemaDescribe()),
		ssot.WithModel(ssot.ModelMultiline(&PipeAttrs.CopyStatement)),
		ssot.WithResourceAssert(),
		ssot.WithShowOutputAssert(ssot.ShowOutputShow(), ssot.ShowOutputDescribe()),
		ssot.WithObjectAssert(ssot.ObjectAssertShow(), ssot.ObjectAssertDescribe()),
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkCreate(ssot.SdkCreateAdditionalValidations()),
			ssot.SdkShow(
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalIn(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDLike, ssot.SdkShowByIDIn),
				ssot.SdkShowAdditionalValidations(),
			),
			ssot.SdkDrop(),
			ssot.SdkAlter(ssot.SdkAlterSetList()),
			ssot.SdkDescribe(ssot.SdkDescribeQueryStructName("DescribePipe")),
			ssot.SdkShowParameters(),
			ssot.SdkShowParametersDetails(),
		),
		ssot.WithObjectParametersAssert(),
		ssot.WithResourceParametersAssert(),
	},
	ComputedAttributes: []ssot.Attribute{
		PipeAttrs.Owner,
		PipeAttrs.NotificationChannel,
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&PipeAttrs.CreatedOn,
			&PipeAttrs.Name,
			&PipeAttrs.DatabaseName,
			&PipeAttrs.SchemaName,
			&PipeAttrs.CopyStatement,
			&PipeAttrs.Owner,
			&PipeAttrs.NotificationChannel,
			&PipeAttrs.Comment,
			&PipeAttrs.Integration,
			&PipeAttrs.Pattern,
			&PipeAttrs.ErrorIntegration,
			&PipeAttrs.OwnerRoleType,
			&PipeAttrs.InvalidReason,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&PipeAttrs.CreatedOn,
			&PipeAttrs.Name,
			&PipeAttrs.DatabaseName,
			&PipeAttrs.SchemaName,
			&PipeAttrs.CopyStatement,
			&PipeAttrs.Owner,
			&PipeAttrs.NotificationChannel,
			&PipeAttrs.Comment,
			&PipeAttrs.Integration,
			&PipeAttrs.Pattern,
			&PipeAttrs.ErrorIntegration,
			&PipeAttrs.OwnerRoleType,
			&PipeAttrs.InvalidReason,
		},
	},
	Attributes: []ssot.Attribute{
		PipeAttrs.Comment,
		PipeAttrs.CopyStatement,
		PipeAttrs.AutoIngest,
		PipeAttrs.AwsSnsTopicArn,
		PipeAttrs.Integration,
		PipeAttrs.ErrorIntegration,
	},
	Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{
		&PipeAttrs.AutoIngest,
		&PipeAttrs.ErrorIntegration,
		&PipeAttrs.AwsSnsTopicArn,
		&PipeAttrs.Integration,
		&PipeAttrs.Comment,
		&PipeAttrs.CopyStatement,
	}},
	Set: []*ssot.Attribute{
		&PipeAttrs.ErrorIntegration,
		ssot.Parameters,
		&PipeAttrs.Comment,
	},
	Unset: []*ssot.Attribute{
		&PipeAttrs.ErrorIntegration,
		ssot.Parameters,
		&PipeAttrs.Comment,
	},
	Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
		ssot.AlterSet(),
		ssot.AlterUnset(),
		ssot.AlterSetTags(),
		ssot.AlterUnsetTags(),
		ssot.AlterKeywordBlock("Refresh", "REFRESH", &pipeRefresh),
	}},
}
