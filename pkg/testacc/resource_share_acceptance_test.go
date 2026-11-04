//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

// TODO [next PRs]: Add the describe output assertions once the share resource has describe_output. A share has no object parameters.
func TestAcc_Share_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)
	account := secondaryTestClient().Account.GetAccountIdentifier(t)
	comment := random.Comment()
	externalComment := random.Comment()

	providerModelImportPermissiveMode := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode)

	basic := model.Share("t", id.Name())
	withOptionals := model.Share("t", id.Name()).
		WithComment(comment).
		WithAccounts(account.Name())
	// secure_objects_only cannot be changed on an existing share, so it is set to a non-default value only here.
	complete := model.Share("t", id.Name()).
		WithComment(comment).
		WithAccounts(account.Name()).
		WithSecureObjectsOnlyBool(false)
	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasSecureObjectsOnlyString(r.BooleanDefault).
			HasCommentEmpty().
			HasAccountsEmpty(),
		resourceshowoutputassert.ShareShowOutput(t, ref).
			HasName(id.Name()).
			HasKind(sdk.ShareKindOutbound).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasListingGlobalName("").
			HasSecureObjectsOnly(true).
			HasToTotal(0),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasKind(sdk.ShareKindOutbound).
			HasOwnerAccountNotEmpty().
			HasComment("").
			HasSecureObjectsOnly(true).
			HasNoTo().
			HasToTotal(0),
	}
	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasSecureObjectsOnlyString(r.BooleanDefault).
			HasCommentString(comment).
			HasAccounts(account.Name()),
		resourceshowoutputassert.ShareShowOutput(t, ref).
			HasName(id.Name()).
			HasKind(sdk.ShareKindOutbound).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasListingGlobalName("").
			HasSecureObjectsOnly(true).
			HasToTotal(0),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasKind(sdk.ShareKindOutbound).
			HasOwnerAccountNotEmpty().
			HasComment(comment).
			HasSecureObjectsOnly(true).
			HasTo(account).
			HasToTotal(0),
	}
	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasSecureObjectsOnlyString("false").
			HasCommentString(comment).
			HasAccounts(account.Name()),
		resourceshowoutputassert.ShareShowOutput(t, ref).
			HasName(id.Name()).
			HasKind(sdk.ShareKindOutbound).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasListingGlobalName("").
			HasSecureObjectsOnly(false).
			HasToTotal(0),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasKind(sdk.ShareKindOutbound).
			HasOwnerAccountNotEmpty().
			HasComment(comment).
			HasSecureObjectsOnly(false).
			HasTo(account).
			HasToTotal(0),
	}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			// 1. Create - without optionals
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// 2. Import - without optionals (enforcing mode)
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, basic),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateVerify:        true,
				// The configuration omits secure_objects_only, so the state has "default" and the import (enforcing mode) the value of the share.
				ImportStateVerifyIgnore: []string{"secure_objects_only"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedShareResource(t, resourceId).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasSecureObjectsOnlyString("true").
						HasCommentEmpty().
						HasAccountsEmpty(),
					resourceshowoutputassert.ImportedShareShowOutput(t, resourceId).
						HasName(id.Name()).
						HasKind(sdk.ShareKindOutbound).
						HasComment("").
						HasSecureObjectsOnly(true).
						HasToTotal(0),
				),
			},
			// 2b. Import - without optionals (permissive mode: the plan after the import is empty)
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   accconfig.FromModels(t, providerModelImportPermissiveMode, basic),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
			// 3. Update - set optionals
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(""), sdk.String(comment)),
						planchecks.ExpectSetChange(ref, "accounts", tfjson.ActionUpdate, nil, []string{account.Name()}),
						planchecks.ExpectNoChangeOnField(ref, "secure_objects_only"),
					},
				},
				Config: accconfig.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// 4. Import - with optionals
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, withOptionals),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateVerify:        true,
				// The configuration omits secure_objects_only, so the state has "default" and the import (enforcing mode) the value of the share.
				ImportStateVerifyIgnore: []string{"secure_objects_only"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedShareResource(t, resourceId).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasSecureObjectsOnlyString("true").
						HasCommentString(comment).
						HasAccounts(account.Name()),
					resourceshowoutputassert.ImportedShareShowOutput(t, resourceId).
						HasName(id.Name()).
						HasKind(sdk.ShareKindOutbound).
						HasComment(comment).
						HasSecureObjectsOnly(true).
						HasToTotal(0),
				),
			},
			// 5. Update - unset optionals
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
						planchecks.ExpectSetElementsChange(ref, "accounts", tfjson.ActionUpdate, nil, []string{account.Name()}),
						planchecks.ExpectNoChangeOnField(ref, "secure_objects_only"),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// 6. Destroy
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, basic),
				Destroy:                  true,
			},
			// 7. Create - with all fields (including the ForceNew secure_objects_only)
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					_, err := testClient().Share.Show(t, id)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			// 8. Import - with all fields (the state equals the imported one, nothing is ignored)
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, complete),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateVerify:        true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedShareResource(t, resourceId).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasSecureObjectsOnlyString("false").
						HasCommentString(comment).
						HasAccounts(account.Name()),
					resourceshowoutputassert.ImportedShareShowOutput(t, resourceId).
						HasName(id.Name()).
						HasKind(sdk.ShareKindOutbound).
						HasComment(comment).
						HasSecureObjectsOnly(false).
						HasToTotal(0),
				),
			},
			// 9. Update - external changes of the updatable fields
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					testClient().Share.Alter(t, sdk.NewAlterShareRequest(id).WithSet(sdk.ShareSetRequest{Comment: sdk.String(externalComment)}))
					testClient().Share.Alter(t, sdk.NewAlterShareRequest(id).WithRemove(sdk.ShareRemoveRequest{Accounts: []sdk.AccountIdentifier{account}}))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
						planchecks.ExpectSetDrift(ref, "accounts", []string{account.Name()}, nil),
						planchecks.ExpectSetElementsChange(ref, "accounts", tfjson.ActionUpdate, []string{account.Name()}, nil),
						planchecks.ExpectNoChangeOnField(ref, "secure_objects_only"),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			// 10. Update - external change of the ForceNew secure_objects_only (the share is replaced outside of Terraform, so it is created again)
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					testClient().Share.CreateShareWithRequest(t, id, sdk.NewCreateShareRequest(id).WithOrReplace(true).WithSecureObjectsOnly(true))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "secure_objects_only", sdk.String("false"), sdk.String("true")),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

// TestAcc_Share_BasicUseCase_Import imports a share that is created outside of Terraform (enforcing mode, the default).
func TestAcc_Share_BasicUseCase_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)

	basic := model.Share("t", id.Name())
	ref := basic.ResourceReference()

	_, cleanup := testClient().Share.CreateShareWithIdentifier(t, id)
	t.Cleanup(cleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			// Import - without optionals
			{
				Config:             accconfig.FromModels(t, basic),
				ResourceName:       ref,
				ImportState:        true,
				ImportStateId:      resourceId,
				ImportStatePersist: true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedShareResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasSecureObjectsOnlyString("true").
						HasCommentEmpty().
						HasAccountsEmpty(),
				),
			},
			// secure_objects_only cannot be changed in place, so omitting it replaces the share; the plan after that is empty
			{
				Config: accconfig.FromModels(t, basic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(ref, "secure_objects_only", tfjson.ActionDelete, sdk.String("true"), sdk.String(r.BooleanDefault)),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(t,
					resourceassert.ShareResource(t, ref).
						HasNameString(id.Name()).
						HasSecureObjectsOnlyString(r.BooleanDefault).
						HasCommentEmpty().
						HasAccountsEmpty(),
				),
			},
		},
	})
}

// TestAcc_Share_SecureObjectsOnly covers secure_objects_only, which Snowflake cannot change on an existing share, so a
// change in the configuration or outside of Terraform recreates the share. BasicUseCase covers the external change when the
// field is set in the configuration; this test covers it when the field is not set, and the import in both modes.
func TestAcc_Share_SecureObjectsOnly(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	unset := model.Share("t", id.Name())
	secureObjectsOnlyFalse := model.Share("t", id.Name()).WithSecureObjectsOnlyBool(false)
	ref := unset.ResourceReference()

	replaceShare := func(secureObjectsOnly bool) {
		testClient().Share.CreateShareWithRequest(t, id, sdk.NewCreateShareRequest(id).WithOrReplace(true).WithSecureObjectsOnly(secureObjectsOnly))
	}
	assertShare := func(configured string, secureObjectsOnly bool) []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceassert.ShareResource(t, ref).
				HasNameString(id.Name()).
				HasSecureObjectsOnlyString(configured),
			resourceshowoutputassert.ShareShowOutput(t, ref).
				HasName(id.Name()).
				HasSecureObjectsOnly(secureObjectsOnly),
			objectassert.Share(t, id).
				HasName(id.Name()).
				HasSecureObjectsOnly(secureObjectsOnly),
		}
	}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			// Create - not set, so Snowflake uses its default (true)
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, unset),
				Check:                    assertThat(t, assertShare(r.BooleanDefault, true)...),
			},
			// External change - not set in the configuration, the share is replaced with false in Snowflake
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig:                func() { replaceShare(false) },
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "secure_objects_only", sdk.String(r.BooleanDefault), sdk.String("false")),
					},
				},
				Config: accconfig.FromModels(t, unset),
				Check:  assertThat(t, assertShare(r.BooleanDefault, true)...),
			},
			// Update - set to false in the configuration (recreate)
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, secureObjectsOnlyFalse),
				Check:  assertThat(t, assertShare("false", false)...),
			},
			// Import - enforcing mode
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, secureObjectsOnlyFalse),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateVerify:        true,
			},
			// Import - permissive mode (the configuration sets the value of the share)
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   accconfig.FromModels(t, providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode), secureObjectsOnlyFalse),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
			// External change - set in the configuration, the share is replaced with true in Snowflake
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig:                func() { replaceShare(true) },
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "secure_objects_only", sdk.String("false"), sdk.String("true")),
					},
				},
				Config: accconfig.FromModels(t, secureObjectsOnlyFalse),
				Check:  assertThat(t, assertShare("false", false)...),
			},
		},
	})
}

// TestAcc_Share_ExternalDeletion checks that a share dropped outside of Terraform is removed from the state and created again.
func TestAcc_Share_ExternalDeletion(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	basic := model.Share("t", id.Name())
	ref := basic.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, resourceassert.ShareResource(t, ref).HasNameString(id.Name())),
			},
			{
				PreConfig: testClient().Share.DropShareFunc(t, id),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check: assertThat(t,
					resourceassert.ShareResource(t, ref).HasNameString(id.Name()),
					objectassert.Share(t, id).HasName(id.Name()),
				),
			},
		},
	})
}

// TestAcc_Share_migrateFromV2_22_1 checks the state upgrader: the id changes from the plain name to the fully qualified name,
// secure_objects_only is filled with its default, and the share is not changed or replaced.
func TestAcc_Share_migrateFromV2_22_1(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	account := secondaryTestClient().Account.GetAccountIdentifier(t)
	comment := random.Comment()

	providerModel := providermodel.SnowflakeProvider().WithPreviewFeaturesEnabled(string(previewfeatures.ShareResource))
	complete := model.Share("t", id.Name()).WithComment(comment).WithAccounts(account.Name())
	ref := complete.ResourceReference()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.22.1"),
				Config:            accconfig.FromModels(t, providerModel, complete),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(ref, "id", id.Name())),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Config: accconfig.FromModels(t, providerModel, complete),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(ref, "id", id.FullyQualifiedName())),
					resourceassert.ShareResource(t, ref).
						HasNameString(id.Name()).
						HasSecureObjectsOnlyString(r.BooleanDefault).
						HasCommentString(comment).
						HasAccounts(account.Name()),
					objectassert.Share(t, id).
						HasName(id.Name()).
						HasSecureObjectsOnly(true).
						HasComment(comment).
						HasTo(account),
				),
			},
		},
	})
}

// TestAcc_Share_Regression_ThreeAccounts adds two recipients to a share and removes them again.
func TestAcc_Share_Regression_ThreeAccounts(t *testing.T) {
	// The prod environment has no setup with three accounts on the same deployment, which is needed to add two recipients.
	if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakeProdEnvironment {
		t.Skip("Skipping for the prod Snowflake environment: three accounts on the same deployment are required")
	}

	account2 := secondaryTestClient().Account.GetAccountIdentifier(t)
	account3 := thirdTestClient().Account.GetAccountIdentifier(t)

	comment := random.Comment()
	id := testClient().Ids.RandomAccountObjectIdentifier()

	noAccounts := model.Share("t", id.Name()).WithComment(comment)
	twoAccounts := model.Share("t", id.Name()).WithComment(comment).WithAccounts(account2.Name(), account3.Name())
	oneAccount := model.Share("t", id.Name()).WithComment(comment).WithAccounts(account2.Name())
	ref := noAccounts.ResourceReference()

	noAccountsAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasCommentString(comment).
			HasAccountsEmpty(),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasComment(comment).
			HasNoTo().
			HasToTotal(0),
	}
	twoAccountsAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasCommentString(comment).
			HasAccounts(account2.Name(), account3.Name()),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasComment(comment).
			HasTo(account2, account3).
			HasToTotal(0),
	}
	oneAccountAssertions := []assert.TestCheckFuncProvider{
		resourceassert.ShareResource(t, ref).
			HasNameString(id.Name()).
			HasCommentString(comment).
			HasAccounts(account2.Name()),
		objectassert.Share(t, id).
			HasName(id.Name()).
			HasComment(comment).
			HasTo(account2).
			HasToTotal(0),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			// Create - without accounts
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, noAccounts),
				Check:  assertThat(t, noAccountsAssertions...),
			},
			// Update - add two accounts
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectSetChange(ref, "accounts", tfjson.ActionUpdate, nil, []string{account2.Name(), account3.Name()}),
					},
				},
				Config: accconfig.FromModels(t, twoAccounts),
				Check:  assertThat(t, twoAccountsAssertions...),
			},
			// Update - remove one account
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectSetElementsChange(ref, "accounts", tfjson.ActionUpdate, nil, []string{account3.Name()}),
					},
				},
				Config: accconfig.FromModels(t, oneAccount),
				Check:  assertThat(t, oneAccountAssertions...),
			},
			// Update - a recipient added outside of the configuration is detected and removed
			{
				PreConfig: func() {
					// Recipients can be added only to a share that has a database, and the provider removes its temporary one.
					testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, testClient().Ids.DatabaseId(), id, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage})
					testClient().Share.Alter(t, sdk.NewAlterShareRequest(id).WithAdd(sdk.ShareAddRequest{
						Accounts:          []sdk.AccountIdentifier{account3},
						ShareRestrictions: sdk.Bool(false),
					}))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectSetElementsDrift(ref, "accounts", []string{account3.Name()}, nil),
						planchecks.ExpectSetElementsChange(ref, "accounts", tfjson.ActionUpdate, nil, []string{account3.Name()}),
					},
				},
				Config: accconfig.FromModels(t, oneAccount),
				Check:  assertThat(t, oneAccountAssertions...),
			},
			// Update - remove all accounts
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectSetChange(ref, "accounts", tfjson.ActionUpdate, []string{account2.Name()}, nil),
					},
				},
				Config: accconfig.FromModels(t, noAccounts),
				Check:  assertThat(t, noAccountsAssertions...),
			},
			// Import
			{
				Config:            accconfig.FromModels(t, noAccounts),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// The configuration omits secure_objects_only, so the state has "default" and the import the value of the share.
				ImportStateVerifyIgnore: []string{"secure_objects_only"},
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedShareResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasCommentString(comment).
						HasAccountsEmpty(),
				),
			},
		},
	})
}

func TestAcc_Share_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	incorrectAccount := model.Share("t", id.Name()).WithComment(comment).WithAccounts("incorrect")
	correctAndIncorrectAccounts := model.Share("t", id.Name()).WithComment(comment).WithAccounts("correct.one", "incorrect")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, incorrectAccount),
				ExpectError: regexp.MustCompile("Unable to parse the account identifier"),
			},
			{
				Config:      accconfig.FromModels(t, correctAndIncorrectAccounts),
				ExpectError: regexp.MustCompile("Unable to parse the account identifier"),
			},
		},
	})
}

// TestAcc_Share_Regression_Issue4398 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/4398 is fixed.
func TestAcc_Share_Regression_Issue4398(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	account := secondaryTestClient().Account.GetAccountIdentifier(t)
	id := testClient().Ids.RandomAccountObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().WithPreviewFeaturesEnabled(string(previewfeatures.ShareResource))
	noAccounts := model.Share("t", id.Name())
	oneAccount := model.Share("t", id.Name()).WithAccounts(account.Name())
	ref := noAccounts.ResourceReference()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			// Step 1: Share with no accounts
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            accconfig.FromModels(t, providerModel, noAccounts),
				Check: assertThat(t,
					resourceassert.ShareResource(t, ref).
						HasNameString(id.Name()).
						HasAccountsEmpty(),
					objectassert.Share(t, id).
						HasName(id.Name()).
						HasNoTo(),
				),
			},
			// Step 2: Grant USAGE on database to the share externally, then add accounts.
			// v2.13.0 fails because setShareAccounts creates a temp database and tries to
			// grant USAGE on it, conflicting with the already-granted USAGE on the real database.
			{
				PreConfig: func() {
					testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, database.ID(), id, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage})
				},
				ExternalProviders: ExternalProviderWithExactVersion("2.13.0"),
				Config:            accconfig.FromModels(t, providerModel, oneAccount),
				ExpectError:       regexp.MustCompile("does not belong to the database that is being shared"),
			},
			// Step 3: Succeeds after fix (skip temp database when share already has a database granted).
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   accconfig.FromModels(t, oneAccount),
				Check: assertThat(t,
					resourceassert.ShareResource(t, ref).
						HasNameString(id.Name()).
						HasAccounts(account.Name()),
					objectassert.Share(t, id).
						HasName(id.Name()).
						HasDatabaseName(database.ID()).
						HasTo(account),
				),
			},
		},
	})
}

// TestAcc_Share_Regression_WithGrantedDatabase proves that adding and updating accounts on a share
// that already has a database granted (via snowflake_grant_privileges_to_share) works correctly.
// This is the regression scenario reported after PR #4174 — previously, the provider would
// erroneously try to create a temp database even when a real DB was already in the share,
// causing a conflict on the USAGE grant.
// See: https://github.com/snowflakedb/terraform-provider-snowflake/pull/4174
func TestAcc_Share_Regression_WithGrantedDatabase(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	account := secondaryTestClient().Account.GetAccountIdentifier(t)
	id := testClient().Ids.RandomAccountObjectIdentifier()

	shareModel := model.ShareWithDefaultMeta(id.Name())
	shareModelWithAccount := model.ShareWithDefaultMeta(id.Name()).WithAccounts(account.Name())
	grantModel := model.GrantPrivilegesToShareWithDefaultMeta([]string{"USAGE"}, id.Name()).
		WithOnDatabase(database.ID().Name()).
		WithDependsOn(shareModel.ResourceReference())
	ref := shareModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Share),
		Steps: []resource.TestStep{
			// Step 1: Share with no accounts + grant DB to share (no accounts yet).
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, shareModel, grantModel),
				Check: assertThat(t,
					resourceassert.ShareResource(t, ref).
						HasNameString(id.Name()).
						HasAccountsEmpty(),
					objectassert.Share(t, id).
						HasName(id.Name()).
						HasDatabaseName(database.ID()).
						HasNoTo(),
				),
			},
			// Step 2: Add an account while DB is already on the share.
			// This exercises the UpdateShare path that directly adds accounts
			// (skipping the temp DB workaround) when share.DatabaseName is set.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, shareModelWithAccount, grantModel),
				Check: assertThat(t,
					resourceassert.ShareResource(t, ref).
						HasNameString(id.Name()).
						HasAccounts(account.Name()),
					objectassert.Share(t, id).
						HasName(id.Name()).
						HasDatabaseName(database.ID()).
						HasTo(account),
				),
			},
			// Step 3: Remove account (DB still granted).
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, shareModel, grantModel),
				Check: assertThat(t,
					resourceassert.ShareResource(t, ref).
						HasNameString(id.Name()).
						HasAccountsEmpty(),
					objectassert.Share(t, id).
						HasName(id.Name()).
						HasDatabaseName(database.ID()).
						HasNoTo(),
				),
			},
		},
	})
}
