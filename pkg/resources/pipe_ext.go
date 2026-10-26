package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
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
	pipeCreateExt = PreviewFeatureCreateContextWrapper(string(previewfeatures.PipeResource), TrackingCreateWrapper(resources.Pipe, CreatePipe))
	pipeReadExt   = PreviewFeatureReadContextWrapper(string(previewfeatures.PipeResource), TrackingReadWrapper(resources.Pipe, ReadPipeFunc(true)))
	pipeUpdateExt = PreviewFeatureUpdateContextWrapper(string(previewfeatures.PipeResource), TrackingUpdateWrapper(resources.Pipe, UpdatePipe))
	pipeDeleteExt = PreviewFeatureDeleteContextWrapper(string(previewfeatures.PipeResource), TrackingDeleteWrapper(resources.Pipe, ResourceDeleteContextFunc(
		pipeParseIdExt,
		pipeDropSafelyInSdkExt,
	)))
	pipeDescriptionExt   = ""
	pipeCustomizeDiffExt schema.CustomizeDiffFunc
	pipeImporterExt      = &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}
)

// =============================================================================
// Request / mapping / Read-Update helpers
// =============================================================================

type pipeAlterRequestsExt struct {
	set   *sdk.PipeSetRequest
	unset *sdk.PipeUnsetRequest
}

func pipeParseIdExt(id string) (sdk.SchemaObjectIdentifier, error) {
	return helpers.DecodeSnowflakeIDErrLegacy[sdk.SchemaObjectIdentifier](id)
}

func pipeApplyCreateOptionalsExt(d *schema.ResourceData, req *sdk.CreatePipeRequest) error {
	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if v, ok := d.GetOk("auto_ingest"); ok && v.(bool) {
		req.WithAutoIngest(true)
	}

	if v, ok := d.GetOk("aws_sns_topic_arn"); ok {
		req.WithAwsSnsTopic(v.(string))
	}

	if v, ok := d.GetOk("integration"); ok {
		req.WithIntegration(v.(string))
	}

	if v, ok := d.GetOk("error_integration"); ok {
		req.WithErrorIntegration(v.(string))
	}
	return nil
}

func pipeSetConfigFieldsExt(d *schema.ResourceData, pipe *sdk.Pipe) error {
	errs := errors.Join(
		d.Set("name", pipe.Name),
		d.Set("database", pipe.DatabaseName),
		d.Set("schema", pipe.SchemaName),
		d.Set("copy_statement", pipe.Definition),
		d.Set("owner", pipe.Owner),
		d.Set("comment", pipe.Comment),
		d.Set("notification_channel", pipe.NotificationChannel),
		d.Set("auto_ingest", pipe.NotificationChannel != ""),
		d.Set("error_integration", pipe.ErrorIntegration),
	)
	if errs != nil {
		return errs
	}
	if notificationChannelIsAwsSnsTopic(pipe.NotificationChannel) {
		return d.Set("aws_sns_topic_arn", pipe.NotificationChannel)
	}
	return nil
}

func pipeApplySetUnsetExt(d *schema.ResourceData) (*pipeAlterRequestsExt, error) {
	reqs := &pipeAlterRequestsExt{set: sdk.NewPipeSetRequest(), unset: sdk.NewPipeUnsetRequest()}

	if d.HasChange("comment") {
		if comment, ok := d.GetOk("comment"); ok {
			reqs.set.WithComment(comment.(string))
		} else {
			reqs.unset.WithComment(true)
		}
	}

	if d.HasChange("error_integration") {
		if errorIntegration, ok := d.GetOk("error_integration"); ok {
			reqs.set.WithErrorIntegration(errorIntegration.(string))
		} else {
			reqs.unset.WithErrorIntegration(true)
		}
	}
	return reqs, nil
}

// =============================================================================
// SDK incision points
// =============================================================================

func pipeCreateInSdkExt(ctx context.Context, meta any, _ sdk.SchemaObjectIdentifier, req *sdk.CreatePipeRequest) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.Pipes.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func pipeAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier, reqs *pipeAlterRequestsExt) error {
	client := meta.(*provider.Context).Client
	if !reflect.DeepEqual(reqs.set, sdk.NewPipeSetRequest()) {
		if err := client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(id).WithSet(*reqs.set)); err != nil {
			return fmt.Errorf("error updating pipe %v: %w", id.Name(), err)
		}
	}
	return nil
}

func pipeAlterUnsetInSdkExt(ctx context.Context, meta any, id sdk.SchemaObjectIdentifier, reqs *pipeAlterRequestsExt) error {
	client := meta.(*provider.Context).Client
	if !reflect.DeepEqual(reqs.unset, sdk.NewPipeUnsetRequest()) {
		if err := client.Pipes.Alter(ctx, sdk.NewAlterPipeRequest(id).WithUnset(*reqs.unset)); err != nil {
			return fmt.Errorf("error updating pipe %v: %w", id.Name(), err)
		}
	}
	return nil
}

func pipeDropSafelyInSdkExt(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
	return client.Pipes.DropSafely
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
