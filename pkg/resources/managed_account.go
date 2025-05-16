package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	SnowflakeReaderAccountType = "READER"
)

var managedAccountSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the managed account; must be unique for your account.",
		ForceNew:    true,
	},
	"admin_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier, as well as login name, for the initial user in the managed account. This user serves as the account administrator for the account.",
		ForceNew:    true,
	},
	"admin_password": {
		Type:        schema.TypeString,
		Required:    true,
		Sensitive:   true,
		Description: "Password for the initial user in the managed account. Check [Snowflake-provided password policy](https://docs.snowflake.com/en/user-guide/admin-user-management#snowflake-provided-password-policy).",
		ForceNew:    true,
	},
	"type": {
		Type:         schema.TypeString,
		Optional:     true,
		Default:      SnowflakeReaderAccountType,
		Description:  "Specifies the type of managed account.",
		ValidateFunc: validation.StringInSlice([]string{SnowflakeReaderAccountType}, true),
		ForceNew:     true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the managed account.",
		ForceNew:    true,
	},
	"cloud": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Cloud in which the managed account is located.",
	},
	"region": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Snowflake Region in which the managed account is located.",
	},
	"locator": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Display name of the managed account.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the managed account was created.",
	},
	"url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "URL for accessing the managed account, particularly through the web interface.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// ManagedAccount returns a pointer to the resource representing a managed account.
func ManagedAccount() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErr[sdk.AccountObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ManagedAccounts.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ManagedAccountResource), TrackingCreateWrapper(resources.ManagedAccount, CreateManagedAccount)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ManagedAccountResource), TrackingReadWrapper(resources.ManagedAccount, ReadManagedAccount)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ManagedAccountResource), TrackingDeleteWrapper(resources.ManagedAccount, deleteFunc)),

		Schema: managedAccountSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateManagedAccount implements schema.CreateFunc.
func CreateManagedAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)

	adminName := d.Get("admin_name").(string)
	adminPassword := d.Get("admin_password").(string)
	createParams := sdk.NewCreateManagedAccountParamsRequest(adminName, adminPassword)

	if v, ok := d.GetOk("comment"); ok {
		createParams.WithComment(v.(string))
	}

	createRequest := sdk.NewCreateManagedAccountRequest(id, *createParams)

	err := client.ManagedAccounts.Create(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadManagedAccount(ctx, d, meta)
}

// ReadManagedAccount implements schema.ReadFunc.
func ReadManagedAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	// We have to wait during the first read, since the locator takes some time to appear.
	// This approach has a downside of not handling correctly the situation where managed account was removed externally.
	// TODO [SNOW-1003380]: discuss it as a provider-wide topic during resources redesign.
	var managedAccount *sdk.ManagedAccount
	var err error
	err = util.Retry(5, 3*time.Second, func() (error, bool) {
		managedAccount, err = client.ManagedAccounts.ShowByIDSafely(ctx, id)
		if err != nil {
			log.Printf("[DEBUG] retryable operation resulted in error: %v", err)
			return nil, false
		}
		return nil, true
	})
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query managed account. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Managed account id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", managedAccount.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cloud", managedAccount.Cloud); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("region", managedAccount.Region); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("locator", managedAccount.Locator); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_on", managedAccount.CreatedOn); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url", managedAccount.URL); err != nil {
		return diag.FromErr(err)
	}

	if managedAccount.IsReader {
		if err := d.Set("type", "READER"); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("unable to determine the account type"))
	}

	if err := d.Set("comment", managedAccount.Comment); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
