package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
)

type TagReferences interface {
	GetForEntity(ctx context.Context, request *GetForEntityTagReferenceRequest) ([]TagReference, error)
}

// getForEntityTagReferenceOptions is based on https://docs.snowflake.com/en/sql-reference/functions/tag_references
type getForEntityTagReferenceOptions struct {
	selectEverythingFrom bool                    `ddl:"static" sql:"SELECT * FROM TABLE"`
	parameters           *tagReferenceParameters `ddl:"list,parentheses,no_comma"`
}

type tagReferenceParameters struct {
	functionFullyQualifiedName bool                           `ddl:"static" sql:"SNOWFLAKE.INFORMATION_SCHEMA.TAG_REFERENCES"`
	arguments                  *tagReferenceFunctionArguments `ddl:"list,parentheses"`
}

type tagReferenceFunctionArguments struct {
	objectName   *string                   `ddl:"keyword,single_quotes"`
	objectDomain *TagReferenceObjectDomain `ddl:"keyword,single_quotes"`
}

type TagReferenceObjectDomain string

const (
	TagReferenceObjectDomainAccount          TagReferenceObjectDomain = "ACCOUNT"
	TagReferenceObjectDomainAlert            TagReferenceObjectDomain = "ALERT"
	TagReferenceObjectDomainColumn           TagReferenceObjectDomain = "COLUMN"
	TagReferenceObjectDomainComputePool      TagReferenceObjectDomain = "COMPUTE POOL"
	TagReferenceObjectDomainDatabase         TagReferenceObjectDomain = "DATABASE"
	TagReferenceObjectDomainDatabaseRole     TagReferenceObjectDomain = "DATABASE ROLE"
	TagReferenceObjectDomainFailoverGroup    TagReferenceObjectDomain = "FAILOVER GROUP"
	TagReferenceObjectDomainFunction         TagReferenceObjectDomain = "FUNCTION"
	TagReferenceObjectDomainIntegration      TagReferenceObjectDomain = "INTEGRATION"
	TagReferenceObjectDomainNetworkPolicy    TagReferenceObjectDomain = "NETWORK POLICY"
	TagReferenceObjectDomainProcedure        TagReferenceObjectDomain = "PROCEDURE"
	TagReferenceObjectDomainReplicationGroup TagReferenceObjectDomain = "REPLICATION GROUP"
	TagReferenceObjectDomainRole             TagReferenceObjectDomain = "ROLE"
	TagReferenceObjectDomainSchema           TagReferenceObjectDomain = "SCHEMA"
	TagReferenceObjectDomainShare            TagReferenceObjectDomain = "SHARE"
	TagReferenceObjectDomainStage            TagReferenceObjectDomain = "STAGE"
	TagReferenceObjectDomainStream           TagReferenceObjectDomain = "STREAM"
	TagReferenceObjectDomainTable            TagReferenceObjectDomain = "TABLE"
	TagReferenceObjectDomainTask             TagReferenceObjectDomain = "TASK"
	TagReferenceObjectDomainUser             TagReferenceObjectDomain = "USER"
	TagReferenceObjectDomainWarehouse        TagReferenceObjectDomain = "WAREHOUSE"
)

var AllTagReferenceObjectDomains = []TagReferenceObjectDomain{
	TagReferenceObjectDomainAccount,
	TagReferenceObjectDomainAlert,
	TagReferenceObjectDomainColumn,
	TagReferenceObjectDomainComputePool,
	TagReferenceObjectDomainDatabase,
	TagReferenceObjectDomainDatabaseRole,
	TagReferenceObjectDomainFailoverGroup,
	TagReferenceObjectDomainFunction,
	TagReferenceObjectDomainIntegration,
	TagReferenceObjectDomainNetworkPolicy,
	TagReferenceObjectDomainProcedure,
	TagReferenceObjectDomainReplicationGroup,
	TagReferenceObjectDomainRole,
	TagReferenceObjectDomainSchema,
	TagReferenceObjectDomainShare,
	TagReferenceObjectDomainStage,
	TagReferenceObjectDomainStream,
	TagReferenceObjectDomainTable,
	TagReferenceObjectDomainTask,
	TagReferenceObjectDomainUser,
	TagReferenceObjectDomainWarehouse,
}

func ToTagReferenceObjectDomain(s string) (TagReferenceObjectDomain, error) {
	d := TagReferenceObjectDomain(strings.ToUpper(s))
	if !slices.Contains(AllTagReferenceObjectDomains, d) {
		return "", fmt.Errorf("invalid TagReferenceObjectDomain: %s", s)
	}
	return d, nil
}

type TagReferenceApplyMethod string

const (
	TagReferenceApplyMethodClassified TagReferenceApplyMethod = "CLASSIFIED"
	TagReferenceApplyMethodInherited  TagReferenceApplyMethod = "INHERITED"
	TagReferenceApplyMethodManual     TagReferenceApplyMethod = "MANUAL"
	TagReferenceApplyMethodPropagated TagReferenceApplyMethod = "PROPAGATED"
)

var AllTagReferenceApplyMethods = []TagReferenceApplyMethod{
	TagReferenceApplyMethodClassified,
	TagReferenceApplyMethodInherited,
	TagReferenceApplyMethodManual,
	TagReferenceApplyMethodPropagated,
}

func ToTagReferenceApplyMethod(s string) (TagReferenceApplyMethod, error) {
	m := TagReferenceApplyMethod(strings.ToUpper(s))
	if !slices.Contains(AllTagReferenceApplyMethods, m) {
		return "", fmt.Errorf("invalid TagReferenceApplyMethod: %s", s)
	}
	return m, nil
}

type TagReference struct {
	TagDatabase    string
	TagSchema      string
	TagName        string
	TagValue       string
	Level          TagReferenceObjectDomain
	ObjectDatabase *string
	ObjectSchema   *string
	ObjectName     string
	Domain         TagReferenceObjectDomain
	ColumnName     *string
	ApplyMethod    TagReferenceApplyMethod
}

type tagReferenceDBRow struct {
	TagDatabase    string         `db:"TAG_DATABASE"`
	TagSchema      string         `db:"TAG_SCHEMA"`
	TagName        string         `db:"TAG_NAME"`
	TagValue       string         `db:"TAG_VALUE"`
	Level          string         `db:"LEVEL"`
	ObjectDatabase sql.NullString `db:"OBJECT_DATABASE"`
	ObjectSchema   sql.NullString `db:"OBJECT_SCHEMA"`
	ObjectName     string         `db:"OBJECT_NAME"`
	Domain         string         `db:"DOMAIN"`
	ColumnName     sql.NullString `db:"COLUMN_NAME"`
	ApplyMethod    string         `db:"APPLY_METHOD"`
}

func (row tagReferenceDBRow) convert() (*TagReference, error) {
	tagReference := TagReference{
		TagDatabase: row.TagDatabase,
		TagSchema:   row.TagSchema,
		TagName:     row.TagName,
		TagValue:    row.TagValue,
		ObjectName:  row.ObjectName,
	}
	mapStringWithMapping(&tagReference.Level, row.Level, ToTagReferenceObjectDomain)
	mapStringWithMapping(&tagReference.Domain, row.Domain, ToTagReferenceObjectDomain)
	mapStringWithMapping(&tagReference.ApplyMethod, row.ApplyMethod, ToTagReferenceApplyMethod)
	mapNullString(&tagReference.ObjectDatabase, row.ObjectDatabase)
	mapNullString(&tagReference.ObjectSchema, row.ObjectSchema)
	mapNullString(&tagReference.ColumnName, row.ColumnName)
	return &tagReference, nil
}
