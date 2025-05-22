package datasources

import (
	"context"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var currentAccountSchema = map[string]*schema.Schema{
	"account": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake Account ID; as returned by CURRENT_ACCOUNT().",
	},

	"region": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake Region; as returned by CURRENT_REGION()",
	},

	"url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake URL.",
	},
}

// CurrentAccount the Snowflake current account resource.
func CurrentAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.CurrentAccountDatasource), TrackingReadWrapper(datasources.CurrentAccount, ReadCurrentAccount)),
		Schema:      currentAccountSchema,
	}
}

// ReadCurrentAccount read the current snowflake account information.
func ReadCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	current, err := client.ContextFunctions.CurrentSessionDetails(ctx)
	if err != nil {
		log.Println("[DEBUG] current_account failed to decode")
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%s.%s", current.Account, current.Region))
	accountErr := d.Set("account", current.Account)
	if accountErr != nil {
		return diag.FromErr(accountErr)
	}
	regionErr := d.Set("region", current.Region)
	if regionErr != nil {
		return diag.FromErr(regionErr)
	}
	url, err := current.AccountURL()
	if err != nil {
		log.Println("[DEBUG] generating snowflake url failed")
		return nil
	}

	urlErr := d.Set("url", url)
	if urlErr != nil {
		return diag.FromErr(urlErr)
	}
	return nil
}
