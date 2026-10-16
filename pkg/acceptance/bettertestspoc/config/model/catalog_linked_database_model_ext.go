package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CatalogLinkedDatabaseModel) WithLinkedCatalog(linkedCatalog []sdk.LinkedCatalogRequest) *CatalogLinkedDatabaseModel {
	if len(linkedCatalog) == 0 {
		return c
	}
	lc := linkedCatalog[0]
	m := map[string]tfconfig.Variable{
		"catalog": tfconfig.StringVariable(lc.Catalog.Name()),
	}
	if len(lc.AllowedNamespaces) > 0 {
		m["allowed_namespaces"] = tfconfig.SetVariable(namespaceWrappersToVariables(lc.AllowedNamespaces)...)
	}
	if len(lc.BlockedNamespaces) > 0 {
		m["blocked_namespaces"] = tfconfig.SetVariable(namespaceWrappersToVariables(lc.BlockedNamespaces)...)
	}
	if lc.AllowedWriteOperations != nil {
		m["allowed_write_operations"] = tfconfig.StringVariable(string(*lc.AllowedWriteOperations))
	}
	if lc.NamespaceMode != nil {
		m["namespace_mode"] = tfconfig.StringVariable(string(*lc.NamespaceMode))
	}
	if lc.NamespaceFlattenDelimiter != nil {
		m["namespace_flatten_delimiter"] = tfconfig.StringVariable(*lc.NamespaceFlattenDelimiter)
	}
	if lc.SyncIntervalSeconds != nil {
		m["sync_interval_seconds"] = tfconfig.IntegerVariable(*lc.SyncIntervalSeconds)
	}
	c.LinkedCatalog = tfconfig.ListVariable(tfconfig.ObjectVariable(m))
	return c
}

func namespaceWrappersToVariables(namespaces []sdk.StringListItemWrapper) []tfconfig.Variable {
	return collections.Map(namespaces, func(namespace sdk.StringListItemWrapper) tfconfig.Variable {
		return tfconfig.StringVariable(namespace.Value)
	})
}
