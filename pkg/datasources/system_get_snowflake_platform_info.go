package datasources

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var systemGetSnowflakePlatformInfoSchema = map[string]*schema.Schema{
	"azure_vnet_subnet_ids": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Computed:    true,
		Description: "Snowflake Azure Virtual Network Subnet IDs",
	},
	"aws_vpc_ids": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Computed:    true,
		Description: "Snowflake AWS Virtual Private Cloud IDs",
	},
}

func SystemGetSnowflakePlatformInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.SystemGetSnowflakePlatformInfoDatasource), TrackingReadWrapper(datasources.SystemGetSnowflakePlatformInfo, ReadSystemGetSnowflakePlatformInfo)),
		Schema:      systemGetSnowflakePlatformInfoSchema,
	}
}

// ReadSystemGetSnowflakePlatformInfo implements schema.ReadFunc.
func ReadSystemGetSnowflakePlatformInfo(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB

	sel := snowflake.SystemGetSnowflakePlatformInfoQuery()
	row := snowflake.QueryRow(db, sel)

	acc, err := client.ContextFunctions.CurrentSessionDetails(context.Background())
	if err != nil {
		// If not found, mark resource to be removed from state file during apply or refresh
		d.SetId("")
		log.Println("[DEBUG] current_account failed to decode")
		return diag.FromErr(fmt.Errorf("error current_account err = %w", err))
	}

	d.SetId(fmt.Sprintf("%s.%s", acc.Account, acc.Region))

	rawInfo, err := snowflake.ScanSnowflakePlatformInfo(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Println("[DEBUG] system_get_snowflake_platform_info not found")
		return diag.FromErr(fmt.Errorf("error system_get_snowflake_platform_info err = %w", err))
	}

	info, err := rawInfo.GetStructuredConfig()
	if err != nil {
		log.Println("[DEBUG] system_get_snowflake_platform_info failed to decode")
		d.SetId("")
		return diag.FromErr(fmt.Errorf("error system_get_snowflake_platform_info err = %w", err))
	}

	if err := d.Set("azure_vnet_subnet_ids", info.AzureVnetSubnetIds); err != nil {
		return diag.FromErr(fmt.Errorf("error system_get_snowflake_platform_info err = %w", err))
	}

	if err := d.Set("aws_vpc_ids", info.AwsVpcIds); err != nil {
		return diag.FromErr(fmt.Errorf("error system_get_snowflake_platform_info err = %w", err))
	}

	return nil
}
