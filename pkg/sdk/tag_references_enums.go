package sdk

import (
	"fmt"
	"slices"
	"strings"
)

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
	s = strings.ToUpper(s)
	if !slices.Contains(AllTagReferenceObjectDomains, TagReferenceObjectDomain(s)) {
		return "", fmt.Errorf("invalid TagReferenceObjectDomain: %s", s)
	}
	return TagReferenceObjectDomain(s), nil
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
	s = strings.ToUpper(s)
	if !slices.Contains(AllTagReferenceApplyMethods, TagReferenceApplyMethod(s)) {
		return "", fmt.Errorf("invalid TagReferenceApplyMethod: %s", s)
	}
	return TagReferenceApplyMethod(s), nil
}
