//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AuthenticationPolicy(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := "This is a test resource"
	basicModel := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment(comment).
		WithAuthenticationMethods(sdk.AuthenticationMethodsPassword).
		WithMfaEnrollmentEnum(sdk.MfaEnrollmentRequired).
		WithClientTypes(sdk.ClientTypesSnowflakeUi).
		WithSecurityIntegrations("ALL")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, basicModel),
				Check: assertThat(t, resourceassert.AuthenticationPolicyResource(t, basicModel.ResourceReference()).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasCommentString(comment).
					HasAuthenticationMethods(sdk.AuthenticationMethodsPassword).
					HasMfaEnrollmentString(string(sdk.MfaEnrollmentRequired)).
					HasClientTypes(sdk.ClientTypesSnowflakeUi).
					HasSecurityIntegrations("ALL"),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, basicModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
				),
			},
			{
				Config:            accconfig.FromModels(t, basicModel),
				ResourceName:      basicModel.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_AuthenticationPolicy_complete(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithAuthenticationMethods(sdk.AuthenticationMethodsPassword).
		WithMfaEnrollmentEnum(sdk.MfaEnrollmentRequired).
		WithClientTypes(sdk.ClientTypesSnowflakeUi).
		WithSecurityIntegrations("ALL").
		WithComment(comment)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.AuthenticationPolicyResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString(comment).
						HasAuthenticationMethods(sdk.AuthenticationMethodsPassword).
						HasMfaEnrollmentString(string(sdk.MfaEnrollmentRequired)).
						HasClientTypes(sdk.ClientTypesSnowflakeUi).
						HasSecurityIntegrations("ALL"),
					resourceshowoutputassert.AuthenticationPolicyShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasKind(string(sdk.PolicyKindAuthenticationPolicy)).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasOptions(""),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.name.0.value", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner.0.value", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.comment.0.value", comment)),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.authentication_methods.0.value", "[PASSWORD]")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.client_types.0.value", "[SNOWFLAKE_UI]")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.security_integrations.0.value", "[ALL]")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.mfa_enrollment.0.value", "OPTIONAL")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.mfa_authentication_methods.0.value", "[PASSWORD]")),
				),
			},
			{
				Config:            accconfig.FromModels(t, modelComplete),
				ResourceName:      modelComplete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_AuthenticationPolicy_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelInvalidAuthenticationMethods := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithAuthenticationMethods("invalid")
	modelInvalidMfaEnrollment := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMfaEnrollmentEnum("invalid")
	modelInvalidClientTypes := model.AuthenticationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithClientTypes("invalid")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelInvalidAuthenticationMethods),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid authentication method: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidMfaEnrollment),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid MFA enrollment option: INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, modelInvalidClientTypes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid client type: INVALID`),
			},
		},
	})
}
