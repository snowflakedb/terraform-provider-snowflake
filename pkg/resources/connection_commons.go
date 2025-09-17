package resources

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DeleteConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Retry is necessary for cases where the changes in Snowflake didn't have enough time to propagate.
	// An example would be having the following setup:
	// 1. Config with primary connection (con1) and secondary connection (con2)
	// 2. Setting secondary connection (con2) as primary
	// 3. Setting previously primary connection (con1) as primary (we will have the connection dependencies as in the first step)
	// 4. Deleting secondary connection (con2; this may fail without waiting a bit for Snowflake to propagate the changes)
	if err := util.Retry(5, 2*time.Second, func() (error, bool) {
		if err := client.Connections.DropSafely(ctx, id); err != nil {
			log.Printf("[DEBUG] Drop secondary connection failed, err = %s", err)
			if strings.Contains(err.Error(), "is currently a primary connection in a replication relationship and cannot be dropped") {
				return nil, false
			}
			return err, true
		}
		return nil, true
	}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
