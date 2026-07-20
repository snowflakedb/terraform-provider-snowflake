package genhelpers

// ShowByParentIdDef groups the three fields needed to generate a constructor
// that fetches an object by a parent identifier (e.g. userId for ProgrammaticAccessToken).
// All three fields must be set together.
type ShowByParentIdDef struct {
	ParentIdType   string
	ClientName     string
	ShowMethodName string
}

// DescribeOverrideDef overrides the default test client and method used in the IsDataSourceOutput
// constructor when the naming convention (TrimSuffix(Name,"Details") + ".Describe") does not match.
type DescribeOverrideDef struct {
	ClientName string
	MethodName string
}

type SdkObjectDetails struct {
	IdType               string
	IsDataSourceOutput   bool
	IsSubStruct          bool
	ObjectTypeName       string
	NoShowById           bool
	NoIdentifiableObject bool
	ShowByParentId       *ShowByParentIdDef
	DescribeOverride     *DescribeOverrideDef
	FromObjectIDExpr     string
	NestedAssertFields   []string
	SkipFields           []string
	StructDetails
}
