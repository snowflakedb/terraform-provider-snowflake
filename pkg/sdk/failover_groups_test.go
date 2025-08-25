package sdk

import (
	"testing"
)

func TestFailoverGroupsCreate(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		opts := &CreateFailoverGroupOptions{
			IfNotExists: Bool(true),
			name:        NewAccountObjectIdentifier("fg1"),
			objectTypes: []PluralObjectType{
				PluralObjectTypeShares,
				PluralObjectTypeDatabases,
			},
			AllowedDatabases: []AccountObjectIdentifier{
				NewAccountObjectIdentifier("db1"),
			},
			AllowedShares: []AccountObjectIdentifier{
				NewAccountObjectIdentifier("share1"),
			},
			allowedAccounts: []AccountIdentifier{
				NewAccountIdentifier("MY_ORG", "MY_ACCOUNT"),
			},
			IgnoreEditionCheck:  Bool(true),
			ReplicationSchedule: String("10 MINUTE"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE FAILOVER GROUP IF NOT EXISTS "fg1" OBJECT_TYPES = SHARES, DATABASES ALLOWED_DATABASES = "db1" ALLOWED_SHARES = "share1" ALLOWED_ACCOUNTS = "MY_ORG"."MY_ACCOUNT" IGNORE EDITION CHECK REPLICATION_SCHEDULE = '10 MINUTE'`)
	})

	t.Run("minimal", func(t *testing.T) {
		opts := &CreateFailoverGroupOptions{
			IfNotExists: Bool(true),
			name:        NewAccountObjectIdentifier("fg1"),
			objectTypes: []PluralObjectType{
				PluralObjectTypeRoles,
			},
			allowedAccounts: []AccountIdentifier{
				NewAccountIdentifier("MY_ORG", "MY_ACCOUNT"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE FAILOVER GROUP IF NOT EXISTS "fg1" OBJECT_TYPES = ROLES ALLOWED_ACCOUNTS = "MY_ORG"."MY_ACCOUNT"`)
	})
}

func TestCreateSecondaryReplicationGroup(t *testing.T) {
	opts := &CreateSecondaryReplicationGroupOptions{
		IfNotExists:          Bool(true),
		name:                 NewAccountObjectIdentifier("fg1"),
		primaryFailoverGroup: NewExternalObjectIdentifierFromFullyQualifiedName("myorg.myaccount.fg1"),
	}
	assertOptsValidAndSQLEquals(t, opts, `CREATE FAILOVER GROUP IF NOT EXISTS "fg1" AS REPLICA OF "myorg"."myaccount"."fg1"`)
}

func TestFailoverGroupAlterSource(t *testing.T) {
	id := NewAccountObjectIdentifier("fg1")

	t.Run("validate: empty unset", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name:  id,
			Unset: &FailoverGroupUnset{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("FailoverGroupUnset", "ReplicationSchedule"))
	})

	t.Run("unset replication schedule", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name: id,
			Unset: &FailoverGroupUnset{
				ReplicationSchedule: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FAILOVER GROUP "fg1" UNSET REPLICATION_SCHEDULE`)
	})

	t.Run("rename", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name:    id,
			NewName: NewAccountObjectIdentifier("myfg1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FAILOVER GROUP "fg1" RENAME TO "myfg1"`)
	})

	t.Run("set object types and replication schedule", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name: id,
			Set: &FailoverGroupSet{
				ObjectTypes:         []PluralObjectType{PluralObjectTypeShares},
				ReplicationSchedule: String("10 MINUTE"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FAILOVER GROUP "fg1" SET OBJECT_TYPES = SHARES REPLICATION_SCHEDULE = '10 MINUTE'`)
	})

	t.Run("add database account object", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name: id,
			Add: &FailoverGroupAdd{
				AllowedDatabases: []AccountObjectIdentifier{
					NewAccountObjectIdentifier("db1"),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FAILOVER GROUP "fg1" ADD "db1" TO ALLOWED_DATABASES`)
	})

	t.Run("remove database account object", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name: id,
			Remove: &FailoverGroupRemove{
				AllowedDatabases: []AccountObjectIdentifier{
					NewAccountObjectIdentifier("db1"),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FAILOVER GROUP "fg1" REMOVE "db1" FROM ALLOWED_DATABASES`)
	})

	t.Run("add allowed accounts", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name: id,
			Add: &FailoverGroupAdd{
				AllowedAccounts: []AccountIdentifier{
					NewAccountIdentifier("MY_ORG", "MY_ACCOUNT"),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FAILOVER GROUP "fg1" ADD "MY_ORG"."MY_ACCOUNT" TO ALLOWED_ACCOUNTS`)
	})

	t.Run("remove allowed accounts", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name: id,
			Remove: &FailoverGroupRemove{
				AllowedAccounts: []AccountIdentifier{
					NewAccountIdentifier("MY_ORG", "MY_ACCOUNT"),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FAILOVER GROUP "fg1" REMOVE "MY_ORG"."MY_ACCOUNT" FROM ALLOWED_ACCOUNTS`)
	})

	t.Run("add share account object", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name: id,
			Add: &FailoverGroupAdd{
				AllowedShares: []AccountObjectIdentifier{
					NewAccountObjectIdentifier("share1"),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FAILOVER GROUP "fg1" ADD "share1" TO ALLOWED_SHARES`)
	})

	t.Run("remove share account object", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name: id,
			Remove: &FailoverGroupRemove{
				AllowedShares: []AccountObjectIdentifier{
					NewAccountObjectIdentifier("share1"),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FAILOVER GROUP "fg1" REMOVE "share1" FROM ALLOWED_SHARES`)
	})

	t.Run("move shares to another failover group", func(t *testing.T) {
		opts := &AlterSourceFailoverGroupOptions{
			name: NewAccountObjectIdentifier("fg1"),
			Move: &FailoverGroupMove{
				Shares: []AccountObjectIdentifier{
					NewAccountObjectIdentifier("share1"),
				},
				To: NewAccountObjectIdentifier("fg2"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FAILOVER GROUP "fg1" MOVE SHARES "share1" TO FAILOVER GROUP "fg2"`)
	})
}

func TestFailoverGroupsAlterTarget(t *testing.T) {
	t.Run("resume", func(t *testing.T) {
		opts := &AlterTargetFailoverGroupOptions{
			name:   NewAccountObjectIdentifier("fg1"),
			Resume: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FAILOVER GROUP "fg1" RESUME`)
	})

	t.Run("primary", func(t *testing.T) {
		opts := &AlterTargetFailoverGroupOptions{
			name:    NewAccountObjectIdentifier("fg1"),
			Primary: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FAILOVER GROUP "fg1" PRIMARY`)
	})
}

func TestFailoverGroupsDrop(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &DropFailoverGroupOptions{
			name: NewAccountObjectIdentifier("fg1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP FAILOVER GROUP "fg1"`)
	})

	t.Run("with IfExists", func(t *testing.T) {
		opts := &DropFailoverGroupOptions{
			name:     NewAccountObjectIdentifier("fg1"),
			IfExists: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP FAILOVER GROUP IF EXISTS "fg1"`)
	})
}

func TestFailoverGroupsShow(t *testing.T) {
	t.Run("without show options", func(t *testing.T) {
		opts := &ShowFailoverGroupOptions{}
		assertOptsValidAndSQLEquals(t, opts, `SHOW FAILOVER GROUPS`)
	})

	t.Run("with show options", func(t *testing.T) {
		opts := &ShowFailoverGroupOptions{
			InAccount: NewAccountIdentifierFromAccountLocator("abcd123"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW FAILOVER GROUPS IN ACCOUNT "abcd123"`)
	})
}

func TestFailoverGroupsShowDatabases(t *testing.T) {
	opts := &showFailoverGroupDatabasesOptions{
		in: NewAccountObjectIdentifier("fg1"),
	}
	assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASES IN FAILOVER GROUP "fg1"`)
}

func TestFailoverGroupsShowShares(t *testing.T) {
	opts := &showFailoverGroupSharesOptions{
		in: NewAccountObjectIdentifier("fg1"),
	}
	assertOptsValidAndSQLEquals(t, opts, `SHOW SHARES IN FAILOVER GROUP "fg1"`)
}
