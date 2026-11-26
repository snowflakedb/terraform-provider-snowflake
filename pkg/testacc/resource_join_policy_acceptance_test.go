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

func TestAcc_JoinPolicy_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	const (
		bodyRequired    = "JOIN_CONSTRAINT(JOIN_REQUIRED => TRUE)"
		bodyNotRequired = "JOIN_CONSTRAINT(JOIN_REQUIRED => FALSE)"
	)

	// body is updated in place with its own SET. comment is SET/UNSET and cannot share that SET.
	basic := model.JoinPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), bodyRequired)
	withOptionals := model.JoinPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), bodyNotRequired).
		WithComment(comment)
	commentUnset := model.JoinPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), bodyNotRequired)
	ref := basic.ResourceReference()

	assertBasic := joinPolicyAssertions(t, ref, id, bodyRequired, "")
	assertWithOptionals := joinPolicyAssertions(t, ref, id, bodyNotRequired, comment)
	assertCommentUnset := joinPolicyAssertions(t, ref, id, bodyNotRequired, "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.JoinPolicy),
		Steps: []resource.TestStep{
			// create with required fields only
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// import with required fields only
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedJoinPolicyResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasBodyString(bodyRequired).
						HasCommentEmpty(),
					resourceshowoutputassert.ImportedJoinPolicyShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasKind("JOIN_POLICY").
						HasComment(""),
					resourceshowoutputassert.ImportedJoinPolicyDescribeOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasSignature("()").
						HasReturnType("JOIN_CONSTRAINT").
						HasBody(bodyRequired),
				),
			},
			// set body and comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "body", tfjson.ActionUpdate, sdk.String(bodyRequired), sdk.String(bodyNotRequired)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, nil, sdk.String(comment)),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// import with optionals set
			{
				Config:            config.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedJoinPolicyResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasBodyString(bodyNotRequired).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedJoinPolicyShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasComment(comment),
					resourceshowoutputassert.ImportedJoinPolicyDescribeOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasBody(bodyNotRequired),
				),
			},
			// unset comment
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
			// destroy
			{
				Config:  config.FromModels(t, commentUnset),
				Destroy: true,
			},
			// create with all fields
			{
				PreConfig: func() {
					_, err := testClient().JoinPolicy.Show(t, id)
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
			// import with all fields
			{
				Config:            config.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// external changes of body and comment are reverted
			{
				PreConfig: func() {
					testClient().JoinPolicy.Alter(t, sdk.NewAlterJoinPolicyRequest(id).
						WithSet(*sdk.NewJoinPolicySetRequest().WithSetBody(bodyRequired)))
					testClient().JoinPolicy.Alter(t, sdk.NewAlterJoinPolicyRequest(id).
						WithSet(*sdk.NewJoinPolicySetRequest().WithComment(externalComment)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "body", sdk.String(bodyNotRequired), sdk.String(bodyRequired)),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "body", tfjson.ActionUpdate, sdk.String(bodyRequired), sdk.String(bodyNotRequired)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
		},
	})
}

func TestAcc_JoinPolicy_Rename(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	renamedId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), testClient().Ids.Alpha())

	const body = "JOIN_CONSTRAINT(JOIN_REQUIRED => TRUE)"

	basic := model.JoinPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), body)
	renamed := model.JoinPolicy("t", id.DatabaseName(), id.SchemaName(), renamedId.Name(), body)
	ref := basic.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.JoinPolicy),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, joinPolicyAssertions(t, ref, id, body, "")...),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, renamed),
				Check:  assertThat(t, joinPolicyAssertions(t, ref, renamedId, body, "")...),
			},
		},
	})
}

func joinPolicyAssertions(t *testing.T, ref string, id sdk.SchemaObjectIdentifier, body string, comment string) []assert.TestCheckFuncProvider {
	t.Helper()
	resourceCheck := resourceassert.JoinPolicyResource(t, ref).
		HasNameString(id.Name()).
		HasDatabaseString(id.DatabaseName()).
		HasSchemaString(id.SchemaName()).
		HasBodyString(body).
		HasFullyQualifiedNameString(id.FullyQualifiedName())
	showCheck := resourceshowoutputassert.JoinPolicyShowOutput(t, ref).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasKind("JOIN_POLICY").
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE").
		HasOptions("").
		HasComment(comment)
	describeCheck := resourceshowoutputassert.JoinPolicyDescribeOutput(t, ref).
		HasName(id.Name()).
		HasSignature("()").
		HasReturnType("JOIN_CONSTRAINT").
		HasBody(body)
	showObject := objectassert.JoinPolicy(t, id).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasKind("JOIN_POLICY").
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE").
		HasOptions("").
		HasComment(comment)
	details := objectassert.JoinPolicyDetails(t, id).
		HasName(id.Name()).
		HasSignature("()").
		HasReturnType("JOIN_CONSTRAINT").
		HasBody(body)
	if comment == "" {
		resourceCheck = resourceCheck.HasCommentEmpty()
	} else {
		resourceCheck = resourceCheck.HasCommentString(comment)
	}
	return []assert.TestCheckFuncProvider{resourceCheck, showCheck, describeCheck, showObject, details}
}
