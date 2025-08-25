package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type PolicyReferences interface {
	GetForEntity(ctx context.Context, request *GetForEntityPolicyReferenceRequest) ([]PolicyReference, error)
}

type getForEntityPolicyReferenceOptions struct {
	selectEverythingFrom bool                       `ddl:"static" sql:"SELECT * FROM TABLE"`
	parameters           *policyReferenceParameters `ddl:"list,parentheses,no_comma"`
}

type policyReferenceParameters struct {
	functionFullyQualifiedName bool                              `ddl:"static" sql:"SNOWFLAKE.INFORMATION_SCHEMA.POLICY_REFERENCES"`
	arguments                  *policyReferenceFunctionArguments `ddl:"list,parentheses"`
}

type PolicyEntityDomain string

const (
	PolicyEntityDomainAccount     PolicyEntityDomain = "ACCOUNT"
	PolicyEntityDomainIntegration PolicyEntityDomain = "INTEGRATION"
	PolicyEntityDomainTable       PolicyEntityDomain = "TABLE"
	PolicyEntityDomainTag         PolicyEntityDomain = "TAG"
	PolicyEntityDomainUser        PolicyEntityDomain = "USER"
	PolicyEntityDomainView        PolicyEntityDomain = "VIEW"
)

var AllPolicyEntityDomains = []PolicyEntityDomain{
	PolicyEntityDomainAccount,
	PolicyEntityDomainIntegration,
	PolicyEntityDomainTable,
	PolicyEntityDomainTag,
	PolicyEntityDomainUser,
	PolicyEntityDomainView,
}

func ToPolicyEntityDomain(s string) (PolicyEntityDomain, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(PolicyEntityDomainAccount):
		return PolicyEntityDomainAccount, nil
	case string(PolicyEntityDomainIntegration):
		return PolicyEntityDomainIntegration, nil
	case string(PolicyEntityDomainTable):
		return PolicyEntityDomainTable, nil
	case string(PolicyEntityDomainTag):
		return PolicyEntityDomainTag, nil
	case string(PolicyEntityDomainUser):
		return PolicyEntityDomainUser, nil
	case string(PolicyEntityDomainView):
		return PolicyEntityDomainView, nil
	default:
		return "", fmt.Errorf("invalid PolicyEntityDomain: %s", s)
	}
}

type policyReferenceFunctionArguments struct {
	refEntityName   []ObjectIdentifier  `ddl:"parameter,single_quotes,arrow_equals" sql:"REF_ENTITY_NAME"`
	refEntityDomain *PolicyEntityDomain `ddl:"parameter,single_quotes,arrow_equals" sql:"REF_ENTITY_DOMAIN"`
}

type PolicyKind string

const (
	PolicyKindAggregationPolicy    PolicyKind = "AGGREGATION_POLICY"
	PolicyKindAuthenticationPolicy PolicyKind = "AUTHENTICATION_POLICY"
	PolicyKindFeaturePolicy        PolicyKind = "FEATURE_POLICY"
	PolicyKindMaskingPolicy        PolicyKind = "MASKING_POLICY"
	PolicyKindPackagesPolicy       PolicyKind = "PACKAGES_POLICY"
	PolicyKindPasswordPolicy       PolicyKind = "PASSWORD_POLICY"
	PolicyKindProjectionPolicy     PolicyKind = "PROJECTION_POLICY"
	PolicyKindRowAccessPolicy      PolicyKind = "ROW_ACCESS_POLICY"
	PolicyKindSessionPolicy        PolicyKind = "SESSION_POLICY"
)

type PolicyReference struct {
	PolicyDb          *string
	PolicySchema      *string
	PolicyName        string
	PolicyKind        PolicyKind
	RefDatabaseName   *string
	RefSchemaName     *string
	RefEntityName     string
	RefEntityDomain   string
	RefColumnName     *string
	RefArgColumnNames *string
	TagDatabase       *string
	TagSchema         *string
	TagName           *string
	PolicyStatus      *string
}

type policyReferenceDBRow struct {
	PolicyDb          sql.NullString `db:"POLICY_DB"`
	PolicySchema      sql.NullString `db:"POLICY_SCHEMA"`
	PolicyName        string         `db:"POLICY_NAME"`
	PolicyKind        string         `db:"POLICY_KIND"`
	RefDatabaseName   sql.NullString `db:"REF_DATABASE_NAME"`
	RefSchemaName     sql.NullString `db:"REF_SCHEMA_NAME"`
	RefEntityName     string         `db:"REF_ENTITY_NAME"`
	RefEntityDomain   string         `db:"REF_ENTITY_DOMAIN"`
	RefColumnName     sql.NullString `db:"REF_COLUMN_NAME"`
	RefArgColumnNames sql.NullString `db:"REF_ARG_COLUMN_NAMES"`
	TagDatabase       sql.NullString `db:"TAG_DATABASE"`
	TagSchema         sql.NullString `db:"TAG_SCHEMA"`
	TagName           sql.NullString `db:"TAG_NAME"`
	PolicyStatus      sql.NullString `db:"POLICY_STATUS"`
}

func (row policyReferenceDBRow) convert() (*PolicyReference, error) {
	policyReference := PolicyReference{
		PolicyName:      row.PolicyName,
		PolicyKind:      PolicyKind(row.PolicyKind),
		RefEntityName:   row.RefEntityName,
		RefEntityDomain: row.RefEntityDomain,
	}
	if row.PolicyDb.Valid {
		policyReference.PolicyDb = &row.PolicyDb.String
	}
	if row.PolicySchema.Valid {
		policyReference.PolicySchema = &row.PolicySchema.String
	}
	if row.RefDatabaseName.Valid {
		policyReference.RefDatabaseName = &row.RefDatabaseName.String
	}
	if row.RefSchemaName.Valid {
		policyReference.RefSchemaName = &row.RefSchemaName.String
	}
	if row.RefColumnName.Valid {
		policyReference.RefColumnName = &row.RefColumnName.String
	}
	if row.RefArgColumnNames.Valid {
		policyReference.RefArgColumnNames = &row.RefArgColumnNames.String
	}
	if row.TagDatabase.Valid {
		policyReference.TagDatabase = &row.TagDatabase.String
	}
	if row.TagSchema.Valid {
		policyReference.TagSchema = &row.TagSchema.String
	}
	if row.TagName.Valid {
		policyReference.TagName = &row.TagName.String
	}
	if row.PolicyStatus.Valid {
		policyReference.PolicyStatus = &row.PolicyStatus.String
	}
	return &policyReference, nil
}
