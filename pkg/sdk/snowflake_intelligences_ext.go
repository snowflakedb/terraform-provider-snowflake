package sdk

func (d *SnowflakeIntelligenceDetails) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(d.Name)
}
