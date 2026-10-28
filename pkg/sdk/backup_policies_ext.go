package sdk

func (d *BackupPolicyDetails) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(d.DatabaseName, d.SchemaName, d.Name)
}
