package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
)

func TestMaskingPolicyCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	signature := []TableColumnSignature{
		{
			Name: "col1",
			Type: dataTypeVarchar,
		},
		{
			Name: "col2",
			Type: dataTypeVarchar,
		},
	}
	expression := "REPLACE('X', 1, 2)"

	t.Run("validation: no body", func(t *testing.T) {
		opts := &CreateMaskingPolicyOptions{
			name:      id,
			signature: signature,
			returns:   dataTypeVarchar,
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateMaskingPolicyOptions", "body"))
	})

	t.Run("validation: no signature", func(t *testing.T) {
		opts := &CreateMaskingPolicyOptions{
			name:    id,
			body:    expression,
			returns: dataTypeVarchar,
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateMaskingPolicyOptions", "signature"))
	})

	t.Run("validation: no returns", func(t *testing.T) {
		opts := &CreateMaskingPolicyOptions{
			name:      id,
			signature: signature,
			body:      expression,
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateMaskingPolicyOptions", "returns"))
	})

	t.Run("validation: both ifNotExists and orReplace present", func(t *testing.T) {
		opts := &CreateMaskingPolicyOptions{
			name:        id,
			signature:   signature,
			body:        expression,
			returns:     dataTypeVarchar,
			IfNotExists: Bool(true),
			OrReplace:   Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateMaskingPolicyOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("only required options", func(t *testing.T) {
		opts := &CreateMaskingPolicyOptions{
			name:      id,
			signature: signature,
			body:      expression,
			returns:   dataTypeVarchar,
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE MASKING POLICY %s AS ("col1" VARCHAR(16777216), "col2" VARCHAR(16777216)) RETURNS %s -> %s`, id.FullyQualifiedName(), dataTypeVarchar.ToSql(), expression)
	})

	t.Run("with complete options", func(t *testing.T) {
		comment := random.Comment()

		opts := &CreateMaskingPolicyOptions{
			OrReplace:           Bool(true),
			name:                id,
			signature:           signature,
			body:                expression,
			returns:             dataTypeVarchar,
			Comment:             String(comment),
			ExemptOtherPolicies: Bool(true),
		}

		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE MASKING POLICY %s AS ("col1" VARCHAR(16777216), "col2" VARCHAR(16777216)) RETURNS %s -> %s COMMENT = '%s' EXEMPT_OTHER_POLICIES = %t`, id.FullyQualifiedName(), dataTypeVarchar.ToSql(), expression, comment, true)
	})
}

func TestMaskingPolicyAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &AlterMaskingPolicyOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no option", func(t *testing.T) {
		opts := &AlterMaskingPolicyOptions{
			name: id,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterMaskingPolicyOptions", "Set", "Unset", "SetTag", "UnsetTag", "NewName"))
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := &AlterMaskingPolicyOptions{
			name: emptySchemaObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: new name from different db", func(t *testing.T) {
		newId := randomSchemaObjectIdentifier()

		opts := &AlterMaskingPolicyOptions{
			NewName: &newId,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrDifferentDatabase)
		assertOptsInvalidJoinedErrors(t, opts, ErrDifferentSchema)
	})

	t.Run("validation: only 1 option allowed at the same time", func(t *testing.T) {
		newID := randomSchemaObjectIdentifierInSchema(id.SchemaId())
		opts := &AlterMaskingPolicyOptions{
			name:    id,
			NewName: &newID,
			Set: &MaskingPolicySet{
				Comment: String("foo"),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterMaskingPolicyOptions", "Set", "Unset", "SetTag", "UnsetTag", "NewName"))
	})

	t.Run("with set", func(t *testing.T) {
		newComment := random.Comment()
		opts := &AlterMaskingPolicyOptions{
			name: id,
			Set: &MaskingPolicySet{
				Comment: String(newComment),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER MASKING POLICY %s SET COMMENT = '%s'", id.FullyQualifiedName(), newComment)
	})

	t.Run("with unset", func(t *testing.T) {
		opts := &AlterMaskingPolicyOptions{
			name: id,
			Unset: &MaskingPolicyUnset{
				Comment: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER MASKING POLICY %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("rename", func(t *testing.T) {
		newID := randomSchemaObjectIdentifierInSchema(id.SchemaId())
		opts := &AlterMaskingPolicyOptions{
			name:    id,
			NewName: &newID,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER MASKING POLICY %s RENAME TO %s", id.FullyQualifiedName(), newID.FullyQualifiedName())
	})

	t.Run("set body", func(t *testing.T) {
		opts := &AlterMaskingPolicyOptions{
			name: id,
			Set: &MaskingPolicySet{
				Body: Pointer("body"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER MASKING POLICY %s SET BODY -> body", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := &AlterMaskingPolicyOptions{
			name:     id,
			IfExists: Pointer(true),
			SetTag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("123"),
					Value: "value-123",
				},
				{
					Name:  NewAccountObjectIdentifier("456"),
					Value: "value-123",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER MASKING POLICY IF EXISTS %s SET TAG "123" = 'value-123', "456" = 'value-123'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := &AlterMaskingPolicyOptions{
			name:     id,
			IfExists: Pointer(true),
			UnsetTag: []ObjectIdentifier{
				NewAccountObjectIdentifier("123"),
				NewAccountObjectIdentifier("456"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER MASKING POLICY IF EXISTS %s UNSET TAG "123", "456"`, id.FullyQualifiedName())
	})
}

func TestMaskingPolicyDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &DropMaskingPolicyOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := &DropMaskingPolicyOptions{
			name: emptySchemaObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &DropMaskingPolicyOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DROP MASKING POLICY %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := &DropMaskingPolicyOptions{
			name:     id,
			IfExists: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "DROP MASKING POLICY IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestMaskingPolicyShow(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &ShowMaskingPolicyOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MASKING POLICIES")
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowMaskingPolicyOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MASKING POLICIES LIKE '%s'", id.Name())
	})

	t.Run("with like and in account", func(t *testing.T) {
		opts := &ShowMaskingPolicyOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &ExtendedIn{
				In: In{
					Account: Pointer(true),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MASKING POLICIES LIKE '%s' IN ACCOUNT", id.Name())
	})

	t.Run("with like and in database", func(t *testing.T) {
		databaseIdentifier := NewAccountObjectIdentifier(id.DatabaseName())
		opts := &ShowMaskingPolicyOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &ExtendedIn{
				In: In{
					Database: id.DatabaseId(),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MASKING POLICIES LIKE '%s' IN DATABASE %s", id.Name(), databaseIdentifier.FullyQualifiedName())
	})

	t.Run("with like and in schema", func(t *testing.T) {
		opts := &ShowMaskingPolicyOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &ExtendedIn{
				In: In{
					Schema: id.SchemaId(),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MASKING POLICIES LIKE '%s' IN SCHEMA %s", id.Name(), id.SchemaId().FullyQualifiedName())
	})

	t.Run("with limit", func(t *testing.T) {
		opts := &ShowMaskingPolicyOptions{
			Limit: &LimitFrom{
				Rows: Int(10),
				From: Pointer("foo"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MASKING POLICIES LIMIT 10 FROM 'foo'")
	})
}

func TestMaskingPolicyDescribe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &describeMaskingPolicyOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &describeMaskingPolicyOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE MASKING POLICY %s", id.FullyQualifiedName())
	})
}
