package resources

import (
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var streamCommonSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the stream; must be unique for the database and schema in which the stream is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the stream."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the stream."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"copy_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: copyGrantsDescription("Retains the access permissions from the original stream when a stream is recreated using the OR REPLACE clause."),
		// Changing ONLY copy grants should have no effect. It is only used as an "option" during CREATE OR REPLACE - when other attributes change, it's not an object state. There is no point in recreating the object when only this field is changed.
		DiffSuppressFunc: IgnoreAfterCreation,
	},
	"stale": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Indicated if the stream is stale. When Terraform detects that the stream is stale, the stream is recreated with `CREATE OR REPLACE`. Read more on stream staleness in Snowflake [docs](https://docs.snowflake.com/en/user-guide/streams-intro#data-retention-period-and-staleness).",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the stream.",
	},
	"stream_type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies a type for the stream. This field is used for checking external changes and recreating the resources if needed.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW STREAMS` for the given stream.",
		Elem: &schema.Resource{
			Schema: schemas.ShowStreamSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE STREAM` for the given stream.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeStreamSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

var atSchema = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	MaxItems:    1,
	Description: externalChangesNotDetectedFieldDescription("This field specifies that the request is inclusive of any changes made by a statement or transaction with a timestamp equal to the specified parameter. Due to Snowflake limitations, the provider does not detect external changes on this field."),
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"timestamp": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies an exact date and time to use for Time Travel. The value must be explicitly cast to a TIMESTAMP, TIMESTAMP_LTZ, TIMESTAMP_NTZ, or TIMESTAMP_TZ data type.",
				ExactlyOneOf: []string{"at.0.timestamp", "at.0.offset", "at.0.statement", "at.0.stream"},
			},
			"offset": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies the difference in seconds from the current time to use for Time Travel, in the form -N where N can be an integer or arithmetic expression (e.g. -120 is 120 seconds, -30*60 is 1800 seconds or 30 minutes).",
				ExactlyOneOf: []string{"at.0.timestamp", "at.0.offset", "at.0.statement", "at.0.stream"},
			},
			"statement": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies the query ID of a statement to use as the reference point for Time Travel. This parameter supports any statement of one of the following types: DML (e.g. INSERT, UPDATE, DELETE), TCL (BEGIN, COMMIT transaction), SELECT.",
				ExactlyOneOf: []string{"at.0.timestamp", "at.0.offset", "at.0.statement", "at.0.stream"},
			},
			"stream": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Specifies the identifier (i.e. name) for an existing stream on the queried table or view. The current offset in the stream is used as the AT point in time for returning change data for the source object.",
				ExactlyOneOf:     []string{"at.0.timestamp", "at.0.offset", "at.0.statement", "at.0.stream"},
				DiffSuppressFunc: suppressIdentifierQuoting,
				ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			},
		},
	},
	ConflictsWith: []string{"before"},
}

var beforeSchema = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	MaxItems:    1,
	Description: externalChangesNotDetectedFieldDescription("This field specifies that the request refers to a point immediately preceding the specified parameter. This point in time is just before the statement, identified by its query ID, is completed.  Due to Snowflake limitations, the provider does not detect external changes on this field."),
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"timestamp": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies an exact date and time to use for Time Travel. The value must be explicitly cast to a TIMESTAMP, TIMESTAMP_LTZ, TIMESTAMP_NTZ, or TIMESTAMP_TZ data type.",
				ExactlyOneOf: []string{"before.0.timestamp", "before.0.offset", "before.0.statement", "before.0.stream"},
			},
			"offset": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies the difference in seconds from the current time to use for Time Travel, in the form -N where N can be an integer or arithmetic expression (e.g. -120 is 120 seconds, -30*60 is 1800 seconds or 30 minutes).",
				ExactlyOneOf: []string{"before.0.timestamp", "before.0.offset", "before.0.statement", "before.0.stream"},
			},
			"statement": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies the query ID of a statement to use as the reference point for Time Travel. This parameter supports any statement of one of the following types: DML (e.g. INSERT, UPDATE, DELETE), TCL (BEGIN, COMMIT transaction), SELECT.",
				ExactlyOneOf: []string{"before.0.timestamp", "before.0.offset", "before.0.statement", "before.0.stream"},
			},
			"stream": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Specifies the identifier (i.e. name) for an existing stream on the queried table or view. The current offset in the stream is used as the AT point in time for returning change data for the source object.",
				ExactlyOneOf:     []string{"before.0.timestamp", "before.0.offset", "before.0.statement", "before.0.stream"},
				DiffSuppressFunc: suppressIdentifierQuoting,
				ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			},
		},
	},
	ConflictsWith: []string{"at"},
}

func handleStreamTimeTravel(d *schema.ResourceData) *sdk.OnStreamRequest {
	if v := d.Get(AtAttributeName).([]any); len(v) > 0 {
		return sdk.NewOnStreamRequest().WithAt(true).WithStatement(handleStreamTimeTravelStatement(v[0].(map[string]any)))
	}
	if v := d.Get(BeforeAttributeName).([]any); len(v) > 0 {
		return sdk.NewOnStreamRequest().WithBefore(true).WithStatement(handleStreamTimeTravelStatement(v[0].(map[string]any)))
	}
	return nil
}

func handleStreamTimeTravelStatement(timeTravelConfig map[string]any) sdk.OnStreamStatementRequest {
	statement := sdk.OnStreamStatementRequest{}
	if v := timeTravelConfig["timestamp"].(string); len(v) > 0 {
		statement.WithTimestamp(v)
	}
	if v := timeTravelConfig["offset"].(string); len(v) > 0 {
		statement.WithOffset(v)
	}
	if v := timeTravelConfig["statement"].(string); len(v) > 0 {
		statement.WithStatement(v)
	}
	if v := timeTravelConfig["stream"].(string); len(v) > 0 {
		statement.WithStream(v)
	}
	return statement
}

var DeleteStreamContext = ResourceDeleteContextFunc(
	sdk.ParseSchemaObjectIdentifier,
	func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] { return client.Streams.DropSafely },
)

func handleStreamRead(d *schema.ResourceData,
	id sdk.SchemaObjectIdentifier,
	stream *sdk.Stream,
	streamDescription *sdk.Stream,
) error {
	return errors.Join(
		d.Set("comment", stream.Comment),
		d.Set("stream_type", stream.SourceType),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.StreamToSchema(stream)}),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.StreamDescriptionToSchema(*streamDescription)}),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set("stale", stream.Stale),
	)
}
