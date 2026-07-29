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
			Value:      String("mycomment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `COMMENT ON SCHEMA %s IS 'mycomment'`, id.FullyQualifiedName())
	})

	t.Run("set if exists", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &SetCommentOptions{
			IfExists:   Bool(true),
			ObjectType: ObjectTypeMaskingPolicy,
			ObjectName: &id,
			Value:      String("mycomment2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `COMMENT IF EXISTS ON MASKING POLICY %s IS 'mycomment2'`, id.FullyQualifiedName())
	})

	t.Run("set column comment", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := &SetColumnCommentOptions{
			Column: id,
			Value:  String("mycomment3"),
		}
		assertOptsValidAndSQLEquals(t, opts, `COMMENT ON COLUMN %s IS 'mycomment3'`, id.FullyQualifiedName())
	})
}
