package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// IdentityStateSets are the d.Set calls of a generated SetFieldsNotSetByRead: the identifier fields, taken from the
// parsed id, in the order name, database, schema. Import runs it because Read does not write them.
// The other fields are written by the Read that follows the import.
func (m ResourceModel) IdentityStateSets() []string {
	if !m.Hooks.SetFieldsNotSetByRead.Generate {
		return nil
	}
	switch m.IdType {
	case "AccountObjectIdentifier":
		return []string{`d.Set("name", id.Name())`}
	case "DatabaseObjectIdentifier":
		return []string{`d.Set("name", id.Name())`, `d.Set("database", id.DatabaseName())`}
	case "SchemaObjectIdentifier":
		return []string{`d.Set("name", id.Name())`, `d.Set("database", id.DatabaseName())`, `d.Set("schema", id.SchemaName())`}
	default:
		panic(fmt.Sprintf("object %q: SetFieldsNotSetByRead is generated only for an account, database, or schema object identifier, not %q", m.Name, m.IdType))
	}
}

// DescribeStateSets are the d.Set calls of a generated SetDescribeOutput: the required attributes that SHOW does not
// return and DESCRIBE does, copied into state on every read. A required attribute is not drift-tracked by the generated
// external changes, so this is how Read keeps it in sync with Snowflake. Only plain string columns are supported.
func (m ResourceModel) DescribeStateSets() []string {
	if !m.Hooks.SetDescribeOutput.Generate {
		return nil
	}
	if !m.Hooks.DescribeInSdk.Present || !m.DescribeFailRead {
		panic(fmt.Sprintf("object %q: a generated SetDescribeOutput needs DescribeInSdk and DescribeFailRead, because it reads the details", m.Name))
	}
	inShow := make(map[string]struct{}, len(m.ShowColumns))
	for _, col := range m.ShowColumns {
		if col != nil {
			inShow[col.Name] = struct{}{}
		}
	}
	required := make(map[string]struct{}, len(m.Attributes))
	for _, attr := range m.Attributes {
		if attr.Required {
			required[attr.Name] = struct{}{}
		}
	}
	identity := m.identityNames()
	var sets []string
	for _, col := range m.DescribeColumns {
		if col == nil {
			continue
		}
		if _, ok := required[col.Name]; !ok {
			continue
		}
		if _, skip := inShow[col.Name]; skip {
			continue
		}
		if _, skip := identity[col.Name]; skip || col.CreateOnly {
			continue
		}
		if col.Type != ssot.AttrString || col.OptionalInShow || len(col.Children) > 0 {
			panic(fmt.Sprintf("object %q: SetDescribeOutput is not generated for column %q: only a plain string is supported", m.Name, col.Name))
		}
		sets = append(sets, fmt.Sprintf("d.Set(%q, %sDetails.%s)", col.Name, m.Prefix(), genhelpers.SnakeCaseToCamel(col.EffectiveShowName())))
	}
	return sets
}
