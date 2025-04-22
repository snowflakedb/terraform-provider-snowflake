package sdk

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateRowAccessPolicyOptions]   = new(CreateRowAccessPolicyRequest)
	_ optionsProvider[AlterRowAccessPolicyOptions]    = new(AlterRowAccessPolicyRequest)
	_ optionsProvider[DropRowAccessPolicyOptions]     = new(DropRowAccessPolicyRequest)
	_ optionsProvider[ShowRowAccessPolicyOptions]     = new(ShowRowAccessPolicyRequest)
	_ optionsProvider[DescribeRowAccessPolicyOptions] = new(DescribeRowAccessPolicyRequest)
)

type CreateRowAccessPolicyRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        SchemaObjectIdentifier             // required
	args        []CreateRowAccessPolicyArgsRequest // required
	body        string                             // required
	Comment     *string
}

type CreateRowAccessPolicyArgsRequest struct {
	Name string // required
	// modified manually
	Type datatypes.DataType // required
}

type AlterRowAccessPolicyRequest struct {
	name         SchemaObjectIdentifier // required
	RenameTo     *SchemaObjectIdentifier
	SetBody      *string
	SetTags      []TagAssociation
	UnsetTags    []ObjectIdentifier
	SetComment   *string
	UnsetComment *bool
}

type DropRowAccessPolicyRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowRowAccessPolicyRequest struct {
	Like  *Like
	In    *ExtendedIn
	Limit *LimitFrom
}

type DescribeRowAccessPolicyRequest struct {
	name SchemaObjectIdentifier // required
}
