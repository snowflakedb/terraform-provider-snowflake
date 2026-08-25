package sdk

import "context"

func (r *CreateOpenflowDeploymentRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (d *OpenflowDeploymentDetails) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(d.Name)
}

// ShowParameters returns the parameters visible at the OPENFLOW DEPLOYMENT level for the given
// deployment. EVENT_TABLE is settable on CREATE and via ALTER SET/UNSET but is absent from both
// SHOW OPENFLOW DEPLOYMENTS and DESCRIBE OPENFLOW DEPLOYMENT, so this is the only way to read it back.
// Mirrors pkg/sdk/hybrid_tables_ext.go's ShowParameters with ParametersIn.OpenflowDeployment.
func (v *openflowDeployments) ShowParameters(ctx context.Context, id AccountObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			OpenflowDeployment: id,
		},
	})
}
