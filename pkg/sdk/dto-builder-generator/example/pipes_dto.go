package example

import (
	"bytes"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

//go:generate go run ../main.go

type CreatePipeRequest struct {
	orReplace        bool
	ifNotExists      bool
	name             sdk.SchemaObjectIdentifier // required
	autoIngest       bool
	errorIntegration string
	awsSnsTopic      string
	integration      string
	comment          string
	copyStatement    string // required
	exampleOfImport  bytes.Buffer
}

type AlterPipeRequest struct {
	ifExists  bool
	name      sdk.SchemaObjectIdentifier // required
	set       PipeSetRequest
	unset     PipeUnsetRequest
	setTags   PipeSetTagsRequest
	unsetTags PipeUnsetTagsRequest
	refresh   PipeRefreshRequest
}

type PipeSetRequest struct {
	errorIntegration    string
	pipeExecutionPaused bool
	comment             string
}

type PipeUnsetRequest struct {
	pipeExecutionPaused bool
	comment             bool
}

type PipeSetTagsRequest struct {
	tag []sdk.TagAssociation // required
}

type PipeUnsetTagsRequest struct {
	tag []sdk.ObjectIdentifier // required
}

type PipeRefreshRequest struct {
	prefix        string
	modifiedAfter string
}
