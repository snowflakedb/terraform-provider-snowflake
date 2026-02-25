package model

import (
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

// TODO(SNOW-1501905): Remove after complex non-list type overrides are handled
func (d *DynamicTableModel) WithTargetLag(targetLag []sdk.TargetLag) *DynamicTableModel {
	if len(targetLag) != 1 {
		log.Panicf("expected exactly one target lag, got %d", len(targetLag))
	}

	if targetLag[0].MaximumDuration != nil {
		return d.WithTargetLagValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"maximum_duration": tfconfig.StringVariable(*targetLag[0].MaximumDuration),
		}))
	}

	if targetLag[0].Downstream != nil {
		return d.WithTargetLagValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"downstream": tfconfig.BoolVariable(*targetLag[0].Downstream),
		}))
	}

	log.Panicf("neither maximum_duration nor downstream is set in target lag: %+v", targetLag[0])
	return nil
}
