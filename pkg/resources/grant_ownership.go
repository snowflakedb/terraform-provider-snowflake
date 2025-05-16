package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var grantOwnershipSchema = map[string]*schema.Schema{
	"account_role_name": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      relatedResourceDescription("The fully qualified name of the account role to which privileges will be granted.", resources.AccountRole),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf: []string{
			"account_role_name",
			"database_role_name",
		},
	},
	"database_role_name": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      relatedResourceDescription("The fully qualified name of the database role to which privileges will be granted.", resources.DatabaseRole),
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf: []string{
			"account_role_name",
			"database_role_name",
		},
	},
	"outbound_privileges": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies whether to remove or transfer all existing outbound privileges on the object when ownership is transferred to a new role. Available options are: REVOKE for removing existing privileges and COPY to transfer them with ownership. For more information head over to [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#optional-parameters).",
		ValidateFunc: validation.StringInSlice([]string{
			"COPY",
			"REVOKE",
		}, true),
	},
	"on": {
		Type:        schema.TypeList,
		Required:    true,
		ForceNew:    true,
		Description: "Configures which object(s) should transfer their ownership to the specified role.",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"object_type": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: fmt.Sprintf("Specifies the type of object on which you are transferring ownership. Available values are: %s", strings.Join(sdk.ValidGrantOwnershipObjectTypesString, " | ")),
					RequiredWith: []string{
						"on.0.object_name",
					},
					ValidateFunc: validation.StringInSlice(sdk.ValidGrantOwnershipObjectTypesString, true),
				},
				"object_name": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Description:      "Specifies the identifier for the object on which you are transferring ownership.",
					DiffSuppressFunc: suppressIdentifierQuoting,
					RequiredWith: []string{
						"on.0.object_type",
					},
					ExactlyOneOf: []string{
						"on.0.object_name",
						"on.0.all",
						"on.0.future",
					},
				},
				"all": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					Description: "Configures the privilege to be granted on all objects in either a database or schema.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: grantOwnershipBulkOperationSchema("all"),
					},
					ExactlyOneOf: []string{
						"on.0.object_name",
						"on.0.all",
						"on.0.future",
					},
				},
				"future": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					Description: "Configures the privilege to be granted on all objects in either a database or schema.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: grantOwnershipBulkOperationSchema("future"),
					},
					ExactlyOneOf: []string{
						"on.0.object_name",
						"on.0.all",
						"on.0.future",
					},
				},
			},
		},
	},
}

func grantOwnershipBulkOperationSchema(branchName string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"object_type_plural": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  fmt.Sprintf("Specifies the type of object in plural form on which you are transferring ownership. Available values are: %s. For more information head over to [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#required-parameters).", strings.Join(sdk.ValidGrantOwnershipPluralObjectTypesString, " | ")),
			ValidateFunc: validation.StringInSlice(sdk.ValidGrantOwnershipPluralObjectTypesString, true),
		},
		"in_database": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			Description:      relatedResourceDescription("The fully qualified name of the database.", resources.Database),
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
			ExactlyOneOf: []string{
				fmt.Sprintf("on.0.%s.0.in_database", branchName),
				fmt.Sprintf("on.0.%s.0.in_schema", branchName),
			},
		},
		"in_schema": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			Description:      relatedResourceDescription("The fully qualified name of the schema.", resources.Schema),
			ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
			ExactlyOneOf: []string{
				fmt.Sprintf("on.0.%s.0.in_database", branchName),
				fmt.Sprintf("on.0.%s.0.in_schema", branchName),
			},
		},
	}
}

func GrantOwnership() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.GrantOwnership, CreateGrantOwnership),
		// There's no Update, because every field is marked as ForceNew
		DeleteContext: TrackingDeleteWrapper(resources.GrantOwnership, DeleteGrantOwnership),
		ReadContext:   TrackingReadWrapper(resources.GrantOwnership, ReadGrantOwnership),

		Schema: grantOwnershipSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.GrantOwnership, ImportGrantOwnership()),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportGrantOwnership() schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
		id, err := ParseGrantOwnershipId(d.Id())
		if err != nil {
			return nil, err
		}

		switch id.GrantOwnershipTargetRoleKind {
		case ToAccountGrantOwnershipTargetRoleKind:
			if err := d.Set("account_role_name", id.AccountRoleName.Name()); err != nil {
				return nil, err
			}
		case ToDatabaseGrantOwnershipTargetRoleKind:
			if err := d.Set("database_role_name", id.DatabaseRoleName.FullyQualifiedName()); err != nil {
				return nil, err
			}
		}

		if id.OutboundPrivilegesBehavior != nil {
			if err := d.Set("outbound_privileges", *id.OutboundPrivilegesBehavior); err != nil {
				return nil, err
			}
		}

		switch id.Kind {
		case OnObjectGrantOwnershipKind:
			data := id.Data.(*OnObjectGrantOwnershipData)

			onObject := make(map[string]any)
			onObject["object_type"] = data.ObjectType.String()
			if objectName, ok := any(data.ObjectName).(sdk.AccountObjectIdentifier); ok {
				onObject["object_name"] = objectName.Name()
			} else {
				onObject["object_name"] = data.ObjectName.FullyQualifiedName()
			}

			if err := d.Set("on", []any{onObject}); err != nil {
				return nil, err
			}
		case OnAllGrantOwnershipKind, OnFutureGrantOwnershipKind:
			data := id.Data.(*BulkOperationGrantData)

			on := make(map[string]any)
			onAllOrFuture := make(map[string]any)
			onAllOrFuture["object_type_plural"] = data.ObjectNamePlural.String()
			switch data.Kind {
			case InDatabaseBulkOperationGrantKind:
				onAllOrFuture["in_database"] = data.Database.Name()
			case InSchemaBulkOperationGrantKind:
				onAllOrFuture["in_schema"] = data.Schema.FullyQualifiedName()
			}

			switch id.Kind {
			case OnAllGrantOwnershipKind:
				on["all"] = []any{onAllOrFuture}
			case OnFutureGrantOwnershipKind:
				on["future"] = []any{onAllOrFuture}
			}

			if err := d.Set("on", []any{on}); err != nil {
				return nil, err
			}
		}

		return []*schema.ResourceData{d}, nil
	}
}

func CreateGrantOwnership(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := createGrantOwnershipIdFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	grantOn, err := getOwnershipGrantOn(d)
	if err != nil {
		return diag.FromErr(err)
	}

	grantTo, err := getOwnershipGrantTo(d)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.Grants.GrantOwnership(
		ctx,
		*grantOn,
		grantTo,
		getOwnershipGrantOpts(id),
	)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "An error occurred during grant ownership",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", id.String(), err),
			},
		}
	}

	d.SetId(id.String())

	return ReadGrantOwnership(ctx, d, meta)
}

func DeleteGrantOwnership(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := ParseGrantOwnershipId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}

	grantOn, err := getOwnershipGrantOn(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if grantOn.Future != nil {
		// TODO (SNOW-1182623): Still waiting for the response on the behavior/SQL syntax we should use here
		log.Printf("[WARN] Unsupported operation, please revoke ownership transfer manually")
	} else {
		accountRoleName, err := client.ContextFunctions.CurrentRole(ctx)
		if err != nil {
			return diag.FromErr(err)
		}

		err = client.Grants.GrantOwnership( // TODO: Should we always set outbound privileges to COPY in delete operation or set it to the config value?
			ctx,
			*grantOn,
			sdk.OwnershipGrantTo{
				AccountRoleName: sdk.Pointer(accountRoleName),
			},
			getOwnershipGrantOpts(id),
		)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "An error occurred when transferring ownership back to the original role",
					Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
				},
			}
		}
	}

	d.SetId("")

	return nil
}

func ReadGrantOwnership(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id, err := ParseGrantOwnershipId(d.Id())
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to parse internal identifier",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}

	opts, expectedGrantedOn := prepareShowGrantsRequestForGrantOwnership(id)
	if opts == nil {
		return nil
	}

	client := meta.(*provider.Context).Client

	grants, err := client.Grants.Show(ctx, opts)
	if err != nil {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to retrieve grants. Marking the resource as removed.",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}

	ownershipFound := false

	for _, grant := range grants {
		if grant.Privilege != "OWNERSHIP" {
			continue
		}

		// Future grants do not have grantedBy, only current grants do.
		// If grantedby is an empty string, it means terraform could not have created the grant
		if (opts.Future == nil || !*opts.Future) && grant.GrantedBy.Name() == "" {
			continue
		}

		// grant_on is for future grants, granted_on is for current grants.
		// They function the same way though in a test for matching the object type
		if expectedGrantedOn != grant.GrantedOn && expectedGrantedOn != grant.GrantOn {
			continue
		}

		switch id.GrantOwnershipTargetRoleKind {
		case ToAccountGrantOwnershipTargetRoleKind:
			if grant.GranteeName.Name() == id.AccountRoleName.Name() {
				ownershipFound = true
			}
		case ToDatabaseGrantOwnershipTargetRoleKind:
			if grant.GranteeName.Name() == id.DatabaseRoleName.Name() {
				ownershipFound = true
			}
		}
	}

	if !ownershipFound {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Couldn't find OWNERSHIP privilege on the target object. Marking the resource as removed.",
				Detail:   fmt.Sprintf("Id: %s", d.Id()),
			},
		}
	}

	return nil
}

// TODO(SNOW-1229218): Make sdk.ObjectType + string objectName to sdk.ObjectIdentifier mapping available in the sdk (for all object types).
func GetOnObjectIdentifier(objectType sdk.ObjectType, objectName string) (sdk.ObjectIdentifier, error) {
	switch objectType {
	case sdk.ObjectTypeComputePool,
		sdk.ObjectTypeDatabase,
		sdk.ObjectTypeExternalVolume,
		sdk.ObjectTypeFailoverGroup,
		sdk.ObjectTypeIntegration,
		sdk.ObjectTypeNetworkPolicy,
		sdk.ObjectTypeReplicationGroup,
		sdk.ObjectTypeResourceMonitor,
		sdk.ObjectTypeRole,
		sdk.ObjectTypeUser,
		sdk.ObjectTypeWarehouse:
		return sdk.ParseAccountObjectIdentifier(objectName)
	case sdk.ObjectTypeDatabaseRole,
		sdk.ObjectTypeSchema:
		return sdk.ParseDatabaseObjectIdentifier(objectName)
	case sdk.ObjectTypeAggregationPolicy,
		sdk.ObjectTypeAlert,
		sdk.ObjectTypeAuthenticationPolicy,
		sdk.ObjectTypeDataMetricFunction,
		sdk.ObjectTypeDynamicTable,
		sdk.ObjectTypeEventTable,
		sdk.ObjectTypeExternalTable,
		sdk.ObjectTypeFileFormat,
		sdk.ObjectTypeGitRepository,
		sdk.ObjectTypeHybridTable,
		sdk.ObjectTypeIcebergTable,
		sdk.ObjectTypeImageRepository,
		sdk.ObjectTypeMaterializedView,
		sdk.ObjectTypeNetworkRule,
		sdk.ObjectTypePackagesPolicy,
		sdk.ObjectTypePipe,
		sdk.ObjectTypeMaskingPolicy,
		sdk.ObjectTypePasswordPolicy,
		sdk.ObjectTypeProjectionPolicy,
		sdk.ObjectTypeRowAccessPolicy,
		sdk.ObjectTypeSessionPolicy,
		sdk.ObjectTypeSecret,
		sdk.ObjectTypeSequence,
		sdk.ObjectTypeStage,
		sdk.ObjectTypeStream,
		sdk.ObjectTypeTable,
		sdk.ObjectTypeTag,
		sdk.ObjectTypeTask,
		sdk.ObjectTypeView:
		return sdk.ParseSchemaObjectIdentifier(objectName)
	case sdk.ObjectTypeFunction,
		sdk.ObjectTypeProcedure,
		sdk.ObjectTypeExternalFunction:
		return sdk.ParseSchemaObjectIdentifierWithArguments(objectName)
	case sdk.ObjectTypeColumn:
		return sdk.ParseTableColumnIdentifier(objectName)

	default:
		return nil, sdk.NewError(fmt.Sprintf("object_type %s is not supported, please create a feature request for the provider if given object_type should be supported", objectType))
	}
}

func getOwnershipGrantOn(d *schema.ResourceData) (*sdk.OwnershipGrantOn, error) {
	ownershipGrantOn := new(sdk.OwnershipGrantOn)

	on := d.Get("on").([]any)[0].(map[string]any)
	onObjectType := on["object_type"].(string)
	onObjectName := on["object_name"].(string)
	onAll := on["all"].([]any)
	onFuture := on["future"].([]any)

	switch {
	case len(onObjectType) > 0 && len(onObjectName) > 0:
		objectType := sdk.ObjectType(strings.ToUpper(onObjectType))
		objectName, err := GetOnObjectIdentifier(objectType, onObjectName)
		if err != nil {
			return nil, err
		}
		ownershipGrantOn.Object = &sdk.Object{
			ObjectType: objectType,
			Name:       objectName,
		}
	case len(onAll) > 0:
		grantOnSchemaObjectIn, err := getGrantOnSchemaObjectIn(onAll[0].(map[string]any))
		if err != nil {
			return nil, err
		}
		ownershipGrantOn.All = grantOnSchemaObjectIn
	case len(onFuture) > 0:
		grantOnSchemaObjectIn, err := getGrantOnSchemaObjectIn(onFuture[0].(map[string]any))
		if err != nil {
			return nil, err
		}
		ownershipGrantOn.Future = grantOnSchemaObjectIn
	}

	return ownershipGrantOn, nil
}

func getOwnershipGrantTo(d *schema.ResourceData) (sdk.OwnershipGrantTo, error) {
	var ownershipGrantTo sdk.OwnershipGrantTo

	if accountRoleName, ok := d.GetOk("account_role_name"); ok {
		accountRoleId, err := sdk.ParseAccountObjectIdentifier(accountRoleName.(string))
		if err != nil {
			return ownershipGrantTo, err
		}
		ownershipGrantTo.AccountRoleName = &accountRoleId
	}

	if databaseRoleName, ok := d.GetOk("database_role_name"); ok {
		databaseRoleId, err := sdk.ParseDatabaseObjectIdentifier(databaseRoleName.(string))
		if err != nil {
			return ownershipGrantTo, err
		}
		ownershipGrantTo.DatabaseRoleName = sdk.Pointer(databaseRoleId)
	}

	return ownershipGrantTo, nil
}

func getOwnershipGrantOpts(id *GrantOwnershipId) *sdk.GrantOwnershipOptions {
	opts := new(sdk.GrantOwnershipOptions)

	if id != nil && id.OutboundPrivilegesBehavior != nil {
		outboundPrivileges := id.OutboundPrivilegesBehavior.ToOwnershipCurrentGrantsOutboundPrivileges()
		if outboundPrivileges != nil {
			opts.CurrentGrants = &sdk.OwnershipCurrentGrants{
				OutboundPrivileges: *outboundPrivileges,
			}
		}
	}

	return opts
}

func prepareShowGrantsRequestForGrantOwnership(id *GrantOwnershipId) (*sdk.ShowGrantOptions, sdk.ObjectType) {
	opts := new(sdk.ShowGrantOptions)
	var expectedGrantedOn sdk.ObjectType

	switch id.Kind {
	case OnObjectGrantOwnershipKind:
		data := id.Data.(*OnObjectGrantOwnershipData)
		switch data.ObjectType {
		case sdk.ObjectTypeDatabaseRole:
			expectedGrantedOn = sdk.ObjectTypeRole
		default:
			expectedGrantedOn = data.ObjectType
		}
		opts.On = &sdk.ShowGrantsOn{
			Object: &sdk.Object{
				ObjectType: data.ObjectType,
				Name:       data.ObjectName,
			},
		}
	case OnAllGrantOwnershipKind: // TODO: discuss if we want to let users do this (lose control over ownership for all objects in x during delete operation - we can also add a flag that would skip delete operation when on_all is set)
		switch data := id.Data.(*BulkOperationGrantData); data.Kind {
		case InDatabaseBulkOperationGrantKind:
			log.Printf("[INFO] Show with on.all option is skipped. No changes in ownership on all %s in database %s in Snowflake will be detected.", data.ObjectNamePlural, data.Database)
		case InSchemaBulkOperationGrantKind:
			log.Printf("[INFO] Show with on.all option is skipped. No changes in ownership on all %s in schema %s in Snowflake will be detected.", data.ObjectNamePlural, data.Schema)
		}
		return nil, ""
	case OnFutureGrantOwnershipKind:
		data := id.Data.(*BulkOperationGrantData)
		expectedGrantedOn = data.ObjectNamePlural.Singular()
		opts.Future = sdk.Bool(true)

		switch data.Kind {
		case InDatabaseBulkOperationGrantKind:
			opts.In = &sdk.ShowGrantsIn{
				Database: data.Database,
			}
		case InSchemaBulkOperationGrantKind:
			opts.In = &sdk.ShowGrantsIn{
				Schema: data.Schema,
			}
		}
	}

	return opts, expectedGrantedOn
}

func createGrantOwnershipIdFromSchema(d *schema.ResourceData) (*GrantOwnershipId, error) {
	id := new(GrantOwnershipId)
	accountRoleName, accountRoleNameOk := d.GetOk("account_role_name")
	databaseRoleName, databaseRoleNameOk := d.GetOk("database_role_name")

	switch {
	case accountRoleNameOk:
		id.GrantOwnershipTargetRoleKind = ToAccountGrantOwnershipTargetRoleKind
		accountRoleId, err := sdk.ParseAccountObjectIdentifier(accountRoleName.(string))
		if err != nil {
			return nil, err
		}
		id.AccountRoleName = accountRoleId
	case databaseRoleNameOk:
		id.GrantOwnershipTargetRoleKind = ToDatabaseGrantOwnershipTargetRoleKind
		databaseRoleId, err := sdk.ParseDatabaseObjectIdentifier(databaseRoleName.(string))
		if err != nil {
			return nil, err
		}
		id.DatabaseRoleName = databaseRoleId
	}

	outboundPrivileges, outboundPrivilegesOk := d.GetOk("outbound_privileges")
	if outboundPrivilegesOk {
		switch OutboundPrivilegesBehavior(outboundPrivileges.(string)) {
		case CopyOutboundPrivilegesBehavior:
			id.OutboundPrivilegesBehavior = sdk.Pointer(CopyOutboundPrivilegesBehavior)
		case RevokeOutboundPrivilegesBehavior:
			id.OutboundPrivilegesBehavior = sdk.Pointer(RevokeOutboundPrivilegesBehavior)
		}
	}

	grantedOn := d.Get("on").([]any)[0].(map[string]any)
	objectType := grantedOn["object_type"].(string)
	objectName := grantedOn["object_name"].(string)
	all := grantedOn["all"].([]any)
	future := grantedOn["future"].([]any)

	switch {
	case len(objectType) > 0 && len(objectName) > 0:
		id.Kind = OnObjectGrantOwnershipKind
		objectType := sdk.ObjectType(objectType)
		objectName, err := GetOnObjectIdentifier(objectType, objectName)
		if err != nil {
			return nil, err
		}
		id.Data = &OnObjectGrantOwnershipData{
			ObjectType: objectType,
			ObjectName: objectName,
		}
	case len(all) > 0:
		id.Kind = OnAllGrantOwnershipKind
		grantOn, err := getGrantOnSchemaObjectIn(all[0].(map[string]any))
		if err != nil {
			return nil, err
		}
		id.Data = getBulkOperationGrantData(grantOn)
	case len(future) > 0:
		id.Kind = OnFutureGrantOwnershipKind
		grantOn, err := getGrantOnSchemaObjectIn(future[0].(map[string]any))
		if err != nil {
			return nil, err
		}
		id.Data = getBulkOperationGrantData(grantOn)
	}

	return id, nil
}
