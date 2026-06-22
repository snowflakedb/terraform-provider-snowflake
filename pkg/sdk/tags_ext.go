package sdk

// TagOnConflictAllowedValuesSequence is the value returned by SHOW TAGS in the on_conflict column
// when the tag propagation on-conflict strategy is set to ALLOWED_VALUES_SEQUENCE.
const TagOnConflictAllowedValuesSequence = "ALLOWED_VALUES_SEQUENCE"

func (r *CreateTagRequest) GetName() SchemaObjectIdentifier {
	return r.name
}
