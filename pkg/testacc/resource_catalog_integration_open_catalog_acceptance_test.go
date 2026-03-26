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

func TestAcc_CatalogIntegrationOpenCatalog_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	catalogUri := "https://testorg-testacc.snowflakecomputing.com/polaris/api/catalog"
	privateCatalogUri := "https://testorg-testacc.privatelink.snowflakecomputing.com/polaris/api/catalog"
	catalogName := random.AlphanumericN(15)
	newCatalogName := random.AlphanumericN(15)

	catalogNamespace := random.AlphanumericN(15)
	newCatalogNamespace := random.AlphanumericN(15)
	externalCatalogNamespace := random.AlphanumericN(15)

	oAuthTokenUri := privateCatalogUri + "/v2/oauth/tokens"
	newOAuthTokenUri := catalogUri + "/v3/oauth/tokens"
	oAuthClientId := random.AlphanumericN(15)
	newOAuthClientId := random.AlphanumericN(15)
	oAuthClientSecret := random.AlphanumericN(15)
	newOAuthClientSecret := random.AlphanumericN(15)
	oAuthAllowedScope := "PRINCIPAL_ROLE:ALL"
	additionalOAuthAllowedScope := "DUMMY"

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
			WithOauthTokenUri(newOAuthTokenUri),
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

	basic := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, basicRestAuth, basicRestConfig)

	altered := model.CatalogIntegrationOpenCatalog("t", id.Name(), true, basicRestAuth, basicRestConfig).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds)

	allAttributes := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, completeRestAuth, completeRestConfig).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(catalogNamespace)

	withChangedCatalogNamespace := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, completeRestAuth, completeRestConfig).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withChangedRestConfig := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, completeRestAuth, changedRestConfig).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withChangedOAuthClientSecret := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, []sdk.OAuthRestAuthenticationRequest{
		*sdk.NewOAuthRestAuthenticationRequest(oAuthClientId, newOAuthClientSecret, []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}}).
			WithOauthTokenUri(oAuthTokenUri),
	}, changedRestConfig).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	withChangedRestAuth := model.CatalogIntegrationOpenCatalog("t", id.Name(), false, changedRestAuth, changedRestConfig).
		WithComment(comment).
		WithRefreshIntervalSeconds(refreshIntervalSeconds).
		WithCatalogNamespace(newCatalogNamespace)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
			HasName(id.Name()).
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
		resourceshowoutputassert.OpenCatalogRestConfigDescribeOutput(t, ref).
			HasCatalogUri(catalogUri).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasCatalogName(catalogName).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
		resourceshowoutputassert.OAuthRestAuthenticationDescribeOutput(t, ref).
			HasOauthTokenUri(catalogUri + "/v1/oauth/tokens").
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(oAuthAllowedScope),
	}

	basicAssertionsWithRefreshIntervalZero := append(
		[]assert.TestCheckFuncProvider{
			resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
				HasName(id.Name()).
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
		resourceshowoutputassert.OpenCatalogRestConfigDescribeOutput(t, ref).
			HasCatalogUri(catalogUri).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasCatalogName(catalogName).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
		resourceshowoutputassert.OAuthRestAuthenticationDescribeOutput(t, ref).
			HasOauthTokenUri(catalogUri + "/v1/oauth/tokens").
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(oAuthAllowedScope),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
			HasName(id.Name()).
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
		resourceshowoutputassert.OpenCatalogRestConfigDescribeOutput(t, ref).
			HasCatalogUri(catalogUri).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasCatalogName(catalogName).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
		resourceshowoutputassert.OAuthRestAuthenticationDescribeOutput(t, ref).
			HasOauthTokenUri(oAuthTokenUri).
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(oAuthAllowedScope),
	}

	forceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
			HasName(id.Name()).
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
		resourceshowoutputassert.OpenCatalogRestConfigDescribeOutput(t, ref).
			HasCatalogUri(catalogUri).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasCatalogName(catalogName).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials),
		resourceshowoutputassert.OAuthRestAuthenticationDescribeOutput(t, ref).
			HasOauthTokenUri(oAuthTokenUri).
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(oAuthAllowedScope),
	}

	moreForceNewAssertions := []assert.TestCheckFuncProvider{
		resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
			HasName(id.Name()).
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
		resourceshowoutputassert.OpenCatalogRestConfigDescribeOutput(t, ref).
			HasCatalogUri(privateCatalogUri).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePrivate).
			HasCatalogName(newCatalogName).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials),
		resourceshowoutputassert.OAuthRestAuthenticationDescribeOutput(t, ref).
			HasOauthTokenUri(oAuthTokenUri).
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(oAuthAllowedScope),
	}

	moreForceNewAssertionsWithChangedSecret := append(
		[]assert.TestCheckFuncProvider{
			resourceassert.CatalogIntegrationOpenCatalogResource(t, ref).
				HasName(id.Name()).
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
					OauthTokenUri:      oAuthTokenUri,
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
				OauthTokenUri:      newOAuthTokenUri,
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
		resourceshowoutputassert.OpenCatalogRestConfigDescribeOutput(t, ref).
			HasCatalogUri(privateCatalogUri).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePrivate).
			HasCatalogName(newCatalogName).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials),
		resourceshowoutputassert.OAuthRestAuthenticationDescribeOutput(t, ref).
			HasOauthTokenUri(newOAuthTokenUri).
			HasOauthClientId(newOAuthClientId).
			HasOauthAllowedScopes(oAuthAllowedScope, additionalOAuthAllowedScope),
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
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import
			{
				Config:                  config.FromModels(t, basic),
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
				ImportStateVerifyIgnore: []string{"refresh_interval_seconds", "rest_config", "rest_authentication"},
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
			// Change force new props in "rest_authentication"
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, withChangedRestAuth),
				Check:  assertThat(t, evenMoreForceNewAssertions...),
			},
			// Change force new props in "rest_authentication" externally
			{
				PreConfig: func() {
					createRequest := sdk.NewCreateCatalogIntegrationRequest(id, false).
						WithOrReplace(true).
						WithOpenCatalogCatalogSourceParams(*sdk.NewOpenCatalogParamsRequest().
							WithRestConfig(changedRestConfig[0]).
							WithRestAuthentication(completeRestAuth[0]))
					testClient().CatalogIntegration.CreateFunc(t, createRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectDrift(ref, "rest_authentication.0.oauth_token_uri", sdk.String(newOAuthTokenUri), sdk.String(oAuthTokenUri)),
						planchecks.ExpectChange(ref, "rest_authentication.0.oauth_token_uri", tfjson.ActionDelete, sdk.String(oAuthTokenUri), sdk.String(newOAuthTokenUri)),
						planchecks.ExpectDrift(ref, "rest_authentication.0.oauth_client_id", sdk.String(newOAuthClientId), sdk.String(oAuthClientId)),
						planchecks.ExpectChange(ref, "rest_authentication.0.oauth_client_id", tfjson.ActionDelete, sdk.String(oAuthClientId), sdk.String(newOAuthClientId)),
						planchecks.ExpectDrift(ref, "rest_authentication.0.oauth_allowed_scopes", sdk.String(fmt.Sprintf("[%s %s]", oAuthAllowedScope, additionalOAuthAllowedScope)), sdk.String(fmt.Sprintf("[%s]", oAuthAllowedScope))),
						planchecks.ExpectChange(ref, "rest_authentication.0.oauth_allowed_scopes", tfjson.ActionDelete, sdk.String(fmt.Sprintf("[%s]", oAuthAllowedScope)), sdk.String(fmt.Sprintf("[%s %s]", oAuthAllowedScope, additionalOAuthAllowedScope))),
					},
				},
				Config: config.FromModels(t, withChangedRestAuth),
				Check:  assertThat(t, evenMoreForceNewAssertions...),
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

func TestAcc_CatalogIntegrationOpenCatalog_ImportValidation(t *testing.T) {
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
				ExpectError:   regexp.MustCompile(fmt.Sprintf(`invalid catalog source type, expected %s, got %s`, sdk.CatalogIntegrationCatalogSourceTypePolaris, sdk.CatalogIntegrationCatalogSourceTypeObjectStorage)),
			},
		},
	})
}
