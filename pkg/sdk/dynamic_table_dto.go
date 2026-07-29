package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[createDynamicTableOptions] = new(CreateDynamicTableRequest)
	_ optionsProvider[alterDynamicTableOptions]  = new(AlterDynamicTableRequest)
	_ optionsProvider[dropDynamicTableOptions]   = new(DropDynamicTableRequest)
	_ optionsProvider[showDynamicTableOptions]   = new(ShowDynamicTableRequest)
)

type CreateDynamicTableRequest struct {
	orReplace bool

	name      SchemaObjectIdentifier  // required
	warehouse AccountObjectIdentifier // required
	targetLag TargetLag               // required
	query     string                  // required

	comment     *string
	refreshMode *DynamicTableRefreshMode
	initialize  *DynamicTableInitialize
}

type AlterDynamicTableRequest struct {
	name SchemaObjectIdentifier // required

	// One of
	suspend                    *bool
	resume                     *bool
	refresh                    *bool
	set                        *DynamicTableSetRequest
	addStorageLifecyclePolicy  *DynamicTableAddStorageLifecyclePolicyRequest
	dropStorageLifecyclePolicy *bool
}

type DynamicTableAddStorageLifecyclePolicyRequest struct {
	StorageLifecyclePolicy SchemaObjectIdentifier // required
	On                     []Column               // required
}

type DynamicTableSetRequest struct {
	targetLag *TargetLag
	warehouse *AccountObjectIdentifier
}

type DropDynamicTableRequest struct {
	name     SchemaObjectIdentifier // required
	IfExists *bool
}

type DescribeDynamicTableRequest struct {
	name SchemaObjectIdentifier // required
}

type ShowDynamicTableRequest struct {
	like       *Like
	in         *In
	startsWith *string
	limit      *LimitFrom
}
