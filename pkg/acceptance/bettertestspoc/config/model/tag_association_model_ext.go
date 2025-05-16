package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (t *TagAssociationModel) WithObjectIdentifiers(objectIdentifiers []sdk.ObjectIdentifier) *TagAssociationModel {
	objectIdentifiersStringVariables := make([]tfconfig.Variable, len(objectIdentifiers))
	for i, v := range objectIdentifiers {
		objectIdentifiersStringVariables[i] = tfconfig.StringVariable(v.FullyQualifiedName())
	}

	t.ObjectIdentifiers = tfconfig.SetVariable(objectIdentifiersStringVariables...)
	return t
}
