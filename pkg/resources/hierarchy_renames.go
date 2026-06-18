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
	ctx context.Context,
	d *schema.ResourceData,
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

// handleDeepHierarchyRename handles the case where both database AND schema change
// for a schema-level object. It determines whether the database/schema were renamed
// or moved and acts accordingly.
func handleDeepHierarchyRename(
	ctx context.Context,
	d *schema.ResourceData,
	client *sdk.Client,
	newDatabaseId, oldDatabaseId sdk.AccountObjectIdentifier,
	oldSchemaName, newSchemaName string,
	objectName string,
	renameFn func(string),
	moveFn func(sdk.SchemaObjectIdentifier, string) diag.Diagnostics,
	leafLevelName string,
) diag.Diagnostics {
	oldDatabaseName := oldDatabaseId.Name()
	newDatabaseName := newDatabaseId.Name()

	// Probe databases
	_, errNewDb := client.Databases.ShowByID(ctx, newDatabaseId)
	newDatabaseExists := errNewDb == nil
	_, errOldDb := client.Databases.ShowByID(ctx, oldDatabaseId)
	oldDatabaseExists := errOldDb == nil

	// Probe schemas (4 combinations)
	schemaInNewDbOldName := sdk.NewDatabaseObjectIdentifier(newDatabaseName, oldSchemaName)
	schemaInNewDbNewName := sdk.NewDatabaseObjectIdentifier(newDatabaseName, newSchemaName)
	schemaInOldDbOldName := sdk.NewDatabaseObjectIdentifier(oldDatabaseName, oldSchemaName)
	schemaInOldDbNewName := sdk.NewDatabaseObjectIdentifier(oldDatabaseName, newSchemaName)

	_, errSchemaNewDbOld := client.Schemas.ShowByID(ctx, schemaInNewDbOldName)
	schemaNewDbOldExists := errSchemaNewDbOld == nil
	_, errSchemaNewDbNew := client.Schemas.ShowByID(ctx, schemaInNewDbNewName)
	schemaNewDbNewExists := errSchemaNewDbNew == nil
	_, errSchemaOldDbOld := client.Schemas.ShowByID(ctx, schemaInOldDbOldName)
	schemaOldDbOldExists := errSchemaOldDbOld == nil
	_, errSchemaOldDbNew := client.Schemas.ShowByID(ctx, schemaInOldDbNewName)
	schemaOldDbNewExists := errSchemaOldDbNew == nil

	switch {
	// Scenario 1: DB rename + Schema rename → just update ID
	case !oldDatabaseExists && newDatabaseExists && !schemaNewDbOldExists && schemaNewDbNewExists:
		renameFn(fmt.Sprintf("Database and schema were both renamed for %s", leafLevelName))
		return nil
	// Scenario 2: DB rename + Schema move → object is at (newDb, oldSchema)
	case !oldDatabaseExists && newDatabaseExists && schemaNewDbOldExists:
		currentId := sdk.NewSchemaObjectIdentifier(newDatabaseName, oldSchemaName, objectName)
		return moveFn(currentId, fmt.Sprintf("Database was renamed, moving %s to different schema", leafLevelName))
	// Scenario 3: DB move + Schema rename → object is at (oldDb, newSchema)
	case oldDatabaseExists && newDatabaseExists && !schemaOldDbOldExists && schemaOldDbNewExists:
		currentId := sdk.NewSchemaObjectIdentifier(oldDatabaseName, newSchemaName, objectName)
		return moveFn(currentId, fmt.Sprintf("Schema was renamed, moving %s to different database", leafLevelName))
	// Scenario 4: DB move + Schema move → object is at (oldDb, oldSchema)
	case oldDatabaseExists && newDatabaseExists && schemaOldDbOldExists:
		currentId := sdk.NewSchemaObjectIdentifier(oldDatabaseName, oldSchemaName, objectName)
		return moveFn(currentId, fmt.Sprintf("Moving %s to different database and schema", leafLevelName))
	default:
		d.Partial(true)
		return diag.FromErr(fmt.Errorf(
			"unknown rename use case: old database %s (exists: %t), new database %s (exists: %t), schema %s (exists: %t), schema %s (exists: %t), schema %s (exists: %t), schema %s (exists: %t). See https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/object_renaming_guide",
			oldDatabaseId.FullyQualifiedName(), oldDatabaseExists,
			newDatabaseId.FullyQualifiedName(), newDatabaseExists,
			schemaInNewDbOldName.FullyQualifiedName(), schemaNewDbOldExists,
			schemaInNewDbNewName.FullyQualifiedName(), schemaNewDbNewExists,
			schemaInOldDbOldName.FullyQualifiedName(), schemaOldDbOldExists,
			schemaInOldDbNewName.FullyQualifiedName(), schemaOldDbNewExists,
		))
	}
}

// handleTwoLevelHierarchyRename handles the full 2-level hierarchy rename/move
// for a database-level object (database → object).
// OR is the return type of objectShowByIDFn.
func handleTwoLevelHierarchyRename[OR any](
	ctx context.Context,
	d *schema.ResourceData,
	client *sdk.Client,
	id *sdk.DatabaseObjectIdentifier,
	objectRenameFn func(sdk.DatabaseObjectIdentifier, sdk.DatabaseObjectIdentifier) func() error,
	objectShowByIDFn func(context.Context, sdk.DatabaseObjectIdentifier) (OR, error),
	encodeIdFn func(sdk.DatabaseObjectIdentifier) string,
	leafLevelName string,
) diag.Diagnostics {
	oldDatabaseNameRaw, newDatabaseNameRaw := d.GetChange("database")
	oldDatabaseName, newDatabaseName := oldDatabaseNameRaw.(string), newDatabaseNameRaw.(string)
	objectName := id.Name()

	oldDatabaseId := sdk.NewAccountObjectIdentifier(oldDatabaseName)
	newDatabaseId := sdk.NewAccountObjectIdentifier(newDatabaseName)
	oldObjectId := sdk.NewDatabaseObjectIdentifier(oldDatabaseName, objectName)
	newObjectId := sdk.NewDatabaseObjectIdentifier(newDatabaseName, objectName)

	renameObj := func(description string) {
		handleHierarchyRenameIdUpdate(d,
			func() string { return encodeIdFn(newObjectId) },
			description)
	}

	moveObj := func(currentId sdk.DatabaseObjectIdentifier, description string) diag.Diagnostics {
		return handleHierarchyMove(d,
			func() string { return encodeIdFn(newObjectId) },
			currentId, newObjectId,
			objectRenameFn,
			description)
	}

	if diags := handleShallowHierarchyRename(ctx, d, client.Databases.ShowByID, objectShowByIDFn, newDatabaseId, oldDatabaseId, newObjectId, oldObjectId, *id, renameObj, moveObj, "database", leafLevelName, leafLevelName); diags != nil {
		return diags
	}

	*id = newObjectId
	return nil
}

// handleThreeLevelHierarchyRename handles the full 3-level hierarchy rename/move
// for a schema-level object (database → schema → object).
// OR is the return type of objectShowByIDFn.
func handleThreeLevelHierarchyRename[OR any](
	ctx context.Context,
	d *schema.ResourceData,
	client *sdk.Client,
	id *sdk.SchemaObjectIdentifier,
	objectRenameFn func(sdk.SchemaObjectIdentifier, sdk.SchemaObjectIdentifier) func() error,
	objectShowByIDFn func(context.Context, sdk.SchemaObjectIdentifier) (OR, error),
	encodeIdFn func(sdk.SchemaObjectIdentifier) string,
	leafLevelName string,
) diag.Diagnostics {
	oldDatabaseNameRaw, newDatabaseNameRaw := d.GetChange("database")
	oldDatabaseName, newDatabaseName := oldDatabaseNameRaw.(string), newDatabaseNameRaw.(string)
	oldSchemaNameRaw, newSchemaNameRaw := d.GetChange("schema")
	oldSchemaName, newSchemaName := oldSchemaNameRaw.(string), newSchemaNameRaw.(string)
	objectName := id.Name()

	databaseChanged := d.HasChange("database")
	schemaChanged := d.HasChange("schema")

	oldDatabaseId := sdk.NewAccountObjectIdentifier(oldDatabaseName)
	newDatabaseId := sdk.NewAccountObjectIdentifier(newDatabaseName)
	oldSchemaId := sdk.NewDatabaseObjectIdentifierInDatabase(oldDatabaseId, oldSchemaName)
	newSchemaId := sdk.NewDatabaseObjectIdentifierInDatabase(newDatabaseId, newSchemaName)
	oldObjectId := sdk.NewSchemaObjectIdentifierInSchema(oldSchemaId, objectName)
	newObjectId := sdk.NewSchemaObjectIdentifierInSchema(newSchemaId, objectName)

	renameObj := func(description string) {
		handleHierarchyRenameIdUpdate(d,
			func() string { return encodeIdFn(newObjectId) },
			description)
	}

	moveObj := func(currentId sdk.SchemaObjectIdentifier, description string) diag.Diagnostics {
		return handleHierarchyMove(d,
			func() string { return encodeIdFn(newObjectId) },
			currentId, newObjectId,
			objectRenameFn,
			description)
	}

	switch {
	case databaseChanged && !schemaChanged:
		if diags := handleShallowHierarchyRename(ctx, d, client.Databases.ShowByID, client.Schemas.ShowByID, newDatabaseId, oldDatabaseId, newSchemaId, oldSchemaId, *id, renameObj, moveObj, "database", "schema", leafLevelName); diags != nil {
			return diags
		}

	case !databaseChanged && schemaChanged:
		if diags := handleShallowHierarchyRename(ctx, d, client.Schemas.ShowByID, objectShowByIDFn, newSchemaId, oldSchemaId, newObjectId, oldObjectId, *id, renameObj, moveObj, "schema", leafLevelName, leafLevelName); diags != nil {
			return diags
		}

	default:
		if diags := handleDeepHierarchyRename(ctx, d, client, newDatabaseId, oldDatabaseId, oldSchemaName, newSchemaName, objectName, renameObj, moveObj, leafLevelName); diags != nil {
			return diags
		}
	}

	*id = newObjectId
	return nil
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
