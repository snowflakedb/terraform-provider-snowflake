//go:build non_account_level_tests

package testacc

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_ProjectionPolicy_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	const (
		bodyDeny  = "PROJECTION_CONSTRAINT(ALLOW => false)"
		bodyAllow = "PROJECTION_CONSTRAINT(ALLOW => true)"
	)

	// body is updated in place with its own SET. comment is SET/UNSET and cannot share that SET.
	basic := model.ProjectionPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), bodyDeny)
	withOptionals := model.ProjectionPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), bodyAllow).
		WithComment(comment)
	commentUnset := model.ProjectionPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), bodyAllow)
	ref := basic.ResourceReference()

	assertBasic := projectionPolicyAssertions(t, ref, id, bodyDeny, "")
	assertWithOptionals := projectionPolicyAssertions(t, ref, id, bodyAllow, comment)
	assertCommentUnset := projectionPolicyAssertions(t, ref, id, bodyAllow, "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ProjectionPolicy),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedProjectionPolicyResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasBodyString(bodyDeny).
						HasCommentEmpty(),
					resourceshowoutputassert.ImportedProjectionPolicyShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasKind("PROJECTION_POLICY").
						HasComment(""),
					resourceshowoutputassert.ImportedProjectionPolicyDescribeOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasSignature("()").
						HasReturnType("PROJECTION_CONSTRAINT").
						HasBody(bodyDeny),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "body", tfjson.ActionUpdate, sdk.String(bodyDeny), sdk.String(bodyAllow)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, nil, sdk.String(comment)),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			{
				Config:            config.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedProjectionPolicyResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasBodyString(bodyAllow).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedProjectionPolicyShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasComment(comment),
					resourceshowoutputassert.ImportedProjectionPolicyDescribeOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasBody(bodyAllow),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
						planchecks.ExpectNoChangeOnField(ref, "body"),
					},
				},
				Config: config.FromModels(t, commentUnset),
				Check:  assertThat(t, assertCommentUnset...),
			},
			{
				Config:  config.FromModels(t, commentUnset),
				Destroy: true,
			},
			{
				PreConfig: func() {
					_, err := testClient().ProjectionPolicy.Show(t, id)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			{
				Config:            config.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					testClient().ProjectionPolicy.CreateWithRequest(t, sdk.NewCreateProjectionPolicyRequest(id, bodyAllow).
						WithOrReplace(true).
						WithComment(externalComment))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
						planchecks.ExpectNoChangeOnField(ref, "body"),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
		},
	})
}

func TestAcc_ProjectionPolicy_Rename(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	renamedId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), testClient().Ids.Alpha())

	const body = "PROJECTION_CONSTRAINT(ALLOW => false)"

	basic := model.ProjectionPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), body)
	renamed := model.ProjectionPolicy("t", id.DatabaseName(), id.SchemaName(), renamedId.Name(), body)
	ref := basic.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ProjectionPolicy),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, projectionPolicyAssertions(t, ref, id, body, "")...),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, renamed),
				Check:  assertThat(t, projectionPolicyAssertions(t, ref, renamedId, body, "")...),
			},
		},
	})
}

func projectionPolicyAssertions(t *testing.T, ref string, id sdk.SchemaObjectIdentifier, body string, comment string) []assert.TestCheckFuncProvider {
	t.Helper()
	resourceCheck := resourceassert.ProjectionPolicyResource(t, ref).
		HasNameString(id.Name()).
		HasDatabaseString(id.DatabaseName()).
		HasSchemaString(id.SchemaName()).
		HasBodyString(body).
		HasFullyQualifiedNameString(id.FullyQualifiedName())
	showCheck := resourceshowoutputassert.ProjectionPolicyShowOutput(t, ref).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasKind("PROJECTION_POLICY").
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE").
		HasOptions("").
		HasComment(comment)
	describeCheck := resourceshowoutputassert.ProjectionPolicyDescribeOutput(t, ref).
		HasName(id.Name()).
		HasSignature("()").
		HasReturnType("PROJECTION_CONSTRAINT").
		HasBody(body)
	showObject := objectassert.ProjectionPolicy(t, id).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasKind("PROJECTION_POLICY").
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE").
		HasOptions("").
		HasComment(comment)
	details := objectassert.ProjectionPolicyDetails(t, id).
		HasName(id.Name()).
		HasSignature("()").
		HasReturnType("PROJECTION_CONSTRAINT").
		HasBody(body)
	if comment == "" {
		resourceCheck = resourceCheck.HasCommentEmpty()
	} else {
		resourceCheck = resourceCheck.HasCommentString(comment)
	}
	return []assert.TestCheckFuncProvider{resourceCheck, showCheck, describeCheck, showObject, details}
}
