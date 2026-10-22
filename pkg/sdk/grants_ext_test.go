package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	roleId := randomAccountObjectIdentifier()
	dbId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifierInDatabase(dbId)
	databaseRoleId := randomDatabaseObjectIdentifierInDatabase(dbId)
	tableId := randomSchemaObjectIdentifierInSchema(schemaId)
	shareId := randomAccountObjectIdentifier()
	packageId := randomAccountObjectIdentifier()
	userId := randomAccountObjectIdentifier()
	tagId := randomSchemaObjectIdentifier()
	dbtProjectId := randomSchemaObjectIdentifierInSchema(schemaId)
	externalVolumeId := NewAccountObjectIdentifier("ex volume")
	computePoolId := NewAccountObjectIdentifier("compute pool")
	connectionId := NewAccountObjectIdentifier("myconn")
	postgresId := NewAccountObjectIdentifier("pg1")

	accountObjectOn := func(objectType ObjectType, name ObjectIdentifier) *AccountRoleGrantOn {
		return &AccountRoleGrantOn{
			AccountObject: &GrantOnAccountObject{
				Object: &Object{ObjectType: objectType, Name: name},
			},
		}
	}
	schemaObjectOn := func() *AccountRoleGrantOn {
		return &AccountRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				SchemaObject: &Object{ObjectType: ObjectTypeTable, Name: tableId},
			},
		}
	}
	futureTablesInDatabaseOn := func() *AccountRoleGrantOn {
		return &AccountRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				Future: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
					InDatabase:       &dbId,
				},
			},
		}
	}
	futureTablesInSchemaOn := func() *AccountRoleGrantOn {
		return &AccountRoleGrantOn{
			SchemaObject: &GrantOnSchemaObject{
				Future: &GrantOnSchemaObjectIn{
					PluralObjectType: PluralObjectTypeTables,
					InSchema:         &schemaId,
				},
			},
		}
	}

	grantsTests.grantPrivilegesToAccountRole.
		withDefaultOpts(func() *grantPrivilegesToAccountRoleOptions {
			return &grantPrivilegesToAccountRoleOptions{
				Privileges:  &AccountRoleGrantPrivileges{GlobalPrivileges: []GlobalPrivilege{GlobalPrivilegeMonitorUsage}},
				On:          &AccountRoleGrantOn{Account: new(true)},
				AccountRole: roleId,
			}
		}).
		withModify(case_Grants_validation_grantPrivilegesToAccountRole_opts_Privileges_ExactlyOneValueSet_MoreThanOneSet, func(opts *grantPrivilegesToAccountRoleOptions) {
			opts.Privileges.AllPrivileges = new(true)
		}).
		withModify(case_Grants_validation_grantPrivilegesToAccountRole_opts_On_SchemaObject_ExactlyOneValueSet_MoreThanOneSet, func(opts *grantPrivilegesToAccountRoleOptions) {
			opts.On = &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{ObjectType: ObjectTypeTable, Name: tableId},
					All: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InDatabase:       &dbId,
					},
				},
			}
		}).
		withAdditionalValidationCase(
			"validation_grantPrivilegesToAccountRole_Privileges_invalidCharacters",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{GlobalPrivileges: []GlobalPrivilege{"MONITOR USAGE; SELECT"}}
			},
			invalidPrivilegeErr("MONITOR USAGE; SELECT"),
		).
		withAdditionalValidationCase(
			"validation_grantPrivilegesToAccountRole_On_AccountObject_Object_invalidObjectType",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{AllPrivileges: new(true)}
				opts.On = accountObjectOn(ObjectType("TABLE; DROP"), postgresId)
			},
			invalidObjectTypeErr("TABLE; DROP"),
		).
		withExpectedSqlf(
			case_Grants_sql_grantPrivilegesToAccountRole_basic,
			`GRANT MONITOR USAGE ON ACCOUNT TO ROLE %s`, roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_specialCharacters",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{
					SchemaPrivileges: []SchemaPrivilege{"CREATE SNOWFLAKE.ML.ANOMALY_DETECTION", "applybudget"},
				}
				opts.On = &AccountRoleGrantOn{Schema: &GrantOnSchema{Schema: &schemaId}}
				opts.WithGrantOption = new(true)
			},
			`GRANT CREATE SNOWFLAKE.ML.ANOMALY_DETECTION, applybudget ON SCHEMA %s TO ROLE %s WITH GRANT OPTION`,
			schemaId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_onAccount_withGrantOption",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{
					GlobalPrivileges: []GlobalPrivilege{GlobalPrivilegeMonitorUsage, GlobalPrivilegeApplyTag},
				}
				opts.WithGrantOption = new(true)
			},
			`GRANT MONITOR USAGE, APPLY TAG ON ACCOUNT TO ROLE %s WITH GRANT OPTION`, roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_onDatabase",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{AllPrivileges: new(true)}
				opts.On = accountObjectOn(ObjectTypeDatabase, dbId)
			},
			`GRANT ALL PRIVILEGES ON DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_onExternalVolume",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{AllPrivileges: new(true)}
				opts.On = accountObjectOn(ObjectTypeExternalVolume, externalVolumeId)
			},
			`GRANT ALL PRIVILEGES ON EXTERNAL VOLUME %s TO ROLE %s`, externalVolumeId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_onComputePool",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{AllPrivileges: new(true)}
				opts.On = accountObjectOn(ObjectTypeComputePool, computePoolId)
			},
			`GRANT ALL PRIVILEGES ON COMPUTE POOL %s TO ROLE %s`, computePoolId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_onConnection",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{AllPrivileges: new(true)}
				opts.On = accountObjectOn(ObjectTypeConnection, connectionId)
			},
			`GRANT ALL PRIVILEGES ON CONNECTION %s TO ROLE %s`, connectionId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_onPostgresInstance",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{AllPrivileges: new(true)}
				opts.On = accountObjectOn(ObjectTypePostgresInstance, postgresId)
			},
			`GRANT ALL PRIVILEGES ON POSTGRES INSTANCE %s TO ROLE %s`, postgresId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_onSchema",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert}}
				opts.On = &AccountRoleGrantOn{Schema: &GrantOnSchema{Schema: &schemaId}}
			},
			`GRANT CREATE ALERT ON SCHEMA %s TO ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_onAllSchemasInDatabase",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert}}
				opts.On = &AccountRoleGrantOn{Schema: &GrantOnSchema{AllSchemasInDatabase: &dbId}}
			},
			`GRANT CREATE ALERT ON ALL SCHEMAS IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_onFutureSchemasInDatabase",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert}}
				opts.On = &AccountRoleGrantOn{Schema: &GrantOnSchema{FutureSchemasInDatabase: &dbId}}
			},
			`GRANT CREATE ALERT ON FUTURE SCHEMAS IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_onSchemaObject",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply}}
				opts.On = schemaObjectOn()
			},
			`GRANT APPLY ON TABLE %s TO ROLE %s`, tableId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_onFutureSchemaObjectInDatabase",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply}}
				opts.On = futureTablesInDatabaseOn()
			},
			`GRANT APPLY ON FUTURE TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToAccountRole_onFutureSchemaObjectInSchema",
			func(opts *grantPrivilegesToAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply}}
				opts.On = futureTablesInSchemaOn()
			},
			`GRANT APPLY ON FUTURE TABLES IN SCHEMA %s TO ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		)

	grantsTests.revokePrivilegesFromAccountRole.
		withDefaultOpts(func() *revokePrivilegesFromAccountRoleOptions {
			return &revokePrivilegesFromAccountRoleOptions{
				Privileges:  &AccountRoleGrantPrivileges{GlobalPrivileges: []GlobalPrivilege{GlobalPrivilegeMonitorUsage}},
				On:          &AccountRoleGrantOn{Account: new(true)},
				AccountRole: roleId,
			}
		}).
		withModify(case_Grants_validation_revokePrivilegesFromAccountRole_opts_Privileges_ExactlyOneValueSet_MoreThanOneSet, func(opts *revokePrivilegesFromAccountRoleOptions) {
			opts.Privileges.AllPrivileges = new(true)
		}).
		withModify(case_Grants_validation_revokePrivilegesFromAccountRole_opts_On_SchemaObject_ExactlyOneValueSet_MoreThanOneSet, func(opts *revokePrivilegesFromAccountRoleOptions) {
			opts.On = &AccountRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{ObjectType: ObjectTypeTable, Name: tableId},
					All: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InDatabase:       &dbId,
					},
				},
			}
		}).
		withAdditionalValidationCase(
			"validation_revokePrivilegesFromAccountRole_Privileges_invalidCharacters",
			func(opts *revokePrivilegesFromAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{GlobalPrivileges: []GlobalPrivilege{"MONITOR USAGE; SELECT"}}
			},
			invalidPrivilegeErr("MONITOR USAGE; SELECT"),
		).
		withExpectedSqlf(
			case_Grants_sql_revokePrivilegesFromAccountRole_basic,
			`REVOKE MONITOR USAGE ON ACCOUNT FROM ROLE %s`, roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromAccountRole_onAccount_multiplePrivileges",
			func(opts *revokePrivilegesFromAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{
					GlobalPrivileges: []GlobalPrivilege{GlobalPrivilegeMonitorUsage, GlobalPrivilegeApplyTag},
				}
			},
			`REVOKE MONITOR USAGE, APPLY TAG ON ACCOUNT FROM ROLE %s`, roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromAccountRole_onDatabase_allPrivileges",
			func(opts *revokePrivilegesFromAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{AllPrivileges: new(true)}
				opts.On = accountObjectOn(ObjectTypeDatabase, dbId)
			},
			`REVOKE ALL PRIVILEGES ON DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromAccountRole_onDatabase",
			func(opts *revokePrivilegesFromAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{
					AccountObjectPrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateDatabaseRole, AccountObjectPrivilegeModify},
				}
				opts.On = accountObjectOn(ObjectTypeDatabase, dbId)
			},
			`REVOKE CREATE DATABASE ROLE, MODIFY ON DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromAccountRole_onSchema",
			func(opts *revokePrivilegesFromAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{
					SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
				}
				opts.On = &AccountRoleGrantOn{Schema: &GrantOnSchema{Schema: &schemaId}}
			},
			`REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON SCHEMA %s FROM ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromAccountRole_onAllSchemasInDatabase_restrict",
			func(opts *revokePrivilegesFromAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{
					SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
				}
				opts.On = &AccountRoleGrantOn{Schema: &GrantOnSchema{AllSchemasInDatabase: &dbId}}
				opts.Restrict = new(true)
			},
			`REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON ALL SCHEMAS IN DATABASE %s FROM ROLE %s RESTRICT`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromAccountRole_onFutureSchemasInDatabase_cascade",
			func(opts *revokePrivilegesFromAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{
					SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
				}
				opts.On = &AccountRoleGrantOn{Schema: &GrantOnSchema{FutureSchemasInDatabase: &dbId}}
				opts.Cascade = new(true)
			},
			`REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON FUTURE SCHEMAS IN DATABASE %s FROM ROLE %s CASCADE`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromAccountRole_onSchemaObject",
			func(opts *revokePrivilegesFromAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{
					SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
				}
				opts.On = schemaObjectOn()
			},
			`REVOKE SELECT, UPDATE ON TABLE %s FROM ROLE %s`, tableId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromAccountRole_onFutureSchemaObjectInDatabase",
			func(opts *revokePrivilegesFromAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{
					SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
				}
				opts.On = futureTablesInDatabaseOn()
			},
			`REVOKE SELECT, UPDATE ON FUTURE TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromAccountRole_onFutureSchemaObjectInSchema",
			func(opts *revokePrivilegesFromAccountRoleOptions) {
				opts.Privileges = &AccountRoleGrantPrivileges{
					SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
				}
				opts.On = futureTablesInSchemaOn()
			},
			`REVOKE SELECT, UPDATE ON FUTURE TABLES IN SCHEMA %s FROM ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		)

	grantsTests.grantPrivilegesToDatabaseRole.
		withDefaultOpts(func() *grantPrivilegesToDatabaseRoleOptions {
			return &grantPrivilegesToDatabaseRoleOptions{
				Privileges:   &DatabaseRoleGrantPrivileges{DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema}},
				On:           &DatabaseRoleGrantOn{Database: &dbId},
				DatabaseRole: databaseRoleId,
			}
		}).
		withModify(case_Grants_validation_grantPrivilegesToDatabaseRole_opts_Privileges_ExactlyOneValueSet_MoreThanOneSet, func(opts *grantPrivilegesToDatabaseRoleOptions) {
			opts.Privileges.SchemaPrivileges = []SchemaPrivilege{SchemaPrivilegeCreateAlert}
		}).
		withModify(case_Grants_validation_grantPrivilegesToDatabaseRole_opts_On_SchemaObject_ExactlyOneValueSet_MoreThanOneSet, func(opts *grantPrivilegesToDatabaseRoleOptions) {
			opts.On = &DatabaseRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{ObjectType: ObjectTypeTable, Name: tableId},
					All: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InDatabase:       &dbId,
					},
				},
			}
		}).
		withAdditionalValidationCase(
			"validation_grantPrivilegesToDatabaseRole_Privileges_invalidCharacters",
			func(opts *grantPrivilegesToDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{DatabasePrivileges: []AccountObjectPrivilege{"CREATE SCHEMA--"}}
			},
			invalidPrivilegeErr("CREATE SCHEMA--"),
		).
		withExpectedSqlf(
			case_Grants_sql_grantPrivilegesToDatabaseRole_basic,
			`GRANT CREATE SCHEMA ON DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToDatabaseRole_onSchema",
			func(opts *grantPrivilegesToDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert}}
				opts.On = &DatabaseRoleGrantOn{Schema: &GrantOnSchema{Schema: &schemaId}}
			},
			`GRANT CREATE ALERT ON SCHEMA %s TO DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToDatabaseRole_onAllSchemasInDatabase",
			func(opts *grantPrivilegesToDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert}}
				opts.On = &DatabaseRoleGrantOn{Schema: &GrantOnSchema{AllSchemasInDatabase: &dbId}}
			},
			`GRANT CREATE ALERT ON ALL SCHEMAS IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToDatabaseRole_onFutureSchemasInDatabase",
			func(opts *grantPrivilegesToDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert}}
				opts.On = &DatabaseRoleGrantOn{Schema: &GrantOnSchema{FutureSchemasInDatabase: &dbId}}
			},
			`GRANT CREATE ALERT ON FUTURE SCHEMAS IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToDatabaseRole_onSchemaObject",
			func(opts *grantPrivilegesToDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply}}
				opts.On = &DatabaseRoleGrantOn{SchemaObject: &GrantOnSchemaObject{SchemaObject: &Object{ObjectType: ObjectTypeTable, Name: tableId}}}
			},
			`GRANT APPLY ON TABLE %s TO DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToDatabaseRole_onFutureSchemaObjectInDatabase",
			func(opts *grantPrivilegesToDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply}}
				opts.On = &DatabaseRoleGrantOn{
					SchemaObject: &GrantOnSchemaObject{
						Future: &GrantOnSchemaObjectIn{PluralObjectType: PluralObjectTypeTables, InDatabase: &dbId},
					},
				}
			},
			`GRANT APPLY ON FUTURE TABLES IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToDatabaseRole_onFutureSchemaObjectInSchema",
			func(opts *grantPrivilegesToDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeApply}}
				opts.On = &DatabaseRoleGrantOn{
					SchemaObject: &GrantOnSchemaObject{
						Future: &GrantOnSchemaObjectIn{PluralObjectType: PluralObjectTypeTables, InSchema: &schemaId},
					},
				}
			},
			`GRANT APPLY ON FUTURE TABLES IN SCHEMA %s TO DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantPrivilegesToDatabaseRole_allPrivileges",
			func(opts *grantPrivilegesToDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{AllPrivileges: new(true)}
				opts.On = &DatabaseRoleGrantOn{SchemaObject: &GrantOnSchemaObject{SchemaObject: &Object{ObjectType: ObjectTypeTable, Name: tableId}}}
			},
			`GRANT ALL PRIVILEGES ON TABLE %s TO DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		)

	grantsTests.revokePrivilegesFromDatabaseRole.
		withDefaultOpts(func() *revokePrivilegesFromDatabaseRoleOptions {
			return &revokePrivilegesFromDatabaseRoleOptions{
				Privileges:   &DatabaseRoleGrantPrivileges{DatabasePrivileges: []AccountObjectPrivilege{AccountObjectPrivilegeCreateSchema}},
				On:           &DatabaseRoleGrantOn{Database: &dbId},
				DatabaseRole: databaseRoleId,
			}
		}).
		withModify(case_Grants_validation_revokePrivilegesFromDatabaseRole_opts_Privileges_ExactlyOneValueSet_MoreThanOneSet, func(opts *revokePrivilegesFromDatabaseRoleOptions) {
			opts.Privileges.SchemaPrivileges = []SchemaPrivilege{SchemaPrivilegeCreateAlert}
		}).
		withModify(case_Grants_validation_revokePrivilegesFromDatabaseRole_opts_On_SchemaObject_ExactlyOneValueSet_MoreThanOneSet, func(opts *revokePrivilegesFromDatabaseRoleOptions) {
			opts.On = &DatabaseRoleGrantOn{
				SchemaObject: &GrantOnSchemaObject{
					SchemaObject: &Object{ObjectType: ObjectTypeTable, Name: tableId},
					All: &GrantOnSchemaObjectIn{
						PluralObjectType: PluralObjectTypeTables,
						InDatabase:       &dbId,
					},
				},
			}
		}).
		withAdditionalValidationCase(
			"validation_revokePrivilegesFromDatabaseRole_Privileges_invalidCharacters",
			func(opts *revokePrivilegesFromDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{DatabasePrivileges: []AccountObjectPrivilege{"CREATE SCHEMA--"}}
			},
			invalidPrivilegeErr("CREATE SCHEMA--"),
		).
		withExpectedSqlf(
			case_Grants_sql_revokePrivilegesFromDatabaseRole_basic,
			`REVOKE CREATE SCHEMA ON DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromDatabaseRole_onSchema",
			func(opts *revokePrivilegesFromDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{
					SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
				}
				opts.On = &DatabaseRoleGrantOn{Schema: &GrantOnSchema{Schema: &schemaId}}
			},
			`REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON SCHEMA %s FROM DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromDatabaseRole_onAllSchemasInDatabase_restrict",
			func(opts *revokePrivilegesFromDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{
					SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
				}
				opts.On = &DatabaseRoleGrantOn{Schema: &GrantOnSchema{AllSchemasInDatabase: &dbId}}
				opts.Restrict = new(true)
			},
			`REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON ALL SCHEMAS IN DATABASE %s FROM DATABASE ROLE %s RESTRICT`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromDatabaseRole_onFutureSchemasInDatabase_cascade",
			func(opts *revokePrivilegesFromDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{
					SchemaPrivileges: []SchemaPrivilege{SchemaPrivilegeCreateAlert, SchemaPrivilegeAddSearchOptimization},
				}
				opts.On = &DatabaseRoleGrantOn{Schema: &GrantOnSchema{FutureSchemasInDatabase: &dbId}}
				opts.Cascade = new(true)
			},
			`REVOKE CREATE ALERT, ADD SEARCH OPTIMIZATION ON FUTURE SCHEMAS IN DATABASE %s FROM DATABASE ROLE %s CASCADE`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromDatabaseRole_onSchemaObject",
			func(opts *revokePrivilegesFromDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{
					SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
				}
				opts.On = &DatabaseRoleGrantOn{SchemaObject: &GrantOnSchemaObject{SchemaObject: &Object{ObjectType: ObjectTypeTable, Name: tableId}}}
			},
			`REVOKE SELECT, UPDATE ON TABLE %s FROM DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromDatabaseRole_onFutureSchemaObjectInDatabase",
			func(opts *revokePrivilegesFromDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{
					SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
				}
				opts.On = &DatabaseRoleGrantOn{
					SchemaObject: &GrantOnSchemaObject{
						Future: &GrantOnSchemaObjectIn{PluralObjectType: PluralObjectTypeTables, InDatabase: &dbId},
					},
				}
			},
			`REVOKE SELECT, UPDATE ON FUTURE TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_revokePrivilegesFromDatabaseRole_onFutureSchemaObjectInSchema",
			func(opts *revokePrivilegesFromDatabaseRoleOptions) {
				opts.Privileges = &DatabaseRoleGrantPrivileges{
					SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeUpdate},
				}
				opts.On = &DatabaseRoleGrantOn{
					SchemaObject: &GrantOnSchemaObject{
						Future: &GrantOnSchemaObjectIn{PluralObjectType: PluralObjectTypeTables, InSchema: &schemaId},
					},
				}
			},
			`REVOKE SELECT, UPDATE ON FUTURE TABLES IN SCHEMA %s FROM DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		)

	grantsTests.GrantPrivilegeToShare.
		withDefaultOpts(func() *GrantPrivilegeToShareOptions {
			return &GrantPrivilegeToShareOptions{
				Privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
				On:         &ShareGrantOn{Database: dbId},
				To:         shareId,
			}
		}).
		withAdditionalValidationCase(
			"validation_GrantPrivilegeToShare_Privileges_invalidCharacters",
			func(opts *GrantPrivilegeToShareOptions) {
				opts.Privileges = []ObjectPrivilege{"USAGE;"}
			},
			invalidPrivilegeErr("USAGE;"),
		).
		withExpectedSqlf(
			case_Grants_sql_GrantPrivilegeToShare_basic,
			`GRANT USAGE ON DATABASE %s TO SHARE %s`, dbId.FullyQualifiedName(), shareId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_GrantPrivilegeToShare_onSchema",
			func(opts *GrantPrivilegeToShareOptions) { opts.On = &ShareGrantOn{Schema: schemaId} },
			`GRANT USAGE ON SCHEMA %s TO SHARE %s`, schemaId.FullyQualifiedName(), shareId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_GrantPrivilegeToShare_onTable",
			func(opts *GrantPrivilegeToShareOptions) { opts.On = &ShareGrantOn{Table: &OnTable{Name: tableId}} },
			`GRANT USAGE ON TABLE %s TO SHARE %s`, tableId.FullyQualifiedName(), shareId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_GrantPrivilegeToShare_onAllTables",
			func(opts *GrantPrivilegeToShareOptions) {
				opts.On = &ShareGrantOn{Table: &OnTable{AllInSchema: schemaId}}
			},
			`GRANT USAGE ON ALL TABLES IN SCHEMA %s TO SHARE %s`, schemaId.FullyQualifiedName(), shareId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_GrantPrivilegeToShare_onView",
			func(opts *GrantPrivilegeToShareOptions) { opts.On = &ShareGrantOn{View: tableId} },
			`GRANT USAGE ON VIEW %s TO SHARE %s`, tableId.FullyQualifiedName(), shareId.FullyQualifiedName(),
		)

	grantsTests.RevokePrivilegeFromShare.
		withDefaultOpts(func() *RevokePrivilegeFromShareOptions {
			return &RevokePrivilegeFromShareOptions{
				Privileges: []ObjectPrivilege{ObjectPrivilegeUsage},
				On:         &ShareGrantOn{Database: dbId},
				From:       shareId,
			}
		}).
		withAdditionalValidationCase(
			"validation_RevokePrivilegeFromShare_Privileges_invalidCharacters",
			func(opts *RevokePrivilegeFromShareOptions) {
				opts.Privileges = []ObjectPrivilege{"USAGE;"}
			},
			invalidPrivilegeErr("USAGE;"),
		).
		withExpectedSqlf(
			case_Grants_sql_RevokePrivilegeFromShare_basic,
			`REVOKE USAGE ON DATABASE %s FROM SHARE %s`, dbId.FullyQualifiedName(), shareId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokePrivilegeFromShare_onSchema",
			func(opts *RevokePrivilegeFromShareOptions) { opts.On = &ShareGrantOn{Schema: schemaId} },
			`REVOKE USAGE ON SCHEMA %s FROM SHARE %s`, schemaId.FullyQualifiedName(), shareId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokePrivilegeFromShare_onTable",
			func(opts *RevokePrivilegeFromShareOptions) { opts.On = &ShareGrantOn{Table: &OnTable{Name: tableId}} },
			`REVOKE USAGE ON TABLE %s FROM SHARE %s`, tableId.FullyQualifiedName(), shareId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokePrivilegeFromShare_onAllTables",
			func(opts *RevokePrivilegeFromShareOptions) {
				opts.On = &ShareGrantOn{Table: &OnTable{AllInSchema: schemaId}}
			},
			`REVOKE USAGE ON ALL TABLES IN SCHEMA %s FROM SHARE %s`, schemaId.FullyQualifiedName(), shareId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokePrivilegeFromShare_onView",
			func(opts *RevokePrivilegeFromShareOptions) { opts.On = &ShareGrantOn{View: tableId} },
			`REVOKE USAGE ON VIEW %s FROM SHARE %s`, tableId.FullyQualifiedName(), shareId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokePrivilegeFromShare_onTag",
			func(opts *RevokePrivilegeFromShareOptions) {
				opts.Privileges = []ObjectPrivilege{ObjectPrivilegeRead}
				opts.On = &ShareGrantOn{Tag: tagId}
			},
			`REVOKE READ ON TAG %s FROM SHARE %s`, tagId.FullyQualifiedName(), shareId.FullyQualifiedName(),
		)

	grantsTests.grantOwnership.
		withDefaultOpts(func() *grantOwnershipOptions {
			return &grantOwnershipOptions{
				On: OwnershipGrantOn{Object: &Object{ObjectType: ObjectTypeTable, Name: tableId}},
				To: OwnershipGrantTo{AccountRoleName: &roleId},
			}
		}).
		withModify(case_Grants_validation_grantOwnership_opts_On_ExactlyOneValueSet_MoreThanOneSet, func(opts *grantOwnershipOptions) {
			opts.On.All = &GrantOnSchemaObjectIn{PluralObjectType: PluralObjectTypeTables, InDatabase: &dbId}
		}).
		withExpectedSqlf(
			case_Grants_sql_grantOwnership_basic,
			`GRANT OWNERSHIP ON TABLE %s TO ROLE %s`, tableId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantOwnership_onDbtProject",
			func(opts *grantOwnershipOptions) {
				opts.On = OwnershipGrantOn{Object: &Object{ObjectType: ObjectTypeDbtProject, Name: dbtProjectId}}
			},
			`GRANT OWNERSHIP ON DBT PROJECT %s TO ROLE %s`, dbtProjectId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantOwnership_toDatabaseRole",
			func(opts *grantOwnershipOptions) {
				opts.To = OwnershipGrantTo{DatabaseRoleName: &databaseRoleId}
			},
			`GRANT OWNERSHIP ON TABLE %s TO DATABASE ROLE %s`, tableId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantOwnership_onFutureInDatabase",
			func(opts *grantOwnershipOptions) {
				opts.On = OwnershipGrantOn{
					Future: &GrantOnSchemaObjectIn{PluralObjectType: PluralObjectTypeTables, InDatabase: &dbId},
				}
			},
			`GRANT OWNERSHIP ON FUTURE TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantOwnership_onAllInSchema",
			func(opts *grantOwnershipOptions) {
				opts.On = OwnershipGrantOn{
					All: &GrantOnSchemaObjectIn{PluralObjectType: PluralObjectTypeTables, InSchema: &schemaId},
				}
			},
			`GRANT OWNERSHIP ON ALL TABLES IN SCHEMA %s TO ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_grantOwnership_currentGrants",
			func(opts *grantOwnershipOptions) {
				opts.CurrentGrants = &OwnershipCurrentGrants{OutboundPrivileges: Copy}
			},
			`GRANT OWNERSHIP ON TABLE %s TO ROLE %s COPY CURRENT GRANTS`, tableId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		)

	grantsTests.RevokeOwnership.
		withDefaultOpts(func() *RevokeOwnershipOptions {
			return &RevokeOwnershipOptions{
				On: RevokeOwnershipGrantOn{
					Future: &GrantOnSchemaObjectIn{PluralObjectType: PluralObjectTypeTables, InDatabase: &dbId},
				},
				From: OwnershipGrantTo{AccountRoleName: &roleId},
			}
		}).
		// On/From are value structs, so valueSet is always true; emptying them surfaces the nested field checks instead.
		withExpectedErr(case_Grants_validation_RevokeOwnership_On_ValidateValueSet, errNotSet("RevokeOwnershipOptions.On", "Future")).
		withExpectedErr(case_Grants_validation_RevokeOwnership_From_ValidateValueSet, errExactlyOneOf("RevokeOwnershipOptions.From", "DatabaseRoleName", "AccountRoleName")).
		withExpectedSqlf(
			case_Grants_sql_RevokeOwnership_basic,
			`REVOKE OWNERSHIP ON FUTURE TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokeOwnership_onFutureInSchema",
			func(opts *RevokeOwnershipOptions) {
				opts.On = RevokeOwnershipGrantOn{
					Future: &GrantOnSchemaObjectIn{PluralObjectType: PluralObjectTypeTables, InSchema: &schemaId},
				}
			},
			`REVOKE OWNERSHIP ON FUTURE TABLES IN SCHEMA %s FROM ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokeOwnership_fromDatabaseRole",
			func(opts *RevokeOwnershipOptions) {
				opts.From = OwnershipGrantTo{DatabaseRoleName: &databaseRoleId}
			},
			`REVOKE OWNERSHIP ON FUTURE TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokeOwnership_cascade",
			func(opts *RevokeOwnershipOptions) { opts.Cascade = new(true) },
			`REVOKE OWNERSHIP ON FUTURE TABLES IN DATABASE %s FROM ROLE %s CASCADE`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		)

	grantsTests.GrantInheritedPrivilegesToAccountRole.
		withDefaultOpts(func() *GrantInheritedPrivilegesToAccountRoleOptions {
			return &GrantInheritedPrivilegesToAccountRoleOptions{
				Privileges:  InheritedAccountRoleGrantPrivileges{SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect}},
				OnAll:       PluralObjectTypeTables,
				In:          InheritedAccountRoleGrantIn{Database: &dbId},
				AccountRole: roleId,
			}
		}).
		withModify(case_Grants_validation_GrantInheritedPrivilegesToAccountRole_opts_Privileges_ExactlyOneValueSet_MoreThanOneSet, func(opts *GrantInheritedPrivilegesToAccountRoleOptions) {
			opts.Privileges.AccountObjectPrivileges = []AccountObjectPrivilege{AccountObjectPrivilegeOperate}
		}).
		withExpectedSqlf(
			case_Grants_sql_GrantInheritedPrivilegesToAccountRole_basic,
			`GRANT INHERITED SELECT ON ALL TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_GrantInheritedPrivilegesToAccountRole_inAccount",
			func(opts *GrantInheritedPrivilegesToAccountRoleOptions) {
				opts.In = InheritedAccountRoleGrantIn{Account: new(true)}
			},
			`GRANT INHERITED SELECT ON ALL TABLES IN ACCOUNT TO ROLE %s`, roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_GrantInheritedPrivilegesToAccountRole_inSchema",
			func(opts *GrantInheritedPrivilegesToAccountRoleOptions) {
				opts.In = InheritedAccountRoleGrantIn{Schema: &schemaId}
			},
			`GRANT INHERITED SELECT ON ALL TABLES IN SCHEMA %s TO ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_GrantInheritedPrivilegesToAccountRole_multiplePrivileges",
			func(opts *GrantInheritedPrivilegesToAccountRoleOptions) {
				opts.Privileges = InheritedAccountRoleGrantPrivileges{
					SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeInsert},
				}
			},
			`GRANT INHERITED SELECT, INSERT ON ALL TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_GrantInheritedPrivilegesToAccountRole_allPrivileges",
			func(opts *GrantInheritedPrivilegesToAccountRoleOptions) {
				opts.Privileges = InheritedAccountRoleGrantPrivileges{AllPrivileges: new(true)}
			},
			`GRANT INHERITED ALL PRIVILEGES ON ALL TABLES IN DATABASE %s TO ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		)

	grantsTests.RevokeInheritedPrivilegesFromAccountRole.
		withDefaultOpts(func() *RevokeInheritedPrivilegesFromAccountRoleOptions {
			return &RevokeInheritedPrivilegesFromAccountRoleOptions{
				Privileges:  InheritedAccountRoleGrantPrivileges{SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect}},
				OnAll:       PluralObjectTypeTables,
				In:          InheritedAccountRoleGrantIn{Database: &dbId},
				AccountRole: roleId,
			}
		}).
		withModify(case_Grants_validation_RevokeInheritedPrivilegesFromAccountRole_opts_Privileges_ExactlyOneValueSet_MoreThanOneSet, func(opts *RevokeInheritedPrivilegesFromAccountRoleOptions) {
			opts.Privileges.AccountObjectPrivileges = []AccountObjectPrivilege{AccountObjectPrivilegeOperate}
		}).
		withExpectedSqlf(
			case_Grants_sql_RevokeInheritedPrivilegesFromAccountRole_basic,
			`REVOKE INHERITED SELECT ON ALL TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokeInheritedPrivilegesFromAccountRole_inAccount",
			func(opts *RevokeInheritedPrivilegesFromAccountRoleOptions) {
				opts.In = InheritedAccountRoleGrantIn{Account: new(true)}
			},
			`REVOKE INHERITED SELECT ON ALL TABLES IN ACCOUNT FROM ROLE %s`, roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokeInheritedPrivilegesFromAccountRole_inSchema",
			func(opts *RevokeInheritedPrivilegesFromAccountRoleOptions) {
				opts.In = InheritedAccountRoleGrantIn{Schema: &schemaId}
			},
			`REVOKE INHERITED SELECT ON ALL TABLES IN SCHEMA %s FROM ROLE %s`, schemaId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokeInheritedPrivilegesFromAccountRole_multiplePrivileges",
			func(opts *RevokeInheritedPrivilegesFromAccountRoleOptions) {
				opts.Privileges = InheritedAccountRoleGrantPrivileges{
					SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeInsert},
				}
			},
			`REVOKE INHERITED SELECT, INSERT ON ALL TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokeInheritedPrivilegesFromAccountRole_allPrivileges",
			func(opts *RevokeInheritedPrivilegesFromAccountRoleOptions) {
				opts.Privileges = InheritedAccountRoleGrantPrivileges{AllPrivileges: new(true)}
			},
			`REVOKE INHERITED ALL PRIVILEGES ON ALL TABLES IN DATABASE %s FROM ROLE %s`, dbId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		)

	grantsTests.GrantInheritedPrivilegesToDatabaseRole.
		withDefaultOpts(func() *GrantInheritedPrivilegesToDatabaseRoleOptions {
			return &GrantInheritedPrivilegesToDatabaseRoleOptions{
				Privileges:   InheritedDatabaseRoleGrantPrivileges{SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect}},
				OnAll:        PluralObjectTypeTables,
				In:           InheritedDatabaseRoleGrantIn{Database: &dbId},
				DatabaseRole: databaseRoleId,
			}
		}).
		withModify(case_Grants_validation_GrantInheritedPrivilegesToDatabaseRole_opts_Privileges_ExactlyOneValueSet_MoreThanOneSet, func(opts *GrantInheritedPrivilegesToDatabaseRoleOptions) {
			opts.Privileges.SchemaPrivileges = []SchemaPrivilege{SchemaPrivilegeCreateTable}
		}).
		withExpectedSqlf(
			case_Grants_sql_GrantInheritedPrivilegesToDatabaseRole_basic,
			`GRANT INHERITED SELECT ON ALL TABLES IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_GrantInheritedPrivilegesToDatabaseRole_inSchema",
			func(opts *GrantInheritedPrivilegesToDatabaseRoleOptions) {
				opts.In = InheritedDatabaseRoleGrantIn{Schema: &schemaId}
			},
			`GRANT INHERITED SELECT ON ALL TABLES IN SCHEMA %s TO DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_GrantInheritedPrivilegesToDatabaseRole_multiplePrivileges",
			func(opts *GrantInheritedPrivilegesToDatabaseRoleOptions) {
				opts.Privileges = InheritedDatabaseRoleGrantPrivileges{
					SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeInsert},
				}
			},
			`GRANT INHERITED SELECT, INSERT ON ALL TABLES IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_GrantInheritedPrivilegesToDatabaseRole_allPrivileges",
			func(opts *GrantInheritedPrivilegesToDatabaseRoleOptions) {
				opts.Privileges = InheritedDatabaseRoleGrantPrivileges{AllPrivileges: new(true)}
			},
			`GRANT INHERITED ALL PRIVILEGES ON ALL TABLES IN DATABASE %s TO DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		)

	grantsTests.RevokeInheritedPrivilegesFromDatabaseRole.
		withDefaultOpts(func() *RevokeInheritedPrivilegesFromDatabaseRoleOptions {
			return &RevokeInheritedPrivilegesFromDatabaseRoleOptions{
				Privileges:   InheritedDatabaseRoleGrantPrivileges{SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect}},
				OnAll:        PluralObjectTypeTables,
				In:           InheritedDatabaseRoleGrantIn{Database: &dbId},
				DatabaseRole: databaseRoleId,
			}
		}).
		withModify(case_Grants_validation_RevokeInheritedPrivilegesFromDatabaseRole_opts_Privileges_ExactlyOneValueSet_MoreThanOneSet, func(opts *RevokeInheritedPrivilegesFromDatabaseRoleOptions) {
			opts.Privileges.SchemaPrivileges = []SchemaPrivilege{SchemaPrivilegeCreateTable}
		}).
		withExpectedSqlf(
			case_Grants_sql_RevokeInheritedPrivilegesFromDatabaseRole_basic,
			`REVOKE INHERITED SELECT ON ALL TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokeInheritedPrivilegesFromDatabaseRole_inSchema",
			func(opts *RevokeInheritedPrivilegesFromDatabaseRoleOptions) {
				opts.In = InheritedDatabaseRoleGrantIn{Schema: &schemaId}
			},
			`REVOKE INHERITED SELECT ON ALL TABLES IN SCHEMA %s FROM DATABASE ROLE %s`, schemaId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokeInheritedPrivilegesFromDatabaseRole_multiplePrivileges",
			func(opts *RevokeInheritedPrivilegesFromDatabaseRoleOptions) {
				opts.Privileges = InheritedDatabaseRoleGrantPrivileges{
					SchemaObjectPrivileges: []SchemaObjectPrivilege{SchemaObjectPrivilegeSelect, SchemaObjectPrivilegeInsert},
				}
			},
			`REVOKE INHERITED SELECT, INSERT ON ALL TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_RevokeInheritedPrivilegesFromDatabaseRole_allPrivileges",
			func(opts *RevokeInheritedPrivilegesFromDatabaseRoleOptions) {
				opts.Privileges = InheritedDatabaseRoleGrantPrivileges{AllPrivileges: new(true)}
			},
			`REVOKE INHERITED ALL PRIVILEGES ON ALL TABLES IN DATABASE %s FROM DATABASE ROLE %s`, dbId.FullyQualifiedName(), databaseRoleId.FullyQualifiedName(),
		)

	grantsTests.showGrants.
		withExpectedSql(
			case_Grants_sql_showGrants_basic,
			`SHOW GRANTS`,
		).
		withAdditionalSqlCasef(
			"sql_showGrants_onAccount",
			func(opts *showGrantsOptions) { opts.On = &ShowGrantsOn{Account: new(true)} },
			`SHOW GRANTS ON ACCOUNT`,
		).
		withAdditionalSqlCasef(
			"sql_showGrants_onDatabase",
			func(opts *showGrantsOptions) {
				opts.On = &ShowGrantsOn{Object: &Object{ObjectType: ObjectTypeDatabase, Name: dbId}}
			},
			`SHOW GRANTS ON DATABASE %s`, dbId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_showGrants_toRole",
			func(opts *showGrantsOptions) { opts.To = &ShowGrantsTo{Role: roleId} },
			`SHOW GRANTS TO ROLE %s`, roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_showGrants_toUser",
			func(opts *showGrantsOptions) { opts.To = &ShowGrantsTo{User: userId} },
			`SHOW GRANTS TO USER %s`, userId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_showGrants_toShare",
			func(opts *showGrantsOptions) { opts.To = &ShowGrantsTo{Share: &ShowGrantsToShare{Name: shareId}} },
			`SHOW GRANTS TO SHARE %s`, shareId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_showGrants_toShareInApplicationPackage",
			func(opts *showGrantsOptions) {
				opts.To = &ShowGrantsTo{Share: &ShowGrantsToShare{Name: shareId, InApplicationPackage: &packageId}}
			},
			`SHOW GRANTS TO SHARE %s IN APPLICATION PACKAGE %s`, shareId.FullyQualifiedName(), packageId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_showGrants_ofRole",
			func(opts *showGrantsOptions) { opts.Of = &ShowGrantsOf{Role: roleId} },
			`SHOW GRANTS OF ROLE %s`, roleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_showGrants_ofDatabaseRole",
			func(opts *showGrantsOptions) { opts.Of = &ShowGrantsOf{DatabaseRole: databaseRoleId} },
			`SHOW GRANTS OF DATABASE ROLE %s`, databaseRoleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_showGrants_ofShare",
			func(opts *showGrantsOptions) { opts.Of = &ShowGrantsOf{Share: shareId} },
			`SHOW GRANTS OF SHARE %s`, shareId.FullyQualifiedName(),
		)
}

func invalidPrivilegeErr(privilege string) error {
	return fmt.Errorf("invalid privilege: %s contains disallowed characters; it must follow this regex: %s", privilege, allowedUnquotedCharactersRegex.String())
}

func invalidObjectTypeErr(objectType string) error {
	return fmt.Errorf("invalid object type: %s contains disallowed characters; it must follow this regex: %s", objectType, allowedUnquotedCharactersRegex.String())
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

		assert.NotEqual(t, onDatabase, inDatabase)
		assert.NotEqual(t, onDatabase, inherited)
		assert.NotEqual(t, onDatabase, future)
		assert.NotEmpty(t, inDatabase)
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
