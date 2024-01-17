// Code generated by dto builder generator; DO NOT EDIT.

package sdk

import ()

func NewCreateRowAccessPolicyRequest(
	name SchemaObjectIdentifier,
	args []CreateRowAccessPolicyArgsRequest,
	body string,
) *CreateRowAccessPolicyRequest {
	s := CreateRowAccessPolicyRequest{}
	s.name = name
	s.args = args
	s.body = body
	return &s
}

func (r *CreateRowAccessPolicyRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (s *CreateRowAccessPolicyRequest) WithOrReplace(OrReplace *bool) *CreateRowAccessPolicyRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateRowAccessPolicyRequest) WithIfNotExists(IfNotExists *bool) *CreateRowAccessPolicyRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateRowAccessPolicyRequest) WithComment(Comment *string) *CreateRowAccessPolicyRequest {
	s.Comment = Comment
	return s
}

func NewCreateRowAccessPolicyArgsRequest(
	Name string,
	Type string,
) *CreateRowAccessPolicyArgsRequest {
	s := CreateRowAccessPolicyArgsRequest{}
	s.Name = Name
	s.Type = Type
	return &s
}

func NewAlterRowAccessPolicyRequest(
	name SchemaObjectIdentifier,
) *AlterRowAccessPolicyRequest {
	s := AlterRowAccessPolicyRequest{}
	s.name = name
	return &s
}

func (s *AlterRowAccessPolicyRequest) WithRenameTo(RenameTo *SchemaObjectIdentifier) *AlterRowAccessPolicyRequest {
	s.RenameTo = RenameTo
	return s
}

func (s *AlterRowAccessPolicyRequest) WithSetBody(SetBody *string) *AlterRowAccessPolicyRequest {
	s.SetBody = SetBody
	return s
}

func (s *AlterRowAccessPolicyRequest) WithSetTags(SetTags []TagAssociation) *AlterRowAccessPolicyRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterRowAccessPolicyRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterRowAccessPolicyRequest {
	s.UnsetTags = UnsetTags
	return s
}

func (s *AlterRowAccessPolicyRequest) WithSetComment(SetComment *string) *AlterRowAccessPolicyRequest {
	s.SetComment = SetComment
	return s
}

func (s *AlterRowAccessPolicyRequest) WithUnsetComment(UnsetComment *bool) *AlterRowAccessPolicyRequest {
	s.UnsetComment = UnsetComment
	return s
}

func NewDropRowAccessPolicyRequest(
	name SchemaObjectIdentifier,
) *DropRowAccessPolicyRequest {
	s := DropRowAccessPolicyRequest{}
	s.name = name
	return &s
}

func NewShowRowAccessPolicyRequest() *ShowRowAccessPolicyRequest {
	return &ShowRowAccessPolicyRequest{}
}

func (s *ShowRowAccessPolicyRequest) WithLike(Like *Like) *ShowRowAccessPolicyRequest {
	s.Like = Like
	return s
}

func (s *ShowRowAccessPolicyRequest) WithIn(In *In) *ShowRowAccessPolicyRequest {
	s.In = In
	return s
}

func NewDescribeRowAccessPolicyRequest(
	name SchemaObjectIdentifier,
) *DescribeRowAccessPolicyRequest {
	s := DescribeRowAccessPolicyRequest{}
	s.name = name
	return &s
}
