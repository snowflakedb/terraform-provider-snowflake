package sdk_test

import (
	"context"
	"encoding/json"
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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

// TODO [SNOW-867247]: move the sweepers outside of the sdk (and sdk_test) package
// TODO [SNOW-867247]: use test client helpers in sweepers?
func TestSweepAll(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableSweep)
	testenvs.AssertEnvSet(t, string(testenvs.TestObjectsSuffix))

	t.Run("sweep after tests", func(t *testing.T) {
		sweepAfterTests := func(client *sdk.Client) {
			assert.NoError(t, SweepAfterIntegrationTests(client, integrationtests.ObjectsSuffix))
			assert.NoError(t, SweepAfterAcceptanceTests(client, acceptancetests.ObjectsSuffix))
		}

		sweepAfterTests(defaultTestClient(t))
		sweepAfterTests(secondaryTestClient(t))
		sweepAfterTests(azureTestClient(t))
		sweepAfterTests(gcpTestClient(t))

		if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakeNonProdEnvironment {
			sweepAfterTests(snowflakeDefaultsTestClient(t))
		}
	})

	t.Run("Send test results", SendTestResults)
}

func SweepAfterIntegrationTests(client *sdk.Client, suffix string) error {
	return sweep(client, suffix)
}

func SweepAfterAcceptanceTests(client *sdk.Client, suffix string) error {
	return sweep(client, suffix)
}

// TODO [SNOW-867247]: sweep all missing account-level objects (like replication groups, connections, compute pools, external volumes, ...)
// TODO [SNOW-867247]: extract sweepers to a separate dir
// TODO [SNOW-867247]: rework the sweepers (funcs -> objects)
// TODO [SNOW-867247]: consider generalization (almost all the sweepers follow the same pattern: show, drop if matches); partially done with nukeAccountObjects
// TODO [SNOW-867247]: consider showing only objects with the given suffix (in almost every sweeper)
func sweep(client *sdk.Client, suffix string) error {
	if suffix == "" {
		return fmt.Errorf("suffix is required to run sweepers")
	}
	sweepers := []func() error{
		nukePostgresInstances(client, suffix),
		nukeSecurityIntegrations(client, suffix),
		getAccountPolicyAttachmentsSweeper(client),
		nukeResourceMonitors(client, suffix),
		nukeNetworkPolicies(client, suffix),
		nukeUsers(client, suffix),
		nukeFailoverGroups(client, suffix),
		nukeShares(client, suffix),
		nukeApplications(client, suffix),
		nukeApplicationPackages(client, suffix),
		nukeListings(client, suffix),
		nukeDatabases(client, "", suffix),
		nukeNotificationIntegrations(client, suffix),
		nukeStorageIntegrations(client, suffix),
		nukeApiIntegrations(client, suffix),
		nukeCatalogIntegrations(client, suffix),
		nukeExternalAccessIntegrations(client, suffix),
		nukeComputePools(client, suffix),
		nukeConnections(client, suffix),
		nukeExternalVolumes(client, suffix),
		nukeWarehouses(client, "", suffix),
		nukeRoles(client, suffix),
	}
	// All the sweepers are run, even if some of them fail; otherwise a single failure would leave the objects handled by the subsequent sweepers behind.
	var errs []error
	for _, sweeper := range sweepers {
		errs = append(errs, sweeper())
	}
	return errors.Join(errs...)
}

func Test_Sweeper_NukeStaleObjects(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableSweep)

	allClients := []*sdk.Client{
		defaultTestClient(t),
		secondaryTestClient(t),
		thirdTestClient(t),
		fourthTestClient(t),
		azureTestClient(t),
		gcpTestClient(t),
	}

	if testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakeNonProdEnvironment {
		allClients = append(
			allClients,
			snowflakeDefaultsTestClient(t),
		)
	}

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

	t.Run("sweep postgres instances", func(t *testing.T) {
		for _, c := range allClients {
			err := nukePostgresInstances(c, "")()
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

	t.Run("sweep applications", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeApplications(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep application packages", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeApplicationPackages(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep listings", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeListings(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep databases", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeDatabases(c, "", "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep notification integrations", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeNotificationIntegrations(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep storage integrations", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeStorageIntegrations(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep api integrations", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeApiIntegrations(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep catalog integrations", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeCatalogIntegrations(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep external access integrations", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeExternalAccessIntegrations(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep compute pools", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeComputePools(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep connections", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeConnections(c, "")()
			assert.NoError(t, err)
		}
	})

	t.Run("sweep external volumes", func(t *testing.T) {
		for _, c := range allClients {
			err := nukeExternalVolumes(c, "")()
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

// Dropping a connection keeps failing until the replication settles; these mirror the polling
// that ConnectionClient.DropFunc does in the tests.
const (
	connectionDropAttempts = 10
	connectionDropInterval = 2 * time.Second
)

// accountObjectSweeperConfig describes an account-level object that follows the usual show -> drop if matches pattern.
type accountObjectSweeperConfig[T any] struct {
	// objectTypeName is a lowercase singular name used in logs and errors, e.g. "notification integration".
	objectTypeName string
	protectedNames []string
	show           func(ctx context.Context) ([]T, error)
	name           func(T) string
	createdOn      func(T) (time.Time, error)
	id             func(T) sdk.AccountObjectIdentifier
	dropSafely     func(ctx context.Context, id sdk.AccountObjectIdentifier) error

	// skip marks the objects that are out of this sweeper's scope for a reason other than their name (e.g. they belong to another account or to a Native App).
	skip func(T) bool

	// stalePeriodOverride is used instead of the package-level stalePeriod when set.
	stalePeriodOverride *time.Duration

	// owner and takeOwnership have to be set together. When they are, the ownership is transferred to ACCOUNTADMIN
	// before dropping the objects owned by any other role.
	// TODO [SNOW-1658402]: handle the ownership problem while handling the better role setup for tests
	owner         func(T) string
	takeOwnership func(ctx context.Context, id sdk.AccountObjectIdentifier) error
}

func (cfg accountObjectSweeperConfig[T]) skips(object T) bool {
	return cfg.skip != nil && cfg.skip(object)
}

func (cfg accountObjectSweeperConfig[T]) effectiveStalePeriod() time.Duration {
	if cfg.stalePeriodOverride != nil {
		return *cfg.stalePeriodOverride
	}
	return stalePeriod
}

// grantOwnershipToAccountadmin builds the takeOwnership function for the objects of the given type.
func grantOwnershipToAccountadmin(client *sdk.Client, objectType sdk.ObjectType) func(ctx context.Context, id sdk.AccountObjectIdentifier) error {
	return func(ctx context.Context, id sdk.AccountObjectIdentifier) error {
		return client.Grants.GrantOwnership(
			ctx,
			sdk.OwnershipGrantOn{Object: &sdk.Object{
				ObjectType: objectType,
				Name:       id,
			}},
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(snowflakeroles.Accountadmin),
			},
			nil,
		)
	}
}

// retryingDropSafely builds a dropSafely function that retries the drop, which is needed for the objects
// that cannot be dropped until an asynchronous operation on Snowflake's side settles.
func retryingDropSafely(
	dropSafely func(ctx context.Context, id sdk.AccountObjectIdentifier) error,
	attempts int,
	interval time.Duration,
) func(ctx context.Context, id sdk.AccountObjectIdentifier) error {
	return func(ctx context.Context, id sdk.AccountObjectIdentifier) error {
		var dropErrs []error
		retryErr := util.Retry(attempts, interval, func() (error, bool) {
			if err := dropSafely(ctx, id); err != nil {
				log.Printf("[DEBUG] Dropping %s failed, err = %v", id.FullyQualifiedName(), err)
				dropErrs = append(dropErrs, err)
				return nil, false
			}
			return nil, true
		})
		if retryErr != nil {
			return errors.Join(append(dropErrs, retryErr)...)
		}
		return nil
	}
}

// externalVolumeCreatedOn reads the external volume's creation time with the undocumented DESCRIBE AS RESOURCE.
// TODO [SNOW-867247]: move this to the test client helpers once the sweepers use them
func externalVolumeCreatedOn(ctx context.Context, client *sdk.Client, id sdk.AccountObjectIdentifier) (time.Time, error) {
	var raw string
	if err := client.QueryOneForTests(ctx, &raw, fmt.Sprintf(`DESCRIBE AS RESOURCE EXTERNAL VOLUME %s`, id.FullyQualifiedName())); err != nil {
		return time.Time{}, err
	}

	// created_on is returned as RFC3339 (e.g. 2024-11-18T13:10:36.721+00:00), so it unmarshals without any custom parsing
	var described struct {
		CreatedOn time.Time `json:"created_on"`
	}
	if err := json.Unmarshal([]byte(raw), &described); err != nil {
		return time.Time{}, err
	}
	return described.CreatedOn, nil
}

// nukeAccountObjects drops the objects matching the given prefix or suffix, or the stale ones when both are empty.
func nukeAccountObjects[T any](prefix string, suffix string, cfg accountObjectSweeperConfig[T]) func() error {
	return func() error {
		ctx := context.Background()

		var dropCondition func(object T) bool
		switch {
		case prefix != "":
			log.Printf("[DEBUG] Sweeping %ss with prefix %s", cfg.objectTypeName, prefix)
			dropCondition = func(object T) bool {
				return strings.HasPrefix(cfg.name(object), prefix)
			}
		case suffix != "":
			log.Printf("[DEBUG] Sweeping %ss with suffix %s", cfg.objectTypeName, suffix)
			dropCondition = func(object T) bool {
				return strings.HasSuffix(cfg.name(object), suffix)
			}
		default:
			log.Printf("[DEBUG] Sweeping stale %ss", cfg.objectTypeName)
			dropCondition = func(object T) bool {
				// this is the only place the creation time is needed, so the objects that have to be queried for it
				// (see externalVolumeCreatedOn) are not queried at all when sweeping by prefix or suffix
				createdOn, err := cfg.createdOn(object)
				if err != nil {
					// without the creation time the staleness can't be told, so the object is left alone
					log.Printf("[DEBUG] Could not read the creation time of %s %s, err = %v", cfg.objectTypeName, cfg.id(object).FullyQualifiedName(), err)
					return false
				}
				log.Printf("[DEBUG] %s %s was created at %s", cfg.objectTypeName, cfg.id(object).FullyQualifiedName(), createdOn.String())
				return createdOn.Before(time.Now().Add(cfg.effectiveStalePeriod()))
			}
		}

		objects, err := cfg.show(ctx)
		if err != nil {
			return fmt.Errorf("showing %ss ended with error, err = %w", cfg.objectTypeName, err)
		}

		log.Printf("[DEBUG] Found %d %ss", len(objects), cfg.objectTypeName)

		var errs []error
		for idx, object := range objects {
			id := cfg.id(object)
			log.Printf("[DEBUG] Processing %s [%d/%d]: %s...", cfg.objectTypeName, idx+1, len(objects), id.FullyQualifiedName())

			if slices.Contains(cfg.protectedNames, cfg.name(object)) || cfg.skips(object) || !dropCondition(object) {
				log.Printf("[DEBUG] Skipping %s %s", cfg.objectTypeName, id.FullyQualifiedName())
				continue
			}

			if cfg.takeOwnership != nil && cfg.owner(object) != snowflakeroles.Accountadmin.Name() {
				log.Printf("[DEBUG] Granting ownership on %s %s, to ACCOUNTADMIN", cfg.objectTypeName, id.FullyQualifiedName())
				if err := cfg.takeOwnership(ctx, id); err != nil {
					errs = append(errs, fmt.Errorf("granting ownership on %s %s ended with error, err = %w", cfg.objectTypeName, id.FullyQualifiedName(), err))
					continue
				}
			}

			log.Printf("[DEBUG] Dropping %s %s", cfg.objectTypeName, id.FullyQualifiedName())
			if err := cfg.dropSafely(ctx, id); err != nil {
				log.Printf("[DEBUG] Dropping %s %s, resulted in error %v", cfg.objectTypeName, id.FullyQualifiedName(), err)
				errs = append(errs, fmt.Errorf("sweeping %s %s ended with error, err = %w", cfg.objectTypeName, id.FullyQualifiedName(), err))
			}
		}

		return errors.Join(errs...)
	}
}

func nukeNotificationIntegrations(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.NotificationIntegration]{
		objectTypeName: "notification integration",
		show: func(ctx context.Context) ([]sdk.NotificationIntegration, error) {
			return client.NotificationIntegrations.Show(ctx, sdk.NewShowNotificationIntegrationRequest())
		},
		name:       func(object sdk.NotificationIntegration) string { return object.Name },
		createdOn:  func(object sdk.NotificationIntegration) (time.Time, error) { return object.CreatedOn, nil },
		id:         func(object sdk.NotificationIntegration) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely: client.NotificationIntegrations.DropSafely,
	})
}

func nukeStorageIntegrations(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.StorageIntegration]{
		objectTypeName: "storage integration",
		protectedNames: []string{
			"S3_STORAGE_INTEGRATION",
			"AZURE_STORAGE_INTEGRATION",
			"GCP_STORAGE_INTEGRATION",
			"TEST_INTEGRATION",
		},
		show: func(ctx context.Context) ([]sdk.StorageIntegration, error) {
			return client.StorageIntegrations.Show(ctx, sdk.NewShowStorageIntegrationRequest())
		},
		name:       func(object sdk.StorageIntegration) string { return object.Name },
		createdOn:  func(object sdk.StorageIntegration) (time.Time, error) { return object.CreatedOn, nil },
		id:         func(object sdk.StorageIntegration) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely: client.StorageIntegrations.DropSafely,
	})
}

func nukeApiIntegrations(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.ApiIntegration]{
		objectTypeName: "api integration",
		show: func(ctx context.Context) ([]sdk.ApiIntegration, error) {
			return client.ApiIntegrations.Show(ctx, sdk.NewShowApiIntegrationRequest())
		},
		name:       func(object sdk.ApiIntegration) string { return object.Name },
		createdOn:  func(object sdk.ApiIntegration) (time.Time, error) { return object.CreatedOn, nil },
		id:         func(object sdk.ApiIntegration) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely: client.ApiIntegrations.DropSafely,
	})
}

func nukeCatalogIntegrations(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.CatalogIntegration]{
		objectTypeName: "catalog integration",
		show: func(ctx context.Context) ([]sdk.CatalogIntegration, error) {
			return client.CatalogIntegrations.Show(ctx, sdk.NewShowCatalogIntegrationRequest())
		},
		name:       func(object sdk.CatalogIntegration) string { return object.Name },
		createdOn:  func(object sdk.CatalogIntegration) (time.Time, error) { return object.CreatedOn, nil },
		id:         func(object sdk.CatalogIntegration) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely: client.CatalogIntegrations.DropSafely,
	})
}

func nukeExternalAccessIntegrations(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.ExternalAccessIntegration]{
		objectTypeName: "external access integration",
		show: func(ctx context.Context) ([]sdk.ExternalAccessIntegration, error) {
			return client.ExternalAccessIntegrations.Show(ctx, sdk.NewShowExternalAccessIntegrationRequest())
		},
		name:       func(object sdk.ExternalAccessIntegration) string { return object.Name },
		createdOn:  func(object sdk.ExternalAccessIntegration) (time.Time, error) { return object.CreatedOn, nil },
		id:         func(object sdk.ExternalAccessIntegration) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely: client.ExternalAccessIntegrations.DropSafely,
	})
}

func nukeApplications(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.Application]{
		objectTypeName: "application",
		protectedNames: []string{
			"SNOWFLAKE",
		},
		show: func(ctx context.Context) ([]sdk.Application, error) {
			return client.Applications.Show(ctx, sdk.NewShowApplicationRequest())
		},
		name:      func(object sdk.Application) string { return object.Name },
		createdOn: func(object sdk.Application) (time.Time, error) { return object.CreatedOn, nil },
		id:        func(object sdk.Application) sdk.AccountObjectIdentifier { return object.ID() },
		// DROP APPLICATION fails when the app owns objects outside itself, and an app can own compute pools,
		// warehouses, and databases. CASCADE drops those too, so it's used instead of DropSafely, which
		// does not support it. Note that nukeComputePools skips the app-owned pools, so without CASCADE
		// neither the app nor its pool would ever be swept.
		dropSafely: func(ctx context.Context, id sdk.AccountObjectIdentifier) error {
			return client.Applications.Drop(ctx, sdk.NewDropApplicationRequest(id).WithIfExists(true).WithCascade(true))
		},
	})
}

func nukeListings(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.Listing]{
		objectTypeName: "listing",
		// no listings are protected for now
		protectedNames: []string{},
		show: func(ctx context.Context) ([]sdk.Listing, error) {
			return client.Listings.Show(ctx, sdk.NewShowListingRequest())
		},
		name:       func(object sdk.Listing) string { return object.Name },
		createdOn:  func(object sdk.Listing) (time.Time, error) { return object.CreatedOn, nil },
		id:         func(object sdk.Listing) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely: client.Listings.DropSafely,
	})
}

func nukeApplicationPackages(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.ApplicationPackage]{
		objectTypeName: "application package",
		// no application packages are protected for now
		protectedNames: []string{},
		show: func(ctx context.Context) ([]sdk.ApplicationPackage, error) {
			return client.ApplicationPackages.Show(ctx, sdk.NewShowApplicationPackageRequest())
		},
		name:       func(object sdk.ApplicationPackage) string { return object.Name },
		createdOn:  func(object sdk.ApplicationPackage) (time.Time, error) { return object.CreatedOn, nil },
		id:         func(object sdk.ApplicationPackage) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely: client.ApplicationPackages.DropSafely,
	})
}

func nukeExternalVolumes(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.ExternalVolume]{
		objectTypeName: "external volume",
		// no external volumes are protected for now
		protectedNames: []string{},
		show: func(ctx context.Context) ([]sdk.ExternalVolume, error) {
			return client.ExternalVolumes.Show(ctx, sdk.NewShowExternalVolumeRequest())
		},
		name: func(object sdk.ExternalVolume) string { return object.Name },
		id:   func(object sdk.ExternalVolume) sdk.AccountObjectIdentifier { return object.ID() },
		// the creation time is not in the SHOW output, so it has to be queried for if needed
		createdOn: func(object sdk.ExternalVolume) (time.Time, error) {
			return externalVolumeCreatedOn(context.Background(), client, object.ID())
		},
		dropSafely: client.ExternalVolumes.DropSafely,
	})
}

func nukeConnections(client *sdk.Client, suffix string) func() error {
	accountLocator := client.GetAccountLocator()

	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.Connection]{
		objectTypeName: "connection",
		// no connections are protected for now
		protectedNames: []string{},
		// the connections replicated from the other accounts are not ours to drop
		skip: func(object sdk.Connection) bool { return object.AccountLocator != accountLocator },
		show: func(ctx context.Context) ([]sdk.Connection, error) {
			return client.Connections.Show(ctx, sdk.NewShowConnectionRequest())
		},
		name:      func(object sdk.Connection) string { return object.Name },
		createdOn: func(object sdk.Connection) (time.Time, error) { return object.CreatedOn, nil },
		id:        func(object sdk.Connection) sdk.AccountObjectIdentifier { return object.ID() },
		// dropping a connection keeps failing until the replication settles, so it's retried the same way
		// the tests clean it up (see ConnectionClient.DropFunc)
		dropSafely: retryingDropSafely(client.Connections.DropSafely, connectionDropAttempts, connectionDropInterval),
	})
}

func nukeComputePools(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.ComputePool]{
		objectTypeName: "compute pool",
		// no compute pools are protected for now
		protectedNames: []string{},
		// the compute pools belonging to a Native App are not ours to drop
		skip: func(object sdk.ComputePool) bool { return object.Application != nil },
		show: func(ctx context.Context) ([]sdk.ComputePool, error) {
			return client.ComputePools.Show(ctx, sdk.NewShowComputePoolRequest())
		},
		name:          func(object sdk.ComputePool) string { return object.Name },
		createdOn:     func(object sdk.ComputePool) (time.Time, error) { return object.CreatedOn, nil },
		id:            func(object sdk.ComputePool) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely:    client.ComputePools.DropSafely,
		owner:         func(object sdk.ComputePool) string { return object.Owner },
		takeOwnership: grantOwnershipToAccountadmin(client, sdk.ObjectTypeComputePool),
	})
}

// TODO [SNOW-867247]: generalize nuke methods (sweepers too)
func nukeWarehouses(client *sdk.Client, prefix string, suffix string) func() error {
	return nukeAccountObjects(prefix, suffix, accountObjectSweeperConfig[sdk.Warehouse]{
		objectTypeName: "warehouse",
		protectedNames: []string{
			"SNOWFLAKE",
			"SYSTEM$STREAMLIT_NOTEBOOK_WH",
		},
		show: func(ctx context.Context) ([]sdk.Warehouse, error) {
			return client.Warehouses.Show(ctx, sdk.NewShowWarehouseRequest())
		},
		name:      func(object sdk.Warehouse) string { return object.Name },
		createdOn: func(object sdk.Warehouse) (time.Time, error) { return object.CreatedOn, nil },
		// TODO [SNOW-1569516]: Use the usual constructors instead.
		id: func(object sdk.Warehouse) sdk.AccountObjectIdentifier {
			return sdk.NewAccountObjectIdentifierNoTrimTestOnly(object.Name)
		},
		dropSafely:    client.Warehouses.DropSafely,
		owner:         func(object sdk.Warehouse) string { return object.Owner },
		takeOwnership: grantOwnershipToAccountadmin(client, sdk.ObjectTypeWarehouse),
	})
}

func nukeDatabases(client *sdk.Client, prefix string, suffix string) func() error {
	return nukeAccountObjects(prefix, suffix, accountObjectSweeperConfig[sdk.Database]{
		objectTypeName: "database",
		protectedNames: []string{
			"SNOWFLAKE",
			"MFA_ENFORCEMENT_POLICY",
			"TERRAFORM_TEST_SETUP_OBJECTS",
			"TEST_RESULTS_DATABASE",
		},
		// SHOW DATABASES returns applications, application packages, and personal databases too, and they can't be
		// handled as databases; applications and application packages have their own sweepers (see nukeApplications
		// and nukeApplicationPackages). Only the kinds this sweeper can drop are handled here.
		skip: func(object sdk.Database) bool {
			if object.Kind == nil {
				log.Printf("[INFO] Skipping database %s, its kind was not recognized", object.Name)
				return true
			}
			return !slices.Contains([]sdk.DatabaseKind{sdk.DatabaseKindStandard, sdk.DatabaseKindImportedDatabase}, *object.Kind)
		},
		show: func(ctx context.Context) ([]sdk.Database, error) {
			return client.Databases.Show(ctx, sdk.NewShowDatabaseRequest())
		},
		name:      func(object sdk.Database) string { return object.Name },
		createdOn: func(object sdk.Database) (time.Time, error) { return object.CreatedOn, nil },
		// TODO [SNOW-1569516]: Use the usual constructors instead.
		id: func(object sdk.Database) sdk.AccountObjectIdentifier {
			return sdk.NewAccountObjectIdentifierNoTrimTestOnly(object.Name)
		},
		dropSafely:    client.Databases.DropSafely,
		owner:         func(object sdk.Database) string { return object.Owner },
		takeOwnership: grantOwnershipToAccountadmin(client, sdk.ObjectTypeDatabase),
	})
}

func nukeUsers(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.User]{
		objectTypeName: "user",
		protectedNames: []string{
			"ARTUR_SAWICKI",
			"JAKUB_MICHALAK",
			"JAN_CIESLAK",
			"KAMIL_WASILEWSKI",
			"PIOTR_CICHON",
			"TEST_CI_SERVICE_USER",
			"PENTESTING_USER_1",
			"PENTESTING_USER_2",
		},
		show: func(ctx context.Context) ([]sdk.User, error) {
			return client.Users.Show(ctx, sdk.NewShowUserRequest())
		},
		name:                func(object sdk.User) string { return object.Name },
		createdOn:           func(object sdk.User) (time.Time, error) { return object.CreatedOn, nil },
		id:                  func(object sdk.User) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely:          client.Users.DropSafely,
		stalePeriodOverride: sdk.Pointer(-15 * time.Minute),
	})
}

func nukePostgresInstances(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.PostgresInstance]{
		objectTypeName: "postgres instance",
		// no postgres instances are protected for now
		protectedNames: []string{},
		show: func(ctx context.Context) ([]sdk.PostgresInstance, error) {
			return client.PostgresInstances.Show(ctx, sdk.NewShowPostgresInstanceRequest())
		},
		name:                func(object sdk.PostgresInstance) string { return object.Name },
		createdOn:           func(object sdk.PostgresInstance) (time.Time, error) { return object.CreatedOn, nil },
		id:                  func(object sdk.PostgresInstance) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely:          client.PostgresInstances.DropSafely,
		stalePeriodOverride: sdk.Pointer(-60 * time.Minute),
	})
}

func nukeSecurityIntegrations(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.SecurityIntegration]{
		objectTypeName: "security integration",
		protectedNames: []string{
			"SNOWFLAKE$LOCAL_APPLICATION",
		},
		show: func(ctx context.Context) ([]sdk.SecurityIntegration, error) {
			return client.SecurityIntegrations.Show(ctx, sdk.NewShowSecurityIntegrationRequest())
		},
		name:                func(object sdk.SecurityIntegration) string { return object.Name },
		createdOn:           func(object sdk.SecurityIntegration) (time.Time, error) { return object.CreatedOn, nil },
		id:                  func(object sdk.SecurityIntegration) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely:          client.SecurityIntegrations.DropSafely,
		stalePeriodOverride: sdk.Pointer(-15 * time.Minute),
	})
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
		snowflakeroles.Restricted,
		snowflakeroles.OktaProvisioner,
		snowflakeroles.AadProvisioner,
		snowflakeroles.GenericScimProvisioner,
	}

	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.Role]{
		objectTypeName: "role",
		protectedNames: collections.Map(protectedRoles, func(id sdk.AccountObjectIdentifier) string { return id.Name() }),
		show: func(ctx context.Context) ([]sdk.Role, error) {
			return client.Roles.Show(ctx, sdk.NewShowRoleRequest())
		},
		name:                func(object sdk.Role) string { return object.Name },
		createdOn:           func(object sdk.Role) (time.Time, error) { return object.CreatedOn, nil },
		id:                  func(object sdk.Role) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely:          client.Roles.DropSafely,
		stalePeriodOverride: sdk.Pointer(-15 * time.Minute),
		owner:               func(object sdk.Role) string { return object.Owner },
		takeOwnership:       grantOwnershipToAccountadmin(client, sdk.ObjectTypeRole),
	})
}

func nukeShares(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.Share]{
		objectTypeName: "share",
		protectedNames: []string{
			// this one is INBOUND but putting it here either way
			"ACCOUNT_USAGE",
		},
		// only the outbound shares are ours to drop
		skip: func(object sdk.Share) bool { return object.Kind != sdk.ShareKindOutbound },
		show: func(ctx context.Context) ([]sdk.Share, error) {
			return client.Shares.Show(ctx, sdk.NewShowShareRequest())
		},
		name:       func(object sdk.Share) string { return object.Name },
		createdOn:  func(object sdk.Share) (time.Time, error) { return object.CreatedOn, nil },
		id:         func(object sdk.Share) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely: client.Shares.DropSafely,
	})
}

// nukeNetworkPolicies was introduced to make sure that network policies created during tests are cleaned up.
// It's required as network policies that have connections to the network rules within databases, block their deletion.
// In Snowflake, the network policies can be removed without unsetting network rules, but the network rules cannot be removed without unsetting network policies.
func nukeNetworkPolicies(client *sdk.Client, suffix string) func() error {
	protectedNetworkPolicies := []string{
		"RESTRICTED_ACCESS",
	}

	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.NetworkPolicy]{
		objectTypeName: "network policy",
		// the protected names are compared case-insensitively, so protectedNames can't be used here
		skip: func(object sdk.NetworkPolicy) bool {
			return slices.Contains(protectedNetworkPolicies, strings.ToUpper(object.Name))
		},
		show: func(ctx context.Context) ([]sdk.NetworkPolicy, error) {
			return client.NetworkPolicies.Show(ctx, sdk.NewShowNetworkPolicyRequest())
		},
		name:       func(object sdk.NetworkPolicy) string { return object.Name },
		createdOn:  func(object sdk.NetworkPolicy) (time.Time, error) { return object.CreatedOn, nil },
		id:         func(object sdk.NetworkPolicy) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely: client.NetworkPolicies.DropSafely,
	})
}

func nukeResourceMonitors(client *sdk.Client, suffix string) func() error {
	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.ResourceMonitor]{
		objectTypeName: "resource monitor",
		// no resource monitors are protected for now
		protectedNames: []string{},
		show: func(ctx context.Context) ([]sdk.ResourceMonitor, error) {
			return client.ResourceMonitors.Show(ctx, sdk.NewShowResourceMonitorRequest())
		},
		name:      func(object sdk.ResourceMonitor) string { return object.Name },
		createdOn: func(object sdk.ResourceMonitor) (time.Time, error) { return object.CreatedOn, nil },
		// TODO [SNOW-1569516]: Use the usual constructors instead.
		id: func(object sdk.ResourceMonitor) sdk.AccountObjectIdentifier {
			return sdk.NewAccountObjectIdentifierNoTrimTestOnly(object.Name)
		},
		dropSafely:    client.ResourceMonitors.DropSafely,
		owner:         func(object sdk.ResourceMonitor) string { return object.Owner },
		takeOwnership: grantOwnershipToAccountadmin(client, sdk.ObjectTypeResourceMonitor),
	})
}

func nukeFailoverGroups(client *sdk.Client, suffix string) func() error {
	accountLocator := client.GetAccountLocator()

	return nukeAccountObjects("", suffix, accountObjectSweeperConfig[sdk.FailoverGroup]{
		objectTypeName: "failover group",
		// no failover groups are protected for now
		protectedNames: []string{},
		// the failover groups replicated from the other accounts are not ours to drop
		skip: func(object sdk.FailoverGroup) bool { return object.AccountLocator != accountLocator },
		show: func(ctx context.Context) ([]sdk.FailoverGroup, error) {
			req := sdk.NewShowFailoverGroupRequest().WithInAccount(sdk.NewAccountIdentifierFromAccountLocator(accountLocator))
			return client.FailoverGroups.Show(ctx, req)
		},
		name:          func(object sdk.FailoverGroup) string { return object.Name },
		createdOn:     func(object sdk.FailoverGroup) (time.Time, error) { return object.CreatedOn, nil },
		id:            func(object sdk.FailoverGroup) sdk.AccountObjectIdentifier { return object.ID() },
		dropSafely:    client.FailoverGroups.DropSafely,
		owner:         func(object sdk.FailoverGroup) string { return object.Owner },
		takeOwnership: grantOwnershipToAccountadmin(client, sdk.ObjectTypeFailoverGroup),
	})
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

func azureTestClient(t *testing.T) *sdk.Client {
	t.Helper()
	return testClient(t, testprofiles.Azure)
}

func gcpTestClient(t *testing.T) *sdk.Client {
	t.Helper()
	return testClient(t, testprofiles.Gcp)
}

func snowflakeDefaultsTestClient(t *testing.T) *sdk.Client {
	t.Helper()
	return testClient(t, testprofiles.SnowflakeDefaults)
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
