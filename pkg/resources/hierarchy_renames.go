package resources

import (
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
func handleHierarchyRenameIdUpdate(d *schema.ResourceData, encodeIdFn func() string, caseDescription string) {
	log.Printf("[DEBUG] %s - no Snowflake modification needed, updating the id...", caseDescription)
	d.SetId(encodeIdFn())
}

// handleHierarchyMove handles the "move" case: performs the ALTER to move the object
// to its new location and updates the Terraform state ID.
func handleHierarchyMove(d *schema.ResourceData, encodeIdFn func() string, currentId, targetId sdk.ObjectIdentifier, alterFn func() error, caseDescription string) diag.Diagnostics {
	log.Printf("[DEBUG] %s - executing ALTER RENAME TO from %s to %s...", caseDescription, currentId.FullyQualifiedName(), targetId.FullyQualifiedName())
	if err := alterFn(); err != nil {
		d.Partial(true)
		return diag.FromErr(fmt.Errorf("failed to move from %s to %s: %w", currentId.FullyQualifiedName(), targetId.FullyQualifiedName(), err))
	}
	d.SetId(encodeIdFn())
	return nil
}
