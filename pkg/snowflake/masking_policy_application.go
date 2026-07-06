package snowflake

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type TableColumnMaskingPolicyApplication struct {
	Table         *SchemaObjectIdentifier
	Column        string
	MaskingPolicy *SchemaObjectIdentifier
}

type TableColumn struct {
	Table  *SchemaObjectIdentifier
	Column string
}

type TableColumnMaskingPolicyApplicationManager struct{}

func NewTableColumnMaskingPolicyApplicationManager() *TableColumnMaskingPolicyApplicationManager {
	return &TableColumnMaskingPolicyApplicationManager{}
}

type TableColumnMaskingPolicyApplicationCreateInput struct {
	TableColumnMaskingPolicyApplication
}

func (m *TableColumnMaskingPolicyApplicationManager) Create(x *TableColumnMaskingPolicyApplicationCreateInput) string {
	column := sdk.DoubleQuotes.Modify(x.Column)
	return fmt.Sprintf(`ALTER TABLE IF EXISTS %s MODIFY COLUMN %s SET MASKING POLICY %s;`, x.Table.QualifiedName(), column, x.MaskingPolicy.QualifiedName())
}

type TableColumnMaskingPolicyApplicationReadInput = TableColumn

func (m *TableColumnMaskingPolicyApplicationManager) Read(x *TableColumnMaskingPolicyApplicationReadInput) string {
	return fmt.Sprintf("DESCRIBE TABLE %s TYPE = COLUMNS;", x.Table.QualifiedName())
}

func (m *TableColumnMaskingPolicyApplicationManager) Parse(rows *sql.Rows, column string) (string, error) {
	var name, sqlType, kind, null, defaultValue, primaryKey, uniqueKey, check, expression, comment, policyName, privacyDomain, schemaEvolutionRecord sql.NullString

	for rows.Next() {
		if err := rows.Scan(&name, &sqlType, &kind, &null, &defaultValue, &primaryKey, &uniqueKey, &check, &expression, &comment, &policyName, &privacyDomain, &schemaEvolutionRecord); err != nil {
			return "", err
		}

		if strings.EqualFold(name.String, column) {
			return policyName.String, nil
		}
	}
	return "", nil
}

type TableColumnMaskingPolicyApplicationDeleteInput struct {
	TableColumn
}

func (m *TableColumnMaskingPolicyApplicationManager) Delete(x *TableColumnMaskingPolicyApplicationDeleteInput) string {
	column := sdk.DoubleQuotes.Modify(x.Column)
	return fmt.Sprintf(`ALTER TABLE IF EXISTS %s MODIFY COLUMN %s UNSET MASKING POLICY;`, x.Table.QualifiedName(), column)
}
