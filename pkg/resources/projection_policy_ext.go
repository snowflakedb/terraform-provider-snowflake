package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	projectionPolicyCreateExt      = TrackingCreateWrapper(resources.ProjectionPolicy, CreateProjectionPolicy)
	projectionPolicyReadExt        = TrackingReadWrapper(resources.ProjectionPolicy, ReadProjectionPolicyFunc(true))
	projectionPolicyUpdateExt      = TrackingUpdateWrapper(resources.ProjectionPolicy, UpdateProjectionPolicy)
	projectionPolicyDescriptionExt = "Resource used to manage projection policies. For more information, check [projection policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-projection-policy)."
	projectionPolicyImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.ProjectionPolicy, importProjectionPolicy),
	}
)

type projectionPolicyAlterRequestsExt struct {
	set   *sdk.ProjectionPolicySetRequest
	unset *sdk.ProjectionPolicyUnsetRequest
}

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

// Name is not ForceNew. RENAME TO updates the id used by the later SET and UNSET.
func projectionPolicyHierarchyRenameExt(ctx context.Context, d *schema.ResourceData, meta any, id *sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if !d.HasChange("name") || d.GetRawState().IsNull() {
		return nil
	}
	newID := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), d.Get("name").(string))
	if err := meta.(*provider.Context).Client.ProjectionPolicies.Alter(ctx, sdk.NewAlterProjectionPolicyRequest(*id).WithRenameTo(newID)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(newID))
	*id = newID
	return nil
}

// SET BODY -> cannot be combined with SET COMMENT. Body is its own ALTER.
func projectionPolicyBeforeAlterExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if !d.HasChange("body") {
		return nil
	}
	request := sdk.NewProjectionPolicySetRequest().WithSetBody(d.Get("body").(string))
	if err := meta.(*provider.Context).Client.ProjectionPolicies.Alter(ctx, sdk.NewAlterProjectionPolicyRequest(id).WithSet(*request)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func projectionPolicyApplySetUnsetExt(d *schema.ResourceData) (*projectionPolicyAlterRequestsExt, error) {
	set, unset := sdk.NewProjectionPolicySetRequest(), sdk.NewProjectionPolicyUnsetRequest()
	errs := errors.Join(
		stringAttributeUpdateBuilder(d, "comment", set.WithComment, unset.WithComment),
	)
	if errs != nil {
		return nil, errs
	}
	return &projectionPolicyAlterRequestsExt{set: set, unset: unset}, nil
}
