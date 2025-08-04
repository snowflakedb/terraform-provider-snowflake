package sdk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"
)

func SweepAfterIntegrationTests(client *Client, suffix string) error {
	return sweep(client, suffix)
}

func SweepAfterAcceptanceTests(client *Client, suffix string) error {
	return sweep(client, suffix)
}

// TODO [SNOW-867247]: move this to test code
// TODO [SNOW-867247]: use if exists/use method from helper for dropping
// TODO [SNOW-867247]: sweep all missing account-level objects (like users, integrations, replication groups, network policies, ...)
// TODO [SNOW-867247]: extract sweepers to a separate dir
// TODO [SNOW-867247]: rework the sweepers (funcs -> objects)
// TODO [SNOW-867247]: consider generalization (almost all the sweepers follow the same pattern: show, drop if matches)
// TODO [SNOW-867247]: consider failing after all sweepers and not with the first error
// TODO [SNOW-867247]: consider showing only objects with the given suffix (in almost every sweeper)
func sweep(client *Client, suffix string) error {
	if suffix == "" {
		return fmt.Errorf("suffix is required to run sweepers")
	}
	sweepers := []func() error{
		getAccountPolicyAttachmentsSweeper(client),
		getResourceMonitorSweeper(client, suffix),
		getNetworkPolicySweeper(client, suffix),
		nukeUsers(client, suffix),
		getFailoverGroupSweeper(client, suffix),
		getShareSweeper(client, suffix),
		getDatabaseSweeper(client, suffix),
		getWarehouseSweeper(client, suffix),
		getRoleSweeper(client, suffix),
	}
	for _, sweeper := range sweepers {
		if err := sweeper(); err != nil {
			return err
		}
	}
	return nil
}

func getAccountPolicyAttachmentsSweeper(client *Client) func() error {
	return func() error {
		log.Printf("[DEBUG] Unsetting password and session policies set on the account level")
		ctx := context.Background()
		_ = client.Accounts.UnsetPolicySafely(ctx, PolicyKindPasswordPolicy)
		_ = client.Accounts.UnsetPolicySafely(ctx, PolicyKindSessionPolicy)
		return nil
	}
}

func getResourceMonitorSweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping resource monitors with suffix %s", suffix)
		ctx := context.Background()

		rms, err := client.ResourceMonitors.Show(ctx, nil)
		if err != nil {
			return fmt.Errorf("sweeping resource monitor ended with error, err = %w", err)
		}
		for _, rm := range rms {
			if strings.HasSuffix(rm.Name, suffix) {
				log.Printf("[DEBUG] Dropping resource monitor %s", rm.ID().FullyQualifiedName())
				if err := client.ResourceMonitors.Drop(ctx, rm.ID(), &DropResourceMonitorOptions{IfExists: Bool(true)}); err != nil {
					return fmt.Errorf("sweeping resource monitor %s ended with error, err = %w", rm.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping resource monitor %s", rm.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}

// getNetworkPolicySweeper was introduced to make sure that network policies created during tests are cleaned up.
// It's required as network policies that have connections to the network rules within databases, block their deletion.
// In Snowflake, the network policies can be removed without unsetting network rules, but the network rules cannot be removed without unsetting network policies.
func getNetworkPolicySweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping network policies with suffix %s", suffix)
		ctx := context.Background()

		nps, err := client.NetworkPolicies.Show(ctx, NewShowNetworkPolicyRequest())
		if err != nil {
			return fmt.Errorf("SHOW NETWORK POLICIES ended with error, err = %w", err)
		}

		for _, np := range nps {
			if strings.HasSuffix(np.Name, suffix) && strings.ToUpper(np.Name) != "RESTRICTED_ACCESS" {
				log.Printf("[DEBUG] Dropping network policy %s", np.ID().FullyQualifiedName())
				if err := client.NetworkPolicies.Drop(ctx, NewDropNetworkPolicyRequest(np.ID()).WithIfExists(true)); err != nil {
					return fmt.Errorf("DROP NETWORK POLICY for %s, ended with error, err = %w", np.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping network policy %s", np.ID().FullyQualifiedName())
			}
		}

		return nil
	}
}

func getFailoverGroupSweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping failover groups with suffix %s", suffix)
		ctx := context.Background()

		currentAccount, err := client.ContextFunctions.CurrentAccount(ctx)
		if err != nil {
			return fmt.Errorf("sweeping failover groups ended with error, err = %w", err)
		}
		opts := &ShowFailoverGroupOptions{
			InAccount: NewAccountIdentifierFromAccountLocator(currentAccount),
		}
		fgs, err := client.FailoverGroups.Show(ctx, opts)
		if err != nil {
			return fmt.Errorf("sweeping failover groups ended with error, err = %w", err)
		}
		for _, fg := range fgs {
			if strings.HasSuffix(fg.Name, suffix) && fg.AccountLocator == currentAccount {
				log.Printf("[DEBUG] Dropping failover group %s", fg.ID().FullyQualifiedName())
				if err := client.FailoverGroups.Drop(ctx, fg.ID(), nil); err != nil {
					return fmt.Errorf("sweeping failover group %s ended with error, err = %w", fg.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping failover group %s", fg.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}

func getShareSweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping shares with suffix %s", suffix)
		ctx := context.Background()

		shares, err := client.Shares.Show(ctx, nil)
		if err != nil {
			return fmt.Errorf("sweeping shares ended with error, err = %w", err)
		}
		for _, share := range shares {
			if share.Kind == ShareKindOutbound && strings.HasSuffix(share.Name.Name(), suffix) {
				log.Printf("[DEBUG] Dropping share %s", share.ID().FullyQualifiedName())
				if err := client.Shares.Drop(ctx, share.ID(), &DropShareOptions{IfExists: Bool(true)}); err != nil {
					return fmt.Errorf("sweeping share %s ended with error, err = %w", share.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping share %s", share.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}

func getDatabaseSweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping databases with suffix %s", suffix)
		ctx := context.Background()

		dbs, err := client.Databases.Show(ctx, nil)
		if err != nil {
			return fmt.Errorf("sweeping databases ended with error, err = %w", err)
		}
		for _, db := range dbs {
			if strings.HasSuffix(db.Name, suffix) && db.Name != "SNOWFLAKE" {
				log.Printf("[DEBUG] Dropping database %s", db.ID().FullyQualifiedName())
				if err := client.Databases.Drop(ctx, db.ID(), nil); err != nil {
					if strings.Contains(err.Error(), "Object found is of type 'APPLICATION', not specified type 'DATABASE'") {
						log.Printf("[DEBUG] Skipping database %s", db.ID().FullyQualifiedName())
					} else {
						return fmt.Errorf("sweeping database %s ended with error, err = %w", db.ID().FullyQualifiedName(), err)
					}
				}
			} else {
				log.Printf("[DEBUG] Skipping database %s", db.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}

func getWarehouseSweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping warehouses with suffix %s", suffix)
		ctx := context.Background()

		whs, err := client.Warehouses.Show(ctx, nil)
		if err != nil {
			return fmt.Errorf("sweeping warehouses ended with error, err = %w", err)
		}
		for _, wh := range whs {
			if strings.HasSuffix(wh.Name, suffix) && wh.Name != "SNOWFLAKE" {
				log.Printf("[DEBUG] Dropping warehouse %s", wh.ID().FullyQualifiedName())
				if err := client.Warehouses.Drop(ctx, wh.ID(), nil); err != nil {
					return fmt.Errorf("sweeping warehouse %s ended with error, err = %w", wh.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping warehouse %s", wh.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}

func getRoleSweeper(client *Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping roles with suffix %s", suffix)
		ctx := context.Background()

		roles, err := client.Roles.Show(ctx, NewShowRoleRequest())
		if err != nil {
			return fmt.Errorf("sweeping roles ended with error, err = %w", err)
		}
		for _, role := range roles {
			if strings.HasSuffix(role.Name, suffix) && !slices.Contains([]string{"ACCOUNTADMIN", "SECURITYADMIN", "SYSADMIN", "ORGADMIN", "USERADMIN", "PUBLIC", "PENTESTING_ROLE"}, role.Name) {
				log.Printf("[DEBUG] Dropping role %s", role.ID().FullyQualifiedName())
				if err := client.Roles.Drop(ctx, NewDropRoleRequest(role.ID())); err != nil {
					return fmt.Errorf("sweeping role %s ended with error, err = %w", role.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping role %s", role.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}

// TODO [SNOW-867247]: generalize nuke methods (sweepers too)
// TODO [SNOW-1658402]: handle the ownership problem while handling the better role setup for tests
func nukeWarehouses(client *Client, prefix string) func() error {
	protectedWarehouses := []string{
		"SNOWFLAKE",
		"SYSTEM$STREAMLIT_NOTEBOOK_WH",
	}

	return func() error {
		log.Printf("[DEBUG] Nuking warehouses with prefix %s", prefix)
		ctx := context.Background()

		var like *Like = nil
		if prefix != "" {
			like = &Like{Pattern: String(prefix)}
		}

		whs, err := client.Warehouses.Show(ctx, &ShowWarehouseOptions{Like: like})
		if err != nil {
			return fmt.Errorf("sweeping warehouses ended with error, err = %w", err)
		}
		var errs []error
		log.Printf("[DEBUG] Found %d warehouses matching search criteria", len(whs))
		for idx, wh := range whs {
			log.Printf("[DEBUG] Processing warehouse [%d/%d]: %s...", idx+1, len(whs), wh.ID().FullyQualifiedName())
			if !slices.Contains(protectedWarehouses, wh.Name) && wh.CreatedOn.Before(time.Now().Add(-8*time.Hour)) {
				if wh.Owner != "ACCOUNTADMIN" {
					log.Printf("[DEBUG] Granting ownership on warehouse %s, to ACCOUNTADMIN", wh.ID().FullyQualifiedName())
					err := client.Grants.GrantOwnership(
						ctx,
						OwnershipGrantOn{Object: &Object{
							ObjectType: ObjectTypeWarehouse,
							Name:       wh.ID(),
						}},
						OwnershipGrantTo{
							AccountRoleName: Pointer(NewAccountObjectIdentifier("ACCOUNTADMIN")),
						},
						nil,
					)
					if err != nil {
						errs = append(errs, fmt.Errorf("granting ownership on warehouse %s ended with error, err = %w", wh.ID().FullyQualifiedName(), err))
						continue
					}
				}

				log.Printf("[DEBUG] Dropping warehouse %s, created at: %s", wh.ID().FullyQualifiedName(), wh.CreatedOn.String())
				if err := client.Warehouses.Drop(ctx, wh.ID(), &DropWarehouseOptions{IfExists: Bool(true)}); err != nil {
					log.Printf("[DEBUG] Dropping warehouse %s, resulted in error %v", wh.ID().FullyQualifiedName(), err)
					errs = append(errs, fmt.Errorf("sweeping warehouse %s ended with error, err = %w", wh.ID().FullyQualifiedName(), err))
				}
			} else {
				log.Printf("[DEBUG] Skipping warehouse %s, created at: %s", wh.ID().FullyQualifiedName(), wh.CreatedOn.String())
			}
		}
		return errors.Join(errs...)
	}
}

func nukeDatabases(client *Client, prefix string) func() error {
	protectedDatabases := []string{
		"SNOWFLAKE",
		"MFA_ENFORCEMENT_POLICY",
	}

	return func() error {
		log.Printf("[DEBUG] Nuking databases with prefix %s", prefix)
		ctx := context.Background()

		var like *Like = nil
		if prefix != "" {
			like = &Like{Pattern: String(prefix)}
		}
		dbs, err := client.Databases.Show(ctx, &ShowDatabasesOptions{Like: like})
		if err != nil {
			return fmt.Errorf("sweeping databases ended with error, err = %w", err)
		}
		var errs []error
		log.Printf("[DEBUG] Found %d databases matching search criteria", len(dbs))
		for idx, db := range dbs {
			if db.Owner != "ACCOUNTADMIN" {
				log.Printf("[DEBUG] Granting ownership on database %s, to ACCOUNTADMIN", db.ID().FullyQualifiedName())
				err := client.Grants.GrantOwnership(
					ctx,
					OwnershipGrantOn{Object: &Object{
						ObjectType: ObjectTypeDatabase,
						Name:       db.ID(),
					}},
					OwnershipGrantTo{
						AccountRoleName: Pointer(NewAccountObjectIdentifier("ACCOUNTADMIN")),
					},
					nil,
				)
				if err != nil {
					errs = append(errs, fmt.Errorf("granting ownership on database %s ended with error, err = %w", db.ID().FullyQualifiedName(), err))
					continue
				}
			}

			log.Printf("[DEBUG] Processing database [%d/%d]: %s...", idx+1, len(dbs), db.ID().FullyQualifiedName())
			if !slices.Contains(protectedDatabases, db.Name) && db.CreatedOn.Before(time.Now().Add(-8*time.Hour)) {
				log.Printf("[DEBUG] Dropping database %s, created at: %s", db.ID().FullyQualifiedName(), db.CreatedOn.String())
				if err := client.Databases.Drop(ctx, db.ID(), &DropDatabaseOptions{IfExists: Bool(true)}); err != nil {
					log.Printf("[DEBUG] Dropping database %s, resulted in error %v", db.ID().FullyQualifiedName(), err)
					errs = append(errs, fmt.Errorf("sweeping database %s ended with error, err = %w", db.ID().FullyQualifiedName(), err))
				}
			} else {
				log.Printf("[DEBUG] Skipping database %s, created at: %s", db.ID().FullyQualifiedName(), db.CreatedOn.String())
			}
		}
		return errors.Join(errs...)
	}
}

func nukeUsers(client *Client, suffix string) func() error {
	protectedUsers := []string{
		"SNOWFLAKE",
		"ARTUR_SAWICKI",
		"ARTUR_SAWICKI_LEGACY",
		"JAKUB_MICHALAK",
		"JAKUB_MICHALAK_LEGACY",
		"JAN_CIESLAK",
		"JAN_CIESLAK_LEGACY",
		"TERRAFORM_SVC_ACCOUNT",
		"TEST_CI_SERVICE_USER",
		"PENTESTING_USER_1",
		"PENTESTING_USER_2",
	}

	return func() error {
		ctx := context.Background()

		var userDropCondition func(u User) bool
		if suffix != "" {
			log.Printf("[DEBUG] Sweeping users with suffix %s", suffix)
			userDropCondition = func(u User) bool {
				return strings.HasSuffix(u.Name, suffix) && !slices.Contains(protectedUsers, u.ID().Name())
			}
		} else {
			log.Println("[DEBUG] Sweeping stale users")
			userDropCondition = func(u User) bool {
				return !slices.Contains(protectedUsers, u.Name) && u.CreatedOn.Before(time.Now().Add(-15*time.Minute))
			}
		}

		urs, err := client.Users.Show(ctx, new(ShowUserOptions))
		if err != nil {
			return fmt.Errorf("SHOW USERS ended with error, err = %w", err)
		}

		log.Printf("[DEBUG] Found %d users", len(urs))

		var errs []error
		for idx, user := range urs {
			log.Printf("[DEBUG] Processing user [%d/%d]: %s...", idx+1, len(urs), user.ID().FullyQualifiedName())

			if userDropCondition(user) {
				log.Printf("[DEBUG] Dropping user %s", user.ID().FullyQualifiedName())
				if err := client.Users.Drop(ctx, user.ID(), &DropUserOptions{IfExists: Bool(true)}); err != nil {
					errs = append(errs, fmt.Errorf("sweeping user %s ended with error, err = %w", user.ID().FullyQualifiedName(), err))
				}
			} else {
				log.Printf("[DEBUG] Skipping user %s", user.ID().FullyQualifiedName())
			}
		}

		return errors.Join(errs...)
	}
}
