package sdk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

func (v *tables) ShowParameters(ctx context.Context, id SchemaObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Table: id,
		},
	})
}

func (r *CreateTableRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// SequenceName is a SET DEFAULT value that renders as <name>.NEXTVAL.
// OptionalEnumLegacy does not emit the type, so it lives here next to String().
type SequenceName string

func (sn SequenceName) String() string {
	return fmt.Sprintf("%s.NEXTVAL", string(sn))
}

// GetClusterByKeys converts the SHOW TABLES result for ClusterBy and converts it to list of keys.
func (v *Table) GetClusterByKeys() []string {
	if v.ClusterBy == "" {
		return nil
	}

	statementWithoutLinear := strings.TrimSuffix(strings.Replace(v.ClusterBy, "LINEAR(", "", 1), ")")
	return collections.Map(splitOuter(statementWithoutLinear, '(', ')', true), strings.TrimSpace)
}

func (opts *CreateTableOptions) additionalValidations() error {
	var errs []error
	if len(opts.ColumnsAndConstraints.Columns) == 0 {
		errs = append(errs, errNotSet("CreateTableOptions", "Columns"))
	}
	if opts.RowAccessPolicy != nil {
		if !ValidObjectIdentifier(opts.RowAccessPolicy.Name) {
			errs = append(errs, errInvalidIdentifier("CreateTableOptions.RowAccessPolicy", "Name"))
		}
	}
	return JoinErrors(errs...)
}

// The ADD branch of ALTER TABLE ... ADD COLUMN for external tables is the only ExternalTableAction branch, so it is required.
// ChildrenRule needs two alternatives, so the SSOT definition cannot express this.
func (opts *AlterTableOptions) additionalValidations() error {
	if valueSet(opts.ExternalTableAction) && !valueSet(opts.ExternalTableAction.Add) {
		return errNotSet("AlterTableOptions.ExternalTableAction", "Add")
	}
	return nil
}

func (opts *ShowTableOptions) additionalValidations() error {
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		return ErrPatternRequiredForLikeKeyword
	}
	return nil
}

// additionalConvert splits COLLATE off the DESCRIBE TABLE type column.
// Type stays nil when the type cannot be parsed; TypeRaw always holds the type without COLLATE.
func (r tableColumnDetailsRow) additionalConvert(result *TableColumnDetails) error {
	typeRaw, collation := r.splitTypeAndCollation()
	result.TypeRaw = typeRaw
	result.Collation = collation
	if dataType, err := datatypes.ParseDataType(typeRaw); err == nil {
		result.Type = dataType
	}
	return nil
}

func (r tableColumnDetailsRow) splitTypeAndCollation() (string, *string) {
	collateRegexp := regexp.MustCompile(`COLLATE +'([a-zA-Z0-9_-]*)'`)
	matches := collateRegexp.FindStringSubmatch(r.Type)

	if len(matches) == 2 {
		collation := matches[1]
		return strings.TrimSpace(collateRegexp.ReplaceAllString(r.Type, "")), &collation
	}
	return r.Type, nil
}

// TableColumnIdentityDetails is an IDENTITY column default as DESCRIBE TABLE returns it, for example
// IDENTITY START 1 INCREMENT 1 ORDER. Order is nil when the value has no ORDER or NOORDER keyword.
type TableColumnIdentityDetails struct {
	Start     int
	Increment int
	Order     *bool
}

// ParseColumnIdentity parses an IDENTITY column default. It returns false for any other default.
func ParseColumnIdentity(raw string) (*TableColumnIdentityDetails, bool) {
	fields := strings.Fields(strings.ToUpper(raw))
	if len(fields) < 5 || fields[0] != "IDENTITY" || fields[1] != "START" || fields[3] != "INCREMENT" {
		return nil, false
	}
	start, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, false
	}
	increment, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, false
	}
	identity := &TableColumnIdentityDetails{Start: start, Increment: increment}
	switch {
	case len(fields) == 5:
	case len(fields) == 6 && fields[5] == "ORDER":
		identity.Order = Bool(true)
	case len(fields) == 6 && fields[5] == "NOORDER":
		identity.Order = Bool(false)
	default:
		return nil, false
	}
	return identity, true
}

// TypeString returns the SQL of the parsed type, falling back to the raw DESCRIBE value when parsing failed.
func (d TableColumnDetails) TypeString() string {
	if d.Type != nil {
		return d.Type.ToSql()
	}
	return d.TypeRaw
}

// TableConstraint is one primary, unique, or foreign key of a table, merged from its SHOW ... KEYS rows.
type TableConstraint = HybridTableConstraint

// GetConstraints reads the keys with SHOW PRIMARY KEYS, SHOW UNIQUE KEYS, and SHOW IMPORTED KEYS, ordered by kind,
// then by column names.
// TODO [next PRs]: the hybrid table SDK implements these commands, and they work the same on any table, so this calls
// it for now. Move the implementation here and make hybridTables call this instead.
func (v *tables) GetConstraints(ctx context.Context, id SchemaObjectIdentifier) ([]TableConstraint, error) {
	return v.client.HybridTables.GetConstraints(ctx, id)
}
