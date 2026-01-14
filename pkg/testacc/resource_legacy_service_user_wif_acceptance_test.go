//go:build account_level_tests

package testacc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_LegacyServiceUser_WIF_OIDC(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	subject := fmt.Sprintf("system:serviceaccount:namespace:%s", random.AlphaN(10))

	userModelWithOIDC := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityOidc(
			testvars.OidcIssuer,
			subject,
		)
	userModelWithOIDCAndAudienceList := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityOidc(
			testvars.OidcIssuer,
			subject,
			testvars.OidcIssuer+"/o/oauth2/auth",
		)
	userModelWithOIDCAndAudienceListTwoItems := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityOidc(
			testvars.OidcIssuer,
			subject,
			testvars.OidcIssuer+"/o/oauth2/auth",
			testvars.OidcIssuer+"/o/oauth2/token",
		)

	userModelWithoutWIF := model.LegacyServiceUser("w", id.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.LegacyServiceUser),
		Steps: []resource.TestStep{
			// create optional
			{
				Config: config.FromModels(t, userModelWithOIDC),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithOIDC.ResourceReference()).
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityOidc(testvars.OidcIssuer, subject),
					resourceshowoutputassert.UserShowOutput(t, userModelWithOIDC.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserWorkloadIdentityAuthenticationMethods(t, id, "DEFAULT").
						HasName("DEFAULT").
						HasType(sdk.WIFTypeOIDC).
						HasNoComment().
						HasLastUsedNotEmpty().
						HasCreatedOnNotEmpty().
						HasOidcAdditionalInfo(sdk.UserWorkloadIdentityAuthenticationMethodsOidcAdditionalInfo{
							Issuer:       testvars.OidcIssuer,
							Subject:      subject,
							AudienceList: []string{},
						}),
				),
			},
			// UPDATE - REMOVE WIF
			{
				Config: config.FromModels(t, userModelWithoutWIF),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithoutWIF.ResourceReference()).
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityEmpty(),
					resourceshowoutputassert.UserShowOutput(t, userModelWithoutWIF.ResourceReference()).
						HasHasWorkloadIdentity(false),
					assert.Check(func(_ *terraform.State) error {
						id := helpers.NewUserWorkloadIdentityAuthenticationMethodsObjectIdentifier(id, "DEFAULT")
						_, err := testClient().User.ShowUserWorkloadIdentityAuthenticationMethodOptions(t, id)
						if !errors.Is(err, collections.ErrObjectNotFound) {
							return err
						}
						return nil
					}),
				),
			},
			// UPDATE - ADD WIF BACK
			{
				Config: config.FromModels(t, userModelWithOIDCAndAudienceList),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithOIDCAndAudienceList.ResourceReference()).
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityOidc(testvars.OidcIssuer, subject, testvars.OidcIssuer+"/o/oauth2/auth"),
					objectassert.UserWorkloadIdentityAuthenticationMethods(t, id, "DEFAULT").
						HasName("DEFAULT").
						HasType(sdk.WIFTypeOIDC).
						HasNoComment().
						HasLastUsedNotEmpty().
						HasCreatedOnNotEmpty().
						HasOidcAdditionalInfo(sdk.UserWorkloadIdentityAuthenticationMethodsOidcAdditionalInfo{
							Issuer:       testvars.OidcIssuer,
							Subject:      subject,
							AudienceList: []string{testvars.OidcIssuer + "/o/oauth2/auth"},
						}),
				),
			},
			// ALTER - CHANGE AUDIENCE LIST EXTERNALLY
			{
				PreConfig: func() {
					testClient().User.Alter(t, id, &sdk.AlterUserOptions{
						Set: &sdk.UserSet{
							ObjectProperties: &sdk.UserAlterObjectProperties{
								UserObjectProperties: sdk.UserObjectProperties{
									WorkloadIdentity: &sdk.UserObjectWorkloadIdentityProperties{
										OidcType: &sdk.UserObjectWorkloadIdentityOidc{
											Issuer:  sdk.String(testvars.OidcIssuer),
											Subject: sdk.String(subject),
											OidcAudienceList: []sdk.StringListItemWrapper{
												{
													Value: testvars.OidcIssuer + "/o/oauth2/changed",
												},
											},
										},
									},
								},
							},
						},
					})
				},
				Config: config.FromModels(t, userModelWithOIDCAndAudienceList),
				Check: assertThat(t,
					objectassert.UserWorkloadIdentityAuthenticationMethods(t, id, "DEFAULT").
						HasOidcAdditionalInfo(sdk.UserWorkloadIdentityAuthenticationMethodsOidcAdditionalInfo{
							Issuer:       testvars.OidcIssuer,
							Subject:      subject,
							AudienceList: []string{testvars.OidcIssuer + "/o/oauth2/auth"},
						}),
				),
			},
			// ALTER - CHANGE AUDIENCE LIST
			{
				Config: config.FromModels(t, userModelWithOIDCAndAudienceListTwoItems),
				Check: assertThat(t,
					objectassert.UserWorkloadIdentityAuthenticationMethods(t, id, "DEFAULT").
						HasOidcAdditionalInfo(sdk.UserWorkloadIdentityAuthenticationMethodsOidcAdditionalInfo{
							Issuer:       testvars.OidcIssuer,
							Subject:      subject,
							AudienceList: []string{testvars.OidcIssuer + "/o/oauth2/auth", testvars.OidcIssuer + "/o/oauth2/token"},
						}),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelWithOIDC.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "mins_to_unlock", "login_name", "display_name", "disabled", "must_change_password", "default_secondary_roles_option"},
			},
		},
	})
}

func TestAcc_LegacyServiceUser_WIF_AWS(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	arn := fmt.Sprintf("arn:aws:iam::%s:role/test-role", random.NumericN(12))
	userModelWithAWS := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityAws(arn)

	userModelWithoutWIF := model.LegacyServiceUser("w", id.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.LegacyServiceUser),
		Steps: []resource.TestStep{
			// CREATE WITH AWS WIF
			{
				Config: config.FromModels(t, userModelWithAWS),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithAWS.ResourceReference()).
						HasNameString(id.Name()),
					resourceshowoutputassert.UserShowOutput(t, userModelWithAWS.ResourceReference()).
						HasHasWorkloadIdentity(true),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelWithAWS.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "mins_to_unlock", "login_name", "display_name", "disabled", "must_change_password", "default_secondary_roles_option", "default_workload_identity"},
			},
			// UPDATE - REMOVE WIF
			{
				Config: config.FromModels(t, userModelWithoutWIF),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithoutWIF.ResourceReference()).
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityEmpty(),
					resourceshowoutputassert.UserShowOutput(t, userModelWithoutWIF.ResourceReference()).
						HasHasWorkloadIdentity(false),
				),
			},
		},
	})
}

func TestAcc_LegacyServiceUser_WIF_GCP(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	subject := random.NumericN(10)

	userModelWithGCP := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityGcp(subject)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.LegacyServiceUser),
		Steps: []resource.TestStep{
			// CREATE WITH GCP WIF
			{
				Config: config.FromModels(t, userModelWithGCP),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithGCP.ResourceReference()).
						HasNameString(id.Name()),
					resourceshowoutputassert.UserShowOutput(t, userModelWithGCP.ResourceReference()).
						HasHasWorkloadIdentity(true),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelWithGCP.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "mins_to_unlock", "login_name", "display_name", "disabled", "must_change_password", "default_secondary_roles_option"},
			},
		},
	})
}

func TestAcc_LegacyServiceUser_WIF_Azure(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	subject := random.AlphaN(10)
	userModelWithAzure := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityAzure(
			testvars.MicrosoftIssuer,
			subject,
		)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.LegacyServiceUser),
		Steps: []resource.TestStep{
			// CREATE WITH Azure WIF
			{
				Config: config.FromModels(t, userModelWithAzure),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithAzure.ResourceReference()).
						HasNameString(id.Name()),
					resourceshowoutputassert.UserShowOutput(t, userModelWithAzure.ResourceReference()).
						HasHasWorkloadIdentity(true),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelWithAzure.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "mins_to_unlock", "login_name", "display_name", "disabled", "must_change_password", "default_secondary_roles_option"},
			},
		},
	})
}

func TestAcc_LegacyServiceUser_WIF_SwitchProvider(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	subject := fmt.Sprintf("system:serviceaccount:namespace:%s", random.AlphaN(10))
	arn := fmt.Sprintf("arn:aws:iam::%s:role/test-role", random.NumericN(12))
	userModelWithOIDC := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityOidc(
			testvars.OidcIssuer,
			subject,
			testvars.OidcIssuer+"/o/oauth2/auth",
		)

	userModelWithAWS := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityAws(arn)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.LegacyServiceUser),
		Steps: []resource.TestStep{
			// Start with OIDC
			{
				Config: config.FromModels(t, userModelWithOIDC),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithOIDC.ResourceReference()).
						HasNameString(id.Name()),
					resourceshowoutputassert.UserShowOutput(t, userModelWithOIDC.ResourceReference()).
						HasHasWorkloadIdentity(true),
				),
			},
			// Switch to AWS
			{
				Config: config.FromModels(t, userModelWithAWS),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithAWS.ResourceReference()).
						HasNameString(id.Name()),
					resourceshowoutputassert.UserShowOutput(t, userModelWithAWS.ResourceReference()).
						HasHasWorkloadIdentity(true),
				),
			},
		},
	})
}

func TestAcc_LegacyServiceUser_WIF_ExternalChange(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	subject := fmt.Sprintf("system:serviceaccount:namespace:%s", random.AlphaN(10))

	userModelWithOIDC := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityOidc(
			testvars.OidcIssuer,
			subject,
			testvars.OidcIssuer+"/o/oauth2/auth",
		)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.LegacyServiceUser),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModelWithOIDC),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithOIDC.ResourceReference()).
						HasNameString(id.Name()),
					resourceshowoutputassert.UserShowOutput(t, userModelWithOIDC.ResourceReference()).
						HasHasWorkloadIdentity(true),
				),
			},
			// External change - unset WIF outside of Terraform
			{
				PreConfig: func() {
					testClient().User.Alter(t, id, &sdk.AlterUserOptions{
						Unset: &sdk.UserUnset{
							ObjectProperties: &sdk.UserObjectPropertiesUnset{
								WorkloadIdentity: sdk.Bool(true),
							},
						},
					})
				},
				Config: config.FromModels(t, userModelWithOIDC),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				// Terraform should detect the drift and re-apply
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithOIDC.ResourceReference()).
						HasNameString(id.Name()),
					resourceshowoutputassert.UserShowOutput(t, userModelWithOIDC.ResourceReference()).
						HasHasWorkloadIdentity(true),
				),
			},
		},
	})
}
