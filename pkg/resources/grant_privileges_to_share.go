package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var grantPrivilegesToShareGrantExactlyOneOfValidation = []string{
	"on_database",
	"on_schema",
	"on_function",
	"on_table",
	"on_all_tables_in_schema",
	"on_tag",
	"on_view",
}

var grantPrivilegesToShareSchema = map[string]*schema.Schema{
	"to_share": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      relatedResourceDescription("The fully qualified name of the share on which privileges will be granted.", resources.Share),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"privileges": {
		Type:        schema.TypeSet,
		Required:    true,
		Description: "The privileges to grant on the share. See available list of privileges: https://docs.snowflake.com/en/sql-reference/sql/grant-privilege-share#syntax",
		Elem:        &schema.Schema{Type: schema.TypeString},
		MinItems:    1,
	},
	"on_database": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      relatedResourceDescription("The fully qualified name of the database on which privileges will be granted.", resources.Database),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"on_schema": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      relatedResourceDescription("The fully qualified name of the schema on which privileges will be granted.", resources.Schema),
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"on_table": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      relatedResourceDescription("The fully qualified name of the table on which privileges will be granted.", resources.Table),
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"on_all_tables_in_schema": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified identifier for the schema for which the specified privilege will be granted for all tables.",
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"on_tag": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      relatedResourceDescription("The fully qualified name of the tag on which privileges will be granted.", resources.Tag),
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"on_view": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      relatedResourceDescription("The fully qualified name of the view on which privileges will be granted.", resources.View),
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
	},
	"on_function": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "The fully qualified name of the function on which privileges will be granted.",
		ExactlyOneOf:     grantPrivilegesToShareGrantExactlyOneOfValidation,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

func GrantPrivilegesToShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.GrantPrivilegesToShare, CreateGrantPrivilegesToShare),
		UpdateContext: TrackingUpdateWrapper(resources.GrantPrivilegesToShare, UpdateGrantPrivilegesToShare),
		DeleteContext: TrackingDeleteWrapper(resources.GrantPrivilegesToShare, DeleteGrantPrivilegesToShare),
		ReadContext:   TrackingReadWrapper(resources.GrantPrivilegesToShare, ReadGrantPrivilegesToShare),

		Schema: grantPrivilegesToShareSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.GrantPrivilegesToShare, ImportGrantPrivilegesToShare()),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportGrantPrivilegesToShare() func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		id, err := ParseGrantPrivilegesToShareId(d.Id())
		if err != nil {
			return nil, err
		}
		if err := d.Set("to_share", id.ShareName.Name()); err != nil {
			return nil, err
		}
		if err := d.Set("privileges", id.Privileges); err != nil {
			return nil, err
		}

		switch id.Kind {
		case OnDatabaseShareGrantKind:
			if err := d.Set("on_database", id.Identifier.Name()); err != nil {
				return nil, err
			}
		case OnSchemaShareGrantKind:
			if err := d.Set("on_schema", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnFunctionShareGrantKind:
			if err := d.Set("on_function", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnTableShareGrantKind:
			if err := d.Set("on_table", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnAllTablesInSchemaShareGrantKind:
			if err := d.Set("on_all_tables_in_schema", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnTagShareGrantKind:
			if err := d.Set("on_tag", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		case OnViewShareGrantKind:
			if err := d.Set("on_view", id.Identifier.FullyQualifiedName()); err != nil {
				return nil, err
			}
		}

		return []*schema.ResourceData{d}, nil
	}
}

func CreateGrantPrivilegesToShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := createGrantPrivilegesToShareIdFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] created identifier from schema: %s", id.String())

	grantOn, err := getShareGrantOn(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Grants.GrantPrivilegeToShare(ctx, getObjectPrivilegesFromSchema(d), grantOn, id.ShareName)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred when granting privileges to share",
				Detail:   fmt.Sprintf("Id: %s\nShare name: %s\nError: %s", id.String(), id.ShareName, err.Error()),
			},
		}
	}

	d.SetId(id.String())

	return ReadGrantPrivilegesToShare(ctx, d, meta)
}

func UpdateGrantPrivilegesToShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := ParseGrantPrivilegesToShareId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
			},
		}
	}

	if d.HasChange("privileges") {
		oldPrivileges, newPrivileges := d.GetChange("privileges")
		privilegesBeforeChange := expandStringList(oldPrivileges.(*schema.Set).List())
		privilegesAfterChange := expandStringList(newPrivileges.(*schema.Set).List())

		var privilegesToAdd, privilegesToRemove []sdk.ObjectPrivilege

		for _, privilegeBeforeChange := range privilegesBeforeChange {
			if !slices.Contains(privilegesAfterChange, privilegeBeforeChange) {
				privilegesToRemove = append(privilegesToRemove, sdk.ObjectPrivilege(privilegeBeforeChange))
			}
		}

		for _, privilegeAfterChange := range privilegesAfterChange {
			if !slices.Contains(privilegesBeforeChange, privilegeAfterChange) {
				privilegesToAdd = append(privilegesToAdd, sdk.ObjectPrivilege(privilegeAfterChange))
			}
		}

		grantOn, err := getShareGrantOn(d)
		if err != nil {
			return diag.FromErr(err)
		}

		if len(privilegesToAdd) > 0 {
			err = client.Grants.GrantPrivilegeToShare(
				ctx,
				privilegesToAdd,
				grantOn,
				id.ShareName,
			)
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
			err = client.Grants.RevokePrivilegeFromShare(
				ctx,
				privilegesToRemove,
				grantOn,
				id.ShareName,
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
		d.SetId(id.String())
	}

	return ReadGrantPrivilegesToShare(ctx, d, meta)
}

func DeleteGrantPrivilegesToShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := ParseGrantPrivilegesToShareId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
			},
		}
	}

	grantOn, err := getShareGrantOn(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Grants.RevokePrivilegeFromShare(ctx, getObjectPrivilegesFromSchema(d), grantOn, id.ShareName)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred when revoking privileges from share",
				Detail:   fmt.Sprintf("Id: %s\nShare name: %s\nError: %s", d.Id(), id.ShareName.FullyQualifiedName(), err.Error()),
			},
		}
	}

	d.SetId("")

	return nil
}

func ReadGrantPrivilegesToShare(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := ParseGrantPrivilegesToShareId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err.Error()),
			},
		}
	}

	opts, grantedOn := prepareShowGrantsRequestForShare(id)
	if opts == nil {
		return nil
	}

	client := meta.(*provider.Context).Client
	if _, err := client.Shares.ShowByID(ctx, id.ShareName); err != nil && errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to retrieve share. Marking the resource as removed.",
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
					Summary:  "Failed to retrieve grants. Object not found. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
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
		if grant.GrantedTo != sdk.ObjectTypeShare {
			continue
		}
		// Only consider privileges that are already present in the ID, so we
		// don't delete privileges managed by other resources.
		if !slices.Contains(id.Privileges, grant.Privilege) {
			continue
		}
		if grant.GranteeName.Name() == id.ShareName.Name() {
			if grantedOn == grant.GrantedOn {
				privileges = append(privileges, grant.Privilege)
			}
		}
	}

	// REFERENCE_USAGE is a special pseudo-privilege that you can grant or revoke,
	// but it won't show up when querying privileges (not returned by show grants ... query).
	// That's why we have to check it manually outside the loop and append it whenever it's specified in the configuration.
	if slices.Contains(id.Privileges, sdk.ObjectPrivilegeReferenceUsage.String()) {
		privileges = append(privileges, sdk.ObjectPrivilegeReferenceUsage.String())
	}

	if err := d.Set("privileges", privileges); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error setting privileges for account role",
				Detail:   fmt.Sprintf("Id: %s\nPrivileges: %v\nError: %s", d.Id(), privileges, err.Error()),
			},
		}
	}

	return nil
}

func createGrantPrivilegesToShareIdFromSchema(d *schema.ResourceData) (id *GrantPrivilegesToShareId, err error) {
	sharedId, err := sdk.ParseAccountObjectIdentifier(d.Get("to_share").(string))
	if err != nil {
		return nil, err
	}
	id = new(GrantPrivilegesToShareId)
	id.ShareName = sharedId
	id.Privileges = expandStringList(d.Get("privileges").(*schema.Set).List())

	databaseName, databaseNameOk := d.GetOk("on_database")
	schemaName, schemaNameOk := d.GetOk("on_schema")
	functionName, functionNameOk := d.GetOk("on_function")
	tableName, tableNameOk := d.GetOk("on_table")
	allTablesInSchema, allTablesInSchemaOk := d.GetOk("on_all_tables_in_schema")
	tagName, tagNameOk := d.GetOk("on_tag")
	viewName, viewNameOk := d.GetOk("on_view")

	switch {
	case databaseNameOk:
		id.Kind = OnDatabaseShareGrantKind
		databaseId, err := sdk.ParseAccountObjectIdentifier(databaseName.(string))
		if err != nil {
			return nil, err
		}
		id.Identifier = databaseId
	case schemaNameOk:
		id.Kind = OnSchemaShareGrantKind
		schemaId, err := sdk.ParseDatabaseObjectIdentifier(schemaName.(string))
		if err != nil {
			return nil, err
		}
		id.Identifier = schemaId
	case functionNameOk:
		id.Kind = OnFunctionShareGrantKind
		parsed, err := sdk.ParseSchemaObjectIdentifierWithArguments(functionName.(string))
		if err != nil {
			return nil, err
		}
		id.Identifier = parsed
	case tableNameOk:
		id.Kind = OnTableShareGrantKind
		tableId, err := sdk.ParseSchemaObjectIdentifier(tableName.(string))
		if err != nil {
			return nil, err
		}
		id.Identifier = tableId
	case allTablesInSchemaOk:
		id.Kind = OnAllTablesInSchemaShareGrantKind
		schemaId, err := sdk.ParseDatabaseObjectIdentifier(allTablesInSchema.(string))
		if err != nil {
			return nil, err
		}
		id.Identifier = schemaId
	case tagNameOk:
		id.Kind = OnTagShareGrantKind
		tagId, err := sdk.ParseSchemaObjectIdentifier(tagName.(string))
		if err != nil {
			return nil, err
		}
		id.Identifier = tagId
	case viewNameOk:
		id.Kind = OnViewShareGrantKind
		viewId, err := sdk.ParseSchemaObjectIdentifier(viewName.(string))
		if err != nil {
			return nil, err
		}
		id.Identifier = viewId
	}

	return id, nil
}

func getObjectPrivilegesFromSchema(d *schema.ResourceData) []sdk.ObjectPrivilege {
	privileges := expandStringList(d.Get("privileges").(*schema.Set).List())
	objectPrivileges := make([]sdk.ObjectPrivilege, len(privileges))
	for i, privilege := range privileges {
		objectPrivileges[i] = sdk.ObjectPrivilege(privilege)
	}
	return objectPrivileges
}

func getShareGrantOn(d *schema.ResourceData) (*sdk.ShareGrantOn, error) {
	grantOn := new(sdk.ShareGrantOn)

	databaseName, databaseNameOk := d.GetOk("on_database")
	schemaName, schemaNameOk := d.GetOk("on_schema")
	functionName, functionNameOk := d.GetOk("on_function")
	tableName, tableNameOk := d.GetOk("on_table")
	allTablesInSchema, allTablesInSchemaOk := d.GetOk("on_all_tables_in_schema")
	tagName, tagNameOk := d.GetOk("on_tag")
	viewName, viewNameOk := d.GetOk("on_view")

	switch {
	case len(databaseName.(string)) > 0 && databaseNameOk:
		databaseId, err := sdk.ParseAccountObjectIdentifier(databaseName.(string))
		if err != nil {
			return nil, err
		}
		grantOn.Database = databaseId
	case len(schemaName.(string)) > 0 && schemaNameOk:
		schemaId, err := sdk.ParseDatabaseObjectIdentifier(schemaName.(string))
		if err != nil {
			return nil, err
		}
		grantOn.Schema = schemaId
	case len(functionName.(string)) > 0 && functionNameOk:
		id, err := sdk.ParseSchemaObjectIdentifierWithArguments(functionName.(string))
		if err != nil {
			return nil, err
		}
		grantOn.Function = id
	case len(tableName.(string)) > 0 && tableNameOk:
		tableId, err := sdk.ParseSchemaObjectIdentifier(tableName.(string))
		if err != nil {
			return nil, err
		}
		grantOn.Table = &sdk.OnTable{
			Name: tableId,
		}
	case len(allTablesInSchema.(string)) > 0 && allTablesInSchemaOk:
		schemaId, err := sdk.ParseDatabaseObjectIdentifier(allTablesInSchema.(string))
		if err != nil {
			return nil, err
		}
		grantOn.Table = &sdk.OnTable{
			AllInSchema: schemaId,
		}
	case len(tagName.(string)) > 0 && tagNameOk:
		tagId, err := sdk.ParseSchemaObjectIdentifier(tagName.(string))
		if err != nil {
			return nil, err
		}
		grantOn.Tag = tagId
	case len(viewName.(string)) > 0 && viewNameOk:
		viewId, err := sdk.ParseSchemaObjectIdentifier(viewName.(string))
		if err != nil {
			return nil, err
		}
		grantOn.View = viewId
	}

	return grantOn, nil
}

func prepareShowGrantsRequestForShare(id GrantPrivilegesToShareId) (*sdk.ShowGrantOptions, sdk.ObjectType) {
	opts := new(sdk.ShowGrantOptions)
	var objectType sdk.ObjectType

	switch id.Kind {
	case OnDatabaseShareGrantKind:
		objectType = sdk.ObjectTypeDatabase
	case OnSchemaShareGrantKind:
		objectType = sdk.ObjectTypeSchema
	case OnTableShareGrantKind:
		objectType = sdk.ObjectTypeTable
	case OnAllTablesInSchemaShareGrantKind:
		log.Printf("[INFO] Show with on_all_tables_in_schema option is skipped. No changes in privileges in Snowflake will be detected.")
		return nil, ""
	case OnTagShareGrantKind:
		objectType = sdk.ObjectTypeTag
	case OnViewShareGrantKind:
		objectType = sdk.ObjectTypeView
	case OnFunctionShareGrantKind:
		objectType = sdk.ObjectTypeFunction
	}

	opts.On = &sdk.ShowGrantsOn{
		Object: &sdk.Object{
			ObjectType: objectType,
			Name:       id.Identifier,
		},
	}

	return opts, objectType
}
