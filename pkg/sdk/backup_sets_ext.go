package sdk

func (d *BackupSetDetails) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(d.DatabaseName, d.SchemaName, d.Name)
}
