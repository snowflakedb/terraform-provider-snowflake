package windows_test

import (
	"fmt"
	"io/fs"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testfiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/manual_tests"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Provider_tomlConfigIsTooPermissive(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableManual)
	if !oswrapper.IsRunningOnWindows() {
		t.Skip("checking file permissions on other platforms is currently done in the provider package")
	}
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	permissions := fs.FileMode(0o755)

	configPath := testfiles.TestFileWithCustomPermissions(t, random.AlphaN(10), random.Bytes(), permissions)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: manual_tests.ManualTestProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, configPath)
				},
				Config:      config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(configPath), datasourceModel()),
				ExpectError: regexp.MustCompile(fmt.Sprintf("could not load config file: config file %s has unsafe permissions - %#o", configPath, permissions)),
			},
		},
	})
}

func datasourceModel() config.DatasourceModel {
	return datasourcemodel.Database("t", "SNOWFLAKE")
}
