// Code generated by config model builder generator; DO NOT EDIT.

package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

type StreamOnTableModel struct {
	AppendOnly         tfconfig.Variable `json:"append_only,omitempty"`
	At                 tfconfig.Variable `json:"at,omitempty"`
	Before             tfconfig.Variable `json:"before,omitempty"`
	Comment            tfconfig.Variable `json:"comment,omitempty"`
	CopyGrants         tfconfig.Variable `json:"copy_grants,omitempty"`
	Database           tfconfig.Variable `json:"database,omitempty"`
	FullyQualifiedName tfconfig.Variable `json:"fully_qualified_name,omitempty"`
	Name               tfconfig.Variable `json:"name,omitempty"`
	Schema             tfconfig.Variable `json:"schema,omitempty"`
	ShowInitialRows    tfconfig.Variable `json:"show_initial_rows,omitempty"`
	Table              tfconfig.Variable `json:"table,omitempty"`

	*config.ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func StreamOnTable(
	resourceName string,
	database string,
	name string,
	schema string,
	table string,
) *StreamOnTableModel {
	s := &StreamOnTableModel{ResourceModelMeta: config.Meta(resourceName, resources.StreamOnTable)}
	s.WithDatabase(database)
	s.WithName(name)
	s.WithSchema(schema)
	s.WithTable(table)
	return s
}

func StreamOnTableWithDefaultMeta(
	database string,
	name string,
	schema string,
	table string,
) *StreamOnTableModel {
	s := &StreamOnTableModel{ResourceModelMeta: config.DefaultMeta(resources.StreamOnTable)}
	s.WithDatabase(database)
	s.WithName(name)
	s.WithSchema(schema)
	s.WithTable(table)
	return s
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

func (s *StreamOnTableModel) WithAppendOnly(appendOnly string) *StreamOnTableModel {
	s.AppendOnly = tfconfig.StringVariable(appendOnly)
	return s
}

// at attribute type is not yet supported, so WithAt can't be generated

// before attribute type is not yet supported, so WithBefore can't be generated

func (s *StreamOnTableModel) WithComment(comment string) *StreamOnTableModel {
	s.Comment = tfconfig.StringVariable(comment)
	return s
}

func (s *StreamOnTableModel) WithCopyGrants(copyGrants bool) *StreamOnTableModel {
	s.CopyGrants = tfconfig.BoolVariable(copyGrants)
	return s
}

func (s *StreamOnTableModel) WithDatabase(database string) *StreamOnTableModel {
	s.Database = tfconfig.StringVariable(database)
	return s
}

func (s *StreamOnTableModel) WithFullyQualifiedName(fullyQualifiedName string) *StreamOnTableModel {
	s.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return s
}

func (s *StreamOnTableModel) WithName(name string) *StreamOnTableModel {
	s.Name = tfconfig.StringVariable(name)
	return s
}

func (s *StreamOnTableModel) WithSchema(schema string) *StreamOnTableModel {
	s.Schema = tfconfig.StringVariable(schema)
	return s
}

func (s *StreamOnTableModel) WithShowInitialRows(showInitialRows string) *StreamOnTableModel {
	s.ShowInitialRows = tfconfig.StringVariable(showInitialRows)
	return s
}

func (s *StreamOnTableModel) WithTable(table string) *StreamOnTableModel {
	s.Table = tfconfig.StringVariable(table)
	return s
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (s *StreamOnTableModel) WithAppendOnlyValue(value tfconfig.Variable) *StreamOnTableModel {
	s.AppendOnly = value
	return s
}

func (s *StreamOnTableModel) WithAtValue(value tfconfig.Variable) *StreamOnTableModel {
	s.At = value
	return s
}

func (s *StreamOnTableModel) WithBeforeValue(value tfconfig.Variable) *StreamOnTableModel {
	s.Before = value
	return s
}

func (s *StreamOnTableModel) WithCommentValue(value tfconfig.Variable) *StreamOnTableModel {
	s.Comment = value
	return s
}

func (s *StreamOnTableModel) WithCopyGrantsValue(value tfconfig.Variable) *StreamOnTableModel {
	s.CopyGrants = value
	return s
}

func (s *StreamOnTableModel) WithDatabaseValue(value tfconfig.Variable) *StreamOnTableModel {
	s.Database = value
	return s
}

func (s *StreamOnTableModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *StreamOnTableModel {
	s.FullyQualifiedName = value
	return s
}

func (s *StreamOnTableModel) WithNameValue(value tfconfig.Variable) *StreamOnTableModel {
	s.Name = value
	return s
}

func (s *StreamOnTableModel) WithSchemaValue(value tfconfig.Variable) *StreamOnTableModel {
	s.Schema = value
	return s
}

func (s *StreamOnTableModel) WithShowInitialRowsValue(value tfconfig.Variable) *StreamOnTableModel {
	s.ShowInitialRows = value
	return s
}

func (s *StreamOnTableModel) WithTableValue(value tfconfig.Variable) *StreamOnTableModel {
	s.Table = value
	return s
}
