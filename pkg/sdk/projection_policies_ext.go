package sdk

import "context"

func (r *CreateProjectionPolicyRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (r projectionPolicyDescRow) additionalConvert(_ *ProjectionPolicyDetails) error {
	// DatabaseName and SchemaName are plain-only fields. DESCRIBE does not return them.
	// DescribeDetails sets them from the identifier.
	return nil
}

func (d *ProjectionPolicyDetails) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(d.DatabaseName, d.SchemaName, d.Name)
}

func (v *projectionPolicies) DescribeDetails(ctx context.Context, id SchemaObjectIdentifier) (*ProjectionPolicyDetails, error) {
	details, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	details.DatabaseName = id.DatabaseName()
	details.SchemaName = id.SchemaName()
	return details, nil
}
