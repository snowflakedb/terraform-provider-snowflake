package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	pipeCreateExt      = PreviewFeatureCreateContextWrapper(string(previewfeatures.PipeResource), TrackingCreateWrapper(resources.Pipe, CreatePipe))
	pipeReadExt        = PreviewFeatureReadContextWrapper(string(previewfeatures.PipeResource), TrackingReadWrapper(resources.Pipe, ReadPipeFunc(true)))
	pipeUpdateExt      = PreviewFeatureUpdateContextWrapper(string(previewfeatures.PipeResource), TrackingUpdateWrapper(resources.Pipe, UpdatePipe))
	pipeDescriptionExt = "Resource used to manage pipes. For more information, check [pipe documentation](https://docs.snowflake.com/en/sql-reference/sql/create-pipe)."
	pipeImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.Pipe, ImportPipe),
	}
)

// =============================================================================
// Request / mapping / Read-Update helpers
// =============================================================================

// name, database, and schema come from the id. integration and error_integration are not set by Read (they are diff-suppressed against show_output),
// so import writes them: the Snowflake values in the enforcing mode, and empty values in the permissive mode.
func pipeSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, pipe *sdk.Pipe, importPermissiveMode bool) error {
	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		importStringAttribute(d, "integration", importPermissiveMode, pipe.Integration),
		importStringAttribute(d, "error_integration", importPermissiveMode, pipe.ErrorIntegration),
	)
}

// auto_ingest, aws_sns_topic_arn, and the deprecated owner and notification_channel fields are derived from the notification channel and the owner.
func pipeSetDescribeOutputExt(d *schema.ResourceData, _ sdk.SchemaObjectIdentifier, details *sdk.PipeDetails, describeErr error, _ bool) error {
	if describeErr != nil {
		return describeErr
	}
	errs := errors.Join(
		d.Set("owner", details.Owner),
		d.Set("notification_channel", details.NotificationChannel),
		d.Set("auto_ingest", details.NotificationChannel != ""),
	)
	if notificationChannelIsAwsSnsTopic(details.NotificationChannel) {
		errs = errors.Join(errs, d.Set("aws_sns_topic_arn", details.NotificationChannel))
	}
	return errs
}

// Snowflake rejects setting another error integration while one is already assigned, so it is cleared first.
// The current value is read from Snowflake, because error_integration is not refreshed into the state.
func pipeBeforeAlterExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if !d.HasChange("error_integration") || d.Get("error_integration").(string) == "" {
		return nil
	}
	client := meta.(*provider.Context).Client
	pipe, err := client.Pipes.ShowByID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if pipe.ErrorIntegration == "" {
		return nil
	}
	if err := client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(id).WithUnset(*sdk.NewPipeUnsetRequest().WithErrorIntegration(true))); err != nil {
		return diag.FromErr(fmt.Errorf("error clearing pipe %v error integration: %w", id.Name(), err))
	}
	return nil
}

// CREATE PIPE has no parameters, so the configured ones are set with ALTER PIPE right after the create.
// The values are read from the configuration, not from the changes, so a configured zero value
// (for example, pipe_execution_paused = false) is set on the pipe too.
func pipeAfterCreateExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	set := sdk.NewPipeSetRequest()
	if diags := JoinDiags(
		handleParameterCreateWithMapping(d, "log_event_level", &set.LogEventLevel, sdk.ToLogLevel),
		handleParameterCreate(d, "pipe_execution_paused", &set.PipeExecutionPaused),
	); diags != nil {
		return diags
	}
	if !reflect.DeepEqual(*set, *sdk.NewPipeSetRequest()) {
		if err := meta.(*provider.Context).Client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

// =============================================================================
// Pre-existing mappers
// =============================================================================

func pipeCopyStatementDiffSuppress(_, o, n string, _ *schema.ResourceData) bool {
	// standardize line endings
	o = strings.ReplaceAll(o, "\r\n", "\n")
	n = strings.ReplaceAll(n, "\r\n", "\n")

	// trim off any trailing line endings and leading/trailing whitespace
	return strings.TrimSpace(strings.TrimRight(o, ";\r\n")) == strings.TrimSpace(strings.TrimRight(n, ";\r\n"))
}

// notificationChannelIsAwsSnsTopic reports whether a pipe notification channel is an SNS topic ARN.
// ":sns:" matches the service field in any AWS partition, including aws-iso.
func notificationChannelIsAwsSnsTopic(notificationChannel string) bool {
	return strings.Contains(notificationChannel, ":sns:")
}
