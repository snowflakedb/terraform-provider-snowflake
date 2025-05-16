package datasources

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var systemGetPrivateLinkConfigSchema = map[string]*schema.Schema{
	"account_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The name of your Snowflake account.",
	},

	"account_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The URL used to connect to Snowflake through AWS PrivateLink or Azure Private Link.",
	},

	"ocsp_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The OCSP URL corresponding to your Snowflake account that uses AWS PrivateLink or Azure Private Link.",
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

	"regionless_account_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The regionless URL to connect to your Snowflake account using AWS PrivateLink, Azure Private Link, or Google Cloud Private Service Connect.",
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

	"internal_stage": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The endpoint to connect to your Snowflake internal stage using AWS PrivateLink or Azure Private Link.",
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
	accURLErr := d.Set("account_url", config.AccountURL)
	if accURLErr != nil {
		return diag.FromErr(accURLErr)
	}
	ocspURLErr := d.Set("ocsp_url", config.OCSPURL)
	if ocspURLErr != nil {
		return diag.FromErr(ocspURLErr)
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

	if config.InternalStage != "" {
		intStgErr := d.Set("internal_stage", config.InternalStage)
		if intStgErr != nil {
			return diag.FromErr(intStgErr)
		}
	}

	if config.SnowsightURL != "" {
		snowSigURLErr := d.Set("snowsight_url", config.SnowsightURL)
		if snowSigURLErr != nil {
			return diag.FromErr(snowSigURLErr)
		}
	}

	if config.RegionlessSnowsightURL != "" {
		reglssSnowURLErr := d.Set("regionless_snowsight_url", config.RegionlessSnowsightURL)
		if reglssSnowURLErr != nil {
			return diag.FromErr(reglssSnowURLErr)
		}
	}

	if config.RegionlessAccountURL != "" {
		reglssAccURLErr := d.Set("regionless_account_url", config.RegionlessAccountURL)
		if reglssAccURLErr != nil {
			return diag.FromErr(reglssAccURLErr)
		}
	}

	return nil
}
