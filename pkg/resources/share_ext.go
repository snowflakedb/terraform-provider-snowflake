package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
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

	shareDescriptionExt = "Resource used to manage share objects. For more information, check [share documentation](https://docs.snowflake.com/en/sql-reference/sql/create-share)."
	shareImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.Share, ImportShare),
	}
)

// =============================================================================
// Request / mapping / Read-Update helpers
// =============================================================================

// shareAlterRequestsExt carries the SET and UNSET requests of the generated update. The accounts are altered in shareBeforeAlterExt.
type shareAlterRequestsExt struct {
	set   *sdk.ShareSetRequest
	unset *sdk.ShareUnsetRequest
}

// shareSetFieldsNotSetByReadExt sets what the read does not: the name and secure_objects_only. The comment and the accounts are set by the read.
// Only an outbound share (created in this account) can be imported.
func shareSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, share *sdk.Share, importPermissiveMode bool) error {
	if share.Kind != sdk.ShareKindOutbound {
		return fmt.Errorf("invalid share kind, expected %s, got %s: only a share created in this account can be imported", sdk.ShareKindOutbound, share.Kind)
	}
	return errors.Join(
		d.Set("name", id.Name()),
		importOptionalBooleanStringAttribute(d, "secure_objects_only", importPermissiveMode, share.SecureObjectsOnly),
	)
}

// shareAfterCreateExt adds the accounts after the share is created. A failure is returned without reading the share.
func shareAfterCreateExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) diag.Diagnostics {
	accounts := expandStringList(d.Get("accounts").(*schema.Set).List())
	if len(accounts) == 0 {
		return nil
	}
	accountIds, err := accountIdentifiersFromSlice(accounts)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := setShareAccounts(ctx, meta.(*provider.Context).Client, id, accountIds); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// shareSetDescribeOutputExt sets describe_output, which is a list of rows. A failed DESCRIBE is only logged, because describing
// a share needs a privilege on it and the rest of the state does not depend on it.
func shareSetDescribeOutputExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, details []sdk.ShareDetails, describeErr error, _ bool) error {
	if describeErr != nil {
		log.Printf("[WARN] describing share: %s, err: %s", id.FullyQualifiedName(), describeErr)
		return nil
	}
	return d.Set(DescribeOutputAttributeName, schemas.ShareDetailsListToSchema(details))
}

func shareSetConfigFieldsExt(_ context.Context, d *schema.ResourceData, _ any, share *sdk.Share) error {
	if err := d.Set("comment", share.Comment); err != nil {
		return err
	}
	return d.Set("accounts", collections.Map(share.To, sdk.AccountIdentifier.Name))
}

// =============================================================================
// SDK incision points
// =============================================================================

// shareDescribeInSdkExt returns a row for every object granted to the share. The read does not fail when DESCRIBE fails.
func shareDescribeInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) ([]sdk.ShareDetails, error) {
	return meta.(*provider.Context).Client.Shares.Describe(ctx, id)
}

// shareAfterShowByIdExt completes the recipients of the SHOW row. A failed lookup of all the recipients is not an error: the
// listed recipients are used and a warning is returned.
func shareAfterShowByIdExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, share *sdk.Share) diag.Diagnostics {
	recipients, err := shareRecipientsExt(ctx, meta.(*provider.Context).Client.Shares.ShowTargetAccounts, id, share)
	share.To = recipients
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Share %s has %d recipients but SHOW SHARES lists only %d", d.Id(), share.ToTotal, len(share.To)),
				Detail: "The full list is read from SNOWFLAKE.ACCOUNT_USAGE.SHARES, which failed (the role needs SELECT on that view, for example through IMPORTED PRIVILEGES on the SNOWFLAKE database). " +
					"Recipients that are not listed are missing from the state and planned as additions. Error: " + err.Error(),
			},
		}
	}
	return nil
}

// shareRecipientsExt returns all the recipients of the share. SHOW SHARES lists at most three recipients and reports the rest
// only as "... (N TOTAL)", so the full list is merged in from SNOWFLAKE.ACCOUNT_USAGE.SHARES. When that lookup fails, the
// listed recipients are returned together with the error.
func shareRecipientsExt(ctx context.Context, showTargetAccounts func(context.Context, sdk.AccountObjectIdentifier) ([]sdk.AccountIdentifier, error), id sdk.AccountObjectIdentifier, share *sdk.Share) ([]sdk.AccountIdentifier, error) {
	if share.ToTotal <= len(share.To) {
		return share.To, nil
	}
	targetAccounts, err := showTargetAccounts(ctx, id)
	if err != nil {
		return share.To, err
	}
	// ACCOUNT_USAGE.SHARES can be up to 3 hours behind; SHOW SHARES is real time, so a recipient added since the
	// view was last refreshed may be among the listed three but not in the view yet.
	// CompactFunc only collapses adjacent duplicates, so sort by fully qualified name first.
	recipients := append(slices.Clone(share.To), targetAccounts...)
	slices.SortFunc(recipients, func(a, b sdk.AccountIdentifier) int {
		return strings.Compare(a.FullyQualifiedName(), b.FullyQualifiedName())
	})
	return slices.CompactFunc(recipients, func(a, b sdk.AccountIdentifier) bool {
		return a.FullyQualifiedName() == b.FullyQualifiedName()
	}), nil
}

// shareBeforeAlterExt makes the recipients equal to the configured accounts: it removes the accounts that are not configured
// anymore and adds the missing ones. SET and UNSET of the other fields run afterwards in the generated update.
func shareBeforeAlterExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) diag.Diagnostics {
	if !d.HasChange("accounts") {
		return nil
	}
	// The read sets the accounts in state to the recipients that Snowflake reports, so the difference between the state
	// and the configuration is what the update has to add and remove.
	o, n := d.GetChange("accounts")
	oldAccounts, err := accountIdentifiersFromSlice(expandStringList(o.(*schema.Set).List()))
	if err != nil {
		return diag.FromErr(err)
	}
	newAccounts, err := accountIdentifiersFromSlice(expandStringList(n.(*schema.Set).List()))
	if err != nil {
		return diag.FromErr(err)
	}
	accountsToAdd := shareAccountsDifference(newAccounts, oldAccounts)
	accountsToRemove := shareAccountsDifference(oldAccounts, newAccounts)

	client := meta.(*provider.Context).Client
	if len(accountsToRemove) > 0 {
		err := client.Shares.Alter(ctx, sdk.NewAlterShareRequest(id).WithRemove(sdk.ShareRemoveRequest{Accounts: accountsToRemove}))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error removing accounts from share (%v) err = %w", id.Name(), err))
		}
	}
	if len(accountsToAdd) > 0 {
		share, err := client.Shares.ShowByID(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error reading share (%v) err = %w", id.Name(), err))
		}
		// If a database is already granted to the share, we can add accounts directly.
		// The setShareAccounts workaround (which creates a temp database) is only needed
		// when no database has been granted yet — Snowflake requires at least one database
		// on a share before accounts can be added.
		if share.DatabaseName.Name() != "" {
			err = client.Shares.Alter(ctx, sdk.NewAlterShareRequest(id).WithAdd(sdk.ShareAddRequest{Accounts: accountsToAdd}))
		} else {
			err = setShareAccounts(ctx, client, id, accountsToAdd)
		}
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

// shareAccountsDifference returns the accounts of a that are not in b. The accounts are compared ignoring case, as the schema does.
func shareAccountsDifference(a, b []sdk.AccountIdentifier) []sdk.AccountIdentifier {
	var result []sdk.AccountIdentifier
	for _, account := range a {
		if !slices.ContainsFunc(b, func(other sdk.AccountIdentifier) bool { return strings.EqualFold(account.Name(), other.Name()) }) {
			result = append(result, account)
		}
	}
	return result
}

// =============================================================================
// Pre-existing mappers
// =============================================================================

// TODO [SNOW-4428377]: Verify if this workaround is still needed. Check if Snowflake still requires a database granted to the share before
// accounts can be added, and if not, add the accounts directly (and remove the temporary database, the grants and the revokes).
func setShareAccounts(ctx context.Context, client *sdk.Client, shareID sdk.AccountObjectIdentifier, accounts []sdk.AccountIdentifier) error {
	// There is nothing to add, so the temporary database is not needed.
	if len(accounts) == 0 {
		return nil
	}

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

func accountIdentifiersFromSlice(accounts []string) ([]sdk.AccountIdentifier, error) {
	return collections.MapErr(accounts, sdk.ParseAccountIdentifier)
}
