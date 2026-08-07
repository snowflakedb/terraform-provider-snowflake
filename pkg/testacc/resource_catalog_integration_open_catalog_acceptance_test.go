//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
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

func TestAcc_CatalogIntegrationOpenCatalog_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountLocator := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogAccountLocator)
	oAuthClientId := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientId)
	oAuthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientSecret)
	newOAuthClientId := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogSecondaryOAuthClientId)
	newOAuthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogSecondaryOAuthClientSecret)
	catalogUri := fmt.Sprintf("https://%s.snowflakecomputing.com/polaris/api/catalog", accountLocator)
	privateCatalogUri := fmt.Sprintf("https://%s.privatelink.snowflakecomputing.com/polaris/api/catalog", accountLocator)
	catalogName := "TEST_CATALOG"
	newCatalogName := "TEST_CATALOG2"

	catalogNamespace := "TEST_NAMESPACE"
	newCatalogNamespace := "TEST_NAMESPACE2"
	externalCatalogNamespace := "TEST_NAMESPACE3"

	oAuthTokenUri := catalogUri + "/v1/oauth/tokens"
	privateOAuthTokenUri := privateCatalogUri + "/v1/oauth/tokens"
	oAuthAllowedScope := "PRINCIPAL_ROLE:PRINCIPAL"
	additionalOAuthAllowedScope := "PRINCIPAL_ROLE:SECONDARY"

	oauthClientSecretVariableName := "oauth_client_secret"
	oauthClientSecretVariableModel, oauthClientSecretConfigVariables := config.SecretStringVariableModelWithConfigVariables(oauthClientSecretVariableName, oAuthClientSecret)
	_, oauthClientSecretConfigVariablesSecondary := config.SecretStringVariableModelWithConfigVariables(oauthClientSecretVariableName, newOAuthClientSecret)

	comment := random.Comment()
	externalComment := random.Comment()

	refreshIntervalSeconds := random.IntRange(30, 86400)
	externalRefreshIntervalSeconds := random.IntRange(30, 86400)

	basicRestAuth := []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, oAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}),
	}
	completeRestAuth := []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, oAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}).
			WithOauthTokenUri(oAuthTokenUri),
	}
	changedRestAuth := []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest(newOAuthClientId, newOAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}, {Value: additionalOAuthAllowedScope}}).
			WithOauthTokenUri(privateOAuthTokenUri),
	}

	basicRestConfig := []sdk.OpenCatalogRestConfigRequest{
		*sdk.NewOpenCatalogRestConfigRequest(catalogUri, catalogName),
	}
	completeRestConfig := []sdk.OpenCatalogRestConfigRequest{
		*sdk.NewOpenCatalogRestConfigRequest(catalogUri, catalogName).
			WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			WithAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
	}
	changedRestConfig := []sdk.OpenCatalogRestConfigRequest{
		*sdk.NewOpenCatalogRestConfigRequest(privateCatalogUri, newCatalogName).
			WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePrivate).
			WithAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials),
	}

	basic := model.CatalogIntegrationOpenCatalogVar("t", id.Name(), false, basicRestAuth, basicRestConfig, oauthClientSecretVariableName)

	altered := model.CatalogIntegrationOpenCatalogVar("t", id.Name(), true, basicRestAuth, basicRestConfig, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds)

	allAttributes := model.CatalogIntegrationOpenCatalogVar("t", id.Name(), false, completeRestAuth, completeRestConfig, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(catalogNamespace)

	withChangedCatalogNamespace := model.CatalogIntegrationOpenCatalogVar("t", id.Name(), false, completeRestAuth, completeRestConfig, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withChangedRestConfig := model.CatalogIntegrationOpenCatalogVar("t", id.Name(), false, basicRestAuth, changedRestConfig, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withChangedOAuthClientSecret := model.CatalogIntegrationOpenCatalogVar("t", id.Name(), false, []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, newOAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}),
	}, changedRestConfig, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withChangedRestAuth := model.CatalogIntegrationOpenCatalogVar("t", id.Name(), false, changedRestAuth, changedRestConfig, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypePolaris)).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasCatalogNamespaceEmpty().
			HasRestConfig(&sdk.OpenCatalogRestConfigDetails{
				CatalogUri:           catalogUri,
				CatalogApiType:       "",
				CatalogName:          catalogName,
				AccessDelegationMode: "",
			}).
			HasRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      "",
				OauthClientId:      oAuthClientId,
				OauthClientSecret:  oAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope},
			}),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(""),
		resourceshowoutputassert.CatalogIntegrationOpenCatalogDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypePolaris).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasCatalogNamespace(""),
		resourceshowoutputassert.CatalogIntegrationOpenCatalogDescribeOutput(t, ref).
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasRestConfigCatalogName(catalogName).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials).
			HasRestAuthenticationOauthTokenUri(catalogUri + "/v1/oauth/tokens").
			HasRestAuthenticationOauthClientId(oAuthClientId).
			HasRestAuthenticationOauthAllowedScopes(oAuthAllowedScope),
	}

	basicAssertionsWithRefreshIntervalZero := append(
		[]assert.TestCheckFuncProvider{
			resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
				HasName(id.Name()).
				HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypePolaris)).
				HasEnabledString(r.BooleanFalse).
				HasCommentEmpty().
				HasRefreshIntervalSeconds(0).
				HasCatalogNamespaceEmpty().
				HasRestConfig(&sdk.OpenCatalogRestConfigDetails{
					CatalogUri:           catalogUri,
					CatalogApiType:       "",
					CatalogName:          catalogName,
					AccessDelegationMode: "",
				}).
				HasRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
					OauthTokenUri:      "",
					OauthClientId:      oAuthClientId,
					OauthClientSecret:  oAuthClientSecret,
					OauthAllowedScopes: []string{oAuthAllowedScope},
				}),
		},
		basicAssertions[1:]...,
	)

	alteredProperties := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypePolaris)).
			HasEnabledString(r.BooleanTrue).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespaceEmpty().
			HasRestConfig(&sdk.OpenCatalogRestConfigDetails{
				CatalogUri:           catalogUri,
				CatalogApiType:       "",
				CatalogName:          catalogName,
				AccessDelegationMode: "",
			}).
			HasRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      "",
				OauthClientId:      oAuthClientId,
				OauthClientSecret:  oAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope},
			}),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationOpenCatalogDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypePolaris).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(true).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(""),
		resourceshowoutputassert.CatalogIntegrationOpenCatalogDescribeOutput(t, ref).
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasRestConfigCatalogName(catalogName).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials).
			HasRestAuthenticationOauthTokenUri(catalogUri + "/v1/oauth/tokens").
			HasRestAuthenticationOauthClientId(oAuthClientId).
			HasRestAuthenticationOauthAllowedScopes(oAuthAllowedScope),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypePolaris)).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(catalogNamespace).
			HasRestConfig(&sdk.OpenCatalogRestConfigDetails{
				CatalogUri:           catalogUri,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePublic,
				CatalogName:          catalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials,
			}).
			HasRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      oAuthTokenUri,
				OauthClientId:      oAuthClientId,
				OauthClientSecret:  oAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope},
			}),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationOpenCatalogDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypePolaris).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(catalogNamespace),
		resourceshowoutputassert.CatalogIntegrationOpenCatalogDescribeOutput(t, ref).
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasRestConfigCatalogName(catalogName).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials).
			HasRestAuthenticationOauthTokenUri(oAuthTokenUri).
			HasRestAuthenticationOauthClientId(oAuthClientId).
			HasRestAuthenticationOauthAllowedScopes(oAuthAllowedScope),
	}

	forceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypePolaris)).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(newCatalogNamespace).
			HasRestConfig(&sdk.OpenCatalogRestConfigDetails{
				CatalogUri:           catalogUri,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePublic,
				CatalogName:          catalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials,
			}).
			HasRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      oAuthTokenUri,
				OauthClientId:      oAuthClientId,
				OauthClientSecret:  oAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope},
			}),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationOpenCatalogDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypePolaris).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(newCatalogNamespace),
		resourceshowoutputassert.CatalogIntegrationOpenCatalogDescribeOutput(t, ref).
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasRestConfigCatalogName(catalogName).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials).
			HasRestAuthenticationOauthTokenUri(oAuthTokenUri).
			HasRestAuthenticationOauthClientId(oAuthClientId).
			HasRestAuthenticationOauthAllowedScopes(oAuthAllowedScope),
	}

	moreForceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypePolaris)).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(newCatalogNamespace).
			HasRestConfig(&sdk.OpenCatalogRestConfigDetails{
				CatalogUri:           privateCatalogUri,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePrivate,
				CatalogName:          newCatalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeVendedCredentials,
			}).
			HasRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      "",
				OauthClientId:      oAuthClientId,
				OauthClientSecret:  oAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope},
			}),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationOpenCatalogDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypePolaris).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(newCatalogNamespace),
		resourceshowoutputassert.CatalogIntegrationOpenCatalogDescribeOutput(t, ref).
			HasRestConfigCatalogUri(privateCatalogUri).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePrivate).
			HasRestConfigCatalogName(newCatalogName).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials).
			HasRestAuthenticationOauthTokenUri(privateOAuthTokenUri).
			HasRestAuthenticationOauthClientId(oAuthClientId).
			HasRestAuthenticationOauthAllowedScopes(oAuthAllowedScope),
	}

	moreForceNewAssertionsWithChangedSecret := append(
		[]assert.TestCheckFuncProvider{
			resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
				HasName(id.Name()).
				HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypePolaris)).
				HasEnabledString(r.BooleanFalse).
				HasComment(comment).
				HasRefreshIntervalSeconds(refreshIntervalSeconds).
				HasCatalogNamespace(newCatalogNamespace).
				HasRestConfig(&sdk.OpenCatalogRestConfigDetails{
					CatalogUri:           privateCatalogUri,
					CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePrivate,
					CatalogName:          newCatalogName,
					AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeVendedCredentials,
				}).
				HasRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
					OauthTokenUri:      "",
					OauthClientId:      oAuthClientId,
					OauthClientSecret:  newOAuthClientSecret,
					OauthAllowedScopes: []string{oAuthAllowedScope},
				}),
		},
		moreForceNewAssertions[1:]...,
	)

	evenMoreForceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypePolaris)).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(newCatalogNamespace).
			HasRestConfig(&sdk.OpenCatalogRestConfigDetails{
				CatalogUri:           privateCatalogUri,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePrivate,
				CatalogName:          newCatalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeVendedCredentials,
			}).
			HasRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      privateOAuthTokenUri,
				OauthClientId:      newOAuthClientId,
				OauthClientSecret:  newOAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope, additionalOAuthAllowedScope},
			}),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationOpenCatalogDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypePolaris).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(newCatalogNamespace),
		resourceshowoutputassert.CatalogIntegrationOpenCatalogDescribeOutput(t, ref).
			HasRestConfigCatalogUri(privateCatalogUri).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePrivate).
			HasRestConfigCatalogName(newCatalogName).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials).
			HasRestAuthenticationOauthTokenUri(privateOAuthTokenUri).
			HasRestAuthenticationOauthClientId(newOAuthClientId).
			HasRestAuthenticationOauthAllowedScopes(oAuthAllowedScope, additionalOAuthAllowedScope),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationOpenCatalog),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config:          config.FromModels(t, basic, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				Check:           assertThat(t, basicAssertions...),
			},
			// Import
			{
				Config:                  config.FromModels(t, basic, oauthClientSecretVariableModel),
				ConfigVariables:         oauthClientSecretConfigVariables,
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rest_config", "rest_authentication"},
			},
			// Change alterable props
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config:          config.FromModels(t, altered, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				Check:           assertThat(t, alteredProperties...),
			},
			// Unset
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config:          config.FromModels(t, basic, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				Check:           assertThat(t, basicAssertionsWithRefreshIntervalZero...),
			},
			// Destroy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroy),
					},
				},
				Config:          config.FromModels(t, basic, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				Destroy:         true,
			},
			// Create with all attributes
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config:          config.FromModels(t, allAttributes, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				Check:           assertThat(t, completeAssertions...),
			},
			// Import
			{
				Config:                  config.FromModels(t, allAttributes, oauthClientSecretVariableModel),
				ConfigVariables:         oauthClientSecretConfigVariables,
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"refresh_interval_seconds", "rest_config", "rest_authentication"},
			},
			// Change alterable props externally
			{
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
				Config:          config.FromModels(t, allAttributes, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				Check:           assertThat(t, completeAssertions...),
			},
			// Change force new "catalog_namespace" prop
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          config.FromModels(t, withChangedCatalogNamespace, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				Check:           assertThat(t, forceNewAssertions...),
			},
			// Change force new "catalog_namespace" prop externally
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithOrReplace(true).
						WithOpenCatalogCatalogSourceParams(*sdk.NewOpenCatalogParamsRequest().
							WithRestConfig(completeRestConfig[0]).
							WithRestAuthentication(completeRestAuth[0]).
							WithCatalogNamespace(externalCatalogNamespace))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "catalog_namespace", sdk.String(newCatalogNamespace), sdk.String(externalCatalogNamespace)),
						planchecks.ExpectChange(ref, "catalog_namespace", tfjson.ActionDelete, sdk.String(externalCatalogNamespace), sdk.String(newCatalogNamespace)),
					},
				},
				Config:          config.FromModels(t, withChangedCatalogNamespace, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				Check:           assertThat(t, forceNewAssertions...),
			},
			// Change force new props in "rest_config"
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          config.FromModels(t, withChangedRestConfig, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				Check:           assertThat(t, moreForceNewAssertions...),
			},
			// Change force new props in "rest_config" externally
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithOrReplace(true).
						WithOpenCatalogCatalogSourceParams(*sdk.NewOpenCatalogParamsRequest().
							WithRestConfig(completeRestConfig[0]).
							WithRestAuthentication(completeRestAuth[0]))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "rest_config.0.catalog_uri", sdk.String(privateCatalogUri), sdk.String(catalogUri)),
						planchecks.ExpectChange(ref, "rest_config.0.catalog_uri", tfjson.ActionDelete, sdk.String(catalogUri), sdk.String(privateCatalogUri)),
						planchecks.ExpectDrift(ref, "rest_config.0.catalog_name", sdk.String(newCatalogName), sdk.String(catalogName)),
						planchecks.ExpectChange(ref, "rest_config.0.catalog_name", tfjson.ActionDelete, sdk.String(catalogName), sdk.String(newCatalogName)),
						planchecks.ExpectDrift(ref, "rest_config.0.catalog_api_type", sdk.String(string(sdk.CatalogIntegrationCatalogApiTypePrivate)), sdk.String(string(sdk.CatalogIntegrationCatalogApiTypePublic))),
						planchecks.ExpectChange(ref, "rest_config.0.catalog_api_type", tfjson.ActionDelete, sdk.String(string(sdk.CatalogIntegrationCatalogApiTypePublic)), sdk.String(string(sdk.CatalogIntegrationCatalogApiTypePrivate))),
						planchecks.ExpectDrift(ref, "rest_config.0.access_delegation_mode", sdk.String(string(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials)), sdk.String(string(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials))),
						planchecks.ExpectChange(ref, "rest_config.0.access_delegation_mode", tfjson.ActionDelete, sdk.String(string(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)), sdk.String(string(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials))),
					},
				},
				Config:          config.FromModels(t, withChangedRestConfig, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				Check:           assertThat(t, moreForceNewAssertions...),
			},
			// Change alterable "oauth_client_secret" prop in "rest_authentication"
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config:          config.FromModels(t, withChangedOAuthClientSecret, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariablesSecondary,
				Check:           assertThat(t, moreForceNewAssertionsWithChangedSecret...),
			},
			// Change force new props in "rest_authentication"
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          config.FromModels(t, withChangedRestAuth, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariablesSecondary,
				Check:           assertThat(t, evenMoreForceNewAssertions...),
			},
			// Change force new props in "rest_authentication" externally
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithOrReplace(true).
						WithOpenCatalogCatalogSourceParams(*sdk.NewOpenCatalogParamsRequest().
							WithRestConfig(changedRestConfig[0]).
							WithRestAuthentication(
								*sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, oAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}).
									WithOauthTokenUri(privateOAuthTokenUri),
							))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					// Open Catalog only accepts one OAuth token URI (<account_url>/v1/oauth/tokens),
					// so we cannot exercise drift detection for rest_authentication.0.oauth_token_uri.
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "rest_authentication.0.oauth_client_id", sdk.String(newOAuthClientId), sdk.String(oAuthClientId)),
						planchecks.ExpectChange(ref, "rest_authentication.0.oauth_client_id", tfjson.ActionDelete, sdk.String(oAuthClientId), sdk.String(newOAuthClientId)),
						planchecks.ExpectDrift(ref, "rest_authentication.0.oauth_allowed_scopes", sdk.String(fmt.Sprintf("[%s %s]", oAuthAllowedScope, additionalOAuthAllowedScope)), sdk.String(fmt.Sprintf("[%s]", oAuthAllowedScope))),
						planchecks.ExpectChange(ref, "rest_authentication.0.oauth_allowed_scopes", tfjson.ActionDelete, sdk.String(fmt.Sprintf("[%s]", oAuthAllowedScope)), sdk.String(fmt.Sprintf("[%s %s]", oAuthAllowedScope, additionalOAuthAllowedScope))),
					},
				},
				Config:          config.FromModels(t, withChangedRestAuth, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariablesSecondary,
				Check:           assertThat(t, evenMoreForceNewAssertions...),
			},
			// Change "catalog_source" externally
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithOrReplace(true).
						WithObjectStorageCatalogSourceParams(*sdk.NewObjectStorageParamsRequest(sdk.CatalogIntegrationTableFormatDelta))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          config.FromModels(t, withChangedRestAuth, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariablesSecondary,
				Check:           assertThat(t, evenMoreForceNewAssertions...),
			},
		},
	})
}

func TestAcc_CatalogIntegrationOpenCatalog_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	catalogUri := "https://testorg-testacc.snowflakecomputing.com/polaris/api/catalog"
	catalogName := "my_catalog_name"
	restConfig := []sdk.OpenCatalogRestConfigRequest{
		*sdk.NewOpenCatalogRestConfigRequest(catalogUri, catalogName),
	}
	restAuth := []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest("my_client_id", "my_client_secret", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}}),
	}

	refreshIntervalNonPositive := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, restAuth, restConfig).
		WithRefreshIntervalSeconds(0)

	emptyCatalogNamespace := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, restAuth, restConfig).
		WithCatalogNamespace("")

	emptyCatalogUri := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, restAuth, []sdk.OpenCatalogRestConfigRequest{
		*sdk.NewOpenCatalogRestConfigRequest("", catalogName),
	})

	emptyCatalogName := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, restAuth, []sdk.OpenCatalogRestConfigRequest{
		*sdk.NewOpenCatalogRestConfigRequest(catalogUri, ""),
	})

	invalidCatalogApiType := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, restAuth, []sdk.OpenCatalogRestConfigRequest{
		*sdk.NewOpenCatalogRestConfigRequest(catalogUri, catalogName).
			WithCatalogApiType("invalid"),
	})

	invalidAccessDelegationMode := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, restAuth, []sdk.OpenCatalogRestConfigRequest{
		*sdk.NewOpenCatalogRestConfigRequest(catalogUri, catalogName).
			WithAccessDelegationMode("invalid"),
	})

	emptyOAuthTokenUri := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest("my_client_id", "my_client_secret", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}}).
			WithOauthTokenUri(""),
	}, restConfig)

	emptyOAuthClientId := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest("", "my_client_secret", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}}),
	}, restConfig)

	emptyOAuthClientSecret := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest("my_client_id", "", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}}),
	}, restConfig)

	emptyOAuthScopes := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest("my_client_id", "my_client_secret", []sdk.StringListItemWrapper{}),
	}, restConfig)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationOpenCatalog),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, refreshIntervalNonPositive),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected refresh_interval_seconds to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, emptyCatalogNamespace),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "catalog_namespace" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptyCatalogUri),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "rest_config\.0\.catalog_uri" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptyCatalogName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "rest_config\.0\.catalog_name" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, invalidCatalogApiType),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid catalog integration catalog api type: INVALID`),
			},
			{
				Config:      config.FromModels(t, invalidAccessDelegationMode),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid catalog integration access delegation mode: INVALID`),
			},
			{
				Config:      config.FromModels(t, emptyOAuthTokenUri),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "rest_authentication\.0\.oauth_token_uri" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptyOAuthClientId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "rest_authentication\.0\.oauth_client_id" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptyOAuthClientSecret),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "rest_authentication\.0\.oauth_client_secret" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptyOAuthScopes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Not enough list items`),
			},
		},
	})
}

func TestAcc_CatalogIntegrationOpenCatalog_BasicUseCase_ImportValidation(t *testing.T) {
	restConfig := []sdk.OpenCatalogRestConfigRequest{
		*sdk.NewOpenCatalogRestConfigRequest("https://testorg-testacc.snowflakecomputing.com/polaris/api/catalog", "my_catalog_name"),
	}
	restAuth := []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest("my_client_id", "my_client_secret", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}}),
	}

	notificationIntegration, notificationIntegrationCleanup := testClient().NotificationIntegration.Create(t)
	t.Cleanup(notificationIntegrationCleanup)

	catalogIntegrationObjectStorageId, catalogIntegrationObjectStorageCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogIntegrationObjectStorageCleanup)

	catalogIntegrationOpenCatalog := model.CatalogIntegrationOpenCatalog("t", notificationIntegration.ID().Name(), false, restAuth, restConfig)
	catalogIntegrationOpenCatalog2 := model.CatalogIntegrationOpenCatalog("t", catalogIntegrationObjectStorageId.Name(), false, restAuth, restConfig)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationOpenCatalog),
		Steps: []resource.TestStep{
			{
				Config:        config.FromModels(t, catalogIntegrationOpenCatalog),
				ResourceName:  catalogIntegrationOpenCatalog.ResourceReference(),
				ImportState:   true,
				ImportStateId: notificationIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile(fmt.Sprintf(`Integration %s is not a CATALOG integration`, notificationIntegration.ID().Name())),
			},
			{
				Config:        config.FromModels(t, catalogIntegrationOpenCatalog2),
				ResourceName:  catalogIntegrationOpenCatalog2.ResourceReference(),
				ImportState:   true,
				ImportStateId: catalogIntegrationObjectStorageId.Name(),
				ExpectError:   regexp.MustCompile(fmt.Sprintf(`invalid catalog source type, expected %s, got %s`, sdk.CatalogIntegrationCatalogSourceTypePolaris, sdk.CatalogIntegrationCatalogSourceTypeObjectStore)),
			},
		},
	})
}

func TestAcc_CatalogIntegrationOpenCatalog_Import(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountLocator := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogAccountLocator)
	oAuthClientId := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientId)
	oAuthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientSecret)
	catalogUri := fmt.Sprintf("https://%s.snowflakecomputing.com/polaris/api/catalog", accountLocator)
	catalogName := "TEST_CATALOG"
	catalogNamespace := "TEST_NAMESPACE"

	oAuthTokenUri := catalogUri + "/v1/oauth/tokens"
	oAuthAllowedScope := "PRINCIPAL_ROLE:PRINCIPAL"

	oauthClientSecretVariableName := "oauth_client_secret"
	oauthClientSecretVariableModel, oauthClientSecretConfigVariables := config.SecretStringVariableModelWithConfigVariables(oauthClientSecretVariableName, oAuthClientSecret)

	comment := random.Comment()
	refreshIntervalSeconds := random.IntRange(30, 86400)

	completeRestAuth := []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, oAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}).
			WithOauthTokenUri(oAuthTokenUri),
	}

	completeRestConfig := []sdk.OpenCatalogRestConfigRequest{
		*sdk.NewOpenCatalogRestConfigRequest(catalogUri, catalogName).
			WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			WithAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
	}

	allAttributes := model.CatalogIntegrationOpenCatalogVar("t", id.Name(), false, completeRestAuth, completeRestConfig, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(catalogNamespace)

	ref := allAttributes.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationOpenCatalog),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithRefreshIntervalSeconds(refreshIntervalSeconds).
						WithComment(comment).
						WithOpenCatalogCatalogSourceParams(*sdk.NewOpenCatalogParamsRequest().
							WithRestConfig(completeRestConfig[0]).
							WithRestAuthentication(completeRestAuth[0]).
							WithCatalogNamespace(catalogNamespace))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				Config:             config.FromModels(t, allAttributes, oauthClientSecretVariableModel),
				ConfigVariables:    oauthClientSecretConfigVariables,
				ResourceName:       ref,
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			{
				Config:          config.FromModels(t, allAttributes, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
