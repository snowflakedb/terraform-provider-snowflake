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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CodeBundle_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	externalComment := random.Comment()

	stage, cleanupStage := testClient().Stage.CreateStage(t)
	t.Cleanup(cleanupStage)
	testClient().Stage.PutOnStageWithContent(t, stage.ID(), "main.py", "print('code bundle')\n")
	from := stage.Location()

	// from is ForceNew. comment is changed in place with SET COMMENT and UNSET COMMENT.
	basic := model.CodeBundle("t", id.DatabaseName(), id.SchemaName(), id.Name(), from)
	complete := model.CodeBundle("t", id.DatabaseName(), id.SchemaName(), id.Name(), from).
		WithComment(comment)
	ref := basic.ResourceReference()

	assertBasic := codeBundleAssertions(t, ref, id, from, stage, "")
	assertComplete := codeBundleAssertions(t, ref, id, from, stage, comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CodeBundle),
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
				Config:                  config.FromModels(t, basic),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"from"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedCodeBundleResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentEmpty(),
					resourceshowoutputassert.ImportedCodeBundleShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner("").
						HasComment(""),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, nil, sdk.String(comment)),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			{
				Config:                  config.FromModels(t, complete),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"from"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedCodeBundleResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedCodeBundleShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasComment(comment),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			{
				PreConfig: func() {
					testClient().CodeBundle.CreateWithRequest(t, sdk.NewCreateCodeBundleRequest(id, from).
						WithOrReplace(true).
						WithComment(externalComment),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(""), sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), nil),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
		},
	})
}

func codeBundleAssertions(t *testing.T, ref string, id sdk.SchemaObjectIdentifier, from string, stage *sdk.Stage, comment string) []assert.TestCheckFuncProvider {
	t.Helper()
	resourceCheck := resourceassert.CodeBundleResource(t, ref).
		HasNameString(id.Name()).
		HasDatabaseString(id.DatabaseName()).
		HasSchemaString(id.SchemaName()).
		HasFromString(from).
		HasFullyQualifiedNameString(id.FullyQualifiedName())
	showCheck := resourceshowoutputassert.CodeBundleShowOutput(t, ref).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwner("")
	describeCheck := resourceshowoutputassert.CodeBundleDescribeOutput(t, ref).
		HasName(id.Name()).
		HasDefaultVersion("LAST").
		HasContainerType("Stage").
		HasContainerName(stage.ID().Name()).
		HasOwnerRole("ACCOUNTADMIN").
		HasOwnerRoleType("ROLE")
	showObject := objectassert.CodeBundle(t, id).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwner("")
	details := objectassert.CodeBundleDetails(t, id).
		HasName(id.Name()).
		HasDefaultVersion("LAST").
		HasDefaultVersionSourceLocationUriContaining(stage.ID().Name()).
		HasContainerType("Stage").
		HasContainerName(stage.ID().Name()).
		HasOwnerRole("ACCOUNTADMIN").
		HasOwnerRoleType("ROLE")
	if comment == "" {
		resourceCheck = resourceCheck.HasCommentEmpty()
		showCheck = showCheck.HasComment("")
		describeCheck = describeCheck.HasComment("")
		showObject = showObject.HasNoComment()
		details = details.HasNoComment()
	} else {
		resourceCheck = resourceCheck.HasCommentString(comment)
		showCheck = showCheck.HasComment(comment)
		describeCheck = describeCheck.HasComment(comment)
		showObject = showObject.HasComment(comment)
		details = details.HasComment(comment)
	}
	return []assert.TestCheckFuncProvider{resourceCheck, showCheck, describeCheck, showObject, details}
}
