package sdk

import "context"

// dropSafelyHook runs before SafeDrop. DROP CODE BUNDLE has no documented pre-drop step.
func (*codeBundles) dropSafelyHook(context.Context, SchemaObjectIdentifier) error {
	return nil
}

func (r codeBundleDescRow) additionalConvert(_ *CodeBundleDetails) error {
	// DatabaseName and SchemaName are plain-only fields. DESCRIBE does not return them.
	// DescribeDetails sets them from the identifier.
	return nil
}

func (d *CodeBundleDetails) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(d.DatabaseName, d.SchemaName, d.Name)
}

func (v *codeBundles) DescribeDetails(ctx context.Context, id SchemaObjectIdentifier) (*CodeBundleDetails, error) {
	details, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	details.DatabaseName = id.DatabaseName()
	details.SchemaName = id.SchemaName()
	return details, nil
}
