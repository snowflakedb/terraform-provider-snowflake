package model

import (
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

// WithOn implements the required constructor method for the generated code.
func (g *GrantOwnershipModel) WithOn(on []sdk.OwnershipGrantOn) *GrantOwnershipModel {
	if len(on) != 1 {
		log.Panicf("expected exactly one on block, got %d", len(on))
	}
	o := on[0]
	switch {
	case o.Object != nil:
		return g.WithOnValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type": tfconfig.StringVariable(o.Object.ObjectType.String()),
			"object_name": tfconfig.StringVariable(o.Object.Name.FullyQualifiedName()),
		}))
	case o.All != nil:
		return g.WithOnValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"all": tfconfig.ListVariable(buildBulkOperationVariable(o.All)),
		}))
	case o.Future != nil:
		return g.WithOnValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"future": tfconfig.ListVariable(buildBulkOperationVariable(o.Future)),
		}))
	default:
		log.Panicf("no valid field set in OwnershipGrantOn: %+v", o)
		return nil
	}
}

func buildBulkOperationVariable(in *sdk.GrantOnSchemaObjectIn) tfconfig.Variable {
	fields := map[string]tfconfig.Variable{
		"object_type_plural": tfconfig.StringVariable(in.PluralObjectType.String()),
	}
	if in.InSchema != nil {
		fields["in_schema"] = tfconfig.StringVariable(in.InSchema.FullyQualifiedName())
	}
	if in.InDatabase != nil {
		fields["in_database"] = tfconfig.StringVariable(in.InDatabase.Name())
	}
	return tfconfig.ObjectVariable(fields)
}
