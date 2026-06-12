package sdk

var (
	_ optionsProvider[createDatabaseRoleOptions] = new(CreateDatabaseRoleRequest)
	_ optionsProvider[alterDatabaseRoleOptions]  = new(AlterDatabaseRoleRequest)
	_ optionsProvider[dropDatabaseRoleOptions]   = new(DropDatabaseRoleRequest)
	_ optionsProvider[showDatabaseRoleOptions]   = new(ShowDatabaseRoleRequest)
)

type CreateDatabaseRoleRequest struct {
	orReplace   bool
	ifNotExists bool
	name        DatabaseObjectIdentifier // required
	comment     *string
}

type AlterDatabaseRoleRequest struct {
	ifExists bool
	name     DatabaseObjectIdentifier // required

	// One of
	rename    *DatabaseObjectIdentifier
	set       *DatabaseRoleSetRequest
	unset     *DatabaseRoleUnsetRequest
	setTags   []TagAssociation
	unsetTags []ObjectIdentifier
}

type DatabaseRoleSetRequest struct {
	comment string // required
}

type DatabaseRoleUnsetRequest struct{}

type DropDatabaseRoleRequest struct {
	ifExists bool
	name     DatabaseObjectIdentifier // required
}

type ShowDatabaseRoleRequest struct {
	like     *Like
	database AccountObjectIdentifier // required
	limit    *LimitFrom
}

type GrantDatabaseRoleRequest struct {
	name DatabaseObjectIdentifier // required

	// One of
	databaseRole *DatabaseObjectIdentifier
	accountRole  *AccountObjectIdentifier
}

type RevokeDatabaseRoleRequest struct {
	name DatabaseObjectIdentifier // required

	// One of
	databaseRole *DatabaseObjectIdentifier
	accountRole  *AccountObjectIdentifier
}

type GrantDatabaseRoleToShareRequest struct {
	name  DatabaseObjectIdentifier // required
	share AccountObjectIdentifier  // required
}

type RevokeDatabaseRoleFromShareRequest struct {
	name  DatabaseObjectIdentifier // required
	share AccountObjectIdentifier  // required
}
