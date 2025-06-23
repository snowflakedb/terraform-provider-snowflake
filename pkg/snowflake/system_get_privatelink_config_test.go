package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemGetPrivateLinkConfigQuery(t *testing.T) {
	r := require.New(t)
	sb := SystemGetPrivateLinkConfigQuery()

	r.Equal(`SELECT SYSTEM$GET_PRIVATELINK_CONFIG() AS "config"`, sb)
}

func TestSystemGetPrivateLinkGetStructuredConfigAws(t *testing.T) {
	r := require.New(t)

	raw := &RawPrivateLinkConfig{
		Config: `{"privatelink-account-name":"ab1234.eu-west-1.privatelink","privatelink-vpce-id":"com.amazonaws.vpce.eu-west-1.vpce-svc-123456789abcdef12","privatelink-account-url":"ab1234.eu-west-1.privatelink.snowflakecomputing.com","privatelink_ocsp-url":"ocsp.ab1234.eu-west-1.privatelink.snowflakecomputing.com","regionless-snowsight-privatelink-url":"app-snoworg-snow_account.privatelink.snowflakecomputing.com","privatelink-internal-stage":"sfc-xx-xx1-12-customer-stage.s3.eu-west-1.amazonaws.com","snowsight-privatelink-url":"app.eu-west-1.privatelink.snowflakecomputing.com","regionless-privatelink-account-url":"snoworg-snow_account.privatelink.snowflakecomputing.com"}`,
	}

	c, e := raw.GetStructuredConfig()
	r.NoError(e)

	r.Equal("ab1234.eu-west-1.privatelink", c.AccountName)
	r.Equal("com.amazonaws.vpce.eu-west-1.vpce-svc-123456789abcdef12", c.AwsVpceID)
	r.Equal("", c.AzurePrivateLinkServiceID)
	r.Equal("ab1234.eu-west-1.privatelink.snowflakecomputing.com", c.AccountURL)
	r.Equal("ocsp.ab1234.eu-west-1.privatelink.snowflakecomputing.com", c.OCSPURL)
	r.Equal("snoworg-snow_account.privatelink.snowflakecomputing.com", c.RegionlessAccountURL)
	r.Equal("app.eu-west-1.privatelink.snowflakecomputing.com", c.SnowsightURL)
	r.Equal("sfc-xx-xx1-12-customer-stage.s3.eu-west-1.amazonaws.com", c.InternalStage)
	r.Equal("app-snoworg-snow_account.privatelink.snowflakecomputing.com", c.RegionlessSnowsightURL)
}

func TestSystemGetPrivateLinkGetStructuredConfigAwsAsPerDocumentation(t *testing.T) {
	r := require.New(t)

	raw := &RawPrivateLinkConfig{
		Config: `{"privatelink-account-name":"ab1234.eu-west-1.privatelink","privatelink-vpce-id":"com.amazonaws.vpce.eu-west-1.vpce-svc-123456789abcdef12","privatelink-account-url":"ab1234.eu-west-1.privatelink.snowflakecomputing.com","privatelink-ocsp-url":"ocsp.ab1234.eu-west-1.privatelink.snowflakecomputing.com"}`,
	}

	c, e := raw.GetStructuredConfig()
	r.NoError(e)

	r.Equal("ab1234.eu-west-1.privatelink", c.AccountName)
	r.Equal("com.amazonaws.vpce.eu-west-1.vpce-svc-123456789abcdef12", c.AwsVpceID)
	r.Equal("", c.AzurePrivateLinkServiceID)
	r.Equal("ab1234.eu-west-1.privatelink.snowflakecomputing.com", c.AccountURL)
	r.Equal("ocsp.ab1234.eu-west-1.privatelink.snowflakecomputing.com", c.OCSPURL)
}

func TestSystemGetPrivateLinkGetStructuredConfigAzure(t *testing.T) {
	r := require.New(t)

	raw := &RawPrivateLinkConfig{
		Config: `{"privatelink-account-name":"ab1234.east-us-2.azure.privatelink","privatelink-pls-id":"sf-pvlinksvc-azeastus2.east-us-2.azure.something","privatelink-account-url":"ab1234.east-us-2.azure.privatelink.snowflakecomputing.com","privatelink_ocsp-url":"ocsp.ab1234.east-us-2.azure.privatelink.snowflakecomputing.com"}`,
	}

	c, e := raw.GetStructuredConfig()
	r.NoError(e)

	r.Equal("ab1234.east-us-2.azure.privatelink", c.AccountName)
	r.Equal("", c.AwsVpceID)
	r.Equal("sf-pvlinksvc-azeastus2.east-us-2.azure.something", c.AzurePrivateLinkServiceID)
	r.Equal("ab1234.east-us-2.azure.privatelink.snowflakecomputing.com", c.AccountURL)
	r.Equal("ocsp.ab1234.east-us-2.azure.privatelink.snowflakecomputing.com", c.OCSPURL)
}
