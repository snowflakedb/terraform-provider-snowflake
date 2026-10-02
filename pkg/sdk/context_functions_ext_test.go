package sdk

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	roleId := randomAccountObjectIdentifier()

	contextFunctionsTests.CurrentAccount.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentAccount_basic,
			"SELECT CURRENT_ACCOUNT() as CURRENT_ACCOUNT",
		)
	contextFunctionsTests.CurrentOrganizationName.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentOrganizationName_basic,
			"SELECT CURRENT_ORGANIZATION_NAME() as CURRENT_ORGANIZATION_NAME",
		)
	contextFunctionsTests.CurrentAccountName.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentAccountName_basic,
			"SELECT CURRENT_ACCOUNT_NAME() as CURRENT_ACCOUNT_NAME",
		)
	contextFunctionsTests.CurrentRole.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentRole_basic,
			"SELECT CURRENT_ROLE() as CURRENT_ROLE",
		)
	contextFunctionsTests.CurrentSecondaryRoles.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentSecondaryRoles_basic,
			"SELECT CURRENT_SECONDARY_ROLES() as CURRENT_ROLES",
		)
	contextFunctionsTests.CurrentRegion.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentRegion_basic,
			"SELECT CURRENT_REGION() AS CURRENT_REGION",
		)
	contextFunctionsTests.CurrentSession.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentSession_basic,
			"SELECT CURRENT_SESSION() as CURRENT_SESSION",
		)
	contextFunctionsTests.CurrentUser.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentUser_basic,
			"SELECT CURRENT_USER() as CURRENT_USER",
		)
	contextFunctionsTests.CurrentHost.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentHost_basic,
			"SELECT VALUE:host::VARCHAR AS HOST FROM TABLE(FLATTEN(INPUT => PARSE_JSON(SYSTEM$ALLOWLIST()))) WHERE VALUE:type::VARCHAR = 'SNOWFLAKE_DEPLOYMENT'",
		)
	contextFunctionsTests.CurrentSessionDetails.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentSessionDetails_basic,
			"SELECT CURRENT_ACCOUNT() as CURRENT_ACCOUNT, CURRENT_ROLE() as CURRENT_ROLE, CURRENT_REGION() AS CURRENT_REGION, CURRENT_SESSION() as CURRENT_SESSION, CURRENT_USER() as CURRENT_USER, CURRENT_ACCOUNT_NAME() as CURRENT_ACCOUNT_NAME, CURRENT_ORGANIZATION_NAME() as CURRENT_ORGANIZATION_NAME",
		)
	contextFunctionsTests.LastQueryId.
		withExpectedSql(
			case_ContextFunctions_sql_LastQueryId_basic,
			"SELECT LAST_QUERY_ID() as LAST_QUERY_ID",
		)
	contextFunctionsTests.CurrentDatabase.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentDatabase_basic,
			"SELECT CURRENT_DATABASE() as CURRENT_DATABASE",
		)
	contextFunctionsTests.CurrentSchema.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentSchema_basic,
			"SELECT CURRENT_SCHEMA() as CURRENT_SCHEMA",
		)
	contextFunctionsTests.CurrentWarehouse.
		withExpectedSql(
			case_ContextFunctions_sql_CurrentWarehouse_basic,
			"SELECT CURRENT_WAREHOUSE() as CURRENT_WAREHOUSE",
		)
	contextFunctionsTests.IsRoleInSession.
		withDefaultOpts(func() *IsRoleInSessionOptions {
			return &IsRoleInSessionOptions{
				Arguments: IsRoleInSessionArguments{Role: roleId},
			}
		}).
		withExpectedSqlf(
			case_ContextFunctions_sql_IsRoleInSession_basic,
			"SELECT IS_ROLE_IN_SESSION ('%s') AS IS_ROLE_IN_SESSION",
			strings.ReplaceAll(roleId.FullyQualifiedName(), `"`, `\"`),
		)
}

func TestCurrentSessionDetails_AccountURL(t *testing.T) {
	testCases := []struct {
		region      string
		expectedURL string
	}{
		// AWS regions
		{region: "aws_us_west_2", expectedURL: "https://TESTACCOUNT.snowflakecomputing.com"},
		{region: "aws_us_gov_west_2", expectedURL: "https://TESTACCOUNT.us-gov-west-2.aws.snowflakecomputing.com"},
		{region: "aws_us_gov_west_1_fhplus", expectedURL: "https://TESTACCOUNT.fhplus.us-gov-west-1.aws.snowflakecomputing.com"},
		{region: "aws_us_gov_west_1_dod", expectedURL: "https://TESTACCOUNT.dod.us-gov-west-1.aws.snowflakecomputing.com"},
		{region: "aws_us_east_2", expectedURL: "https://TESTACCOUNT.us-east-2.aws.snowflakecomputing.com"},
		{region: "aws_us_east_1", expectedURL: "https://TESTACCOUNT.us-east-1.snowflakecomputing.com"},
		{region: "aws_us_gov_east_1", expectedURL: "https://TESTACCOUNT.us-gov-east-1.aws.snowflakecomputing.com"},
		{region: "aws_us_gov_east_1_fhplus", expectedURL: "https://TESTACCOUNT.fhplus.us-gov-east-1.aws.snowflakecomputing.com"},
		{region: "aws_ca_central_1", expectedURL: "https://TESTACCOUNT.ca-central-1.aws.snowflakecomputing.com"},
		{region: "aws_sa_east_1", expectedURL: "https://TESTACCOUNT.sa-east-1.aws.snowflakecomputing.com"},
		{region: "aws_af_south_1", expectedURL: "https://TESTACCOUNT.af-south-1.aws.snowflakecomputing.com"},
		{region: "aws_eu_west_1", expectedURL: "https://TESTACCOUNT.eu-west-1.snowflakecomputing.com"},
		{region: "aws_eu_west_2", expectedURL: "https://TESTACCOUNT.eu-west-2.aws.snowflakecomputing.com"},
		{region: "aws_eu_west_3", expectedURL: "https://TESTACCOUNT.eu-west-3.aws.snowflakecomputing.com"},
		{region: "aws_eu_central_1", expectedURL: "https://TESTACCOUNT.eu-central-1.snowflakecomputing.com"},
		{region: "aws_eu_central_2", expectedURL: "https://TESTACCOUNT.eu-central-2.aws.snowflakecomputing.com"},
		{region: "aws_eu_north_1", expectedURL: "https://TESTACCOUNT.eu-north-1.aws.snowflakecomputing.com"},
		{region: "aws_ap_northeast_1", expectedURL: "https://TESTACCOUNT.ap-northeast-1.aws.snowflakecomputing.com"},
		{region: "aws_ap_northeast_2", expectedURL: "https://TESTACCOUNT.ap-northeast-2.aws.snowflakecomputing.com"},
		{region: "aws_ap_northeast_3", expectedURL: "https://TESTACCOUNT.ap-northeast-3.aws.snowflakecomputing.com"},
		{region: "aws_ap_south_1", expectedURL: "https://TESTACCOUNT.ap-south-1.aws.snowflakecomputing.com"},
		{region: "aws_ap_southeast_1", expectedURL: "https://TESTACCOUNT.ap-southeast-1.snowflakecomputing.com"},
		{region: "aws_ap_southeast_2", expectedURL: "https://TESTACCOUNT.ap-southeast-2.snowflakecomputing.com"},
		{region: "aws_ap_southeast_3", expectedURL: "https://TESTACCOUNT.ap-southeast-3.aws.snowflakecomputing.com"},
		{region: "aws_cn_northwest_1", expectedURL: "https://TESTACCOUNT.cn-northwest-1.aws.snowflakecomputing.cn"},

		// GCP regions
		{region: "gcp_us_central1", expectedURL: "https://TESTACCOUNT.us-central1.gcp.snowflakecomputing.com"},
		{region: "gcp_us_east4", expectedURL: "https://TESTACCOUNT.us-east4.gcp.snowflakecomputing.com"},
		{region: "gcp_europe_west2", expectedURL: "https://TESTACCOUNT.europe-west2.gcp.snowflakecomputing.com"},
		{region: "gcp_europe_west3", expectedURL: "https://TESTACCOUNT.europe-west3.gcp.snowflakecomputing.com"},
		{region: "gcp_europe_west4", expectedURL: "https://TESTACCOUNT.europe-west4.gcp.snowflakecomputing.com"},
		{region: "gcp_me_central2", expectedURL: "https://TESTACCOUNT.me-central2.gcp.snowflakecomputing.com"},

		// Azure regions
		{region: "azure_westus2", expectedURL: "https://TESTACCOUNT.west-us-2.azure.snowflakecomputing.com"},
		{region: "azure_centralus", expectedURL: "https://TESTACCOUNT.central-us.azure.snowflakecomputing.com"},
		{region: "azure_southcentralus", expectedURL: "https://TESTACCOUNT.south-central-us.azure.snowflakecomputing.com"},
		{region: "azure_eastus2", expectedURL: "https://TESTACCOUNT.east-us-2.azure.snowflakecomputing.com"},
		{region: "azure_usgovvirginia", expectedURL: "https://TESTACCOUNT.us-gov-virginia.azure.snowflakecomputing.com"},
		{region: "azure_usgovvirginia_fhplus", expectedURL: "https://TESTACCOUNT.fhplus.us-gov-virginia.azure.snowflakecomputing.com"},
		{region: "azure_canadacentral", expectedURL: "https://TESTACCOUNT.canada-central.azure.snowflakecomputing.com"},
		{region: "azure_mexicocentral", expectedURL: "https://TESTACCOUNT.mexicocentral.azure.snowflakecomputing.com"},
		{region: "azure_uksouth", expectedURL: "https://TESTACCOUNT.uk-south.azure.snowflakecomputing.com"},
		{region: "azure_northeurope", expectedURL: "https://TESTACCOUNT.north-europe.azure.snowflakecomputing.com"},
		{region: "azure_swedencentral", expectedURL: "https://TESTACCOUNT.sweden-central.azure.snowflakecomputing.com"},
		{region: "azure_westeurope", expectedURL: "https://TESTACCOUNT.west-europe.azure.snowflakecomputing.com"},
		{region: "azure_southeastasia", expectedURL: "https://TESTACCOUNT.southeast-asia.azure.snowflakecomputing.com"},
		{region: "azure_switzerlandnorth", expectedURL: "https://TESTACCOUNT.switzerland-north.azure.snowflakecomputing.com"},
		{region: "azure_uaenorth", expectedURL: "https://TESTACCOUNT.uae-north.azure.snowflakecomputing.com"},
		{region: "azure_centralindia", expectedURL: "https://TESTACCOUNT.central-india.azure.snowflakecomputing.com"},
		{region: "azure_japaneast", expectedURL: "https://TESTACCOUNT.japan-east.azure.snowflakecomputing.com"},
		{region: "azure_koreacentral", expectedURL: "https://TESTACCOUNT.korea-central.azure.snowflakecomputing.com"},
		{region: "azure_australiaeast", expectedURL: "https://TESTACCOUNT.australia-east.azure.snowflakecomputing.com"},

		{region: "AWS_US_WEST_2", expectedURL: "https://TESTACCOUNT.snowflakecomputing.com"}, // case insensitive

		// prefixed forms returned by CURRENT_REGION() in newer Snowflake versions
		{region: "PUBLIC.AWS_US_EAST_1", expectedURL: "https://TESTACCOUNT.us-east-1.snowflakecomputing.com"},
		{region: "AWS_US_GOV_EAST_1_FHPLUS.AWS_US_GOV_EAST_1_FHPLUS", expectedURL: "https://TESTACCOUNT.fhplus.us-gov-east-1.aws.snowflakecomputing.com"},
		{region: "PUBLIC.AZURE_EASTUS2", expectedURL: "https://TESTACCOUNT.east-us-2.azure.snowflakecomputing.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.region, func(t *testing.T) {
			details := &CurrentSessionDetails{
				Account: "TESTACCOUNT",
				Region:  tc.region,
			}
			url, err := details.AccountURL()
			require.NoError(t, err)
			assert.Equal(t, tc.expectedURL, url)
		})
	}

	invalidTestCases := []struct {
		region      string
		expectedErr string
	}{
		{region: "invalid_region", expectedErr: "failed to map Snowflake account region invalid_region to a region_id"},
		{region: "", expectedErr: "failed to map Snowflake account region  to a region_id"},
	}

	for _, tc := range invalidTestCases {
		t.Run(tc.region, func(t *testing.T) {
			details := &CurrentSessionDetails{
				Account: "TESTACCOUNT",
				Region:  tc.region,
			}
			_, err := details.AccountURL()
			require.ErrorContains(t, err, tc.expectedErr)
		})
	}
}
