package sdk

func NewCreateDynamicTableRequest(
	name SchemaObjectIdentifier,
	warehouse AccountObjectIdentifier,
	targetLag TargetLag,
	query string,
) *CreateDynamicTableRequest {
	s := CreateDynamicTableRequest{}
	s.name = name
	s.warehouse = warehouse
	s.targetLag = targetLag
	s.query = query
	return &s
}

func (s *CreateDynamicTableRequest) WithOrReplace(orReplace bool) *CreateDynamicTableRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateDynamicTableRequest) WithComment(comment *string) *CreateDynamicTableRequest {
	s.comment = comment
	return s
}

func (s *CreateDynamicTableRequest) WithRefreshMode(refreshMode DynamicTableRefreshMode) *CreateDynamicTableRequest {
	s.refreshMode = &refreshMode
	return s
}

func (s *CreateDynamicTableRequest) WithInitialize(initialize DynamicTableInitialize) *CreateDynamicTableRequest {
	s.initialize = &initialize
	return s
}

func NewAlterDynamicTableRequest(
	name SchemaObjectIdentifier,
) *AlterDynamicTableRequest {
	s := AlterDynamicTableRequest{}
	s.name = name
	return &s
}

func (s *AlterDynamicTableRequest) WithSuspend(suspend *bool) *AlterDynamicTableRequest {
	s.suspend = suspend
	return s
}

func (s *AlterDynamicTableRequest) WithResume(resume *bool) *AlterDynamicTableRequest {
	s.resume = resume
	return s
}

func (s *AlterDynamicTableRequest) WithRefresh(refresh *bool) *AlterDynamicTableRequest {
	s.refresh = refresh
	return s
}

func (s *AlterDynamicTableRequest) WithSet(set *DynamicTableSetRequest) *AlterDynamicTableRequest {
	s.set = set
	return s
}

func (s *AlterDynamicTableRequest) WithAddStorageLifecyclePolicy(addStorageLifecyclePolicy *DynamicTableAddStorageLifecyclePolicyRequest) *AlterDynamicTableRequest {
	s.addStorageLifecyclePolicy = addStorageLifecyclePolicy
	return s
}

func (s *AlterDynamicTableRequest) WithDropStorageLifecyclePolicy(dropStorageLifecyclePolicy *bool) *AlterDynamicTableRequest {
	s.dropStorageLifecyclePolicy = dropStorageLifecyclePolicy
	return s
}

func NewDynamicTableAddStorageLifecyclePolicyRequest(
	storageLifecyclePolicy SchemaObjectIdentifier,
	on []Column,
) *DynamicTableAddStorageLifecyclePolicyRequest {
	s := DynamicTableAddStorageLifecyclePolicyRequest{}
	s.StorageLifecyclePolicy = storageLifecyclePolicy
	s.On = on
	return &s
}

func NewDynamicTableSetRequest() *DynamicTableSetRequest {
	return &DynamicTableSetRequest{}
}

func (s *DynamicTableSetRequest) WithTargetLag(targetLag TargetLag) *DynamicTableSetRequest {
	s.targetLag = &targetLag
	return s
}

func (s *DynamicTableSetRequest) WithWarehouse(warehouse AccountObjectIdentifier) *DynamicTableSetRequest {
	s.warehouse = &warehouse
	return s
}

func NewDropDynamicTableRequest(
	name SchemaObjectIdentifier,
) *DropDynamicTableRequest {
	s := DropDynamicTableRequest{}
	s.name = name
	return &s
}

func (s *DropDynamicTableRequest) WithIfExists(ifExists bool) *DropDynamicTableRequest {
	s.IfExists = &ifExists
	return s
}

func NewDescribeDynamicTableRequest(
	name SchemaObjectIdentifier,
) *DescribeDynamicTableRequest {
	s := DescribeDynamicTableRequest{}
	s.name = name
	return &s
}

func NewShowDynamicTableRequest() *ShowDynamicTableRequest {
	return &ShowDynamicTableRequest{}
}

func (s *ShowDynamicTableRequest) WithLike(like *Like) *ShowDynamicTableRequest {
	s.like = like
	return s
}

func (s *ShowDynamicTableRequest) WithIn(in *In) *ShowDynamicTableRequest {
	s.in = in
	return s
}

func (s *ShowDynamicTableRequest) WithStartsWith(startsWith *string) *ShowDynamicTableRequest {
	s.startsWith = startsWith
	return s
}

func (s *ShowDynamicTableRequest) WithLimit(limit *LimitFrom) *ShowDynamicTableRequest {
	s.limit = limit
	return s
}
