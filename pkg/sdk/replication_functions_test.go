package sdk

import "testing"

func TestReplicationFunctions_ShowReplicationDatabases(t *testing.T) {
	externalId := randomExternalObjectIdentifier()

	// Minimal valid ShowReplicationDatabasesOptions
	defaultOpts := func() *ShowReplicationDatabasesOptions {
		return &ShowReplicationDatabasesOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowReplicationDatabasesOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.WithPrimary]", func(t *testing.T) {
		opts := defaultOpts()
		opts.WithPrimary = Pointer(NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator(""), emptyAccountObjectIdentifier))
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW REPLICATION DATABASES")
	})

	t.Run("with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("mydb"),
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW REPLICATION DATABASES LIKE 'mydb'")
	})

	t.Run("with primary", func(t *testing.T) {
		opts := defaultOpts()
		opts.WithPrimary = &externalId
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW REPLICATION DATABASES WITH PRIMARY %s", externalId.FullyQualifiedName())
	})
}
