package sdk

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// openflowRuntimeNodeTypeTiers is the only place that knows which tiers each node type has. Snowflake
// rejects a tier that does not belong to the runtime's node type. The first tier of a node type is the one
// CREATE OPENFLOW RUNTIME stores when NODE_TYPE_TIER is omitted.
var openflowRuntimeNodeTypeTiers = map[OpenflowRuntimeNodeType][]OpenflowRuntimeNodeTypeTier{
	OpenflowRuntimeNodeTypeSmall:  {OpenflowRuntimeNodeTypeTierS1, OpenflowRuntimeNodeTypeTierS2, OpenflowRuntimeNodeTypeTierS3},
	OpenflowRuntimeNodeTypeMedium: {OpenflowRuntimeNodeTypeTierM4, OpenflowRuntimeNodeTypeTierM6},
	OpenflowRuntimeNodeTypeLarge:  {OpenflowRuntimeNodeTypeTierL8},
}

// OpenflowRuntimeNodeTypeTiersFor returns the tiers that belong to the node type, or nil for an unknown node type.
func OpenflowRuntimeNodeTypeTiersFor(nodeType OpenflowRuntimeNodeType) []OpenflowRuntimeNodeTypeTier {
	return slices.Clone(openflowRuntimeNodeTypeTiers[nodeType])
}

// BelongsTo reports whether the tier is legal for the node type.
func (t OpenflowRuntimeNodeTypeTier) BelongsTo(nodeType OpenflowRuntimeNodeType) bool {
	return slices.Contains(openflowRuntimeNodeTypeTiers[nodeType], t)
}

// ValidateOpenflowRuntimeNodeTypeTier returns an error listing the valid tiers when the tier does not belong to the node type.
func ValidateOpenflowRuntimeNodeTypeTier(nodeType OpenflowRuntimeNodeType, tier OpenflowRuntimeNodeTypeTier) error {
	if tier.BelongsTo(nodeType) {
		return nil
	}
	tiers := openflowRuntimeNodeTypeTiers[nodeType]
	if len(tiers) == 0 {
		return fmt.Errorf("node type tier %s does not belong to node type %s: the node type has no tiers", tier, nodeType)
	}
	return fmt.Errorf("node type tier %s does not belong to node type %s: valid tiers are %s", tier, nodeType, collections.JoinStrings(tiers, ", "))
}

// ALTER does not validate the tier against the node type: the node type cannot be changed, and the request does not
// carry it, so Snowflake checks the tier against the runtime's node type. NODE_TYPE_TIER has no UNSET.
func (opts *CreateOpenflowRuntimeOptions) additionalValidations() error {
	return validateOpenflowRuntimeNodeTypeTier(opts.NodeType, opts.NodeTypeTier)
}

func validateOpenflowRuntimeNodeTypeTier(nodeType OpenflowRuntimeNodeType, tier *OpenflowRuntimeNodeTypeTier) error {
	if tier == nil {
		return nil
	}
	return ValidateOpenflowRuntimeNodeTypeTier(nodeType, *tier)
}

// DefaultOpenflowRuntimeNodeTypeTier is the tier CREATE OPENFLOW RUNTIME stores when NODE_TYPE_TIER is
// omitted. ALTER has no UNSET for the property, so this is also the value used to put a runtime back to that default.
func DefaultOpenflowRuntimeNodeTypeTier(nodeType OpenflowRuntimeNodeType) (OpenflowRuntimeNodeTypeTier, bool) {
	tiers := openflowRuntimeNodeTypeTiers[nodeType]
	if len(tiers) == 0 {
		return "", false
	}
	return tiers[0], true
}

// ParseOpenflowRuntimeExternalAccessIntegrations parses the external_access_integrations column of
// SHOW and DESCRIBE OPENFLOW RUNTIME, which Snowflake returns as a JSON array of quoted names.
//
// A runtime with no integrations reports the column three different ways depending on how it got there,
// and all three have to mean "empty". A freshly created runtime returns SQL NULL, but removing the last
// integration with ALTER ... REMOVE EXTERNAL_ACCESS_INTEGRATIONS leaves the literal string `null`, and an
// empty array is plausible too. Handing any of those to the comma-separated parser yields one bogus
// integration, in the `null` case one actually named "null", which then shows up as a permanent diff.
func ParseOpenflowRuntimeExternalAccessIntegrations(value string) ([]AccountObjectIdentifier, error) {
	switch strings.TrimSpace(value) {
	case "", "null", "[]":
		return nil, nil
	}
	return ParseCommaSeparatedAccountObjectIdentifierArray(value)
}

// OpenflowRuntimeFailureStatuses are terminal: a runtime in one of these will never move on its own.
var OpenflowRuntimeFailureStatuses = []OpenflowRuntimeStatus{
	OpenflowRuntimeStatusCreateFailed,
	OpenflowRuntimeStatusUpdateFailed,
	OpenflowRuntimeStatusSuspendFailed,
	OpenflowRuntimeStatusActivateFailed,
	OpenflowRuntimeStatusDeleteFailed,
	OpenflowRuntimeStatusRestartFailed,
	OpenflowRuntimeStatusUpgradeFailed,
	OpenflowRuntimeStatusMigrationFailed,
	OpenflowRuntimeStatusRollbackFailed,
}

// OpenflowRuntimeTransientStatuses are still in flight. Snowflake refuses SUSPEND, TERMINATE and DROP
// while a runtime is in one of these.
var OpenflowRuntimeTransientStatuses = []OpenflowRuntimeStatus{
	OpenflowRuntimeStatusCreating,
	OpenflowRuntimeStatusUpdating,
	OpenflowRuntimeStatusSuspending,
	OpenflowRuntimeStatusActivating,
	OpenflowRuntimeStatusRestarting,
	OpenflowRuntimeStatusUpgrading,
	OpenflowRuntimeStatusCleaningUp,
	OpenflowRuntimeStatusCancelRequested,
	OpenflowRuntimeStatusGeneratingDiagnosticBundle,
}

func (r *CreateOpenflowRuntimeRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// additionalConvert is required because OpenflowRuntimeDetails carries a plain-only Id field. There is
// nothing to derive from the row: Id is set by the caller after Describe returns.
func (r openflowRuntimeDetailsRow) additionalConvert(_ *OpenflowRuntimeDetails) error {
	return nil
}

// ID returns the identifier threaded in by the caller. DESCRIBE OPENFLOW RUNTIME does not return
// database_name or schema_name, so it cannot be derived from the response.
func (d *OpenflowRuntimeDetails) ID() SchemaObjectIdentifier {
	return d.Id
}
