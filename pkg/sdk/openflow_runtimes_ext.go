package sdk

import "strings"

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
