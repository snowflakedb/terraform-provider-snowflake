package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Accounts interface {
	Create(ctx context.Context, request *CreateAccountRequest) error
	Alter(ctx context.Context, request *AlterAccountRequest) error
	Show(ctx context.Context, request *ShowAccountRequest) ([]Account, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Account, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Account, error)
	Drop(ctx context.Context, request *DropAccountRequest) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Undrop(ctx context.Context, request *UndropAccountRequest) error
	ShowParameters(ctx context.Context) ([]*Parameter, error)
	UnsetAllParameters(ctx context.Context) error
	// UnsetAllPoliciesSafely calls UnsetPolicySafely for every policy that can be unset from the current account.
	UnsetAllPoliciesSafely(ctx context.Context) error
	// UnsetPolicySafely unsets a policy on the current account by a given supported kind.
	// It ignores an error that occurs on the Snowflake side whenever you try to unset policy which is already unset.
	UnsetPolicySafely(ctx context.Context, kind PolicyKind) error
	// UnsetAll unsets all policies and parameters that can be attached to the current account.
	UnsetAll(ctx context.Context) error
}

// CreateAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-account.
type CreateAccountOptions struct {
	create  bool                    `ddl:"static" sql:"CREATE"`
	account bool                    `ddl:"static" sql:"ACCOUNT"`
	name    AccountObjectIdentifier `ddl:"identifier"`

	AdminName                string         `ddl:"parameter,single_quotes" sql:"ADMIN_NAME"`
	AdminPassword            *string        `ddl:"parameter,single_quotes" sql:"ADMIN_PASSWORD"`
	AdminRSAPublicKey        *string        `ddl:"parameter,single_quotes" sql:"ADMIN_RSA_PUBLIC_KEY"`
	AdminUserType            *UserType      `ddl:"parameter" sql:"ADMIN_USER_TYPE"`
	FirstName                *string        `ddl:"parameter,single_quotes" sql:"FIRST_NAME"`
	LastName                 *string        `ddl:"parameter,single_quotes" sql:"LAST_NAME"`
	Email                    string         `ddl:"parameter,single_quotes" sql:"EMAIL"`
	MustChangePassword       *bool          `ddl:"parameter" sql:"MUST_CHANGE_PASSWORD"`
	Edition                  AccountEdition `ddl:"parameter" sql:"EDITION"`
	RegionGroup              *string        `ddl:"parameter,single_quotes" sql:"REGION_GROUP"`
	Region                   *string        `ddl:"parameter,single_quotes" sql:"REGION"`
	Comment                  *string        `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ConsumptionBillingEntity *string        `ddl:"parameter,double_quotes" sql:"CONSUMPTION_BILLING_ENTITY"`
	Polaris                  *bool          `ddl:"parameter" sql:"POLARIS"`
}

// AlterAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-account.
type AlterAccountOptions struct {
	alter   bool                     `ddl:"static" sql:"ALTER"`
	account bool                     `ddl:"static" sql:"ACCOUNT"`
	Name    *AccountObjectIdentifier `ddl:"identifier"`

	Set      *AccountSet        `ddl:"keyword" sql:"SET"`
	Unset    *AccountUnset      `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTag   []TagAssociation   `ddl:"keyword" sql:"SET TAG"`
	UnsetTag []ObjectIdentifier `ddl:"keyword" sql:"UNSET TAG"`
	Rename   *AccountRename     `ddl:"-"`
	Drop     *AccountDrop       `ddl:"-"`
}

type AccountLevelParameters struct {
	AccountParameters *LegacyAccountParameters `ddl:"list,no_parentheses"`
	SessionParameters *SessionParameters       `ddl:"list,no_parentheses"`
	ObjectParameters  *ObjectParameters        `ddl:"list,no_parentheses"`
	UserParameters    *UserParameters          `ddl:"list,no_parentheses"`
}

type AccountSet struct {
	Parameters               *AccountParameters       `ddl:"list,no_parentheses"`
	LegacyParameters         *AccountLevelParameters  `ddl:"list,no_parentheses"`
	ResourceMonitor          *AccountObjectIdentifier `ddl:"identifier,equals" sql:"RESOURCE_MONITOR"`
	PackagesPolicy           *SchemaObjectIdentifier  `ddl:"identifier" sql:"PACKAGES POLICY"`
	PasswordPolicy           *SchemaObjectIdentifier  `ddl:"identifier" sql:"PASSWORD POLICY"`
	SessionPolicy            *SchemaObjectIdentifier  `ddl:"identifier" sql:"SESSION POLICY"`
	AuthenticationPolicy     *SchemaObjectIdentifier  `ddl:"identifier" sql:"AUTHENTICATION POLICY"`
	FeaturePolicySet         *AccountFeaturePolicySet `ddl:"keyword"`
	ConsumptionBillingEntity *string                  `ddl:"parameter,double_quotes" sql:"CONSUMPTION_BILLING_ENTITY"`
	OrgAdmin                 *bool                    `ddl:"parameter" sql:"IS_ORG_ADMIN"`
	Force                    *bool                    `ddl:"keyword" sql:"FORCE"`
}

type AccountFeaturePolicySet struct {
	FeaturePolicy      *SchemaObjectIdentifier `ddl:"identifier" sql:"FEATURE POLICY"`
	forAllApplications bool                    `ddl:"static" sql:"FOR ALL APPLICATIONS"`
}

type AccountLevelParametersUnset struct {
	AccountParameters *LegacyAccountParametersUnset `ddl:"list,no_parentheses"`
	SessionParameters *SessionParametersUnset       `ddl:"list,no_parentheses"`
	ObjectParameters  *ObjectParametersUnset        `ddl:"list,no_parentheses"`
	UserParameters    *UserParametersUnset          `ddl:"list,no_parentheses"`
}

type AccountUnset struct {
	Parameters               *AccountParametersUnset      `ddl:"list,no_parentheses"`
	LegacyParameters         *AccountLevelParametersUnset `ddl:"list,no_parentheses"`
	AuthenticationPolicy     *bool                        `ddl:"keyword" sql:"AUTHENTICATION POLICY"`
	FeaturePolicyUnset       *AccountFeaturePolicyUnset   `ddl:"keyword"`
	PackagesPolicy           *bool                        `ddl:"keyword" sql:"PACKAGES POLICY"`
	PasswordPolicy           *bool                        `ddl:"keyword" sql:"PASSWORD POLICY"`
	SessionPolicy            *bool                        `ddl:"keyword" sql:"SESSION POLICY"`
	ResourceMonitor          *bool                        `ddl:"keyword" sql:"RESOURCE_MONITOR"`
	ConsumptionBillingEntity *bool                        `ddl:"keyword" sql:"CONSUMPTION_BILLING_ENTITY"`
}

type AccountFeaturePolicyUnset struct {
	FeaturePolicy      *bool `ddl:"keyword" sql:"FEATURE POLICY"`
	forAllApplications bool  `ddl:"static" sql:"FOR ALL APPLICATIONS"`
}

type AccountRename struct {
	NewName    AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SaveOldURL *bool                   `ddl:"parameter" sql:"SAVE_OLD_URL"`
}

type AccountDrop struct {
	OldUrl             *bool `ddl:"keyword" sql:"DROP OLD URL"`
	OldOrganizationUrl *bool `ddl:"keyword" sql:"DROP OLD ORGANIZATION URL"`
}

// ShowAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-organisation-accounts.
type ShowAccountOptions struct {
	show     bool  `ddl:"static" sql:"SHOW"`
	accounts bool  `ddl:"static" sql:"ACCOUNTS"`
	History  *bool `ddl:"keyword" sql:"HISTORY"`
	Like     *Like `ddl:"keyword" sql:"LIKE"`
}

// DropAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-account.
type DropAccountOptions struct {
	drop              bool                    `ddl:"static" sql:"DROP"`
	account           bool                    `ddl:"static" sql:"ACCOUNT"`
	IfExists          *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name              AccountObjectIdentifier `ddl:"identifier"`
	GracePeriodInDays *int                    `ddl:"parameter" sql:"GRACE_PERIOD_IN_DAYS"`
}

// UndropAccountOptions is based on https://docs.snowflake.com/en/sql-reference/sql/undrop-account.
type UndropAccountOptions struct {
	undrop  bool                    `ddl:"static" sql:"UNDROP"`
	account bool                    `ddl:"static" sql:"ACCOUNT"`
	name    AccountObjectIdentifier `ddl:"identifier"`
}

type Account struct {
	OrganizationName                     string
	AccountName                          string
	SnowflakeRegion                      string
	RegionGroup                          *string
	Edition                              *AccountEdition
	AccountURL                           *string
	CreatedOn                            *time.Time
	Comment                              *string
	AccountLocator                       string
	AccountLocatorUrl                    *string
	ManagedAccounts                      *int
	ConsumptionBillingEntityName         *string
	MarketplaceConsumerBillingEntityName *string
	MarketplaceProviderBillingEntityName *string
	OldAccountURL                        *string
	IsOrgAdmin                           *bool
	AccountOldUrlSavedOn                 *time.Time
	AccountOldUrlLastUsed                *time.Time
	OrganizationOldUrl                   *string
	OrganizationOldUrlSavedOn            *time.Time
	OrganizationOldUrlLastUsed           *time.Time
	IsEventsAccount                      *bool
	IsOrganizationAccount                bool
	DroppedOn                            *time.Time
	ScheduledDeletionTime                *time.Time
	RestoredOn                           *time.Time
	MovedToOrganization                  *string
	MovedOn                              *string
	OrganizationUrlExpirationOn          *time.Time
}

type accountDBRow struct {
	OrganizationName                     string         `db:"organization_name"`
	AccountName                          string         `db:"account_name"`
	RegionGroup                          sql.NullString `db:"region_group"`
	SnowflakeRegion                      string         `db:"snowflake_region"`
	Edition                              sql.NullString `db:"edition"`
	AccountURL                           sql.NullString `db:"account_url"`
	CreatedOn                            sql.NullTime   `db:"created_on"`
	Comment                              sql.NullString `db:"comment"`
	AccountLocator                       string         `db:"account_locator"`
	AccountLocatorURL                    sql.NullString `db:"account_locator_url"`
	ManagedAccounts                      sql.NullInt32  `db:"managed_accounts"`
	ConsumptionBillingEntityName         sql.NullString `db:"consumption_billing_entity_name"`
	MarketplaceConsumerBillingEntityName sql.NullString `db:"marketplace_consumer_billing_entity_name"`
	MarketplaceProviderBillingEntityName sql.NullString `db:"marketplace_provider_billing_entity_name"`
	OldAccountURL                        sql.NullString `db:"old_account_url"`
	IsOrgAdmin                           sql.NullBool   `db:"is_org_admin"`
	AccountOldUrlSavedOn                 sql.NullTime   `db:"account_old_url_saved_on"`
	AccountOldUrlLastUsed                sql.NullTime   `db:"account_old_url_last_used"`
	OrganizationOldUrl                   sql.NullString `db:"organization_old_url"`
	OrganizationOldUrlSavedOn            sql.NullTime   `db:"organization_old_url_saved_on"`
	OrganizationOldUrlLastUsed           sql.NullTime   `db:"organization_old_url_last_used"`
	IsEventsAccount                      sql.NullBool   `db:"is_events_account"`
	IsOrganizationAccount                bool           `db:"is_organization_account"`
	DroppedOn                            sql.NullTime   `db:"dropped_on"`
	ScheduledDeletionTime                sql.NullTime   `db:"scheduled_deletion_time"`
	RestoredOn                           sql.NullTime   `db:"restored_on"`
	MovedToOrganization                  sql.NullString `db:"moved_to_organization"`
	MovedOn                              sql.NullString `db:"moved_on"`
	OrganizationUrlExpirationOn          sql.NullTime   `db:"organization_URL_expiration_on"`
}
