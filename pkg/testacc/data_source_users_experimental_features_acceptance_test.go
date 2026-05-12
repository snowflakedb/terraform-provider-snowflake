//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Experimental_Users_ParametersReducedOutput(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("w", userId.Name())
	usersModel := datasourcemodel.Users("test").
		WithLike(userId.Name()).
		WithDependsOn(userModel.ResourceReference())
	usersModelWithoutParameters := datasourcemodel.Users("test").
		WithLike(userId.Name()).
		WithWithParameters(false).
		WithDependsOn(userModel.ResourceReference())
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
				Config:                   config.FromModels(t, userModel, usersModel, providerModel),
				Check: assertThat(t,
					resourceparametersassert.UsersDatasourceParameters(t, usersModel.DatasourceReference()).
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
				ProtoV6ProviderFactories: providerFactoryUsingCache("TestAcc_Experimental_Users_ParametersReducedOutput"),
				Config:                   config.FromModels(t, userModel, usersModel, providerModelWithExperimentEnabled),
				Check: assertThat(t,
					resourceparametersassert.UsersDatasourceParameters(t, usersModel.DatasourceReference()).
						HasEnableUnredactedQuerySyntaxErrorValueDefault().
						HasEnableUnredactedQuerySyntaxErrorLevel(sdk.ParameterTypeSnowflakeDefault).
						HasEnableUnredactedQuerySyntaxErrorKeyEmpty().
						HasEnableUnredactedQuerySyntaxErrorDefaultEmpty().
						HasEnableUnredactedQuerySyntaxErrorDescriptionEmpty(),
				),
			},
			// output without parameters
			{
				ProtoV6ProviderFactories: providerFactoryUsingCache("TestAcc_Experimental_Users_ParametersReducedOutput"),
				Config:                   config.FromModels(t, userModel, usersModelWithoutParameters, providerModelWithExperimentEnabled),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(usersModelWithoutParameters.DatasourceReference(), "users.0.parameters.#", "0")),
				),
			},
			// output with parameters
			{
				ProtoV6ProviderFactories: providerFactoryUsingCache("TestAcc_Experimental_Users_ParametersReducedOutput"),
				Config:                   config.FromModels(t, userModel, usersModel, providerModelWithExperimentEnabled),
				Check: assertThat(t,
					resourceparametersassert.UsersDatasourceParameters(t, usersModel.DatasourceReference()).
						HasEnableUnredactedQuerySyntaxErrorValueDefault().
						HasEnableUnredactedQuerySyntaxErrorLevel(sdk.ParameterTypeSnowflakeDefault).
						HasEnableUnredactedQuerySyntaxErrorKeyEmpty().
						HasEnableUnredactedQuerySyntaxErrorDefaultEmpty().
						HasEnableUnredactedQuerySyntaxErrorDescriptionEmpty(),
				),
			},
			// turning the reduced output off
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				Config:                   config.FromModels(t, userModel, usersModel, providerModel),
				Check: assertThat(t,
					resourceparametersassert.UsersDatasourceParameters(t, usersModel.DatasourceReference()).
						HasEnableUnredactedQuerySyntaxErrorValueDefault().
						HasEnableUnredactedQuerySyntaxErrorLevel(sdk.ParameterTypeSnowflakeDefault).
						HasEnableUnredactedQuerySyntaxErrorKey().
						HasEnableUnredactedQuerySyntaxErrorDefault().
						HasEnableUnredactedQuerySyntaxErrorDescriptionNotEmpty(),
				),
			},
		},
	})
}
