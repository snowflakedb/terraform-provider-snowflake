package sdk

import (
	"fmt"
	"testing"
)

func TestGrantPrivilegesToAccountRole(t *testing.T) {
	t.Run("validation: privilege with disallowed characters", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				GlobalPrivileges: []GlobalPrivilege{"MONITOR USAGE; SELECT"},
			},
			on: &AccountRoleGrantOn{
				Account: new(true),
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", "MONITOR USAGE; SELECT", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("privileges with certain special characters are allowed", func(t *testing.T) {
		schemaId := randomDatabaseObjectIdentifier()
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{"CREATE SNOWFLAKE.ML.ANOMALY_DETECTION", "applybudget"},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: new(schemaId),
				},
			},
			accountRole:     NewAccountObjectIdentifier("role1"),
			WithGrantOption: new(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE SNOWFLAKE.ML.ANOMALY_DETECTION, applybudget ON SCHEMA %s TO ROLE "role1" WITH GRANT OPTION`, schemaId.FullyQualifiedName())
	})

	t.Run("on account", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				GlobalPrivileges: []GlobalPrivilege{GlobalPrivilegeMonitorUsage, GlobalPrivilegeApplyTag},
			},
			on: &AccountRoleGrantOn{
				Account: Bool(true),
			},
			accountRole:     NewAccountObjectIdentifier("role1"),
			WithGrantOption: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT MONITOR USAGE, APPLY TAG ON ACCOUNT TO ROLE "role1" WITH GRANT OPTION`)
	})

	t.Run("on account object", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			on: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Database: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT ALL PRIVILEGES ON DATABASE "db1" TO ROLE "role1"`)
	})

	t.Run("on account object - external volume", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			on: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					ExternalVolume: Pointer(NewAccountObjectIdentifier("ex volume")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT ALL PRIVILEGES ON EXTERNAL VOLUME "ex volume" TO ROLE "role1"`)
	})

	t.Run("on account object - compute pool", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			on: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					ComputePool: Pointer(NewAccountObjectIdentifier("compute pool")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT ALL PRIVILEGES ON COMPUTE POOL "compute pool" TO ROLE "role1"`)
	})

	t.Run("on account object - connection", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			on: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Connection: Pointer(NewAccountObjectIdentifier("myconn")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT ALL PRIVILEGES ON CONNECTION "myconn" TO ROLE "role1"`)
	})

	t.Run("on account object - exactly one of validation", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			on: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Database:    Pointer(NewAccountObjectIdentifier("database")),
					ComputePool: Pointer(NewAccountObjectIdentifier("pool")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("GrantOnAccountObject", "User", "ResourceMonitor", "Warehouse", "ComputePool", "Database", "Integration", "Connection", "FailoverGroup", "ReplicationGroup", "ExternalVolume", "SnowflakeIntelligence"))
	})

	t.Run("on account object - exactly one of validation - empty options", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			on: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("GrantOnAccountObject", "User", "ResourceMonitor", "Warehouse", "ComputePool", "Database", "Integration", "Connection", "FailoverGroup", "ReplicationGroup", "ExternalVolume", "SnowflakeIntelligence"))
	})

	t.Run("on schema", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(id),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON SCHEMA %s TO ROLE "role1"`, id.FullyQualifiedName())
	})

	t.Run("on all schemas in database", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					AllSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON ALL SCHEMAS IN DATABASE "db1" TO ROLE "role1"`)
	})

	t.Run("on all future schemas in database", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					FutureSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON FUTURE SCHEMAS IN DATABASE "db1" TO ROLE "role1"`)
	})

	t.Run("on schema object", func(t *testing.T) {
		tableId := randomSchemaObjectIdentifier()
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			on: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{
						ObjectType: ObjectTypeTable,
						Name:       tableId,
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON TABLE %s TO ROLE "role1"`, tableId.FullyQualifiedName())
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			on: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InDatabase:       Pointer(NewAccountObjectIdentifier("db1")),
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON FUTURE TABLES IN DATABASE "db1" TO ROLE "role1"`)
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := &GrantPrivilegesToAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			on: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InSchema:         Pointer(id),
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON FUTURE TABLES IN SCHEMA %s TO ROLE "role1"`, id.FullyQualifiedName())
	})
}

func TestRevokePrivilegesFromAccountRole(t *testing.T) {
	schemaId := randomDatabaseObjectIdentifier()

	t.Run("validation: privilege with disallowed characters", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				GlobalPrivileges: []GlobalPrivilege{"MONITOR USAGE; SELECT"},
			},
			on: &AccountRoleGrantOn{
				Account: new(true),
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", "MONITOR USAGE; SELECT", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("on account", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				GlobalPrivileges: []GlobalPrivilege{GlobalPrivilegeMonitorUsage, GlobalPrivilegeApplyTag},
			},
			on: &AccountRoleGrantOn{
				Account: Bool(true),
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE MONITOR USAGE, APPLY TAG ON ACCOUNT FROM ROLE "role1"`)
	})

	t.Run("on account object", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				AllPrivileges: Bool(true),
			},
			on: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Database: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE ALL PRIVILEGES ON DATABASE "db1" FROM ROLE "role1"`)
	})

	t.Run("on account object", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				AccountObjectPrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateDatabaseRole, AccountObjectPrivilegeModify},
			},
			on: &AccountRoleGrantOn{
				AccountObject: &GrantOnAccountObject{
					Database: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE DATABASE ROLE, MODIFY ON DATABASE "db1" FROM ROLE "role1"`)
	})

	t.Run("on schema", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(schemaId),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON SCHEMA %s FROM ROLE "role1"`, schemaId.FullyQualifiedName())
	})

	t.Run("on all schemas in database + restrict", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					AllSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
			Restrict:    Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON ALL SCHEMAS IN DATABASE "db1" FROM ROLE "role1" RESTRICT`)
	})

	t.Run("on all future schemas in database + cascade", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			on: &AccountRoleGrantOn{
				Schema: &GrantOnSchema{
					FutureSchemasInDatabase: Pointer(NewAccountObjectIdentifier("db1")),
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
			Cascade:     Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON FUTURE SCHEMAS IN DATABASE "db1" FROM ROLE "role1" CASCADE`)
	})

	t.Run("on schema object", func(t *testing.T) {
		tableId := randomSchemaObjectIdentifier()
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
			},
			on: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{
						ObjectType: ObjectTypeTable,
						Name:       tableId,
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE SELECT, UPDATE ON TABLE %s FROM ROLE "role1"`, tableId.FullyQualifiedName())
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
			},
			on: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InDatabase:       Pointer(NewAccountObjectIdentifier("db1")),
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE SELECT, UPDATE ON FUTURE TABLES IN DATABASE "db1" FROM ROLE "role1"`)
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := &RevokePrivilegesFromAccountRoleOptions{
			privileges: &AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
			},
			on: &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					Future: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InSchema:         Pointer(id),
					},
				},
			},
			accountRole: NewAccountObjectIdentifier("role1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE SELECT, UPDATE ON FUTURE TABLES IN SCHEMA %s FROM ROLE "role1"`, id.FullyQualifiedName())
	})
}

func TestGrants_GrantPrivilegesToDatabaseRole(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	databaseRoleId := randomDatabaseObjectIdentifierInDatabase(dbId)
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)

	defaultGrantsForDb := func() *GrantPrivilegesToDatabaseRoleOptions {
		return &GrantPrivilegesToDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			},
			on: &DatabaseRoleGrantOn{
				Database: &dbId,
			},
			databaseRole: databaseRoleId,
		}
	}

	defaultGrantsForSchema := func() *GrantPrivilegesToDatabaseRoleOptions {
		return &GrantPrivilegesToDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert},
			},
			on: &DatabaseRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(schemaId),
				},
			},
			databaseRole: databaseRoleId,
		}
	}
	tableId := randomSchemaObjectIdentifier()
	defaultGrantsForSchemaObject := func() *GrantPrivilegesToDatabaseRoleOptions {
		return &GrantPrivilegesToDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply},
			},
			on: &DatabaseRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{
						ObjectType: ObjectTypeTable,
						Name:       tableId,
					},
				},
			},
			databaseRole: databaseRoleId,
		}
	}

	t.Run("validation: privilege with disallowed characters", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = &DatabaseRoleGrantPrivileges{
			DatabasePrivileges: []AccountObjectPrivilege{"CREATE SCHEMA--"},
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", "CREATE SCHEMA--", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("validation: nil privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GrantPrivilegesToDatabaseRoleOptions", "privileges"))
	})

	t.Run("validation: no privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = &DatabaseRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantPrivileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges", "AllPrivileges"))
	})

	t.Run("validation: too many privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = &DatabaseRoleGrantPrivileges{
			DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			SchemaPrivileges:   []SchemaPrivilege{SchemaPrivilegeCreateAlert},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantPrivileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges", "AllPrivileges"))
	})

	t.Run("validation: no on set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.on = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("GrantPrivilegesToDatabaseRoleOptions", "on"))
	})

	t.Run("validation: no on set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.on = &DatabaseRoleGrantOn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantOn", "Database", "Schema", "SchemaObject"))
	})

	t.Run("validation: too many ons set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.on = &DatabaseRoleGrantOn{
			Database: &dbId,
			Schema: &GrantOnSchema{
				Schema: Pointer(schemaId),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantOn", "Database", "Schema", "SchemaObject"))
	})

	t.Run("validation: grant on schema", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.on.Schema = &GrantOnSchema{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchema", "Schema", "AllSchemasInDatabase", "FutureSchemasInDatabase"))
	})

	t.Run("validation: grant on schema object", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on.SchemaObject = &GrantOnSchemaObject{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObject", "SchemaObject", "All", "Future"))
	})

	t.Run("validation: grant on schema object - all", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on = &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				All: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema"))
	})

	t.Run("validation: grant on schema object - future", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on = &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				Future: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema"))
	})

	t.Run("on database", func(t *testing.T) {
		opts := defaultGrantsForDb()
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE SCHEMA ON DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on schema", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON SCHEMA %s TO DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on all schemas in database", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.on.Schema = &GrantOnSchema{
			AllSchemasInDatabase: Pointer(dbId),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON ALL SCHEMAS IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on all future schemas in database", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.on.Schema = &GrantOnSchema{
			FutureSchemasInDatabase: Pointer(dbId),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT CREATE ALERT ON FUTURE SCHEMAS IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on schema object", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON TABLE %s TO DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on.SchemaObject = &GrantOnSchemaObject{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InDatabase:       Pointer(dbId),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON FUTURE TABLES IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on.SchemaObject = &GrantOnSchemaObject{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InSchema:         Pointer(schemaId),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT APPLY ON FUTURE TABLES IN SCHEMA %s TO DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("grant all privileges", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.privileges = &DatabaseRoleGrantPrivileges{
			AllPrivileges: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT ALL PRIVILEGES ON TABLE %s TO DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})
}

func TestGrants_RevokePrivilegesFromDatabaseRoleRole(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	databaseRoleId := randomDatabaseObjectIdentifierInDatabase(dbId)
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	tableId := randomSchemaObjectIdentifierInSchema(schemaId)

	defaultGrantsForDb := func() *RevokePrivilegesFromDatabaseRoleOptions {
		return &RevokePrivilegesFromDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			},
			on: &DatabaseRoleGrantOn{
				Database: &dbId,
			},
			databaseRole: databaseRoleId,
		}
	}

	defaultGrantsForSchema := func() *RevokePrivilegesFromDatabaseRoleOptions {
		return &RevokePrivilegesFromDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
			},
			on: &DatabaseRoleGrantOn{
				Schema: &GrantOnSchema{
					Schema: Pointer(schemaId),
				},
			},
			databaseRole: databaseRoleId,
		}
	}

	defaultGrantsForSchemaObject := func() *RevokePrivilegesFromDatabaseRoleOptions {
		return &RevokePrivilegesFromDatabaseRoleOptions{
			privileges: &DatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
			},
			on: &DatabaseRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{
						ObjectType: ObjectTypeTable,
						Name:       tableId,
					},
				},
			},
			databaseRole: databaseRoleId,
		}
	}

	t.Run("validation: privilege with disallowed characters", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = &DatabaseRoleGrantPrivileges{
			DatabasePrivileges: []AccountObjectPrivilege{"CREATE SCHEMA--"},
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", "CREATE SCHEMA--", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("validation: nil privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("RevokePrivilegesFromDatabaseRoleOptions", "privileges"))
	})

	t.Run("validation: no privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = &DatabaseRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantPrivileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges", "AllPrivileges"))
	})

	t.Run("validation: too many privileges set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.privileges = &DatabaseRoleGrantPrivileges{
			DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema},
			SchemaPrivileges:   []SchemaPrivilege{SchemaPrivilegeCreateAlert},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantPrivileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges", "AllPrivileges"))
	})

	t.Run("validation: nil on set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.on = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("RevokePrivilegesFromDatabaseRoleOptions", "on"))
	})

	t.Run("validation: no on set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.on = &DatabaseRoleGrantOn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantOn", "Database", "Schema", "SchemaObject"))
	})

	t.Run("validation: too many ons set", func(t *testing.T) {
		opts := defaultGrantsForDb()
		opts.on = &DatabaseRoleGrantOn{
			Database: &dbId,
			Schema: &GrantOnSchema{
				Schema: Pointer(schemaId),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("DatabaseRoleGrantOn", "Database", "Schema", "SchemaObject"))
	})

	t.Run("validation: grant on schema", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.on.Schema = &GrantOnSchema{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchema", "Schema", "AllSchemasInDatabase", "FutureSchemasInDatabase"))
	})

	t.Run("validation: grant on schema object", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on.SchemaObject = &GrantOnSchemaObject{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObject", "SchemaObject", "All", "Future"))
	})

	t.Run("validation: grant on schema object - all", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on = &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				All: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema"))
	})

	t.Run("validation: grant on schema object - future", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on = &DatabaseRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				Future: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema"))
	})

	t.Run("on database", func(t *testing.T) {
		opts := defaultGrantsForDb()
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE SCHEMA ON DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on schema", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON SCHEMA %s FROM DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on all schemas in database + restrict", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.on.Schema = &GrantOnSchema{
			AllSchemasInDatabase: Pointer(dbId),
		}
		opts.Restrict = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON ALL SCHEMAS IN DATABASE %s FROM DATABASE ROLE %s RESTRICT`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on all future schemas in database + cascade", func(t *testing.T) {
		opts := defaultGrantsForSchema()
		opts.on.Schema = &GrantOnSchema{
			FutureSchemasInDatabase: Pointer(dbId),
		}
		opts.Cascade = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON FUTURE SCHEMAS IN DATABASE %s FROM DATABASE ROLE %s CASCADE`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on schema object", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		assertOptsValidAndSQLEquals(t, opts, `REVOKE SELECT, UPDATE ON TABLE %s FROM DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on.SchemaObject = &GrantOnSchemaObject{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InDatabase:       Pointer(dbId),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE SELECT, UPDATE ON FUTURE TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object in schema", func(t *testing.T) {
		opts := defaultGrantsForSchemaObject()
		opts.on.SchemaObject = &GrantOnSchemaObject{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InSchema:         Pointer(schemaId),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE SELECT, UPDATE ON FUTURE TABLES IN SCHEMA %s FROM DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})
}

func TestGrantPrivilegeToShare(t *testing.T) {
	id := randomAccountObjectIdentifier()
	t.Run("validation: privilege with disallowed characters", func(t *testing.T) {
		opts := &grantPrivilegeToShareOptions{
			privileges: []ObjectPrivilege{"USAGE;"},
			On: &ShareGrantOn{
				Database: randomAccountObjectIdentifier(),
			},
			to: id,
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", "USAGE;", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("on database", func(t *testing.T) {
		otherID := randomAccountObjectIdentifier()
		opts := &grantPrivilegeToShareOptions{
			privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Database: otherID,
			},
			to: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "GRANT USAGE ON DATABASE %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on schema", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier()
		opts := &grantPrivilegeToShareOptions{
			privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Schema: otherID,
			},
			to: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "GRANT USAGE ON SCHEMA %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on table", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier()
		opts := &grantPrivilegeToShareOptions{
			privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Table: &OnTable{
					Name: otherID,
				},
			},
			to: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "GRANT USAGE ON TABLE %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on all tables", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier()
		opts := &grantPrivilegeToShareOptions{
			privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Table: &OnTable{
					AllInSchema: otherID,
				},
			},
			to: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "GRANT USAGE ON ALL TABLES IN SCHEMA %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on view", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier()
		opts := &grantPrivilegeToShareOptions{
			privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				View: otherID,
			},
			to: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "GRANT USAGE ON VIEW %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})
}

func TestRevokePrivilegeFromShare(t *testing.T) {
	id := randomAccountObjectIdentifier()
	t.Run("validation: privilege with disallowed characters", func(t *testing.T) {
		opts := &revokePrivilegeFromShareOptions{
			privileges: []ObjectPrivilege{"USAGE;"},
			On: &ShareGrantOn{
				Database: randomAccountObjectIdentifier(),
			},
			from: id,
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", "USAGE;", allowedUnquotedCharactersRegex.String()))
	})

	t.Run("on database", func(t *testing.T) {
		otherID := randomAccountObjectIdentifier()
		opts := &revokePrivilegeFromShareOptions{
			privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Database: otherID,
			},
			from: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "REVOKE USAGE ON DATABASE %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on schema", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier()
		opts := &revokePrivilegeFromShareOptions{
			privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Schema: otherID,
			},
			from: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "REVOKE USAGE ON SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on table", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier()
		opts := &revokePrivilegeFromShareOptions{
			privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Table: &OnTable{
					Name: otherID,
				},
			},
			from: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "REVOKE USAGE ON TABLE %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on all tables", func(t *testing.T) {
		otherID := randomDatabaseObjectIdentifier()
		opts := &revokePrivilegeFromShareOptions{
			privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				Table: &OnTable{
					AllInSchema: otherID,
				},
			},
			from: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "REVOKE USAGE ON ALL TABLES IN SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on view", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier()
		opts := &revokePrivilegeFromShareOptions{
			privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
			On: &ShareGrantOn{
				View: otherID,
			},
			from: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "REVOKE USAGE ON VIEW %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
	})

	t.Run("on tag", func(t *testing.T) {
		tagId := randomSchemaObjectIdentifier()
		opts := &revokePrivilegeFromShareOptions{
			privileges: []ObjectPrivilege{ObjectPrivilegeRead},
			On: &ShareGrantOn{
				Tag: tagId,
			},
			from: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "REVOKE READ ON TAG %s FROM SHARE %s", tagId.FullyQualifiedName(), id.FullyQualifiedName())
	})
}

func TestGrants_GrantOwnership(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	roleId := randomAccountObjectIdentifier()
	databaseRoleId := randomDatabaseObjectIdentifierInDatabase(dbId)
	tableId := randomSchemaObjectIdentifierInSchema(schemaId)

	defaultOpts := func() *GrantOwnershipOptions {
		return &GrantOwnershipOptions{
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
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("OwnershipGrantOn", "Object", "AllIn", "Future"))
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
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("OwnershipGrantOn", "Object", "AllIn", "Future"))
	})

	t.Run("validation: grant on schema object - all", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			All: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema"))
	})

	t.Run("validation: grant on schema object - future", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema"))
	})

	t.Run("validation: grant to empty", func(t *testing.T) {
		opts := defaultOpts()
		opts.To = OwnershipGrantTo{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("OwnershipGrantTo", "databaseRoleName", "accountRoleName"))
	})

	t.Run("validation: grant to role and database role", func(t *testing.T) {
		opts := defaultOpts()
		opts.To = OwnershipGrantTo{
			DatabaseRoleName: Pointer(databaseRoleId),
			AccountRoleName:  Pointer(roleId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("OwnershipGrantTo", "databaseRoleName", "accountRoleName"))
	})

	t.Run("on schema object to role", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `GRANT OWNERSHIP ON TABLE %s TO ROLE %s`, tableId.FullyQualifiedName(), roleId.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, `GRANT OWNERSHIP ON DBT PROJECT %s TO ROLE %s`, dbtProjectId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on schema object to database role", func(t *testing.T) {
		opts := defaultOpts()
		opts.To = OwnershipGrantTo{
			DatabaseRoleName: Pointer(databaseRoleId),
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT OWNERSHIP ON TABLE %s TO DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object in database", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InDatabase:       Pointer(dbId),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT OWNERSHIP ON FUTURE TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on all schema objects in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = OwnershipGrantOn{
			All: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InSchema:         Pointer(schemaId),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT OWNERSHIP ON ALL TABLES IN SCHEMA %s TO ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on schema object with current grants", func(t *testing.T) {
		opts := defaultOpts()
		opts.CurrentGrants = &OwnershipCurrentGrants{
			OutboundPrivileges: Copy,
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT OWNERSHIP ON TABLE %s TO ROLE %s COPY CURRENT GRANTS`, tableId.FullyQualifiedName(), roleId.FullyQualifiedName())
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
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("RevokeOwnershipGrantOn", "Future"))
	})

	t.Run("validation: revoke from empty", func(t *testing.T) {
		opts := defaultOpts()
		opts.From = OwnershipGrantTo{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("OwnershipGrantTo", "databaseRoleName", "accountRoleName"))
	})

	t.Run("validation: revoke from role and database role", func(t *testing.T) {
		opts := defaultOpts()
		opts.From = OwnershipGrantTo{
			DatabaseRoleName: Pointer(databaseRoleId),
			AccountRoleName:  Pointer(roleId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("OwnershipGrantTo", "databaseRoleName", "accountRoleName"))
	})

	t.Run("validation: restrict and cascade", func(t *testing.T) {
		opts := defaultOpts()
		opts.Restrict = Bool(true)
		opts.Cascade = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("RevokeOwnershipOptions", "Restrict", "Cascade"))
	})

	t.Run("on future schema object in database to role", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `REVOKE OWNERSHIP ON FUTURE TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on future schema object in schema to role", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = RevokeOwnershipGrantOn{
			Future: &GrantOnSchemaObjectIn{
				PluralObjectType: PluralObjectTypeTables,
				InSchema:         Pointer(schemaId),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE OWNERSHIP ON FUTURE TABLES IN SCHEMA %s FROM ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on future schema object in database to database role", func(t *testing.T) {
		opts := defaultOpts()
		opts.From = OwnershipGrantTo{
			DatabaseRoleName: Pointer(databaseRoleId),
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE OWNERSHIP ON FUTURE TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on future schema object with cascade", func(t *testing.T) {
		opts := defaultOpts()
		opts.Cascade = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `REVOKE OWNERSHIP ON FUTURE TABLES IN DATABASE %s FROM ROLE %s CASCADE`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})
}

func TestGrantShow(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		opts := &ShowGrantOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS")
	})

	t.Run("on account", func(t *testing.T) {
		opts := &ShowGrantOptions{
			On: &ShowGrantsOn{
				Account: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS ON ACCOUNT")
	})

	t.Run("on database", func(t *testing.T) {
		dbID := randomAccountObjectIdentifier()
		opts := &ShowGrantOptions{
			On: &ShowGrantsOn{
				Object: &Object{
					ObjectType: ObjectTypeDatabase,
					Name:       dbID,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS ON DATABASE %s", dbID.FullyQualifiedName())
	})

	t.Run("to role", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier()
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS TO ROLE %s", roleID.FullyQualifiedName())
	})

	t.Run("to user", func(t *testing.T) {
		userID := randomAccountObjectIdentifier()
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				User: userID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS TO USER %s", userID.FullyQualifiedName())
	})

	t.Run("to share", func(t *testing.T) {
		shareID := randomAccountObjectIdentifier()
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				Share: &ShowGrantsToShare{
					Name: shareID,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS TO SHARE %s", shareID.FullyQualifiedName())
	})

	t.Run("to share in application package", func(t *testing.T) {
		shareID := randomAccountObjectIdentifier()
		packageId := randomAccountObjectIdentifier()
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				Share: &ShowGrantsToShare{
					Name:                 shareID,
					InApplicationPackage: &packageId,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS TO SHARE %s IN APPLICATION PACKAGE %s", shareID.FullyQualifiedName(), packageId.FullyQualifiedName())
	})

	t.Run("of role", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier()
		opts := &ShowGrantOptions{
			Of: &ShowGrantsOf{
				Role: roleID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS OF ROLE %s", roleID.FullyQualifiedName())
	})

	t.Run("of database role", func(t *testing.T) {
		roleID := randomDatabaseObjectIdentifier()
		opts := &ShowGrantOptions{
			Of: &ShowGrantsOf{
				DatabaseRole: roleID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS OF DATABASE ROLE %s", roleID.FullyQualifiedName())
	})

	t.Run("of share", func(t *testing.T) {
		shareID := randomAccountObjectIdentifier()
		opts := &ShowGrantOptions{
			Of: &ShowGrantsOf{
				Share: shareID,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GRANTS OF SHARE %s", shareID.FullyQualifiedName())
	})
}

func TestGrantInheritedPrivilegesToAccountRole(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	roleId := randomAccountObjectIdentifier()

	defaultOpts := func() *grantInheritedPrivilegesToAccountRoleOptions {
		return &grantInheritedPrivilegesToAccountRoleOptions{
			privileges: InheritedAccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
			},
			onAll:       PluralObjectTypeTables,
			in:          InheritedAccountRoleGrantIn{Database: new(dbId)},
			accountRole: roleId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *grantInheritedPrivilegesToAccountRoleOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: at least one of the fields [opts.privileges.AllPrivileges opts.privileges.AccountObjectPrivileges opts.privileges.SchemaPrivileges opts.privileges.SchemaObjectPrivileges] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedAccountRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedAccountRoleGrantPrivileges", "AllPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: at least one of the fields [opts.privileges.AllPrivileges opts.privileges.AccountObjectPrivileges opts.privileges.SchemaPrivileges opts.privileges.SchemaObjectPrivileges] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedAccountRoleGrantPrivileges{
			AccountObjectPrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeOperate},
			SchemaObjectPrivileges:  []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedAccountRoleGrantPrivileges", "AllPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: [opts.onAll] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.onAll = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("grantInheritedPrivilegesToAccountRoleOptions", "onAll"))
	})

	t.Run("validation: at least one of the fields [opts.in.Account opts.in.Database opts.in.Schema] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedAccountRoleGrantIn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedAccountRoleGrantIn", "Account", "Database", "Schema"))
	})

	t.Run("validation: at least one of the fields [opts.in.Account opts.in.Database opts.in.Schema] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedAccountRoleGrantIn{
			Account:  new(true),
			Database: new(dbId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedAccountRoleGrantIn", "Account", "Database", "Schema"))
	})

	t.Run("validation: valid identifier for [opts.in.Database]", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedAccountRoleGrantIn{
			Database: new(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.in.Schema]", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedAccountRoleGrantIn{
			Schema: new(emptyDatabaseObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.accountRole]", func(t *testing.T) {
		opts := defaultOpts()
		opts.accountRole = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("on all tables in account", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedAccountRoleGrantIn{Account: new(true)}
		assertOptsValidAndSQLEquals(t, opts, `GRANT INHERITED SELECT ON ALL TABLES IN ACCOUNT TO ROLE %s`, roleId.FullyQualifiedName())
	})

	t.Run("on all tables in database", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `GRANT INHERITED SELECT ON ALL TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on all tables in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedAccountRoleGrantIn{Schema: new(schemaId)}
		assertOptsValidAndSQLEquals(t, opts, `GRANT INHERITED SELECT ON ALL TABLES IN SCHEMA %s TO ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("multiple privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedAccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeInsert},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT INHERITED SELECT, INSERT ON ALL TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedNameEscaped())
	})

	t.Run("all privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedAccountRoleGrantPrivileges{AllPrivileges: new(true)}
		assertOptsValidAndSQLEquals(t, opts, `GRANT INHERITED ALL PRIVILEGES ON ALL TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})
}

func TestRevokeInheritedPrivilegesFromAccountRole(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	roleId := randomAccountObjectIdentifier()

	defaultOpts := func() *revokeInheritedPrivilegesFromAccountRoleOptions {
		return &revokeInheritedPrivilegesFromAccountRoleOptions{
			privileges: InheritedAccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
			},
			onAll:       PluralObjectTypeTables,
			in:          InheritedAccountRoleGrantIn{Database: new(dbId)},
			accountRole: roleId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *revokeInheritedPrivilegesFromAccountRoleOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: at least one of the fields [opts.privileges.AllPrivileges opts.privileges.AccountObjectPrivileges opts.privileges.SchemaPrivileges opts.privileges.SchemaObjectPrivileges] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedAccountRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedAccountRoleGrantPrivileges", "AllPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: at least one of the fields [opts.privileges.AllPrivileges opts.privileges.AccountObjectPrivileges opts.privileges.SchemaPrivileges opts.privileges.SchemaObjectPrivileges] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedAccountRoleGrantPrivileges{
			AccountObjectPrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeOperate},
			SchemaObjectPrivileges:  []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedAccountRoleGrantPrivileges", "AllPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: [opts.onAll] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.onAll = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("revokeInheritedPrivilegesFromAccountRoleOptions", "onAll"))
	})

	t.Run("validation: at least one of the fields [opts.in.Account opts.in.Database opts.in.Schema] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedAccountRoleGrantIn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedAccountRoleGrantIn", "Account", "Database", "Schema"))
	})

	t.Run("validation: at least one of the fields [opts.in.Account opts.in.Database opts.in.Schema] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedAccountRoleGrantIn{
			Account:  new(true),
			Database: new(dbId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedAccountRoleGrantIn", "Account", "Database", "Schema"))
	})

	t.Run("validation: valid identifier for [opts.in.Database]", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedAccountRoleGrantIn{
			Database: new(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.in.Schema]", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedAccountRoleGrantIn{
			Schema: new(emptyDatabaseObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.accountRole]", func(t *testing.T) {
		opts := defaultOpts()
		opts.accountRole = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("on all tables in account", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedAccountRoleGrantIn{Account: new(true)}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE INHERITED SELECT ON ALL TABLES IN ACCOUNT FROM ROLE %s`, roleId.FullyQualifiedName())
	})

	t.Run("on all tables in database", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `REVOKE INHERITED SELECT ON ALL TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("on all tables in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedAccountRoleGrantIn{Schema: new(schemaId)}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE INHERITED SELECT ON ALL TABLES IN SCHEMA %s FROM ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})

	t.Run("multiple privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedAccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeInsert},
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE INHERITED SELECT, INSERT ON ALL TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedNameEscaped())
	})

	t.Run("all privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedAccountRoleGrantPrivileges{AllPrivileges: new(true)}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE INHERITED ALL PRIVILEGES ON ALL TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName())
	})
}

func TestGrantInheritedPrivilegesToDatabaseRole(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	databaseRoleId := randomDatabaseObjectIdentifier()

	defaultOpts := func() *grantInheritedPrivilegesToDatabaseRoleOptions {
		return &grantInheritedPrivilegesToDatabaseRoleOptions{
			privileges: InheritedDatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
			},
			onAll:        PluralObjectTypeTables,
			in:           InheritedDatabaseRoleGrantIn{Database: new(dbId)},
			databaseRole: databaseRoleId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *grantInheritedPrivilegesToDatabaseRoleOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: at least one of the fields [opts.privileges.AllPrivileges opts.privileges.SchemaPrivileges opts.privileges.SchemaObjectPrivileges] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedDatabaseRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedDatabaseRoleGrantPrivileges", "AllPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: at least one of the fields [opts.privileges.AllPrivileges opts.privileges.SchemaPrivileges opts.privileges.SchemaObjectPrivileges] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedDatabaseRoleGrantPrivileges{
			SchemaPrivileges:       []SchemaPrivilege{SchemaPrivilegeCreateTable},
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedDatabaseRoleGrantPrivileges", "AllPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: [opts.onAll] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.onAll = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("grantInheritedPrivilegesToDatabaseRoleOptions", "onAll"))
	})

	t.Run("validation: at least one of the fields [opts.in.Database opts.in.Schema] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedDatabaseRoleGrantIn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedDatabaseRoleGrantIn", "Database", "Schema"))
	})

	t.Run("validation: at least one of the fields [opts.in.Database opts.in.Schema] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedDatabaseRoleGrantIn{
			Database: new(dbId),
			Schema:   new(schemaId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedDatabaseRoleGrantIn", "Database", "Schema"))
	})

	t.Run("validation: valid identifier for [opts.in.Database]", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedDatabaseRoleGrantIn{
			Database: new(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.in.Schema]", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedDatabaseRoleGrantIn{
			Schema: new(emptyDatabaseObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.databaseRole]", func(t *testing.T) {
		opts := defaultOpts()
		opts.databaseRole = emptyDatabaseObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("on all tables in database", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `GRANT INHERITED SELECT ON ALL TABLES IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on all tables in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedDatabaseRoleGrantIn{Schema: new(schemaId)}
		assertOptsValidAndSQLEquals(t, opts, `GRANT INHERITED SELECT ON ALL TABLES IN SCHEMA %s TO DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("multiple privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedDatabaseRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeInsert},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT INHERITED SELECT, INSERT ON ALL TABLES IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedNameEscaped())
	})

	t.Run("all privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedDatabaseRoleGrantPrivileges{AllPrivileges: new(true)}
		assertOptsValidAndSQLEquals(t, opts, `GRANT INHERITED ALL PRIVILEGES ON ALL TABLES IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})
}

func TestRevokeInheritedPrivilegesFromDatabaseRole(t *testing.T) {
	dbId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	databaseRoleId := randomDatabaseObjectIdentifier()

	defaultOpts := func() *revokeInheritedPrivilegesFromDatabaseRoleOptions {
		return &revokeInheritedPrivilegesFromDatabaseRoleOptions{
			privileges: InheritedDatabaseRoleGrantPrivileges{
				SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
			},
			onAll:        PluralObjectTypeTables,
			in:           InheritedDatabaseRoleGrantIn{Database: new(dbId)},
			databaseRole: databaseRoleId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *revokeInheritedPrivilegesFromDatabaseRoleOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: at least one of the fields [opts.privileges.AllPrivileges opts.privileges.SchemaPrivileges opts.privileges.SchemaObjectPrivileges] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedDatabaseRoleGrantPrivileges{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedDatabaseRoleGrantPrivileges", "AllPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: at least one of the fields [opts.privileges.AllPrivileges opts.privileges.SchemaPrivileges opts.privileges.SchemaObjectPrivileges] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedDatabaseRoleGrantPrivileges{
			SchemaPrivileges:       []SchemaPrivilege{SchemaPrivilegeCreateTable},
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedDatabaseRoleGrantPrivileges", "AllPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges"))
	})

	t.Run("validation: [opts.onAll] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.onAll = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("revokeInheritedPrivilegesFromDatabaseRoleOptions", "onAll"))
	})

	t.Run("validation: at least one of the fields [opts.in.Database opts.in.Schema] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedDatabaseRoleGrantIn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedDatabaseRoleGrantIn", "Database", "Schema"))
	})

	t.Run("validation: at least one of the fields [opts.in.Database opts.in.Schema] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedDatabaseRoleGrantIn{
			Database: new(dbId),
			Schema:   new(schemaId),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("InheritedDatabaseRoleGrantIn", "Database", "Schema"))
	})

	t.Run("validation: valid identifier for [opts.in.Database]", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedDatabaseRoleGrantIn{
			Database: new(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.in.Schema]", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedDatabaseRoleGrantIn{
			Schema: new(emptyDatabaseObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.databaseRole]", func(t *testing.T) {
		opts := defaultOpts()
		opts.databaseRole = emptyDatabaseObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("on all tables in database", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `REVOKE INHERITED SELECT ON ALL TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("on all tables in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.in = InheritedDatabaseRoleGrantIn{Schema: new(schemaId)}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE INHERITED SELECT ON ALL TABLES IN SCHEMA %s FROM DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("multiple privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedDatabaseRoleGrantPrivileges{
			SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeInsert},
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE INHERITED SELECT, INSERT ON ALL TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedNameEscaped())
	})

	t.Run("all privileges", func(t *testing.T) {
		opts := defaultOpts()
		opts.privileges = InheritedDatabaseRoleGrantPrivileges{AllPrivileges: new(true)}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE INHERITED ALL PRIVILEGES ON ALL TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})
}
