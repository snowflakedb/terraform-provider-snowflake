//go:build non_account_level_tests

package testacc

import (
	"maps"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func catalogLinkedDatabaseShowOutputAssertions(t *testing.T, ref string, name string, comment string) assert.TestCheckFuncProvider {
	t.Helper()
	return resourceshowoutputassert.DatabaseShowOutput(t, ref).
		HasCreatedOnNotEmpty().
		HasName(name).
		HasIsDefault(false).
		HasOwner("ACCOUNTADMIN").
		HasComment(comment).
		HasTransient(false).
		HasKind(sdk.DatabaseKindCatalogLinkedDatabase).
		HasOwnerRoleType("ROLE")
}

// The tests need a preconfigured external Iceberg REST catalog integration and external volume, so they
// run only with the TEST_SF_TF_CATALOG_LINKED_DATABASE_* environment variables set.
// TODO(SNOW-3725859): Provide them dynamically and unskip these tests in CI.
// TODO(next PRs): Cover the variant without external_volume (credential-vending catalogs).
func TestAcc_CatalogLinkedDatabase_BasicUseCase(t *testing.T) {
	restCatalogId := sdk.NewAccountObjectIdentifier(testenvs.GetOrSkipTest(t, testenvs.CatalogLinkedDatabaseCatalogIntegration))
	externalVolumeId := sdk.NewAccountObjectIdentifier(testenvs.GetOrSkipTest(t, testenvs.CatalogLinkedDatabaseExternalVolume))

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	basicLinkedCatalog := []sdk.LinkedCatalogRequest{
		*sdk.NewLinkedCatalogRequest().WithCatalog(restCatalogId),
	}
	completeLinkedCatalog := []sdk.LinkedCatalogRequest{
		*sdk.NewLinkedCatalogRequest().
			WithCatalog(restCatalogId).
			WithAllowedNamespaces([]sdk.StringListItemWrapper{{Value: "ns1"}, {Value: "ns2"}}).
			WithAllowedWriteOperations(sdk.CatalogLinkedDatabaseAllowedWriteOperationsAll).
			WithSyncIntervalSeconds(60),
	}

	basic := model.CatalogLinkedDatabaseWithDefaultMeta(id.Name(), basicLinkedCatalog).
		WithExternalVolume(externalVolumeId.Name())

	complete := model.CatalogLinkedDatabaseWithDefaultMeta(id.Name(), completeLinkedCatalog).
		WithExternalVolume(externalVolumeId.Name()).
		WithComment(comment)

	ref := basic.ResourceReference()

	catalogLinkedDatabaseAssertions := func(allowedNamespaces []string, syncIntervalSeconds int, comment string) []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceassert.CatalogLinkedDatabaseResource(t, ref).
				HasNameString(id.Name()).
				HasFullyQualifiedNameString(id.FullyQualifiedName()).
				HasLinkedCatalogCatalog(restCatalogId).
				HasLinkedCatalogAllowedNamespaces(allowedNamespaces...).
				HasLinkedCatalogBlockedNamespaces().
				HasLinkedCatalogAllowedWriteOperations(sdk.CatalogLinkedDatabaseAllowedWriteOperationsAll).
				HasLinkedCatalogNamespaceMode(sdk.CatalogLinkedDatabaseNamespaceModeIgnoreNestedNamespace).
				HasLinkedCatalogNamespaceFlattenDelimiter("").
				HasLinkedCatalogSyncIntervalSeconds(syncIntervalSeconds).
				HasExternalVolumeString(externalVolumeId.Name()).
				HasCatalogCaseSensitivityString(string(sdk.DatabaseCatalogCaseSensitivityCaseInsensitive)).
				HasCommentString(comment),
			catalogLinkedDatabaseShowOutputAssertions(t, ref, id.Name(), comment),
		}
	}

	// Defaults come from the docs; they may need adjusting after a run against a live external catalog.
	basicAssertions := catalogLinkedDatabaseAssertions(nil, 30, "")
	completeAssertions := catalogLinkedDatabaseAssertions([]string{"ns1", "ns2"}, 60, comment)
	// Snowflake exposes no UNSET for allowed_write_operations and sync_interval_seconds, so they keep the last values.
	backToBasicAssertions := catalogLinkedDatabaseAssertions(nil, 60, "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogLinkedDatabase),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import - without optionals
			{
				Config:            accconfig.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, basic),
				Check:  assertThat(t, backToBasicAssertions...),
			},
		},
	})
}

// The ForceNew attributes can only be exercised at create, hence a separate test.
func TestAcc_CatalogLinkedDatabase_CompleteUseCase(t *testing.T) {
	restCatalogId := sdk.NewAccountObjectIdentifier(testenvs.GetOrSkipTest(t, testenvs.CatalogLinkedDatabaseCatalogIntegration))
	externalVolumeId := sdk.NewAccountObjectIdentifier(testenvs.GetOrSkipTest(t, testenvs.CatalogLinkedDatabaseExternalVolume))

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	completeLinkedCatalog := []sdk.LinkedCatalogRequest{
		*sdk.NewLinkedCatalogRequest().
			WithCatalog(restCatalogId).
			WithAllowedNamespaces([]sdk.StringListItemWrapper{{Value: "ns1"}, {Value: "ns2"}}).
			WithBlockedNamespaces([]sdk.StringListItemWrapper{{Value: "ns3"}}).
			WithAllowedWriteOperations(sdk.CatalogLinkedDatabaseAllowedWriteOperationsAll).
			WithNamespaceMode(sdk.CatalogLinkedDatabaseNamespaceModeFlattenNestedNamespace).
			WithNamespaceFlattenDelimiter("_").
			WithSyncIntervalSeconds(60),
	}

	complete := model.CatalogLinkedDatabaseWithDefaultMeta(id.Name(), completeLinkedCatalog).
		WithExternalVolume(externalVolumeId.Name()).
		WithCatalogCaseSensitivity(string(sdk.DatabaseCatalogCaseSensitivityCaseInsensitive)).
		WithComment(comment)

	ref := complete.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogLinkedDatabase),
		Steps: []resource.TestStep{
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: accconfig.FromModels(t, complete),
				Check: assertThat(t,
					resourceassert.CatalogLinkedDatabaseResource(t, ref).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasLinkedCatalogCatalog(restCatalogId).
						HasLinkedCatalogAllowedNamespaces("ns1", "ns2").
						HasLinkedCatalogBlockedNamespaces("ns3").
						HasLinkedCatalogAllowedWriteOperations(sdk.CatalogLinkedDatabaseAllowedWriteOperationsAll).
						HasLinkedCatalogNamespaceMode(sdk.CatalogLinkedDatabaseNamespaceModeFlattenNestedNamespace).
						HasLinkedCatalogNamespaceFlattenDelimiter("_").
						HasLinkedCatalogSyncIntervalSeconds(60).
						HasExternalVolumeString(externalVolumeId.Name()).
						HasCatalogCaseSensitivityString(string(sdk.DatabaseCatalogCaseSensitivityCaseInsensitive)).
						HasCommentString(comment),
					catalogLinkedDatabaseShowOutputAssertions(t, ref, id.Name(), comment),
				),
			},
			{
				Config:            accconfig.FromModels(t, complete),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_CatalogLinkedDatabase_Rename(t *testing.T) {
	restCatalogId := sdk.NewAccountObjectIdentifier(testenvs.GetOrSkipTest(t, testenvs.CatalogLinkedDatabaseCatalogIntegration))
	externalVolumeId := sdk.NewAccountObjectIdentifier(testenvs.GetOrSkipTest(t, testenvs.CatalogLinkedDatabaseExternalVolume))

	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()

	linkedCatalog := []sdk.LinkedCatalogRequest{
		*sdk.NewLinkedCatalogRequest().WithCatalog(restCatalogId),
	}

	initial := model.CatalogLinkedDatabaseWithDefaultMeta(id.Name(), linkedCatalog).
		WithExternalVolume(externalVolumeId.Name())
	renamed := model.CatalogLinkedDatabaseWithDefaultMeta(newId.Name(), linkedCatalog).
		WithExternalVolume(externalVolumeId.Name())

	ref := initial.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogLinkedDatabase),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, initial),
				Check: assertThat(t,
					resourceassert.CatalogLinkedDatabaseResource(t, ref).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// Renaming is handled in place through ALTER ... RENAME TO
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, renamed),
				Check: assertThat(t,
					resourceassert.CatalogLinkedDatabaseResource(t, ref).
						HasNameString(newId.Name()).
						HasFullyQualifiedNameString(newId.FullyQualifiedName()),
				),
			},
		},
	})
}

// Plan-time only: needs no Snowflake object and no external catalog fixtures, so it runs in CI.
func TestAcc_CatalogLinkedDatabase_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	linkedCatalogValue := func(extra map[string]tfconfig.Variable) tfconfig.Variable {
		m := map[string]tfconfig.Variable{
			"catalog": tfconfig.StringVariable("SOME_CATALOG"),
		}
		maps.Copy(m, extra)
		return tfconfig.ListVariable(tfconfig.ObjectVariable(m))
	}
	modelWithLinkedCatalog := func(extra map[string]tfconfig.Variable) *model.CatalogLinkedDatabaseModel {
		return model.CatalogLinkedDatabaseWithDefaultMeta(id.Name(), nil).
			WithLinkedCatalogValue(linkedCatalogValue(extra))
	}

	invalidWriteOperations := modelWithLinkedCatalog(map[string]tfconfig.Variable{"allowed_write_operations": tfconfig.StringVariable("INVALID")})
	invalidNamespaceMode := modelWithLinkedCatalog(map[string]tfconfig.Variable{"namespace_mode": tfconfig.StringVariable("INVALID")})
	invalidCaseSensitivity := modelWithLinkedCatalog(nil).WithCatalogCaseSensitivity("INVALID")
	invalidCatalogId := modelWithLinkedCatalog(map[string]tfconfig.Variable{"catalog": tfconfig.StringVariable("a.b.c")})
	invalidExternalVolumeId := modelWithLinkedCatalog(nil).WithExternalVolume("a.b.c")
	invalidSyncInterval := modelWithLinkedCatalog(map[string]tfconfig.Variable{"sync_interval_seconds": tfconfig.IntegerVariable(0)})
	delimiterWithoutFlattenMode := modelWithLinkedCatalog(map[string]tfconfig.Variable{"namespace_flatten_delimiter": tfconfig.StringVariable("_")})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogLinkedDatabase),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, invalidWriteOperations),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("invalid catalog linked database allowed write operations: INVALID"),
			},
			{
				Config:      accconfig.FromModels(t, invalidNamespaceMode),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("invalid catalog linked database namespace mode: INVALID"),
			},
			{
				Config:      accconfig.FromModels(t, invalidCaseSensitivity),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("invalid database catalog case sensitivity: INVALID"),
			},
			{
				Config:      accconfig.FromModels(t, invalidCatalogId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid identifier type"),
			},
			{
				Config:      accconfig.FromModels(t, invalidExternalVolumeId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid identifier type"),
			},
			{
				Config:      accconfig.FromModels(t, invalidSyncInterval),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected linked_catalog\.0\.sync_interval_seconds to be at least \(1\), got 0`),
			},
			{
				Config:      accconfig.FromModels(t, delimiterWithoutFlattenMode),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("namespace_flatten_delimiter can only be set when namespace_mode is FLATTEN_NESTED_NAMESPACE"),
			},
		},
	})
}

// Needs no external catalog fixtures, so it runs in CI.
func TestAcc_CatalogLinkedDatabase_Import_WrongKind(t *testing.T) {
	// A standard database created outside of Terraform is used as the import target.
	standardDatabase, standardDatabaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(standardDatabaseCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	catalogLinkedDatabaseModel := model.CatalogLinkedDatabaseWithDefaultMeta(id.Name(), []sdk.LinkedCatalogRequest{
		*sdk.NewLinkedCatalogRequest().WithCatalog(sdk.NewAccountObjectIdentifier("SOME_CATALOG")),
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogLinkedDatabase),
		Steps: []resource.TestStep{
			{
				Config:        accconfig.FromModels(t, catalogLinkedDatabaseModel),
				ResourceName:  catalogLinkedDatabaseModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: standardDatabase.ID().Name(),
				ExpectError:   regexp.MustCompile("is not a catalog-linked database"),
			},
		},
	})
}
