package sdk

func (d *CortexSearchServiceDetails) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(d.DatabaseName, d.SchemaName, d.Name)
}
