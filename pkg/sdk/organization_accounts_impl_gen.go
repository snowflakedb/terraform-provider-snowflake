package sdk

import (
	"context"
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

// ShowParameters added manually
func (v *organizationAccounts) ShowParameters(ctx context.Context) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Account: Bool(true),
		},
	})
}

// UnsetAllParameters added manually
func (v *organizationAccounts) UnsetAllParameters(ctx context.Context) error {
	return v.client.Accounts.UnsetAllParameters(ctx)
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
		}
	}
	if r.Unset != nil {
		opts.Unset = &OrganizationAccountUnset{
			Parameters:      r.Unset.Parameters,
			ResourceMonitor: r.Unset.ResourceMonitor,
			PasswordPolicy:  r.Unset.PasswordPolicy,
			SessionPolicy:   r.Unset.SessionPolicy,
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
		Comment:                              r.Comment,
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
	mapStringIfNotNil(&oa.MarketplaceConsumerBillingEntityName, r.MarketplaceConsumerBillingEntityName)
	mapStringIfNotNil(&oa.OldAccountUrl, r.OldAccountUrl)
	mapStringIfNotNil(&oa.AccountOldUrlSavedOn, r.AccountOldUrlSavedOn)
	mapStringIfNotNil(&oa.AccountOldUrlLastUsed, r.AccountOldUrlLastUsed)
	mapStringIfNotNil(&oa.OrganizationOldUrl, r.OrganizationOldUrl)
	mapStringIfNotNil(&oa.OrganizationOldUrlSavedOn, r.OrganizationOldUrlSavedOn)
	mapStringIfNotNil(&oa.OrganizationOldUrlLastUsed, r.OrganizationOldUrlLastUsed)
	return oa
}
