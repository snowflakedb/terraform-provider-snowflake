package sdk

import (
	"testing"
)

func TestDatabaseRoleCreate(t *testing.T) {
	id := randomDatabaseObjectIdentifier()

	defaultOpts := func() *createDatabaseRoleOptions {
		return &createDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createDatabaseRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyDatabaseObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: both ifNotExists and orReplace present", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("createDatabaseRoleOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: multiple errors", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyDatabaseObjectIdentifier
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier, errOneOf("createDatabaseRoleOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE ROLE %s`, id.FullyQualifiedName())
	})

	t.Run("all optional", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(false)
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE ROLE IF NOT EXISTS %s COMMENT = 'some comment'`, id.FullyQualifiedName())
	})
}

func TestDatabaseRoleAlter(t *testing.T) {
	id := randomDatabaseObjectIdentifier()

	defaultOpts := func() *alterDatabaseRoleOptions {
		return &alterDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *alterDatabaseRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyDatabaseObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("alterDatabaseRoleOptions", "Rename", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: multiple alter actions", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseRoleSet{
			Comment: "new comment",
		}
		opts.Unset = &DatabaseRoleUnset{
			Comment: true,
		}
		opts.SetTags = []TagAssociation{}
		opts.UnsetTags = []ObjectIdentifier{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("alterDatabaseRoleOptions", "Rename", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: invalid new name", func(t *testing.T) {
		opts := defaultOpts()
		opts.Rename = &emptyDatabaseObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("alterDatabaseRoleOptions", "Rename"))
	})

	t.Run("validation: new name from different db", func(t *testing.T) {
		newId := randomDatabaseObjectIdentifier()

		opts := defaultOpts()
		opts.Rename = &newId
		assertOptsInvalidJoinedErrors(t, opts, ErrDifferentDatabase)
	})

	t.Run("validation: no property to unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &DatabaseRoleUnset{
			Comment: false,
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("alterDatabaseRoleOptions.Unset", "Comment"))
	})

	t.Run("rename", func(t *testing.T) {
		newId := randomDatabaseObjectIdentifierInDatabase(id.DatabaseId())

		opts := defaultOpts()
		opts.Rename = &newId
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE %s RENAME TO %s`, id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Set = &DatabaseRoleSet{
			Comment: "new comment",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE IF EXISTS %s SET COMMENT = 'new comment'`, id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("123"),
				Value: "value-123",
			},
			{
				Name:  NewAccountObjectIdentifier("456"),
				Value: "value-123",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE IF EXISTS %s SET TAG "123" = 'value-123', "456" = 'value-123'`, id.FullyQualifiedName())
	})

	t.Run("set comment to empty", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Set = &DatabaseRoleSet{
			Comment: "",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE IF EXISTS %s SET COMMENT = ''`, id.FullyQualifiedName())
	})

	t.Run("unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Unset = &DatabaseRoleUnset{
			Comment: true,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE IF EXISTS %s UNSET COMMENT`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("123"),
			NewAccountObjectIdentifier("456"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE ROLE IF EXISTS %s UNSET TAG "123", "456"`, id.FullyQualifiedName())
	})
}

func TestDatabaseRoleDrop(t *testing.T) {
	id := randomDatabaseObjectIdentifier()

	defaultOpts := func() *dropDatabaseRoleOptions {
		return &dropDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *dropDatabaseRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyDatabaseObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE ROLE %s`, id.FullyQualifiedName())
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE ROLE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestDatabaseRolesShow(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *showDatabaseRoleOptions {
		return &showDatabaseRoleOptions{
			Database: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *showDatabaseRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.Database = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, ErrPatternRequiredForLikeKeyword)
	})

	t.Run("show", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASE ROLES IN DATABASE %s`, id.FullyQualifiedName())
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASE ROLES LIKE '%s' IN DATABASE %s`, id.Name(), id.FullyQualifiedName())
	})

	t.Run("show with like and limit from", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		opts.Limit = &LimitFrom{
			Rows: Int(10),
			From: String("name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASE ROLES LIKE '%s' IN DATABASE %s LIMIT 10 FROM 'name'`, id.Name(), id.FullyQualifiedName())
	})
}

func TestDatabaseRoles_Grant(t *testing.T) {
	id := randomDatabaseObjectIdentifier()
	databaseRoleId := randomDatabaseObjectIdentifier()
	accountRoleId := randomAccountObjectIdentifier()

	setUpOpts := func() *grantDatabaseRoleOptions {
		return &grantDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *grantDatabaseRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = emptyDatabaseObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no role", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("DatabaseRoleName", "AccountRoleName"))
	})

	t.Run("validation: multiple roles", func(t *testing.T) {
		opts := setUpOpts()
		opts.ParentRole.DatabaseRoleName = &databaseRoleId
		opts.ParentRole.AccountRoleName = &accountRoleId
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("DatabaseRoleName", "AccountRoleName"))
	})

	t.Run("grant to database role", func(t *testing.T) {
		opts := setUpOpts()
		opts.ParentRole.DatabaseRoleName = &databaseRoleId

		assertOptsValidAndSQLEquals(t, opts, `GRANT DATABASE ROLE %s TO DATABASE ROLE %s`, id.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("grant to account role", func(t *testing.T) {
		opts := setUpOpts()
		opts.ParentRole.AccountRoleName = &accountRoleId

		assertOptsValidAndSQLEquals(t, opts, `GRANT DATABASE ROLE %s TO ROLE %s`, id.FullyQualifiedName(), accountRoleId.FullyQualifiedName())
	})
}

func TestDatabaseRoles_Revoke(t *testing.T) {
	id := randomDatabaseObjectIdentifier()
	databaseRoleId := randomDatabaseObjectIdentifier()
	accountRoleId := randomAccountObjectIdentifier()

	setUpOpts := func() *revokeDatabaseRoleOptions {
		return &revokeDatabaseRoleOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *revokeDatabaseRoleOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = emptyDatabaseObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no role", func(t *testing.T) {
		opts := setUpOpts()
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("DatabaseRoleName", "AccountRoleName"))
	})

	t.Run("validation: multiple roles", func(t *testing.T) {
		opts := setUpOpts()
		opts.ParentRole.DatabaseRoleName = &databaseRoleId
		opts.ParentRole.AccountRoleName = &accountRoleId
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("DatabaseRoleName", "AccountRoleName"))
	})

	t.Run("revoke from database role", func(t *testing.T) {
		opts := setUpOpts()
		opts.ParentRole.DatabaseRoleName = &databaseRoleId

		assertOptsValidAndSQLEquals(t, opts, `REVOKE DATABASE ROLE %s FROM DATABASE ROLE %s`, id.FullyQualifiedName(), databaseRoleId.FullyQualifiedName())
	})

	t.Run("revoke from account role", func(t *testing.T) {
		opts := setUpOpts()
		opts.ParentRole.AccountRoleName = &accountRoleId

		assertOptsValidAndSQLEquals(t, opts, `REVOKE DATABASE ROLE %s FROM ROLE %s`, id.FullyQualifiedName(), accountRoleId.FullyQualifiedName())
	})
}

func TestDatabaseRoles_GrantToShare(t *testing.T) {
	id := randomDatabaseObjectIdentifier()
	share := randomAccountObjectIdentifier()

	setUpOpts := func() *grantDatabaseRoleToShareOptions {
		return &grantDatabaseRoleToShareOptions{
			name:  id,
			Share: share,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *grantDatabaseRoleToShareOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = emptyDatabaseObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid share identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.Share = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("grant to share", func(t *testing.T) {
		opts := setUpOpts()

		assertOptsValidAndSQLEquals(t, opts, `GRANT DATABASE ROLE %s TO SHARE %s`, id.FullyQualifiedName(), share.FullyQualifiedName())
	})
}

func TestDatabaseRoles_RevokeFromShare(t *testing.T) {
	id := randomDatabaseObjectIdentifier()
	share := randomAccountObjectIdentifier()

	setUpOpts := func() *revokeDatabaseRoleFromShareOptions {
		return &revokeDatabaseRoleFromShareOptions{
			name:  id,
			Share: share,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *revokeDatabaseRoleFromShareOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.name = emptyDatabaseObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid share identifier", func(t *testing.T) {
		opts := setUpOpts()
		opts.Share = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("revoke from share", func(t *testing.T) {
		opts := setUpOpts()

		assertOptsValidAndSQLEquals(t, opts, `REVOKE DATABASE ROLE %s FROM SHARE %s`, id.FullyQualifiedName(), share.FullyQualifiedName())
	})
}
