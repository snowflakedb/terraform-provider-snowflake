package resources

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	shareCreateExt = PreviewFeatureCreateContextWrapper(string(previewfeatures.ShareResource), TrackingCreateWrapper(resources.Share, CreateShare))
	shareReadExt   = PreviewFeatureReadContextWrapper(string(previewfeatures.ShareResource), TrackingReadWrapper(resources.Share, ReadShareFunc(true)))
	shareUpdateExt = PreviewFeatureUpdateContextWrapper(string(previewfeatures.ShareResource), TrackingUpdateWrapper(resources.Share, UpdateShare))
	// The delete func decodes the id with the error variant of the legacy decoder, unlike read and update.
	shareDeleteExt = PreviewFeatureDeleteContextWrapper(string(previewfeatures.ShareResource), TrackingDeleteWrapper(resources.Share, ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErrLegacy[sdk.AccountObjectIdentifier],
		shareDropSafelyInSdkExt,
	)))
	shareDescriptionExt   = ""
	shareCustomizeDiffExt schema.CustomizeDiffFunc
	shareImporterExt      = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.Share, ImportName[sdk.AccountObjectIdentifier]),
	}
)

// =============================================================================
// Request / mapping / Read-Update helpers
// =============================================================================

// shareAlterRequestsExt carries what changed, because the accounts are altered with several statements.
type shareAlterRequestsExt struct {
	accountsChanged bool
	oldAccounts     []string
	newAccounts     []string
	commentChanged  bool
	comment         string
}

func shareParseIdFromConfigExt(d *schema.ResourceData) (sdk.AccountObjectIdentifier, error) {
	return sdk.NewAccountObjectIdentifier(d.Get("name").(string)), nil
}

func shareParseIdExt(id string) (sdk.AccountObjectIdentifier, error) {
	return helpers.DecodeSnowflakeIDLegacy(id).(sdk.AccountObjectIdentifier), nil
}

func shareNewCreateRequestExt(_ *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateShareRequest, error) {
	return sdk.NewCreateShareRequest(id), nil
}

func shareApplyCreateOptionalsExt(d *schema.ResourceData, req *sdk.CreateShareRequest) error {
	if comment := d.Get("comment").(string); comment != "" {
		req.WithComment(comment)
	}
	return nil
}

// shareAfterCreateExt adds the accounts after the share is created. A failure is returned without reading the share.
func shareAfterCreateExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) diag.Diagnostics {
	accounts := expandStringList(d.Get("accounts").([]any))
	if len(accounts) == 0 {
		return nil
	}
	if err := setShareAccounts(ctx, meta.(*provider.Context).Client, id, accountIdentifiersFromSlice(accounts)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func shareSetConfigFieldsExt(d *schema.ResourceData, share *sdk.Share) error {
	if err := d.Set("comment", share.Comment); err != nil {
		return err
	}
	accounts := collections.Map(share.To, sdk.AccountIdentifier.Name)

	currentAccount := d.Get("accounts")
	if currentAccount != nil {
		currentAccounts := expandStringList(currentAccount.([]any))
		// reorder the accounts so they match the order in the config
		// this is to avoid unnecessary diffs
		accounts = reorderStringList(currentAccounts, accounts)
	}
	return d.Set("accounts", accounts)
}

func shareApplySetUnsetExt(d *schema.ResourceData) (*shareAlterRequestsExt, error) {
	reqs := &shareAlterRequestsExt{}
	if d.HasChange("accounts") {
		o, n := d.GetChange("accounts")
		reqs.accountsChanged = true
		reqs.oldAccounts = expandStringList(o.([]any))
		reqs.newAccounts = expandStringList(n.([]any))
	}
	if d.HasChange("comment") {
		reqs.commentChanged = true
		reqs.comment = d.Get("comment").(string)
	}
	return reqs, nil
}

// =============================================================================
// SDK incision points
// =============================================================================

func shareCreateInSdkExt(ctx context.Context, meta any, _ sdk.AccountObjectIdentifier, req *sdk.CreateShareRequest) diag.Diagnostics {
	if err := meta.(*provider.Context).Client.Shares.Create(ctx, req); err != nil {
		// The id is not set yet at this point, so the message has always had an empty id.
		return diag.FromErr(fmt.Errorf("error creating share (%v) err = %w", "", err))
	}
	return nil
}

// shareShowByIdSafelyInSdkExt keeps the read of the handwritten resource: a share that is not found is an error, not a
// removal from the state. SHOW SHARES lists at most three recipients and reports the rest only as "... (N TOTAL)", so
// the full list is merged in from SNOWFLAKE.ACCOUNT_USAGE.SHARES. A failed lookup is returned as a warning next to the share.
func shareShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (*sdk.Share, diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	share, err := client.Shares.ShowByID(ctx, id)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error reading share (%v) err = %w", d.Id(), err))
	}
	if share.ToTotal > len(share.To) {
		targetAccounts, err := client.Shares.ShowTargetAccounts(ctx, id)
		if err != nil {
			// Not an error: the first three recipients are still read, and the warning is returned with the share.
			return share, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Share %s has %d recipients but SHOW SHARES lists only %d", d.Id(), share.ToTotal, len(share.To)),
					Detail: "The full list is read from SNOWFLAKE.ACCOUNT_USAGE.SHARES, which failed (the role needs SELECT on that view, for example through IMPORTED PRIVILEGES on the SNOWFLAKE database). " +
						"Recipients that are not listed are missing from the state and planned as additions. Error: " + err.Error(),
				},
			}
		}
		// ACCOUNT_USAGE.SHARES can be up to 3 hours behind; SHOW SHARES is real time, so a recipient added since the
		// view was last refreshed may be among the listed three but not in the view yet.
		// CompactFunc only collapses adjacent duplicates, so sort by fully qualified name first.
		recipients := append(slices.Clone(share.To), targetAccounts...)
		slices.SortFunc(recipients, func(a, b sdk.AccountIdentifier) int {
			return strings.Compare(a.FullyQualifiedName(), b.FullyQualifiedName())
		})
		share.To = slices.CompactFunc(recipients, func(a, b sdk.AccountIdentifier) bool {
			return a.FullyQualifiedName() == b.FullyQualifiedName()
		})
	}
	return share, nil
}

func shareAlterSetInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier, reqs *shareAlterRequestsExt) error {
	client := meta.(*provider.Context).Client

	if reqs.accountsChanged {
		if len(reqs.newAccounts) == 0 {
			accountIdentifiers := accountIdentifiersFromSlice(reqs.oldAccounts)
			err := client.Shares.Alter(ctx, sdk.NewAlterShareRequest(id).WithRemove(sdk.ShareRemoveRequest{Accounts: accountIdentifiers}))
			if err != nil {
				return fmt.Errorf("error removing accounts from share (%v) err = %w", id.Name(), err)
			}
		} else {
			accountIdentifiers := accountIdentifiersFromSlice(reqs.newAccounts)

			share, err := client.Shares.ShowByID(ctx, id)
			if err != nil {
				return fmt.Errorf("error reading share (%v) err = %w", id.Name(), err)
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
				return err
			}
		}
	}
	if reqs.commentChanged {
		err := client.Shares.Alter(ctx, sdk.NewAlterShareRequest(id).WithSet(sdk.ShareSetRequest{Comment: sdk.String(reqs.comment)}))
		if err != nil {
			return fmt.Errorf("error updating share (%v) comment err = %w", id.Name(), err)
		}
	}
	return nil
}

func shareDropSafelyInSdkExt(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
	return client.Shares.DropSafely
}

// =============================================================================
// Pre-existing mappers
// =============================================================================

func setShareAccounts(ctx context.Context, client *sdk.Client, shareID sdk.AccountObjectIdentifier, accounts []sdk.AccountIdentifier) error {
	// There is a race condition where error accounts cannot be added to a
	// share until after a database is added to the share. Since a database
	// grant is dependent on the share itself, this is a hack to get the
	// thing working.

	// 1. Create new temporary DB
	tempName := fmt.Sprintf("TEMP_%v_%d", shareID.Name(), time.Now().Unix())
	tempDatabaseID := sdk.NewAccountObjectIdentifier(tempName)
	err := client.Databases.Create(ctx, sdk.NewCreateDatabaseRequest(tempDatabaseID))
	if err != nil {
		return fmt.Errorf("error creating temporary DB %v err = %w", tempName, err)
	}
	defer func() {
		// drop the temporary DB during cleanup
		err = client.Databases.Drop(ctx, sdk.NewDropDatabaseRequest(tempDatabaseID))
		if err != nil {
			log.Printf("[WARN] error dropping temporary DB %v err = %v", tempName, err)
		}
	}()
	// 2. Create temporary DB grant to the share
	// USAGE can only be granted to one database - granting USAGE on the temp db here
	// conflicts (and errors) with having a database already shared (i.e. when you
	// already have a share and are just adding or removing accounts). Instead, use
	// USAGE and REFERENCE_USAGE which is intended for multi-database sharing as per Snowflake
	// documentation here:
	// https://docs.snowflake.com/en/sql-reference/sql/grant-privilege-share.html#usage-notes
	// Note however that USAGE will be granted automatically on the temp db for the
	// case where the main db doesn't already exist, so it will need to be revoked
	// before deleting the temp db. Where USAGE hasn't been already granted it is not
	// an error to revoke it, so it's ok to just do the revoke every time.
	err = client.Grants.GrantPrivilegeToShare(ctx, sdk.NewGrantPrivilegeToShareRequest(shareID).
		WithPrivileges([]sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage, sdk.ObjectPrivilegeReferenceUsage}).
		WithOn(*sdk.NewShareGrantOnRequest().WithDatabase(tempDatabaseID)))
	if err != nil {
		return fmt.Errorf("error granting privilege to share (%v) err = %w", shareID.Name(), err)
	}
	defer func() {
		// revoke the REFERENCE_USAGE privilege during cleanup
		err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.NewRevokePrivilegeFromShareRequest(shareID).
			WithPrivileges([]sdk.ObjectPrivilege{sdk.ObjectPrivilegeReferenceUsage}).
			WithOn(*sdk.NewShareGrantOnRequest().WithDatabase(tempDatabaseID)))
		if err != nil {
			log.Printf("[WARN] error revoking privilege from share (%v) err = %v", shareID.Name(), err)
		}
		// revoke the maybe automatically granted USAGE privilege during cleanup
		err = client.Grants.RevokePrivilegeFromShare(ctx, sdk.NewRevokePrivilegeFromShareRequest(shareID).
			WithPrivileges([]sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}).
			WithOn(*sdk.NewShareGrantOnRequest().WithDatabase(tempDatabaseID)))
		if err != nil {
			log.Printf("[WARN] error revoking privilege from share (%v) err = %v", shareID.Name(), err)
		}
	}()
	// 3. Add accounts to the share
	err = client.Shares.Alter(ctx, sdk.NewAlterShareRequest(shareID).WithAdd(sdk.ShareAddRequest{Accounts: accounts}))
	return err
}

func accountIdentifiersFromSlice(accounts []string) []sdk.AccountIdentifier {
	accountIdentifiers := make([]sdk.AccountIdentifier, len(accounts))
	for i, account := range accounts {
		parts := strings.Split(account, ".")
		orgName := parts[0]
		accountName := parts[1]
		accountIdentifiers[i] = sdk.NewAccountIdentifier(orgName, accountName)
	}
	return accountIdentifiers
}
