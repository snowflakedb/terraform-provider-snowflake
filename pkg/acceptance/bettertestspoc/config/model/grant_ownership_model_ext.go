package model

import (
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

// GrantOwnershipWithRawOn creates a GrantOwnershipModel without calling WithOn, allowing callers to set On via WithOnValue.
func GrantOwnershipWithRawOn(resourceName string) *GrantOwnershipModel {
	return &GrantOwnershipModel{ResourceModelMeta: config.Meta(resourceName, resources.GrantOwnership)}
}

// WithOnObject sets the on block to a specific object type and name.
func (g *GrantOwnershipModel) WithOnObject(objectType, objectName string) *GrantOwnershipModel {
	return g.WithOnValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"object_type": tfconfig.StringVariable(objectType),
		"object_name": tfconfig.StringVariable(objectName),
	}))
}

// WithOnAllInDatabase sets the on block to all objects of the given plural type in a database.
func (g *GrantOwnershipModel) WithOnAllInDatabase(pluralObjectType, databaseName string) *GrantOwnershipModel {
	return g.WithOnValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"all": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type_plural": tfconfig.StringVariable(pluralObjectType),
			"in_database":        tfconfig.StringVariable(databaseName),
		})),
	}))
}

// WithOnAllInSchema sets the on block to all objects of the given plural type in a schema.
func (g *GrantOwnershipModel) WithOnAllInSchema(pluralObjectType, schemaFQN string) *GrantOwnershipModel {
	return g.WithOnValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"all": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type_plural": tfconfig.StringVariable(pluralObjectType),
			"in_schema":          tfconfig.StringVariable(schemaFQN),
		})),
	}))
}

// WithOnFutureInDatabase sets the on block to future objects of the given plural type in a database.
func (g *GrantOwnershipModel) WithOnFutureInDatabase(pluralObjectType, databaseName string) *GrantOwnershipModel {
	return g.WithOnValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"future": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type_plural": tfconfig.StringVariable(pluralObjectType),
			"in_database":        tfconfig.StringVariable(databaseName),
		})),
	}))
}

// WithOnFutureInSchema sets the on block to future objects of the given plural type in a schema.
func (g *GrantOwnershipModel) WithOnFutureInSchema(pluralObjectType, schemaFQN string) *GrantOwnershipModel {
	return g.WithOnValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"future": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"object_type_plural": tfconfig.StringVariable(pluralObjectType),
			"in_schema":          tfconfig.StringVariable(schemaFQN),
		})),
	}))
}

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
