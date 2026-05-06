package sdk

func (r *CreateOpenflowRuntimeRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// OpenflowRuntime intentionally has no ID() method. SHOW OPENFLOW RUNTIMES does not return
// database_name or schema_name columns, so a SchemaObjectIdentifier cannot be reconstructed
// from the row alone. ShowByID threads the caller's identifier context to work around this.

