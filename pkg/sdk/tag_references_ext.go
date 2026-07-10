package sdk

func (t TagReference) TagId() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(t.TagDatabase, t.TagSchema, t.TagName)
}

func (t TagReference) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(t.TagDatabase, t.TagSchema, t.TagName)
}
