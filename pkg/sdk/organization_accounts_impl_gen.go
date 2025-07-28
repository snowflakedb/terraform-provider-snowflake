package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ OrganizationAccounts = (*organizationAccounts)(nil)

type organizationAccounts struct {
	client *Client
}

func (v *organizationAccounts) Create(ctx context.Context, request *CreateOrganizationAccountRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *organizationAccounts) Alter(ctx context.Context, request *AlterOrganizationAccountRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *organizationAccounts) Show(ctx context.Context, request *ShowOrganizationAccountRequest) ([]OrganizationAccount, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[organizationAccountDbRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[organizationAccountDbRow, OrganizationAccount](dbRows)
	return resultList, nil
}

func (v *organizationAccounts) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*OrganizationAccount, error) {
	request := NewShowOrganizationAccountRequest().
		WithLike(Like{Pattern: String(id.Name())})
	organizationAccounts, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(organizationAccounts, func(r OrganizationAccount) bool { return r.AccountName == id.Name() })
}

func (v *organizationAccounts) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*OrganizationAccount, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (r *CreateOrganizationAccountRequest) toOpts() *CreateOrganizationAccountOptions {
	opts := &CreateOrganizationAccountOptions{
		name:               r.name,
		AdminName:          r.AdminName,
		AdminPassword:      r.AdminPassword,
		AdminRsaPublicKey:  r.AdminRsaPublicKey,
		FirstName:          r.FirstName,
		LastName:           r.LastName,
		Email:              r.Email,
		MustChangePassword: r.MustChangePassword,
		Edition:            r.Edition,
		RegionGroup:        r.RegionGroup,
		Region:             r.Region,
		Comment:            r.Comment,
	}
	return opts
}

func (r *AlterOrganizationAccountRequest) toOpts() *AlterOrganizationAccountOptions {
	opts := &AlterOrganizationAccountOptions{
		Name: r.Name,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,

		DropOldUrl: r.DropOldUrl,
	}
	if r.Set != nil {
		opts.Set = &OrganizationAccountSet{
			Parameters:      r.Set.Parameters,
			ResourceMonitor: r.Set.ResourceMonitor,
			PasswordPolicy:  r.Set.PasswordPolicy,
			SessionPolicy:   r.Set.SessionPolicy,
			Comment:         r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &OrganizationAccountUnset{
			Parameters:      r.Unset.Parameters,
			ResourceMonitor: r.Unset.ResourceMonitor,
			PasswordPolicy:  r.Unset.PasswordPolicy,
			SessionPolicy:   r.Unset.SessionPolicy,
			Comment:         r.Unset.Comment,
		}
	}
	if r.RenameTo != nil {
		opts.RenameTo = &OrganizationAccountRename{
			NewName:    r.RenameTo.NewName,
			SaveOldUrl: r.RenameTo.SaveOldUrl,
		}
	}
	return opts
}

func (r *ShowOrganizationAccountRequest) toOpts() *ShowOrganizationAccountOptions {
	opts := &ShowOrganizationAccountOptions{
		Like: r.Like,
	}
	return opts
}

func (r organizationAccountDbRow) convert() *OrganizationAccount {
	oa := &OrganizationAccount{
		OrganizationName:                     r.OrganizationName,
		AccountName:                          r.AccountName,
		SnowflakeRegion:                      r.SnowflakeRegion,
		AccountUrl:                           r.AccountUrl,
		CreatedOn:                            r.CreatedOn,
		AccountLocator:                       r.AccountLocator,
		AccountLocatorUrl:                    r.AccountLocatorUrl,
		ManagedAccounts:                      r.ManagedAccounts,
		ConsumptionBillingEntityName:         r.ConsumptionBillingEntityName,
		MarketplaceProviderBillingEntityName: r.MarketplaceProviderBillingEntityName,
		IsOrgAdmin:                           r.IsOrgAdmin,
		IsEventsAccount:                      r.IsEventsAccount,
		IsOrganizationAccount:                r.IsOrganizationAccount,
	}
	mapStringWithMapping(&oa.Edition, r.Edition, ToOrganizationAccountEdition)
	mapNullString(&oa.Comment, r.Comment)
	mapNullString(&oa.MarketplaceConsumerBillingEntityName, r.MarketplaceConsumerBillingEntityName)
	mapNullString(&oa.OldAccountUrl, r.OldAccountUrl)
	mapNullString(&oa.AccountOldUrlSavedOn, r.AccountOldUrlSavedOn)
	mapNullString(&oa.AccountOldUrlLastUsed, r.AccountOldUrlLastUsed)
	mapNullString(&oa.OrganizationOldUrl, r.OrganizationOldUrl)
	mapNullString(&oa.OrganizationOldUrlSavedOn, r.OrganizationOldUrlSavedOn)
	mapNullString(&oa.OrganizationOldUrlLastUsed, r.OrganizationOldUrlLastUsed)
	return oa
}
