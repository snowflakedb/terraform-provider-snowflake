// todo: remove this
// https://snowflakecomputing.atlassian.net/browse/SNOW-1205579
package snowflake

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/jmoiron/sqlx"
)

// TagBuilder abstracts the creation of SQL queries for a Snowflake tag.
type TagBuilder struct {
	tagId           sdk.SchemaObjectIdentifier
	maskingPolicyId sdk.SchemaObjectIdentifier
}

func (tb *TagBuilder) WithMaskingPolicy(mid sdk.SchemaObjectIdentifier) *TagBuilder {
	tb.maskingPolicyId = mid
	return tb
}

// Tag returns a pointer to a Builder that abstracts the DDL operations for a tag.
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/object-tagging.html)
func NewTagBuilder(tid sdk.SchemaObjectIdentifier) *TagBuilder {
	return &TagBuilder{
		tagId: tid,
	}
}

// Returns sql to show a tag with a specific policy attached to it.
func (tb *TagBuilder) ShowAttachedPolicy() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`SELECT * from table ("%v".information_schema.policy_references(ref_entity_name => '%v', ref_entity_domain => 'TAG')) where policy_db='%v' and policy_schema='%v' and policy_name='%v'`, tb.tagId.DatabaseName(), tb.tagId.FullyQualifiedName(), tb.maskingPolicyId.DatabaseName(), tb.maskingPolicyId.SchemaName(), tb.maskingPolicyId.Name()))

	return q.String()
}

type Tag struct {
	Name          sql.NullString `db:"name"`
	DatabaseName  sql.NullString `db:"database_name"`
	SchemaName    sql.NullString `db:"schema_name"`
	Comment       sql.NullString `db:"comment"`
	AllowedValues sql.NullString `db:"allowed_values"`
}

type TagPolicyAttachment struct {
	PolicyDB        sql.NullString `db:"POLICY_DB"`
	PolicySchema    sql.NullString `db:"POLICY_SCHEMA"`
	PolicyName      sql.NullString `db:"POLICY_NAME"`
	PolicyKind      sql.NullString `db:"POLICY_KIND"`
	RefDB           sql.NullString `db:"REF_DATABASE_NAME"`
	RefSchema       sql.NullString `db:"REF_SCHEMA_NAME"`
	RefEntity       sql.NullString `db:"REF_ENTITY_NAME"`
	RefEntityDomain sql.NullString `db:"REF_ENTITY_DOMAIN"`
}

type TagValue struct {
	Name     string
	Database string
	Schema   string
	Value    string
}

func ScanTagPolicy(row *sqlx.Row) (*TagPolicyAttachment, error) {
	r := &TagPolicyAttachment{}
	err := row.StructScan(r)
	return r, err
}
