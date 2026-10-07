package sdk

func (r *CreateSnowflakeIntelligenceRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (d *SnowflakeIntelligenceDetails) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(d.Name)
}

func (a *SnowflakeIntelligenceAgent) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(a.DatabaseName, a.SchemaName, a.Name)
}
