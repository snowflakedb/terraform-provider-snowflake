package sdk

import (
	"testing"
)

func TestComments(t *testing.T) {
	t.Run("set on schema", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := &SetCommentOptions{
			ObjectType: ObjectTypeSchema,
			ObjectName: &id,
			Value:      new("mycomment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `COMMENT ON SCHEMA %s IS 'mycomment'`, id.FullyQualifiedName())
	})

	t.Run("set if exists", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &SetCommentOptions{
			IfExists:   new(true),
			ObjectType: ObjectTypeMaskingPolicy,
			ObjectName: &id,
			Value:      new("mycomment2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `COMMENT IF EXISTS ON MASKING POLICY %s IS 'mycomment2'`, id.FullyQualifiedName())
	})

	t.Run("set column comment", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := &SetColumnCommentOptions{
			Column: id,
			Value:  new("mycomment3"),
		}
		assertOptsValidAndSQLEquals(t, opts, `COMMENT ON COLUMN %s IS 'mycomment3'`, id.FullyQualifiedName())
	})
}
