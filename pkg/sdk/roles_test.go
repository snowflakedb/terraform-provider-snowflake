package sdk

import (
	"testing"
)

func TestRolesCreate(t *testing.T) {
	t.Run("if not exists", func(t *testing.T) {
		opts := &CreateRoleOptions{
			name:        NewAccountObjectIdentifier("new_role"),
			IfNotExists: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE ROLE IF NOT EXISTS "new_role"`)
	})

	t.Run("all options", func(t *testing.T) {
		tagId := randomSchemaObjectIdentifier()
		opts := &CreateRoleOptions{
			name:      NewAccountObjectIdentifier("new_role"),
			OrReplace: Bool(true),
			Tag: []TagAssociation{
				{
					Name:  tagId,
					Value: "v1",
				},
			},
			Comment: String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE ROLE "new_role" COMMENT = 'comment' TAG (%s = 'v1')`, tagId.FullyQualifiedName())
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := &CreateRoleOptions{
			name: emptyAccountObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: one of OrReplace, IfNotExists", func(t *testing.T) {
		opts := &CreateRoleOptions{
			name:        randomAccountObjectIdentifier(),
			IfNotExists: Bool(true),
			OrReplace:   Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateRoleOptions", "OrReplace", "IfNotExists"))
	})
}

func TestRolesDrop(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		opts := &DropRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP ROLE "new_role"`)
	})

	t.Run("if exists", func(t *testing.T) {
		opts := &DropRoleOptions{
			name:     NewAccountObjectIdentifier("new_role"),
			IfExists: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP ROLE IF EXISTS "new_role"`)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := &DropRoleOptions{
			name: emptyAccountObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
}

func TestRolesAlter(t *testing.T) {
	t.Run("rename to", func(t *testing.T) {
		newID := NewAccountObjectIdentifier("new_role123")
		opts := &AlterRoleOptions{
			name:     NewAccountObjectIdentifier("new_role"),
			RenameTo: &newID,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ROLE "new_role" RENAME TO "new_role123"`)
	})

	t.Run("set comment", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name:       NewAccountObjectIdentifier("new_role"),
			SetComment: String("some comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ROLE "new_role" SET COMMENT = 'some comment'`)
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name:         NewAccountObjectIdentifier("new_role"),
			UnsetComment: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ROLE "new_role" UNSET COMMENT`)
	})

	t.Run("set tags", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
			SetTags: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag-name"),
					Value: "tag-value",
				},
				{
					Name:  NewAccountObjectIdentifier("tag-name2"),
					Value: "tag-value2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ROLE "new_role" SET TAG "tag-name" = 'tag-value', "tag-name2" = 'tag-value2'`)
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
			UnsetTags: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag-name"),
				NewAccountObjectIdentifier("tag-name2"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ROLE "new_role" UNSET TAG "tag-name", "tag-name2"`)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name:         emptyAccountObjectIdentifier,
			UnsetComment: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action specified", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name: randomAccountObjectIdentifier(),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterRoleOptions", "RenameTo", "SetComment", "SetTags", "UnsetComment", "UnsetTags"))
	})

	t.Run("validation: more than one alter action specified", func(t *testing.T) {
		opts := &AlterRoleOptions{
			name:         randomAccountObjectIdentifier(),
			SetComment:   String("comment"),
			UnsetComment: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterRoleOptions", "RenameTo", "SetComment", "SetTags", "UnsetComment", "UnsetTags"))
	})
}

func TestRolesShow(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		assertOptsValidAndSQLEquals(t, &ShowRoleOptions{}, `SHOW ROLES`)
	})

	t.Run("like", func(t *testing.T) {
		opts := &ShowRoleOptions{
			Like: &Like{
				Pattern: String("new_role"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW ROLES LIKE 'new_role'`)
	})

	t.Run("in class", func(t *testing.T) {
		class := randomSchemaObjectIdentifier()
		opts := &ShowRoleOptions{
			InClass: &RolesInClass{
				Class: class,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW ROLES IN CLASS %s`, class.FullyQualifiedName())
	})

	t.Run("validation: like with no pattern", func(t *testing.T) {
		opts := &ShowRoleOptions{
			Like: &Like{},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrPatternRequiredForLikeKeyword)
	})

	t.Run("validation: invalid class name", func(t *testing.T) {
		opts := &ShowRoleOptions{
			InClass: &RolesInClass{
				Class: emptySchemaObjectIdentifier,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
}

func TestRolesGrant(t *testing.T) {
	t.Run("user grant", func(t *testing.T) {
		opts := &GrantRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Grant: GrantRoleTo{
				User: &AccountObjectIdentifier{name: "some_user"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT ROLE "new_role" TO USER "some_user"`)
	})

	t.Run("role grant", func(t *testing.T) {
		opts := &GrantRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Grant: GrantRoleTo{
				Role: &AccountObjectIdentifier{name: "parent_role"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `GRANT ROLE "new_role" TO ROLE "parent_role"`)
	})

	t.Run("validation: invalid object identifier and no grant option", func(t *testing.T) {
		opts := &GrantRoleOptions{
			name: emptyAccountObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier, errExactlyOneOf("GrantRoleOptions.Grant", "Role", "User"))
	})

	t.Run("validation: invalid object identifier for granted role", func(t *testing.T) {
		id := emptyAccountObjectIdentifier
		opts := &GrantRoleOptions{
			name: randomAccountObjectIdentifier(),
			Grant: GrantRoleTo{
				Role: &id,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("GrantRoleOptions.Grant", "Role"))
	})

	t.Run("validation: invalid object identifier for granted user", func(t *testing.T) {
		id := emptyAccountObjectIdentifier
		opts := &GrantRoleOptions{
			name: randomAccountObjectIdentifier(),
			Grant: GrantRoleTo{
				User: &id,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("GrantRoleOptions.Grant", "User"))
	})
}

func TestRolesRevoke(t *testing.T) {
	t.Run("revoke user", func(t *testing.T) {
		opts := &RevokeRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Revoke: RevokeRoleFrom{
				User: &AccountObjectIdentifier{name: "some_user"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE ROLE "new_role" FROM USER "some_user"`)
	})

	t.Run("revoke role", func(t *testing.T) {
		opts := &RevokeRoleOptions{
			name: NewAccountObjectIdentifier("new_role"),
			Revoke: RevokeRoleFrom{
				Role: &AccountObjectIdentifier{name: "parent_role"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `REVOKE ROLE "new_role" FROM ROLE "parent_role"`)
	})

	t.Run("validation: invalid object identifier and no option set", func(t *testing.T) {
		opts := &RevokeRoleOptions{
			name: emptyAccountObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier, errExactlyOneOf("RevokeRoleOptions.Revoke", "Role", "User"))
	})
}
