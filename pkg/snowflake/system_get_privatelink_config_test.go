package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemGetPrivateLinkConfigQuery(t *testing.T) {
	r := require.New(t)
	sb := SystemGetPrivateLinkConfigQuery()

	r.Equal(`SELECT SYSTEM$GET_PRIVATELINK_CONFIG() AS "CONFIG"`, sb)
}

// TestSystemGetPrivateLinkGetStructuredConfigAws mirrors an AWS response where Snowflake
// uses the key "privatelink_ocsp-url" (underscore before "ocsp") instead of the documented "privatelink-ocsp-url".
func TestSystemGetPrivateLinkGetStructuredConfigAws(t *testing.T) {
	r := require.New(t)

	raw := &RawPrivateLinkConfig{
		Config: `{
			"privatelink-account-name": "testaccount.us-east-1.privatelink",
			"privatelink-vpce-id": "com.amazonaws.vpce.us-east-1.vpce-svc-00000000000000000",
			"privatelink-account-url": "testaccount.us-east-1.privatelink.snowflakecomputing.com",
			"privatelink_ocsp-url": "ocsp.testaccount.us-east-1.privatelink.snowflakecomputing.com",
			"privatelink-account-principal": "arn:aws:iam::000000000000:root",
			"app-service-privatelink-url": "*.testaccount.us-east-1.privatelink.snowflakecomputing.com",
			"regionless-snowsight-privatelink-url": "app-testorg-testaccount.privatelink.snowflakecomputing.com",
			"snowsight-privatelink-url": "app.us-east-1.privatelink.snowflakecomputing.com",
			"regionless-privatelink-account-url": "testorg-testaccount.privatelink.snowflakecomputing.com",
			"regionless-privatelink-ocsp-url": "ocsp.testorg-testaccount.privatelink.snowflakecomputing.com",
			"privatelink-dashed-urls-for-duo": "[testaccount.us-east-1.privatelink.snowflakecomputing.com, app-testaccount.us-east-1.privatelink.snowflakecomputing.com]"
		}`,
	}

	c, e := raw.GetStructuredConfig()
	r.NoError(e)

	// Common fields
	r.Equal("testaccount.us-east-1.privatelink", c.AccountName)
	r.Equal("testaccount.us-east-1.privatelink.snowflakecomputing.com", c.AccountURL)
	r.Equal("ocsp.testaccount.us-east-1.privatelink.snowflakecomputing.com", c.OCSPURL)
	r.Equal("testorg-testaccount.privatelink.snowflakecomputing.com", c.RegionlessAccountURL)
	r.Equal("app-testorg-testaccount.privatelink.snowflakecomputing.com", c.RegionlessSnowsightURL)
	r.Equal("app.us-east-1.privatelink.snowflakecomputing.com", c.SnowsightURL)
	r.Equal("ocsp.testorg-testaccount.privatelink.snowflakecomputing.com", c.RegionlessOCSPURL)
	r.Equal("[testaccount.us-east-1.privatelink.snowflakecomputing.com, app-testaccount.us-east-1.privatelink.snowflakecomputing.com]", c.DashedDuoURLs)
	// AWS-specific fields
	r.Equal("com.amazonaws.vpce.us-east-1.vpce-svc-00000000000000000", c.AwsVpceID)
	r.Equal("arn:aws:iam::000000000000:root", c.AccountPrincipal)
	r.Equal("*.testaccount.us-east-1.privatelink.snowflakecomputing.com", c.AppServiceURL)
	// Azure-specific fields absent for AWS
	r.Equal("", c.AzurePrivateLinkServiceID)
	r.Equal("", c.InternalStage)
	r.Equal("", c.AzureStorageVolumeNFS)
	r.Equal("", c.AzureStorageVolumeFS)
	// GCP-specific fields absent for AWS
	r.Equal("", c.GCPServiceAttachment)
	// Client redirect fields absent when client redirect is not configured
	r.Equal("", c.ConnectionURLs)
	r.Equal("", c.ConnectionOCSPURLs)
}

// TestSystemGetPrivateLinkGetStructuredConfigAwsAsPerDocumentation tests parsing of all AWS fields
// as documented at https://docs.snowflake.com/en/sql-reference/functions/system_get_privatelink_config,
// using the standard "privatelink-ocsp-url" key and including client redirect connection URLs.
func TestSystemGetPrivateLinkGetStructuredConfigAwsAsPerDocumentation(t *testing.T) {
	r := require.New(t)

	raw := &RawPrivateLinkConfig{
		Config: `{
			"privatelink-account-name": "testaccount.us-east-1.privatelink",
			"privatelink-vpce-id": "com.amazonaws.vpce.us-east-1.vpce-svc-00000000000000000",
			"privatelink-account-url": "testaccount.us-east-1.privatelink.snowflakecomputing.com",
			"privatelink-ocsp-url": "ocsp.testaccount.us-east-1.privatelink.snowflakecomputing.com",
			"privatelink-account-principal": "arn:aws:iam::000000000000:root",
			"app-service-privatelink-url": "*.testaccount.us-east-1.privatelink.snowflakecomputing.com",
			"privatelink-connection-urls": "testaccount.us-east-1.privatelink.snowflakecomputing.com:443",
			"privatelink-connection-ocsp-urls": "ocsp.testaccount.us-east-1.privatelink.snowflakecomputing.com:443",
			"regionless-privatelink-ocsp-url": "ocsp.testorg-testaccount.privatelink.snowflakecomputing.com",
			"privatelink-dashed-urls-for-duo": "[testaccount.us-east-1.privatelink.snowflakecomputing.com]"
		}`,
	}

	c, e := raw.GetStructuredConfig()
	r.NoError(e)

	// Common fields
	r.Equal("testaccount.us-east-1.privatelink", c.AccountName)
	r.Equal("testaccount.us-east-1.privatelink.snowflakecomputing.com", c.AccountURL)
	r.Equal("ocsp.testaccount.us-east-1.privatelink.snowflakecomputing.com", c.OCSPURL)
	r.Equal("ocsp.testorg-testaccount.privatelink.snowflakecomputing.com", c.RegionlessOCSPURL)
	r.Equal("[testaccount.us-east-1.privatelink.snowflakecomputing.com]", c.DashedDuoURLs)
	r.Equal("testaccount.us-east-1.privatelink.snowflakecomputing.com:443", c.ConnectionURLs)
	r.Equal("ocsp.testaccount.us-east-1.privatelink.snowflakecomputing.com:443", c.ConnectionOCSPURLs)
	// AWS-specific fields
	r.Equal("com.amazonaws.vpce.us-east-1.vpce-svc-00000000000000000", c.AwsVpceID)
	r.Equal("arn:aws:iam::000000000000:root", c.AccountPrincipal)
	r.Equal("*.testaccount.us-east-1.privatelink.snowflakecomputing.com", c.AppServiceURL)
	// Azure-specific fields absent for AWS
	r.Equal("", c.AzurePrivateLinkServiceID)
}

func TestSystemGetPrivateLinkGetStructuredConfigAzure(t *testing.T) {
	r := require.New(t)

	raw := &RawPrivateLinkConfig{
		Config: `{
			"privatelink-account-name": "testaccount.east-us-2.azure.privatelink",
			"privatelink-pls-id": "sf-pvlinksvc-azeastus2.east-us-2.azure.example",
			"privatelink-account-url": "testaccount.east-us-2.azure.privatelink.snowflakecomputing.com",
			"privatelink_ocsp-url": "ocsp.testaccount.east-us-2.azure.privatelink.snowflakecomputing.com",
			"privatelink-internal-stage": "sfcteststorage.blob.core.windows.net",
			"privatelink-snowflake-managed-storage-volume-nfs": "nfsstorage.blob.core.windows.net",
			"privatelink-snowflake-managed-storage-volume-fs": "fsstorage.blob.core.windows.net",
			"regionless-privatelink-ocsp-url": "ocsp.testorg-testaccount.privatelink.snowflakecomputing.com",
			"privatelink-dashed-urls-for-duo": "[testaccount.east-us-2.azure.privatelink.snowflakecomputing.com]"
		}`,
	}

	c, e := raw.GetStructuredConfig()
	r.NoError(e)

	// Common fields
	r.Equal("testaccount.east-us-2.azure.privatelink", c.AccountName)
	r.Equal("testaccount.east-us-2.azure.privatelink.snowflakecomputing.com", c.AccountURL)
	r.Equal("ocsp.testaccount.east-us-2.azure.privatelink.snowflakecomputing.com", c.OCSPURL)
	r.Equal("ocsp.testorg-testaccount.privatelink.snowflakecomputing.com", c.RegionlessOCSPURL)
	r.Equal("[testaccount.east-us-2.azure.privatelink.snowflakecomputing.com]", c.DashedDuoURLs)
	// Azure-specific fields
	r.Equal("sf-pvlinksvc-azeastus2.east-us-2.azure.example", c.AzurePrivateLinkServiceID)
	r.Equal("sfcteststorage.blob.core.windows.net", c.InternalStage)
	r.Equal("nfsstorage.blob.core.windows.net", c.AzureStorageVolumeNFS)
	r.Equal("fsstorage.blob.core.windows.net", c.AzureStorageVolumeFS)
	// AWS-specific fields absent for Azure
	r.Equal("", c.AwsVpceID)
	r.Equal("", c.AccountPrincipal)
	r.Equal("", c.AppServiceURL)
	// GCP-specific fields absent for Azure
	r.Equal("", c.GCPServiceAttachment)
	// Client redirect fields absent when client redirect is not configured
	r.Equal("", c.ConnectionURLs)
	r.Equal("", c.ConnectionOCSPURLs)
}

func TestSystemGetPrivateLinkGetStructuredConfigGcp(t *testing.T) {
	r := require.New(t)

	raw := &RawPrivateLinkConfig{
		Config: `{
			"privatelink-account-name": "testaccount.us-central1.gcp.privatelink",
			"privatelink-account-url": "testaccount.us-central1.gcp.privatelink.snowflakecomputing.com",
			"privatelink-ocsp-url": "ocsp.testaccount.us-central1.gcp.privatelink.snowflakecomputing.com",
			"privatelink-gcp-service-attachment": "projects/snowflake/regions/us-central1/serviceAttachments/testServiceAttachment",
			"regionless-privatelink-account-url": "testorg-testaccount.privatelink.snowflakecomputing.com",
			"regionless-privatelink-ocsp-url": "ocsp.testorg-testaccount.privatelink.snowflakecomputing.com",
			"privatelink-dashed-urls-for-duo": "[testaccount-us-central1-gcp.snowflakecomputing.com]"
		}`,
	}

	c, e := raw.GetStructuredConfig()
	r.NoError(e)

	// Common fields
	r.Equal("testaccount.us-central1.gcp.privatelink", c.AccountName)
	r.Equal("testaccount.us-central1.gcp.privatelink.snowflakecomputing.com", c.AccountURL)
	r.Equal("ocsp.testaccount.us-central1.gcp.privatelink.snowflakecomputing.com", c.OCSPURL)
	r.Equal("testorg-testaccount.privatelink.snowflakecomputing.com", c.RegionlessAccountURL)
	r.Equal("ocsp.testorg-testaccount.privatelink.snowflakecomputing.com", c.RegionlessOCSPURL)
	r.Equal("[testaccount-us-central1-gcp.snowflakecomputing.com]", c.DashedDuoURLs)
	// GCP-specific fields
	r.Equal("projects/snowflake/regions/us-central1/serviceAttachments/testServiceAttachment", c.GCPServiceAttachment)
	// AWS-specific fields absent for GCP
	r.Equal("", c.AwsVpceID)
	r.Equal("", c.AccountPrincipal)
	r.Equal("", c.AppServiceURL)
	// Azure-specific fields absent for GCP
	r.Equal("", c.AzurePrivateLinkServiceID)
	r.Equal("", c.InternalStage)
	r.Equal("", c.AzureStorageVolumeNFS)
	r.Equal("", c.AzureStorageVolumeFS)
	// Client redirect fields absent when client redirect is not configured
	r.Equal("", c.ConnectionURLs)
	r.Equal("", c.ConnectionOCSPURLs)
}
