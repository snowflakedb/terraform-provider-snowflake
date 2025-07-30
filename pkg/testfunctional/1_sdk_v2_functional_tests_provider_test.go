package testfunctional_test

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SdkV2FunctionalTestsProviderName = "snowflake-sdk-v2-functional-tests"

// sdkV2FunctionalTestsProvider returns a Terraform Provider used for our SDKv2 functional tests.
func sdkV2FunctionalTestsProvider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"test_name": {
				Type:        schema.TypeString,
				Description: "Specifies the name of the test used to instantiate the provider.",
				Optional:    true,
			},
		},
		ResourcesMap:         testResources(),
		DataSourcesMap:       testDataSources(),
		ConfigureContextFunc: configureTestProvider,
		ProviderMetaSchema:   map[string]*schema.Schema{},
	}
}

func testResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"snowflake_test_resource_object_renaming":              testfunctional.TestResourceObjectRenamingListsAndSets(),
		"snowflake_test_resource_data_type_diff_handling":      testfunctional.TestResourceDataTypeDiffHandling(),
		"snowflake_test_resource_data_type_diff_handling_list": testfunctional.TestResourceDataTypeDiffHandlingList(),
	}
}

func testDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

func configureTestProvider(_ context.Context, s *schema.ResourceData) (any, diag.Diagnostics) {
	var testName string
	if v, ok := s.GetOk("test_name"); ok {
		testName = v.(string)
	}

	return common.NewTestProviderContext(testName, server.URL), nil
}
