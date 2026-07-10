package sdk

import "strings"

func (d *StorageLifecyclePolicyDetails) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(d.DatabaseName, d.SchemaName, d.Name)
}

var StorageLifecyclePolicySupportedTableTypes = []PolicyEntityDomain{
	PolicyEntityDomainTable,
	PolicyEntityDomainDynamicTable,
}

// normalizeStorageLifecyclePolicyArchiveTier normalizes the archive tier value returned by
// DESCRIBE STORAGE LIFECYCLE POLICY. When no archive tier is set, Snowflake returns the literal
// string "NULL".
// We normalize it to an empty string so that callers do not have to special-case the "NULL" literal.
func normalizeStorageLifecyclePolicyArchiveTier(archiveTier string) string {
	if strings.EqualFold(archiveTier, "NULL") {
		return ""
	}
	return archiveTier
}
