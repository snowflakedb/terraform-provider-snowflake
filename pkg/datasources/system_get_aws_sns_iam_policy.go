package datasources

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var systemGetAWSSNSIAMPolicySchema = map[string]*schema.Schema{
	"aws_sns_topic_arn": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Amazon Resource Name (ARN) of the SNS topic for your S3 bucket",
	},
	"aws_sns_topic_policy_json": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "IAM policy for Snowflake’s SQS queue to subscribe to this topic",
	},
}

func SystemGetAWSSNSIAMPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.SystemGetAWSSNSIAMPolicyDatasource), TrackingReadWrapper(datasources.SystemGetAwsSnsIamPolicy, ReadSystemGetAWSSNSIAMPolicy)),
		Schema:      systemGetAWSSNSIAMPolicySchema,
	}
}

// ReadSystemGetAWSSNSIAMPolicy implements schema.ReadFunc.
func ReadSystemGetAWSSNSIAMPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB
	awsSNSTopicArn := d.Get("aws_sns_topic_arn").(string)

	sel := snowflake.NewSystemGetAWSSNSIAMPolicyBuilder(awsSNSTopicArn).Select()
	row := snowflake.QueryRow(db, sel)
	policy, err := snowflake.ScanAWSSNSIAMPolicy(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] system_get_aws_sns_iam_policy (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(awsSNSTopicArn)
	return diag.FromErr(d.Set("aws_sns_topic_policy_json", policy.Policy))
}
