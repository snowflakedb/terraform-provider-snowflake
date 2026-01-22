//go:build account_level_tests

package testacc

import (
	"errors"
	"fmt"
	"regexp"
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
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
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
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
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
					testClient().User.SetOidcWorkloadIdentity(t, id, testvars.OidcIssuer, subject, testvars.OidcIssuer+"/o/oauth2/changed")
				},
				Config: config.FromModels(t, userModelWithOIDCAndAudienceList),
				Check: assertThat(t,
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
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
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
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
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityAws(arn),
					resourceshowoutputassert.UserShowOutput(t, userModelWithAWS.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
						HasName("DEFAULT").
						HasType(sdk.WIFTypeAWS).
						HasNoComment().
						HasLastUsedNotEmpty().
						HasCreatedOnNotEmpty(),
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
	changedSubject := random.NumericN(10)

	userModelWithGCP := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityGcp(subject)

	userModelWithoutWIF := model.LegacyServiceUser("w", id.Name())

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
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityGcp(subject),
					resourceshowoutputassert.UserShowOutput(t, userModelWithGCP.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
						HasName("DEFAULT").
						HasType(sdk.WIFTypeGCP).
						HasNoComment().
						HasLastUsedNotEmpty().
						HasCreatedOnNotEmpty().
						HasGcpAdditionalInfo(sdk.UserWorkloadIdentityAuthenticationMethodsGcpAdditionalInfo{
							Subject: subject,
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
						wifId := helpers.NewUserWorkloadIdentityAuthenticationMethodsObjectIdentifier(id, "DEFAULT")
						_, err := testClient().User.ShowUserWorkloadIdentityAuthenticationMethodOptions(t, wifId)
						if !errors.Is(err, collections.ErrObjectNotFound) {
							return err
						}
						return nil
					}),
				),
			},
			// UPDATE - ADD WIF BACK
			{
				Config: config.FromModels(t, userModelWithGCP),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithGCP.ResourceReference()).
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityGcp(subject),
					resourceshowoutputassert.UserShowOutput(t, userModelWithGCP.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
						HasName("DEFAULT").
						HasType(sdk.WIFTypeGCP).
						HasNoComment().
						HasLastUsedNotEmpty().
						HasCreatedOnNotEmpty().
						HasGcpAdditionalInfo(sdk.UserWorkloadIdentityAuthenticationMethodsGcpAdditionalInfo{
							Subject: subject,
						}),
				),
			},
			// External change - change subject externally
			{
				PreConfig: func() {
					testClient().User.SetGcpWorkloadIdentity(t, id, changedSubject)
				},
				Config: config.FromModels(t, userModelWithGCP),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithGCP.ResourceReference()).
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityGcp(subject),
					resourceshowoutputassert.UserShowOutput(t, userModelWithGCP.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
						HasName("DEFAULT").
						HasType(sdk.WIFTypeGCP).
						HasNoComment().
						HasLastUsedNotEmpty().
						HasCreatedOnNotEmpty().
						HasGcpAdditionalInfo(sdk.UserWorkloadIdentityAuthenticationMethodsGcpAdditionalInfo{
							Subject: subject,
						}),
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
	changedSubject := random.AlphaN(10)

	userModelWithAzure := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityAzure(
			testvars.MicrosoftIssuer,
			subject,
		)

	userModelWithoutWIF := model.LegacyServiceUser("w", id.Name())

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
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityAzure(testvars.MicrosoftIssuer, subject),
					resourceshowoutputassert.UserShowOutput(t, userModelWithAzure.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
						HasName("DEFAULT").
						HasType(sdk.WIFTypeAzure).
						HasNoComment().
						HasLastUsedNotEmpty().
						HasCreatedOnNotEmpty().
						HasAzureAdditionalInfo(sdk.UserWorkloadIdentityAuthenticationMethodsAzureAdditionalInfo{
							Issuer:  testvars.MicrosoftIssuer,
							Subject: subject,
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
						wifId := helpers.NewUserWorkloadIdentityAuthenticationMethodsObjectIdentifier(id, "DEFAULT")
						_, err := testClient().User.ShowUserWorkloadIdentityAuthenticationMethodOptions(t, wifId)
						if !errors.Is(err, collections.ErrObjectNotFound) {
							return err
						}
						return nil
					}),
				),
			},
			// UPDATE - ADD WIF BACK
			{
				Config: config.FromModels(t, userModelWithAzure),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithAzure.ResourceReference()).
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityAzure(testvars.MicrosoftIssuer, subject),
					resourceshowoutputassert.UserShowOutput(t, userModelWithAzure.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
						HasName("DEFAULT").
						HasType(sdk.WIFTypeAzure).
						HasNoComment().
						HasLastUsedNotEmpty().
						HasCreatedOnNotEmpty().
						HasAzureAdditionalInfo(sdk.UserWorkloadIdentityAuthenticationMethodsAzureAdditionalInfo{
							Issuer:  testvars.MicrosoftIssuer,
							Subject: subject,
						}),
				),
			},
			// External change - change subject externally
			{
				PreConfig: func() {
					testClient().User.SetAzureWorkloadIdentity(t, id, testvars.MicrosoftIssuer, changedSubject)
				},
				Config: config.FromModels(t, userModelWithAzure),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithAzure.ResourceReference()).
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityAzure(testvars.MicrosoftIssuer, subject),
					resourceshowoutputassert.UserShowOutput(t, userModelWithAzure.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
						HasName("DEFAULT").
						HasType(sdk.WIFTypeAzure).
						HasNoComment().
						HasLastUsedNotEmpty().
						HasCreatedOnNotEmpty().
						HasAzureAdditionalInfo(sdk.UserWorkloadIdentityAuthenticationMethodsAzureAdditionalInfo{
							Issuer:  testvars.MicrosoftIssuer,
							Subject: subject,
						}),
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
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityOidc(testvars.OidcIssuer, subject, testvars.OidcIssuer+"/o/oauth2/auth"),
					resourceshowoutputassert.UserShowOutput(t, userModelWithOIDC.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
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
			// Switch to AWS
			{
				Config: config.FromModels(t, userModelWithAWS),
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithAWS.ResourceReference()).
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityAws(arn),
					resourceshowoutputassert.UserShowOutput(t, userModelWithAWS.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
						HasName("DEFAULT").
						HasType(sdk.WIFTypeAWS).
						HasNoComment().
						HasLastUsedNotEmpty().
						HasCreatedOnNotEmpty(),
				),
			},
		},
	})
}

func TestAcc_LegacyServiceUser_WIF_ExternalChange(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	subject := fmt.Sprintf("system:serviceaccount:namespace:%s", random.AlphaN(10))
	gcpSubject := random.NumericN(10)

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
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityOidc(testvars.OidcIssuer, subject, testvars.OidcIssuer+"/o/oauth2/auth"),
					resourceshowoutputassert.UserShowOutput(t, userModelWithOIDC.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
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
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityOidc(testvars.OidcIssuer, subject, testvars.OidcIssuer+"/o/oauth2/auth"),
					resourceshowoutputassert.UserShowOutput(t, userModelWithOIDC.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
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
			// External change - WIF changed to GCP type externally
			{
				PreConfig: func() {
					testClient().User.SetGcpWorkloadIdentity(t, id, gcpSubject)
				},
				Config: config.FromModels(t, userModelWithOIDC),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				// Terraform should detect the type change and re-apply OIDC config
				Check: assertThat(t,
					resourceassert.LegacyServiceUserResource(t, userModelWithOIDC.ResourceReference()).
						HasNameString(id.Name()).
						HasDefaultWorkloadIdentityOidc(testvars.OidcIssuer, subject, testvars.OidcIssuer+"/o/oauth2/auth"),
					resourceshowoutputassert.UserShowOutput(t, userModelWithOIDC.ResourceReference()).
						HasHasWorkloadIdentity(true),
					objectassert.UserDefaultWorkloadIdentityAuthenticationMethods(t, id).
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
		},
	})
}

func TestAcc_LegacyServiceUser_WIF_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	userModelWithAwsEmpty := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityAwsEmpty()
	userModelWithGcpEmpty := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityGcpEmpty()
	userModelWithAzureEmpty := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityAzureEmpty()
	userModelWithOidcEmpty := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityOidcEmpty()
	userModelWithMultipleProviders := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityMultipleProviders()
	userModelWithEmptyBlock := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityEmpty()
	userModelWithoutEnabledFlag := model.LegacyServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityAzure(
			testvars.MicrosoftIssuer,
			"subject",
		)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, userModelWithAwsEmpty),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The argument "arn" is required, but no definition was found.`),
			},
			{
				Config:      config.FromModels(t, userModelWithGcpEmpty),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The argument "subject" is required, but no definition was found.`),
			},
			{
				Config:      config.FromModels(t, userModelWithAzureEmpty),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The argument "issuer" is required, but no definition was found.`),
			},
			{
				Config:      config.FromModels(t, userModelWithOidcEmpty),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The argument "issuer" is required, but no definition was found.`),
			},
			{
				Config:      config.FromModels(t, userModelWithMultipleProviders),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("`default_workload_identity.0.aws,default_workload_identity.0.gcp` were\nspecified"),
			},
			{
				Config:      config.FromModels(t, userModelWithEmptyBlock),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("one of\n`default_workload_identity.0.aws,default_workload_identity.0.azure,default_workload_identity.0.gcp,default_workload_identity.0.oidc`\nmust be specified"),
			},
			{
				Config: config.FromModels(t, userModelWithoutEnabledFlag),
				// PlanOnly is not set because the validation happens during resource operations.
				ExpectError: regexp.MustCompile("to use `default_workload_identity`, you need to first specify the `USER_ENABLE_DEFAULT_WORKLOAD_IDENTITY` feature in the `experimental_features_enabled` field at the provider level"),
			},
		},
	})
}
