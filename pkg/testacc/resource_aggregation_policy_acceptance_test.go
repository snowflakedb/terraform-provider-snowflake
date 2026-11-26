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

func TestAcc_AggregationPolicy_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	const (
		bodyMinGroupSize5 = "AGGREGATION_CONSTRAINT(MIN_GROUP_SIZE => 5)"
		bodyConditional   = "CASE WHEN CURRENT_ROLE() = 'ADMIN' THEN NO_AGGREGATION_CONSTRAINT() ELSE AGGREGATION_CONSTRAINT(MIN_GROUP_SIZE => 2) END"
	)

	// body is updated in place with its own SET. comment is SET/UNSET and cannot share that SET.
	basic := model.AggregationPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), bodyMinGroupSize5)
	withOptionals := model.AggregationPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), bodyConditional).
		WithComment(comment)
	commentUnset := model.AggregationPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), bodyConditional)
	ref := basic.ResourceReference()

	assertBasic := aggregationPolicyAssertions(t, ref, id, bodyMinGroupSize5, "")
	assertWithOptionals := aggregationPolicyAssertions(t, ref, id, bodyConditional, comment)
	assertCommentUnset := aggregationPolicyAssertions(t, ref, id, bodyConditional, "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AggregationPolicy),
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
					resourceassert.ImportedAggregationPolicyResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasBodyString(bodyMinGroupSize5).
						HasCommentEmpty(),
					resourceshowoutputassert.ImportedAggregationPolicyShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasKind("AGGREGATION_POLICY").
						HasComment(""),
					resourceshowoutputassert.ImportedAggregationPolicyDescribeOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasSignature("()").
						HasReturnType("AGGREGATION_CONSTRAINT").
						HasBody(bodyMinGroupSize5),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "body", tfjson.ActionUpdate, sdk.String(bodyMinGroupSize5), sdk.String(bodyConditional)),
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
					resourceassert.ImportedAggregationPolicyResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasBodyString(bodyConditional).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedAggregationPolicyShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasComment(comment),
					resourceshowoutputassert.ImportedAggregationPolicyDescribeOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasBody(bodyConditional),
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
					_, err := testClient().AggregationPolicy.Show(t, id)
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
					testClient().AggregationPolicy.CreateWithRequest(t, sdk.NewCreateAggregationPolicyRequest(id, bodyMinGroupSize5).
						WithOrReplace(true).
						WithComment(externalComment))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
						planchecks.ExpectDrift(ref, "body", sdk.String(bodyConditional), sdk.String(bodyMinGroupSize5)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
						planchecks.ExpectChange(ref, "body", tfjson.ActionUpdate, sdk.String(bodyMinGroupSize5), sdk.String(bodyConditional)),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
		},
	})
}

func aggregationPolicyAssertions(t *testing.T, ref string, id sdk.SchemaObjectIdentifier, body string, comment string) []assert.TestCheckFuncProvider {
	t.Helper()
	resourceCheck := resourceassert.AggregationPolicyResource(t, ref).
		HasNameString(id.Name()).
		HasDatabaseString(id.DatabaseName()).
		HasSchemaString(id.SchemaName()).
		HasBodyString(body).
		HasFullyQualifiedNameString(id.FullyQualifiedName())
	showCheck := resourceshowoutputassert.AggregationPolicyShowOutput(t, ref).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasKind("AGGREGATION_POLICY").
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE").
		HasOptions("").
		HasComment(comment)
	describeCheck := resourceshowoutputassert.AggregationPolicyDescribeOutput(t, ref).
		HasName(id.Name()).
		HasSignature("()").
		HasReturnType("AGGREGATION_CONSTRAINT").
		HasBody(body)
	showObject := objectassert.AggregationPolicy(t, id).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasKind("AGGREGATION_POLICY").
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE").
		HasOptions("").
		HasComment(comment)
	details := objectassert.AggregationPolicyDetails(t, id).
		HasName(id.Name()).
		HasSignature("()").
		HasReturnType("AGGREGATION_CONSTRAINT").
		HasBody(body)
	if comment == "" {
		resourceCheck = resourceCheck.HasCommentEmpty()
	} else {
		resourceCheck = resourceCheck.HasCommentString(comment)
	}
	return []assert.TestCheckFuncProvider{resourceCheck, showCheck, describeCheck, showObject, details}
}
