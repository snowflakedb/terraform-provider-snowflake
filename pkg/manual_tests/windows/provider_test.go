package windows_test

import (
	"fmt"
	"io/fs"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/snowflakedb/terraform-provider-snowflake/v2/pkg/testhelpers"
)

func TestAcc_Provider_tomlConfigIsTooPermissive(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	if !oswrapper.IsRunningOnWindows() {
		t.Skip("checking file permissions on other platforms is currently done in the provider package")
	}
	acc.TestAccPreCheck(t)
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	permissions := fs.FileMode(0o755)

	configPath := testhelpers.TestFileWithCustomPermissions(t, random.AlphaN(10), random.Bytes(), permissions)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
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
	return datasourcemodel.Database("t", acc.TestDatabaseName)
}
