package manual

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// This is a manual test for authenticating with Okta.
func TestAcc_Provider_OktaAuth(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableManual)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.User)
			testenvs.AssertEnvNotSet(t, snowflakeenvs.Password)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfigWithAuthenticator(t, Okta, sdk.AuthenticationTypeOkta),
			},
		},
	})
}

// This test requires manual action due to MFA. Make sure the user does not have a positive `mins_to_bypass_mfa` in `SHOW USERS`.
func TestAcc_Provider_UsernamePasswordMfaAuth(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableManual)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// ensure MFA is checked here - accept login on your MFA device
			{
				Config: providerConfigWithAuthenticator(t, testprofiles.Default, sdk.AuthenticationTypeUsernamePasswordMfa),
			},
			// check that MFA login is cached - this step should not require manual action
			{
				Config: providerConfigWithAuthenticator(t, testprofiles.Default, sdk.AuthenticationTypeUsernamePasswordMfa),
			},
		},
	})
}

// This test requires manual action due to MFA. Make sure the user does not have a positive `mins_to_bypass_mfa` in `SHOW USERS`.
func TestAcc_Provider_UsernamePasswordMfaAuthWithPasscode(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableManual)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// ensure MFA is checked here - accept access to keychain on your device
			{
				Config: providerConfigWithAuthenticator(t, DefaultWithPasscode, sdk.AuthenticationTypeUsernamePasswordMfa),
			},
			// check that MFA login is cached - this step should not require manual action
			{
				Config: providerConfigWithAuthenticator(t, DefaultWithPasscode, sdk.AuthenticationTypeUsernamePasswordMfa),
			},
		},
	})
}

func providerConfigWithAuthenticator(t *testing.T, profile string, authenticator sdk.AuthenticationType) string {
	t.Helper()
	return config.FromModels(t,
		providermodel.SnowflakeProvider().WithProfile(profile).WithAuthenticator(string(authenticator)),
		model.Database("t", "SNOWFLAKE"),
	)
}
