package sdk

func (r *CreateMaterializedViewRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// TypeString returns the parsed column type when available, otherwise the raw DESC output.
func (d MaterializedViewDetails) TypeString() string {
	if d.Type != nil {
		return d.Type.ToSql()
	}
	return d.DataTypeRaw
}
