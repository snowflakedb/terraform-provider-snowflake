package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

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
func handleHierarchyMove[T sdk.ObjectIdentifier](d *schema.ResourceData, encodeIdFn func() string, currentId, targetId T, renameFn func(T, T) func() error, caseDescription string) diag.Diagnostics {
	log.Printf("[DEBUG] %s - executing ALTER RENAME TO from %s to %s...", caseDescription, currentId.FullyQualifiedName(), targetId.FullyQualifiedName())
	if err := renameFn(currentId, targetId)(); err != nil {
		d.Partial(true)
		return diag.FromErr(fmt.Errorf("failed to move from %s to %s: %w", currentId.FullyQualifiedName(), targetId.FullyQualifiedName(), err))
	}
	d.SetId(encodeIdFn())
	return nil
}

// handleShallowHierarchyRename handles the complete single-level rename/move detection.
// P is the parent identifier type, O is the probed object identifier type, M is the moved object type.
// PR and OR are the return types of the ShowByID functions.
func handleShallowHierarchyRename[P, O, M sdk.ObjectIdentifier, PR, OR any](
	d *schema.ResourceData,
	ctx context.Context,
	parentShowByIDFn func(context.Context, P) (PR, error),
	objectShowByIDFn func(context.Context, O) (OR, error),
	newParentId, oldParentId P,
	targetObjectId, currentObjectId O,
	movedObjectId M,
	renameFn func(string),
	moveFn func(M, string) diag.Diagnostics,
	parentLevelName string,
	childLevelName string,
	leafLevelName string,
) diag.Diagnostics {
	_, errNewParent := parentShowByIDFn(ctx, newParentId)
	newParentExists := errNewParent == nil
	_, errOldParent := parentShowByIDFn(ctx, oldParentId)
	oldParentExists := errOldParent == nil
	_, errTargetObject := objectShowByIDFn(ctx, targetObjectId)
	objectAtTargetExists := errTargetObject == nil
	_, errCurrentObject := objectShowByIDFn(ctx, currentObjectId)
	objectAtSourceExists := errCurrentObject == nil

	switch {
	case isRenameOfTheGivenLevelInTheHierarchy(newParentExists, oldParentExists, objectAtTargetExists):
		renameFn(fmt.Sprintf("%s was renamed for %s", capitalize(parentLevelName), leafLevelName))
		return nil
	case isMoveToADifferentObjectOnTheGivenLevelInTheHierarchy(newParentExists, oldParentExists, objectAtSourceExists):
		return moveFn(movedObjectId, fmt.Sprintf("Moving %s to different %s", leafLevelName, parentLevelName))
	default:
		d.Partial(true)
		return diag.FromErr(fmt.Errorf(
			"unknown rename use case: old %s %s (exists: %t), new %s %s (exists: %t), old %s %s (exists: %t), new %s %s (exists: %t). See https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/object_renaming_guide",
			parentLevelName, oldParentId.FullyQualifiedName(), oldParentExists,
			parentLevelName, newParentId.FullyQualifiedName(), newParentExists,
			childLevelName, currentObjectId.FullyQualifiedName(), objectAtSourceExists,
			childLevelName, targetObjectId.FullyQualifiedName(), objectAtTargetExists,
		))
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
