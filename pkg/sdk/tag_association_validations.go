package sdk

import "slices"

var (
	// based on https://docs.snowflake.com/en/user-guide/object-tagging.html#supported-objects
	TagAssociationAllowedObjectTypes = []ObjectType{
		// organization level
		ObjectTypeAccount,

		// account level
		ObjectTypeApplication,
		ObjectTypeApplicationPackage,
		ObjectTypeComputePool,
		ObjectTypeDatabase,
		ObjectTypeFailoverGroup,
		ObjectTypeIntegration,
		ObjectTypeNetworkPolicy,
		ObjectTypeReplicationGroup,
		ObjectTypeRole,
		ObjectTypeShare,
		ObjectTypeUser,
		ObjectTypeWarehouse,

		// database level
		ObjectTypeDatabaseRole,
		ObjectTypeSchema,

		// schema level
		ObjectTypeAlert,
		ObjectTypeBudget,
		ObjectTypeClassification,
		ObjectTypeExternalFunction,
		ObjectTypeExternalTable,
		ObjectTypeFunction,
		ObjectTypeImageRepository,
		ObjectTypeGitRepository,
		ObjectTypeIcebergTable,
		ObjectTypeMaterializedView,
		ObjectTypePipe,
		ObjectTypeMaskingPolicy,
		ObjectTypePasswordPolicy,
		ObjectTypeRowAccessPolicy,
		ObjectTypeSessionPolicy,
		ObjectTypePrivacyPolicy,
		ObjectTypeProcedure,
		ObjectTypeService,
		ObjectTypeStage,
		ObjectTypeStream,
		ObjectTypeTable,
		ObjectTypeTask,
		ObjectTypeView,

		// table or column level
		ObjectTypeColumn,
		ObjectTypeEventTable,
	}
	// TODO(SNOW-1229218): Object types should be able tell their id structure and tagAssociationAllowedObjectTypes should be used to filter correct object types.
	TagAssociationTagObjectTypeIsSchemaObjectType = []ObjectType{
		ObjectTypeAlert,
		ObjectTypeExternalFunction,
		ObjectTypeExternalTable,
		ObjectTypeGitRepository,
		ObjectTypeIcebergTable,
		ObjectTypeMaterializedView,
		ObjectTypeColumn, // As a workaround for object_name can be specified as `table_name.column_name`
		ObjectTypePipe,
		ObjectTypeProcedure,
		ObjectTypeStage,
		ObjectTypeStream,
		ObjectTypeTable,
		ObjectTypeTask,
		ObjectTypeView,
		ObjectTypeEventTable,
	}
	TagAssociationAllowedObjectTypesString = make([]string, len(TagAssociationAllowedObjectTypes))
)

func canBeAssociatedWithTag(o ObjectType) bool {
	return slices.Contains(TagAssociationAllowedObjectTypes, o)
}

func init() {
	for i, v := range TagAssociationAllowedObjectTypes {
		TagAssociationAllowedObjectTypesString[i] = v.String()
	}
}
