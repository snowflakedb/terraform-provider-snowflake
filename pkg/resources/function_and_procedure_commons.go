package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readFunctionOrProcedureArguments(d *schema.ResourceData, args []sdk.NormalizedArgument) error {
	if len(args) == 0 {
		// TODO [before V1]: handle empty list
		return nil
	}
	// We do it the unusual way because the default values are not returned by SF.
	// We update what we have - leaving the defaults unchanged.
	if currentArgs, ok := d.Get("arguments").([]map[string]any); !ok {
		return fmt.Errorf("arguments must be a list")
	} else {
		for i, arg := range args {
			currentArgs[i]["arg_name"] = arg.Name
			currentArgs[i]["arg_data_type"] = arg.DataType.ToSql()
		}
		return d.Set("arguments", currentArgs)
	}
}

func importFunctionOrProcedureArguments(d *schema.ResourceData, args []sdk.NormalizedArgument) error {
	currentArgs := make([]map[string]any, len(args))
	for i, arg := range args {
		currentArgs[i] = map[string]any{
			"arg_name":      arg.Name,
			"arg_data_type": arg.DataType.ToSql(),
		}
	}
	return d.Set("arguments", currentArgs)
}

func readFunctionOrProcedureImports(d *schema.ResourceData, imports []sdk.NormalizedPath) error {
	if len(imports) == 0 {
		// don't do anything if imports not present
		return nil
	}
	imps := collections.Map(imports, func(imp sdk.NormalizedPath) map[string]any {
		return map[string]any{
			"stage_location": imp.StageLocation,
			"path_on_stage":  imp.PathOnStage,
		}
	})
	return d.Set("imports", imps)
}

// TODO: extract?
func readFunctionOrProcedureExternalAccessIntegrations(d *schema.ResourceData, externalAccessIntegrations []sdk.AccountObjectIdentifier) error {
	return d.Set("external_access_integrations", collections.Map(externalAccessIntegrations, func(id sdk.AccountObjectIdentifier) string { return id.Name() }))
}

func readFunctionOrProcedureSecrets(d *schema.ResourceData, secrets map[string]sdk.SchemaObjectIdentifier) error {
	all := make([]map[string]any, 0)
	for k, v := range secrets {
		all = append(all, map[string]any{
			"secret_variable_name": k,
			"secret_id":            v.FullyQualifiedName(),
		})
	}
	return d.Set("secrets", all)
}

func readFunctionOrProcedureTargetPath(d *schema.ResourceData, normalizedPath *sdk.NormalizedPath) error {
	if normalizedPath == nil {
		// don't do anything if imports not present
		return nil
	}
	tp := make([]map[string]any, 1)
	tp[0] = map[string]any{
		"stage_location": normalizedPath.StageLocation,
		"path_on_stage":  normalizedPath.PathOnStage,
	}
	return d.Set("target_path", tp)
}

func setExternalAccessIntegrationsInBuilder[T any](d *schema.ResourceData, setIntegrations func([]sdk.AccountObjectIdentifier) T) error {
	integrations, err := parseExternalAccessIntegrationsCommon(d)
	if err != nil {
		return err
	}
	setIntegrations(integrations)
	return nil
}

func setSecretsInBuilder[T any](d *schema.ResourceData, setSecrets func([]sdk.SecretReference) T) error {
	secrets, err := parseSecretsCommon(d)
	if err != nil {
		return err
	}
	setSecrets(secrets)
	return nil
}

func parseExternalAccessIntegrationsCommon(d *schema.ResourceData) ([]sdk.AccountObjectIdentifier, error) {
	integrations := make([]sdk.AccountObjectIdentifier, 0)
	if v, ok := d.GetOk("external_access_integrations"); ok {
		for _, i := range v.(*schema.Set).List() {
			id, err := sdk.ParseAccountObjectIdentifier(i.(string))
			if err != nil {
				return nil, err
			}
			integrations = append(integrations, id)
		}
	}
	return integrations, nil
}

func parseSecretsCommon(d *schema.ResourceData) ([]sdk.SecretReference, error) {
	secretReferences := make([]sdk.SecretReference, 0)
	if v, ok := d.GetOk("secrets"); ok {
		for _, s := range v.(*schema.Set).List() {
			name := s.(map[string]any)["secret_variable_name"].(string)
			idRaw := s.(map[string]any)["secret_id"].(string)
			id, err := sdk.ParseSchemaObjectIdentifier(idRaw)
			if err != nil {
				return nil, err
			}
			secretReferences = append(secretReferences, sdk.SecretReference{VariableName: name, Name: id})
		}
	}
	return secretReferences, nil
}
