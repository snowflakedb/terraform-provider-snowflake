package sdk

import "testing"

func TestCatalogIntegrations_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOptsObjectStore := func() *CreateCatalogIntegrationOptions {
		return &CreateCatalogIntegrationOptions{
			name:    id,
			Enabled: true,
			ObjectStoreParams: &ObjectStoreParams{
				TableFormat: TableFormatIceberg,
			},
		}
	}

	defaultOptsGlue := func() *CreateCatalogIntegrationOptions {
		return &CreateCatalogIntegrationOptions{
			name:    id,
			Enabled: true,
			GlueParams: &GlueParams{
				TableFormat:    TableFormatIceberg,
				GlueAwsRoleArn: "arn:aws:iam::123456789012:role/myrole",
				GlueCatalogId:  "123456789012",
			},
		}
	}

	defaultOpts := defaultOptsObjectStore

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateCatalogIntegrationOptions)(nil)
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
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateCatalogIntegrationOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("validation: exactly one field from source params should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ObjectStoreParams = nil
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateCatalogIntegrationOptions", "ObjectStoreParams", "GlueParams", "IcebergRestParams", "PolarisParams", "SapBdcParams"))
	})

	t.Run("validation: exactly one field from source params should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.GlueParams = &GlueParams{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateCatalogIntegrationOptions", "ObjectStoreParams", "GlueParams", "IcebergRestParams", "PolarisParams", "SapBdcParams"))
	})

	t.Run("basic - object store iceberg", func(t *testing.T) {
		opts := defaultOptsObjectStore()
		assertOptsValidAndSQLEquals(t, opts, "CREATE CATALOG INTEGRATION %s ENABLED = true CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = ICEBERG", id.FullyQualifiedName())
	})

	t.Run("all options - object store delta", func(t *testing.T) {
		opts := defaultOptsObjectStore()
		opts.IfNotExists = Bool(true)
		opts.ObjectStoreParams.TableFormat = TableFormatDelta
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE CATALOG INTEGRATION IF NOT EXISTS %s ENABLED = true CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = DELTA COMMENT = 'some comment'", id.FullyQualifiedName())
	})

	t.Run("all options - glue", func(t *testing.T) {
		opts := defaultOptsGlue()
		opts.IfNotExists = Bool(true)
		opts.GlueParams.GlueRegion = String("us-east-1")
		opts.CatalogNamespace = String("my_namespace")
		opts.Comment = String("glue comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE CATALOG INTEGRATION IF NOT EXISTS %s ENABLED = true CATALOG_SOURCE = GLUE TABLE_FORMAT = ICEBERG GLUE_AWS_ROLE_ARN = 'arn:aws:iam::123456789012:role/myrole' GLUE_CATALOG_ID = '123456789012' GLUE_REGION = 'us-east-1' CATALOG_NAMESPACE = 'my_namespace' COMMENT = 'glue comment'", id.FullyQualifiedName())
	})

	t.Run("or replace - object store", func(t *testing.T) {
		opts := defaultOptsObjectStore()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE CATALOG INTEGRATION %s ENABLED = true CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = ICEBERG", id.FullyQualifiedName())
	})
}

func TestCatalogIntegrations_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *AlterCatalogIntegrationOptions {
		return &AlterCatalogIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterCatalogIntegrationOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterCatalogIntegrationOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the fields in Set should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &CatalogIntegrationSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterCatalogIntegrationOptions.Set", "Enabled", "GlueAwsRoleArn", "GlueCatalogId", "GlueRegion", "RestConfig", "RestAuthentication", "CatalogNamespace", "Comment"))
	})

	t.Run("validation: at least one of the fields in Unset should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &CatalogIntegrationUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterCatalogIntegrationOptions.Unset", "Enabled", "CatalogNamespace", "Comment"))
	})

	t.Run("set enabled and comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &CatalogIntegrationSet{
			Enabled: Bool(true),
			Comment: String("new comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CATALOG INTEGRATION %s SET ENABLED = true COMMENT = 'new comment'", id.FullyQualifiedName())
	})

	t.Run("set with if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Set = &CatalogIntegrationSet{
			Enabled: Bool(false),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CATALOG INTEGRATION IF EXISTS %s SET ENABLED = false", id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &CatalogIntegrationUnset{
			Comment: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CATALOG INTEGRATION %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("unset multiple", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &CatalogIntegrationUnset{
			Enabled:          Bool(true),
			CatalogNamespace: Bool(true),
			Comment:          Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER CATALOG INTEGRATION %s UNSET ENABLED, CATALOG_NAMESPACE, COMMENT", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("name"),
				Value: "value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER CATALOG INTEGRATION %s SET TAG "name" = 'value'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER CATALOG INTEGRATION %s UNSET TAG "name"`, id.FullyQualifiedName())
	})
}

func TestCatalogIntegrations_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *DropCatalogIntegrationOptions {
		return &DropCatalogIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DropCatalogIntegrationOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP CATALOG INTEGRATION %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP CATALOG INTEGRATION IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestCatalogIntegrations_Show(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *ShowCatalogIntegrationOptions {
		return &ShowCatalogIntegrationOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowCatalogIntegrationOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW CATALOG INTEGRATIONS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW CATALOG INTEGRATIONS LIKE '%s'", id.Name())
	})
}

func TestCatalogIntegrations_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *DescribeCatalogIntegrationOptions {
		return &DescribeCatalogIntegrationOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeCatalogIntegrationOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE CATALOG INTEGRATION %s", id.FullyQualifiedName())
	})
}
