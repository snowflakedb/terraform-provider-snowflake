package sdk_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random/acceptancetests"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random/integrationtests"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

// TODO [SNOW-867247]: move the sweepers outside of the sdk (and sdk_test) package
// TODO [SNOW-867247]: use test client helpers in sweepers?
func TestSweepAll(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableSweep)
	testenvs.AssertEnvSet(t, string(testenvs.TestObjectsSuffix))

	t.Run("sweep after tests", func(t *testing.T) {
		client := defaultTestClient(t)
		secondaryClient := secondaryTestClient(t)

		err := SweepAfterIntegrationTests(client, integrationtests.ObjectsSuffix)
		assert.NoError(t, err)

		err = SweepAfterIntegrationTests(secondaryClient, integrationtests.ObjectsSuffix)
		assert.NoError(t, err)

		err = SweepAfterAcceptanceTests(client, acceptancetests.ObjectsSuffix)
		assert.NoError(t, err)

		err = SweepAfterAcceptanceTests(secondaryClient, acceptancetests.ObjectsSuffix)
		assert.NoError(t, err)
	})
}

func SweepAfterIntegrationTests(client *sdk.Client, suffix string) error {
	return sweep(client, suffix)
}

func SweepAfterAcceptanceTests(client *sdk.Client, suffix string) error {
	return sweep(client, suffix)
}

// TODO [SNOW-867247]: use if exists/use method from helper for dropping
// TODO [SNOW-867247]: sweep all missing account-level objects (like users, integrations, replication groups, network policies, ...)
// TODO [SNOW-867247]: extract sweepers to a separate dir
// TODO [SNOW-867247]: rework the sweepers (funcs -> objects)
// TODO [SNOW-867247]: consider generalization (almost all the sweepers follow the same pattern: show, drop if matches)
// TODO [SNOW-867247]: consider failing after all sweepers and not with the first error
// TODO [SNOW-867247]: consider showing only objects with the given suffix (in almost every sweeper)
func sweep(client *sdk.Client, suffix string) error {
	if suffix == "" {
		return fmt.Errorf("suffix is required to run sweepers")
	}
	sweepers := []func() error{
		nukeSecurityIntegrations(client, suffix),
		getAccountPolicyAttachmentsSweeper(client),
		nukeResourceMonitors(client, suffix),
		nukeNetworkPolicies(client, suffix),
		nukeUsers(client, suffix),
		nukeFailoverGroups(client, suffix),
		nukeShares(client, suffix),
		nukeDatabases(client, "", suffix),
		nukeWarehouses(client, "", suffix),
		nukeRoles(client, suffix),
	}
	for _, sweeper := range sweepers {
		if err := sweeper(); err != nil {
			return err
		}
	}
	return nil
}

func Test_Sweeper_NukeStaleObjects(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableSweep)

	client := defaultTestClient(t)
	secondaryClient := secondaryTestClient(t)
	thirdClient := thirdTestClient(t)
	fourthClient := fourthTestClient(t)

	allClients := []*sdk.Client{client, secondaryClient, thirdClient, fourthClient}

	// can't use extracted IntegrationTestPrefix and AcceptanceTestPrefix until sweepers reside in the SDK package (cyclic)
	const integrationTestPrefix = "int_test_"
	const acceptanceTestPrefix = "acc_test_"

	t.Run("sweep integration test precreated objects", func(t *testing.T) {
		integrationTestWarehousesPrefix := fmt.Sprintf("%swh_%%", integrationTestPrefix)
		integrationTestDatabasesPrefix := fmt.Sprintf("%sdb_%%", integrationTestPrefix)

		for _, c := range allClients {
			err := nukeWarehouses(c, integrationTestWarehousesPrefix, "")()
			assert.NoError(t, err)

			err = nukeDatabases(c, integrationTestDatabasesPrefix, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep acceptance tests precreated objects", func(t *testing.T) {
		acceptanceTestWarehousesPrefix := fmt.Sprintf("%swh_%%", acceptanceTestPrefix)
		acceptanceTestDatabasesPrefix := fmt.Sprintf("%sdb_%%", acceptanceTestPrefix)

		for _, c := range allClients {
			err := nukeWarehouses(c, acceptanceTestWarehousesPrefix, "")()
			assert.NoError(t, err)

			err = nukeDatabases(c, acceptanceTestDatabasesPrefix, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep security integrations", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeSecurityIntegrations(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep resource monitors", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeResourceMonitors(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep network policies", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeNetworkPolicies(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep users", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeUsers(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep failover groups", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeFailoverGroups(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep roles", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeRoles(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep shares", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeShares(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep databases", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeDatabases(c, "", "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep warehouses", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeWarehouses(c, "", "")()
			assert.NoError(t, err)
		}
	})

	// TODO [SNOW-867247]: nuke stale objects (e.g. created more than 2 weeks ago)

	// TODO [SNOW-867247]: nuke external oauth integrations because of errors like
	// Error: 003524 (22023): SQL execution error: An integration with the given issuer already exists for this account
}

// TODO [SNOW-867247]: longer time for now; validate the timezone behavior during sweepers rework
var stalePeriod = -12 * time.Hour

// TODO [SNOW-867247]: generalize nuke methods (sweepers too)
// TODO [SNOW-1658402]: handle the ownership problem while handling the better role setup for tests
func nukeWarehouses(client *sdk.Client, prefix string, suffix string) func() error {
	protectedWarehouses := []string{
		"SNOWFLAKE",
		"SYSTEM$STREAMLIT_NOTEBOOK_WH",
	}

	return func() error {
		ctx := context.Background()

		var whDropCondition func(wh sdk.Warehouse) bool
		switch {
		case prefix != "":
			log.Printf("[DEBUG] Sweeping warehouses with prefix %s", prefix)
			whDropCondition = func(wh sdk.Warehouse) bool {
				return strings.HasPrefix(wh.Name, prefix)
			}
		case suffix != "":
			log.Printf("[DEBUG] Sweeping warehouses with suffix %s", suffix)
			whDropCondition = func(wh sdk.Warehouse) bool {
				return strings.HasSuffix(wh.Name, suffix)
			}
		default:
			log.Println("[DEBUG] Sweeping stale warehouses")
			whDropCondition = func(wh sdk.Warehouse) bool {
				return wh.CreatedOn.Before(time.Now().Add(stalePeriod))
			}
		}

		whs, err := client.Warehouses.Show(ctx, new(sdk.ShowWarehouseOptions))
		if err != nil {
			return fmt.Errorf("SHOW WAREHOUSES ended with error, err = %w", err)
		}

		log.Printf("[DEBUG] Found %d warehouses", len(whs))

		var errs []error
		for idx, wh := range whs {
			log.Printf("[DEBUG] Processing warehouse [%d/%d]: %s...", idx+1, len(whs), wh.ID().FullyQualifiedName())
			if !slices.Contains(protectedWarehouses, wh.Name) && whDropCondition(wh) {
				if wh.Owner != snowflakeroles.Accountadmin.Name() {
					log.Printf("[DEBUG] Granting ownership on warehouse %s, to ACCOUNTADMIN", wh.ID().FullyQualifiedName())
					err := client.Grants.GrantOwnership(
						ctx,
						sdk.OwnershipGrantOn{Object: &sdk.Object{
							ObjectType: sdk.ObjectTypeWarehouse,
							Name:       wh.ID(),
						}},
						sdk.OwnershipGrantTo{
							AccountRoleName: sdk.Pointer(snowflakeroles.Accountadmin),
						},
						nil,
					)
					if err != nil {
						errs = append(errs, fmt.Errorf("granting ownership on warehouse %s ended with error, err = %w", wh.ID().FullyQualifiedName(), err))
						continue
					}
				}

				log.Printf("[DEBUG] Dropping warehouse %s, created at: %s", wh.ID().FullyQualifiedName(), wh.CreatedOn.String())
				// to handle identifiers with containing `"` - we do not escape them currently in the SDK SQL generation
				whId := wh.ID()
				if strings.Contains(whId.Name(), `"`) {
					whId = sdk.NewAccountObjectIdentifier(strings.ReplaceAll(whId.Name(), `"`, `""`))
				}
				if err := client.Warehouses.DropSafely(ctx, whId); err != nil {
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

func nukeDatabases(client *sdk.Client, prefix string, suffix string) func() error {
	protectedDatabases := []string{
		"SNOWFLAKE",
		"MFA_ENFORCEMENT_POLICY",
		"TERRAFORM_TEST_SETUP_OBJECTS",
		"TEST_RESULTS_DATABASE",
	}

	return func() error {
		ctx := context.Background()

		var dbDropCondition func(db sdk.Database) bool
		switch {
		case prefix != "":
			log.Printf("[DEBUG] Sweeping databases with prefix %s", prefix)
			dbDropCondition = func(db sdk.Database) bool {
				return strings.HasPrefix(db.Name, prefix)
			}
		case suffix != "":
			log.Printf("[DEBUG] Sweeping databases with suffix %s", suffix)
			dbDropCondition = func(db sdk.Database) bool {
				return strings.HasSuffix(db.Name, suffix)
			}
		default:
			log.Println("[DEBUG] Sweeping stale databases")
			dbDropCondition = func(db sdk.Database) bool {
				return db.CreatedOn.Before(time.Now().Add(stalePeriod))
			}
		}

		dbs, err := client.Databases.Show(ctx, new(sdk.ShowDatabasesOptions))
		if err != nil {
			return fmt.Errorf("SHOW DATABASES ended with error, err = %w", err)
		}

		log.Printf("[DEBUG] Found %d databases", len(dbs))

		var errs []error
		for idx, db := range dbs {
			log.Printf("[DEBUG] Processing database [%d/%d]: %s...", idx+1, len(dbs), db.ID().FullyQualifiedName())
			if !slices.Contains(protectedDatabases, db.Name) && dbDropCondition(db) {
				if db.Owner != snowflakeroles.Accountadmin.Name() {
					log.Printf("[DEBUG] Granting ownership on database %s, to ACCOUNTADMIN", db.ID().FullyQualifiedName())
					err := client.Grants.GrantOwnership(
						ctx,
						sdk.OwnershipGrantOn{Object: &sdk.Object{
							ObjectType: sdk.ObjectTypeDatabase,
							Name:       db.ID(),
						}},
						sdk.OwnershipGrantTo{
							AccountRoleName: sdk.Pointer(snowflakeroles.Accountadmin),
						},
						nil,
					)
					if err != nil {
						errs = append(errs, fmt.Errorf("granting ownership on database %s ended with error, err = %w", db.ID().FullyQualifiedName(), err))
						continue
					}
				}

				log.Printf("[DEBUG] Dropping database %s, created at: %s", db.ID().FullyQualifiedName(), db.CreatedOn.String())
				if err := client.Databases.DropSafely(ctx, db.ID()); err != nil {
					if strings.Contains(err.Error(), "Object found is of type 'APPLICATION', not specified type 'DATABASE'") {
						log.Printf("[DEBUG] Skipping database %s as it's an application, err: %v", db.ID().FullyQualifiedName(), err)
						continue
					}
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

func nukeUsers(client *sdk.Client, suffix string) func() error {
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

		var userDropCondition func(u sdk.User) bool
		if suffix != "" {
			log.Printf("[DEBUG] Sweeping users with suffix %s", suffix)
			userDropCondition = func(u sdk.User) bool {
				return strings.HasSuffix(u.Name, suffix)
			}
		} else {
			log.Println("[DEBUG] Sweeping stale users")
			userDropCondition = func(u sdk.User) bool {
				return u.CreatedOn.Before(time.Now().Add(-15 * time.Minute))
			}
		}

		urs, err := client.Users.Show(ctx, new(sdk.ShowUserOptions))
		if err != nil {
			return fmt.Errorf("SHOW USERS ended with error, err = %w", err)
		}

		log.Printf("[DEBUG] Found %d users", len(urs))

		var errs []error
		for idx, user := range urs {
			log.Printf("[DEBUG] Processing user [%d/%d]: %s...", idx+1, len(urs), user.ID().FullyQualifiedName())

			if !slices.Contains(protectedUsers, user.Name) && userDropCondition(user) {
				log.Printf("[DEBUG] Dropping user %s", user.ID().FullyQualifiedName())
				if err := client.Users.Drop(ctx, user.ID(), &sdk.DropUserOptions{IfExists: sdk.Bool(true)}); err != nil {
					errs = append(errs, fmt.Errorf("sweeping user %s ended with error, err = %w", user.ID().FullyQualifiedName(), err))
				}
			} else {
				log.Printf("[DEBUG] Skipping user %s", user.ID().FullyQualifiedName())
			}
		}

		return errors.Join(errs...)
	}
}

func nukeSecurityIntegrations(client *sdk.Client, suffix string) func() error {
	return func() error {
		ctx := context.Background()

		var integrationDropCondition func(u sdk.SecurityIntegration) bool
		if suffix != "" {
			log.Printf("[DEBUG] Sweeping security integrations with suffix %s", suffix)
			integrationDropCondition = func(u sdk.SecurityIntegration) bool {
				return strings.HasSuffix(u.Name, suffix)
			}
		} else {
			log.Println("[DEBUG] Sweeping stale security integrations")
			integrationDropCondition = func(u sdk.SecurityIntegration) bool {
				return u.CreatedOn.Before(time.Now().Add(-15 * time.Minute))
			}
		}
		integrations, err := client.SecurityIntegrations.Show(ctx, sdk.NewShowSecurityIntegrationRequest())
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Found %d security integrations", len(integrations))

		var errs []error
		for idx, integration := range integrations {
			log.Printf("[DEBUG] Processing security integration [%d/%d]: %s...", idx+1, len(integrations), integration.Name)

			if integrationDropCondition(integration) {
				log.Printf("[DEBUG] Dropping security integration %s", integration.Name)
				if err := client.SecurityIntegrations.DropSafely(ctx, integration.ID()); err != nil {
					errs = append(errs, fmt.Errorf("sweeping security integration %s ended with an error, err = %w", integration.Name, err))
				}
			} else {
				log.Printf("[DEBUG] Skipping security integration %s", integration.Name)
			}
		}

		return errors.Join(errs...)
	}
}

func nukeRoles(client *sdk.Client, suffix string) func() error {
	protectedRoles := []sdk.AccountObjectIdentifier{
		snowflakeroles.GlobalOrgAdmin,
		snowflakeroles.Orgadmin,
		snowflakeroles.Accountadmin,
		snowflakeroles.SecurityAdmin,
		snowflakeroles.SysAdmin,
		snowflakeroles.UserAdmin,
		snowflakeroles.Public,
		snowflakeroles.PentestingRole,
		snowflakeroles.OktaProvisioner,
		snowflakeroles.AadProvisioner,
		snowflakeroles.GenericScimProvisioner,
	}

	return func() error {
		ctx := context.Background()

		var roleDropCondition func(r sdk.Role) bool
		if suffix != "" {
			log.Printf("[DEBUG] Sweeping roles with suffix %s", suffix)
			roleDropCondition = func(r sdk.Role) bool {
				return strings.HasSuffix(r.Name, suffix)
			}
		} else {
			log.Println("[DEBUG] Sweeping stale roles")
			roleDropCondition = func(r sdk.Role) bool {
				return r.CreatedOn.Before(time.Now().Add(-15 * time.Minute))
			}
		}

		rs, err := client.Roles.Show(ctx, sdk.NewShowRoleRequest())
		if err != nil {
			return fmt.Errorf("SHOW ROLES ended with error, err = %w", err)
		}

		log.Printf("[DEBUG] Found %d roles", len(rs))

		var errs []error
		for idx, role := range rs {
			log.Printf("[DEBUG] Processing role [%d/%d]: %s...", idx+1, len(rs), role.ID().FullyQualifiedName())

			if !slices.Contains(protectedRoles, role.ID()) && roleDropCondition(role) {
				if role.Owner != snowflakeroles.Accountadmin.Name() {
					log.Printf("[DEBUG] Granting ownership on role %s, to ACCOUNTADMIN", role.ID().FullyQualifiedName())
					err := client.Grants.GrantOwnership(
						ctx,
						sdk.OwnershipGrantOn{Object: &sdk.Object{
							ObjectType: sdk.ObjectTypeRole,
							Name:       role.ID(),
						}},
						sdk.OwnershipGrantTo{
							AccountRoleName: sdk.Pointer(snowflakeroles.Accountadmin),
						},
						nil,
					)
					if err != nil {
						errs = append(errs, fmt.Errorf("granting ownership on role %s ended with error, err = %w", role.ID().FullyQualifiedName(), err))
						continue
					}
				}

				log.Printf("[DEBUG] Dropping role %s", role.ID().FullyQualifiedName())
				if err := client.Roles.DropSafely(ctx, role.ID()); err != nil {
					errs = append(errs, fmt.Errorf("sweeping role %s ended with error, err = %w", role.ID().FullyQualifiedName(), err))
				}
			} else {
				log.Printf("[DEBUG] Skipping role %s", role.ID().FullyQualifiedName())
			}
		}

		return errors.Join(errs...)
	}
}

func nukeShares(client *sdk.Client, suffix string) func() error {
	protectedShares := []string{
		// this one is INBOUND but putting it here either way
		"ACCOUNT_USAGE",
	}

	return func() error {
		ctx := context.Background()

		var shareDropCondition func(s sdk.Share) bool
		if suffix != "" {
			log.Printf("[DEBUG] Sweeping shares with suffix %s", suffix)
			shareDropCondition = func(s sdk.Share) bool {
				return strings.HasSuffix(s.Name.Name(), suffix)
			}
		} else {
			log.Println("[DEBUG] Sweeping stale shares")
			shareDropCondition = func(s sdk.Share) bool {
				return s.CreatedOn.Before(time.Now().Add(stalePeriod))
			}
		}

		shares, err := client.Shares.Show(ctx, new(sdk.ShowShareOptions))
		if err != nil {
			return fmt.Errorf("SHOW SHARES ended with error, err = %w", err)
		}

		log.Printf("[DEBUG] Found %d shares", len(shares))

		var errs []error
		for idx, share := range shares {
			log.Printf("[DEBUG] Processing share [%d/%d]: %s...", idx+1, len(shares), share.ID().FullyQualifiedName())

			if !slices.Contains(protectedShares, share.Name.Name()) && shareDropCondition(share) && share.Kind == sdk.ShareKindOutbound {
				log.Printf("[DEBUG] Dropping share %s", share.ID().FullyQualifiedName())
				if err := client.Shares.DropSafely(ctx, share.ID()); err != nil {
					errs = append(errs, fmt.Errorf("sweeping share %s ended with error, err = %w", share.ID().FullyQualifiedName(), err))
				}
			} else {
				log.Printf("[DEBUG] Skipping share %s", share.ID().FullyQualifiedName())
			}
		}
		return errors.Join(errs...)
	}
}

// nukeNetworkPolicies was introduced to make sure that network policies created during tests are cleaned up.
// It's required as network policies that have connections to the network rules within databases, block their deletion.
// In Snowflake, the network policies can be removed without unsetting network rules, but the network rules cannot be removed without unsetting network policies.
func nukeNetworkPolicies(client *sdk.Client, suffix string) func() error {
	protectedNetworkPolicies := []string{
		"RESTRICTED_ACCESS",
	}

	return func() error {
		ctx := context.Background()

		var networkPolicyDropCondition func(n sdk.NetworkPolicy) bool
		if suffix != "" {
			log.Printf("[DEBUG] Sweeping network policies with suffix %s", suffix)
			networkPolicyDropCondition = func(n sdk.NetworkPolicy) bool {
				return strings.HasSuffix(n.Name, suffix)
			}
		} else {
			log.Println("[DEBUG] Sweeping stale network policies")
			networkPolicyDropCondition = func(n sdk.NetworkPolicy) bool {
				// CreatedOn in network policy is string and not time
				createdOn, err := time.Parse(time.RFC3339, n.CreatedOn)
				if err != nil {
					log.Printf("[DEBUG] Could not parse created on: '%s' for network policy %s", n.CreatedOn, n.ID().FullyQualifiedName())
					return false
				}
				return createdOn.Before(time.Now().Add(stalePeriod))
			}
		}

		nps, err := client.NetworkPolicies.Show(ctx, sdk.NewShowNetworkPolicyRequest())
		if err != nil {
			return fmt.Errorf("SHOW NETWORK POLICIES ended with error, err = %w", err)
		}

		log.Printf("[DEBUG] Found %d network policies", len(nps))

		var errs []error
		for idx, np := range nps {
			log.Printf("[DEBUG] Processing network policy [%d/%d]: %s...", idx+1, len(nps), np.ID().FullyQualifiedName())
			if !slices.Contains(protectedNetworkPolicies, strings.ToUpper(np.Name)) && networkPolicyDropCondition(np) {
				log.Printf("[DEBUG] Dropping network policy %s", np.ID().FullyQualifiedName())
				if err := client.NetworkPolicies.DropSafely(ctx, np.ID()); err != nil {
					errs = append(errs, fmt.Errorf("sweeping network policy %s ended with error, err = %w", np.ID().FullyQualifiedName(), err))
				}
			} else {
				log.Printf("[DEBUG] Skipping network policy %s", np.ID().FullyQualifiedName())
			}
		}

		return errors.Join(errs...)
	}
}

func nukeResourceMonitors(client *sdk.Client, suffix string) func() error {
	var protectedResourceMonitors []string

	return func() error {
		ctx := context.Background()

		var rmDropCondition func(rm sdk.ResourceMonitor) bool
		switch {
		case suffix != "":
			log.Printf("[DEBUG] Sweeping resource monitors with suffix %s", suffix)
			rmDropCondition = func(rm sdk.ResourceMonitor) bool {
				return strings.HasSuffix(rm.Name, suffix)
			}
		default:
			log.Println("[DEBUG] Sweeping stale resource monitors")
			rmDropCondition = func(rm sdk.ResourceMonitor) bool {
				return rm.CreatedOn.Before(time.Now().Add(stalePeriod))
			}
		}

		rms, err := client.ResourceMonitors.Show(ctx, new(sdk.ShowResourceMonitorOptions))
		if err != nil {
			return fmt.Errorf("SHOW RESOURCE MONITORS ended with error, err = %w", err)
		}

		log.Printf("[DEBUG] Found %d resource monitors", len(rms))

		var errs []error
		for idx, rm := range rms {
			log.Printf("[DEBUG] Processing resurce monitor [%d/%d]: %s...", idx+1, len(rms), rm.ID().FullyQualifiedName())

			if !slices.Contains(protectedResourceMonitors, rm.Name) && rmDropCondition(rm) {
				if rm.Owner != snowflakeroles.Accountadmin.Name() {
					log.Printf("[DEBUG] Granting ownership on resource monitor %s, to ACCOUNTADMIN", rm.ID().FullyQualifiedName())
					err := client.Grants.GrantOwnership(
						ctx,
						sdk.OwnershipGrantOn{Object: &sdk.Object{
							ObjectType: sdk.ObjectTypeResourceMonitor,
							Name:       rm.ID(),
						}},
						sdk.OwnershipGrantTo{
							AccountRoleName: sdk.Pointer(snowflakeroles.Accountadmin),
						},
						nil,
					)
					if err != nil {
						errs = append(errs, fmt.Errorf("granting ownership on resource monitor %s ended with error, err = %w", rm.ID().FullyQualifiedName(), err))
						continue
					}
				}

				log.Printf("[DEBUG] Dropping resource monitor %s, created at: %s", rm.ID().FullyQualifiedName(), rm.CreatedOn.String())
				// to handle identifiers with containing `"` - we do not escape them currently in the SDK SQL generation
				rmId := rm.ID()
				if strings.Contains(rmId.Name(), `"`) {
					rmId = sdk.NewAccountObjectIdentifier(strings.ReplaceAll(rmId.Name(), `"`, `""`))
				}
				if err := client.ResourceMonitors.DropSafely(ctx, rmId); err != nil {
					log.Printf("[DEBUG] Dropping resource monitor %s, resulted in error %v", rm.ID().FullyQualifiedName(), err)
					errs = append(errs, fmt.Errorf("sweeping resource monitor %s ended with error, err = %w", rm.ID().FullyQualifiedName(), err))
				}
			} else {
				log.Printf("[DEBUG] Skipping resource monitor %s, created at: %s", rm.ID().FullyQualifiedName(), rm.CreatedOn.String())
			}
		}
		return errors.Join(errs...)
	}
}

func nukeFailoverGroups(client *sdk.Client, suffix string) func() error {
	var protectedFailoverGroups []string

	return func() error {
		ctx := context.Background()

		var fgDropCondition func(fg sdk.FailoverGroup) bool
		switch {
		case suffix != "":
			log.Printf("[DEBUG] Sweeping failover groups with suffix %s", suffix)
			fgDropCondition = func(fg sdk.FailoverGroup) bool {
				return strings.HasSuffix(fg.Name, suffix)
			}
		default:
			log.Println("[DEBUG] Sweeping stale failover groups")
			fgDropCondition = func(fg sdk.FailoverGroup) bool {
				return fg.CreatedOn.Before(time.Now().Add(stalePeriod))
			}
		}

		accountLocator := client.GetAccountLocator()
		opts := &sdk.ShowFailoverGroupOptions{
			InAccount: sdk.NewAccountIdentifierFromAccountLocator(accountLocator),
		}

		fgs, err := client.FailoverGroups.Show(ctx, opts)
		if err != nil {
			return fmt.Errorf("SHOW FAILOVER GROUPS ended with error, err = %w", err)
		}

		log.Printf("[DEBUG] Found %d failover groups", len(fgs))

		var errs []error
		for idx, fg := range fgs {
			log.Printf("[DEBUG] Processing failover group [%d/%d]: %s...", idx+1, len(fgs), fg.ID().FullyQualifiedName())

			if fg.Owner != snowflakeroles.Accountadmin.Name() {
				log.Printf("[DEBUG] Granting ownership on failover group %s, to ACCOUNTADMIN", fg.ID().FullyQualifiedName())
				err := client.Grants.GrantOwnership(
					ctx,
					sdk.OwnershipGrantOn{Object: &sdk.Object{
						ObjectType: sdk.ObjectTypeFailoverGroup,
						Name:       fg.ID(),
					}},
					sdk.OwnershipGrantTo{
						AccountRoleName: sdk.Pointer(snowflakeroles.Accountadmin),
					},
					nil,
				)
				if err != nil {
					errs = append(errs, fmt.Errorf("granting ownership on failover group %s ended with error, err = %w", fg.ID().FullyQualifiedName(), err))
					continue
				}
			}

			if !slices.Contains(protectedFailoverGroups, fg.Name) && fgDropCondition(fg) && fg.AccountLocator == accountLocator {
				log.Printf("[DEBUG] Dropping failover group %s, created at: %s", fg.ID().FullyQualifiedName(), fg.CreatedOn.String())
				if err := client.FailoverGroups.DropSafely(ctx, fg.ID()); err != nil {
					log.Printf("[DEBUG] Dropping failover group %s, resulted in error %v", fg.ID().FullyQualifiedName(), err)
					errs = append(errs, fmt.Errorf("sweeping failover group %s ended with error, err = %w", fg.ID().FullyQualifiedName(), err))
				}
			} else {
				log.Printf("[DEBUG] Skipping failover group %s, created at: %s", fg.ID().FullyQualifiedName(), fg.CreatedOn.String())
			}
		}
		return errors.Join(errs...)
	}
}

func defaultTestClient(t *testing.T) *sdk.Client {
	t.Helper()
	return testClient(t, testprofiles.Default)
}

func secondaryTestClient(t *testing.T) *sdk.Client {
	t.Helper()
	return testClient(t, testprofiles.Secondary)
}

func thirdTestClient(t *testing.T) *sdk.Client {
	t.Helper()
	return testClient(t, testprofiles.Third)
}

func fourthTestClient(t *testing.T) *sdk.Client {
	t.Helper()
	return testClient(t, testprofiles.Fourth)
}

func testClient(t *testing.T, profile string) *sdk.Client {
	t.Helper()

	config, err := sdk.ProfileConfig(profile)
	if err != nil {
		t.Skipf("Snowflake %s profile not configured. Must be set in ~/.snowflake/config", profile)
	}
	client, err := sdk.NewClient(config)
	if err != nil {
		t.Skipf("Snowflake %s profile not configured. Must be set in ~/.snowflake/config", profile)
	}

	return client
}
