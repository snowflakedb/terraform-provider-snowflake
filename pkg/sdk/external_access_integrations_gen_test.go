package sdk

import (
	"testing"
)

func TestExternalAccessIntegrations_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()
	networkRuleId := randomSchemaObjectIdentifier()
	secretId := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateExternalAccessIntegrationOptions {
		return &CreateExternalAccessIntegrationOptions{
			name:                id,
			AllowedNetworkRules: []SchemaObjectIdentifier{networkRuleId},
			Enabled:             true,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateExternalAccessIntegrationOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateExternalAccessIntegrationOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("validation: AllowedNetworkRules must be non-empty", func(t *testing.T) {
		opts := defaultOpts()
		opts.AllowedNetworkRules = nil
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateExternalAccessIntegrationOptions", "AllowedNetworkRules"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s) ENABLED = true`, id.FullyQualifiedName(), networkRuleId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.AllowedAuthenticationSecrets = []SchemaObjectIdentifier{secretId}
		opts.Enabled = false
		opts.Comment = String("my comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s) ALLOWED_AUTHENTICATION_SECRETS = (%s) ENABLED = false COMMENT = 'my comment'`, id.FullyQualifiedName(), networkRuleId.FullyQualifiedName(), secretId.FullyQualifiedName())
	})

	t.Run("if not exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL ACCESS INTEGRATION IF NOT EXISTS %s ALLOWED_NETWORK_RULES = (%s) ENABLED = true`, id.FullyQualifiedName(), networkRuleId.FullyQualifiedName())
	})
}

func TestExternalAccessIntegrations_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()
	networkRuleId := randomSchemaObjectIdentifier()
	secretId := randomSchemaObjectIdentifier()

	defaultOpts := func() *AlterExternalAccessIntegrationOptions {
		return &AlterExternalAccessIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterExternalAccessIntegrationOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalAccessIntegrationOptions", "Set", "Unset"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset] should be present - both present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalAccessIntegrationSet{Enabled: Bool(true)}
		opts.Unset = &ExternalAccessIntegrationUnset{Comment: Bool(true)}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterExternalAccessIntegrationOptions", "Set", "Unset"))
	})

	t.Run("validation: at least one field in Set must be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalAccessIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterExternalAccessIntegrationOptions.Set", "AllowedNetworkRules", "AllowedAuthenticationSecrets", "Enabled", "Comment"))
	})

	t.Run("validation: at least one field in Unset must be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ExternalAccessIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterExternalAccessIntegrationOptions.Unset", "AllowedAuthenticationSecrets", "Comment"))
	})

	t.Run("set - network rules", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalAccessIntegrationSet{
			AllowedNetworkRules: []SchemaObjectIdentifier{networkRuleId},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL ACCESS INTEGRATION %s SET ALLOWED_NETWORK_RULES = (%s)`, id.FullyQualifiedName(), networkRuleId.FullyQualifiedName())
	})

	t.Run("set - all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ExternalAccessIntegrationSet{
			AllowedNetworkRules:          []SchemaObjectIdentifier{networkRuleId},
			AllowedAuthenticationSecrets: []SchemaObjectIdentifier{secretId},
			Enabled:                      Bool(true),
			Comment:                      String("updated comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL ACCESS INTEGRATION %s SET ALLOWED_NETWORK_RULES = (%s) ALLOWED_AUTHENTICATION_SECRETS = (%s) ENABLED = true COMMENT = 'updated comment'`, id.FullyQualifiedName(), networkRuleId.FullyQualifiedName(), secretId.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ExternalAccessIntegrationUnset{
			Comment: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL ACCESS INTEGRATION %s UNSET COMMENT`, id.FullyQualifiedName())
	})

	t.Run("unset allowed_authentication_secrets", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ExternalAccessIntegrationUnset{
			AllowedAuthenticationSecrets: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL ACCESS INTEGRATION %s UNSET ALLOWED_AUTHENTICATION_SECRETS`, id.FullyQualifiedName())
	})

	t.Run("unset multiple", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ExternalAccessIntegrationUnset{
			AllowedAuthenticationSecrets: Bool(true),
			Comment:                      Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL ACCESS INTEGRATION %s UNSET ALLOWED_AUTHENTICATION_SECRETS, COMMENT`, id.FullyQualifiedName())
	})
}

func TestExternalAccessIntegrations_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *DropExternalAccessIntegrationOptions {
		return &DropExternalAccessIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DropExternalAccessIntegrationOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP EXTERNAL ACCESS INTEGRATION %s`, id.FullyQualifiedName())
	})

	t.Run("if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP EXTERNAL ACCESS INTEGRATION IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestExternalAccessIntegrations_Show(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *ShowExternalAccessIntegrationOptions {
		return &ShowExternalAccessIntegrationOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowExternalAccessIntegrationOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW EXTERNAL ACCESS INTEGRATIONS`)
	})

	t.Run("with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW EXTERNAL ACCESS INTEGRATIONS LIKE '%s'`, id.Name())
	})
}

func TestExternalAccessIntegrations_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *DescribeExternalAccessIntegrationOptions {
		return &DescribeExternalAccessIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeExternalAccessIntegrationOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE EXTERNAL ACCESS INTEGRATION %s`, id.FullyQualifiedName())
	})
}
