// Code generated by dto builder generator; DO NOT EDIT.

package sdk

import ()

func NewCreateMaterializedViewRequest(
	name SchemaObjectIdentifier,
	sql string,
) *CreateMaterializedViewRequest {
	s := CreateMaterializedViewRequest{}
	s.name = name
	s.sql = sql
	return &s
}

func (s *CreateMaterializedViewRequest) WithOrReplace(OrReplace *bool) *CreateMaterializedViewRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateMaterializedViewRequest) WithSecure(Secure *bool) *CreateMaterializedViewRequest {
	s.Secure = Secure
	return s
}

func (s *CreateMaterializedViewRequest) WithIfNotExists(IfNotExists *bool) *CreateMaterializedViewRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateMaterializedViewRequest) WithCopyGrants(CopyGrants *bool) *CreateMaterializedViewRequest {
	s.CopyGrants = CopyGrants
	return s
}

func (s *CreateMaterializedViewRequest) WithColumns(Columns []MaterializedViewColumnRequest) *CreateMaterializedViewRequest {
	s.Columns = Columns
	return s
}

func (s *CreateMaterializedViewRequest) WithColumnsMaskingPolicies(ColumnsMaskingPolicies []MaterializedViewColumnMaskingPolicyRequest) *CreateMaterializedViewRequest {
	s.ColumnsMaskingPolicies = ColumnsMaskingPolicies
	return s
}

func (s *CreateMaterializedViewRequest) WithComment(Comment *string) *CreateMaterializedViewRequest {
	s.Comment = Comment
	return s
}

func (s *CreateMaterializedViewRequest) WithRowAccessPolicy(RowAccessPolicy *MaterializedViewRowAccessPolicyRequest) *CreateMaterializedViewRequest {
	s.RowAccessPolicy = RowAccessPolicy
	return s
}

func (s *CreateMaterializedViewRequest) WithTag(Tag []TagAssociation) *CreateMaterializedViewRequest {
	s.Tag = Tag
	return s
}

func (s *CreateMaterializedViewRequest) WithClusterBy(ClusterBy []string) *CreateMaterializedViewRequest {
	s.ClusterBy = ClusterBy
	return s
}

func NewMaterializedViewColumnRequest(
	Name string,
) *MaterializedViewColumnRequest {
	s := MaterializedViewColumnRequest{}
	s.Name = Name
	return &s
}

func (s *MaterializedViewColumnRequest) WithComment(Comment *string) *MaterializedViewColumnRequest {
	s.Comment = Comment
	return s
}

func NewMaterializedViewColumnMaskingPolicyRequest(
	Name string,
	MaskingPolicy SchemaObjectIdentifier,
) *MaterializedViewColumnMaskingPolicyRequest {
	s := MaterializedViewColumnMaskingPolicyRequest{}
	s.Name = Name
	s.MaskingPolicy = MaskingPolicy
	return &s
}

func (s *MaterializedViewColumnMaskingPolicyRequest) WithUsing(Using []string) *MaterializedViewColumnMaskingPolicyRequest {
	s.Using = Using
	return s
}

func (s *MaterializedViewColumnMaskingPolicyRequest) WithTag(Tag []TagAssociation) *MaterializedViewColumnMaskingPolicyRequest {
	s.Tag = Tag
	return s
}

func NewMaterializedViewRowAccessPolicyRequest(
	RowAccessPolicy SchemaObjectIdentifier,
	On []string,
) *MaterializedViewRowAccessPolicyRequest {
	s := MaterializedViewRowAccessPolicyRequest{}
	s.RowAccessPolicy = RowAccessPolicy
	s.On = On
	return &s
}

func NewAlterMaterializedViewRequest(
	name SchemaObjectIdentifier,
) *AlterMaterializedViewRequest {
	s := AlterMaterializedViewRequest{}
	s.name = name
	return &s
}

func (s *AlterMaterializedViewRequest) WithRenameTo(RenameTo *SchemaObjectIdentifier) *AlterMaterializedViewRequest {
	s.RenameTo = RenameTo
	return s
}

func (s *AlterMaterializedViewRequest) WithClusterBy(ClusterBy []string) *AlterMaterializedViewRequest {
	s.ClusterBy = ClusterBy
	return s
}

func (s *AlterMaterializedViewRequest) WithDropClusteringKey(DropClusteringKey *bool) *AlterMaterializedViewRequest {
	s.DropClusteringKey = DropClusteringKey
	return s
}

func (s *AlterMaterializedViewRequest) WithSuspendRecluster(SuspendRecluster *bool) *AlterMaterializedViewRequest {
	s.SuspendRecluster = SuspendRecluster
	return s
}

func (s *AlterMaterializedViewRequest) WithResumeRecluster(ResumeRecluster *bool) *AlterMaterializedViewRequest {
	s.ResumeRecluster = ResumeRecluster
	return s
}

func (s *AlterMaterializedViewRequest) WithSuspend(Suspend *bool) *AlterMaterializedViewRequest {
	s.Suspend = Suspend
	return s
}

func (s *AlterMaterializedViewRequest) WithResume(Resume *bool) *AlterMaterializedViewRequest {
	s.Resume = Resume
	return s
}

func (s *AlterMaterializedViewRequest) WithSet(Set *MaterializedViewSetRequest) *AlterMaterializedViewRequest {
	s.Set = Set
	return s
}

func (s *AlterMaterializedViewRequest) WithUnset(Unset *MaterializedViewUnsetRequest) *AlterMaterializedViewRequest {
	s.Unset = Unset
	return s
}

func NewMaterializedViewSetRequest() *MaterializedViewSetRequest {
	return &MaterializedViewSetRequest{}
}

func (s *MaterializedViewSetRequest) WithSecure(Secure *bool) *MaterializedViewSetRequest {
	s.Secure = Secure
	return s
}

func (s *MaterializedViewSetRequest) WithComment(Comment *string) *MaterializedViewSetRequest {
	s.Comment = Comment
	return s
}

func NewMaterializedViewUnsetRequest() *MaterializedViewUnsetRequest {
	return &MaterializedViewUnsetRequest{}
}

func (s *MaterializedViewUnsetRequest) WithSecure(Secure *bool) *MaterializedViewUnsetRequest {
	s.Secure = Secure
	return s
}

func (s *MaterializedViewUnsetRequest) WithComment(Comment *bool) *MaterializedViewUnsetRequest {
	s.Comment = Comment
	return s
}

func NewDropMaterializedViewRequest(
	name SchemaObjectIdentifier,
) *DropMaterializedViewRequest {
	s := DropMaterializedViewRequest{}
	s.name = name
	return &s
}

func (s *DropMaterializedViewRequest) WithIfExists(IfExists *bool) *DropMaterializedViewRequest {
	s.IfExists = IfExists
	return s
}

func NewShowMaterializedViewRequest() *ShowMaterializedViewRequest {
	return &ShowMaterializedViewRequest{}
}

func (s *ShowMaterializedViewRequest) WithLike(Like *Like) *ShowMaterializedViewRequest {
	s.Like = Like
	return s
}

func (s *ShowMaterializedViewRequest) WithIn(In *In) *ShowMaterializedViewRequest {
	s.In = In
	return s
}

func NewDescribeMaterializedViewRequest(
	name SchemaObjectIdentifier,
) *DescribeMaterializedViewRequest {
	s := DescribeMaterializedViewRequest{}
	s.name = name
	return &s
}
