package datasources

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var systemGetPrivateLinkConfigSchema = map[string]*schema.Schema{
	"account_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The name of your Snowflake account.",
	},

	"account_principal": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The AWS principal ARN to allow for outbound private connections to your VPC endpoint services.",
	},

	"account_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The URL to connect to your Snowflake account using AWS PrivateLink, Azure Private Link, or Google Cloud Private Service Connect.",
	},

	"app_service_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The PrivateLink endpoint URL used to route traffic to Snowflake-hosted app services, such as Streamlit or Notebooks.",
	},

	"aws_vpce_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The AWS VPCE ID for your account.",
	},

	"azure_pls_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Microsoft Azure Private Link Service ID for your account identifier in the format of an alias.",
	},

	"azure_storage_volume_fs": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The endpoint for failsafe Snowflake-managed storage volumes when using Azure Private Link.",
	},

	"azure_storage_volume_nfs": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The endpoint for non-failsafe Snowflake-managed storage volumes when using Azure Private Link.",
	},

	"dashed_duo_urls": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The list of dashed variant URLs for Duo Multi-Factor Authentication, shown only when the hostname contains an underscore.",
	},

	"gcp_service_attachment": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The endpoint for the Snowflake service when using Google Cloud Private Service Connect.",
	},

	"internal_stage": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The endpoint to connect to your Snowflake internal stage using AWS PrivateLink or Azure Private Link.",
	},

	"ocsp_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The OCSP URL corresponding to your Snowflake account identifier.",
	},

	"connection_ocsp_urls": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The list of OCSP URLs for use with redirecting client connections when using client redirect.",
	},

	"connection_urls": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The private connectivity connection URLs for your account when using client redirect.",
	},

	"regionless_account_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The regionless URL to connect to your Snowflake account using AWS PrivateLink, Azure Private Link, or Google Cloud Private Service Connect.",
	},

	"regionless_ocsp_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The regionless OCSP URL to connect to Snowflake OCSP using private connectivity.",
	},

	"regionless_snowsight_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The URL for your organization to access Snowsight using Private Connectivity to the Snowflake Service.",
	},

	"snowsight_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The URL containing the cloud region to access Snowsight and the Snowflake Marketplace using Private Connectivity to the Snowflake Service.",
	},
}

func SystemGetPrivateLinkConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.SystemGetPrivateLinkConfigDatasource), TrackingReadWrapper(datasources.SystemGetPrivateLinkConfig, ReadSystemGetPrivateLinkConfig)),
		Schema:      systemGetPrivateLinkConfigSchema,
	}
}

// ReadSystemGetPrivateLinkConfig implements schema.ReadFunc.
func ReadSystemGetPrivateLinkConfig(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB

	sel := snowflake.SystemGetPrivateLinkConfigQuery()
	row := snowflake.QueryRow(db, sel)
	rawConfig, err := snowflake.ScanPrivateLinkConfig(row)

	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Println("[DEBUG] system_get_privatelink_config not found")
		d.SetId("")
		return nil
	}

	config, err := rawConfig.GetStructuredConfig()
	if err != nil {
		log.Println("[DEBUG] system_get_privatelink_config failed to decode")
		d.SetId("")
		return nil
	}

	d.SetId(config.AccountName)

	errs := errors.Join(
		d.Set("account_name", config.AccountName),
		d.Set("account_principal", config.AccountPrincipal),
		d.Set("account_url", config.AccountURL),
		d.Set("app_service_url", config.AppServiceURL),
		d.Set("aws_vpce_id", config.AwsVpceID),
		d.Set("azure_pls_id", config.AzurePrivateLinkServiceID),
		d.Set("azure_storage_volume_fs", config.AzureStorageVolumeFS),
		d.Set("azure_storage_volume_nfs", config.AzureStorageVolumeNFS),
		d.Set("dashed_duo_urls", config.DashedDuoURLs),
		d.Set("gcp_service_attachment", config.GCPServiceAttachment),
		d.Set("internal_stage", config.InternalStage),
		d.Set("ocsp_url", config.OCSPURL),
		d.Set("connection_ocsp_urls", config.ConnectionOCSPURLs),
		d.Set("connection_urls", config.ConnectionURLs),
		d.Set("regionless_account_url", config.RegionlessAccountURL),
		d.Set("regionless_ocsp_url", config.RegionlessOCSPURL),
		d.Set("regionless_snowsight_url", config.RegionlessSnowsightURL),
		d.Set("snowsight_url", config.SnowsightURL),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	return nil
}
