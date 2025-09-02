package sdk

import (
	"context"
	"log"
	"strings"
	"time"
)

type Grants interface {
	GrantPrivilegesToAccountRole(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, role AccountObjectIdentifier, opts *GrantPrivilegesToAccountRoleOptions) error
	RevokePrivilegesFromAccountRole(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, role AccountObjectIdentifier, opts *RevokePrivilegesFromAccountRoleOptions) error
	GrantPrivilegesToDatabaseRole(ctx context.Context, privileges *DatabaseRoleGrantPrivileges, on *DatabaseRoleGrantOn, role DatabaseObjectIdentifier, opts *GrantPrivilegesToDatabaseRoleOptions) error
	RevokePrivilegesFromDatabaseRole(ctx context.Context, privileges *DatabaseRoleGrantPrivileges, on *DatabaseRoleGrantOn, role DatabaseObjectIdentifier, opts *RevokePrivilegesFromDatabaseRoleOptions) error
	GrantPrivilegeToShare(ctx context.Context, privileges []ObjectPrivilege, on *ShareGrantOn, to AccountObjectIdentifier) error
	RevokePrivilegeFromShare(ctx context.Context, privileges []ObjectPrivilege, on *ShareGrantOn, from AccountObjectIdentifier) error
	GrantOwnership(ctx context.Context, on OwnershipGrantOn, to OwnershipGrantTo, opts *GrantOwnershipOptions) error

	Show(ctx context.Context, opts *ShowGrantOptions) ([]Grant, error)
}

// GrantPrivilegesToAccountRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/grant-privilege#syntax.
type GrantPrivilegesToAccountRoleOptions struct {
	grant           bool                        `ddl:"static" sql:"GRANT"`
	privileges      *AccountRoleGrantPrivileges `ddl:"-"`
	on              *AccountRoleGrantOn         `ddl:"keyword" sql:"ON"`
	accountRole     AccountObjectIdentifier     `ddl:"identifier" sql:"TO ROLE"`
	WithGrantOption *bool                       `ddl:"keyword" sql:"WITH GRANT OPTION"`
}

type AccountRoleGrantPrivileges struct {
	GlobalPrivileges        []GlobalPrivilege        `ddl:"-"`
	AccountObjectPrivileges []AccountObjectPrivilege `ddl:"-"`
	SchemaPrivileges        []SchemaPrivilege        `ddl:"-"`
	SchemaObjectPrivileges  []SchemaObjectPrivilege  `ddl:"-"`
	AllPrivileges           *bool                    `ddl:"keyword" sql:"ALL PRIVILEGES"`
}

type AccountRoleGrantOn struct {
	Account       *bool                 `ddl:"keyword" sql:"ACCOUNT"`
	AccountObject *GrantOnAccountObject `ddl:"-"`
	Schema        *GrantOnSchema        `ddl:"-"`
	SchemaObject  *GrantOnSchemaObject  `ddl:"-"`
}

type GrantOnAccountObject struct {
	User             *AccountObjectIdentifier `ddl:"identifier" sql:"USER"`
	ResourceMonitor  *AccountObjectIdentifier `ddl:"identifier" sql:"RESOURCE MONITOR"`
	Warehouse        *AccountObjectIdentifier `ddl:"identifier" sql:"WAREHOUSE"`
	ComputePool      *AccountObjectIdentifier `ddl:"identifier" sql:"COMPUTE POOL"`
	Database         *AccountObjectIdentifier `ddl:"identifier" sql:"DATABASE"`
	Integration      *AccountObjectIdentifier `ddl:"identifier" sql:"INTEGRATION"`
	FailoverGroup    *AccountObjectIdentifier `ddl:"identifier" sql:"FAILOVER GROUP"`
	ReplicationGroup *AccountObjectIdentifier `ddl:"identifier" sql:"REPLICATION GROUP"`
	ExternalVolume   *AccountObjectIdentifier `ddl:"identifier" sql:"EXTERNAL VOLUME"`
}

type GrantOnSchema struct {
	Schema                  *DatabaseObjectIdentifier `ddl:"identifier" sql:"SCHEMA"`
	AllSchemasInDatabase    *AccountObjectIdentifier  `ddl:"identifier" sql:"ALL SCHEMAS IN DATABASE"`
	FutureSchemasInDatabase *AccountObjectIdentifier  `ddl:"identifier" sql:"FUTURE SCHEMAS IN DATABASE"`
}

type GrantOnSchemaObject struct {
	SchemaObject *Object                `ddl:"-"`
	All          *GrantOnSchemaObjectIn `ddl:"keyword" sql:"ALL"`
	Future       *GrantOnSchemaObjectIn `ddl:"keyword" sql:"FUTURE"`
}

type GrantOnSchemaObjectIn struct {
	PluralObjectType PluralObjectType          `ddl:"keyword" sql:"ALL"`
	InDatabase       *AccountObjectIdentifier  `ddl:"identifier" sql:"IN DATABASE"`
	InSchema         *DatabaseObjectIdentifier `ddl:"identifier" sql:"IN SCHEMA"`
}

// RevokePrivilegesFromAccountRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege#syntax.
type RevokePrivilegesFromAccountRoleOptions struct {
	revoke         bool                        `ddl:"static" sql:"REVOKE"`
	GrantOptionFor *bool                       `ddl:"keyword" sql:"GRANT OPTION FOR"`
	privileges     *AccountRoleGrantPrivileges `ddl:"-"`
	on             *AccountRoleGrantOn         `ddl:"keyword" sql:"ON"`
	accountRole    AccountObjectIdentifier     `ddl:"identifier" sql:"FROM ROLE"`
	Restrict       *bool                       `ddl:"keyword" sql:"RESTRICT"`
	Cascade        *bool                       `ddl:"keyword" sql:"CASCADE"`
}

// GrantPrivilegesToDatabaseRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/grant-privilege#syntax.
type GrantPrivilegesToDatabaseRoleOptions struct {
	grant           bool                         `ddl:"static" sql:"GRANT"`
	privileges      *DatabaseRoleGrantPrivileges `ddl:"-"`
	on              *DatabaseRoleGrantOn         `ddl:"keyword" sql:"ON"`
	databaseRole    DatabaseObjectIdentifier     `ddl:"identifier" sql:"TO DATABASE ROLE"`
	WithGrantOption *bool                        `ddl:"keyword" sql:"WITH GRANT OPTION"`
}

type DatabaseRoleGrantPrivileges struct {
	DatabasePrivileges     []AccountObjectPrivilege `ddl:"-"`
	SchemaPrivileges       []SchemaPrivilege        `ddl:"-"`
	SchemaObjectPrivileges []SchemaObjectPrivilege  `ddl:"-"`
	AllPrivileges          *bool                    `ddl:"keyword" sql:"ALL PRIVILEGES"`
}

type DatabaseRoleGrantOn struct {
	Database     *AccountObjectIdentifier `ddl:"identifier" sql:"DATABASE"`
	Schema       *GrantOnSchema           `ddl:"-"`
	SchemaObject *GrantOnSchemaObject     `ddl:"-"`
}

// RevokePrivilegesFromDatabaseRoleOptions is based on https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege#syntax.
type RevokePrivilegesFromDatabaseRoleOptions struct {
	revoke         bool                         `ddl:"static" sql:"REVOKE"`
	GrantOptionFor *bool                        `ddl:"keyword" sql:"GRANT OPTION FOR"`
	privileges     *DatabaseRoleGrantPrivileges `ddl:"-"`
	on             *DatabaseRoleGrantOn         `ddl:"keyword" sql:"ON"`
	databaseRole   DatabaseObjectIdentifier     `ddl:"identifier" sql:"FROM DATABASE ROLE"`
	Restrict       *bool                        `ddl:"keyword" sql:"RESTRICT"`
	Cascade        *bool                        `ddl:"keyword" sql:"CASCADE"`
}

// grantPrivilegeToShareOptions is based on https://docs.snowflake.com/en/sql-reference/sql/grant-privilege-share.
type grantPrivilegeToShareOptions struct {
	grant      bool                    `ddl:"static" sql:"GRANT"`
	privileges []ObjectPrivilege       `ddl:"-"`
	On         *ShareGrantOn           `ddl:"keyword" sql:"ON"`
	to         AccountObjectIdentifier `ddl:"identifier" sql:"TO SHARE"`
}

type ShareGrantOn struct {
	Database AccountObjectIdentifier             `ddl:"identifier" sql:"DATABASE"`
	Schema   DatabaseObjectIdentifier            `ddl:"identifier" sql:"SCHEMA"`
	Function SchemaObjectIdentifierWithArguments `ddl:"identifier" sql:"FUNCTION"`
	Table    *OnTable                            `ddl:"-"`
	Tag      SchemaObjectIdentifier              `ddl:"identifier" sql:"TAG"`
	View     SchemaObjectIdentifier              `ddl:"identifier" sql:"VIEW"`
}

type OnTable struct {
	Name        SchemaObjectIdentifier   `ddl:"identifier" sql:"TABLE"`
	AllInSchema DatabaseObjectIdentifier `ddl:"identifier" sql:"ALL TABLES IN SCHEMA"`
}

// revokePrivilegeFromShareOptions is based on https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege-share.
type revokePrivilegeFromShareOptions struct {
	revoke     bool                    `ddl:"static" sql:"REVOKE"`
	privileges []ObjectPrivilege       `ddl:"-"`
	On         *ShareGrantOn           `ddl:"keyword" sql:"ON"`
	from       AccountObjectIdentifier `ddl:"identifier" sql:"FROM SHARE"`
}

type OnView struct {
	Name        SchemaObjectIdentifier   `ddl:"identifier" sql:"VIEW"`
	AllInSchema DatabaseObjectIdentifier `ddl:"identifier" sql:"ALL VIEWS IN SCHEMA"`
}

// ShowGrantOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-grants.
type ShowGrantOptions struct {
	show   bool          `ddl:"static" sql:"SHOW"`
	Future *bool         `ddl:"keyword" sql:"FUTURE"`
	grants bool          `ddl:"static" sql:"GRANTS"`
	On     *ShowGrantsOn `ddl:"keyword" sql:"ON"`
	To     *ShowGrantsTo `ddl:"keyword" sql:"TO"`
	Of     *ShowGrantsOf `ddl:"keyword" sql:"OF"`
	In     *ShowGrantsIn `ddl:"keyword" sql:"IN"`
}

type ShowGrantsIn struct {
	Schema   *DatabaseObjectIdentifier `ddl:"identifier" sql:"SCHEMA"`
	Database *AccountObjectIdentifier  `ddl:"identifier" sql:"DATABASE"`
}

type ShowGrantsOn struct {
	Account *bool `ddl:"keyword" sql:"ACCOUNT"`
	Object  *Object
}

type ShowGrantsTo struct {
	Application     AccountObjectIdentifier  `ddl:"identifier" sql:"APPLICATION"`
	ApplicationRole DatabaseObjectIdentifier `ddl:"identifier" sql:"APPLICATION ROLE"`
	Role            AccountObjectIdentifier  `ddl:"identifier" sql:"ROLE"`
	User            AccountObjectIdentifier  `ddl:"identifier" sql:"USER"`
	Share           *ShowGrantsToShare       `ddl:"-"`
	DatabaseRole    DatabaseObjectIdentifier `ddl:"identifier" sql:"DATABASE ROLE"`
}

type ShowGrantsToShare struct {
	Name                 AccountObjectIdentifier  `ddl:"identifier" sql:"SHARE"`
	InApplicationPackage *AccountObjectIdentifier `ddl:"identifier" sql:"IN APPLICATION PACKAGE"`
}

type ShowGrantsOf struct {
	ApplicationRole DatabaseObjectIdentifier `ddl:"identifier" sql:"APPLICATION ROLE"`
	Role            AccountObjectIdentifier  `ddl:"identifier" sql:"ROLE"`
	DatabaseRole    DatabaseObjectIdentifier `ddl:"identifier" sql:"DATABASE ROLE"`
	Share           AccountObjectIdentifier  `ddl:"identifier" sql:"SHARE"`
}

type grantRow struct {
	CreatedOn   time.Time `db:"created_on"`
	Privilege   string    `db:"privilege"`
	GrantedOn   string    `db:"granted_on"`
	GrantOn     string    `db:"grant_on"`
	Name        string    `db:"name"`
	GrantedTo   string    `db:"granted_to"`
	GrantTo     string    `db:"grant_to"`
	GranteeName string    `db:"grantee_name"`
	GrantOption bool      `db:"grant_option"`
	GrantedBy   string    `db:"granted_by"`
}

type Grant struct {
	CreatedOn   time.Time
	Privilege   string
	GrantedOn   ObjectType
	GrantOn     ObjectType
	Name        ObjectIdentifier
	Role        ObjectIdentifier
	GrantedTo   ObjectType
	GrantTo     ObjectType
	GranteeName ObjectIdentifier
	GrantOption bool
	GrantedBy   AccountObjectIdentifier
}

func (v *Grant) ID() ObjectIdentifier {
	return v.Name
}

// TODO(SNOW-2097063): Improve SHOW GRANTS implementation
func (row grantRow) convert() (*Grant, error) {
	grantedTo := ObjectType(strings.ReplaceAll(row.GrantedTo, "_", " "))
	grantTo := ObjectType(strings.ReplaceAll(row.GrantTo, "_", " "))
	var grantedOn ObjectType
	// true for current grants
	if row.GrantedOn != "" {
		grantedOn = ObjectType(strings.ReplaceAll(row.GrantedOn, "_", " "))
	}
	if row.GrantedOn == "VOLUME" {
		grantedOn = ObjectTypeExternalVolume
	}
	if row.GrantedOn == "MODULE" {
		grantedOn = ObjectTypeModel
	}

	var grantOn ObjectType
	// true for future grants
	if row.GrantOn != "" {
		grantOn = ObjectType(strings.ReplaceAll(row.GrantOn, "_", " "))
	}
	if row.GrantOn == "VOLUME" {
		grantOn = ObjectTypeExternalVolume
	}
	if row.GrantOn == "MODULE" {
		grantOn = ObjectTypeModel
	}

	var name ObjectIdentifier
	var err error
	// TODO(SNOW-1569535): use a mapper from object type to parsing function
	if ObjectType(row.GrantedOn).IsWithArguments() {
		name, err = ParseSchemaObjectIdentifierWithArgumentsAndReturnType(row.Name)
	} else {
		name, err = ParseObjectIdentifierString(row.Name)
	}
	if err != nil {
		log.Printf("[DEBUG] Failed to parse identifier [%s], err = \"%s\"; falling back to fully qualified name conversion", row.Name, err)
		name = NewObjectIdentifierFromFullyQualifiedName(row.Name)
	}

	return &Grant{
		CreatedOn: row.CreatedOn,
		Privilege: row.Privilege,
		GrantedOn: grantedOn,
		GrantOn:   grantOn,
		GrantedTo: grantedTo,
		GrantTo:   grantTo,
		Name:      name,
		// GranteeName is computed in Show operation. Its format is depending on the grant request options.
		GrantOption: row.GrantOption,
		GrantedBy:   NewAccountObjectIdentifier(row.GrantedBy),
	}, nil
}

// GrantOwnershipOptions is based on https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#syntax.
// Description is a bit misleading, ownership can be given not only to schema objects but also to account level objects.
type GrantOwnershipOptions struct {
	grantOwnership bool                    `ddl:"static" sql:"GRANT OWNERSHIP"`
	On             OwnershipGrantOn        `ddl:"keyword" sql:"ON"`
	To             OwnershipGrantTo        `ddl:"keyword" sql:"TO"`
	CurrentGrants  *OwnershipCurrentGrants `ddl:"-"`
}

type OwnershipGrantOn struct {
	// One of
	Object *Object                `ddl:"-"`
	All    *GrantOnSchemaObjectIn `ddl:"keyword" sql:"ALL"`
	Future *GrantOnSchemaObjectIn `ddl:"keyword" sql:"FUTURE"`
}

type OwnershipGrantTo struct {
	// One of
	DatabaseRoleName *DatabaseObjectIdentifier `ddl:"identifier" sql:"DATABASE ROLE"`
	AccountRoleName  *AccountObjectIdentifier  `ddl:"identifier" sql:"ROLE"`
}

type OwnershipCurrentGrants struct {
	OutboundPrivileges OwnershipCurrentGrantsOutboundPrivileges `ddl:"keyword"`
	currentGrants      bool                                     `ddl:"static" sql:"CURRENT GRANTS"`
}

type OwnershipCurrentGrantsOutboundPrivileges string

const (
	Revoke OwnershipCurrentGrantsOutboundPrivileges = "REVOKE"
	Copy   OwnershipCurrentGrantsOutboundPrivileges = "COPY"
)
