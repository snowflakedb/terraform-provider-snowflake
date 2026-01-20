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
		Description: "The URL used to connect to Snowflake through AWS PrivateLink or Azure Private Link.",
	},

	"app_service_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The wildcard URL required for routing Streamlit applications and Snowpark Container Services through AWS PrivateLink or Azure Private Link.",
	},

	"aws_vpce_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The AWS VPCE ID for your account.",
	},

	"azure_pls_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Azure Private Link Service ID for your account.",
	},

	"client_redirect_urls": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The client redirect connection URLs for your account using private connectivity.",
	},

	"dashed_duo_urls": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The dashed URLs for Duo integration using private connectivity.",
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
		Description: "The OCSP URL corresponding to your Snowflake account that uses AWS PrivateLink or Azure Private Link.",
	},

	"openflow_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The OpenFlow URL to connect to Snowflake Openflow using AWS PrivateLink or Azure Private Link.",
	},

	"openflow_telemetry_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The OpenFlow telemetry URL to connect to Snowflake Openflow telemetry using AWS PrivateLink or Azure Private Link.",
	},

	"regionless_account_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The regionless URL to connect to your Snowflake account using AWS PrivateLink, Azure Private Link, or Google Cloud Private Service Connect.",
	},

	"regionless_ocsp_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The OCSP URL for your account identifier.",
	},

	"spcs_auth_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowpark Container Services authentication URL using private connectivity.",
	},

	"spcs_registry_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowpark Container Services registry URL using private connectivity.",
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

	accNameErr := d.Set("account_name", config.AccountName)
	if accNameErr != nil {
		return diag.FromErr(accNameErr)
	}

	if config.AccountPrincipal != "" {
		accPrincipalErr := d.Set("account_principal", config.AccountPrincipal)
		if accPrincipalErr != nil {
			return diag.FromErr(accPrincipalErr)
		}
	}

	accURLErr := d.Set("account_url", config.AccountURL)
	if accURLErr != nil {
		return diag.FromErr(accURLErr)
	}

	if config.AppServiceURL != "" {
		appServiceURLErr := d.Set("app_service_url", config.AppServiceURL)
		if appServiceURLErr != nil {
			return diag.FromErr(appServiceURLErr)
		}
	}

	if config.AwsVpceID != "" {
		awsVpceIDErr := d.Set("aws_vpce_id", config.AwsVpceID)
		if awsVpceIDErr != nil {
			return diag.FromErr(awsVpceIDErr)
		}
	}

	if config.AzurePrivateLinkServiceID != "" {
		azurePlsIDErr := d.Set("azure_pls_id", config.AzurePrivateLinkServiceID)
		if azurePlsIDErr != nil {
			return diag.FromErr(azurePlsIDErr)
		}
	}

	if config.ConnectionURLs != "" {
		connURLsErr := d.Set("client_redirect_connection_urls", config.ConnectionURLs)
		if connURLsErr != nil {
			return diag.FromErr(connURLsErr)
		}
	}

	dashedDusURLsErr := d.Set("dashed_duo_urls", config.DashedDuoURLs)
	if dashedDusURLsErr != nil {
		return diag.FromErr(dashedDusURLsErr)
	}

	if config.GCPServiceAttachment != "" {
		gcpSvcAttachErr := d.Set("gcp_service_attachment", config.GCPServiceAttachment)
		if gcpSvcAttachErr != nil {
			return diag.FromErr(gcpSvcAttachErr)
		}
	}

	if config.InternalStage != "" {
		intStgErr := d.Set("internal_stage", config.InternalStage)
		if intStgErr != nil {
			return diag.FromErr(intStgErr)
		}
	}

	if config.OCSPURL != "" {
		ocspURLErr := d.Set("ocsp_url", config.OCSPURL)
		if ocspURLErr != nil {
			return diag.FromErr(ocspURLErr)
		}
	}

	if config.OpenflowURL != "" {
		openflowURLErr := d.Set("openflow_url", config.OpenflowURL)
		if openflowURLErr != nil {
			return diag.FromErr(openflowURLErr)
		}
	}

	if config.OpenflowTelemetryURL != "" {
		openflowTelemetryURLErr := d.Set("openflow_telemetry_url", config.OpenflowTelemetryURL)
		if openflowTelemetryURLErr != nil {
			return diag.FromErr(openflowTelemetryURLErr)
		}
	}

	if config.RegionlessAccountURL != "" {
		reglssAccURLErr := d.Set("regionless_account_url", config.RegionlessAccountURL)
		if reglssAccURLErr != nil {
			return diag.FromErr(reglssAccURLErr)
		}
	}

	reglssOCSPURLErr := d.Set("regionless_ocsp_url", config.RegionlessOCSPURL)
	if reglssOCSPURLErr != nil {
		return diag.FromErr(reglssOCSPURLErr)
	}

	if config.RegionlessSnowsightURL != "" {
		reglssSnowURLErr := d.Set("regionless_snowsight_url", config.RegionlessSnowsightURL)
		if reglssSnowURLErr != nil {
			return diag.FromErr(reglssSnowURLErr)
		}
	}
	if config.SnowparkCSAuthURL != "" {
		spcsAuthURLErr := d.Set("spcs_auth_url", config.SnowparkCSAuthURL)
		if spcsAuthURLErr != nil {
			return diag.FromErr(spcsAuthURLErr)
		}
	}

	if config.SnowparkCSRegistryURL != "" {
		spcsRegURLErr := d.Set("spcs_registry_url", config.SnowparkCSRegistryURL)
		if spcsRegURLErr != nil {
			return diag.FromErr(spcsRegURLErr)
		}
	}

	if config.SnowsightURL != "" {
		snowSigURLErr := d.Set("snowsight_url", config.SnowsightURL)
		if snowSigURLErr != nil {
			return diag.FromErr(snowSigURLErr)
		}
	}

	return nil
}
