package datasources

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var systemGenerateSCIMAccesstokenSchema = map[string]*schema.Schema{
	"integration_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "SCIM Integration Name",
	},
	"access_token": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "SCIM Access Token",
		Sensitive:   true,
	},
}

func SystemGenerateSCIMAccessToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.SystemGenerateSCIMAccessTokenDatasource), TrackingReadWrapper(datasources.SystemGenerateScimAccessToken, ReadSystemGenerateSCIMAccessToken)),
		Schema:      systemGenerateSCIMAccesstokenSchema,
	}
}

// ReadSystemGetAWSSNSIAMPolicy implements schema.ReadFunc.
func ReadSystemGenerateSCIMAccessToken(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB

	integrationName := sdk.NewAccountObjectIdentifier(d.Get("integration_name").(string)).Name()

	sel := snowflake.NewSystemGenerateSCIMAccessTokenBuilder(integrationName).Select()
	row := snowflake.QueryRow(db, sel)
	accessToken, err := snowflake.ScanSCIMAccessToken(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] system_generate_scim_access_token (%s) not found", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		log.Printf("[DEBUG] system_generate_scim_access_token (%s) failed to generate (%q)", d.Id(), err.Error())
		d.SetId("")
		return nil
	}

	d.SetId(integrationName)
	return diag.FromErr(d.Set("access_token", accessToken.Token))
}
