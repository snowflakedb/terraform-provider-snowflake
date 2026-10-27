package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CatalogLinkedDatabaseResourceAssert) HasLinkedCatalogCatalog(expected sdk.AccountObjectIdentifier) *CatalogLinkedDatabaseResourceAssert {
	c.ValueSet("linked_catalog.#", "1")
	c.ValueSet("linked_catalog.0.catalog", expected.Name())
	return c
}

func (c *CatalogLinkedDatabaseResourceAssert) HasLinkedCatalogAllowedNamespaces(expected ...string) *CatalogLinkedDatabaseResourceAssert {
	c.SetContainsExactlyStringValues("linked_catalog.0.allowed_namespaces", expected...)
	return c
}

func (c *CatalogLinkedDatabaseResourceAssert) HasLinkedCatalogBlockedNamespaces(expected ...string) *CatalogLinkedDatabaseResourceAssert {
	c.SetContainsExactlyStringValues("linked_catalog.0.blocked_namespaces", expected...)
	return c
}

func (c *CatalogLinkedDatabaseResourceAssert) HasLinkedCatalogAllowedWriteOperations(expected sdk.CatalogLinkedDatabaseAllowedWriteOperations) *CatalogLinkedDatabaseResourceAssert {
	c.ValueSet("linked_catalog.0.allowed_write_operations", string(expected))
	return c
}

func (c *CatalogLinkedDatabaseResourceAssert) HasLinkedCatalogNamespaceMode(expected sdk.CatalogLinkedDatabaseNamespaceMode) *CatalogLinkedDatabaseResourceAssert {
	c.ValueSet("linked_catalog.0.namespace_mode", string(expected))
	return c
}

func (c *CatalogLinkedDatabaseResourceAssert) HasLinkedCatalogNamespaceFlattenDelimiter(expected string) *CatalogLinkedDatabaseResourceAssert {
	c.ValueSet("linked_catalog.0.namespace_flatten_delimiter", expected)
	return c
}

func (c *CatalogLinkedDatabaseResourceAssert) HasLinkedCatalogSyncIntervalSeconds(expected int) *CatalogLinkedDatabaseResourceAssert {
	c.IntValueSet("linked_catalog.0.sync_interval_seconds", expected)
	return c
}
