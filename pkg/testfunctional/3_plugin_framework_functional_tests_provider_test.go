package testfunctional_test

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/computednestedlist"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/httpserver"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const PluginFrameworkFunctionalTestsProviderName = "snowflake-plugin-framework-functional-tests"

// Ensure the implementation satisfies the provider.Provider interface.
var _ provider.Provider = &pluginFrameworkFunctionalTestsProvider{}

type pluginFrameworkFunctionalTestsProvider struct{}

type pluginFrameworkFunctionalTestsProviderModelV0 struct {
	TestName types.String `tfsdk:"test_name"`
}

func (p *pluginFrameworkFunctionalTestsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = PluginFrameworkFunctionalTestsProviderName
	response.Version = "dev"
}

func (p *pluginFrameworkFunctionalTestsProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"test_name": schema.StringAttribute{
				Description: "Specifies the name of the test used to instantiate the provider.",
				Optional:    true,
			},
		},
	}
}

func (p *pluginFrameworkFunctionalTestsProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var configModel pluginFrameworkFunctionalTestsProviderModelV0

	response.Diagnostics.Append(request.Config.Get(ctx, &configModel)...)

	var testName string
	if !configModel.TestName.IsNull() {
		testName = configModel.TestName.ValueString()
	}

	if response.Diagnostics.HasError() {
		return
	}

	providerCtx := common.NewTestProviderContext(testName, server.URL)
	response.DataSourceData = providerCtx
	response.ResourceData = providerCtx
}

func (p *pluginFrameworkFunctionalTestsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *pluginFrameworkFunctionalTestsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		testfunctional.NewSomeResource,
		testfunctional.NewZeroValuesResource,
		computednestedlist.NewComputedNestedListResource,
		httpserver.NewHttpServerResource,
		testfunctional.NewStringWithMetadataResource,
		testfunctional.NewOptionalWithBackingFieldResource,
		testfunctional.NewParameterHandlingResourcePlanModifierResource,
		testfunctional.NewParameterHandlingReadLogicResource,
		testfunctional.NewParameterHandlingBackingFieldResource,
		testfunctional.NewParameterHandlingPrivateResource,
		testfunctional.NewEnumHandlingResource,
	}
}

func New() provider.Provider {
	return &pluginFrameworkFunctionalTestsProvider{}
}
