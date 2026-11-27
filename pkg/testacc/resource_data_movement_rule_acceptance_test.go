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

func TestAcc_DataMovementRule_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	renamedId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), testClient().Ids.Alpha())
	comment := random.Comment()
	externalComment := random.Comment()

	const (
		movementType         = sdk.DataMovementRuleTypeCopyIntoExternalStage
		internalMovementType = sdk.DataMovementRuleTypeCopyIntoInternalStage
		snowsightType        = sdk.DataMovementRuleTypeSnowsightUi
		maxRows              = "(100)"
		otherMaxRows         = "CASE WHEN SYS_CONTEXT('SNOWFLAKE$SESSION', 'ROLE') = 'COMP_TEAM_ROLE' THEN NULL ELSE 0 END"
		externalValue        = "(7)"
	)

	// max_rows is updated in place with its own SET. comment is SET/UNSET and cannot share that SET.
	basic := model.DataMovementRule("t", id.DatabaseName(), id.SchemaName(), id.Name(), maxRows, string(movementType))
	withOptionals := model.DataMovementRule("t", id.DatabaseName(), id.SchemaName(), id.Name(), otherMaxRows, string(movementType)).
		WithComment(comment)
	commentUnset := model.DataMovementRule("t", id.DatabaseName(), id.SchemaName(), id.Name(), otherMaxRows, string(movementType))
	withInternalType := model.DataMovementRule("t", id.DatabaseName(), id.SchemaName(), id.Name(), maxRows, string(internalMovementType))
	withSnowsightType := model.DataMovementRule("t", id.DatabaseName(), id.SchemaName(), id.Name(), maxRows, string(snowsightType))
	renamed := model.DataMovementRule("t", id.DatabaseName(), id.SchemaName(), renamedId.Name(), maxRows, string(snowsightType))
	ref := basic.ResourceReference()

	assertBasic := dataMovementRuleAssertions(t, ref, id, movementType, maxRows, "")
	assertWithOptionals := dataMovementRuleAssertions(t, ref, id, movementType, otherMaxRows, comment)
	assertCommentUnset := dataMovementRuleAssertions(t, ref, id, movementType, otherMaxRows, "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.DataMovementRule),
		Steps: []resource.TestStep{
			// Create the rule with only the required fields.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import the rule with only the required fields configured.
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedDataMovementRuleResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasTypeString(string(movementType)).
						HasMaxRowsString(maxRows).
						HasCommentEmpty(),
					resourceshowoutputassert.ImportedDataMovementRuleShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasType(movementType).
						HasComment(""),
					resourceshowoutputassert.ImportedDataMovementRuleDescribeOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasType(movementType).
						HasMaxRowsSignature("()").
						HasMaxRowsBody(maxRows),
				),
			},
			// Replace the external-stage rule with an internal-stage rule.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, withInternalType),
				Check:  assertThat(t, dataMovementRuleAssertions(t, ref, id, internalMovementType, maxRows, "")...),
			},
			// Replace the internal-stage rule with the original external-stage rule.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Update max_rows and set the comment in place.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "max_rows", tfjson.ActionUpdate, sdk.String(maxRows), sdk.String(otherMaxRows)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, nil, sdk.String(comment)),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// Import the rule with the updated max_rows expression and comment.
			{
				Config:            config.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedDataMovementRuleResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasTypeString(string(movementType)).
						HasMaxRowsString(otherMaxRows).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedDataMovementRuleShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasComment(comment),
					resourceshowoutputassert.ImportedDataMovementRuleDescribeOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasMaxRowsBody(otherMaxRows),
				),
			},
			// Remove the comment while preserving max_rows.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
						planchecks.ExpectNoChangeOnField(ref, "max_rows"),
					},
				},
				Config: config.FromModels(t, commentUnset),
				Check:  assertThat(t, assertCommentUnset...),
			},
			// Destroy the rule.
			{
				Config:  config.FromModels(t, commentUnset),
				Destroy: true,
			},
			// Verify the rule is gone, then recreate it with all fields configured.
			{
				PreConfig: func() {
					_, err := testClient().DataMovementRule.Show(t, id)
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
			// Verify the recreated rule imports with all fields.
			{
				Config:            config.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Apply the configured max_rows and comment after changing both outside Terraform.
			{
				PreConfig: func() {
					testClient().DataMovementRule.Alter(t, sdk.NewAlterDataMovementRuleRequest(id).
						WithSet(*sdk.NewDataMovementRuleSetRequest().WithSetMaxRows(externalValue)))
					testClient().DataMovementRule.Alter(t, sdk.NewAlterDataMovementRuleRequest(id).
						WithSet(*sdk.NewDataMovementRuleSetRequest().WithComment(externalComment)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "max_rows", sdk.String(otherMaxRows), sdk.String(externalValue)),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "max_rows", tfjson.ActionUpdate, sdk.String(externalValue), sdk.String(otherMaxRows)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// Replace the rule after an out-of-band change to its movement type.
			{
				PreConfig: func() {
					testClient().DataMovementRule.CreateWithRequest(t, sdk.NewCreateDataMovementRuleRequest(id, internalMovementType, otherMaxRows).
						WithOrReplace(true).
						WithComment(comment))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "type", sdk.String(string(movementType)), sdk.String(string(internalMovementType))),
						planchecks.ExpectChange(ref, "type", tfjson.ActionDelete, sdk.String(string(internalMovementType)), sdk.String(string(movementType))),
						planchecks.ExpectNoChangeOnField(ref, "max_rows"),
						planchecks.ExpectNoChangeOnField(ref, "comment"),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// Replace the external-stage rule with a Snowsight rule.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, withSnowsightType),
				Check:  assertThat(t, dataMovementRuleAssertions(t, ref, id, snowsightType, maxRows, "")...),
			},
			// Rename the Snowsight rule in place.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, renamed),
				Check:  assertThat(t, dataMovementRuleAssertions(t, ref, renamedId, snowsightType, maxRows, "")...),
			},
		},
	})
}

func dataMovementRuleAssertions(t *testing.T, ref string, id sdk.SchemaObjectIdentifier, movementType sdk.DataMovementRuleType, maxRows string, comment string) []assert.TestCheckFuncProvider {
	t.Helper()
	resourceCheck := resourceassert.DataMovementRuleResource(t, ref).
		HasNameString(id.Name()).
		HasDatabaseString(id.DatabaseName()).
		HasSchemaString(id.SchemaName()).
		HasTypeString(string(movementType)).
		HasMaxRowsString(maxRows).
		HasFullyQualifiedNameString(id.FullyQualifiedName())
	showCheck := resourceshowoutputassert.DataMovementRuleShowOutput(t, ref).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasComment(comment).
		HasType(movementType).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE")
	describeCheck := resourceshowoutputassert.DataMovementRuleDescribeOutput(t, ref).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasComment(comment).
		HasType(movementType).
		HasMaxRowsSignature("()").
		HasMaxRowsBody(maxRows)
	showObject := objectassert.DataMovementRule(t, id).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasComment(comment).
		HasType(movementType).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasOwnerRoleType("ROLE")
	details := objectassert.DataMovementRuleDetails(t, id).
		HasName(id.Name()).
		HasType(movementType).
		HasMaxRowsSignature("()").
		HasMaxRowsBody(maxRows).
		HasComment(comment)
	if comment == "" {
		resourceCheck = resourceCheck.HasCommentEmpty()
	} else {
		resourceCheck = resourceCheck.HasCommentString(comment)
	}
	return []assert.TestCheckFuncProvider{resourceCheck, showCheck, describeCheck, showObject, details}
}
