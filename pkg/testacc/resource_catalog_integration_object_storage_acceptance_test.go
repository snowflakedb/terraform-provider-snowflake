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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CatalogIntegrationObjectStorage_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()
	newComment := random.Comment()
	externalComment := random.Comment()

	refreshIntervalSeconds := random.IntRange(30, 86400)
	newRefreshIntervalSeconds := random.IntRange(30, 86400)
	externalRefreshIntervalSeconds := random.IntRange(30, 86400)

	tableFormat := string(sdk.CatalogIntegrationTableFormatDelta)
	basic := model.CatalogIntegrationObjectStorage("t", id.Name(), false, tableFormat)

	altered := model.CatalogIntegrationObjectStorage("t", id.Name(), true, tableFormat).
		WithComment(newComment).
		WithRefreshIntervalSeconds(newRefreshIntervalSeconds)

	allAttributes := model.CatalogIntegrationObjectStorage("t", id.Name(), false, tableFormat).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds)

	withChangedTableFormat := model.CatalogIntegrationObjectStorage("t", id.Name(), false, string(sdk.CatalogIntegrationTableFormatIceberg))

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationObjectStorageResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeObjectStore)).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasTableFormat(tableFormat),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(""),
		resourceshowoutputassert.CatalogIntegrationObjectStorageDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeObjectStore).
			HasTableFormat(sdk.CatalogIntegrationTableFormatDelta).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment(""),
	}

	basicAssertionsWithRefreshIntervalZero := append(
		[]assert.TestCheckFuncProvider{
			resourceassert.CatalogIntegrationObjectStorageResource(t, ref).
				HasName(id.Name()).
				HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeObjectStore)).
				HasEnabledString(r.BooleanFalse).
				HasCommentEmpty().
				HasRefreshIntervalSeconds(0).
				HasTableFormat(tableFormat),
		},
		basicAssertions[1:]...,
	)

	alteredProperties := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationObjectStorageResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeObjectStore)).
			HasEnabledString(r.BooleanTrue).
			HasComment(newComment).
			HasRefreshIntervalSeconds(newRefreshIntervalSeconds).
			HasTableFormat(tableFormat),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(true).
			HasComment(newComment),
		resourceshowoutputassert.CatalogIntegrationObjectStorageDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeObjectStore).
			HasTableFormat(sdk.CatalogIntegrationTableFormatDelta).
			HasEnabled(true).
			HasRefreshIntervalSeconds(newRefreshIntervalSeconds).
			HasComment(newComment),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationObjectStorageResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeObjectStore)).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasTableFormat(tableFormat),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationObjectStorageDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeObjectStore).
			HasTableFormat(sdk.CatalogIntegrationTableFormatDelta).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment),
	}

	forceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationObjectStorageResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeObjectStore)).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasTableFormat(string(sdk.CatalogIntegrationTableFormatIceberg)),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(""),
		resourceshowoutputassert.CatalogIntegrationObjectStorageDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeObjectStore).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment(""),
	}

	providerModelImportPermissiveMode := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationObjectStorage),
		Steps: []resource.TestStep{
			// Create
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import - permissive mode
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   config.FromModels(t, providerModelImportPermissiveMode, basic),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
			// Change alterable props
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, altered),
				Check:  assertThat(t, alteredProperties...),
			},
			// Unset
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertionsWithRefreshIntervalZero...),
			},
			// Destroy
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
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
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Change alterable props externally
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					alterRequest := sdk.NewAlterCatalogIntegrationRequest(id).WithSet(
						*sdk.NewCatalogIntegrationSetRequest().
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
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Change force new "table_format" prop
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, withChangedTableFormat),
				Check:  assertThat(t, forceNewAssertions...),
			},
			// Change force new "table_format" prop externally
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					createRequest := sdk.NewCreateObjectStorageCatalogIntegrationRequest(id, sdk.CatalogIntegrationTableFormatDelta, false).
						WithOrReplace(true)
					testClient().CatalogIntegration.CreateObjectStorageFunc(t, createRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "table_format", sdk.String(string(sdk.CatalogIntegrationTableFormatIceberg)), sdk.String(string(sdk.CatalogIntegrationTableFormatDelta))),
						planchecks.ExpectChange(ref, "table_format", tfjson.ActionDelete, sdk.String(string(sdk.CatalogIntegrationTableFormatDelta)), sdk.String(string(sdk.CatalogIntegrationTableFormatIceberg))),
					},
				},
				Config: config.FromModels(t, withChangedTableFormat),
				Check:  assertThat(t, forceNewAssertions...),
			},
			// Change "catalog_source" externally
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					createRequest := sdk.NewCreateAwsGlueCatalogIntegrationRequest(id, "arn:aws:iam::123456789012:role/sqsAccess", random.NumericN(15), false).
						WithOrReplace(true)
					testClient().CatalogIntegration.CreateAwsGlueFunc(t, createRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, withChangedTableFormat),
				Check:  assertThat(t, forceNewAssertions...),
			},
		},
	})
}

func TestAcc_CatalogIntegrationObjectStorage_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()
	refreshIntervalSeconds := random.IntRange(30, 86400)
	tableFormat := string(sdk.CatalogIntegrationTableFormatDelta)

	complete := model.CatalogIntegrationObjectStorage("t", id.Name(), false, tableFormat).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds)

	ref := complete.ResourceReference()

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationObjectStorageResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeObjectStore)).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasTableFormat(tableFormat).
			HasFullyQualifiedName(id.FullyQualifiedName()),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationObjectStorageDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeObjectStore).
			HasTableFormat(sdk.CatalogIntegrationTableFormatDelta).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment),
		objectassert.CatalogIntegration(t, id).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment).
			HasCreatedOnNotEmpty(),
		objectassert.CatalogIntegrationObjectStorageDetails(t, id).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeObjectStore).
			HasTableFormat(sdk.CatalogIntegrationTableFormatDelta).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment),
	}

	providerModelImportPermissiveMode := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationObjectStorage),
		Steps: []resource.TestStep{
			// Create - with all optionals (including optional force-new fields)
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import - with all optionals
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, complete),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateVerify:        true,
			},
			// Import - with all optionals, permissive mode
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   config.FromModels(t, providerModelImportPermissiveMode, complete),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
		},
	})
}

func TestAcc_CatalogIntegrationObjectStorage_BasicUseCase_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	tableFormat := sdk.CatalogIntegrationTableFormatDelta
	refreshIntervalSeconds := random.IntRange(31, 86400)

	basicModel := model.CatalogIntegrationObjectStorage("t", id.Name(), false, string(tableFormat))
	ref := basicModel.ResourceReference()

	_, cleanup := testClient().CatalogIntegration.CreateObjectStorageFunc(t, sdk.NewCreateObjectStorageCatalogIntegrationRequest(id, tableFormat, false).
		WithRefreshIntervalSeconds(refreshIntervalSeconds))
	t.Cleanup(cleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationObjectStorage),
		Steps: []resource.TestStep{
			// Import - without optionals
			{
				Config:             config.FromModels(t, basicModel),
				ResourceName:       ref,
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			// Expect empty plan after the plan after import is applied
			{
				Config: config.FromModels(t, basicModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "refresh_interval_seconds", tfjson.ActionUpdate, new(strconv.Itoa(refreshIntervalSeconds)), nil),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t,
					resourceassert.CatalogIntegrationObjectStorageResource(t, ref).
						HasRefreshIntervalSeconds(0),
				),
			},
		},
	})
}

func TestAcc_CatalogIntegrationObjectStorage_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	tableFormat := string(sdk.CatalogIntegrationTableFormatDelta)

	refreshIntervalNonPositive := model.CatalogIntegrationObjectStorage("t", id.Name(), false, tableFormat).
		WithRefreshIntervalSeconds(0)

	invalidTableFormat := model.CatalogIntegrationObjectStorage("t", id.Name(), false, "INVALID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationObjectStorage),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, refreshIntervalNonPositive),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected refresh_interval_seconds to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, invalidTableFormat),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid catalog integration table format: INVALID`),
			},
		},
	})
}

func TestAcc_CatalogIntegrationObjectStorage_ImportValidation(t *testing.T) {
	notificationIntegration, notificationIntegrationCleanup := testClient().NotificationIntegration.Create(t)
	t.Cleanup(notificationIntegrationCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	catalogIntegrationAwsGlue, catalogIntegrationAwsGlueCleanup := testClient().CatalogIntegration.CreateAwsGlueFunc(t,
		sdk.NewCreateAwsGlueCatalogIntegrationRequest(id, "arn:aws:iam::123456789012:role/sqsAccess", random.NumericN(15), false))
	t.Cleanup(catalogIntegrationAwsGlueCleanup)

	tableFormat := string(sdk.CatalogIntegrationTableFormatDelta)
	catalogIntegrationObjectStorage := model.CatalogIntegrationObjectStorage("t", notificationIntegration.ID().Name(), false, tableFormat)
	catalogIntegrationObjectStorage2 := model.CatalogIntegrationObjectStorage("t", catalogIntegrationAwsGlue.Name(), false, tableFormat)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationObjectStorage),
		Steps: []resource.TestStep{
			// import a different integration category
			{
				Config:        config.FromModels(t, catalogIntegrationObjectStorage),
				ResourceName:  catalogIntegrationObjectStorage.ResourceReference(),
				ImportState:   true,
				ImportStateId: notificationIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile(fmt.Sprintf(`Integration %s is not a CATALOG integration`, notificationIntegration.ID().Name())),
			},
			// import a different catalog source type
			{
				Config:        config.FromModels(t, catalogIntegrationObjectStorage2),
				ResourceName:  catalogIntegrationObjectStorage2.ResourceReference(),
				ImportState:   true,
				ImportStateId: catalogIntegrationAwsGlue.Name(),
				ExpectError:   regexp.MustCompile(fmt.Sprintf(`invalid catalog source type, expected %s, got %s`, sdk.CatalogIntegrationCatalogSourceTypeObjectStore, sdk.CatalogIntegrationCatalogSourceTypeGlue)),
			},
		},
	})
}

func TestAcc_CatalogIntegrationObjectStorage_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	tableFormat := string(sdk.CatalogIntegrationTableFormatDelta)
	comment := random.Comment()
	refreshIntervalSeconds := random.IntRange(30, 86400)

	allAttributes := model.CatalogIntegrationObjectStorage("t", id.Name(), false, tableFormat).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds)

	ref := allAttributes.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationObjectStorage),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					tableFormatEnum, _ := sdk.ToCatalogIntegrationTableFormat(tableFormat)
					createRequest := sdk.NewCreateObjectStorageCatalogIntegrationRequest(id, tableFormatEnum, false).
						WithRefreshIntervalSeconds(refreshIntervalSeconds).
						WithComment(comment)
					testClient().CatalogIntegration.CreateObjectStorageFunc(t, createRequest)
				},
				Config:             config.FromModels(t, allAttributes),
				ResourceName:       ref,
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			{
				Config: config.FromModels(t, allAttributes),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
