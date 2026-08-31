package sdk

import "context"

// OpenflowDeploymentFailureStatuses are terminal: a deployment in one of these will never move on its own.
var OpenflowDeploymentFailureStatuses = []OpenflowDeploymentStatus{
	OpenflowDeploymentStatusCreateFailed,
	OpenflowDeploymentStatusDeleteFailed,
	OpenflowDeploymentStatusUpgradeFailed,
	OpenflowDeploymentStatusMigrationFailed,
	OpenflowDeploymentStatusRollbackFailed,
}

// OpenflowDeploymentTransientStatuses are still coming up. Only a settled deployment accepts ALTER SET,
// TERMINATE or DROP. SPCS goes CREATING -> PROVISIONING -> ACTIVE, so PROVISIONING is transient too; BYOC
// goes straight to INACTIVE.
var OpenflowDeploymentTransientStatuses = []OpenflowDeploymentStatus{
	OpenflowDeploymentStatusCreating,
	OpenflowDeploymentStatusProvisioning,
}

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
