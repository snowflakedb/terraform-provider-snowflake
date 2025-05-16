package helpers

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

type AccountClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewAccountClient(context *TestClientContext, idsGenerator *IdsGenerator) *AccountClient {
	return &AccountClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *AccountClient) client() sdk.Accounts {
	return c.context.client.Accounts
}

// GetAccountIdentifier gets the account identifier from Snowflake API, by fetching the account locator
// and by filtering the list of accounts in replication accounts by it (because there is no direct way to get).
func (c *AccountClient) GetAccountIdentifier(t *testing.T) sdk.AccountIdentifier {
	t.Helper()
	ctx := context.Background()

	currentAccountLocator, err := c.context.client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)

	replicationAccounts, err := c.context.client.ReplicationFunctions.ShowReplicationAccounts(ctx)
	require.NoError(t, err)

	for _, replicationAccount := range replicationAccounts {
		if replicationAccount.AccountLocator == currentAccountLocator {
			return sdk.NewAccountIdentifier(replicationAccount.OrganizationName, replicationAccount.AccountName)
		}
	}
	t.Fatal("could not find the account identifier for the locator")
	return sdk.AccountIdentifier{}
}

func (c *AccountClient) Create(t *testing.T) (*sdk.Account, func()) {
	t.Helper()
	id := c.ids.RandomSensitiveAccountObjectIdentifier()
	name := random.AdminName()
	email := random.Email()
	publicKey, _ := random.GenerateRSAPublicKey(t)

	return c.CreateWithRequest(t, id, &sdk.CreateAccountOptions{
		AdminName:         name,
		AdminRSAPublicKey: sdk.String(publicKey),
		Email:             email,
		Edition:           sdk.EditionStandard,
	})
}

func (c *AccountClient) CreateWithRequest(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.CreateAccountOptions) (*sdk.Account, func()) {
	t.Helper()
	_, err := c.client().Create(context.Background(), id, opts)
	require.NoError(t, err)

	account, err := c.client().ShowByID(context.Background(), id)
	require.NoError(t, err)

	return account, c.DropFunc(t, id)
}

func (c *AccountClient) Alter(t *testing.T, opts *sdk.AlterAccountOptions) {
	t.Helper()
	err := c.client().Alter(context.Background(), opts)
	require.NoError(t, err)
}

func (c *AccountClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	return func() {
		require.NoError(t, c.Drop(t, id))
	}
}

func (c *AccountClient) Drop(t *testing.T, id sdk.AccountObjectIdentifier) error {
	t.Helper()
	ctx := context.Background()

	return c.client().Drop(ctx, id, 3, &sdk.DropAccountOptions{IfExists: sdk.Bool(true)})
}

type Region struct {
	SnowflakeRegion string `db:"snowflake_region"`
	Cloud           string `db:"cloud"`
	Region          string `db:"region"`
	DisplayName     string `db:"display_name"`
}

func (c *AccountClient) ShowRegions(t *testing.T) []Region {
	t.Helper()

	var regions []Region
	err := c.context.client.QueryForTests(context.Background(), &regions, "SHOW REGIONS")
	require.NoError(t, err)

	return regions
}

func (c *AccountClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Account, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}

func (c *AccountClient) CreateAndLogIn(t *testing.T) (*sdk.Account, *sdk.Client, func()) {
	t.Helper()
	id := c.ids.RandomSensitiveAccountObjectIdentifier()
	name := random.AdminName()
	privateKey := random.GenerateRSAPrivateKey(t)
	publicKey, _ := random.GenerateRSAPublicKeyFromPrivateKey(t, privateKey)
	email := random.Email()

	account, accountCleanup := c.CreateWithRequest(t, id, &sdk.CreateAccountOptions{
		AdminName:         name,
		AdminRSAPublicKey: sdk.String(publicKey),
		AdminUserType:     sdk.Pointer(sdk.UserTypeService),
		Email:             email,
		Edition:           sdk.EditionStandard,
	})

	var client *sdk.Client
	require.Eventually(t, func() bool {
		newClient, err := sdk.NewClient(&gosnowflake.Config{
			Account:       fmt.Sprintf("%s-%s", account.OrganizationName, account.AccountName),
			User:          name,
			Host:          strings.TrimPrefix(*account.AccountLocatorUrl, `https://`),
			Authenticator: gosnowflake.AuthTypeJwt,
			PrivateKey:    privateKey,
			Role:          snowflakeroles.Accountadmin.Name(),
		})
		client = newClient
		return err == nil
	}, 2*time.Minute, time.Second*15)

	return account, client, accountCleanup
}
