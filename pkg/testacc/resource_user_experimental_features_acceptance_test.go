//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Experimental_User_ParametersReducedOutput_UpdateExisting(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("w", userId.Name())
	userModelWithParameterSet := model.User("w", userId.Name()).
		WithEnableUnredactedQuerySyntaxError(true)
	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
	providerModelWithExperimentEnabled := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary).
		WithExperimentalFeaturesEnabled(experimentalfeatures.ParametersReducedOutput)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		// TODO [SNOW-1653619]: check destroy for secondary account
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// setting up initially with the whole output
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				Config:                   config.FromModels(t, providerModel, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasNameString(userId.Name()),
					resourceparametersassert.UserResourceParameters(t, userModel.ResourceReference()).
						// we just check a single parameter to verify logic
						HasEnableUnredactedQuerySyntaxErrorValueDefault().
						HasEnableUnredactedQuerySyntaxErrorLevel(sdk.ParameterTypeSnowflakeDefault).
						HasEnableUnredactedQuerySyntaxErrorKey().
						HasEnableUnredactedQuerySyntaxErrorDefault().
						HasEnableUnredactedQuerySyntaxErrorDescriptionNotEmpty(),
				),
			},
			// turning the reduced output on
			{
				ProtoV6ProviderFactories: providerFactoryUsingCache("TestAcc_Experimental_User_ParametersReducedOutput"),
				Config:                   config.FromModels(t, providerModelWithExperimentEnabled, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasNameString(userId.Name()),
					resourceparametersassert.UserResourceParameters(t, userModel.ResourceReference()).
						HasEnableUnredactedQuerySyntaxErrorValueDefault().
						HasEnableUnredactedQuerySyntaxErrorLevel(sdk.ParameterTypeSnowflakeDefault).
						HasEnableUnredactedQuerySyntaxErrorKeyEmpty().
						HasEnableUnredactedQuerySyntaxErrorDefaultEmpty().
						HasEnableUnredactedQuerySyntaxErrorDescriptionEmpty(),
				),
			},
			// changing the value in config
			{
				ProtoV6ProviderFactories: providerFactoryUsingCache("TestAcc_Experimental_User_ParametersReducedOutput"),
				Config:                   config.FromModels(t, providerModelWithExperimentEnabled, userModelWithParameterSet),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasNameString(userId.Name()),
					resourceparametersassert.UserResourceParameters(t, userModel.ResourceReference()).
						HasEnableUnredactedQuerySyntaxError(true).
						HasEnableUnredactedQuerySyntaxErrorLevel(sdk.ParameterTypeUser).
						HasEnableUnredactedQuerySyntaxErrorKeyEmpty().
						HasEnableUnredactedQuerySyntaxErrorDefaultEmpty().
						HasEnableUnredactedQuerySyntaxErrorDescriptionEmpty(),
				),
			},
			// changing the value externally
			{
				ProtoV6ProviderFactories: providerFactoryUsingCache("TestAcc_Experimental_User_ParametersReducedOutput"),
				PreConfig: func() {
					secondaryTestClient().User.UpdateEnableUnredactedQuerySyntaxError(t, userId, false)
				},
				Config: config.FromModels(t, providerModelWithExperimentEnabled, userModelWithParameterSet),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasNameString(userId.Name()),
					resourceparametersassert.UserResourceParameters(t, userModel.ResourceReference()).
						HasEnableUnredactedQuerySyntaxError(true).
						HasEnableUnredactedQuerySyntaxErrorLevel(sdk.ParameterTypeUser).
						HasEnableUnredactedQuerySyntaxErrorKeyEmpty().
						HasEnableUnredactedQuerySyntaxErrorDefaultEmpty().
						HasEnableUnredactedQuerySyntaxErrorDescriptionEmpty(),
				),
			},
			// turning the reduced output off
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				Config:                   config.FromModels(t, providerModel, userModelWithParameterSet),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasNameString(userId.Name()),
					resourceparametersassert.UserResourceParameters(t, userModel.ResourceReference()).
						HasEnableUnredactedQuerySyntaxError(true).
						HasEnableUnredactedQuerySyntaxErrorLevel(sdk.ParameterTypeUser).
						HasEnableUnredactedQuerySyntaxErrorKey().
						HasEnableUnredactedQuerySyntaxErrorDefault().
						HasEnableUnredactedQuerySyntaxErrorDescriptionNotEmpty(),
				),
			},
		},
	})
}

func TestAcc_Experimental_User_ParametersReducedOutput_CreateNew(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("w", userId.Name())
	providerModelWithExperimentEnabled := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary).
		WithExperimentalFeaturesEnabled(experimentalfeatures.ParametersReducedOutput)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactoryUsingCache("TestAcc_Experimental_User_ParametersReducedOutput"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		// TODO [SNOW-1653619]: check destroy for secondary account
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModelWithExperimentEnabled, userModel),
				Check: assertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasNameString(userId.Name()),
					resourceparametersassert.UserResourceParameters(t, userModel.ResourceReference()).
						HasEnableUnredactedQuerySyntaxErrorValueDefault().
						HasEnableUnredactedQuerySyntaxErrorLevel(sdk.ParameterTypeSnowflakeDefault).
						HasEnableUnredactedQuerySyntaxErrorKeyEmpty().
						HasEnableUnredactedQuerySyntaxErrorDefaultEmpty().
						HasEnableUnredactedQuerySyntaxErrorDescriptionEmpty(),
				),
			},
		},
	})
}
