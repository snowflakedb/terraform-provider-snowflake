package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

// WithObject sets object_type and object_identifier. The fields of the identifier depend on the type of the passed identifier.
func (o *ObjectParameterModel) WithObject(objectType sdk.ObjectType, id sdk.ObjectIdentifier) *ObjectParameterModel {
	identifier := map[string]tfconfig.Variable{"name": tfconfig.StringVariable(id.Name())}
	switch id := id.(type) {
	case sdk.DatabaseObjectIdentifier:
		identifier["database"] = tfconfig.StringVariable(id.DatabaseName())
	case sdk.SchemaObjectIdentifier:
		identifier["database"] = tfconfig.StringVariable(id.DatabaseName())
		identifier["schema"] = tfconfig.StringVariable(id.SchemaName())
	}
	return o.WithObjectType(objectType.String()).WithObjectIdentifierValue(tfconfig.ListVariable(tfconfig.ObjectVariable(identifier)))
}
