package genhelpers

type SdkObjectDetails struct {
	IdType             string
	IsDataSourceOutput bool
	IsSubStruct        bool
	ObjectTypeName     string
	NoShowById         bool
	StructDetails
}
