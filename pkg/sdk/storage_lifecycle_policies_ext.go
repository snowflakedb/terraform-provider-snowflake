package sdk

import "strings"

func (r describeStorageLifecyclePolicyDBRow) additionalConvert(_ *StorageLifecyclePolicyDetails) error {
	// additionalConvert is generated as DatabaseName and SchemaName are plain only fields.
	// They can't be set here as they are not returned by DESCRIBE; they are populated from the ID in the test helper.
	return nil
}

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
