package resources

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-1348325,SNOW-1501905]: Add this parameter to SDK
const ReplicableWithFailoverGroups = "REPLICABLE_WITH_FAILOVER_GROUPS"

var objectParameterSchema = map[string]*schema.Schema{
	"key": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		Description: fmt.Sprintf("Name of object parameter. Valid values are those in [object parameters](https://docs.snowflake.com/en/sql-reference/parameters.html#object-parameters). "+
			"With the `OBJECT_PARAMETER_PARSE_ON_SET` experiment (enabled by default), the provider supports the following parameters, and it parses and quotes their values itself: %s. "+
			"Any other parameter, and any parameter on any other object type, is sent to Snowflake as is: with `value_type`, the value is rendered by it, and without it, the value is sent as written.",
			objectParameterSupportedParametersDescription()),
	},
	"value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Value of object parameter, as a string. Constraints are the same as those for the parameters in Snowflake documentation.",
	},
	"value_type": parameterValueTypeSchema("the `key` field description, per object type", experimentalfeatures.ObjectParameterParseOnSet, suppressValueTypeForSupportedObjectParameter),
	"on_account": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If true, the object parameter will be set on the account level.",
	},
	"object_type": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      fmt.Sprintf("Type of object to which the parameter applies. Valid values are those in [object types](https://docs.snowflake.com/en/sql-reference/parameters.html#object-types). If no value is provided, then the resource will default to setting the object parameter at account level. The provider supports the parameters of the following object types: %s (see `key`). Other object types are sent as is, with the value rendered by `value_type`.", possibleValuesListed(sdk.AsStringList(objectParameterSupportedObjectTypes))),
		DiffSuppressFunc: ignoreCaseSuppressFunc,
		RequiredWith:     []string{"object_identifier"},
	},
	"object_identifier": {
		Type:         schema.TypeList,
		Optional:     true,
		MinItems:     1,
		Description:  "Specifies the object identifier for the object parameter. If no value is provided, then the resource will default to setting the object parameter at account level.",
		RequiredWith: []string{"object_type"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Name of the object to set the parameter for.",
				},
				"database": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the database that the object was created in.",
				},
				"schema": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the schema that the object was created in.",
				},
			},
		},
	},
}

func ObjectParameter() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ObjectParameterResource), TrackingCreateWrapper(resources.ObjectParameter, CreateObjectParameter)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ObjectParameterResource), TrackingReadWrapper(resources.ObjectParameter, ReadObjectParameter)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ObjectParameterResource), TrackingUpdateWrapper(resources.ObjectParameter, UpdateObjectParameter)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ObjectParameterResource), TrackingDeleteWrapper(resources.ObjectParameter, DeleteObjectParameter)),

		Schema: objectParameterSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateObjectParameter implements schema.CreateFunc.
func CreateObjectParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	key := d.Get("key").(string)
	value := d.Get("value").(string)

	parameter := sdk.ObjectParameter(key)

	o := sdk.Object{}
	if _, ok := d.GetOk("object_identifier"); ok {
		var err error
		if o, err = parameterObjectFromConfig(d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Empty when value_type is not set.
	var kind sdk.ParameterValueKind
	if valueType := d.Get("value_type").(string); valueType != "" {
		var err error
		if kind, err = sdk.ToParameterValueKind(valueType); err != nil {
			return diag.FromErr(err)
		}
	}

	onAccount := d.Get("on_account").(bool)
	parseOnSetEnabled := providerCtx.Experiments.IsEnabled(experimentalfeatures.ObjectParameterParseOnSet)

	var err error
	if onAccount {
		err = setParameterOnAccount(ctx, client, parseOnSetEnabled, key, value, kind)
	} else {
		err = setParameterOnObject(ctx, client, parseOnSetEnabled, o, parameter, value, kind)
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting object parameter err = %w", err))
	}

	var id string
	if onAccount {
		id = fmt.Sprintf("%v||", key)
	} else {
		id = fmt.Sprintf("%v|%v|%v", key, o.ObjectType, o.Name.FullyQualifiedName())
	}

	d.SetId(id)
	var p *sdk.Parameter
	if onAccount {
		p, err = client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(key))
	} else {
		p, err = client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameter(key), o)
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading object parameter err = %w", err))
	}
	err = d.Set("value", p.Value)
	if err != nil {
		return diag.FromErr(err)
	}
	return ReadObjectParameter(ctx, d, meta)
}

// ReadObjectParameter implements schema.ReadFunc.
func ReadObjectParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := d.Id()
	parts := strings.Split(id, "|")
	if len(parts) != 3 {
		parts = strings.Split(id, "❄️") // for backwards compatibility
	}
	if len(parts) != 3 {
		return diag.FromErr(fmt.Errorf("unexpected format of ID (%v), expected key|object_type|object_identifier", id))
	}
	key := parts[0]
	var p *sdk.Parameter
	var err error
	if parts[1] == "" {
		p, err = client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(key))
	} else {
		var objectIdentifier sdk.ObjectIdentifier
		objectIdentifier, err = sdk.ParseObjectIdentifierString(parts[2])
		if err != nil {
			return diag.FromErr(err)
		}
		p, err = client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameter(key), sdk.Object{
			ObjectType: sdk.ObjectType(strings.ToUpper(parts[1])),
			Name:       objectIdentifier,
		})
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading object parameter err = %w", err))
	}
	if err := d.Set("value", p.Value); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// UpdateObjectParameter implements schema.UpdateFunc.
func UpdateObjectParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return CreateObjectParameter(ctx, d, meta)
}

// DeleteObjectParameter implements schema.DeleteFunc.
func DeleteObjectParameter(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	key := d.Get("key").(string)
	// Empty when value_type is not set.
	var kind sdk.ParameterValueKind
	if valueType := d.Get("value_type").(string); valueType != "" {
		var err error
		if kind, err = sdk.ToParameterValueKind(valueType); err != nil {
			return diag.FromErr(err)
		}
	}
	parseOnSet := providerCtx.Experiments.IsEnabled(experimentalfeatures.ObjectParameterParseOnSet)
	unsetOnDelete := providerCtx.Experiments.IsEnabled(experimentalfeatures.ObjectParameterUnsetOnDelete)

	var err error
	if d.Get("on_account").(bool) {
		err = removeParameterOnAccount(ctx, client, parseOnSet, unsetOnDelete, key, kind)
	} else {
		o, objectErr := parameterObjectFromConfig(d)
		if objectErr != nil {
			return diag.FromErr(objectErr)
		}
		err = removeParameterOnObject(ctx, client, parseOnSet, unsetOnDelete, o, sdk.ObjectParameter(key), kind)
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting object parameter err = %w", err))
	}

	d.SetId("")
	return nil
}

// parameterObjectFromConfig builds the object from the object_type and object_identifier fields.
func parameterObjectFromConfig(d *schema.ResourceData) (sdk.Object, error) {
	database, schemaName, name := expandObjectIdentifier(d.Get("object_identifier").([]any))
	id, err := sdk.ParseObjectIdentifierString(FormatFullyQualifiedObjectID(database, schemaName, name))
	if err != nil {
		return sdk.Object{}, err
	}
	return sdk.Object{ObjectType: sdk.ObjectType(strings.ToUpper(d.Get("object_type").(string))), Name: id}, nil
}

// setParameterOnAccount sets a parameter on the account level (on_account = true).
func setParameterOnAccount(ctx context.Context, client *sdk.Client, parseOnSet bool, key, value string, kind sdk.ParameterValueKind) error {
	parameter := sdk.ObjectParameter(key)
	switch {
	// Legacy path: the experiment is disabled, so the value is formatted by the frozen v2.22.x table.
	case !parseOnSet:
		return withExperimentHint(client.Parameters.SetObjectParameterOnAccountLegacy(ctx, parameter, value), experimentalfeatures.ObjectParameterParseOnSet)
	// Hybrid path: the catalog sends the parameter if it supports it, otherwise the value is rendered by value_type and sent as is.
	case parseOnSet && kind != "":
		return client.Parameters.SetAccountParameterWithGenericFallback(ctx, key, kind, value)
	}

	// Catalog path with the legacy fallback: no value_type, so a parameter that the catalog does not handle goes through the legacy table.
	handled, err := client.Parameters.SetObjectParameterOnAccountViaCatalog(ctx, parameter, value)
	if err != nil || handled {
		return err
	}
	return client.Parameters.SetObjectParameterOnAccountLegacy(ctx, parameter, value)
}

// setParameterOnObject sets a parameter on a specific object.
func setParameterOnObject(ctx context.Context, client *sdk.Client, parseOnSet bool, o sdk.Object, parameter sdk.ObjectParameter, value string, kind sdk.ParameterValueKind) error {
	switch {
	// Legacy path: the experiment is disabled, so the value is sent as written.
	case !parseOnSet:
		return client.Parameters.SetObjectParameterOnObjectRaw(ctx, sdk.NewSetObjectParameterOnObjectRawParameterRequest(o.ObjectType, o.Name, string(parameter), value))
	// Hybrid path: the catalog sends the parameter if it supports the object type and the parameter, otherwise the value is rendered by value_type and sent as is.
	case parseOnSet && kind != "":
		return client.Parameters.SetObjectParameterOnObjectWithGenericFallback(ctx, o, parameter, kind, value)
	}

	// Catalog path with the legacy fallback: no value_type, so what the catalog does not handle is sent as written.
	handled, err := client.Parameters.SetObjectParameterOnObjectViaCatalog(ctx, o, parameter, value)
	if err != nil || handled {
		return err
	}
	return client.Parameters.SetObjectParameterOnObjectRaw(ctx, sdk.NewSetObjectParameterOnObjectRawParameterRequest(o.ObjectType, o.Name, string(parameter), value))
}

// unsetParameterOnAccount unsets a parameter on the account level.
func unsetParameterOnAccount(ctx context.Context, client *sdk.Client, parseOnSet, hasValueType bool, key string) error {
	parameter := sdk.ObjectParameter(key)
	switch {
	// Legacy path: the experiment is disabled.
	case !parseOnSet:
		return client.Parameters.UnsetObjectParameterOnAccountLegacy(ctx, parameter)
	// Hybrid path: value_type tells that the parameter may be unknown to the provider, so it is unset as is when the catalog does not handle it.
	case parseOnSet && hasValueType:
		return client.Parameters.UnsetAccountParameterWithGenericFallback(ctx, key)
	}

	// Catalog path with the legacy fallback.
	handled, err := client.Parameters.UnsetObjectParameterOnAccountViaCatalog(ctx, parameter)
	if err != nil || handled {
		return err
	}
	return client.Parameters.UnsetObjectParameterOnAccountLegacy(ctx, parameter)
}

// removeParameterOnAccount unsets a parameter on the account level, or sets it back to its default value.
func removeParameterOnAccount(ctx context.Context, client *sdk.Client, parseOnSet, unsetOnDelete bool, key string, kind sdk.ParameterValueKind) error {
	if unsetOnDelete {
		return unsetParameterOnAccount(ctx, client, parseOnSet, kind != "", key)
	}
	defaultParameter, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameter(key))
	if err != nil {
		return err
	}
	switch {
	// An empty default value cannot be sent, so the parameter is unset.
	case defaultParameter.Default == "":
		return unsetParameterOnAccount(ctx, client, parseOnSet, kind != "", key)
	// Legacy path: the experiment is disabled, so the default value is formatted by the table of the account parameters, like in v2.22.x.
	case !parseOnSet:
		return client.Parameters.SetAccountParameterLegacy(ctx, sdk.AccountParameter(key), defaultParameter.Default)
	}

	// Catalog path: the default value is set when the catalog handles the parameter.
	handled, err := client.Parameters.SetObjectParameterOnAccountViaCatalog(ctx, sdk.ObjectParameter(key), defaultParameter.Default)
	if err != nil || handled {
		return err
	}
	// Hybrid fallback: a parameter set with value_type is unknown to the provider, so it is unset as is.
	if kind != "" {
		return client.Parameters.UnsetAccountParameterRaw(ctx, sdk.NewUnsetAccountParameterRawParameterRequest(key))
	}
	return client.Parameters.SetAccountParameterLegacy(ctx, sdk.AccountParameter(key), defaultParameter.Default)
}

// removeParameterOnObject unsets a parameter on a specific object, or sets it back to its default value.
func removeParameterOnObject(ctx context.Context, client *sdk.Client, parseOnSet, unsetOnDelete bool, o sdk.Object, parameter sdk.ObjectParameter, kind sdk.ParameterValueKind) error {
	if unsetOnDelete {
		return unsetParameterOnObject(ctx, client, parseOnSet, o, parameter)
	}
	defaultParameter, err := client.Parameters.ShowObjectParameter(ctx, parameter, o)
	if err != nil {
		return err
	}
	defaultValue := defaultParameter.Default
	// The default value of REPLICABLE_WITH_FAILOVER_GROUPS is not settable in Snowflake, so the parameter is set to YES.
	if string(parameter) == ReplicableWithFailoverGroups {
		defaultValue = "YES"
	}
	// An empty default value cannot be sent, so the parameter is unset.
	if defaultValue == "" {
		return unsetParameterOnObject(ctx, client, parseOnSet, o, parameter)
	}
	return setParameterOnObject(ctx, client, parseOnSet, o, parameter, defaultValue, kind)
}

// unsetParameterOnObject unsets a parameter on a specific object.
func unsetParameterOnObject(ctx context.Context, client *sdk.Client, parseOnSet bool, o sdk.Object, parameter sdk.ObjectParameter) error {
	// Catalog path, when the experiment is enabled.
	if parseOnSet {
		handled, err := client.Parameters.UnsetObjectParameterOnObjectViaCatalog(ctx, o, parameter)
		if err != nil || handled {
			return err
		}
	}

	// Legacy path: what the catalog does not handle is unset as is.
	return client.Parameters.UnsetObjectParameterOnObjectRaw(ctx, sdk.NewUnsetObjectParameterOnObjectRawParameterRequest(o.ObjectType, o.Name, string(parameter)))
}

// Functions and procedures have catalog writers, but this resource does not capture argument types, so they are not listed.
var objectParameterSupportedObjectTypes = []sdk.ObjectType{
	sdk.ObjectTypeDatabase,
	sdk.ObjectTypeSchema,
	sdk.ObjectTypeWarehouse,
	sdk.ObjectTypeUser,
	sdk.ObjectTypeTask,
	sdk.ObjectTypeService,
	sdk.ObjectTypeHybridTable,
	sdk.ObjectTypeIcebergTable,
}

func objectParameterSupportedParametersDescription() string {
	scopes := make([]string, 0, 1+len(objectParameterSupportedObjectTypes))
	scopes = append(scopes, fmt.Sprintf("on the account level (`on_account = true`): %s", possibleValuesListed(sdk.AsStringList(accountParameterSupportedParameters))))
	for _, objectType := range objectParameterSupportedObjectTypes {
		scopes = append(scopes, fmt.Sprintf("`%s`: %s", objectType, possibleValuesListed(sdk.AsStringList(sdk.SupportedObjectParameters(objectType)))))
	}
	return strings.Join(scopes, "; ")
}

// TODO(next prs): The generated per-object-type parameter lists also contain names that the catalog writers do not support, so sdk.SupportedObjectParameters is used here and above.
// Once the catalog levels are fixed, both uses switch to a map from the object type to its generated All*Parameters list.
// suppressValueTypeForSupportedObjectParameter hides value_type diffs once the provider supports the parameter.
func suppressValueTypeForSupportedObjectParameter(_, _, _ string, d *schema.ResourceData) bool {
	key := strings.ToUpper(d.Get("key").(string))
	if d.Get("on_account").(bool) {
		return slices.Contains(accountParameterSupportedParameters, sdk.AccountParameter(key))
	}
	objectType := sdk.ObjectType(strings.ToUpper(d.Get("object_type").(string)))
	return slices.Contains(objectParameterSupportedObjectTypes, objectType) && slices.Contains(sdk.SupportedObjectParameters(objectType), sdk.ObjectParameter(key))
}
