package sdk

func (r *CreateNetworkRuleRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (d *NetworkRuleDetails) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(d.DatabaseName, d.SchemaName, d.Name)
}
