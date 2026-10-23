package resources

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

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
