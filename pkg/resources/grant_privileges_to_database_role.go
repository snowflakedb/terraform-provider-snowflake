package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO: Handle IMPORTED PRIVILEGES privilege (after second account will be added - SNOW-976501)

var grantPrivilegesToDatabaseRoleSchema = map[string]*schema.Schema{
	"database_role_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      relatedResourceDescription("The fully qualified name of the database role to which privileges will be granted.", resources.DatabaseRole),
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"privileges": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "The privileges to grant on the database role.",
		MinItems:    1,
		ExactlyOneOf: []string{
			"privileges",
			"all_privileges",
		},
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: isNotOwnershipGrant(),
		},
	},
	"all_privileges": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Grant all privileges on the database role.",
		ExactlyOneOf: []string{
			"privileges",
			"all_privileges",
		},
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		ForceNew:    true,
		Description: "If specified, allows the recipient role to grant the privileges to other roles.",
	},
	"always_apply": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If true, the resource will always produce a “plan” and on “apply” it will re-grant defined privileges. It is supposed to be used only in “grant privileges on all X’s in database / schema Y” or “grant all privileges to X” scenarios to make sure that every new object in a given database / schema is granted by the account role and every new privilege is granted to the database role. Important note: this flag is not compliant with the Terraform assumptions of the config being eventually convergent (producing an empty plan).",
	},
	"always_apply_trigger": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "This is a helper field and should not be set. Its main purpose is to help to achieve the functionality described by the always_apply field.",
	},
	"on_database": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      relatedResourceDescription("The fully qualified name of the database on which privileges will be granted.", resources.Database),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf: []string{
			"on_database",
			"on_schema",
			"on_schema_object",
		},
	},
	"on_schema": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the schema on which privileges will be granted.",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"on_database",
			"on_schema",
			"on_schema_object",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"schema_name": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Description:      "The fully qualified name of the schema.",
					ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
					DiffSuppressFunc: suppressIdentifierQuoting,
					ExactlyOneOf: []string{
						"on_schema.0.schema_name",
						"on_schema.0.all_schemas_in_database",
						"on_schema.0.future_schemas_in_database",
					},
				},
				"all_schemas_in_database": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Description:      "The fully qualified name of the database.",
					ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
					DiffSuppressFunc: suppressIdentifierQuoting,
					ExactlyOneOf: []string{
						"on_schema.0.schema_name",
						"on_schema.0.all_schemas_in_database",
						"on_schema.0.future_schemas_in_database",
					},
				},
				"future_schemas_in_database": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Description:      "The fully qualified name of the database.",
					ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
					DiffSuppressFunc: suppressIdentifierQuoting,
					ExactlyOneOf: []string{
						"on_schema.0.schema_name",
						"on_schema.0.all_schemas_in_database",
						"on_schema.0.future_schemas_in_database",
					},
				},
			},
		},
	},
	"on_schema_object": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the schema object on which privileges will be granted.",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"on_database",
			"on_schema",
			"on_schema_object",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"object_type": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: fmt.Sprintf("The object type of the schema object on which privileges will be granted. Valid values are: %s", strings.Join(sdk.ValidGrantToSchemaObjectTypesString, " | ")),
					RequiredWith: []string{
						"on_schema_object.0.object_name",
					},
					ConflictsWith: []string{
						"on_schema_object.0.all",
						"on_schema_object.0.future",
					},
					ValidateDiagFunc: StringInSlice(sdk.ValidGrantToSchemaObjectTypesString, true),
				},
				"object_name": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "The fully qualified name of the object on which privileges will be granted.",
					RequiredWith: []string{
						"on_schema_object.0.object_type",
					},
					ExactlyOneOf: []string{
						"on_schema_object.0.object_name",
						"on_schema_object.0.all",
						"on_schema_object.0.future",
					},
					DiffSuppressFunc: suppressIdentifierQuoting,
				},
				"all": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					Description: "Configures the privilege to be granted on all objects in either a database or schema.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: getGrantPrivilegesOnDatabaseRoleBulkOperationSchema(sdk.ValidGrantToPluralObjectTypesString),
					},
					ConflictsWith: []string{
						"on_schema_object.0.object_type",
					},
					ExactlyOneOf: []string{
						"on_schema_object.0.object_name",
						"on_schema_object.0.all",
						"on_schema_object.0.future",
					},
				},
				"future": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					Description: "Configures the privilege to be granted on future objects in either a database or schema.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: getGrantPrivilegesOnDatabaseRoleBulkOperationSchema(sdk.ValidGrantToFuturePluralObjectTypesString),
					},
					ConflictsWith: []string{
						"on_schema_object.0.object_type",
					},
					ExactlyOneOf: []string{
						"on_schema_object.0.object_name",
						"on_schema_object.0.all",
						"on_schema_object.0.future",
					},
				},
			},
		},
	},
}

func getGrantPrivilegesOnDatabaseRoleBulkOperationSchema(validGrantToObjectTypes []string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"object_type_plural": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			Description:      fmt.Sprintf("The plural object type of the schema object on which privileges will be granted. Valid values are: %s.", strings.Join(validGrantToObjectTypes, " | ")),
			ValidateDiagFunc: StringInSlice(validGrantToObjectTypes, true),
		},
		"in_database": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			Description:      "The fully qualified name of the database.",
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"in_schema": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			Description:      "The fully qualified name of the schema.",
			ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
	}
}

func GrantPrivilegesToDatabaseRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.GrantPrivilegesToDatabaseRole, CreateGrantPrivilegesToDatabaseRole),
		UpdateContext: TrackingUpdateWrapper(resources.GrantPrivilegesToDatabaseRole, UpdateGrantPrivilegesToDatabaseRole),
		DeleteContext: TrackingDeleteWrapper(resources.GrantPrivilegesToDatabaseRole, DeleteGrantPrivilegesToDatabaseRole),
		ReadContext:   TrackingReadWrapper(resources.GrantPrivilegesToDatabaseRole, ReadGrantPrivilegesToDatabaseRole),

		Schema: grantPrivilegesToDatabaseRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.GrantPrivilegesToDatabaseRole, ImportGrantPrivilegesToDatabaseRole),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportGrantPrivilegesToDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	id, err := ParseGrantPrivilegesToDatabaseRoleId(d.Id())
	if err != nil {
		return nil, err
	}
	if err := d.Set("database_role_name", id.DatabaseRoleName.FullyQualifiedName()); err != nil {
		return nil, err
	}
	if err := d.Set("with_grant_option", id.WithGrantOption); err != nil {
		return nil, err
	}
	if err := d.Set("always_apply", id.AlwaysApply); err != nil {
		return nil, err
	}
	if err := d.Set("all_privileges", id.AllPrivileges); err != nil {
		return nil, err
	}
	if err := d.Set("privileges", id.Privileges); err != nil {
		return nil, err
	}

	switch id.Kind {
	case OnDatabaseDatabaseRoleGrantKind:
		if err := d.Set("on_database", id.Data.(*OnDatabaseGrantData).DatabaseName.FullyQualifiedName()); err != nil {
			return nil, err
		}
	case OnSchemaDatabaseRoleGrantKind:
		data := id.Data.(*OnSchemaGrantData)
		onSchema := make(map[string]any)

		switch data.Kind {
		case OnSchemaSchemaGrantKind:
			onSchema["schema_name"] = data.SchemaName.FullyQualifiedName()
		case OnAllSchemasInDatabaseSchemaGrantKind:
			onSchema["all_schemas_in_database"] = data.DatabaseName.FullyQualifiedName()
		case OnFutureSchemasInDatabaseSchemaGrantKind:
			onSchema["future_schemas_in_database"] = data.DatabaseName.FullyQualifiedName()
		}

		if err := d.Set("on_schema", []any{onSchema}); err != nil {
			return nil, err
		}
	case OnSchemaObjectDatabaseRoleGrantKind:
		data := id.Data.(*OnSchemaObjectGrantData)
		onSchemaObject := make(map[string]any)

		switch data.Kind {
		case OnObjectSchemaObjectGrantKind:
			onSchemaObject["object_type"] = data.Object.ObjectType.String()
			onSchemaObject["object_name"] = data.Object.Name.FullyQualifiedName()
		case OnAllSchemaObjectGrantKind:
			onAll := make(map[string]any)

			onAll["object_type_plural"] = data.OnAllOrFuture.ObjectNamePlural.String()
			switch data.OnAllOrFuture.Kind {
			case InDatabaseBulkOperationGrantKind:
				onAll["in_database"] = data.OnAllOrFuture.Database.FullyQualifiedName()
			case InSchemaBulkOperationGrantKind:
				onAll["in_schema"] = data.OnAllOrFuture.Schema.FullyQualifiedName()
			}

			onSchemaObject["all"] = []any{onAll}
		case OnFutureSchemaObjectGrantKind:
			onFuture := make(map[string]any)

			onFuture["object_type_plural"] = data.OnAllOrFuture.ObjectNamePlural.String()
			switch data.OnAllOrFuture.Kind {
			case InDatabaseBulkOperationGrantKind:
				onFuture["in_database"] = data.OnAllOrFuture.Database.FullyQualifiedName()
			case InSchemaBulkOperationGrantKind:
				onFuture["in_schema"] = data.OnAllOrFuture.Schema.FullyQualifiedName()
			}

			onSchemaObject["future"] = []any{onFuture}
		}

		if err := d.Set("on_schema_object", []any{onSchemaObject}); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func CreateGrantPrivilegesToDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := createGrantPrivilegesToDatabaseRoleIdFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	grantOn, err := getDatabaseRoleGrantOn(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Grants.GrantPrivilegesToDatabaseRole(
		ctx,
		getDatabaseRolePrivilegesFromSchema(d),
		grantOn,
		sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(d.Get("database_role_name").(string)),
		&sdk.GrantPrivilegesToDatabaseRoleOptions{
			WithGrantOption: sdk.Bool(d.Get("with_grant_option").(bool)),
		},
	)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred when granting privileges to database role",
				Detail:   fmt.Sprintf("Id: %s\nDatabase role name: %s\nError: %s", id.String(), id.DatabaseRoleName, err.Error()),
			},
		}
	}

	d.SetId(id.String())

	return ReadGrantPrivilegesToDatabaseRole(ctx, d, meta)
}

func UpdateGrantPrivilegesToDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := ParseGrantPrivilegesToDatabaseRoleId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
			},
		}
	}

	if d.HasChange("with_grant_option") {
		id.WithGrantOption = d.Get("with_grant_option").(bool)
	}

	// handle all_privileges -> privileges change (revoke all privileges)
	if d.HasChange("all_privileges") {
		_, allPrivileges := d.GetChange("all_privileges")

		if !allPrivileges.(bool) {
			grantOn, err := getDatabaseRoleGrantOn(d)
			if err != nil {
				return diag.FromErr(err)
			}
			err = client.Grants.RevokePrivilegesFromDatabaseRole(ctx, &sdk.DatabaseRoleGrantPrivileges{
				AllPrivileges: sdk.Bool(true),
			},
				grantOn,
				id.DatabaseRoleName,
				new(sdk.RevokePrivilegesFromDatabaseRoleOptions),
			)
			if err != nil {
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to revoke all privileges",
						Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
					},
				}
			}
		}

		id.AllPrivileges = allPrivileges.(bool)
	}

	if d.HasChange("privileges") {
		shouldHandlePrivilegesChange := true

		// Skip if all_privileges was set to true
		if d.HasChange("all_privileges") {
			if _, allPrivileges := d.GetChange("all_privileges"); allPrivileges.(bool) {
				shouldHandlePrivilegesChange = false
				id.Privileges = []string{}
			}
		}

		if shouldHandlePrivilegesChange {
			before, after := d.GetChange("privileges")
			privilegesBeforeChange := expandStringList(before.(*schema.Set).List())
			privilegesAfterChange := expandStringList(after.(*schema.Set).List())

			var privilegesToAdd, privilegesToRemove []string

			for _, privilegeBeforeChange := range privilegesBeforeChange {
				if !slices.Contains(privilegesAfterChange, privilegeBeforeChange) {
					privilegesToRemove = append(privilegesToRemove, privilegeBeforeChange)
				}
			}

			for _, privilegeAfterChange := range privilegesAfterChange {
				if !slices.Contains(privilegesBeforeChange, privilegeAfterChange) {
					privilegesToAdd = append(privilegesToAdd, privilegeAfterChange)
				}
			}

			grantOn, err := getDatabaseRoleGrantOn(d)
			if err != nil {
				return diag.FromErr(err)
			}

			if len(privilegesToAdd) > 0 {
				privilegesToGrant := getDatabaseRolePrivileges(
					false,
					privilegesToAdd,
					id.Kind == OnDatabaseDatabaseRoleGrantKind,
					id.Kind == OnSchemaDatabaseRoleGrantKind,
					id.Kind == OnSchemaObjectDatabaseRoleGrantKind,
				)

				if !id.WithGrantOption {
					if err = client.Grants.RevokePrivilegesFromDatabaseRole(ctx, privilegesToGrant, grantOn, id.DatabaseRoleName, &sdk.RevokePrivilegesFromDatabaseRoleOptions{
						GrantOptionFor: sdk.Bool(true),
					}); err != nil {
						return diag.Diagnostics{
							diag.Diagnostic{
								Severity: diag.Error,
								Summary:  "Failed to revoke privileges to add",
								Detail:   fmt.Sprintf("Id: %s\nPrivileges to add: %v\nError: %s", d.Id(), privilegesToAdd, err.Error()),
							},
						}
					}
				}

				err = client.Grants.GrantPrivilegesToDatabaseRole(ctx, privilegesToGrant, grantOn, id.DatabaseRoleName, &sdk.GrantPrivilegesToDatabaseRoleOptions{WithGrantOption: sdk.Bool(id.WithGrantOption)})
				if err != nil {
					return diag.Diagnostics{
						diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Failed to grant added privileges",
							Detail:   fmt.Sprintf("Id: %s\nPrivileges to add: %v\nError: %s", d.Id(), privilegesToAdd, err.Error()),
						},
					}
				}
			}

			if len(privilegesToRemove) > 0 {
				err = client.Grants.RevokePrivilegesFromDatabaseRole(
					ctx,
					getDatabaseRolePrivileges(
						false,
						privilegesToRemove,
						id.Kind == OnDatabaseDatabaseRoleGrantKind,
						id.Kind == OnSchemaDatabaseRoleGrantKind,
						id.Kind == OnSchemaObjectDatabaseRoleGrantKind,
					),
					grantOn,
					id.DatabaseRoleName,
					new(sdk.RevokePrivilegesFromDatabaseRoleOptions),
				)
				if err != nil {
					return diag.Diagnostics{
						diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Failed to revoke removed privileges",
							Detail:   fmt.Sprintf("Id: %s\nPrivileges to remove: %v\nError: %s", d.Id(), privilegesToRemove, err.Error()),
						},
					}
				}
			}

			id.Privileges = privilegesAfterChange
		}
	}

	// handle privileges -> all_privileges change (grant all privileges)
	if d.HasChange("all_privileges") {
		_, allPrivileges := d.GetChange("all_privileges")

		if allPrivileges.(bool) {
			grantOn, err := getDatabaseRoleGrantOn(d)
			if err != nil {
				return diag.FromErr(err)
			}
			err = client.Grants.GrantPrivilegesToDatabaseRole(ctx, &sdk.DatabaseRoleGrantPrivileges{
				AllPrivileges: sdk.Bool(true),
			},
				grantOn,
				id.DatabaseRoleName,
				new(sdk.GrantPrivilegesToDatabaseRoleOptions),
			)
			if err != nil {
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to grant all privileges",
						Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
					},
				}
			}
		}

		id.AllPrivileges = allPrivileges.(bool)
	}

	if d.HasChange("always_apply") {
		id.AlwaysApply = d.Get("always_apply").(bool)
	}

	if id.AlwaysApply {
		grantOn, err := getDatabaseRoleGrantOn(d)
		if err != nil {
			return diag.FromErr(err)
		}
		err = client.Grants.GrantPrivilegesToDatabaseRole(
			ctx,
			getDatabaseRolePrivilegesFromSchema(d),
			grantOn,
			id.DatabaseRoleName,
			&sdk.GrantPrivilegesToDatabaseRoleOptions{
				WithGrantOption: &id.WithGrantOption,
			},
		)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Always apply. An error occurred when granting privileges to database role",
					Detail:   fmt.Sprintf("Id: %s\nDatabase role name: %s\nError: %s", d.Id(), id.DatabaseRoleName, err.Error()),
				},
			}
		}
	}

	d.SetId(id.String())

	return ReadGrantPrivilegesToDatabaseRole(ctx, d, meta)
}

func DeleteGrantPrivilegesToDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := ParseGrantPrivilegesToDatabaseRoleId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
			},
		}
	}

	grantOn, err := getDatabaseRoleGrantOn(d)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.Grants.RevokePrivilegesFromDatabaseRole(
		ctx,
		getDatabaseRolePrivilegesFromSchema(d),
		grantOn,
		id.DatabaseRoleName,
		&sdk.RevokePrivilegesFromDatabaseRoleOptions{},
	)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred when revoking privileges from database role",
				Detail:   fmt.Sprintf("Id: %s\nDatabase role name: %s\nError: %s", d.Id(), id.DatabaseRoleName, err.Error()),
			},
		}
	}

	d.SetId("")

	return nil
}

func ReadGrantPrivilegesToDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := ParseGrantPrivilegesToDatabaseRoleId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
			},
		}
	}

	if id.AlwaysApply {
		// The Trigger is a string rather than boolean that would be flipped on every terraform apply
		// because it's easier to think about and not to worry about edge cases that may occur with 1bit values.
		// The only place to have the "flip" is Read operation, because there we can set value and produce a plan
		// that later on will be executed in the Update operation.
		//
		// The following example shows that we can end up with the same value as before, which may lead to empty plans:
		// 1. Create configuration with always_apply = false (let's say trigger will be false by default)
		// 2. terraform apply: Create (Read will update it to false)
		// 3. Update config so that always_apply = true
		// 4. terraform apply: Read (updated trigger to false) -> change is not detected (no plan; no Update)
		triggerId, err := uuid.GenerateUUID()
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to generate UUID",
					Detail:   fmt.Sprintf("Original error: %s", err.Error()),
				},
			}
		}

		// Change the value of always_apply_trigger to produce a plan
		if err := d.Set("always_apply_trigger", triggerId); err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Error setting always_apply_trigger for database role",
					Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
				},
			}
		}
	}

	if id.AllPrivileges {
		log.Printf("[INFO] Show with all_privileges option is skipped. No changes in privileges in Snowflake will be detected. Consider specifying all privileges in 'privileges' block.")
		return nil
	}

	opts, grantedOn := prepareShowGrantsRequest(id)
	if opts == nil {
		return nil
	}

	client := meta.(*provider.Context).Client
	if _, err := client.DatabaseRoles.ShowByID(ctx, id.DatabaseRoleName); err != nil && errors.Is(err, sdk.ErrObjectNotFound) {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to retrieve database role. Marking the resource as removed.",
				Detail:   fmt.Sprintf("Id: %s", d.Id()),
			},
		}
	}

	grants, err := client.Grants.Show(ctx, opts)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve grants. Target object not found. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Id: %s", d.Id()),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to retrieve grants",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
			},
		}
	}

	var privileges []string

	for _, grant := range grants {
		// Accept only DATABASE ROLEs
		if grant.GrantTo != sdk.ObjectTypeDatabaseRole && grant.GrantedTo != sdk.ObjectTypeDatabaseRole {
			continue
		}
		// Only consider privileges that are already present in the ID, so we
		// don't delete privileges managed by other resources.
		if !slices.Contains(id.Privileges, grant.Privilege) {
			continue
		}
		if id.WithGrantOption == grant.GrantOption && id.DatabaseRoleName.Name() == grant.GranteeName.Name() {
			// Future grants do not have grantedBy, only current grants do.
			// If grantedby is an empty string, it means terraform could not have created the grant
			if (opts.Future == nil || !*opts.Future) && grant.GrantedBy.Name() == "" {
				continue
			}
			// grant_on is for future grants, granted_on is for current grants.
			// They function the same way though in a test for matching the object type
			if grantedOn == grant.GrantedOn || grantedOn == grant.GrantOn {
				privileges = append(privileges, grant.Privilege)
			}
		}
	}

	if err := d.Set("privileges", privileges); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error setting privileges for database role",
				Detail:   fmt.Sprintf("Id: %s\nPrivileges: %v\nError: %s", d.Id(), privileges, err.Error()),
			},
		}
	}

	return nil
}

func prepareShowGrantsRequest(id GrantPrivilegesToDatabaseRoleId) (*sdk.ShowGrantOptions, sdk.ObjectType) {
	opts := new(sdk.ShowGrantOptions)
	var grantedOn sdk.ObjectType

	switch id.Kind {
	case OnDatabaseDatabaseRoleGrantKind:
		grantedOn = sdk.ObjectTypeDatabase
		data := id.Data.(*OnDatabaseGrantData)
		opts.On = &sdk.ShowGrantsOn{
			Object: &sdk.Object{
				ObjectType: sdk.ObjectTypeDatabase,
				Name:       data.DatabaseName,
			},
		}
	case OnSchemaDatabaseRoleGrantKind:
		grantedOn = sdk.ObjectTypeSchema
		data := id.Data.(*OnSchemaGrantData)

		switch data.Kind {
		case OnSchemaSchemaGrantKind:
			opts.On = &sdk.ShowGrantsOn{
				Object: &sdk.Object{
					ObjectType: sdk.ObjectTypeSchema,
					Name:       data.SchemaName,
				},
			}
		case OnAllSchemasInDatabaseSchemaGrantKind:
			log.Printf("[INFO] Show with on_schema.all_schemas_in_database option is skipped. No changes in privileges in Snowflake will be detected.")
			return nil, ""
		case OnFutureSchemasInDatabaseSchemaGrantKind:
			opts.Future = sdk.Bool(true)
			opts.In = &sdk.ShowGrantsIn{
				Database: data.DatabaseName,
			}
		}
	case OnSchemaObjectDatabaseRoleGrantKind:
		data := id.Data.(*OnSchemaObjectGrantData)

		switch data.Kind {
		case OnObjectSchemaObjectGrantKind:
			grantedOn = data.Object.ObjectType
			opts.On = &sdk.ShowGrantsOn{
				Object: data.Object,
			}
		case OnAllSchemaObjectGrantKind:
			log.Printf("[INFO] Show with on_schema_object.on_all option is skipped. No changes in privileges in Snowflake will be detected.")
			return nil, ""
		case OnFutureSchemaObjectGrantKind:
			grantedOn = data.OnAllOrFuture.ObjectNamePlural.Singular()
			opts.Future = sdk.Bool(true)

			switch data.OnAllOrFuture.Kind {
			case InDatabaseBulkOperationGrantKind:
				opts.In = &sdk.ShowGrantsIn{
					Database: data.OnAllOrFuture.Database,
				}
			case InSchemaBulkOperationGrantKind:
				opts.In = &sdk.ShowGrantsIn{
					Schema: data.OnAllOrFuture.Schema,
				}
			}
		}
	}

	return opts, grantedOn
}

func getDatabaseRolePrivilegesFromSchema(d *schema.ResourceData) *sdk.DatabaseRoleGrantPrivileges {
	_, onDatabaseOk := d.GetOk("on_database")
	_, onSchemaOk := d.GetOk("on_schema")
	_, onSchemaObjectOk := d.GetOk("on_schema_object")

	return getDatabaseRolePrivileges(
		d.Get("all_privileges").(bool),
		expandStringList(d.Get("privileges").(*schema.Set).List()),
		onDatabaseOk,
		onSchemaOk,
		onSchemaObjectOk,
	)
}

func getDatabaseRolePrivileges(allPrivileges bool, privileges []string, onDatabase bool, onSchema bool, onSchemaObject bool) *sdk.DatabaseRoleGrantPrivileges {
	databaseRoleGrantPrivileges := new(sdk.DatabaseRoleGrantPrivileges)

	if allPrivileges {
		databaseRoleGrantPrivileges.AllPrivileges = sdk.Bool(true)
		return databaseRoleGrantPrivileges
	}

	switch {
	case onDatabase:
		databasePrivileges := make([]sdk.AccountObjectPrivilege, len(privileges))
		for i, privilege := range privileges {
			databasePrivileges[i] = sdk.AccountObjectPrivilege(privilege)
		}
		databaseRoleGrantPrivileges.DatabasePrivileges = databasePrivileges
	case onSchema:
		schemaPrivileges := make([]sdk.SchemaPrivilege, len(privileges))
		for i, privilege := range privileges {
			schemaPrivileges[i] = sdk.SchemaPrivilege(privilege)
		}
		databaseRoleGrantPrivileges.SchemaPrivileges = schemaPrivileges
	case onSchemaObject:
		schemaObjectPrivileges := make([]sdk.SchemaObjectPrivilege, len(privileges))
		for i, privilege := range privileges {
			schemaObjectPrivileges[i] = sdk.SchemaObjectPrivilege(privilege)
		}
		databaseRoleGrantPrivileges.SchemaObjectPrivileges = schemaObjectPrivileges
	}

	return databaseRoleGrantPrivileges
}

func getDatabaseRoleGrantOn(d *schema.ResourceData) (*sdk.DatabaseRoleGrantOn, error) {
	onDatabase, onDatabaseOk := d.GetOk("on_database")
	onSchemaBlock, onSchemaOk := d.GetOk("on_schema")
	onSchemaObjectBlock, onSchemaObjectOk := d.GetOk("on_schema_object")
	on := new(sdk.DatabaseRoleGrantOn)

	switch {
	case onDatabaseOk:
		databaseId, err := sdk.ParseAccountObjectIdentifier(onDatabase.(string))
		if err != nil {
			return nil, err
		}
		on.Database = sdk.Pointer(databaseId)
	case onSchemaOk:
		onSchema := onSchemaBlock.([]any)[0].(map[string]any)

		grantOnSchema := new(sdk.GrantOnSchema)

		schemaName := onSchema["schema_name"].(string)
		schemaNameOk := len(schemaName) > 0

		allSchemasInDatabase := onSchema["all_schemas_in_database"].(string)
		allSchemasInDatabaseOk := len(allSchemasInDatabase) > 0

		futureSchemasInDatabase := onSchema["future_schemas_in_database"].(string)
		futureSchemasInDatabaseOk := len(futureSchemasInDatabase) > 0

		switch {
		case schemaNameOk:
			schemaId, err := sdk.ParseDatabaseObjectIdentifier(schemaName)
			if err != nil {
				return nil, err
			}
			grantOnSchema.Schema = sdk.Pointer(schemaId)
		case allSchemasInDatabaseOk:
			databaseId, err := sdk.ParseAccountObjectIdentifier(allSchemasInDatabase)
			if err != nil {
				return nil, err
			}
			grantOnSchema.AllSchemasInDatabase = sdk.Pointer(databaseId)
		case futureSchemasInDatabaseOk:
			databaseId, err := sdk.ParseAccountObjectIdentifier(futureSchemasInDatabase)
			if err != nil {
				return nil, err
			}
			grantOnSchema.FutureSchemasInDatabase = sdk.Pointer(databaseId)
		}

		on.Schema = grantOnSchema
	case onSchemaObjectOk:
		onSchemaObject := onSchemaObjectBlock.([]any)[0].(map[string]any)

		grantOnSchemaObject := new(sdk.GrantOnSchemaObject)

		objectType := onSchemaObject["object_type"].(string)
		objectTypeOk := len(objectType) > 0

		objectName := onSchemaObject["object_name"].(string)
		objectNameOk := len(objectName) > 0

		all := onSchemaObject["all"].([]any)
		allOk := len(all) > 0

		future := onSchemaObject["future"].([]any)
		futureOk := len(future) > 0

		switch {
		case objectTypeOk && objectNameOk:
			objectType := sdk.ObjectType(objectType)
			var id sdk.ObjectIdentifier
			var err error
			// TODO(SNOW-1569535): use a mapper from object type to parsing function
			if objectType.IsWithArguments() {
				id, err = sdk.ParseSchemaObjectIdentifierWithArguments(objectName)
				if err != nil {
					return nil, err
				}
			} else {
				id, err = sdk.ParseSchemaObjectIdentifier(objectName)
				if err != nil {
					return nil, err
				}
			}
			grantOnSchemaObject.SchemaObject = &sdk.Object{
				ObjectType: objectType,
				Name:       id,
			}
		case allOk:
			grantOnSchemaObjectIn, err := getGrantOnSchemaObjectIn(all[0].(map[string]any))
			if err != nil {
				return nil, err
			}
			grantOnSchemaObject.All = grantOnSchemaObjectIn
		case futureOk:
			grantOnSchemaObjectIn, err := getGrantOnSchemaObjectIn(future[0].(map[string]any))
			if err != nil {
				return nil, err
			}
			grantOnSchemaObject.Future = grantOnSchemaObjectIn
		}

		on.SchemaObject = grantOnSchemaObject
	}

	return on, nil
}

func createGrantPrivilegesToDatabaseRoleIdFromSchema(d *schema.ResourceData) (id *GrantPrivilegesToDatabaseRoleId, err error) {
	id = new(GrantPrivilegesToDatabaseRoleId)
	roleId, err := sdk.ParseDatabaseObjectIdentifier(d.Get("database_role_name").(string))
	if err != nil {
		return nil, err
	}
	id.DatabaseRoleName = roleId
	id.AllPrivileges = d.Get("all_privileges").(bool)
	if p, ok := d.GetOk("privileges"); ok {
		id.Privileges = expandStringList(p.(*schema.Set).List())
	}
	id.WithGrantOption = d.Get("with_grant_option").(bool)
	id.AlwaysApply = d.Get("always_apply").(bool)

	on, err := getDatabaseRoleGrantOn(d)
	if err != nil {
		return nil, err
	}
	switch {
	case on.Database != nil:
		id.Kind = OnDatabaseDatabaseRoleGrantKind
		id.Data = &OnDatabaseGrantData{
			DatabaseName: *on.Database,
		}
	case on.Schema != nil:
		onSchemaGrantData := new(OnSchemaGrantData)

		switch {
		case on.Schema.Schema != nil:
			onSchemaGrantData.Kind = OnSchemaSchemaGrantKind
			onSchemaGrantData.SchemaName = on.Schema.Schema
		case on.Schema.AllSchemasInDatabase != nil:
			onSchemaGrantData.Kind = OnAllSchemasInDatabaseSchemaGrantKind
			onSchemaGrantData.DatabaseName = on.Schema.AllSchemasInDatabase
		case on.Schema.FutureSchemasInDatabase != nil:
			onSchemaGrantData.Kind = OnFutureSchemasInDatabaseSchemaGrantKind
			onSchemaGrantData.DatabaseName = on.Schema.FutureSchemasInDatabase
		}

		id.Kind = OnSchemaDatabaseRoleGrantKind
		id.Data = onSchemaGrantData
	case on.SchemaObject != nil:
		onSchemaObjectGrantData := new(OnSchemaObjectGrantData)

		switch {
		case on.SchemaObject.SchemaObject != nil:
			onSchemaObjectGrantData.Kind = OnObjectSchemaObjectGrantKind
			onSchemaObjectGrantData.Object = on.SchemaObject.SchemaObject
		case on.SchemaObject.All != nil:
			onSchemaObjectGrantData.Kind = OnAllSchemaObjectGrantKind
			onSchemaObjectGrantData.OnAllOrFuture = getBulkOperationGrantData(on.SchemaObject.All)
		case on.SchemaObject.Future != nil:
			onSchemaObjectGrantData.Kind = OnFutureSchemaObjectGrantKind
			onSchemaObjectGrantData.OnAllOrFuture = getBulkOperationGrantData(on.SchemaObject.Future)
		}

		id.Kind = OnSchemaObjectDatabaseRoleGrantKind
		id.Data = onSchemaObjectGrantData
	}

	return id, nil
}
