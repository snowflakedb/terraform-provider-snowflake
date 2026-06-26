//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_ApiIntegrationGitRepositoryPrivateLink_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	apiProvider := string(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi)

	comment := random.Comment()
	externalComment := random.Comment()

	certSecretId := testClient().Ids.RandomSchemaObjectIdentifier()
	_, cleanupCert := testClient().Secret.CreateWithGenericString(t, certSecretId, random.GenerateX509(t))
	t.Cleanup(cleanupCert)

	basic := model.ApiIntegrationGitRepositoryPrivateLink("t", id.Name(), []string{gitAllowedPrefix}, true, true)
	withOptionals := model.ApiIntegrationGitRepositoryPrivateLink("t", id.Name(), []string{gitAllowedPrefix}, true, true).
		WithAllAllowedAuthenticationSecrets(true).
		WithApiBlockedPrefixes([]string{gitBlockedPrefix}).
		WithTlsTrustedCertificates([]string{certSecretId.FullyQualifiedName()}).
		WithComment(comment)

	ref := basic.ResourceReference()

	basicCommonAsserts := func() []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
				HasName(id.Name()).
				HasEnabled(true).
				HasComment(""),
			resourceshowoutputassert.ApiIntegrationGitRepositoryPrivateLinkDescribeOutput(t, ref).
				HasEnabled(true).
				HasApiProvider(apiProvider).
				HasAllowedAuthenticationSecrets("").
				HasUsePrivatelinkEndpoint(true).
				HasNoTlsTrustedCertificates().
				HasAllowedPrefixes(gitAllowedPrefix).
				HasNoBlockedPrefixes().
				HasComment(""),
			objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
				HasEnabled(true).
				HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
				HasNoUserAuthType().
				HasAllowedAuthenticationSecrets("").
				HasUsePrivatelinkEndpoint(true).
				HasAllowedPrefixes(gitAllowedPrefix).
				HasNoBlockedPrefixes().
				HasNoTlsTrustedCertificates().
				HasComment(""),
		}
	}
	basicResourceBaseAsserts := func() *resourceassert.ApiIntegrationGitRepositoryPrivateLinkResourceAssert {
		return resourceassert.ApiIntegrationGitRepositoryPrivateLinkResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasUsePrivatelinkEndpointString(r.BooleanTrue).
			HasNoNoAllowedAuthenticationSecrets().
			HasAllowedAuthenticationSecretsEmpty().
			HasApiBlockedPrefixesEmpty().
			HasTlsTrustedCertificatesEmpty().
			HasCommentEmpty()
	}
	assertBasic := append(
		[]assert.TestCheckFuncProvider{
			basicResourceBaseAsserts().HasNoAllAllowedAuthenticationSecrets(),
		},
		basicCommonAsserts()...,
	)
	assertBasicAfterUpdate := append(
		[]assert.TestCheckFuncProvider{
			basicResourceBaseAsserts().HasAllAllowedAuthenticationSecrets(false),
		},
		basicCommonAsserts()...,
	)

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.ApiIntegrationGitRepositoryPrivateLinkResource(t, ref).
			HasNameString(id.Name()).
			HasEnabledString(r.BooleanTrue).
			HasUsePrivatelinkEndpointString(r.BooleanTrue).
			HasAllAllowedAuthenticationSecrets(true).
			HasNoNoAllowedAuthenticationSecrets().
			HasAllowedAuthenticationSecretsEmpty().
			HasApiBlockedPrefixes(gitBlockedPrefix).
			HasTlsTrustedCertificates(certSecretId.FullyQualifiedName()).
			HasCommentString(comment),
		resourceshowoutputassert.ApiIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.ApiIntegrationGitRepositoryPrivateLinkDescribeOutput(t, ref).
			HasEnabled(true).
			HasApiProvider(apiProvider).
			HasAllowedAuthenticationSecrets("").
			HasUsePrivatelinkEndpoint(true).
			HasTlsTrustedCertificates(fmt.Sprintf(`"%s"."%s".%s`, certSecretId.DatabaseName(), certSecretId.SchemaName(), certSecretId.Name())).
			HasAllowedPrefixes(gitAllowedPrefix).
			HasBlockedPrefixes(gitBlockedPrefix).
			HasComment(comment),
		objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets("").
			HasUsePrivatelinkEndpoint(true).
			HasAllowedPrefixes(gitAllowedPrefix).
			HasBlockedPrefixes(gitBlockedPrefix).
			HasTlsTrustedCertificates(fmt.Sprintf(`"%s"."%s".%s`, certSecretId.DatabaseName(), certSecretId.SchemaName(), certSecretId.Name())).
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryPrivateLink),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals (basic has no auth secrets, so import round-trips cleanly)
			{
				Config:            config.FromModels(t, basic),
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
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// Import - with optionals (auth secrets cannot be read back from Snowflake, so ignore those fields)
			{
				Config:                  config.FromModels(t, withOptionals),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"all_allowed_authentication_secrets", "no_allowed_authentication_secrets", "allowed_authentication_secrets"},
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasicAfterUpdate...),
			},
			// Update - external changes
			{
				PreConfig: func() {
					testClient().ApiIntegration.Alter(t, sdk.NewAlterApiIntegrationRequest(id).WithSet(
						*sdk.NewApiIntegrationSetRequest().WithComment(externalComment),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasicAfterUpdate...),
			},
			// Destroy
			{
				Destroy: true,
				Config:  config.FromModels(t, basic),
			},
			// Create - with optionals
			{
				PreConfig: func() {
					_, err := testClient().ApiIntegration.Show(t, id)
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
		},
	})
}

// TestAcc_ApiIntegrationGitRepositoryPrivateLink_AllowedSecrets_Update covers all directed transitions
// between the four valid states: not-set, ALL, NONE, and a specific secrets list.
func TestAcc_ApiIntegrationGitRepositoryPrivateLink_AllowedSecrets_Update(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	secretId, cleanupSecret := testClient().Secret.CreateRandomPasswordSecret(t)
	t.Cleanup(cleanupSecret)

	withNone := model.ApiIntegrationGitRepositoryPrivateLink("t", id.Name(), []string{gitAllowedPrefix}, true, true).
		WithNoAllowedAuthenticationSecrets(true)
	withAll := model.ApiIntegrationGitRepositoryPrivateLink("t", id.Name(), []string{gitAllowedPrefix}, true, true).
		WithAllAllowedAuthenticationSecrets(true)
	withList := model.ApiIntegrationGitRepositoryPrivateLink("t", id.Name(), []string{gitAllowedPrefix}, true, true).
		WithAllowedAuthenticationSecrets([]string{secretId.FullyQualifiedName()})
	withNotSet := model.ApiIntegrationGitRepositoryPrivateLink("t", id.Name(), []string{gitAllowedPrefix}, true, true)

	ref := withNotSet.ResourceReference()

	assertNotSet := func() []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceassert.ApiIntegrationGitRepositoryPrivateLinkResource(t, ref).
				HasNoAllAllowedAuthenticationSecrets().
				HasNoNoAllowedAuthenticationSecrets().
				HasAllowedAuthenticationSecretsEmpty(),
			resourceshowoutputassert.ApiIntegrationGitRepositoryPrivateLinkDescribeOutput(t, ref).
				HasAllowedAuthenticationSecrets(""),
		}
	}
	assertNotSetAfterTransition := func() []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceassert.ApiIntegrationGitRepositoryPrivateLinkResource(t, ref).
				HasAllAllowedAuthenticationSecrets(false).
				HasNoAllowedAuthenticationSecrets(false).
				HasAllowedAuthenticationSecretsEmpty(),
			resourceshowoutputassert.ApiIntegrationGitRepositoryPrivateLinkDescribeOutput(t, ref).
				HasAllowedAuthenticationSecrets(""),
		}
	}
	assertAll := func() []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceassert.ApiIntegrationGitRepositoryPrivateLinkResource(t, ref).
				HasAllAllowedAuthenticationSecrets(true).
				HasNoNoAllowedAuthenticationSecrets().
				HasAllowedAuthenticationSecretsEmpty(),
			resourceshowoutputassert.ApiIntegrationGitRepositoryPrivateLinkDescribeOutput(t, ref).
				HasAllowedAuthenticationSecrets(""),
		}
	}
	assertAllAfterTransition := func() []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceassert.ApiIntegrationGitRepositoryPrivateLinkResource(t, ref).
				HasAllAllowedAuthenticationSecrets(true).
				HasNoAllowedAuthenticationSecrets(false).
				HasAllowedAuthenticationSecretsEmpty(),
			resourceshowoutputassert.ApiIntegrationGitRepositoryPrivateLinkDescribeOutput(t, ref).
				HasAllowedAuthenticationSecrets(""),
		}
	}
	assertNone := func() []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceassert.ApiIntegrationGitRepositoryPrivateLinkResource(t, ref).
				HasAllAllowedAuthenticationSecrets(false).
				HasNoAllowedAuthenticationSecrets(true).
				HasAllowedAuthenticationSecretsEmpty(),
			resourceshowoutputassert.ApiIntegrationGitRepositoryPrivateLinkDescribeOutput(t, ref).
				HasAllowedAuthenticationSecrets(""),
		}
	}
	assertList := func() []assert.TestCheckFuncProvider {
		return []assert.TestCheckFuncProvider{
			resourceassert.ApiIntegrationGitRepositoryPrivateLinkResource(t, ref).
				HasAllAllowedAuthenticationSecrets(false).
				HasNoAllowedAuthenticationSecrets(false).
				HasAllowedAuthenticationSecrets(secretId.FullyQualifiedName()),
			resourceshowoutputassert.ApiIntegrationGitRepositoryPrivateLinkDescribeOutput(t, ref).
				HasAllowedAuthenticationSecrets(""),
		}
	}

	expectUpdate := func(m *model.ApiIntegrationGitRepositoryPrivateLinkModel, checks func() []assert.TestCheckFuncProvider) resource.TestStep {
		return resource.TestStep{
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
				},
			},
			Config: config.FromModels(t, m),
			Check:  assertThat(t, checks()...),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryPrivateLink),
		Steps: []resource.TestStep{
			// Create with not-set state
			{
				Config: config.FromModels(t, withNotSet),
				Check:  assertThat(t, assertNotSet()...),
			},
			// not-set → ALL
			expectUpdate(withAll, assertAll),
			// ALL → NONE
			expectUpdate(withNone, assertNone),
			// NONE → list
			expectUpdate(withList, assertList),
			// list → ALL
			expectUpdate(withAll, assertAllAfterTransition),
			// ALL → not-set
			expectUpdate(withNotSet, assertNotSetAfterTransition),
			// not-set → NONE
			expectUpdate(withNone, assertNone),
			// NONE → not-set
			expectUpdate(withNotSet, assertNotSetAfterTransition),
			// not-set → list
			expectUpdate(withList, assertList),
			// list → NONE
			expectUpdate(withNone, assertNone),
			// NONE → ALL
			expectUpdate(withAll, assertAllAfterTransition),
			// ALL → list
			expectUpdate(withList, assertList),
			// list → not-set
			expectUpdate(withNotSet, assertNotSetAfterTransition),
		},
	})
}

// TestAcc_ApiIntegrationGitRepositoryPrivateLink_Import verifies that importing a resource created outside Terraform
// produces no destroy-before-create plan.
func TestAcc_ApiIntegrationGitRepositoryPrivateLink_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()

	testModel := model.ApiIntegrationGitRepositoryPrivateLink("t", id.Name(), []string{gitAllowedPrefix}, true, true).
		WithApiBlockedPrefixes([]string{gitBlockedPrefix}).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryPrivateLink),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, cleanup := testClient().ApiIntegration.CreateWithRequest(
						t,
						sdk.NewCreateApiIntegrationRequest(id,
							[]sdk.ApiIntegrationEndpointPrefix{{Path: gitAllowedPrefix}}, true).
							WithComment(comment).
							WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitBlockedPrefix}}).
							WithGitHttpsApiPrivateLinkProviderParams(*sdk.NewGitHttpsApiPrivateLinkParamsRequest(true)),
					)
					t.Cleanup(cleanup)
				},
				Config:             config.FromModels(t, testModel),
				ResourceName:       testModel.ResourceReference(),
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			{
				Config: config.FromModels(t, testModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func TestAcc_ApiIntegrationGitRepositoryPrivateLink_Import_WrongProviderType(t *testing.T) {
	awsIntegration, awsCleanup := testClient().ApiIntegration.CreateAws(t)
	t.Cleanup(awsCleanup)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	dummyModel := model.ApiIntegrationGitRepositoryPrivateLink("t", id.Name(),
		[]string{gitAllowedPrefix}, true, true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ApiIntegrationGitRepositoryPrivateLink),
		Steps: []resource.TestStep{
			{
				Config:        config.FromModels(t, dummyModel),
				ResourceName:  dummyModel.ResourceReference(),
				ImportState:   true,
				ImportStateId: awsIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile("not compatible with snowflake_api_integration_git_repository_private_link"),
			},
		},
	})
}
