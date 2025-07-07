//go:build exclude

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/example"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

var definitionMapping = map[string]*generator.Interface{
	"database_role_def.go":            example.DatabaseRole,
	"to_opts_optional_example_def.go": example.ToOptsOptionalExample,

	"network_policies_def.go":                sdk.NetworkPoliciesDef,
	"session_policies_def.go":                sdk.SessionPoliciesDef,
	"tasks_def.go":                           sdk.TasksDef,
	"streams_def.go":                         sdk.StreamsDef,
	"application_roles_def.go":               sdk.ApplicationRolesDef,
	"views_def.go":                           sdk.ViewsDef,
	"stages_def.go":                          sdk.StagesDef,
	"functions_def.go":                       sdk.FunctionsDef,
	"procedures_def.go":                      sdk.ProceduresDef,
	"event_tables_def.go":                    sdk.EventTablesDef,
	"application_packages_def.go":            sdk.ApplicationPackagesDef,
	"storage_integration_def.go":             sdk.StorageIntegrationDef,
	"managed_accounts_def.go":                sdk.ManagedAccountsDef,
	"row_access_policies_def.go":             sdk.RowAccessPoliciesDef,
	"applications_def.go":                    sdk.ApplicationsDef,
	"sequences_def.go":                       sdk.SequencesDef,
	"materialized_views_def.go":              sdk.MaterializedViewsDef,
	"api_integrations_def.go":                sdk.ApiIntegrationsDef,
	"notification_integrations_def.go":       sdk.NotificationIntegrationsDef,
	"external_functions_def.go":              sdk.ExternalFunctionsDef,
	"streamlits_def.go":                      sdk.StreamlitsDef,
	"network_rule_def.go":                    sdk.NetworkRuleDef,
	"security_integrations_def.go":           sdk.SecurityIntegrationsDef,
	"cortex_search_services_def.go":          sdk.CortexSearchServiceDef,
	"data_metric_function_references_def.go": sdk.DataMetricFunctionReferenceDef,
	"external_volumes_def.go":                sdk.ExternalVolumesDef,
	"authentication_policies_def.go":         sdk.AuthenticationPoliciesDef,
	"secrets_def.go":                         sdk.SecretsDef,
	"connections_def.go":                     sdk.ConnectionDef,
	"image_repository_def.go":                sdk.ImageRepositoriesDef,
	"compute_pools_def.go":                   sdk.ComputePoolsDef,
	"git_repository_def.go":                  sdk.GitRepositoriesDef,
	"services_def.go":                        sdk.ServicesDef,
	"user_programmatic_access_tokens_def.go": sdk.UserProgrammaticAccessTokensDef,
	"listings_def.go":                        sdk.ListingsDef,
}

func main() {
	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])
	definition := getDefinition(file)

	// runAllTemplatesToStdOut(definition)
	runAllTemplatesAndSave(definition, file)
	fmt.Println("Integration tests should be added manually to the pkg/sdk/testint/ directory")
}

func getDefinition(file string) *generator.Interface {
	def, ok := definitionMapping[file]
	if !ok {
		log.Panicf("Definition for key %s not found", file)
	}
	preprocessDefinition(def)
	return def
}

// preprocessDefinition is needed because current simple builder is not ideal, should be removed later
func preprocessDefinition(definition *generator.Interface) {
	for _, o := range definition.Operations {
		o.ObjectInterface = definition
		if o.OptsField != nil {
			o.OptsField.Name = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
			o.OptsField.Kind = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
			setParent(o.OptsField)
		}
	}
}

func setParent(field *generator.Field) {
	for _, f := range field.Fields {
		f.Parent = field
		setParent(f)
	}
}

func runAllTemplatesToStdOut(definition *generator.Interface) {
	writer := os.Stdout
	generator.GenerateInterface(writer, definition)
	generator.GenerateDtos(writer, definition)
	generator.GenerateImplementation(writer, definition)
	generator.GenerateUnitTests(writer, definition)
	generator.GenerateValidations(writer, definition)
}

func runAllTemplatesAndSave(definition *generator.Interface, file string) {
	fileWithoutSuffix, _ := strings.CutSuffix(file, "_def.go")
	runTemplateAndSave(definition, generator.GenerateInterface, filenameFor(fileWithoutSuffix, ""))
	runTemplateAndSave(definition, generator.GenerateDtos, filenameFor(fileWithoutSuffix, "_dto"))
	runTemplateAndSave(definition, generator.GenerateImplementation, filenameFor(fileWithoutSuffix, "_impl"))
	runTemplateAndSave(definition, generator.GenerateUnitTests, filename(fileWithoutSuffix, "_gen", "_test.go"))
	runTemplateAndSave(definition, generator.GenerateValidations, filenameFor(fileWithoutSuffix, "_validations"))
}

func runTemplateAndSave(def *generator.Interface, genFunc func(io.Writer, *generator.Interface), fileName string) {
	buffer := bytes.Buffer{}
	genFunc(&buffer, def)
	if err := genhelpers.WriteCodeToFile(&buffer, fileName); err != nil {
		log.Panicln(err)
	}
}

func filenameFor(prefix string, part string) string {
	return filename(prefix, part, "_gen.go")
}

func filename(prefix string, part string, suffix string) string {
	return fmt.Sprintf("%s%s%s", prefix, part, suffix)
}
