package sdk

func (v *Stream) IsAppendOnly() bool {
	return v != nil && v.Mode != nil && *v.Mode == StreamModeAppendOnly
}

func (v *Stream) IsInsertOnly() bool {
	return v != nil && v.Mode != nil && *v.Mode == StreamModeInsertOnly
}

func (r *CreateOnTableStreamRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (r *CreateOnViewStreamRequest) GetName() SchemaObjectIdentifier {
	return r.name
}
