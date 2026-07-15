package genhelpers

// ShowByParentIdDef groups the three fields needed to generate a constructor
// that fetches an object by a parent identifier (e.g. userId for ProgrammaticAccessToken).
// All three fields must be set together.
type ShowByParentIdDef struct {
	ParentIdType   string
	ClientName     string
	ShowMethodName string
}

type SdkObjectDetails struct {
	IdType             string
	IsDataSourceOutput bool
	IsSubStruct        bool
	ObjectTypeName     string
	NoShowById         bool
	ShowByParentId     *ShowByParentIdDef
	NestedAssertFields []string
	StructDetails
}
