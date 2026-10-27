package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// PipeAttrs holds reusable pipe fields.
//
// Omitted from Attributes and added in pipeAttributesModificationExt:
//   - owner and notification_channel: computed only. A Readonly attribute must not sit on Attributes.
//
// pipeAttributesModificationExt also sets the Default of auto_ingest and the custom DiffSuppressFunc of copy_statement.
// integration and error_integration stay strings: the resource never suppressed identifier quoting on them.
// The resource has no show_output and no describe_output.
var PipeAttrs = struct {
	Comment          ssot.Attribute
	CopyStatement    ssot.Attribute
	AutoIngest       ssot.Attribute
	AwsSnsTopicArn   ssot.Attribute
	Integration      ssot.Attribute
	ErrorIntegration ssot.Attribute
}{
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the pipe.",
	},
	CopyStatement: ssot.Attribute{
		Name:        "copy_statement",
		Type:        ssot.AttrString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the copy statement for the pipe.",
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
	},
	Integration: ssot.Attribute{
		Name:        "integration",
		Type:        ssot.AttrString,
		Description: "Specifies an integration for the pipe.",
	},
	ErrorIntegration: ssot.Attribute{
		Name:        "error_integration",
		Type:        ssot.AttrString,
		Description: "Specifies the name of the notification integration used for error notifications.",
	},
}

var Pipe = ssot.Object{
	Name:                   "Pipe",
	NamePlural:             "Pipes",
	IdType:                 "SchemaObjectIdentifier",
	LegacyIdEncoding:       true,
	AttributesModification: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
	},
	// The importer is a passthrough (pipeImporterExt). With ShowByIdInSdk omitted, no Import function is generated.
	// ApplyCreateOptionals and ApplySetUnset stay in the extension: aws_sns_topic_arn maps to WithAwsSnsTopic,
	// auto_ingest has no create mapping, and integration is not part of SET/UNSET.
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithNewCreateRequest(ssot.Generated()),
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithSetStateToValuesFromConfig(ssot.Omitted()),
		ssot.WithShowByIdInSdk(ssot.Omitted()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
		ssot.WithSetFieldsNotSetByRead(ssot.Omitted()),
		ssot.WithIdentitySchema(ssot.ExtCall()),
		ssot.WithShowOutputSchema(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
	},
	Attributes: []ssot.Attribute{
		PipeAttrs.Comment,
		PipeAttrs.CopyStatement,
		PipeAttrs.AutoIngest,
		PipeAttrs.AwsSnsTopicArn,
		PipeAttrs.Integration,
		PipeAttrs.ErrorIntegration,
	},
}
