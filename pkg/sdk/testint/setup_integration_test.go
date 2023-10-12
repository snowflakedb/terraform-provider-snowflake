package testint

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

var itc integrationTestContext

func TestMain(m *testing.M) {
	exitVal := execute(m)
	os.Exit(exitVal)
}

func execute(m *testing.M) int {
	defer timer("tests")()
	setup()
	exitVal := m.Run()
	cleanup()
	return exitVal
}

func setup() {
	log.Println("Running integration tests setup")

	err := itc.initialize()
	if err != nil {
		log.Printf("Integration test context initialisation failed with %s\n", err)
		os.Exit(1)
	}
}

func cleanup() {
	log.Println("Running integration tests cleanup")

}

type integrationTestContext struct {
	client *sdk.Client
	ctx    context.Context

	database *sdk.Database
}

func (itc *integrationTestContext) initialize() error {
	log.Println("Initializing integration test context")
	var err error
	itc.client, err = sdk.NewDefaultClient()
	itc.ctx = context.Background()

	db, dbCleanup := createDatabase(itc.client)
	t.Cleanup(databaseCleanup)
	return err
}

func createDatabase(client *sdk.Client) (*sdk.Database, func()) {
	return createDatabaseWithOptions(client, sdk.RandomAccountObjectIdentifier(), &sdk.CreateDatabaseOptions{})
}

func createDatabaseWithOptions(client *sdk.Client, id sdk.AccountObjectIdentifier) (*sdk.Database, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.Databases.Create(ctx, id, nil)
	require.NoError(t, err)
	database, err := client.Databases.ShowByID(ctx, id)
	require.NoError(t, err)
	return database, func() {
		err := client.Databases.Drop(ctx, id, nil)
		require.NoError(t, err)
	}
}

// timer measures time from invocation point to the end of method.
// It's supposed to be used like:
//
//	defer timer("something to measure name")()
func timer(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("[DEBUG] %s took %v\n", name, time.Since(start))
	}
}

// TODO: Discuss after this initial change is merged.
// This is temporary way to move all integration tests to this package without doing revolution in a single PR.
func testClient(t *testing.T) *sdk.Client {
	t.Helper()
	return itc.client
}

// TODO: Discuss after this initial change is merged.
// This is temporary way to move all integration tests to this package without doing revolution in a single PR.
func testContext(t *testing.T) context.Context {
	t.Helper()
	return itc.ctx
}
