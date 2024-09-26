// Code generated by config model builder generator; DO NOT EDIT.

package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

type SecretWithClientCredentialsModel struct {
	ApiAuthentication  tfconfig.Variable `json:"api_authentication,omitempty"`
	Comment            tfconfig.Variable `json:"comment,omitempty"`
	Database           tfconfig.Variable `json:"database,omitempty"`
	FullyQualifiedName tfconfig.Variable `json:"fully_qualified_name,omitempty"`
	Name               tfconfig.Variable `json:"name,omitempty"`
	OauthScopes        tfconfig.Variable `json:"oauth_scopes,omitempty"`
	Schema             tfconfig.Variable `json:"schema,omitempty"`

	*config.ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func SecretWithClientCredentials(
	resourceName string,
	apiAuthentication string,
	database string,
	schema string,
	name string,
	oauthScopes []string,
) *SecretWithClientCredentialsModel {
	s := &SecretWithClientCredentialsModel{ResourceModelMeta: config.Meta(resourceName, resources.SecretWithClientCredentials)}
	s.WithApiAuthentication(apiAuthentication)
	s.WithDatabase(database)
	s.WithName(name)
	s.WithOauthScopes(oauthScopes)
	s.WithSchema(schema)
	return s
}

func SecretWithClientCredentialsWithDefaultMeta(
	apiAuthentication string,
	database string,
	name string,
	oauthScopes []string,
	schema string,
) *SecretWithClientCredentialsModel {
	s := &SecretWithClientCredentialsModel{ResourceModelMeta: config.DefaultMeta(resources.SecretWithClientCredentials)}
	s.WithApiAuthentication(apiAuthentication)
	s.WithDatabase(database)
	s.WithName(name)
	s.WithOauthScopes(oauthScopes)
	s.WithSchema(schema)
	return s
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

func (s *SecretWithClientCredentialsModel) WithApiAuthentication(apiAuthentication string) *SecretWithClientCredentialsModel {
	s.ApiAuthentication = tfconfig.StringVariable(apiAuthentication)
	return s
}

func (s *SecretWithClientCredentialsModel) WithComment(comment string) *SecretWithClientCredentialsModel {
	s.Comment = tfconfig.StringVariable(comment)
	return s
}

func (s *SecretWithClientCredentialsModel) WithDatabase(database string) *SecretWithClientCredentialsModel {
	s.Database = tfconfig.StringVariable(database)
	return s
}

func (s *SecretWithClientCredentialsModel) WithFullyQualifiedName(fullyQualifiedName string) *SecretWithClientCredentialsModel {
	s.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return s
}

func (s *SecretWithClientCredentialsModel) WithName(name string) *SecretWithClientCredentialsModel {
	s.Name = tfconfig.StringVariable(name)
	return s
}

// oauth_scopes attribute type is not yet supported, so WithOauthScopes can't be generated

func (s *SecretWithClientCredentialsModel) WithSchema(schema string) *SecretWithClientCredentialsModel {
	s.Schema = tfconfig.StringVariable(schema)
	return s
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (s *SecretWithClientCredentialsModel) WithApiAuthenticationValue(value tfconfig.Variable) *SecretWithClientCredentialsModel {
	s.ApiAuthentication = value
	return s
}

func (s *SecretWithClientCredentialsModel) WithCommentValue(value tfconfig.Variable) *SecretWithClientCredentialsModel {
	s.Comment = value
	return s
}

func (s *SecretWithClientCredentialsModel) WithDatabaseValue(value tfconfig.Variable) *SecretWithClientCredentialsModel {
	s.Database = value
	return s
}

func (s *SecretWithClientCredentialsModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *SecretWithClientCredentialsModel {
	s.FullyQualifiedName = value
	return s
}

func (s *SecretWithClientCredentialsModel) WithNameValue(value tfconfig.Variable) *SecretWithClientCredentialsModel {
	s.Name = value
	return s
}

func (s *SecretWithClientCredentialsModel) WithOauthScopesValue(value tfconfig.Variable) *SecretWithClientCredentialsModel {
	s.OauthScopes = value
	return s
}

func (s *SecretWithClientCredentialsModel) WithSchemaValue(value tfconfig.Variable) *SecretWithClientCredentialsModel {
	s.Schema = value
	return s
}
func (s*SecretWithClientCredentialsModel) WithOauthScopes(oauthScopes []string) *SecretWithClientCredentialsModel {
	oauthScopesStringVariables := make([]tfconfig.Variable, len(oauthScopes))
	for i, v := range oauthScopes {
		oauthScopesStringVariables[i] = tfconfig.StringVariable(v)
	}

	s.OauthScopes = tfconfig.SetVariable(oauthScopesStringVariables...)
	return s
}
