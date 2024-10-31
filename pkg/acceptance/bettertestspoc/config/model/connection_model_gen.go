// Code generated by config model builder generator; DO NOT EDIT.

package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

type ConnectionModel struct {
	Comment                  tfconfig.Variable `json:"comment,omitempty"`
	EnableFailoverToAccounts tfconfig.Variable `json:"enable_failover_to_accounts,omitempty"`
	FullyQualifiedName       tfconfig.Variable `json:"fully_qualified_name,omitempty"`
	Name                     tfconfig.Variable `json:"name,omitempty"`

	*config.ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func Connection(
	resourceName string,
	name string,
) *ConnectionModel {
	c := &ConnectionModel{ResourceModelMeta: config.Meta(resourceName, resources.Connection)}
	c.WithName(name)
	return c
}

func ConnectionWithDefaultMeta(
	name string,
) *ConnectionModel {
	c := &ConnectionModel{ResourceModelMeta: config.DefaultMeta(resources.Connection)}
	c.WithName(name)
	return c
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

func (c *ConnectionModel) WithComment(comment string) *ConnectionModel {
	c.Comment = tfconfig.StringVariable(comment)
	return c
}

// enable_failover_to_accounts attribute type is not yet supported, so WithEnableFailoverToAccounts can't be generated

func (c *ConnectionModel) WithFullyQualifiedName(fullyQualifiedName string) *ConnectionModel {
	c.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return c
}

func (c *ConnectionModel) WithName(name string) *ConnectionModel {
	c.Name = tfconfig.StringVariable(name)
	return c
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (c *ConnectionModel) WithCommentValue(value tfconfig.Variable) *ConnectionModel {
	c.Comment = value
	return c
}

func (c *ConnectionModel) WithEnableFailoverToAccountsValue(value tfconfig.Variable) *ConnectionModel {
	c.EnableFailoverToAccounts = value
	return c
}

func (c *ConnectionModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *ConnectionModel {
	c.FullyQualifiedName = value
	return c
}

func (c *ConnectionModel) WithNameValue(value tfconfig.Variable) *ConnectionModel {
	c.Name = value
	return c
}
