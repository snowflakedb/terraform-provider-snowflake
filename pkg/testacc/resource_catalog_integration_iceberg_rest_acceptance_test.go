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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CatalogIntegrationIcebergRest_BasicUseCaseOAuth(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountLocator := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogAccountLocator)
	oAuthClientId := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientId)
	oAuthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientSecret)
	newOAuthClientId := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogSecondaryOAuthClientId)
	newOAuthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogSecondaryOAuthClientSecret)
	catalogUri := fmt.Sprintf("https://%s.snowflakecomputing.com/polaris/api/catalog", accountLocator)
	catalogName := "TEST_CATALOG"
	newCatalogName := "TEST_CATALOG2"
	prefix := "prefix"
	newPrefix := "prefix2"

	catalogNamespace := "TEST_NAMESPACE"
	newCatalogNamespace := "TEST_NAMESPACE2"
	externalCatalogNamespace := "TEST_NAMESPACE3"

	oAuthTokenUri := catalogUri + "/v1/oauth/tokens"
	oAuthAllowedScope := "PRINCIPAL_ROLE:PRINCIPAL"
	additionalOAuthAllowedScope := "PRINCIPAL_ROLE:SECONDARY"
	bearerToken := testClient().CatalogIntegration.RetrieveOpenCatalogBearerToken(t, catalogUri, oAuthClientId, oAuthClientSecret, oAuthAllowedScope)

	oauthClientSecretVariableName := "oauth_client_secret"
	bearerTokenVariableName := "bearer_token"
	oauthClientSecretVariableModel, oauthClientSecretConfigVariables := config.SecretStringVariableModelWithConfigVariables(oauthClientSecretVariableName, oAuthClientSecret)
	_, oauthClientSecretConfigVariablesSecondary := config.SecretStringVariableModelWithConfigVariables(oauthClientSecretVariableName, newOAuthClientSecret)
	bearerTokenVariableModel, bearerTokenConfigVariables := config.SecretStringVariableModelWithConfigVariables(bearerTokenVariableName, bearerToken)

	comment := random.Comment()
	externalComment := random.Comment()

	refreshIntervalSeconds := random.IntRange(30, 86400)
	externalRefreshIntervalSeconds := random.IntRange(30, 86400)

	basicRestAuth := *sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, oAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}})
	completeRestAuth := *sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, oAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}).
		WithOauthTokenUri(oAuthTokenUri)
	changedRestAuth := *sdk.NewOAuthRestAuthenticationRequest(newOAuthClientId, newOAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}, {Value: additionalOAuthAllowedScope}}).
		WithOauthTokenUri(oAuthTokenUri)

	basicRestConfig := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithCatalogName(catalogName)
	completeRestConfig := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithPrefix(prefix).
		WithCatalogName(catalogName).
		WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
		WithAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)
	changedRestConfig := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithPrefix(newPrefix).
		WithCatalogName(newCatalogName).
		WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
		WithAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials)

	basic := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, basicRestConfig, basicRestAuth, oauthClientSecretVariableName)

	altered := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), true, basicRestConfig, basicRestAuth, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds)

	allAttributes := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, completeRestConfig, completeRestAuth, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(catalogNamespace)

	withChangedCatalogNamespace := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, completeRestConfig, completeRestAuth, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withChangedRestConfig := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, changedRestConfig, basicRestAuth, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withChangedOAuthClientSecret := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, changedRestConfig, *sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, newOAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}), oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withChangedRestAuth := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, changedRestConfig, changedRestAuth, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withBearerToken := model.CatalogIntegrationIcebergRestBearer("t", id.Name(), false, changedRestConfig, bearerTokenVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasCatalogNamespaceEmpty().
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               "",
				CatalogApiType:       "",
				CatalogName:          catalogName,
				AccessDelegationMode: "",
			}).
			HasOauthRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      "",
				OauthClientId:      oAuthClientId,
				OauthClientSecret:  oAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope},
			}).
			HasBearerRestAuthenticationEmpty().
			HasSigv4RestAuthenticationEmpty(),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(""),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasCatalogNamespace(""),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigPrefix("").
			HasRestConfigCatalogName(catalogName).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials).
			HasOAuthRestAuthenticationOauthTokenUri(catalogUri + "/v1/oauth/tokens").
			HasOAuthRestAuthenticationOauthClientId(oAuthClientId).
			HasOAuthRestAuthenticationOauthAllowedScopes(oAuthAllowedScope),
	}

	basicAssertionsWithRefreshIntervalZero := append(
		[]assert.TestCheckFuncProvider{
			resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
				HasName(id.Name()).
				HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
				HasEnabledString(r.BooleanFalse).
				HasCommentEmpty().
				HasRefreshIntervalSeconds(0).
				HasCatalogNamespaceEmpty().
				HasRestConfig(&sdk.IcebergRestRestConfigDetails{
					CatalogUri:           catalogUri,
					Prefix:               "",
					CatalogApiType:       "",
					CatalogName:          catalogName,
					AccessDelegationMode: "",
				}).
				HasOauthRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
					OauthTokenUri:      "",
					OauthClientId:      oAuthClientId,
					OauthClientSecret:  oAuthClientSecret,
					OauthAllowedScopes: []string{oAuthAllowedScope},
				}).
				HasSigv4RestAuthenticationEmpty().
				HasBearerRestAuthenticationEmpty(),
		},
		basicAssertions[1:]...,
	)

	alteredProperties := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanTrue).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespaceEmpty().
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               "",
				CatalogApiType:       "",
				CatalogName:          catalogName,
				AccessDelegationMode: "",
			}).
			HasOauthRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      "",
				OauthClientId:      oAuthClientId,
				OauthClientSecret:  oAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope},
			}).
			HasBearerRestAuthenticationEmpty().
			HasSigv4RestAuthenticationEmpty(),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(true).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(true).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(""),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigPrefix("").
			HasRestConfigCatalogName(catalogName).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials).
			HasOAuthRestAuthenticationOauthTokenUri(catalogUri + "/v1/oauth/tokens").
			HasOAuthRestAuthenticationOauthClientId(oAuthClientId).
			HasOAuthRestAuthenticationOauthAllowedScopes(oAuthAllowedScope),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(catalogNamespace).
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               prefix,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePublic,
				CatalogName:          catalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials,
			}).
			HasOauthRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      oAuthTokenUri,
				OauthClientId:      oAuthClientId,
				OauthClientSecret:  oAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope},
			}).
			HasBearerRestAuthenticationEmpty().
			HasSigv4RestAuthenticationEmpty(),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(catalogNamespace),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigPrefix(prefix).
			HasRestConfigCatalogName(catalogName).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials).
			HasOAuthRestAuthenticationOauthTokenUri(oAuthTokenUri).
			HasOAuthRestAuthenticationOauthClientId(oAuthClientId).
			HasOAuthRestAuthenticationOauthAllowedScopes(oAuthAllowedScope),
	}

	forceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(newCatalogNamespace).
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               prefix,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePublic,
				CatalogName:          catalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials,
			}).
			HasOauthRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      oAuthTokenUri,
				OauthClientId:      oAuthClientId,
				OauthClientSecret:  oAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope},
			}).
			HasBearerRestAuthenticationEmpty().
			HasSigv4RestAuthenticationEmpty(),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(newCatalogNamespace),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigPrefix(prefix).
			HasRestConfigCatalogName(catalogName).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials).
			HasOAuthRestAuthenticationOauthTokenUri(oAuthTokenUri).
			HasOAuthRestAuthenticationOauthClientId(oAuthClientId).
			HasOAuthRestAuthenticationOauthAllowedScopes(oAuthAllowedScope),
	}

	moreForceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(newCatalogNamespace).
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               newPrefix,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePublic,
				CatalogName:          newCatalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeVendedCredentials,
			}).
			HasOauthRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      "",
				OauthClientId:      oAuthClientId,
				OauthClientSecret:  oAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope},
			}).
			HasBearerRestAuthenticationEmpty().
			HasSigv4RestAuthenticationEmpty(),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(newCatalogNamespace),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigPrefix(newPrefix).
			HasRestConfigCatalogName(newCatalogName).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials).
			HasOAuthRestAuthenticationOauthTokenUri(oAuthTokenUri).
			HasOAuthRestAuthenticationOauthClientId(oAuthClientId).
			HasOAuthRestAuthenticationOauthAllowedScopes(oAuthAllowedScope),
	}

	moreForceNewAssertionsWithChangedSecret := append(
		[]assert.TestCheckFuncProvider{
			resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
				HasName(id.Name()).
				HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
				HasEnabledString(r.BooleanFalse).
				HasComment(comment).
				HasRefreshIntervalSeconds(refreshIntervalSeconds).
				HasCatalogNamespace(newCatalogNamespace).
				HasRestConfig(&sdk.IcebergRestRestConfigDetails{
					CatalogUri:           catalogUri,
					Prefix:               newPrefix,
					CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePublic,
					CatalogName:          newCatalogName,
					AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeVendedCredentials,
				}).
				HasOauthRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
					OauthTokenUri:      "",
					OauthClientId:      oAuthClientId,
					OauthClientSecret:  newOAuthClientSecret,
					OauthAllowedScopes: []string{oAuthAllowedScope},
				}).
				HasSigv4RestAuthenticationEmpty().
				HasBearerRestAuthenticationEmpty(),
		},
		moreForceNewAssertions[1:]...,
	)

	evenMoreForceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(newCatalogNamespace).
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               newPrefix,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePublic,
				CatalogName:          newCatalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeVendedCredentials,
			}).
			HasOauthRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      oAuthTokenUri,
				OauthClientId:      newOAuthClientId,
				OauthClientSecret:  newOAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope, additionalOAuthAllowedScope},
			}).
			HasBearerRestAuthenticationEmpty().
			HasSigv4RestAuthenticationEmpty(),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(newCatalogNamespace),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigPrefix(newPrefix).
			HasRestConfigCatalogName(newCatalogName).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasOAuthRestAuthenticationOauthTokenUri(oAuthTokenUri).
			HasOAuthRestAuthenticationOauthClientId(newOAuthClientId).
			HasOAuthRestAuthenticationOauthAllowedScopes(oAuthAllowedScope, additionalOAuthAllowedScope),
	}

	withBearerTokenAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(newCatalogNamespace).
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               newPrefix,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePublic,
				CatalogName:          newCatalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeVendedCredentials,
			}).
			HasOauthRestAuthenticationEmpty().
			HasBearerRestAuthentication(&sdk.BearerRestAuthenticationDetails{BearerToken: bearerToken}).
			HasSigv4RestAuthenticationEmpty(),
	}, evenMoreForceNewAssertions[1:4]...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationIcebergRest),
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
						WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
							WithRestConfig(completeRestConfig).
							WithOAuthRestAuthentication(completeRestAuth).
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
						WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
							WithRestConfig(completeRestConfig).
							WithOAuthRestAuthentication(basicRestAuth))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "rest_config.0.prefix", sdk.String(newPrefix), sdk.String(prefix)),
						planchecks.ExpectChange(ref, "rest_config.0.prefix", tfjson.ActionDelete, sdk.String(prefix), sdk.String(newPrefix)),
						planchecks.ExpectDrift(ref, "rest_config.0.catalog_name", sdk.String(newCatalogName), sdk.String(catalogName)),
						planchecks.ExpectChange(ref, "rest_config.0.catalog_name", tfjson.ActionDelete, sdk.String(catalogName), sdk.String(newCatalogName)),
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
			// Change force new props in "oauth_rest_authentication"
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
			// Change force new props in "oauth_rest_authentication" externally
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithOrReplace(true).
						WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
							WithRestConfig(changedRestConfig).
							WithOAuthRestAuthentication(completeRestAuth))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					// Open Catalog only accepts one OAuth token URI (<account_url>/v1/oauth/tokens),
					// so we cannot exercise drift detection for oauth_rest_authentication.0.oauth_token_uri.
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "oauth_rest_authentication.0.oauth_client_id", sdk.String(newOAuthClientId), sdk.String(oAuthClientId)),
						planchecks.ExpectChange(ref, "oauth_rest_authentication.0.oauth_client_id", tfjson.ActionDelete, sdk.String(oAuthClientId), sdk.String(newOAuthClientId)),
						planchecks.ExpectDrift(ref, "oauth_rest_authentication.0.oauth_allowed_scopes", sdk.String(fmt.Sprintf("[%s %s]", oAuthAllowedScope, additionalOAuthAllowedScope)), sdk.String(fmt.Sprintf("[%s]", oAuthAllowedScope))),
						planchecks.ExpectChange(ref, "oauth_rest_authentication.0.oauth_allowed_scopes", tfjson.ActionDelete, sdk.String(fmt.Sprintf("[%s]", oAuthAllowedScope)), sdk.String(fmt.Sprintf("[%s %s]", oAuthAllowedScope, additionalOAuthAllowedScope))),
					},
				},
				Config:          config.FromModels(t, withChangedRestAuth, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariablesSecondary,
				Check:           assertThat(t, evenMoreForceNewAssertions...),
			},
			// Change to different authentication type
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          config.FromModels(t, withBearerToken, bearerTokenVariableModel),
				ConfigVariables: bearerTokenConfigVariables,
				Check:           assertThat(t, withBearerTokenAssertions...),
			},
			// Change catalog source externally
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
				Config:          config.FromModels(t, withBearerToken, bearerTokenVariableModel),
				ConfigVariables: bearerTokenConfigVariables,
				Check:           assertThat(t, withBearerTokenAssertions...),
			},
		},
	})
}

func TestAcc_CatalogIntegrationIcebergRest_ImportOAuth(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountLocator := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogAccountLocator)
	oAuthClientId := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientId)
	oAuthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientSecret)
	catalogUri := fmt.Sprintf("https://%s.snowflakecomputing.com/polaris/api/catalog", accountLocator)
	catalogName := "TEST_CATALOG"
	prefix := "prefix"
	catalogNamespace := "TEST_NAMESPACE"

	oAuthTokenUri := catalogUri + "/v1/oauth/tokens"
	oAuthAllowedScope := "PRINCIPAL_ROLE:PRINCIPAL"

	oauthClientSecretVariableName := "oauth_client_secret"
	oauthClientSecretVariableModel, oauthClientSecretConfigVariables := config.SecretStringVariableModelWithConfigVariables(oauthClientSecretVariableName, oAuthClientSecret)

	comment := random.Comment()
	refreshIntervalSeconds := random.IntRange(30, 86400)

	completeRestAuth := *sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, oAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}).
		WithOauthTokenUri(oAuthTokenUri)

	completeRestConfig := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithPrefix(prefix).
		WithCatalogName(catalogName).
		WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
		WithAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)

	allAttributes := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, completeRestConfig, completeRestAuth, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(catalogNamespace)

	ref := allAttributes.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationIcebergRest),
		Steps: []resource.TestStep{
			// Import the externally created resource
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithRefreshIntervalSeconds(refreshIntervalSeconds).
						WithComment(comment).
						WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
							WithRestConfig(completeRestConfig).
							WithOAuthRestAuthentication(completeRestAuth).
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
			// After import, only oauth_client_secret is unset in state (write-only field Snowflake never returns).
			// All ForceNew fields (rest_config, oauth_client_id, oauth_allowed_scopes, etc.) must already match
			// the config — only an in-place update for oauth_client_secret should be planned.
			{
				Config:          config.FromModels(t, allAttributes, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(ref, "rest_config", "refresh_interval_seconds"),
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

func TestAcc_CatalogIntegrationIcebergRest_ImportBearer(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountLocator := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogAccountLocator)
	oAuthClientId := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientId)
	oAuthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientSecret)
	catalogUri := fmt.Sprintf("https://%s.snowflakecomputing.com/polaris/api/catalog", accountLocator)
	catalogName := "TEST_CATALOG"
	prefix := "prefix"
	catalogNamespace := "TEST_NAMESPACE"

	bearerToken := testClient().CatalogIntegration.RetrieveOpenCatalogBearerToken(t, catalogUri, oAuthClientId, oAuthClientSecret, "PRINCIPAL_ROLE:PRINCIPAL")

	bearerTokenVariableName := "bearer_token"
	bearerTokenVariableModel, bearerTokenConfigVariables := config.SecretStringVariableModelWithConfigVariables(bearerTokenVariableName, bearerToken)

	comment := random.Comment()
	refreshIntervalSeconds := random.IntRange(30, 86400)

	completeRestAuth := *sdk.NewBearerRestAuthenticationRequest(bearerToken)

	completeRestConfig := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithPrefix(prefix).
		WithCatalogName(catalogName).
		WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
		WithAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)

	allAttributes := model.CatalogIntegrationIcebergRestBearer("t", id.Name(), false, completeRestConfig, bearerTokenVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(catalogNamespace)

	ref := allAttributes.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationIcebergRest),
		Steps: []resource.TestStep{
			// Import the externally created resource
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithRefreshIntervalSeconds(refreshIntervalSeconds).
						WithComment(comment).
						WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
							WithRestConfig(completeRestConfig).
							WithBearerRestAuthentication(completeRestAuth).
							WithCatalogNamespace(catalogNamespace))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				Config:             config.FromModels(t, allAttributes, bearerTokenVariableModel),
				ConfigVariables:    bearerTokenConfigVariables,
				ResourceName:       ref,
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			// After import, only bearer_token is unset in state (write-only field Snowflake never returns).
			// All ForceNew fields (rest_config, etc.) must already match the config — only an in-place
			// update for bearer_token should be planned.
			{
				Config:          config.FromModels(t, allAttributes, bearerTokenVariableModel),
				ConfigVariables: bearerTokenConfigVariables,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(ref, "rest_config", "refresh_interval_seconds"),
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

func TestAcc_CatalogIntegrationIcebergRest_ImportSigV4(t *testing.T) {
	t.Skip("SNOW-3888393: no real Iceberg REST catalog with SigV4 auth available; fake credentials fail when CREATE validates the catalog")
	id := testClient().Ids.RandomAccountObjectIdentifier()

	catalogUri := "https://api.tabular.io/ws"
	catalogName := random.AlphanumericN(15)
	prefix := "prefix"
	catalogNamespace := random.AlphanumericN(15)

	sigV4IamRole := "arn:aws:iam::123456789012:role/sigv4-role"
	sigV4SigningRegion := "us-west-2"
	sigV4ExternalId := random.AlphanumericN(15)

	comment := random.Comment()
	refreshIntervalSeconds := random.IntRange(30, 86400)

	completeRestAuth := *sdk.NewSigV4RestAuthenticationRequest(sigV4IamRole).
		WithSigv4SigningRegion(sigV4SigningRegion).
		WithSigv4ExternalId(sigV4ExternalId)

	completeRestConfig := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithPrefix(prefix).
		WithCatalogName(catalogName).
		WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway).
		WithAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)

	allAttributes := model.CatalogIntegrationIcebergRestSigV4("t", id.Name(), false, completeRestConfig, completeRestAuth).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(catalogNamespace)

	ref := allAttributes.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationIcebergRest),
		Steps: []resource.TestStep{
			// Import the externally created resource
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithRefreshIntervalSeconds(refreshIntervalSeconds).
						WithComment(comment).
						WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
							WithRestConfig(completeRestConfig).
							WithSigV4RestAuthentication(completeRestAuth).
							WithCatalogNamespace(catalogNamespace))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				Config:             config.FromModels(t, allAttributes),
				ResourceName:       ref,
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			// After import, all ForceNew fields (rest_config, sigv4_rest_authentication, etc.) must already
			// match the config. sigv4_external_id is copied from configuration into state during import because
			// Snowflake does not return it — no recreation plan should be produced.
			{
				Config: config.FromModels(t, allAttributes),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.PrintPlanDetails(ref, "rest_config", "refresh_interval_seconds"),
						// SIGV4_EXTERNAL_ID is ForceNew and is not returned from Snowflake, so we need to destroy and create the resource.
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

func TestAcc_CatalogIntegrationIcebergRest_BasicUseCaseBearer(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountLocator := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogAccountLocator)
	oAuthClientId := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientId)
	oAuthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientSecret)
	newOAuthClientId := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogSecondaryOAuthClientId)
	newOAuthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogSecondaryOAuthClientSecret)
	catalogUri := fmt.Sprintf("https://%s.snowflakecomputing.com/polaris/api/catalog", accountLocator)
	catalogName := "TEST_CATALOG"
	oAuthAllowedScope := "PRINCIPAL_ROLE:PRINCIPAL"
	additionalOAuthAllowedScope := "PRINCIPAL_ROLE:SECONDARY"

	bearerToken1 := testClient().CatalogIntegration.RetrieveOpenCatalogBearerToken(t, catalogUri, oAuthClientId, oAuthClientSecret, oAuthAllowedScope)
	bearerToken2 := testClient().CatalogIntegration.RetrieveOpenCatalogBearerToken(t, catalogUri, newOAuthClientId, newOAuthClientSecret, additionalOAuthAllowedScope)

	oauthClientSecretVariableName := "oauth_client_secret"
	bearerTokenVariableName := "bearer_token"
	bearerTokenVariableModel, bearerTokenConfigVariables := config.SecretStringVariableModelWithConfigVariables(bearerTokenVariableName, bearerToken1)
	_, bearerTokenConfigVariablesSecondary := config.SecretStringVariableModelWithConfigVariables(bearerTokenVariableName, bearerToken2)
	oauthClientSecretVariableModel, oauthClientSecretConfigVariables := config.SecretStringVariableModelWithConfigVariables(oauthClientSecretVariableName, oAuthClientSecret)

	basicRestCfg := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithCatalogName(catalogName)

	basic := model.CatalogIntegrationIcebergRestBearer("t", id.Name(), false, basicRestCfg, bearerTokenVariableName)
	updated := model.CatalogIntegrationIcebergRestBearer("t", id.Name(), false, basicRestCfg, bearerTokenVariableName)
	withOAuth := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, basicRestCfg, *sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, oAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}), oauthClientSecretVariableName)

	ref := basic.ResourceReference()

	commonAssertions := []assert.TestCheckFuncProvider{
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(""),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasCatalogNamespace("").
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigPrefix("").
			HasRestConfigCatalogName(catalogName).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
	}

	basicAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasCatalogNamespaceEmpty().
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               "",
				CatalogApiType:       "",
				CatalogName:          catalogName,
				AccessDelegationMode: "",
			}).
			HasOauthRestAuthenticationEmpty().
			HasBearerRestAuthentication(&sdk.BearerRestAuthenticationDetails{BearerToken: bearerToken1}).
			HasSigv4RestAuthenticationEmpty(),
	}, commonAssertions...)

	updatedAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasCatalogNamespaceEmpty().
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               "",
				CatalogApiType:       "",
				CatalogName:          catalogName,
				AccessDelegationMode: "",
			}).
			HasOauthRestAuthenticationEmpty().
			HasBearerRestAuthentication(&sdk.BearerRestAuthenticationDetails{BearerToken: bearerToken2}).
			HasSigv4RestAuthenticationEmpty(),
	}, commonAssertions...)

	oauthAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasCatalogNamespaceEmpty().
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               "",
				CatalogApiType:       "",
				CatalogName:          catalogName,
				AccessDelegationMode: "",
			}).
			HasOauthRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      "",
				OauthClientId:      oAuthClientId,
				OauthClientSecret:  oAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope},
			}).
			HasBearerRestAuthenticationEmpty().
			HasSigv4RestAuthenticationEmpty(),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasOAuthRestAuthenticationOauthTokenUri(catalogUri + "/v1/oauth/tokens").
			HasOAuthRestAuthenticationOauthClientId(oAuthClientId).
			HasOAuthRestAuthenticationOauthAllowedScopes(oAuthAllowedScope),
	}, commonAssertions...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationIcebergRest),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config:          config.FromModels(t, basic, bearerTokenVariableModel),
				ConfigVariables: bearerTokenConfigVariables,
				Check:           assertThat(t, basicAssertions...),
			},
			// Change alterable "bearer_token" prop in "bearer_rest_authentication"
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config:          config.FromModels(t, updated, bearerTokenVariableModel),
				ConfigVariables: bearerTokenConfigVariablesSecondary,
				Check:           assertThat(t, updatedAssertions...),
			},
			// Change to different authentication type
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          config.FromModels(t, withOAuth, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				Check:           assertThat(t, oauthAssertions...),
			},
		},
	})
}

func TestAcc_CatalogIntegrationIcebergRest_BasicUseCaseSigV4(t *testing.T) {
	t.Skip("SNOW-3888393: no real Iceberg REST catalog with SigV4 auth available; fake credentials fail when CREATE validates the catalog")
	id := testClient().Ids.RandomAccountObjectIdentifier()

	catalogUri := "https://api.tabular.io/ws"

	sigV4IamRole := "arn:aws:iam::123456789012:role/sigv4-role-1"
	sigV4SigningRegion := "us-west-2"
	newSigV4IamRole := "arn:aws:iam::123456789012:role/sigv4-role-2"
	newSigV4SigningRegion := "eu-west-1"
	newSigV4ExternalId := "external-id-2"

	bearerToken := random.AlphanumericN(32)

	bearerTokenVariableName := "bearer_token"
	bearerTokenVariableModel, bearerTokenConfigVariables := config.SecretStringVariableModelWithConfigVariables(bearerTokenVariableName, bearerToken)

	basicRestCfg := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway)

	basicSigV4Auth := *sdk.NewSigV4RestAuthenticationRequest(sigV4IamRole).
		WithSigv4SigningRegion(sigV4SigningRegion)

	updatedSigV4Auth := *sdk.NewSigV4RestAuthenticationRequest(newSigV4IamRole).
		WithSigv4SigningRegion(newSigV4SigningRegion).
		WithSigv4ExternalId(newSigV4ExternalId)

	basic := model.CatalogIntegrationIcebergRestSigV4("t", id.Name(), false, basicRestCfg, basicSigV4Auth)

	updated := model.CatalogIntegrationIcebergRestSigV4("t", id.Name(), false, basicRestCfg, updatedSigV4Auth)

	withBearerToken := model.CatalogIntegrationIcebergRestBearer("t", id.Name(), false, basicRestCfg, bearerTokenVariableName)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasCatalogNamespaceEmpty().
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               "",
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway,
				CatalogName:          "",
				AccessDelegationMode: "",
			}).
			HasOauthRestAuthenticationEmpty().
			HasBearerRestAuthenticationEmpty().
			HasSigV4RestAuthentication(&sdk.SigV4RestAuthenticationDetails{
				Sigv4IamRole:       sigV4IamRole,
				Sigv4SigningRegion: sigV4SigningRegion,
				Sigv4ExternalId:    "",
			}),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasSigv4RestAuthenticationSigv4IamRole(sigV4IamRole).
			HasSigv4RestAuthenticationSigv4SigningRegion(sigV4SigningRegion),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasCatalogNamespace(""),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigPrefix("").
			HasRestConfigCatalogName("").
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
	}

	updatedAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasCatalogNamespaceEmpty().
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               "",
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway,
				CatalogName:          "",
				AccessDelegationMode: "",
			}).
			HasOauthRestAuthenticationEmpty().
			HasBearerRestAuthenticationEmpty().
			HasSigV4RestAuthentication(&sdk.SigV4RestAuthenticationDetails{
				Sigv4IamRole:       newSigV4IamRole,
				Sigv4SigningRegion: newSigV4SigningRegion,
				Sigv4ExternalId:    newSigV4ExternalId,
			}),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasSigv4RestAuthenticationSigv4IamRole(newSigV4IamRole).
			HasSigv4RestAuthenticationSigv4SigningRegion(newSigV4SigningRegion),
	}, basicAssertions[2:]...)

	withBearerTokenAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasCatalogNamespaceEmpty().
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               "",
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway,
				CatalogName:          "",
				AccessDelegationMode: "",
			}).
			HasOauthRestAuthenticationEmpty().
			HasBearerRestAuthentication(&sdk.BearerRestAuthenticationDetails{BearerToken: bearerToken}).
			HasSigv4RestAuthenticationEmpty(),
	}, basicAssertions[2:]...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationIcebergRest),
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
			// Change force new props in "sigv4_rest_authentication"
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, updated),
				Check:  assertThat(t, updatedAssertions...),
			},
			// Change force new props in "sigv4_rest_authentication" externally
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithOrReplace(true).
						WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
							WithRestConfig(basicRestCfg).
							WithSigV4RestAuthentication(basicSigV4Auth))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "sigv4_rest_authentication.0.sigv4_iam_role", sdk.String(newSigV4IamRole), sdk.String(sigV4IamRole)),
						planchecks.ExpectChange(ref, "sigv4_rest_authentication.0.sigv4_iam_role", tfjson.ActionDelete, sdk.String(sigV4IamRole), sdk.String(newSigV4IamRole)),
						planchecks.ExpectDrift(ref, "sigv4_rest_authentication.0.sigv4_signing_region", sdk.String(newSigV4SigningRegion), sdk.String(sigV4SigningRegion)),
						planchecks.ExpectChange(ref, "sigv4_rest_authentication.0.sigv4_signing_region", tfjson.ActionDelete, sdk.String(sigV4SigningRegion), sdk.String(newSigV4SigningRegion)),
					},
				},
				Config: config.FromModels(t, updated),
				Check:  assertThat(t, updatedAssertions...),
			},
			// Change to different authentication type
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:          config.FromModels(t, withBearerToken, bearerTokenVariableModel),
				ConfigVariables: bearerTokenConfigVariables,
				Check:           assertThat(t, withBearerTokenAssertions...),
			},
		},
	})
}

func TestAcc_CatalogIntegrationIcebergRest_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	accountLocator := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogAccountLocator)
	oAuthClientId := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientId)
	oAuthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientSecret)
	catalogUri := fmt.Sprintf("https://%s.snowflakecomputing.com/polaris/api/catalog", accountLocator)
	catalogName := "TEST_CATALOG"
	prefix := "prefix"
	catalogNamespace := "TEST_NAMESPACE"

	oAuthTokenUri := catalogUri + "/v1/oauth/tokens"
	oAuthAllowedScope := "PRINCIPAL_ROLE:PRINCIPAL"

	oauthClientSecretVariableName := "oauth_client_secret"
	oauthClientSecretVariableModel, oauthClientSecretConfigVariables := config.SecretStringVariableModelWithConfigVariables(oauthClientSecretVariableName, oAuthClientSecret)

	comment := random.Comment()
	refreshIntervalSeconds := random.IntRange(30, 86400)

	completeRestAuth := *sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, oAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}).
		WithOauthTokenUri(oAuthTokenUri)

	completeRestConfig := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithPrefix(prefix).
		WithCatalogName(catalogName).
		WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
		WithAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)

	complete := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, completeRestConfig, completeRestAuth, oauthClientSecretVariableName).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(catalogNamespace)

	ref := complete.ResourceReference()

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasCatalogSource(string(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest)).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(catalogNamespace).
			HasFullyQualifiedName(id.FullyQualifiedName()).
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               prefix,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePublic,
				CatalogName:          catalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials,
			}).
			HasOauthRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      oAuthTokenUri,
				OauthClientId:      oAuthClientId,
				OauthClientSecret:  oAuthClientSecret,
				OauthAllowedScopes: []string{oAuthAllowedScope},
			}).
			HasBearerRestAuthenticationEmpty().
			HasSigv4RestAuthenticationEmpty(),
		resourceshowoutputassert.CatalogIntegrationShowOutput(t, ref).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(catalogNamespace),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasRestConfigCatalogUri(catalogUri).
			HasRestConfigPrefix(prefix).
			HasRestConfigCatalogName(catalogName).
			HasRestConfigCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasRestConfigAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials).
			HasOAuthRestAuthenticationOauthTokenUri(oAuthTokenUri).
			HasOAuthRestAuthenticationOauthClientId(oAuthClientId).
			HasOAuthRestAuthenticationOauthAllowedScopes(oAuthAllowedScope),
		objectassert.CatalogIntegration(t, id).
			HasName(id.Name()).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasEnabled(false).
			HasComment(comment).
			HasCreatedOnNotEmpty(),
		objectassert.CatalogIntegrationIcebergRestDetails(t, id).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(catalogNamespace).
			HasRestConfigWith(objectassert.NewIcebergRestRestConfigDetailsAssert().
				HasCatalogUri(catalogUri).
				HasPrefix(prefix).
				HasCatalogName(catalogName).
				HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
				HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)).
			HasOAuthRestAuthenticationWith(objectassert.NewOAuthRestAuthenticationDetailsAssert().
				HasOauthTokenUri(oAuthTokenUri).
				HasOauthClientId(oAuthClientId).
				HasOauthAllowedScopes(oAuthAllowedScope)),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationIcebergRest),
		Steps: []resource.TestStep{
			// Create - with all optionals (including optional force-new fields)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config:          config.FromModels(t, complete, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				Check:           assertThat(t, completeAssertions...),
			},
			// Import - with all optionals.
			// ImportOAuth has no ImportStateVerifyIgnore list (it uses ImportStatePersist instead).
			// These fields cannot round-trip through import:
			{
				Config:            config.FromModels(t, complete, oauthClientSecretVariableModel),
				ConfigVariables:   oauthClientSecretConfigVariables,
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"refresh_interval_seconds",                        // unreadable: Read skips this attribute; value is only compared via describe_output
					"oauth_rest_authentication.0.oauth_client_secret", // write-only: Snowflake never returns oauth_client_secret (see ImportOAuth)
				},
			},
		},
	})
}

func TestAcc_CatalogIntegrationIcebergRest_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	catalogUri := "https://api.tabular.io/ws"
	catalogName := "my_catalog_name"
	restConfig := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithCatalogName(catalogName)
	restAuth := *sdk.NewOAuthRestAuthenticationRequest("my_client_id", "my_client_secret", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}})

	oauthClientSecretVariableName := "oauth_client_secret"
	bearerTokenVariableName := "bearer_token"
	oauthClientSecretVariableModel, oauthClientSecretConfigVariables := config.SecretStringVariableModelWithConfigVariables(oauthClientSecretVariableName, "my_client_secret")
	_, oauthClientSecretConfigVariablesEmpty := config.SecretStringVariableModelWithConfigVariables(oauthClientSecretVariableName, "")
	bearerTokenVariableModel, bearerTokenConfigVariables := config.SecretStringVariableModelWithConfigVariables(bearerTokenVariableName, "token")
	_, bearerTokenConfigVariablesEmpty := config.SecretStringVariableModelWithConfigVariables(bearerTokenVariableName, "")

	refreshIntervalNonPositive := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, restAuth, oauthClientSecretVariableName).
		WithRefreshIntervalSeconds(0)

	emptyCatalogNamespace := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, restAuth, oauthClientSecretVariableName).
		WithCatalogNamespace("")

	emptyCatalogUri := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, *sdk.NewIcebergRestRestConfigRequest(""), restAuth, oauthClientSecretVariableName)

	emptyCatalogName := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithCatalogName(""), restAuth, oauthClientSecretVariableName)

	invalidCatalogApiType := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithCatalogName(catalogName).
		WithCatalogApiType("invalid"), restAuth, oauthClientSecretVariableName)

	invalidAccessDelegationMode := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithCatalogName(catalogName).
		WithAccessDelegationMode("invalid"), restAuth, oauthClientSecretVariableName)

	emptyOAuthTokenUri := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, *sdk.NewOAuthRestAuthenticationRequest("my_client_id", "my_client_secret", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}}).
		WithOauthTokenUri(""), oauthClientSecretVariableName)

	emptyOAuthClientId := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, *sdk.NewOAuthRestAuthenticationRequest("", "my_client_secret", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}}), oauthClientSecretVariableName)

	emptyOAuthClientSecret := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, *sdk.NewOAuthRestAuthenticationRequest("my_client_id", "", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}}), oauthClientSecretVariableName)

	emptyOAuthScopes := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, *sdk.NewOAuthRestAuthenticationRequest("my_client_id", "my_client_secret", []sdk.StringListItemWrapper{}), oauthClientSecretVariableName)

	emptyBearerToken := model.CatalogIntegrationIcebergRestBearer("t", id.Name(), false, restConfig, bearerTokenVariableName)

	emptySigV4IamRole := model.CatalogIntegrationIcebergRestSigV4("t", id.Name(), false, restConfig, *sdk.NewSigV4RestAuthenticationRequest(""))

	emptySigV4SigningRegion := model.CatalogIntegrationIcebergRestSigV4(
		"t", id.Name(), false, restConfig, *sdk.NewSigV4RestAuthenticationRequest("arn:aws:iam::123456789012:role/role").
			WithSigv4SigningRegion(""),
	)

	emptySigV4ExternalId := model.CatalogIntegrationIcebergRestSigV4(
		"t", id.Name(), false, restConfig, *sdk.NewSigV4RestAuthenticationRequest("arn:aws:iam::123456789012:role/role").
			WithSigv4ExternalId(""),
	)

	noAuthentication := model.CatalogIntegrationIcebergRest("t", id.Name(), false, []sdk.IcebergRestRestConfigRequest{restConfig})

	oauthAndBearer := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, restAuth, oauthClientSecretVariableName).
		WithBearerRestAuthentication(bearerTokenVariableName)

	oauthAndSigv4 := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, restAuth, oauthClientSecretVariableName).
		WithSigV4RestAuthentication(*sdk.NewSigV4RestAuthenticationRequest("arn:aws:iam::123456789012:role/role"))

	bearerAndSigv4 := model.CatalogIntegrationIcebergRest("t", id.Name(), false, []sdk.IcebergRestRestConfigRequest{restConfig}).
		WithBearerRestAuthentication(bearerTokenVariableName).
		WithSigV4RestAuthentication(*sdk.NewSigV4RestAuthenticationRequest("arn:aws:iam::123456789012:role/role"))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationIcebergRest),
		Steps: []resource.TestStep{
			{
				Config:          config.FromModels(t, refreshIntervalNonPositive, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`expected refresh_interval_seconds to be at least \(1\), got 0`),
			},
			{
				Config:          config.FromModels(t, emptyCatalogNamespace, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`expected "catalog_namespace" to not be an empty string`),
			},
			{
				Config:          config.FromModels(t, emptyCatalogUri, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`expected "rest_config\.0\.catalog_uri" to not be an empty string`),
			},
			{
				Config:          config.FromModels(t, emptyCatalogName, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`expected "rest_config\.0\.catalog_name" to not be an empty string`),
			},
			{
				Config:          config.FromModels(t, invalidCatalogApiType, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`invalid catalog integration catalog api type: INVALID`),
			},
			{
				Config:          config.FromModels(t, invalidAccessDelegationMode, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`invalid catalog integration access delegation mode: INVALID`),
			},
			{
				Config:          config.FromModels(t, emptyOAuthTokenUri, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`expected "oauth_rest_authentication.0\.oauth_token_uri" to not be an empty string`),
			},
			{
				Config:          config.FromModels(t, emptyOAuthClientId, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`expected "oauth_rest_authentication.0\.oauth_client_id" to not be an empty string`),
			},
			{
				Config:          config.FromModels(t, emptyOAuthClientSecret, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariablesEmpty,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`expected "oauth_rest_authentication.0\.oauth_client_secret" to not be an empty string`),
			},
			{
				Config:          config.FromModels(t, emptyOAuthScopes, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`Not enough list items`),
			},
			{
				Config:          config.FromModels(t, emptyBearerToken, bearerTokenVariableModel),
				ConfigVariables: bearerTokenConfigVariablesEmpty,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`expected "bearer_rest_authentication.0\.bearer_token" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptySigV4IamRole),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "sigv4_rest_authentication.0\.sigv4_iam_role" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptySigV4SigningRegion),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "sigv4_rest_authentication.0\.sigv4_signing_region" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptySigV4ExternalId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "sigv4_rest_authentication.0\.sigv4_external_id" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, noAuthentication),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
			{
				Config: config.FromModels(t, oauthAndBearer, oauthClientSecretVariableModel, bearerTokenVariableModel),
				ConfigVariables: tfconfig.Variables{
					oauthClientSecretVariableName: tfconfig.StringVariable("my_client_secret"),
					bearerTokenVariableName:       tfconfig.StringVariable("token"),
				},
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
			{
				Config:          config.FromModels(t, oauthAndSigv4, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`Invalid combination of arguments`),
			},
			{
				Config:          config.FromModels(t, bearerAndSigv4, bearerTokenVariableModel),
				ConfigVariables: bearerTokenConfigVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`Invalid combination of arguments`),
			},
		},
	})
}

func TestAcc_CatalogIntegrationIcebergRest_BasicUseCase_ImportValidation(t *testing.T) {
	restConfig := *sdk.NewIcebergRestRestConfigRequest("https://api.tabular.io/ws").
		WithCatalogName("my_catalog_name")
	restAuth := *sdk.NewOAuthRestAuthenticationRequest("my_client_id", "my_client_secret", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}})
	oauthClientSecretVariableName := "oauth_client_secret"
	oauthClientSecretVariableModel, oauthClientSecretConfigVariables := config.SecretStringVariableModelWithConfigVariables(oauthClientSecretVariableName, "my_client_secret")

	notificationIntegration, notificationIntegrationCleanup := testClient().NotificationIntegration.Create(t)
	t.Cleanup(notificationIntegrationCleanup)

	catalogIntegrationObjectStorageId, catalogIntegrationObjectStorageCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogIntegrationObjectStorageCleanup)

	catalogIntegrationIcebergRest := model.CatalogIntegrationIcebergRestOAuth("t", notificationIntegration.ID().Name(), false, restConfig, restAuth, oauthClientSecretVariableName)
	catalogIntegrationIcebergRest2 := model.CatalogIntegrationIcebergRestOAuth("t", catalogIntegrationObjectStorageId.Name(), false, restConfig, restAuth, oauthClientSecretVariableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationIcebergRest),
		Steps: []resource.TestStep{
			{
				Config:          config.FromModels(t, catalogIntegrationIcebergRest, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				ResourceName:    catalogIntegrationIcebergRest.ResourceReference(),
				ImportState:     true,
				ImportStateId:   notificationIntegration.ID().Name(),
				ExpectError:     regexp.MustCompile(fmt.Sprintf(`Integration %s is not a CATALOG integration`, notificationIntegration.ID().Name())),
			},
			{
				Config:          config.FromModels(t, catalogIntegrationIcebergRest2, oauthClientSecretVariableModel),
				ConfigVariables: oauthClientSecretConfigVariables,
				ResourceName:    catalogIntegrationIcebergRest2.ResourceReference(),
				ImportState:     true,
				ImportStateId:   catalogIntegrationObjectStorageId.Name(),
				ExpectError:     regexp.MustCompile(fmt.Sprintf(`invalid catalog source type, expected %s, got %s`, sdk.CatalogIntegrationCatalogSourceTypeIcebergRest, sdk.CatalogIntegrationCatalogSourceTypeObjectStore)),
			},
		},
	})
}
