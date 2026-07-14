package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ Accounts                = (*accounts)(nil)
	_ convertibleRow[Account] = new(accountDBRow)
)

type accounts struct {
	client *Client
}

func (c *accounts) Create(ctx context.Context, request *CreateAccountRequest) error {
	opts := request.toOpts()
	return validateAndExec(c.client, ctx, opts)
}

func (c *accounts) Alter(ctx context.Context, request *AlterAccountRequest) error {
	opts := request.toOpts()
	return validateAndExec(c.client, ctx, opts)
}

func (c *accounts) Show(ctx context.Context, request *ShowAccountRequest) ([]Account, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[accountDBRow](c.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[accountDBRow, Account](dbRows)
}

func (c *accounts) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Account, error) {
	request := NewShowAccountRequest().
		WithLike(Like{Pattern: String(id.Name())})
	accounts, err := c.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(accounts, func(account Account) bool {
		return account.AccountName == id.Name()
	})
}

func (c *accounts) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Account, error) {
	return SafeShowById(c.client, c.ShowByID, ctx, id)
}

func (c *accounts) Drop(ctx context.Context, request *DropAccountRequest) error {
	opts := request.toOpts()
	return validateAndExec(c.client, ctx, opts)
}

func (c *accounts) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(c.client, func() error { return c.Drop(ctx, NewDropAccountRequest(id).WithIfExists(true)) }, ctx, id)
}

func (c *accounts) Undrop(ctx context.Context, request *UndropAccountRequest) error {
	opts := request.toOpts()
	return validateAndExec(c.client, ctx, opts)
}

func (row accountDBRow) convert() (*Account, error) {
	acc := &Account{
		OrganizationName:      row.OrganizationName,
		AccountName:           row.AccountName,
		SnowflakeRegion:       row.SnowflakeRegion,
		AccountLocator:        row.AccountLocator,
		IsOrganizationAccount: row.IsOrganizationAccount,
	}
	if row.RegionGroup.Valid {
		acc.RegionGroup = &row.RegionGroup.String
	}
	if row.Edition.Valid {
		acc.Edition = Pointer(AccountEdition(row.Edition.String))
	}
	if row.AccountURL.Valid {
		acc.AccountURL = &row.AccountURL.String
	}
	if row.CreatedOn.Valid {
		acc.CreatedOn = &row.CreatedOn.Time
	}
	if row.Comment.Valid {
		acc.Comment = &row.Comment.String
	}
	if row.AccountLocatorURL.Valid {
		acc.AccountLocatorUrl = &row.AccountLocatorURL.String
	}
	if row.ConsumptionBillingEntityName.Valid {
		acc.ConsumptionBillingEntityName = &row.ConsumptionBillingEntityName.String
	}
	if row.OldAccountURL.Valid {
		acc.OldAccountURL = &row.OldAccountURL.String
	}
	if row.IsOrgAdmin.Valid {
		acc.IsOrgAdmin = &row.IsOrgAdmin.Bool
	}
	if row.OrganizationOldUrl.Valid {
		acc.OrganizationOldUrl = &row.OrganizationOldUrl.String
	}
	if row.IsEventsAccount.Valid {
		acc.IsEventsAccount = &row.IsEventsAccount.Bool
	}
	if row.MarketplaceConsumerBillingEntityName.Valid {
		acc.MarketplaceConsumerBillingEntityName = &row.MarketplaceConsumerBillingEntityName.String
	}
	if row.MarketplaceProviderBillingEntityName.Valid {
		acc.MarketplaceProviderBillingEntityName = &row.MarketplaceProviderBillingEntityName.String
	}
	if row.AccountOldUrlSavedOn.Valid {
		acc.AccountOldUrlSavedOn = &row.AccountOldUrlSavedOn.Time
	}
	if row.AccountOldUrlLastUsed.Valid {
		acc.AccountOldUrlLastUsed = &row.AccountOldUrlLastUsed.Time
	}
	if row.OrganizationOldUrlSavedOn.Valid {
		acc.OrganizationOldUrlSavedOn = &row.OrganizationOldUrlSavedOn.Time
	}
	if row.OrganizationOldUrlLastUsed.Valid {
		acc.OrganizationOldUrlLastUsed = &row.OrganizationOldUrlLastUsed.Time
	}
	if row.DroppedOn.Valid {
		acc.DroppedOn = &row.DroppedOn.Time
	}
	if row.ScheduledDeletionTime.Valid {
		acc.ScheduledDeletionTime = &row.ScheduledDeletionTime.Time
	}
	if row.RestoredOn.Valid {
		acc.RestoredOn = &row.RestoredOn.Time
	}
	if row.MovedToOrganization.Valid {
		acc.MovedToOrganization = &row.MovedToOrganization.String
	}
	if row.MovedOn.Valid {
		acc.MovedOn = &row.MovedOn.String
	}
	if row.OrganizationUrlExpirationOn.Valid {
		acc.OrganizationUrlExpirationOn = &row.OrganizationUrlExpirationOn.Time
	}
	return row.additionalConvert(acc), nil
}
