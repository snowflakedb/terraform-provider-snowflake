package resources

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func isRenameOfTheGivenLevelInTheHierarchy(newParentExists, oldParentExists, objectAtTargetLocationExists bool) bool {
	return newParentExists && !oldParentExists && objectAtTargetLocationExists
}

func isMoveToADifferentObjectOnTheGivenLevelInTheHierarchy(newParentExists, oldParentExists, objectAtSourceLocationExists bool) bool {
	return newParentExists && oldParentExists && objectAtSourceLocationExists
}

// handleHierarchyRenameIdUpdate handles the "rename" case: a parent was renamed,
// so only the Terraform state ID needs updating (no Snowflake ALTER needed).
func handleHierarchyRenameIdUpdate(d *schema.ResourceData, encodeIdFn func() string, logMsg string) {
	log.Printf("[DEBUG] %s", logMsg)
	d.SetId(encodeIdFn())
}

// handleHierarchyMove handles the "move" case: performs the ALTER to move the object
// to its new location and updates the Terraform state ID.
func handleHierarchyMove(d *schema.ResourceData, encodeIdFn func() string, currentFQN, targetFQN string, alterFn func() error, logMsg string) diag.Diagnostics {
	log.Printf("[DEBUG] %s", logMsg)
	if err := alterFn(); err != nil {
		d.Partial(true)
		return diag.FromErr(fmt.Errorf("failed to move from %s to %s: %w", currentFQN, targetFQN, err))
	}
	d.SetId(encodeIdFn())
	return nil
}
