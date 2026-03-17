//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CatalogIntegrationAwsGlue_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	glueAwsRoleArn := "arn:aws:iam::123456789012:role/sqsAccess"
	newGlueAwsRoleArn := "arn:aws:iam::123456789012:role/fakeRole"
	externalGlueAwsRoleArn := "arn:aws:iam::123456789012:role/fakeRole2"
	glueCatalogId := random.NumericN(15)
	externalGlueCatalogId := random.NumericN(15)

	comment := random.Comment()
	newComment := random.Comment()
	externalComment := random.Comment()

	refreshIntervalSeconds := random.IntRange(30, 86400)
	newRefreshIntervalSeconds := random.IntRange(30, 86400)
	externalRefreshIntervalSeconds := random.IntRange(30, 86400)

	glueRegion := "us-east-1"
	newGlueRegion := "eu-west-1"
	externalGlueRegion := "eu-west-2"
	catalogNamespace := random.AlphanumericN(15)
	externalCatalogNamespace := random.AlphanumericN(15)

	catalogIntegrationAwsGlueBasic := model.CatalogIntegrationAwsGlue("t", id.Name(), false, glueAwsRoleArn, glueCatalogId)

	catalogIntegrationAwsGlueAltered := model.CatalogIntegrationAwsGlue("t", id.Name(), true, glueAwsRoleArn, glueCatalogId).
		WithComment(newComment).
		WithRefreshIntervalSeconds(newRefreshIntervalSeconds)

	catalogIntegrationAwsGlueAllAttributes := model.CatalogIntegrationAwsGlue("t", id.Name(), false, glueAwsRoleArn, glueCatalogId).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithGlueRegion(glueRegion).
		WithCatalogNamespace(catalogNamespace)

	catalogIntegrationAwsGlueAllAttributesWithExternalChanges := model.CatalogIntegrationAwsGlue("t", id.Name(), true, glueAwsRoleArn, glueCatalogId).
		WithComment(externalComment).
		WithRefreshIntervalSeconds(externalRefreshIntervalSeconds).
		WithGlueRegion(glueRegion).
		WithCatalogNamespace(catalogNamespace)

	catalogIntegrationAwsGlueWithChangedForceNewAttributes := model.CatalogIntegrationAwsGlue("t", id.Name(), false, newGlueAwsRoleArn, glueCatalogId)

	catalogIntegrationAwsGlueWithMoreChangedForceNewAttributes := model.CatalogIntegrationAwsGlue("t", id.Name(), false, newGlueAwsRoleArn, glueCatalogId).
		WithGlueRegion(newGlueRegion)

	catalogIntegrationAwsGlueWithMoreExternalChanges := model.CatalogIntegrationAwsGlue("t", id.Name(), false, externalGlueAwsRoleArn, externalGlueCatalogId).
		WithGlueRegion(externalGlueRegion).
		WithCatalogNamespace(externalCatalogNamespace)

	ref := catalogIntegrationAwsGlueBasic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationAwsGlueResource(t, ref).
			HasName(id.Name()).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasNoGlueRegion().
			HasCatalogNamespaceEmpty(),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(""),
		resourceshowoutputassert.CatalogIntegrationAwsGlueDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeAWSGlue).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			// Don't check glue_region, as its default value depends on the current region name
			HasCatalogNamespace(""),
	}

	basicAssertionsWithRefreshIntervalZero := append(
		[]assert.TestCheckFuncProvider{
			resourceassert.CatalogIntegrationAwsGlueResource(t, ref).
				HasName(id.Name()).
				HasEnabledString(r.BooleanFalse).
				HasCommentEmpty().
				HasRefreshIntervalSeconds(0).
				HasGlueAwsRoleArn(glueAwsRoleArn).
				HasGlueCatalogId(glueCatalogId).
				HasNoGlueRegion().
				HasCatalogNamespaceEmpty(),
		},
		basicAssertions[1:]...,
	)

	alteredProperties := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationAwsGlueResource(t, ref).
			HasName(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasComment(newComment).
			HasRefreshIntervalSeconds(newRefreshIntervalSeconds).
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasNoGlueRegion().
			HasCatalogNamespaceEmpty(),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(true).
			HasComment(newComment),
		resourceshowoutputassert.CatalogIntegrationAwsGlueDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeAWSGlue).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(true).
			HasRefreshIntervalSeconds(newRefreshIntervalSeconds).
			HasComment(newComment).
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			// Don't check glue_region, as its default value depends on the current region name
			HasCatalogNamespace(""),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationAwsGlueResource(t, ref).
			HasName(id.Name()).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasGlueRegion(glueRegion).
			HasCatalogNamespace(catalogNamespace),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationAwsGlueDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeAWSGlue).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasGlueRegion(glueRegion).
			HasCatalogNamespace(catalogNamespace),
	}

	externallyAlteredProperties := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationAwsGlueResource(t, ref).
			HasName(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasComment(externalComment).
			HasRefreshIntervalSeconds(externalRefreshIntervalSeconds).
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasGlueRegion(glueRegion).
			HasCatalogNamespace(catalogNamespace),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(true).
			HasComment(externalComment),
		resourceshowoutputassert.CatalogIntegrationAwsGlueDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeAWSGlue).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(true).
			HasRefreshIntervalSeconds(externalRefreshIntervalSeconds).
			HasComment(externalComment).
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasGlueRegion(glueRegion).
			HasCatalogNamespace(catalogNamespace),
	}

	forceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationAwsGlueResource(t, ref).
			HasName(id.Name()).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasGlueAwsRoleArn(newGlueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasNoGlueRegion().
			HasCatalogNamespaceEmpty(),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(""),
		resourceshowoutputassert.CatalogIntegrationAwsGlueDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeAWSGlue).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasGlueAwsRoleArn(newGlueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			// Don't check glue_region, as its default value depends on the current region name
			HasCatalogNamespace(""),
	}

	moreForceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationAwsGlueResource(t, ref).
			HasName(id.Name()).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasGlueAwsRoleArn(newGlueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasGlueRegion(newGlueRegion).
			HasCatalogNamespaceEmpty(),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(""),
		resourceshowoutputassert.CatalogIntegrationAwsGlueDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeAWSGlue).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasGlueAwsRoleArn(newGlueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasGlueRegion(newGlueRegion).
			HasCatalogNamespace(""),
	}

	moreExternallyAlteredProperties := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationAwsGlueResource(t, ref).
			HasName(id.Name()).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasGlueAwsRoleArn(externalGlueAwsRoleArn).
			HasGlueCatalogId(externalGlueCatalogId).
			HasGlueRegion(externalGlueRegion).
			HasCatalogNamespace(externalCatalogNamespace),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(""),
		resourceshowoutputassert.CatalogIntegrationAwsGlueDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeAWSGlue).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasGlueAwsRoleArn(externalGlueAwsRoleArn).
			HasGlueCatalogId(externalGlueCatalogId).
			HasGlueRegion(externalGlueRegion).
			HasCatalogNamespace(externalCatalogNamespace),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationAwsGlue),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, catalogIntegrationAwsGlueBasic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import
			{
				Config:            config.FromModels(t, catalogIntegrationAwsGlueBasic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change alterable props
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, catalogIntegrationAwsGlueAltered),
				Check:  assertThat(t, alteredProperties...),
			},
			// Unset
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, catalogIntegrationAwsGlueBasic),
				Check:  assertThat(t, basicAssertionsWithRefreshIntervalZero...),
			},
			// Destroy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroy),
					},
				},
				Config:  config.FromModels(t, catalogIntegrationAwsGlueBasic),
				Destroy: true,
			},
			// Create with all attributes
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, catalogIntegrationAwsGlueAllAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Change alterable props externally
			{
				PreConfig: func() {
					alterRequest := sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*sdk.NewCatalogIntegrationSetRequest().
						WithEnabled(true).
						WithComment(sdk.StringAllowEmpty{Value: externalComment}).
						WithRefreshIntervalSeconds(externalRefreshIntervalSeconds),
					)
					testClient().CatalogIntegration.Alter(t, alterRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "enabled", sdk.String("false"), sdk.String("true")),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
						planchecks.ExpectDrift(ref, "refresh_interval_seconds", sdk.String(strconv.Itoa(refreshIntervalSeconds)), sdk.String(strconv.Itoa(externalRefreshIntervalSeconds))),
						planchecks.ExpectChange(ref, "enabled", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
						planchecks.ExpectChange(ref, "refresh_interval_seconds", tfjson.ActionUpdate, sdk.String(strconv.Itoa(externalRefreshIntervalSeconds)), sdk.String(strconv.Itoa(refreshIntervalSeconds))),
					},
				},
				Config: config.FromModels(t, catalogIntegrationAwsGlueAllAttributes),
			},
			// Adjust config to the current Snowflake values
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: config.FromModels(t, catalogIntegrationAwsGlueAllAttributesWithExternalChanges),
				Check:  assertThat(t, externallyAlteredProperties...),
			},
			// Change force new "glue_aws_role_arn" prop
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, catalogIntegrationAwsGlueWithChangedForceNewAttributes),
				Check:  assertThat(t, forceNewAssertions...),
			},
			// Change force new "glue_region" prop
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, catalogIntegrationAwsGlueWithMoreChangedForceNewAttributes),
				Check:  assertThat(t, moreForceNewAssertions...),
			},
			// Change force new props externally
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithOrReplace(true).
						WithAwsGlueCatalogSourceParams(*sdk.NewAwsGlueParamsRequest(externalGlueAwsRoleArn, externalGlueCatalogId).
							WithGlueRegion(externalGlueRegion).
							WithCatalogNamespace(externalCatalogNamespace))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "glue_aws_role_arn", sdk.String(newGlueAwsRoleArn), sdk.String(externalGlueAwsRoleArn)),
						planchecks.ExpectDrift(ref, "glue_catalog_id", sdk.String(glueCatalogId), sdk.String(externalGlueCatalogId)),
						planchecks.ExpectDrift(ref, "glue_region", sdk.String(newGlueRegion), sdk.String(externalGlueRegion)),
						planchecks.ExpectDrift(ref, "catalog_namespace", sdk.String(""), sdk.String(externalCatalogNamespace)),
						planchecks.ExpectChange(ref, "glue_aws_role_arn", tfjson.ActionDelete, sdk.String(externalGlueAwsRoleArn), sdk.String(newGlueAwsRoleArn)),
						planchecks.ExpectChange(ref, "glue_catalog_id", tfjson.ActionDelete, sdk.String(externalGlueCatalogId), sdk.String(glueCatalogId)),
						planchecks.ExpectChange(ref, "glue_region", tfjson.ActionDelete, sdk.String(externalGlueRegion), sdk.String(newGlueRegion)),
						planchecks.ExpectChange(ref, "catalog_namespace", tfjson.ActionDelete, sdk.String(externalCatalogNamespace), nil),
					},
				},
				Config: config.FromModels(t, catalogIntegrationAwsGlueWithMoreChangedForceNewAttributes),
			},
			// Adjust config to the current Snowflake values
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: config.FromModels(t, catalogIntegrationAwsGlueWithMoreExternalChanges),
				Check:  assertThat(t, moreExternallyAlteredProperties...),
			},
		},
	})
}

func TestAcc_CatalogIntegrationAwsGlue_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	glueAwsRoleArn := "arn:aws:iam::123456789012:role/sqsAccess"
	glueCatalogId := random.NumericN(15)

	catalogIntegrationAwsGlueRefreshIntervalNonPositive := model.CatalogIntegrationAwsGlue("t", id.Name(), false, glueAwsRoleArn, glueCatalogId).
		WithRefreshIntervalSeconds(0)

	catalogIntegrationAwsGlueEmptyAwsRoleArn := model.CatalogIntegrationAwsGlue("t", id.Name(), false, "", glueCatalogId)
	catalogIntegrationAwsGlueEmptyCatalogId := model.CatalogIntegrationAwsGlue("t", id.Name(), false, glueAwsRoleArn, "")
	catalogIntegrationAwsGlueEmptyRegion := model.CatalogIntegrationAwsGlue("t", id.Name(), false, glueAwsRoleArn, glueCatalogId).
		WithGlueRegion("")
	catalogIntegrationAwsGlueEmptyCatalogNamespace := model.CatalogIntegrationAwsGlue("t", id.Name(), false, glueAwsRoleArn, glueCatalogId).
		WithCatalogNamespace("")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationAwsGlue),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, catalogIntegrationAwsGlueRefreshIntervalNonPositive),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected refresh_interval_seconds to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, catalogIntegrationAwsGlueEmptyAwsRoleArn),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "glue_aws_role_arn" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, catalogIntegrationAwsGlueEmptyCatalogId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "glue_catalog_id" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, catalogIntegrationAwsGlueEmptyRegion),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "glue_region" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, catalogIntegrationAwsGlueEmptyCatalogNamespace),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "catalog_namespace" to not be an empty string`),
			},
		},
	})
}

func TestAcc_CatalogIntegrationAwsGlue_ImportValidation(t *testing.T) {
	glueAwsRoleArn := "arn:aws:iam::123456789012:role/sqsAccess"
	glueCatalogId := random.NumericN(15)

	notificationIntegration, notificationIntegrationCleanup := testClient().NotificationIntegration.Create(t)
	t.Cleanup(notificationIntegrationCleanup)

	catalogIntegrationObjectStorage, catalogIntegrationObjectStorageCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogIntegrationObjectStorageCleanup)

	catalogIntegrationAwsGlue := model.CatalogIntegrationAwsGlue("t", notificationIntegration.ID().Name(), false, glueAwsRoleArn, glueCatalogId)
	catalogIntegrationAwsGlue2 := model.CatalogIntegrationAwsGlue("t", catalogIntegrationObjectStorage.Name(), false, glueAwsRoleArn, glueCatalogId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationAwsGlue),
		Steps: []resource.TestStep{
			// import a different integration category
			{
				Config:        config.FromModels(t, catalogIntegrationAwsGlue),
				ResourceName:  catalogIntegrationAwsGlue.ResourceReference(),
				ImportState:   true,
				ImportStateId: notificationIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile(fmt.Sprintf(`Integration %s is not a CATALOG integration`, notificationIntegration.ID().Name())),
			},
			// import a different catalog source type
			{
				Config:        config.FromModels(t, catalogIntegrationAwsGlue2),
				ResourceName:  catalogIntegrationAwsGlue2.ResourceReference(),
				ImportState:   true,
				ImportStateId: catalogIntegrationObjectStorage.Name(),
				ExpectError:   regexp.MustCompile(fmt.Sprintf(`invalid catalog source type, expected %s, got %s`, sdk.CatalogIntegrationCatalogSourceTypeAWSGlue, sdk.CatalogIntegrationCatalogSourceTypeObjectStorage)),
			},
		},
	})
}
