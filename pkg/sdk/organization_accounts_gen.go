package sdk

import (
	"context"
	"database/sql"
)

type OrganizationAccounts interface {
	Create(ctx context.Context, request *CreateOrganizationAccountRequest) error
	Alter(ctx context.Context, request *AlterOrganizationAccountRequest) error
	Show(ctx context.Context, request *ShowOrganizationAccountRequest) ([]OrganizationAccount, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*OrganizationAccount, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*OrganizationAccount, error)
	// ShowParameters added manually
	ShowParameters(ctx context.Context) ([]*Parameter, error)
	// UnsetAllParameters added manually
	UnsetAllParameters(ctx context.Context) error
}

// CreateOrganizationAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-organization-account.
type CreateOrganizationAccountOptions struct {
	create              bool                       `ddl:"static" sql:"CREATE"`
	organizationAccount bool                       `ddl:"static" sql:"ORGANIZATION ACCOUNT"`
	name                AccountObjectIdentifier    `ddl:"identifier"`
	AdminName           string                     `ddl:"parameter,no_quotes" sql:"ADMIN_NAME"`
	AdminPassword       *string                    `ddl:"parameter,single_quotes" sql:"ADMIN_PASSWORD"`
	AdminRsaPublicKey   *string                    `ddl:"parameter,single_quotes" sql:"ADMIN_RSA_PUBLIC_KEY"`
	FirstName           *string                    `ddl:"parameter,single_quotes" sql:"FIRST_NAME"`
	LastName            *string                    `ddl:"parameter,single_quotes" sql:"LAST_NAME"`
	Email               string                     `ddl:"parameter,single_quotes" sql:"EMAIL"`
	MustChangePassword  *bool                      `ddl:"parameter" sql:"MUST_CHANGE_PASSWORD"`
	Edition             OrganizationAccountEdition `ddl:"parameter,no_quotes" sql:"EDITION"`
	RegionGroup         *string                    `ddl:"parameter,double_quotes" sql:"REGION_GROUP"`
	Region              *string                    `ddl:"parameter,double_quotes" sql:"REGION"`
	Comment             *string                    `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterOrganizationAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-organization-account.
type AlterOrganizationAccountOptions struct {
	alter               bool                       `ddl:"static" sql:"ALTER"`
	organizationAccount bool                       `ddl:"static" sql:"ORGANIZATION ACCOUNT"`
	Name                *AccountObjectIdentifier   `ddl:"identifier"`
	Set                 *OrganizationAccountSet    `ddl:"keyword" sql:"SET"`
	Unset               *OrganizationAccountUnset  `ddl:"keyword" sql:"UNSET"`
	SetTags             []TagAssociation           `ddl:"keyword" sql:"SET TAG"`
	UnsetTags           []ObjectIdentifier         `ddl:"keyword" sql:"UNSET TAG"`
	RenameTo            *OrganizationAccountRename `ddl:"keyword"`
	DropOldUrl          *bool                      `ddl:"keyword" sql:"DROP OLD URL"`
}

type OrganizationAccountSet struct {
	Parameters      *AccountParameters       `ddl:"list,no_parentheses"`
	ResourceMonitor *AccountObjectIdentifier `ddl:"identifier,equals" sql:"RESOURCE_MONITOR"`
	PasswordPolicy  *SchemaObjectIdentifier  `ddl:"identifier" sql:"PASSWORD POLICY"`
	SessionPolicy   *SchemaObjectIdentifier  `ddl:"identifier" sql:"SESSION POLICY"`
	Comment         *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type OrganizationAccountUnset struct {
	Parameters      *AccountParametersUnset `ddl:"list,no_parentheses"`
	ResourceMonitor *bool                   `ddl:"keyword" sql:"RESOURCE_MONITOR"`
	PasswordPolicy  *bool                   `ddl:"keyword" sql:"PASSWORD POLICY"`
	SessionPolicy   *bool                   `ddl:"keyword" sql:"SESSION POLICY"`
	Comment         *bool                   `ddl:"keyword" sql:"COMMENT"`
}

type OrganizationAccountRename struct {
	NewName    *AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SaveOldUrl *bool                    `ddl:"parameter" sql:"SAVE_OLD_URL"`
}

// ShowOrganizationAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-organization-accounts.
type ShowOrganizationAccountOptions struct {
	show                 bool  `ddl:"static" sql:"SHOW"`
	organizationAccounts bool  `ddl:"static" sql:"ORGANIZATION ACCOUNTS"`
	Like                 *Like `ddl:"keyword" sql:"LIKE"`
}

type organizationAccountDbRow struct {
	OrganizationName                     string         `db:"organization_name"`
	AccountName                          string         `db:"account_name"`
	SnowflakeRegion                      string         `db:"snowflake_region"`
	Edition                              string         `db:"edition"`
	AccountUrl                           string         `db:"account_url"`
	CreatedOn                            string         `db:"created_on"`
	Comment                              sql.NullString `db:"comment"`
	AccountLocator                       string         `db:"account_locator"`
	AccountLocatorUrl                    string         `db:"account_locator_url"`
	ManagedAccounts                      int            `db:"managed_accounts"`
	ConsumptionBillingEntityName         string         `db:"consumption_billing_entity_name"`
	MarketplaceConsumerBillingEntityName sql.NullString `db:"marketplace_consumer_billing_entity_name"`
	MarketplaceProviderBillingEntityName string         `db:"marketplace_provider_billing_entity_name"`
	OldAccountUrl                        sql.NullString `db:"old_account_url"`
	IsOrgAdmin                           bool           `db:"is_org_admin"`
	AccountOldUrlSavedOn                 sql.NullString `db:"account_old_url_saved_on"`
	AccountOldUrlLastUsed                sql.NullString `db:"account_old_url_last_used"`
	OrganizationOldUrl                   sql.NullString `db:"organization_old_url"`
	OrganizationOldUrlSavedOn            sql.NullString `db:"organization_old_url_saved_on"`
	OrganizationOldUrlLastUsed           sql.NullString `db:"organization_old_url_last_used"`
	IsEventsAccount                      bool           `db:"is_events_account"`
	IsOrganizationAccount                bool           `db:"is_organization_account"`
}

type OrganizationAccount struct {
	OrganizationName                     string
	AccountName                          string
	SnowflakeRegion                      string
	Edition                              OrganizationAccountEdition
	AccountUrl                           string
	CreatedOn                            string
	Comment                              *string
	AccountLocator                       string
	AccountLocatorUrl                    string
	ManagedAccounts                      int
	ConsumptionBillingEntityName         string
	MarketplaceConsumerBillingEntityName *string
	MarketplaceProviderBillingEntityName string
	OldAccountUrl                        *string
	IsOrgAdmin                           bool
	AccountOldUrlSavedOn                 *string
	AccountOldUrlLastUsed                *string
	OrganizationOldUrl                   *string
	OrganizationOldUrlSavedOn            *string
	OrganizationOldUrlLastUsed           *string
	IsEventsAccount                      bool
	IsOrganizationAccount                bool
}

func (v *OrganizationAccount) ID() AccountIdentifier {
	return NewAccountIdentifier(v.OrganizationName, v.AccountName)
}

func (v *OrganizationAccount) ObjectType() ObjectType {
	return ObjectTypeAccount
}
