//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ServiceUser_WIF_OIDC(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	userModelWithOIDC := model.ServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityOidc(
			"https://accounts.google.com",
			"system:serviceaccount:namespace:sa-name",
			[]string{"https://accounts.google.com/o/oauth2/auth"},
		)

	userModelWithoutWIF := model.ServiceUser("w", id.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ServiceUser),
		Steps: []resource.TestStep{
			// CREATE WITH OIDC WIF
			{
				Config: config.FromModels(t, userModelWithOIDC),
				Check: assertThat(t,
					resourceassert.ServiceUserResource(t, userModelWithOIDC.ResourceReference()).
						HasNameString(id.Name()),
					objectassert.User(t, id).
						HasHasWorkloadIdentity(true),
				),
			},
			// UPDATE - REMOVE WIF
			{
				Config: config.FromModels(t, userModelWithoutWIF),
				Check: assertThat(t,
					objectassert.User(t, id).
						HasHasWorkloadIdentity(false),
				),
			},
			// UPDATE - ADD WIF BACK
			{
				Config: config.FromModels(t, userModelWithOIDC),
				Check: assertThat(t,
					objectassert.User(t, id).
						HasHasWorkloadIdentity(true),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelWithOIDC.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "mins_to_unlock", "login_name", "display_name", "disabled", "default_secondary_roles_option"},
			},
		},
	})
}

func TestAcc_ServiceUser_WIF_AWS(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	userModelWithAWS := model.ServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityAws("arn:aws:iam::123456789012:role/test-role")

	userModelWithoutWIF := model.ServiceUser("w", id.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ServiceUser),
		Steps: []resource.TestStep{
			// CREATE WITH AWS WIF
			{
				Config: config.FromModels(t, userModelWithAWS),
				Check: assertThat(t,
					resourceassert.ServiceUserResource(t, userModelWithAWS.ResourceReference()).
						HasNameString(id.Name()),
					objectassert.User(t, id).
						HasHasWorkloadIdentity(true),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelWithAWS.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "mins_to_unlock", "login_name", "display_name", "disabled", "default_secondary_roles_option"},
			},
			// UPDATE - REMOVE WIF
			{
				Config: config.FromModels(t, userModelWithoutWIF),
				Check: assertThat(t,
					objectassert.User(t, id).
						HasHasWorkloadIdentity(false),
				),
			},
		},
	})
}

func TestAcc_ServiceUser_WIF_GCP(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	userModelWithGCP := model.ServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityGcp("projects/my-project/locations/global/workloadIdentityPools/my-pool/subject/my-subject")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ServiceUser),
		Steps: []resource.TestStep{
			// CREATE WITH GCP WIF
			{
				Config: config.FromModels(t, userModelWithGCP),
				Check: assertThat(t,
					objectassert.User(t, id).
						HasHasWorkloadIdentity(true),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelWithGCP.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "mins_to_unlock", "login_name", "display_name", "disabled", "default_secondary_roles_option"},
			},
		},
	})
}

func TestAcc_ServiceUser_WIF_Azure(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	userModelWithAzure := model.ServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityAzure(
			"https://login.microsoftonline.com/tenant-id/v2.0",
			"subject-identifier",
		)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ServiceUser),
		Steps: []resource.TestStep{
			// CREATE WITH Azure WIF
			{
				Config: config.FromModels(t, userModelWithAzure),
				Check: assertThat(t,
					objectassert.User(t, id).
						HasHasWorkloadIdentity(true),
				),
			},
			// IMPORT
			{
				ResourceName:            userModelWithAzure.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"days_to_expiry", "mins_to_unlock", "login_name", "display_name", "disabled", "default_secondary_roles_option"},
			},
		},
	})
}

func TestAcc_ServiceUser_WIF_SwitchProvider(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	userModelWithOIDC := model.ServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityOidc(
			"https://accounts.google.com",
			"system:serviceaccount:namespace:sa-name",
			[]string{"https://accounts.google.com/o/oauth2/auth"},
		)

	userModelWithAWS := model.ServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityAws("arn:aws:iam::123456789012:role/test-role")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ServiceUser),
		Steps: []resource.TestStep{
			// Start with OIDC
			{
				Config: config.FromModels(t, userModelWithOIDC),
				Check: assertThat(t,
					objectassert.User(t, id).
						HasHasWorkloadIdentity(true),
				),
			},
			// Switch to AWS
			{
				Config: config.FromModels(t, userModelWithAWS),
				Check: assertThat(t,
					objectassert.User(t, id).
						HasHasWorkloadIdentity(true),
				),
			},
		},
	})
}

func TestAcc_ServiceUser_WIF_ExternalChange(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	userModelWithOIDC := model.ServiceUser("w", id.Name()).
		WithDefaultWorkloadIdentityOidc(
			"https://accounts.google.com",
			"system:serviceaccount:namespace:sa-name",
			[]string{"https://accounts.google.com/o/oauth2/auth"},
		)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ServiceUser),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, userModelWithOIDC),
				Check: assertThat(t,
					objectassert.User(t, id).
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
					objectassert.User(t, id).
						HasHasWorkloadIdentity(true),
				),
			},
		},
	})
}
