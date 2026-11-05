package sdk

import (
	"context"
	"encoding/csv"
	"fmt"
	"slices"
	"strings"
	"unicode"
)

// OpenflowDeploymentFailureStatuses are terminal: a deployment in one of these will never move on its own.
var OpenflowDeploymentFailureStatuses = []OpenflowDeploymentStatus{
	OpenflowDeploymentStatusCreateFailed,
	OpenflowDeploymentStatusDeleteFailed,
	OpenflowDeploymentStatusUpgradeFailed,
	OpenflowDeploymentStatusMigrationFailed,
	OpenflowDeploymentStatusRollbackFailed,
}

// OpenflowDeploymentTransientStatuses are still coming up. Only a settled deployment accepts ALTER SET,
// TERMINATE or DROP. SPCS goes CREATING -> PROVISIONING -> ACTIVE, so PROVISIONING is transient too; BYOC
// goes straight to INACTIVE.
var OpenflowDeploymentTransientStatuses = []OpenflowDeploymentStatus{
	OpenflowDeploymentStatusCreating,
	OpenflowDeploymentStatusProvisioning,
}

func (r *CreateOpenflowDeploymentRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (d *OpenflowDeploymentDetails) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(d.Name)
}

// ShowParameters returns the parameters visible at the OPENFLOW DEPLOYMENT level for the given
// deployment. EVENT_TABLE is settable on CREATE and via ALTER SET/UNSET but is absent from both
// SHOW OPENFLOW DEPLOYMENTS and DESCRIBE OPENFLOW DEPLOYMENT, so this is the only way to read it back.
// Mirrors pkg/sdk/hybrid_tables_ext.go's ShowParameters with ParametersIn.OpenflowDeployment.
func (v *openflowDeployments) ShowParameters(ctx context.Context, id AccountObjectIdentifier) ([]*Parameter, error) {
	params, err := v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			OpenflowDeployment: id,
		},
	})
	if err != nil {
		return nil, err
	}
	for _, p := range params {
		if p.Key != "EVENT_TABLE" {
			continue
		}
		if err := quoteOpenflowEventTableParameter(p); err != nil {
			return nil, err
		}
	}
	return params, nil
}

// quoteOpenflowEventTableParameter is temporary: SHOW PARAMETERS returns EVENT_TABLE as
// snowflake.db.schema or snowflake.db."schema", and ParseSchemaObjectIdentifier would keep
// unquoted parts lowercase. Remove when Snowflake returns a quoted fully qualified name.
// Empty and NONE are left as-is; the standard parser cannot read NONE.
func quoteOpenflowEventTableParameter(p *Parameter) error {
	value, err := quoteOpenflowEventTableValue(p.Value)
	if err != nil {
		return fmt.Errorf("parsing parameter %q value: %w", p.Key, err)
	}
	def, err := quoteOpenflowEventTableValue(p.Default)
	if err != nil {
		return fmt.Errorf("parsing parameter %q default: %w", p.Key, err)
	}
	p.Value = value
	p.Default = def
	return nil
}

func quoteOpenflowEventTableValue(value string) (string, error) {
	if value == "" || strings.EqualFold(value, "NONE") {
		return value, nil
	}
	parts, err := ParseIdentifierStringWithOpts(foldUnquotedIdentifierParts(value), func(r *csv.Reader) {
		r.Comma = IdDelimiter
	})
	if err != nil {
		return "", err
	}
	if len(parts) != 3 {
		return "", fmt.Errorf(`unexpected number of parts %d in identifier %s, expected 3 in a form of "<database_name>.<schema_name>.<schema_object_name>"`, len(parts), value)
	}
	if slices.Contains(parts, "") {
		return "", fmt.Errorf("unable to parse identifier: %s, empty part", value)
	}
	return NewSchemaObjectIdentifier(parts[0], parts[1], parts[2]).FullyQualifiedName(), nil
}

// foldUnquotedIdentifierParts uppercases characters outside double quotes, as Snowflake does for unquoted identifiers.
// Every quote toggles the state, so an escaped quote ("") inside a quoted part keeps it quoted.
func foldUnquotedIdentifierParts(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	quoted := false
	for _, r := range s {
		if r == '"' {
			quoted = !quoted
		} else if !quoted {
			r = unicode.ToUpper(r)
		}
		b.WriteRune(r)
	}
	return b.String()
}
