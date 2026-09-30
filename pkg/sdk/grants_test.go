package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrantPrivilegesToAccountRole(t *testing.T) {
	t.Run("validation: privilege with disallowed characters", func(t *testing.T) {
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				GlobalPrivileges: []GlobalPrivilege{"MONITOR USAGE; SELECT"},
			},
			On: &AccountRoleGrantOn{
				Account: new(true),
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", "MONITOR USAGE; SELECT", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("privileges with certain special characters are allowed", func(t *testing.T) {
		schemaId := randomDatabaseObjectIdentifier()
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{"CREATE SNOWFLAKE.ML.ANOMALY_DETECTION", "applybudget"},
			},
			On: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: new(schemaId),
				},
			},
			AccountRole:     NewAccountObjectIdentifier("role1"),
			WithGrantOption: new(true),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT CREATE SNOWFLAKE.ML.ANOMALY_DETECTION, applybudget ON SCHEMA %s TO ROLE "role1" WITH GRANT OPTION`, schemaId.FullyQualifiedName())
	})

	t.Run("on account", func(t *testing.T) {
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				GlobalPrivileges: []GlobalPrivilege{GlobalPrivilegeMonitorUsage, GlobalPrivilegeApplyTag},
			},
			On: &AccountRoleGrantOn{
				Account: Bool(true),
			},
			AccountRole:     NewAccountObjectIdentifier("role1"),
			WithGrantOption: Bool(true),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT MONITOR USAGE, APPLY TAG ON ACCOUNT TO ROLE "role1" WITH GRANT OPTION`)
	})

	t.Run("on account object", func(t *testing.T) {
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			On: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Object: &Object{
						ObjectType: ObjectTypeDatabase,
						Name:       NewAccountObjectIdentifier("db1"),
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT ALL PRIVILEGES ON DATABASE "db1" TO ROLE "role1"`)
	})

	t.Run("on account object - external volume", func(t *testing.T) {
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			On: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Object: &Object{
						ObjectType: ObjectTypeExternalVolume,
						Name:       NewAccountObjectIdentifier("ex volume"),
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT ALL PRIVILEGES ON EXTERNAL VOLUME "ex volume" TO ROLE "role1"`)
	})

	t.Run("on account object - compute pool", func(t *testing.T) {
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			On: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Object: &Object{
						ObjectType: ObjectTypeComputePool,
						Name:       NewAccountObjectIdentifier("compute pool"),
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT ALL PRIVILEGES ON COMPUTE POOL "compute pool" TO ROLE "role1"`)
	})

	t.Run("on account object - connection", func(t *testing.T) {
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			On: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Object: &Object{
						ObjectType: ObjectTypeConnection,
						Name:       NewAccountObjectIdentifier("myconn"),
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT ALL PRIVILEGES ON CONNECTION "myconn" TO ROLE "role1"`)
	})

	t.Run("on account object - object is required", func(t *testing.T) {
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			On: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("grantPrivilegesToAccountRoleOptions.On.AccountObject", "Object"))
	})

	t.Run("on account object - unknown type fallback", func(t *testing.T) {
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			On: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Object: &Object{
						ObjectType: ObjectTypePostgresInstance,
						Name:       NewAccountObjectIdentifier("pg1"),
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT ALL PRIVILEGES ON POSTGRES INSTANCE "pg1" TO ROLE "role1"`)
	})

	t.Run("on account object - unknown type fallback rejects injection", func(t *testing.T) {
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			On: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Object: &Object{
						ObjectType: ObjectType("TABLE; DROP"),
						Name:       NewAccountObjectIdentifier("pg1"),
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid object type: %s contains disallowed characters; it must follow this regex: %s", "TABLE; DROP", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("on schema", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			On: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(id),
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT CREATE ALERT ON SCHEMA %s TO ROLE "role1"`, id.FullyQualifiedName())
	})

	t.Run("on all schemas in database", func(t *testing.T) {
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			On: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					AllSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT CREATE ALERT ON ALL SCHEMAS IN DATABASE "db1" TO ROLE "role1"`)
	})

	t.Run("on all future schemas in database", func(t *testing.T) {
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			On: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					FutureSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT CREATE ALERT ON FUTURE SCHEMAS IN DATABASE "db1" TO ROLE "role1"`)
	})

	t.Run("on schema object", func(t *testing.T) {
		tableId := randomSchemaObjectIdentifier()
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			On: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{
						ObjectType: ObjectTypeTable,
						Name:       tableId,
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT APPLY ON TABLE %s TO ROLE "role1"`, tableId.FullyQualifiedName())
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			On: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InDatabase:       Pointer(NewAccountObjectIdentifier("db1")),
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT APPLY ON FUTURE TABLES IN DATABASE "db1" TO ROLE "role1"`)
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := &grantPrivilegesToAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			On: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InSchema:         Pointer(id),
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT APPLY ON FUTURE TABLES IN SCHEMA %s TO ROLE "role1"`, id.FullyQualifiedName())
	})
}

func TestRevokePrivilegesFromAccountRole(t *testing.T) {
	schemaId := randomDatabaseObjectIdentifier()

	t.Run("validation: privilege with disallowed characters", func(t *testing.T) {
		opts := &revokePrivilegesFromAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				GlobalPrivileges: []GlobalPrivilege{"MONITOR USAGE; SELECT"},
			},
			On: &AccountRoleGrantOn{
				Account: new(true),
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", "MONITOR USAGE; SELECT", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("on account", func(t *testing.T) {
		opts := &revokePrivilegesFromAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				GlobalPrivileges: []GlobalPrivilege{GlobalPrivilegeMonitorUsage, GlobalPrivilegeApplyTag},
			},
			On: &AccountRoleGrantOn{
				Account: Bool(true),
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE MONITOR USAGE, APPLY TAG ON ACCOUNT FROM ROLE "role1"`)
	})

	t.Run("on account object", func(t *testing.T) {
		opts := &revokePrivilegesFromAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			On: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Object: &Object{
						ObjectType: ObjectTypeDatabase,
						Name:       NewAccountObjectIdentifier("db1"),
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE ALL PRIVILEGES ON DATABASE "db1" FROM ROLE "role1"`)
	})

	t.Run("on account object", func(t *testing.T) {
		opts := &revokePrivilegesFromAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				AccountObjectPrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateDatabaseRole, AccountObjectPrivilegeModify},
			},
			On: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Object: &Object{
						ObjectType: ObjectTypeDatabase,
						Name:       NewAccountObjectIdentifier("db1"),
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE CREATE DATABASE ROLE, MODIFY ON DATABASE "db1" FROM ROLE "role1"`)
	})

	t.Run("on schema", func(t *testing.T) {
		opts := &revokePrivilegesFromAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			On: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(schemaId),
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON SCHEMA %s FROM ROLE "role1"`, schemaId.FullyQualifiedName())
	})

	t.Run("on all schemas in database + restrict", func(t *testing.T) {
		opts := &revokePrivilegesFromAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			On: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					AllSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
			Restrict:    Bool(true),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON ALL SCHEMAS IN DATABASE "db1" FROM ROLE "role1" RESTRICT`)
	})

	t.Run("on all future schemas in database + cascade", func(t *testing.T) {
		opts := &revokePrivilegesFromAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			On: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					FutureSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
			Cascade:     Bool(true),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON FUTURE SCHEMAS IN DATABASE "db1" FROM ROLE "role1" CASCADE`)
	})

	t.Run("on schema object", func(t *testing.T) {
		tableId := randomSchemaObjectIdentifier()
		opts := &revokePrivilegesFromAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
			},
			On: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{
						ObjectType: ObjectTypeTable,
						Name:       tableId,
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE SELECT, UPDATE ON TABLE %s FROM ROLE "role1"`, tableId.FullyQualifiedName())
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := &revokePrivilegesFromAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
			},
			On: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InDatabase:       Pointer(NewAccountObjectIdentifier("db1")),
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE SELECT, UPDATE ON FUTURE TABLES IN DATABASE "db1" FROM ROLE "role1"`)
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := &revokePrivilegesFromAccountRoleOptions{
			Privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
			},
			On: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InSchema:         Pointer(id),
					},
				},
			},
			AccountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE SELECT, UPDATE ON FUTURE TABLES IN SCHEMA %s FROM ROLE "role1"`, id.FullyQualifiedName())
	})
}

func TestGrants_GrantPrivilegesToDatabaseRole(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	databaseRoleId := randomDatabaseObjectIdentifierInDatabase(dbId)
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)

	defaultGrantsForDb := func() *grantPrivilegesToDatabaseRoleOptions {
		return &grantPrivilegesToDatabaseRoleOptions{
			Privileges: &DatabaseRoleGrantPrivileges{
				DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			},
			On: &DatabaseRoleGrantOn{
				Database: &dbId,
			},
			DatabaseRole: databaseRoleId,
		}
	}

	defaultGrantsForSchema := func() *grantPrivilegesToDatabaseRoleOptions {
		return &grantPrivilegesToDatabaseRoleOptions{
			Privileges: &DatabaseRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			On: &DatabaseRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(schemaId),
				},
			},
			DatabaseRole: databaseRoleId,
		}
	}
	tableId := randomSchemaObjectIdentifier()
	defaultGrantsForSchemaObject := func() *grantPrivilegesToDatabaseRoleOptions {
		return &grantPrivilegesToDatabaseRoleOptions{
			Privileges: &DatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			On: &DatabaseRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{
						ObjectType: ObjectTypeTable,
						Name:       tableId,
					},
				},
			},
			DatabaseRole: databaseRoleId,
		}
	}

	t.Run("validation: privilege with disallowed characters", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.Privileges = &DatabaseRoleGrantPrivileges{
			DatabasePrivileges: []AccountObjectPrivilege{"CREATE SCHEMA--"},
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", "CREATE SCHEMA--", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("validation: nil privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.Privileges = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("grantPrivilegesToDatabaseRoleOptions", "Privileges"))
	})

	t.Run("validation: no privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.Privileges = &DatabaseRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantPrivilegesToDatabaseRoleOptions.Privileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges", "AllPrivileges"))
	})

	t.Run("validation: too many privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.Privileges = &DatabaseRoleGrantPrivileges{
			DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			SchemaPrivileges:   []SchemaPrivilege{SchemaPrivilegeCreateAlert},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantPrivilegesToDatabaseRoleOptions.Privileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges", "AllPrivileges"))
	})

	t.Run("validation: no on set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.On = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("grantPrivilegesToDatabaseRoleOptions", "On"))
	})

	t.Run("validation: no on set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.On = &DatabaseRoleGrantOn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantPrivilegesToDatabaseRoleOptions.On", "Database", "Schema", "SchemaObject"))
	})

	t.Run("validation: too many ons set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.On = &DatabaseRoleGrantOn{
			Database: &dbId,
			Schema: &GrantOnSchema{
				Schema: Pointer(schemaId),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantPrivilegesToDatabaseRoleOptions.On", "Database", "Schema", "SchemaObject"))
	})

	t.Run("validation: grant on schema", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.On.Schema = &GrantOnSchema{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantPrivilegesToDatabaseRoleOptions.On.Schema", "Schema", "AllSchemasInDatabase", "FutureSchemasInDatabase"))
	})

	t.Run("validation: grant on schema object", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.On.SchemaObject = &GrantOnSchemaObject{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantPrivilegesToDatabaseRoleOptions.On.SchemaObject", "SchemaObject", "All", "Future"))
	})

	t.Run("validation: grant on schema object - all", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.On = &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				All: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantPrivilegesToDatabaseRoleOptions.On.SchemaObject.All", "InDatabase", "InSchema"))
	})

	t.Run("validation: grant on schema object - future", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.On = &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				Future: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantPrivilegesToDatabaseRoleOptions.On.SchemaObject.Future", "InDatabase", "InSchema"))
	})

	t.Run("on database", func(t *testing.T) {
		opts := defaultGrantsForDb()
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT CREATE SCHEMA ON DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on schema", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT CREATE ALERT ON SCHEMA %s TO DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on all schemas in database", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.On.Schema = &GrantOnSchema{
			AllSchemasInDatabase: Pointer(dbId),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT CREATE ALERT ON ALL SCHEMAS IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on all future schemas in database", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.On.Schema = &GrantOnSchema{
			FutureSchemasInDatabase: Pointer(dbId),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT CREATE ALERT ON FUTURE SCHEMAS IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on schema object", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT APPLY ON TABLE %s TO DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.On.SchemaObject = &GrantOnSchemaObject{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InDatabase:       Pointer(dbId),
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT APPLY ON FUTURE TABLES IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.On.SchemaObject = &GrantOnSchemaObject{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InSchema:         Pointer(schemaId),
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT APPLY ON FUTURE TABLES IN SCHEMA %s TO DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("grant all privileges", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.Privileges = &DatabaseRoleGrantPrivileges{
			AllPrivileges: Bool(true),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT ALL PRIVILEGES ON TABLE %s TO DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})
}

func TestGrants_RevokePrivilegesFromDatabaseRoleRole(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	databaseRoleId := randomDatabaseObjectIdentifierInDatabase(dbId)
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	tableId := randomSchemaObjectIdentifierInSchema(schemaId)

	defaultGrantsForDb := func() *revokePrivilegesFromDatabaseRoleOptions {
		return &revokePrivilegesFromDatabaseRoleOptions{
			Privileges: &DatabaseRoleGrantPrivileges{
				DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			},
			On: &DatabaseRoleGrantOn{
				Database: &dbId,
			},
			DatabaseRole: databaseRoleId,
		}
	}

	defaultGrantsForSchema := func() *revokePrivilegesFromDatabaseRoleOptions {
		return &revokePrivilegesFromDatabaseRoleOptions{
			Privileges: &DatabaseRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			On: &DatabaseRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(schemaId),
				},
			},
			DatabaseRole: databaseRoleId,
		}
	}

	defaultGrantsForSchemaObject := func() *revokePrivilegesFromDatabaseRoleOptions {
		return &revokePrivilegesFromDatabaseRoleOptions{
			Privileges: &DatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
			},
			On: &DatabaseRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{
						ObjectType: ObjectTypeTable,
						Name:       tableId,
					},
				},
			},
			DatabaseRole: databaseRoleId,
		}
	}

	t.Run("validation: privilege with disallowed characters", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.Privileges = &DatabaseRoleGrantPrivileges{
			DatabasePrivileges: []AccountObjectPrivilege{"CREATE SCHEMA--"},
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", "CREATE SCHEMA--", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("validation: nil privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.Privileges = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("revokePrivilegesFromDatabaseRoleOptions", "Privileges"))
	})

	t.Run("validation: no privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.Privileges = &DatabaseRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("revokePrivilegesFromDatabaseRoleOptions.Privileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges", "AllPrivileges"))
	})

	t.Run("validation: too many privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.Privileges = &DatabaseRoleGrantPrivileges{
			DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			SchemaPrivileges:   []SchemaPrivilege{SchemaPrivilegeCreateAlert},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("revokePrivilegesFromDatabaseRoleOptions.Privileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges", "AllPrivileges"))
	})

	t.Run("validation: nil on set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.On = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("revokePrivilegesFromDatabaseRoleOptions", "On"))
	})

	t.Run("validation: no on set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.On = &DatabaseRoleGrantOn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("revokePrivilegesFromDatabaseRoleOptions.On", "Database", "Schema", "SchemaObject"))
	})

	t.Run("validation: too many ons set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.On = &DatabaseRoleGrantOn{
			Database: &dbId,
			Schema: &GrantOnSchema{
				Schema: Pointer(schemaId),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("revokePrivilegesFromDatabaseRoleOptions.On", "Database", "Schema", "SchemaObject"))
	})

	t.Run("validation: grant on schema", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.On.Schema = &GrantOnSchema{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("revokePrivilegesFromDatabaseRoleOptions.On.Schema", "Schema", "AllSchemasInDatabase", "FutureSchemasInDatabase"))
	})

	t.Run("validation: grant on schema object", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.On.SchemaObject = &GrantOnSchemaObject{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("revokePrivilegesFromDatabaseRoleOptions.On.SchemaObject", "SchemaObject", "All", "Future"))
	})

	t.Run("validation: grant on schema object - all", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.On = &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				All: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("revokePrivilegesFromDatabaseRoleOptions.On.SchemaObject.All", "InDatabase", "InSchema"))
	})

	t.Run("validation: grant on schema object - future", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.On = &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				Future: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("revokePrivilegesFromDatabaseRoleOptions.On.SchemaObject.Future", "InDatabase", "InSchema"))
	})

	t.Run("on database", func(t *testing.T) {
		opts := defaultGrantsForDb()
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE CREATE SCHEMA ON DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on schema", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON SCHEMA %s FROM DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on all schemas in database + restrict", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.On.Schema = &GrantOnSchema{
			AllSchemasInDatabase: Pointer(dbId),
		}
		opts.Restrict = Bool(true)
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON ALL SCHEMAS IN DATABASE %s FROM DATABASE ROLE %s RESTRICT`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on all future schemas in database + cascade", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.On.Schema = &GrantOnSchema{
			FutureSchemasInDatabase: Pointer(dbId),
		}
		opts.Cascade = Bool(true)
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON FUTURE SCHEMAS IN DATABASE %s FROM DATABASE ROLE %s CASCADE`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on schema object", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE SELECT, UPDATE ON TABLE %s FROM DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.On.SchemaObject = &GrantOnSchemaObject{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InDatabase:       Pointer(dbId),
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE SELECT, UPDATE ON FUTURE TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.On.SchemaObject = &GrantOnSchemaObject{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InSchema:         Pointer(schemaId),
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE SELECT, UPDATE ON FUTURE TABLES IN SCHEMA %s FROM DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})
}

func TestGrantPrivilegeToShare(t *testing.T) {
	id := randomAccountObjectIdentifier()
	t.Run("validation: privilege with disallowed characters", func(t *testing.T) {
		opts := &GrantPrivilegeToShareOptions{
			Privileges: []ObjectPrivilege{"USAGE;"},
			On: &ShareGrantOn{
				Database: randomAccountObjectIdentifier(),
			},
			To: id,
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", "USAGE;", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("on database", func(t *testing.T) {
		otherID := randomAccountObjectIdentifier()
		opts := &GrantPrivilegeToShareOptions{
			Privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Database: otherID,
			},
			To: id,
		}
		assertOptsValidAndSqlEqualsf(t, opts, "GRANT USAGE ON DATABASE %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on schema", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier()
		opts := &GrantPrivilegeToShareOptions{
			Privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Schema: otherID,
			},
			To: id,
		}
		assertOptsValidAndSqlEqualsf(t, opts, "GRANT USAGE ON SCHEMA %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on table", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier()
		opts := &GrantPrivilegeToShareOptions{
			Privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Table: &OnTable{
					Name: otherID,
				},
			},
			To: id,
		}
		assertOptsValidAndSqlEqualsf(t, opts, "GRANT USAGE ON TABLE %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on all tables", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier()
		opts := &GrantPrivilegeToShareOptions{
			Privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Table: &OnTable{
					AllInSchema: otherID,
				},
			},
			To: id,
		}
		assertOptsValidAndSqlEqualsf(t, opts, "GRANT USAGE ON ALL TABLES IN SCHEMA %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on view", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier()
		opts := &GrantPrivilegeToShareOptions{
			Privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				View: otherID,
			},
			To: id,
		}
		assertOptsValidAndSqlEqualsf(t, opts, "GRANT USAGE ON VIEW %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})
}

func TestRevokePrivilegeFromShare(t *testing.T) {
	id := randomAccountObjectIdentifier()
	t.Run("validation: privilege with disallowed characters", func(t *testing.T) {
		opts := &RevokePrivilegeFromShareOptions{
			Privileges: []ObjectPrivilege{"USAGE;"},
			On: &ShareGrantOn{
				Database: randomAccountObjectIdentifier(),
			},
			From: id,
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", "USAGE;", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("on database", func(t *testing.T) {
		otherID := randomAccountObjectIdentifier()
		opts := &RevokePrivilegeFromShareOptions{
			Privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Database: otherID,
			},
			From: id,
		}
		assertOptsValidAndSqlEqualsf(t, opts, "REVOKE USAGE ON DATABASE %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on schema", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier()
		opts := &RevokePrivilegeFromShareOptions{
			Privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Schema: otherID,
			},
			From: id,
		}
		assertOptsValidAndSqlEqualsf(t, opts, "REVOKE USAGE ON SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on table", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier()
		opts := &RevokePrivilegeFromShareOptions{
			Privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Table: &OnTable{
					Name: otherID,
				},
			},
			From: id,
		}
		assertOptsValidAndSqlEqualsf(t, opts, "REVOKE USAGE ON TABLE %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on all tables", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier()
		opts := &RevokePrivilegeFromShareOptions{
			Privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Table: &OnTable{
					AllInSchema: otherID,
				},
			},
			From: id,
		}
		assertOptsValidAndSqlEqualsf(t, opts, "REVOKE USAGE ON ALL TABLES IN SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on view", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier()
		opts := &RevokePrivilegeFromShareOptions{
			Privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				View: otherID,
			},
			From: id,
		}
		assertOptsValidAndSqlEqualsf(t, opts, "REVOKE USAGE ON VIEW %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on tag", func(t *testing.T) {
		tagId := randomSchemaObjectIdentifier()
		opts := &RevokePrivilegeFromShareOptions{
			Privileges: []ObjectPrivilege{ObjectPrivilegeRead},
			On: &ShareGrantOn{
				Tag: tagId,
			},
			From: id,
		}
		assertOptsValidAndSqlEqualsf(t, opts, "REVOKE READ ON TAG %s FROM SHARE %s", tagId.FullyQualifiedName(), id.FullyQualifiedName())
	})
}

func TestGrants_GrantOwnership(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	roleId := randomAccountObjectIdentifier()
	databaseRoleId := randomDatabaseObjectIdentifierInDatabase(dbId)
	tableId := randomSchemaObjectIdentifierInSchema(schemaId)

	defaultOpts := func() *grantOwnershipOptions {
		return &grantOwnershipOptions{
			On: OwnershipGrantOn{
				Object: &Object{
					ObjectType: ObjectTypeTable,
					Name:       tableId,
				},
			},
			To: OwnershipGrantTo{
				AccountRoleName: Pointer(roleId),
			},
		}
	}

	t.Run("validation: grant on empty", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantOwnershipOptions.On", "Object", "All", "Future"))
	})

	t.Run("validation: grant on too many", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			Object: &Object{
				ObjectType: ObjectTypeTable,
				Name:       tableId,
			},
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InDatabase:       Pointer(dbId),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantOwnershipOptions.On", "Object", "All", "Future"))
	})

	t.Run("validation: grant on schema object - all", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			All: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantOwnershipOptions.On.All", "InDatabase", "InSchema"))
	})

	t.Run("validation: grant on schema object - future", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantOwnershipOptions.On.Future", "InDatabase", "InSchema"))
	})

	t.Run("validation: grant to empty", func(t *testing.T) {
		opts := defaultOpts()
		opts.To = OwnershipGrantTo{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantOwnershipOptions.To", "DatabaseRoleName", "AccountRoleName"))
	})

	t.Run("validation: grant to role and database role", func(t *testing.T) {
		opts := defaultOpts()
		opts.To = OwnershipGrantTo{
			DatabaseRoleName: Pointer(databaseRoleId),
			AccountRoleName:  Pointer(roleId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("grantOwnershipOptions.To", "DatabaseRoleName", "AccountRoleName"))
	})

	t.Run("on schema object to role", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT OWNERSHIP ON TABLE %s TO ROLE %s`, tableId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on dbt project to role", func(t *testing.T) {
		dbtProjectId := randomSchemaObjectIdentifierInSchema(schemaId)
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			Object: &Object{
				ObjectType: ObjectTypeDbtProject,
				Name:       dbtProjectId,
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT OWNERSHIP ON DBT PROJECT %s TO ROLE %s`, dbtProjectId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on schema object to database role", func(t *testing.T) {
		opts := defaultOpts()
		opts.To = OwnershipGrantTo{
			DatabaseRoleName: Pointer(databaseRoleId),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT OWNERSHIP ON TABLE %s TO DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InDatabase:       Pointer(dbId),
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT OWNERSHIP ON FUTURE TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on all schema objects in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			All: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InSchema:         Pointer(schemaId),
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT OWNERSHIP ON ALL TABLES IN SCHEMA %s TO ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on schema object with current grants", func(t *testing.T) {
		opts := defaultOpts()
		opts.CurrentGrants = &OwnershipCurrentGrants{
			OutboundPrivileges: Copy,
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT OWNERSHIP ON TABLE %s TO ROLE %s COPY CURRENT GRANTS`, tableId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})
}

func TestGrants_RevokeOwnership(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	roleId := randomAccountObjectIdentifier()
	databaseRoleId := randomDatabaseObjectIdentifierInDatabase(dbId)

	defaultOpts := func() *RevokeOwnershipOptions {
		return &RevokeOwnershipOptions{
			On: RevokeOwnershipGrantOn{
				Future: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
					InDatabase:       Pointer(dbId),
				},
			},
			From: OwnershipGrantTo{
				AccountRoleName: Pointer(roleId),
			},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *RevokeOwnershipOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: revoke on empty (future not set)", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = RevokeOwnershipGrantOn{}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("RevokeOwnershipOptions.On", "Future"))
	})

	t.Run("validation: revoke from empty", func(t *testing.T) {
		opts := defaultOpts()
		opts.From = OwnershipGrantTo{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RevokeOwnershipOptions.From", "DatabaseRoleName", "AccountRoleName"))
	})

	t.Run("validation: revoke from role and database role", func(t *testing.T) {
		opts := defaultOpts()
		opts.From = OwnershipGrantTo{
			DatabaseRoleName: Pointer(databaseRoleId),
			AccountRoleName:  Pointer(roleId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RevokeOwnershipOptions.From", "DatabaseRoleName", "AccountRoleName"))
	})

	t.Run("validation: restrict and cascade", func(t *testing.T) {
		opts := defaultOpts()
		opts.Restrict = Bool(true)
		opts.Cascade = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("RevokeOwnershipOptions", "Restrict", "Cascade"))
	})

	t.Run("on future schema object in database to role", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE OWNERSHIP ON FUTURE TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on future schema object in schema to role", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = RevokeOwnershipGrantOn{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InSchema:         Pointer(schemaId),
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE OWNERSHIP ON FUTURE TABLES IN SCHEMA %s FROM ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on future schema object in database to database role", func(t *testing.T) {
		opts := defaultOpts()
		opts.From = OwnershipGrantTo{
			DatabaseRoleName: Pointer(databaseRoleId),
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE OWNERSHIP ON FUTURE TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object with cascade", func(t *testing.T) {
		opts := defaultOpts()
		opts.Cascade = Bool(true)
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE OWNERSHIP ON FUTURE TABLES IN DATABASE %s FROM ROLE %s CASCADE`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})
}

func TestGrantShow(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		opts := &showGrantsOptions{}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW GRANTS")
	})

	t.Run("on account", func(t *testing.T) {
		opts := &showGrantsOptions{
			On: &ShowGrantsOn{
				Account: Bool(true),
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW GRANTS ON ACCOUNT")
	})

	t.Run("on database", func(t *testing.T) {
		dbID := randomAccountObjectIdentifier()
		opts := &showGrantsOptions{
			On: &ShowGrantsOn{
				Object: &Object{
					ObjectType: ObjectTypeDatabase,
					Name:       dbID,
				},
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW GRANTS ON DATABASE %s", dbID.FullyQualifiedName())
	})

	t.Run("to role", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier()
		opts := &showGrantsOptions{
			To: &ShowGrantsTo{
				Role: roleID,
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW GRANTS TO ROLE %s", roleID.FullyQualifiedName())
	})

	t.Run("to user", func(t *testing.T) {
		userID := randomAccountObjectIdentifier()
		opts := &showGrantsOptions{
			To: &ShowGrantsTo{
				User: userID,
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW GRANTS TO USER %s", userID.FullyQualifiedName())
	})

	t.Run("to share", func(t *testing.T) {
		shareID := randomAccountObjectIdentifier()
		opts := &showGrantsOptions{
			To: &ShowGrantsTo{
				Share: &ShowGrantsToShare{
					Name: shareID,
				},
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW GRANTS TO SHARE %s", shareID.FullyQualifiedName())
	})

	t.Run("to share in application package", func(t *testing.T) {
		shareID := randomAccountObjectIdentifier()
		packageId := randomAccountObjectIdentifier()
		opts := &showGrantsOptions{
			To: &ShowGrantsTo{
				Share: &ShowGrantsToShare{
					Name:                 shareID,
					InApplicationPackage: &packageId,
				},
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW GRANTS TO SHARE %s IN APPLICATION PACKAGE %s", shareID.FullyQualifiedName(), packageId.FullyQualifiedName())
	})

	t.Run("of role", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier()
		opts := &showGrantsOptions{
			Of: &ShowGrantsOf{
				Role: roleID,
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW GRANTS OF ROLE %s", roleID.FullyQualifiedName())
	})

	t.Run("of database role", func(t *testing.T) {
		roleID := randomDatabaseObjectIdentifier()
		opts := &showGrantsOptions{
			Of: &ShowGrantsOf{
				DatabaseRole: roleID,
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW GRANTS OF DATABASE ROLE %s", roleID.FullyQualifiedName())
	})

	t.Run("of share", func(t *testing.T) {
		shareID := randomAccountObjectIdentifier()
		opts := &showGrantsOptions{
			Of: &ShowGrantsOf{
				Share: shareID,
			},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW GRANTS OF SHARE %s", shareID.FullyQualifiedName())
	})
}

func TestNormalizeShareGranteeName(t *testing.T) {
	accountLocator := "AB12345"
	testCases := []struct {
		name        string
		granteeName string
		expected    string
	}{
		{name: "account-prefixed share", granteeName: "AB12345.MY_SHARE", expected: "MY_SHARE"},
		{name: "account-prefixed dotted share", granteeName: `AB12345."MY.SHARE"`, expected: `"MY.SHARE"`},
		{name: "quoted account prefix is left unchanged", granteeName: `"AB12345"."MY.SHARE"`, expected: `"AB12345"."MY.SHARE"`},
		{name: "unprefixed dotted share", granteeName: "MY.SHARE", expected: "MY.SHARE"},
		{name: "case-insensitive account prefix", granteeName: "ab12345.MY_SHARE", expected: "MY_SHARE"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, normalizeShareGranteeName(testCase.granteeName, accountLocator))
		})
	}
}

// TestShowGrantsRequest_SQLKey covers the property the SHOW GRANTS cache relies on: the cache key
// (ShowGrantsRequest.SQLKey) must render exactly the statement Grants.Show executes, and distinct
// queries must not collide.
func TestShowGrantsRequest_SQLKey(t *testing.T) {
	dbId := randomAccountObjectIdentifier()

	t.Run("key equals the executed SHOW GRANTS statement", func(t *testing.T) {
		req := &ShowGrantsRequest{
			On: &ShowGrantsOnRequest{
				Object: &Object{ObjectType: ObjectTypeDatabase, Name: dbId},
			},
		}

		key, err := req.SQLKey()
		require.NoError(t, err)

		executed, err := structToSQL(req.toOpts())
		require.NoError(t, err)

		assert.Equal(t, executed, key)
		assert.Contains(t, key, "SHOW GRANTS ON DATABASE")
	})

	t.Run("distinct queries produce distinct keys (no collisions)", func(t *testing.T) {
		keyFor := func(req *ShowGrantsRequest) string {
			key, err := req.SQLKey()
			require.NoError(t, err)
			return key
		}

		onDatabase := keyFor(&ShowGrantsRequest{On: &ShowGrantsOnRequest{Object: &Object{ObjectType: ObjectTypeDatabase, Name: dbId}}})
		inDatabase := keyFor(&ShowGrantsRequest{In: &ShowGrantsInRequest{Database: &dbId}})
		inherited := keyFor(&ShowGrantsRequest{Inherited: Bool(true), On: &ShowGrantsOnRequest{Object: &Object{ObjectType: ObjectTypeDatabase, Name: dbId}}})
		future := keyFor(&ShowGrantsRequest{Future: Bool(true), On: &ShowGrantsOnRequest{Object: &Object{ObjectType: ObjectTypeDatabase, Name: dbId}}})

		// ON vs IN, plain vs INHERITED, and plain vs FUTURE must all differ.
		assert.NotEqual(t, onDatabase, inDatabase)
		assert.NotEqual(t, onDatabase, inherited)
		assert.NotEqual(t, onDatabase, future)
		assert.NotEmpty(t, inDatabase)
	})
}

func TestGrantInheritedPrivilegesToAccountRole(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	roleId := randomAccountObjectIdentifier()

	defaultOpts := func() *GrantInheritedPrivilegesToAccountRoleOptions {
		return &GrantInheritedPrivilegesToAccountRoleOptions{
			Privileges: InheritedAccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
			},
			OnAll:       PluralObjectTypeTables,
			In:          InheritedAccountRoleGrantIn{Database: new(dbId)},
			AccountRole: roleId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *GrantInheritedPrivilegesToAccountRoleOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: at least one of the fields [opts.Privileges.AllPrivileges opts.Privileges.AccountObjectPrivileges opts.Privileges.SchemaPrivileges opts.Privileges.SchemaObjectPrivileges] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedAccountRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantInheritedPrivilegesToAccountRoleOptions.Privileges", "AllPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: at least one of the fields [opts.Privileges.AllPrivileges opts.Privileges.AccountObjectPrivileges opts.Privileges.SchemaPrivileges opts.Privileges.SchemaObjectPrivileges] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedAccountRoleGrantPrivileges{
			AccountObjectPrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeOperate},
			SchemaObjectPrivileges:  []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantInheritedPrivilegesToAccountRoleOptions.Privileges", "AllPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: [opts.OnAll] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.OnAll = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GrantInheritedPrivilegesToAccountRoleOptions", "OnAll"))
	})

	t.Run("validation: at least one of the fields [opts.In.Account opts.In.Database opts.In.Schema] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedAccountRoleGrantIn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantInheritedPrivilegesToAccountRoleOptions.In", "Account", "Database", "Schema"))
	})

	t.Run("validation: at least one of the fields [opts.In.Account opts.In.Database opts.In.Schema] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedAccountRoleGrantIn{
			Account:  new(true),
			Database: new(dbId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantInheritedPrivilegesToAccountRoleOptions.In", "Account", "Database", "Schema"))
	})

	t.Run("validation: valid identifier for [opts.In.Database]", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedAccountRoleGrantIn{
			Database: new(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.In.Schema]", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedAccountRoleGrantIn{
			Schema: new(emptyDatabaseObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.AccountRole]", func(t *testing.T) {
		opts := defaultOpts()
		opts.AccountRole = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("on all tables in account", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedAccountRoleGrantIn{Account: new(true)}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT INHERITED SELECT ON ALL TABLES IN ACCOUNT TO ROLE %s`, roleId.FullyQualifiedName())
	})

	t.Run("on all tables in database", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT INHERITED SELECT ON ALL TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on all tables in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedAccountRoleGrantIn{Schema: new(schemaId)}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT INHERITED SELECT ON ALL TABLES IN SCHEMA %s TO ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("multiple privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedAccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeInsert},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT INHERITED SELECT, INSERT ON ALL TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedNameEscaped())
	})

	t.Run("all privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedAccountRoleGrantPrivileges{AllPrivileges: new(true)}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT INHERITED ALL PRIVILEGES ON ALL TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})
}

func TestRevokeInheritedPrivilegesFromAccountRole(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	roleId := randomAccountObjectIdentifier()

	defaultOpts := func() *RevokeInheritedPrivilegesFromAccountRoleOptions {
		return &RevokeInheritedPrivilegesFromAccountRoleOptions{
			Privileges: InheritedAccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
			},
			OnAll:       PluralObjectTypeTables,
			In:          InheritedAccountRoleGrantIn{Database: new(dbId)},
			AccountRole: roleId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *RevokeInheritedPrivilegesFromAccountRoleOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: at least one of the fields [opts.Privileges.AllPrivileges opts.Privileges.AccountObjectPrivileges opts.Privileges.SchemaPrivileges opts.Privileges.SchemaObjectPrivileges] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedAccountRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RevokeInheritedPrivilegesFromAccountRoleOptions.Privileges", "AllPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: at least one of the fields [opts.Privileges.AllPrivileges opts.Privileges.AccountObjectPrivileges opts.Privileges.SchemaPrivileges opts.Privileges.SchemaObjectPrivileges] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedAccountRoleGrantPrivileges{
			AccountObjectPrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeOperate},
			SchemaObjectPrivileges:  []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RevokeInheritedPrivilegesFromAccountRoleOptions.Privileges", "AllPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: [opts.OnAll] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.OnAll = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("RevokeInheritedPrivilegesFromAccountRoleOptions", "OnAll"))
	})

	t.Run("validation: at least one of the fields [opts.In.Account opts.In.Database opts.In.Schema] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedAccountRoleGrantIn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RevokeInheritedPrivilegesFromAccountRoleOptions.In", "Account", "Database", "Schema"))
	})

	t.Run("validation: at least one of the fields [opts.In.Account opts.In.Database opts.In.Schema] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedAccountRoleGrantIn{
			Account:  new(true),
			Database: new(dbId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RevokeInheritedPrivilegesFromAccountRoleOptions.In", "Account", "Database", "Schema"))
	})

	t.Run("validation: valid identifier for [opts.In.Database]", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedAccountRoleGrantIn{
			Database: new(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.In.Schema]", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedAccountRoleGrantIn{
			Schema: new(emptyDatabaseObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.AccountRole]", func(t *testing.T) {
		opts := defaultOpts()
		opts.AccountRole = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("on all tables in account", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedAccountRoleGrantIn{Account: new(true)}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE INHERITED SELECT ON ALL TABLES IN ACCOUNT FROM ROLE %s`, roleId.FullyQualifiedName())
	})

	t.Run("on all tables in database", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE INHERITED SELECT ON ALL TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on all tables in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedAccountRoleGrantIn{Schema: new(schemaId)}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE INHERITED SELECT ON ALL TABLES IN SCHEMA %s FROM ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("multiple privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedAccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeInsert},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE INHERITED SELECT, INSERT ON ALL TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedNameEscaped())
	})

	t.Run("all privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedAccountRoleGrantPrivileges{AllPrivileges: new(true)}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE INHERITED ALL PRIVILEGES ON ALL TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})
}

func TestGrantInheritedPrivilegesToDatabaseRole(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	databaseRoleId := randomDatabaseObjectIdentifier()

	defaultOpts := func() *GrantInheritedPrivilegesToDatabaseRoleOptions {
		return &GrantInheritedPrivilegesToDatabaseRoleOptions{
			Privileges: InheritedDatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
			},
			OnAll:        PluralObjectTypeTables,
			In:           InheritedDatabaseRoleGrantIn{Database: new(dbId)},
			DatabaseRole: databaseRoleId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *GrantInheritedPrivilegesToDatabaseRoleOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: at least one of the fields [opts.Privileges.AllPrivileges opts.Privileges.SchemaPrivileges opts.Privileges.SchemaObjectPrivileges] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedDatabaseRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantInheritedPrivilegesToDatabaseRoleOptions.Privileges", "AllPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: at least one of the fields [opts.Privileges.AllPrivileges opts.Privileges.SchemaPrivileges opts.Privileges.SchemaObjectPrivileges] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedDatabaseRoleGrantPrivileges{
			SchemaPrivileges:       []SchemaPrivilege{SchemaPrivilegeCreateTable},
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantInheritedPrivilegesToDatabaseRoleOptions.Privileges", "AllPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: [opts.OnAll] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.OnAll = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GrantInheritedPrivilegesToDatabaseRoleOptions", "OnAll"))
	})

	t.Run("validation: at least one of the fields [opts.In.Database opts.In.Schema] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedDatabaseRoleGrantIn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantInheritedPrivilegesToDatabaseRoleOptions.In", "Database", "Schema"))
	})

	t.Run("validation: at least one of the fields [opts.In.Database opts.In.Schema] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedDatabaseRoleGrantIn{
			Database: new(dbId),
			Schema:   new(schemaId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantInheritedPrivilegesToDatabaseRoleOptions.In", "Database", "Schema"))
	})

	t.Run("validation: valid identifier for [opts.In.Database]", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedDatabaseRoleGrantIn{
			Database: new(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.In.Schema]", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedDatabaseRoleGrantIn{
			Schema: new(emptyDatabaseObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.DatabaseRole]", func(t *testing.T) {
		opts := defaultOpts()
		opts.DatabaseRole = emptyDatabaseObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("on all tables in database", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT INHERITED SELECT ON ALL TABLES IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on all tables in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedDatabaseRoleGrantIn{Schema: new(schemaId)}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT INHERITED SELECT ON ALL TABLES IN SCHEMA %s TO DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("multiple privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedDatabaseRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeInsert},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT INHERITED SELECT, INSERT ON ALL TABLES IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedNameEscaped())
	})

	t.Run("all privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedDatabaseRoleGrantPrivileges{AllPrivileges: new(true)}
		assertOptsValidAndSqlEqualsf(t, opts, `GRANT INHERITED ALL PRIVILEGES ON ALL TABLES IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})
}

func TestRevokeInheritedPrivilegesFromDatabaseRole(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	databaseRoleId := randomDatabaseObjectIdentifier()

	defaultOpts := func() *RevokeInheritedPrivilegesFromDatabaseRoleOptions {
		return &RevokeInheritedPrivilegesFromDatabaseRoleOptions{
			Privileges: InheritedDatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
			},
			OnAll:        PluralObjectTypeTables,
			In:           InheritedDatabaseRoleGrantIn{Database: new(dbId)},
			DatabaseRole: databaseRoleId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *RevokeInheritedPrivilegesFromDatabaseRoleOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: at least one of the fields [opts.Privileges.AllPrivileges opts.Privileges.SchemaPrivileges opts.Privileges.SchemaObjectPrivileges] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedDatabaseRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RevokeInheritedPrivilegesFromDatabaseRoleOptions.Privileges", "AllPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: at least one of the fields [opts.Privileges.AllPrivileges opts.Privileges.SchemaPrivileges opts.Privileges.SchemaObjectPrivileges] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedDatabaseRoleGrantPrivileges{
			SchemaPrivileges:       []SchemaPrivilege{SchemaPrivilegeCreateTable},
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RevokeInheritedPrivilegesFromDatabaseRoleOptions.Privileges", "AllPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: [opts.OnAll] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.OnAll = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("RevokeInheritedPrivilegesFromDatabaseRoleOptions", "OnAll"))
	})

	t.Run("validation: at least one of the fields [opts.In.Database opts.In.Schema] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedDatabaseRoleGrantIn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RevokeInheritedPrivilegesFromDatabaseRoleOptions.In", "Database", "Schema"))
	})

	t.Run("validation: at least one of the fields [opts.In.Database opts.In.Schema] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedDatabaseRoleGrantIn{
			Database: new(dbId),
			Schema:   new(schemaId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("RevokeInheritedPrivilegesFromDatabaseRoleOptions.In", "Database", "Schema"))
	})

	t.Run("validation: valid identifier for [opts.In.Database]", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedDatabaseRoleGrantIn{
			Database: new(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.In.Schema]", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedDatabaseRoleGrantIn{
			Schema: new(emptyDatabaseObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.DatabaseRole]", func(t *testing.T) {
		opts := defaultOpts()
		opts.DatabaseRole = emptyDatabaseObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("on all tables in database", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE INHERITED SELECT ON ALL TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on all tables in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = InheritedDatabaseRoleGrantIn{Schema: new(schemaId)}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE INHERITED SELECT ON ALL TABLES IN SCHEMA %s FROM DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("multiple privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedDatabaseRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeInsert},
		}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE INHERITED SELECT, INSERT ON ALL TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedNameEscaped())
	})

	t.Run("all privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.Privileges = InheritedDatabaseRoleGrantPrivileges{AllPrivileges: new(true)}
		assertOptsValidAndSqlEqualsf(t, opts, `REVOKE INHERITED ALL PRIVILEGES ON ALL TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})
}

func TestObjectTypeFromShowGrants(t *testing.T) {
	tests := []struct {
		raw  string
		want ObjectType
	}{
		{raw: "", want: ""},
		{raw: "DATABASE", want: ObjectTypeDatabase},
		{raw: "EXTERNAL_VOLUME", want: ObjectTypeExternalVolume},
		{raw: "VOLUME", want: ObjectTypeExternalVolume},
		{raw: "POSTGRES", want: ObjectTypePostgresInstance},
		{raw: "POSTGRES_INSTANCE", want: ObjectTypePostgresInstance},
		{raw: "MODULE", want: ObjectTypeModel},
		{raw: "CORTEX_AGENT", want: ObjectTypeAgent},
		{raw: "CORTEX_AGENT_SERVER", want: ObjectTypeMcpServer},
		{raw: "QUALITY_MONITOR", want: ObjectTypeModelMonitor},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.raw), func(t *testing.T) {
			actual, _ := ObjectTypeFromShowGrants(tt.raw)
			assert.Equal(t, tt.want, actual)
		})
	}
}
