package sdk

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
		ObjectTypeAgent,
		ObjectTypeAlert,
		ObjectTypeBudget,
		ObjectTypeClassification,
		ObjectTypeExternalFunction,
		ObjectTypeExternalTable,
		ObjectTypeFunction,
		ObjectTypeImageRepository,
		ObjectTypeGitRepository,
		ObjectTypeIcebergTable,
		ObjectTypeInteractiveTable,
		ObjectTypeMaterializedView,
		ObjectTypePipe,
		ObjectTypeMaskingPolicy,
		ObjectTypePasswordPolicy,
		ObjectTypeRowAccessPolicy,
		ObjectTypeSessionPolicy,
		ObjectTypeStorageLifecyclePolicy,
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
		ObjectTypeIcebergTableColumn,
	}
	TagAssociationAllowedObjectTypesString = make([]string, len(TagAssociationAllowedObjectTypes))
)

func init() {
	for i, v := range TagAssociationAllowedObjectTypes {
		TagAssociationAllowedObjectTypesString[i] = v.String()
	}
}
