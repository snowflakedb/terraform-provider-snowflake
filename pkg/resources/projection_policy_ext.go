package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	projectionPolicyCreateExt      = PreviewFeatureCreateContextWrapper(string(previewfeatures.ProjectionPolicyResource), TrackingCreateWrapper(resources.ProjectionPolicy, CreateProjectionPolicy))
	projectionPolicyReadExt        = PreviewFeatureReadContextWrapper(string(previewfeatures.ProjectionPolicyResource), TrackingReadWrapper(resources.ProjectionPolicy, ReadProjectionPolicyFunc(true)))
	projectionPolicyUpdateExt      = PreviewFeatureUpdateContextWrapper(string(previewfeatures.ProjectionPolicyResource), TrackingUpdateWrapper(resources.ProjectionPolicy, UpdateProjectionPolicy))
	projectionPolicyDescriptionExt = "Resource used to manage projection policies. For more information, check [projection policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-projection-policy)."
	projectionPolicyImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.ProjectionPolicy, importProjectionPolicy),
	}
)

// SHOW does not return body. Import reads it from DESCRIBE after the generated show import.
func importProjectionPolicy(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	data, err := ImportProjectionPolicy(ctx, d, meta)
	if err != nil {
		return nil, err
	}
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	details, err := meta.(*provider.Context).Client.ProjectionPolicies.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := d.Set("body", details.Body); err != nil {
		return nil, err
	}
	return data, nil
}

func projectionPolicySetDescribeOutputExt(d *schema.ResourceData, _ sdk.SchemaObjectIdentifier, details *sdk.ProjectionPolicyDetails, describeErr error, _ bool) error {
	if describeErr != nil {
		return describeErr
	}
	return d.Set("body", details.Body)
}

func projectionPolicySetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, projectionPolicy *sdk.ProjectionPolicy, _ bool) error {
	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		d.Set("comment", projectionPolicy.Comment),
	)
}
