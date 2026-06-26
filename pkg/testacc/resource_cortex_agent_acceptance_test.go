//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CortexAgent_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	response := "You are a helpful assistant"
	hclSpec := model.SampleSpecAsYamlencodeHCL(response)
	normalizedSpec, err := sdk.NormalizeCortexAgentSpecification(
		testClient().CortexAgent.SampleSpecWithResponse(t, response),
	)
	require.NoError(t, err)

	newResponse := "You will respond in a friendly but concise manner"
	newHclSpec := model.SampleSpecAsYamlencodeHCL(newResponse)
	newNormalizedSpec, err := sdk.NormalizeCortexAgentSpecification(
		testClient().CortexAgent.SampleSpecWithResponse(t, newResponse),
	)
	require.NoError(t, err)

	comment := random.Comment()
	externalComment := random.Comment()

	emptyProfile := sdk.CortexAgentProfile{}
	completeProfile := sdk.CortexAgentProfile{
		DisplayName: sdk.String("My Helpful Assistant"),
		Avatar:      sdk.String("business-icon.png"),
		Color:       sdk.String("red"),
	}
	partialProfile := sdk.CortexAgentProfile{
		Color: sdk.String("green"),
	}
	externalProfile := sdk.CortexAgentProfile{
		DisplayName: sdk.String("My Friendly Assistant"),
		Color:       sdk.String("red"),
	}

	basic := model.CortexAgentWithSpecification("t", id.DatabaseName(), id.SchemaName(), id.Name(), hclSpec)

	complete := model.CortexAgentWithSpecification("t", id.DatabaseName(), id.SchemaName(), id.Name(), newHclSpec).
		WithComment(comment).
		WithProfile(completeProfile)

	withPartialProfile := model.CortexAgentWithSpecification("t", id.DatabaseName(), id.SchemaName(), id.Name(), newHclSpec).
		WithComment(comment).
		WithProfile(partialProfile)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CortexAgentResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasSpecification(normalizedSpec).
			HasCommentEmpty().
			HasProfileEmpty(),
		resourceshowoutputassert.CortexAgentShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasProfile(emptyProfile),
		resourceshowoutputassert.CortexAgentDescribeOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasProfile(emptyProfile).
			HasAgentSpec(normalizedSpec).
			HasCreatedOnNotEmpty().
			HasDefaultVersionName("LAST").
			HasVersions(`["VERSION$1"]`).
			HasAliases(`{"DEFAULT":"VERSION$1","FIRST":"VERSION$1","LAST":"VERSION$1"}`),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CortexAgentResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasSpecification(newNormalizedSpec).
			HasComment(comment).
			HasProfile(completeProfile),
		resourceshowoutputassert.CortexAgentShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasProfile(completeProfile),
		resourceshowoutputassert.CortexAgentDescribeOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasProfile(completeProfile).
			HasAgentSpec(newNormalizedSpec).
			HasCreatedOnNotEmpty().
			HasDefaultVersionName("LAST").
			HasVersions(`["VERSION$1"]`).
			HasAliases(`{"DEFAULT":"VERSION$1","FIRST":"VERSION$1","LAST":"VERSION$1"}`),
	}

	withPartialProfileAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CortexAgentResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasSpecification(newNormalizedSpec).
			HasComment(comment).
			HasProfile(partialProfile),
		resourceshowoutputassert.CortexAgentShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasProfile(partialProfile),
		resourceshowoutputassert.CortexAgentDescribeOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasProfile(partialProfile).
			HasAgentSpec(newNormalizedSpec).
			HasCreatedOnNotEmpty().
			HasDefaultVersionName("LAST").
			HasVersions(`["VERSION$1"]`).
			HasAliases(`{"DEFAULT":"VERSION$1","FIRST":"VERSION$1","LAST":"VERSION$1"}`),
	}

	basicAssertionsWithEmptyValues := []assert.TestCheckFuncProvider{
		resourceassert.CortexAgentResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasSpecification(normalizedSpec).
			HasCommentEmpty().
			HasProfileEmpty(),
		resourceshowoutputassert.CortexAgentShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasProfile(emptyProfile),
		resourceshowoutputassert.CortexAgentDescribeOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasProfile(emptyProfile).
			HasAgentSpec(normalizedSpec).
			HasCreatedOnNotEmpty().
			HasDefaultVersionName("LAST").
			HasVersions(`["VERSION$1"]`).
			HasAliases(`{"DEFAULT":"VERSION$1","FIRST":"VERSION$1","LAST":"VERSION$1"}`),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CortexAgent),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Set comment + full profile + new specification
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Set partial profile
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withPartialProfile),
				Check:  assertThat(t, withPartialProfileAssertions...),
			},
			// Unset
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertionsWithEmptyValues...),
			},
			// Destroy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroy),
					},
				},
				Config:  config.FromModels(t, basic),
				Destroy: true,
			},
			// Create with all attributes
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import
			{
				Config:            config.FromModels(t, complete),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change all props externally
			{
				PreConfig: func() {
					alterRequest := sdk.NewAlterCortexAgentRequest(id).WithModifyLiveVersionSet(
						*sdk.NewCortexAgentModifyLiveVersionSetRequest(normalizedSpec),
					)
					testClient().CortexAgent.Alter(t, alterRequest)

					externalProfileAsJson, err := sdk.MarshalCortexAgentProfile(externalProfile)
					require.NoError(t, err)
					alterRequest = sdk.NewAlterCortexAgentRequest(id).WithSet(
						*sdk.NewCortexAgentSetRequest().
							WithComment(sdk.StringAllowEmpty{Value: externalComment}).
							WithProfile(externalProfileAsJson),
					)
					testClient().CortexAgent.Alter(t, alterRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: func() []plancheck.PlanCheck {
						return []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
							planchecks.ExpectDrift(ref, "specification", sdk.String(newNormalizedSpec), sdk.String(normalizedSpec)),
							// ExpectChange on specification is omitted: the after change value comes from yamlencode in
							// the configuration, and there is no string literal we can pass to ExpectChange that
							// matches the planned "after" exactly.
							planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
							planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
							planchecks.ExpectDrift(ref, "profile.0.display_name", sdk.String("My Helpful Assistant"), sdk.String("My Friendly Assistant")),
							planchecks.ExpectChange(ref, "profile.0.display_name", tfjson.ActionUpdate, sdk.String("My Friendly Assistant"), sdk.String("My Helpful Assistant")),
							planchecks.ExpectDrift(ref, "profile.0.avatar", sdk.String("business-icon.png"), sdk.String("")),
							planchecks.ExpectChange(ref, "profile.0.avatar", tfjson.ActionUpdate, sdk.String(""), sdk.String("business-icon.png")),
							planchecks.ExpectNoChangeOnField(ref, "profile.0.color"),
						}
					}(),
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
		},
	})
}

func TestAcc_CortexAgent_CompleteUseCase_EmptyAndNullCommentsHandling(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	response := "You are a helpful assistant"
	hclSpec := model.SampleSpecAsYamlencodeHCL(response)
	normalizedSpec, err := sdk.NormalizeCortexAgentSpecification(
		testClient().CortexAgent.SampleSpecWithResponse(t, response),
	)
	require.NoError(t, err)
	emptyProfile := sdk.CortexAgentProfile{}

	basic := model.CortexAgentWithSpecification("t", id.DatabaseName(), id.SchemaName(), id.Name(), hclSpec)

	basicWithEmptyComment := model.CortexAgentWithSpecification("t", id.DatabaseName(), id.SchemaName(), id.Name(), hclSpec).
		WithComment("")

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CortexAgentResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasSpecification(normalizedSpec).
			HasCommentEmpty().
			HasProfileEmpty(),
		resourceshowoutputassert.CortexAgentShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasProfile(emptyProfile),
		resourceshowoutputassert.CortexAgentDescribeOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasProfile(emptyProfile).
			HasAgentSpec(normalizedSpec).
			HasCreatedOnNotEmpty().
			HasDefaultVersionName("LAST").
			HasVersions(`["VERSION$1"]`).
			HasAliases(`{"DEFAULT":"VERSION$1","FIRST":"VERSION$1","LAST":"VERSION$1"}`),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CortexAgent),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Set empty comment externally and expect no drift
			{
				PreConfig: func() {
					alterRequest := sdk.NewAlterCortexAgentRequest(id).WithSet(
						*sdk.NewCortexAgentSetRequest().
							WithComment(sdk.StringAllowEmpty{Value: ""}),
					)
					testClient().CortexAgent.Alter(t, alterRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Set empty comment and expect empty plan
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: config.FromModels(t, basicWithEmptyComment),
				Check:  assertThat(t, basicAssertions...),
			},
			// Set comment to NULL externally and expect no drift
			{
				PreConfig: func() {
					// There's no way to set a comment to NULL other than recreating an object.
					replaceRequest := sdk.NewCreateCortexAgentRequest(id, normalizedSpec).
						WithOrReplace(true)
					testClient().CortexAgent.CreateWithRequest(t, replaceRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: config.FromModels(t, basicWithEmptyComment),
				Check:  assertThat(t, basicAssertions...),
			},
		},
	})
}

func TestAcc_CortexAgent_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	emptySpec := model.CortexAgent("t", id.DatabaseName(), id.SchemaName(), id.Name(), "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CortexAgent),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, emptySpec),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "specification" to not be an empty string`),
			},
		},
	})
}
