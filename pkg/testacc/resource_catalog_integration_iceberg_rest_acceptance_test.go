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

func TestAcc_CatalogIntegrationIcebergRest_BasicUseCaseOAuth(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	catalogUri := "https://api.tabular.io/ws"
	newCatalogUri := "https://onelake.table.fabric.microsoft.com/iceberg"
	catalogName := random.AlphanumericN(15)
	newCatalogName := random.AlphanumericN(15)
	prefix := "prefix"
	newPrefix := "prefix2"

	catalogNamespace := random.AlphanumericN(15)
	newCatalogNamespace := random.AlphanumericN(15)
	externalCatalogNamespace := random.AlphanumericN(15)

	oAuthTokenUri := catalogUri + "/v2/oauth/tokens"
	newOAuthTokenUri := newCatalogUri + "/v3/oauth/tokens"
	oAuthClientId := random.AlphanumericN(15)
	newOAuthClientId := random.AlphanumericN(15)
	oAuthClientSecret := random.AlphanumericN(15)
	newOAuthClientSecret := random.AlphanumericN(15)
	oAuthAllowedScope := "PRINCIPAL_ROLE:ALL"
	additionalOAuthAllowedScope := "DUMMY"

	bearerToken := random.AlphanumericN(32)

	comment := random.Comment()
	externalComment := random.Comment()

	refreshIntervalSeconds := random.IntRange(30, 86400)
	externalRefreshIntervalSeconds := random.IntRange(30, 86400)

	basicRestAuth := *sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, oAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}})
	completeRestAuth := *sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, oAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}).
		WithOauthTokenUri(oAuthTokenUri)
	changedRestAuth := *sdk.NewOAuthRestAuthenticationRequest(newOAuthClientId, newOAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}, {Value: additionalOAuthAllowedScope}}).
		WithOauthTokenUri(newOAuthTokenUri)

	basicRestConfig := *sdk.NewIcebergRestRestConfigRequest(catalogUri)
	completeRestConfig := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithPrefix(prefix).
		WithCatalogName(catalogName).
		WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
		WithAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)
	changedRestConfig := *sdk.NewIcebergRestRestConfigRequest(newCatalogUri).
		WithPrefix(newPrefix).
		WithCatalogName(newCatalogName).
		WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePrivate).
		WithAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials)

	basic := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, basicRestConfig, basicRestAuth)

	altered := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), true, basicRestConfig, basicRestAuth).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds)

	allAttributes := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, completeRestConfig, completeRestAuth).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(catalogNamespace)

	withChangedCatalogNamespace := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, completeRestConfig, completeRestAuth).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withChangedRestConfig := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, changedRestConfig, completeRestAuth).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withChangedOAuthClientSecret := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, changedRestConfig, *sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, newOAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}).
		WithOauthTokenUri(oAuthTokenUri)).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withChangedRestAuth := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, changedRestConfig, changedRestAuth).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withBearerToken := model.CatalogIntegrationIcebergRestBearer("t", id.Name(), false, changedRestConfig, *sdk.NewBearerRestAuthenticationRequest(bearerToken)).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasEnabledString(r.BooleanFalse).
			HasCommentEmpty().
			HasNoRefreshIntervalSeconds().
			HasCatalogNamespaceEmpty().
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               "",
				CatalogApiType:       "",
				CatalogName:          "",
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
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergREST).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasCatalogNamespace(""),
		resourceshowoutputassert.IcebergRestRestConfigDescribeOutput(t, ref).
			HasCatalogUri(catalogUri).
			HasPrefix("").
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasCatalogName("").
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
		resourceshowoutputassert.OAuthRestAuthenticationDescribeOutput(t, ref, "oauth_rest_authentication").
			HasOauthTokenUri(catalogUri + "/v1/oauth/tokens").
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(oAuthAllowedScope),
	}

	basicAssertionsWithRefreshIntervalZero := append(
		[]assert.TestCheckFuncProvider{
			resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
				HasName(id.Name()).
				HasEnabledString(r.BooleanFalse).
				HasCommentEmpty().
				HasRefreshIntervalSeconds(0).
				HasCatalogNamespaceEmpty().
				HasRestConfig(&sdk.IcebergRestRestConfigDetails{
					CatalogUri:           catalogUri,
					Prefix:               "",
					CatalogApiType:       "",
					CatalogName:          "",
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
			HasEnabledString(r.BooleanTrue).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespaceEmpty().
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           catalogUri,
				Prefix:               "",
				CatalogApiType:       "",
				CatalogName:          "",
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
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergREST).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(true).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(""),
		resourceshowoutputassert.IcebergRestRestConfigDescribeOutput(t, ref).
			HasCatalogUri(catalogUri).
			HasPrefix("").
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasCatalogName("").
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
		resourceshowoutputassert.OAuthRestAuthenticationDescribeOutput(t, ref, "oauth_rest_authentication").
			HasOauthTokenUri(catalogUri + "/v1/oauth/tokens").
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(oAuthAllowedScope),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
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
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergREST).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(catalogNamespace),
		resourceshowoutputassert.IcebergRestRestConfigDescribeOutput(t, ref).
			HasCatalogUri(catalogUri).
			HasPrefix(prefix).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasCatalogName(catalogName).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
		resourceshowoutputassert.OAuthRestAuthenticationDescribeOutput(t, ref, "oauth_rest_authentication").
			HasOauthTokenUri(oAuthTokenUri).
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(oAuthAllowedScope),
	}

	forceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
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
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergREST).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(newCatalogNamespace),
		resourceshowoutputassert.IcebergRestRestConfigDescribeOutput(t, ref).
			HasCatalogUri(catalogUri).
			HasPrefix(prefix).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasCatalogName(catalogName).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
		resourceshowoutputassert.OAuthRestAuthenticationDescribeOutput(t, ref, "oauth_rest_authentication").
			HasOauthTokenUri(oAuthTokenUri).
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(oAuthAllowedScope),
	}

	moreForceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(newCatalogNamespace).
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           newCatalogUri,
				Prefix:               newPrefix,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePrivate,
				CatalogName:          newCatalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeVendedCredentials,
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
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergREST).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(newCatalogNamespace),
		resourceshowoutputassert.IcebergRestRestConfigDescribeOutput(t, ref).
			HasCatalogUri(newCatalogUri).
			HasPrefix(newPrefix).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePrivate).
			HasCatalogName(newCatalogName).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials),
		resourceshowoutputassert.OAuthRestAuthenticationDescribeOutput(t, ref, "oauth_rest_authentication").
			HasOauthTokenUri(oAuthTokenUri).
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(oAuthAllowedScope),
	}

	moreForceNewAssertionsWithChangedSecret := append(
		[]assert.TestCheckFuncProvider{
			resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
				HasName(id.Name()).
				HasEnabledString(r.BooleanFalse).
				HasComment(comment).
				HasRefreshIntervalSeconds(refreshIntervalSeconds).
				HasCatalogNamespace(newCatalogNamespace).
				HasRestConfig(&sdk.IcebergRestRestConfigDetails{
					CatalogUri:           newCatalogUri,
					Prefix:               newPrefix,
					CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePrivate,
					CatalogName:          newCatalogName,
					AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeVendedCredentials,
				}).
				HasOauthRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
					OauthTokenUri:      oAuthTokenUri,
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
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(newCatalogNamespace).
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           newCatalogUri,
				Prefix:               newPrefix,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePrivate,
				CatalogName:          newCatalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeVendedCredentials,
			}).
			HasOauthRestAuthentication(&sdk.OAuthRestAuthenticationDetails{
				OauthTokenUri:      newOAuthTokenUri,
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
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergREST).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasComment(comment).
			HasCatalogNamespace(newCatalogNamespace),
		resourceshowoutputassert.IcebergRestRestConfigDescribeOutput(t, ref).
			HasCatalogUri(newCatalogUri).
			HasPrefix(newPrefix).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePrivate).
			HasCatalogName(newCatalogName).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials),
		resourceshowoutputassert.OAuthRestAuthenticationDescribeOutput(t, ref, "oauth_rest_authentication").
			HasOauthTokenUri(newOAuthTokenUri).
			HasOauthClientId(newOAuthClientId).
			HasOauthAllowedScopes(oAuthAllowedScope, additionalOAuthAllowedScope),
	}

	withBearerTokenAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
			HasEnabledString(r.BooleanFalse).
			HasComment(comment).
			HasRefreshIntervalSeconds(refreshIntervalSeconds).
			HasCatalogNamespace(newCatalogNamespace).
			HasRestConfig(&sdk.IcebergRestRestConfigDetails{
				CatalogUri:           newCatalogUri,
				Prefix:               newPrefix,
				CatalogApiType:       sdk.CatalogIntegrationCatalogApiTypePrivate,
				CatalogName:          newCatalogName,
				AccessDelegationMode: sdk.CatalogIntegrationAccessDelegationModeVendedCredentials,
			}).
			HasOauthRestAuthenticationEmpty().
			HasBearerRestAuthentication(&sdk.BearerRestAuthenticationDetails{bearerToken}).
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
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import
			{
				Config:                  config.FromModels(t, basic),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rest_config", "oauth_rest_authentication"},
			},
			// Change alterable props
			{
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
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import
			{
				Config:                  config.FromModels(t, allAttributes),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"refresh_interval_seconds", "rest_config", "oauth_rest_authentication"},
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
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Change force new "catalog_namespace" prop
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, withChangedCatalogNamespace),
				Check:  assertThat(t, forceNewAssertions...),
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
				Config: config.FromModels(t, withChangedCatalogNamespace),
				Check:  assertThat(t, forceNewAssertions...),
			},
			// Change force new props in "rest_config"
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, withChangedRestConfig),
				Check:  assertThat(t, moreForceNewAssertions...),
			},
			// Change force new props in "rest_config" externally
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithOrReplace(true).
						WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
							WithRestConfig(completeRestConfig).
							WithOAuthRestAuthentication(completeRestAuth))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "rest_config.0.catalog_uri", sdk.String(newCatalogUri), sdk.String(catalogUri)),
						planchecks.ExpectChange(ref, "rest_config.0.catalog_uri", tfjson.ActionDelete, sdk.String(catalogUri), sdk.String(newCatalogUri)),
						planchecks.ExpectDrift(ref, "rest_config.0.prefix", sdk.String(newPrefix), sdk.String(prefix)),
						planchecks.ExpectChange(ref, "rest_config.0.prefix", tfjson.ActionDelete, sdk.String(prefix), sdk.String(newPrefix)),
						planchecks.ExpectDrift(ref, "rest_config.0.catalog_name", sdk.String(newCatalogName), sdk.String(catalogName)),
						planchecks.ExpectChange(ref, "rest_config.0.catalog_name", tfjson.ActionDelete, sdk.String(catalogName), sdk.String(newCatalogName)),
						planchecks.ExpectDrift(ref, "rest_config.0.catalog_api_type", sdk.String(string(sdk.CatalogIntegrationCatalogApiTypePrivate)), sdk.String(string(sdk.CatalogIntegrationCatalogApiTypePublic))),
						planchecks.ExpectChange(ref, "rest_config.0.catalog_api_type", tfjson.ActionDelete, sdk.String(string(sdk.CatalogIntegrationCatalogApiTypePublic)), sdk.String(string(sdk.CatalogIntegrationCatalogApiTypePrivate))),
						planchecks.ExpectDrift(ref, "rest_config.0.access_delegation_mode", sdk.String(string(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials)), sdk.String(string(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials))),
						planchecks.ExpectChange(ref, "rest_config.0.access_delegation_mode", tfjson.ActionDelete, sdk.String(string(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)), sdk.String(string(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials))),
					},
				},
				Config: config.FromModels(t, withChangedRestConfig),
				Check:  assertThat(t, moreForceNewAssertions...),
			},
			// Change alterable "oauth_client_secret" prop in "rest_authentication"
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withChangedOAuthClientSecret),
				Check:  assertThat(t, moreForceNewAssertionsWithChangedSecret...),
			},
			// Change force new props in "oauth_rest_authentication"
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, withChangedRestAuth),
				Check:  assertThat(t, evenMoreForceNewAssertions...),
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
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "oauth_rest_authentication.0.oauth_token_uri", sdk.String(newOAuthTokenUri), sdk.String(oAuthTokenUri)),
						planchecks.ExpectChange(ref, "oauth_rest_authentication.0.oauth_token_uri", tfjson.ActionDelete, sdk.String(oAuthTokenUri), sdk.String(newOAuthTokenUri)),
						planchecks.ExpectDrift(ref, "oauth_rest_authentication.0.oauth_client_id", sdk.String(newOAuthClientId), sdk.String(oAuthClientId)),
						planchecks.ExpectChange(ref, "oauth_rest_authentication.0.oauth_client_id", tfjson.ActionDelete, sdk.String(oAuthClientId), sdk.String(newOAuthClientId)),
						planchecks.ExpectDrift(ref, "oauth_rest_authentication.0.oauth_allowed_scopes", sdk.String(fmt.Sprintf("[%s %s]", oAuthAllowedScope, additionalOAuthAllowedScope)), sdk.String(fmt.Sprintf("[%s]", oAuthAllowedScope))),
						planchecks.ExpectChange(ref, "oauth_rest_authentication.0.oauth_allowed_scopes", tfjson.ActionDelete, sdk.String(fmt.Sprintf("[%s]", oAuthAllowedScope)), sdk.String(fmt.Sprintf("[%s %s]", oAuthAllowedScope, additionalOAuthAllowedScope))),
					},
				},
				Config: config.FromModels(t, withChangedRestAuth),
				Check:  assertThat(t, evenMoreForceNewAssertions...),
			},
			// Change to different authentication type
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, withBearerToken),
				Check:  assertThat(t, withBearerTokenAssertions...),
			},
		},
	})
}

func TestAcc_CatalogIntegrationIcebergRest_BasicUseCaseBearer(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	catalogUri := "https://api.tabular.io/ws"

	bearerToken1 := random.AlphanumericN(32)
	bearerToken2 := random.AlphanumericN(32)
	sigV4IamRole := "arn:aws:iam::123456789012:role/sigv4-role"

	basicRestCfg := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway)

	basic := model.CatalogIntegrationIcebergRestBearer("t", id.Name(), false, basicRestCfg, *sdk.NewBearerRestAuthenticationRequest(bearerToken1))
	updated := model.CatalogIntegrationIcebergRestBearer("t", id.Name(), false, basicRestCfg, *sdk.NewBearerRestAuthenticationRequest(bearerToken2))
	withSigV4 := model.CatalogIntegrationIcebergRestSigV4("t", id.Name(), false, basicRestCfg, *sdk.NewSigV4RestAuthenticationRequest(sigV4IamRole))

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
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergREST).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasCatalogNamespace(""),
		resourceshowoutputassert.IcebergRestRestConfigDescribeOutput(t, ref).
			HasCatalogUri(catalogUri).
			HasPrefix("").
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway).
			HasCatalogName("").
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
	}

	basicAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
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
			HasBearerRestAuthentication(&sdk.BearerRestAuthenticationDetails{BearerToken: bearerToken1}).
			HasSigv4RestAuthenticationEmpty(),
	}, commonAssertions...)

	updatedAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
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
			HasBearerRestAuthentication(&sdk.BearerRestAuthenticationDetails{BearerToken: bearerToken2}).
			HasSigv4RestAuthenticationEmpty(),
	}, commonAssertions...)

	sigV4Assertions := append([]assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
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
				Sigv4SigningRegion: "",
				Sigv4ExternalId:    "",
			}),
		resourceshowoutputassert.SigV4RestAuthenticationDescribeOutput(t, ref).
			// Don't check sigv4_signing_region, as its default value depends on the current region name
			HasSigv4IamRole(sigV4IamRole),
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
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Change alterable "bearer_token" prop in "bearer_rest_authentication"
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
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
				Config: config.FromModels(t, withSigV4),
				Check:  assertThat(t, sigV4Assertions...),
			},
		},
	})
}

func TestAcc_CatalogIntegrationIcebergRest_BasicUseCaseSigV4(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	catalogUri := "https://api.tabular.io/ws"

	sigV4IamRole := "arn:aws:iam::123456789012:role/sigv4-role-1"
	sigV4SigningRegion := "us-west-2"
	newSigV4IamRole := "arn:aws:iam::123456789012:role/sigv4-role-2"
	newSigV4SigningRegion := "eu-west-1"
	newSigV4ExternalId := "external-id-2"

	bearerToken := random.AlphanumericN(32)

	basicRestCfg := *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithCatalogApiType(sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway)

	basicSigV4Auth := *sdk.NewSigV4RestAuthenticationRequest(sigV4IamRole).
		WithSigv4SigningRegion(sigV4SigningRegion)

	updatedSigV4Auth := *sdk.NewSigV4RestAuthenticationRequest(newSigV4IamRole).
		WithSigv4SigningRegion(newSigV4SigningRegion).
		WithSigv4ExternalId(newSigV4ExternalId)

	basic := model.CatalogIntegrationIcebergRestSigV4("t", id.Name(), false, basicRestCfg, basicSigV4Auth)

	updated := model.CatalogIntegrationIcebergRestSigV4("t", id.Name(), false, basicRestCfg, updatedSigV4Auth)

	withBearerToken := model.CatalogIntegrationIcebergRestBearer("t", id.Name(), false, basicRestCfg, *sdk.NewBearerRestAuthenticationRequest(bearerToken))

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
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
		resourceshowoutputassert.SigV4RestAuthenticationDescribeOutput(t, ref).
			HasSigv4IamRole(sigV4IamRole).
			HasSigv4SigningRegion(sigV4SigningRegion),
		resourceshowoutputassert.CatalogIntegrationIcebergRestDescribeOutput(t, ref).
			HasId(id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergREST).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasCatalogNamespace(""),
		resourceshowoutputassert.IcebergRestRestConfigDescribeOutput(t, ref).
			HasCatalogUri(catalogUri).
			HasPrefix("").
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway).
			HasCatalogName("").
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
	}

	updatedAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
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
		resourceshowoutputassert.SigV4RestAuthenticationDescribeOutput(t, ref).
			HasSigv4IamRole(newSigV4IamRole).
			HasSigv4SigningRegion(newSigV4SigningRegion),
	}, basicAssertions[2:]...)

	withBearerTokenAssertions := append([]assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationIcebergRestResource(t, ref).
			HasName(id.Name()).
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
			HasBearerRestAuthentication(&sdk.BearerRestAuthenticationDetails{bearerToken}).
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
				Config: config.FromModels(t, withBearerToken),
				Check:  assertThat(t, withBearerTokenAssertions...),
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

	refreshIntervalNonPositive := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, restAuth).
		WithRefreshIntervalSeconds(0)

	emptyCatalogNamespace := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, restAuth).
		WithCatalogNamespace("")

	emptyCatalogUri := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, *sdk.NewIcebergRestRestConfigRequest(""), restAuth)

	emptyCatalogName := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithCatalogName(""), restAuth)

	invalidCatalogApiType := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithCatalogName(catalogName).
		WithCatalogApiType("invalid"), restAuth)

	invalidAccessDelegationMode := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, *sdk.NewIcebergRestRestConfigRequest(catalogUri).
		WithCatalogName(catalogName).
		WithAccessDelegationMode("invalid"), restAuth)

	emptyOAuthTokenUri := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, *sdk.NewOAuthRestAuthenticationRequest("my_client_id", "my_client_secret", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}}).
		WithOauthTokenUri(""))

	emptyOAuthClientId := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, *sdk.NewOAuthRestAuthenticationRequest("", "my_client_secret", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}}))

	emptyOAuthClientSecret := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, *sdk.NewOAuthRestAuthenticationRequest("my_client_id", "", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}}))

	emptyOAuthScopes := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, *sdk.NewOAuthRestAuthenticationRequest("my_client_id", "my_client_secret", []sdk.StringListItemWrapper{}))

	emptyBearerToken := model.CatalogIntegrationIcebergRestBearer("t", id.Name(), false, restConfig, *sdk.NewBearerRestAuthenticationRequest(""))

	emptySigV4IamRole := model.CatalogIntegrationIcebergRestSigV4("t", id.Name(), false, restConfig, *sdk.NewSigV4RestAuthenticationRequest(""))

	emptySigV4SigningRegion := model.CatalogIntegrationIcebergRestSigV4("t", id.Name(), false, restConfig, *sdk.NewSigV4RestAuthenticationRequest("arn:aws:iam::123456789012:role/role").
		WithSigv4SigningRegion(""),
	)

	emptySigV4ExternalId := model.CatalogIntegrationIcebergRestSigV4("t", id.Name(), false, restConfig, *sdk.NewSigV4RestAuthenticationRequest("arn:aws:iam::123456789012:role/role").
		WithSigv4ExternalId(""),
	)

	noAuthentication := model.CatalogIntegrationIcebergRest("t", id.Name(), false, []sdk.IcebergRestRestConfigRequest{restConfig})

	oauthAndBearer := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, restAuth).
		WithBearerRestAuthentication(*sdk.NewBearerRestAuthenticationRequest("token"))

	oauthAndSigv4 := model.CatalogIntegrationIcebergRestOAuth("t", id.Name(), false, restConfig, restAuth).
		WithSigV4RestAuthentication(*sdk.NewSigV4RestAuthenticationRequest("arn:aws:iam::123456789012:role/role"))

	bearerAndSigv4 := model.CatalogIntegrationIcebergRest("t", id.Name(), false, []sdk.IcebergRestRestConfigRequest{restConfig}).
		WithBearerRestAuthentication(*sdk.NewBearerRestAuthenticationRequest("token")).
		WithSigV4RestAuthentication(*sdk.NewSigV4RestAuthenticationRequest("arn:aws:iam::123456789012:role/role"))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationIcebergRest),
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
				ExpectError: regexp.MustCompile(`invalid catalog api type: INVALID`),
			},
			{
				Config:      config.FromModels(t, invalidAccessDelegationMode),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid access delegation mode: INVALID`),
			},
			{
				Config:      config.FromModels(t, emptyOAuthTokenUri),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "oauth_rest_authentication.0\.oauth_token_uri" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptyOAuthClientId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "oauth_rest_authentication.0\.oauth_client_id" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptyOAuthClientSecret),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "oauth_rest_authentication.0\.oauth_client_secret" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, emptyOAuthScopes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Not enough list items`),
			},
			{
				Config:      config.FromModels(t, emptyBearerToken),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "bearer_rest_authentication.0\.bearer_token" to not be an empty string`),
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
				Config:      config.FromModels(t, oauthAndBearer),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
			{
				Config:      config.FromModels(t, oauthAndSigv4),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
			{
				Config:      config.FromModels(t, bearerAndSigv4),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
		},
	})
}

func TestAcc_CatalogIntegrationIcebergRest_ImportValidation(t *testing.T) {
	restConfig := *sdk.NewIcebergRestRestConfigRequest("https://api.tabular.io/ws").
		WithCatalogName("my_catalog_name")
	restAuth := *sdk.NewOAuthRestAuthenticationRequest("my_client_id", "my_client_secret", []sdk.StringListItemWrapper{{Value: "PRINCIPAL_ROLE:ALL"}})

	notificationIntegration, notificationIntegrationCleanup := testClient().NotificationIntegration.Create(t)
	t.Cleanup(notificationIntegrationCleanup)

	catalogIntegrationObjectStorageId, catalogIntegrationObjectStorageCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogIntegrationObjectStorageCleanup)

	catalogIntegrationIcebergRest := model.CatalogIntegrationIcebergRestOAuth("t", notificationIntegration.ID().Name(), false, restConfig, restAuth)
	catalogIntegrationIcebergRest2 := model.CatalogIntegrationIcebergRestOAuth("t", catalogIntegrationObjectStorageId.Name(), false, restConfig, restAuth)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.CatalogIntegrationIcebergRest),
		Steps: []resource.TestStep{
			{
				Config:        config.FromModels(t, catalogIntegrationIcebergRest),
				ResourceName:  catalogIntegrationIcebergRest.ResourceReference(),
				ImportState:   true,
				ImportStateId: notificationIntegration.ID().Name(),
				ExpectError:   regexp.MustCompile(fmt.Sprintf(`Integration %s is not a CATALOG integration`, notificationIntegration.ID().Name())),
			},
			{
				Config:        config.FromModels(t, catalogIntegrationIcebergRest2),
				ResourceName:  catalogIntegrationIcebergRest2.ResourceReference(),
				ImportState:   true,
				ImportStateId: catalogIntegrationObjectStorageId.Name(),
				ExpectError:   regexp.MustCompile(fmt.Sprintf(`invalid catalog source type, expected %s, got %s`, sdk.CatalogIntegrationCatalogSourceTypeIcebergREST, sdk.CatalogIntegrationCatalogSourceTypeObjectStorage)),
			},
		},
	})
}
