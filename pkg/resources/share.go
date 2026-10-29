package resources

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Share() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErrLegacy[sdk.AccountObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] { return client.Shares.DropSafely },
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ShareResource), TrackingCreateWrapper(resources.Share, CreateShare)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ShareResource), TrackingReadWrapper(resources.Share, ReadShare)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ShareResource), TrackingUpdateWrapper(resources.Share, UpdateShare)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ShareResource), TrackingDeleteWrapper(resources.Share, deleteFunc)),

		Schema: shareSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Share, ImportName[sdk.AccountObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateShare implements schema.CreateFunc.
func CreateShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)

	comment := d.Get("comment").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	req := sdk.NewCreateShareRequest(id)
	if comment != "" {
		req.WithComment(comment)
	}
	if err := client.Shares.Create(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating share (%v) err = %w", d.Id(), err))
	}
	d.SetId(name)

	accounts := expandStringList(d.Get("accounts").([]any))
	if len(accounts) > 0 {
		shareID := sdk.NewAccountObjectIdentifier(name)
		accountIdentifiers := make([]sdk.AccountIdentifier, len(accounts))
		for i, account := range accounts {
			parts := strings.Split(account, ".")
			orgName := parts[0]
			accountName := parts[1]
			accountIdentifiers[i] = sdk.NewAccountIdentifier(orgName, accountName)
		}
		err := setShareAccounts(ctx, client, shareID, accountIdentifiers)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadShare(ctx, d, meta)
}

// ReadShare implements schema.ReadFunc.
func ReadShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.AccountObjectIdentifier)

	share, err := client.Shares.ShowByID(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading share (%v) err = %w", d.Id(), err))
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("comment", share.Comment); err != nil {
		return diag.FromErr(err)
	}
	recipients := share.To
	var diags diag.Diagnostics
	// SHOW SHARES lists at most three recipients and reports the rest only as "... (N TOTAL)".
	if share.ToTotal > len(share.To) {
		targetAccounts, err := client.Shares.ShowTargetAccounts(ctx, id)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Share %s has %d recipients but SHOW SHARES lists only %d", d.Id(), share.ToTotal, len(share.To)),
				Detail: "The full list is read from SNOWFLAKE.ACCOUNT_USAGE.SHARES, which failed (the role needs SELECT on that view, for example through IMPORTED PRIVILEGES on the SNOWFLAKE database). " +
					"Recipients that are not listed are missing from the state and planned as additions. Error: " + err.Error(),
			})
		} else {
			// ACCOUNT_USAGE.SHARES can be up to 3 hours behind; SHOW SHARES is real time, so a recipient added since the
			// view was last refreshed may be among the listed three but not in the view yet.
			// CompactFunc only collapses adjacent duplicates, so sort by fully qualified name first.
			recipients = append(slices.Clone(share.To), targetAccounts...)
			slices.SortFunc(recipients, func(a, b sdk.AccountIdentifier) int {
				return strings.Compare(a.FullyQualifiedName(), b.FullyQualifiedName())
			})
			recipients = slices.CompactFunc(recipients, func(a, b sdk.AccountIdentifier) bool {
				return a.FullyQualifiedName() == b.FullyQualifiedName()
			})
		}
	}
	accounts := collections.Map(recipients, sdk.AccountIdentifier.Name)

	currentAccount := d.Get("accounts")
	if currentAccount != nil {
		currentAccounts := expandStringList(currentAccount.([]any))
		// reorder the accounts so they match the order in the config
		// this is to avoid unnecessary diffs
		accounts = reorderStringList(currentAccounts, accounts)
	}
	if err := d.Set("accounts", accounts); err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateShare implements schema.UpdateFunc.
func UpdateShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.AccountObjectIdentifier)
	client := meta.(*provider.Context).Client

	if d.HasChange("accounts") {
		o, n := d.GetChange("accounts")
		oldAccounts := expandStringList(o.([]any))
		newAccounts := expandStringList(n.([]any))
		if len(newAccounts) == 0 {
			accountIdentifiers := accountIdentifiersFromSlice(oldAccounts)
			err := client.Shares.Alter(ctx, sdk.NewAlterShareRequest(id).WithRemove(sdk.ShareRemoveRequest{Accounts: accountIdentifiers}))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error removing accounts from share (%v) err = %w", d.Id(), err))
			}
		} else {
			accountIdentifiers := accountIdentifiersFromSlice(newAccounts)

			share, err := client.Shares.ShowByID(ctx, id)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error reading share (%v) err = %w", d.Id(), err))
			}

			// If a database is already granted to the share, we can add accounts directly.
			// The setShareAccounts workaround (which creates a temp database) is only needed
			// when no database has been granted yet — Snowflake requires at least one database
			// on a share before accounts can be added.
			if share.DatabaseName.Name() != "" {
				err = client.Shares.Alter(ctx, sdk.NewAlterShareRequest(id).WithAdd(sdk.ShareAddRequest{Accounts: accountIdentifiers}))
			} else {
				err = setShareAccounts(ctx, client, id, accountIdentifiers)
			}
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		err := client.Shares.Alter(ctx, sdk.NewAlterShareRequest(id).WithSet(sdk.ShareSetRequest{Comment: sdk.String(comment)}))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating share (%v) comment err = %w", d.Id(), err))
		}
	}

	return ReadShare(ctx, d, meta)
}
