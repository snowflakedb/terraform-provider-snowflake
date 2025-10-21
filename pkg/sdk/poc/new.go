package poc

import (
	"fmt"
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

// TODO [SNOW-2324252]: rename this file and move it (can't be moved currently due to import cycle: sdk needs gen fro definition and , generator needs all the definitions list

type SdkObjectDef struct {
	name       string
	definition *generator.Interface
}

func GetSdkDefinitions() []*generator.Interface {
	allDefinitions := allSdkObjectDefinitions
	interfaces := make([]*generator.Interface, len(allDefinitions))
	for idx, s := range allDefinitions {
		preprocessDefinition(s.definition)
		interfaces[idx] = s.definition
	}
	return interfaces
}

// preprocessDefinition is needed because current simple builder is not ideal, should be removed later
func preprocessDefinition(definition *generator.Interface) {
	generatedStructs := make([]string, 0)
	generatedDtos := make([]string, 0)

	for _, o := range definition.Operations {
		o.ObjectInterface = definition
		if o.OptsField != nil {
			o.OptsField.Name = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
			o.OptsField.Kind = fmt.Sprintf("%s%sOptions", o.Name, o.ObjectInterface.NameSingular)
			setParent(o.OptsField)

			// TODO [SNOW-2324252]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
			// fill out StructsToGenerate; it replaces the old generateOptionsStruct and generateStruct
			structsToGenerate := make([]*generator.Field, 0)
			for _, f := range o.HelperStructs {
				if !slices.Contains(generatedStructs, f.KindNoPtr()) {
					structsToGenerate, generatedStructs = addStructToGenerate(f, structsToGenerate, generatedStructs)
				}
			}
			for _, f := range o.OptsField.Fields {
				if len(f.Fields) > 0 && !slices.Contains(generatedStructs, f.KindNoPtr()) {
					structsToGenerate, generatedStructs = addStructToGenerate(f, structsToGenerate, generatedStructs)
				}
			}
			log.Printf("[DEBUG] Structs to generate (length: %d): %v", len(structsToGenerate), structsToGenerate)
			o.StructsToGenerate = structsToGenerate

			// TODO [SNOW-2324252]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
			// fill out ObjectIdMethod and ObjectIdType; it replaces the old template executors logic
			if o.Name == string(generator.OperationKindShow) {
				// TODO [SNOW-2324252]: do we really conversion logic? The definition file should handle this.
				idKind, err := generator.ToObjectIdentifierKind(definition.IdentifierKind)
				if err != nil {
					log.Printf("[WARN] for showObjectIdMethod: %v", err)
				}
				if generator.CheckRequiredFieldsForIdMethod(definition.NameSingular, o.HelperStructs, idKind) {
					o.ObjectIdMethod = generator.NewShowObjectIDMethod(definition.NameSingular, idKind)
				}

				o.ObjectTypeMethod = generator.NewShowObjectTypeMethod(definition.NameSingular)
			}

			// TODO [SNOW-2324252]: this logic is currently the old logic adjusted. Let's clean it after new generation is working.
			// fill out DtosToGenerate; it replaces the old GenerateDtos and generateDtoDecls logic
			dtosToGenerate := make([]*generator.Field, 0)
			dtosToGenerate, generatedDtos = addDtoToGenerate(o.OptsField, dtosToGenerate, generatedDtos)
			log.Printf("[DEBUG] Dtos to generate (length: %d): %v", len(dtosToGenerate), dtosToGenerate)
			o.DtosToGenerate = dtosToGenerate
		}
	}
}

func setParent(field *generator.Field) {
	for _, f := range field.Fields {
		f.Parent = field
		setParent(f)
	}
}

func addStructToGenerate(field *generator.Field, structsToGenerate []*generator.Field, generatedStructs []string) ([]*generator.Field, []string) {
	if !slices.Contains(generatedStructs, field.KindNoPtr()) {
		log.Printf("[DEBUG] Adding %s to structs to be generated", field.KindNoPtr())
		structsToGenerate = append(structsToGenerate, field)
		generatedStructs = append(generatedStructs, field.KindNoPtr())
	}

	for _, f := range field.Fields {
		if len(f.Fields) > 0 && !slices.Contains(generatedStructs, f.Name) {
			structsToGenerate, generatedStructs = addStructToGenerate(f, structsToGenerate, generatedStructs)
		}
	}
	return structsToGenerate, generatedStructs
}

func addDtoToGenerate(field *generator.Field, dtosToGenerate []*generator.Field, generatedDtos []string) ([]*generator.Field, []string) {
	if !slices.Contains(generatedDtos, field.DtoDecl()) {
		log.Printf("[DEBUG] Adding %s to structs to be generated", field.DtoDecl())
		dtosToGenerate = append(dtosToGenerate, field)
		generatedDtos = append(generatedDtos, field.DtoDecl())

		for _, f := range field.Fields {
			if f.IsStruct() {
				dtosToGenerate, generatedDtos = addDtoToGenerate(f, dtosToGenerate, generatedDtos)
			}
		}
	}
	return dtosToGenerate, generatedDtos
}

func ExtendInterface(path string) func(*generator.Interface, *genhelpers.PreambleModel) *generator.Interface {
	return func(i *generator.Interface, preamble *genhelpers.PreambleModel) *generator.Interface {
		i.PreambleModel = preamble
		i.PathToDtoBuilderGen = path
		return i
	}
}

var allSdkObjectDefinitions = []SdkObjectDef{
	{
		name:       "NetworkPolicies",
		definition: sdk.NetworkPoliciesDef,
	},
	{
		name:       "SessionPolicies",
		definition: sdk.SessionPoliciesDef,
	},
	{
		name:       "Tasks",
		definition: sdk.TasksDef,
	},
	{
		name:       "Streams",
		definition: sdk.StreamsDef,
	},
	{
		name:       "ApplicationRoles",
		definition: sdk.ApplicationRolesDef,
	},
	{
		name:       "Views",
		definition: sdk.ViewsDef,
	},
	{
		name:       "Stages",
		definition: sdk.StagesDef,
	},
	{
		name:       "Functions",
		definition: sdk.FunctionsDef,
	},
	{
		name:       "Procedures",
		definition: sdk.ProceduresDef,
	},
	{
		name:       "EventTables",
		definition: sdk.EventTablesDef,
	},
	{
		name:       "ApplicationPackages",
		definition: sdk.ApplicationPackagesDef,
	},
	{
		name:       "StorageIntegration",
		definition: sdk.StorageIntegrationDef,
	},
	{
		name:       "ManagedAccounts",
		definition: sdk.ManagedAccountsDef,
	},
	{
		name:       "RowAccessPolicies",
		definition: sdk.RowAccessPoliciesDef,
	},
	{
		name:       "Applications",
		definition: sdk.ApplicationsDef,
	},
	{
		name:       "Sequences",
		definition: sdk.SequencesDef,
	},
	{
		name:       "MaterializedViews",
		definition: sdk.MaterializedViewsDef,
	},
	{
		name:       "ApiIntegrations",
		definition: sdk.ApiIntegrationsDef,
	},
	{
		name:       "NotificationIntegrations",
		definition: sdk.NotificationIntegrationsDef,
	},
	{
		name:       "ExternalFunctions",
		definition: sdk.ExternalFunctionsDef,
	},
	{
		name:       "Streamlits",
		definition: sdk.StreamlitsDef,
	},
	{
		name:       "NetworkRule",
		definition: sdk.NetworkRuleDef,
	},
	{
		name:       "SecurityIntegrations",
		definition: sdk.SecurityIntegrationsDef,
	},
	{
		name:       "CortexSearchService",
		definition: sdk.CortexSearchServiceDef,
	},
	{
		name:       "DataMetricFunctionReference",
		definition: sdk.DataMetricFunctionReferenceDef,
	},
	{
		name:       "ExternalVolumes",
		definition: sdk.ExternalVolumesDef,
	},
	{
		name:       "AuthenticationPolicies",
		definition: sdk.AuthenticationPoliciesDef,
	},
	{
		name:       "Secrets",
		definition: sdk.SecretsDef,
	},
	{
		name:       "Connection",
		definition: sdk.ConnectionDef,
	},
	{
		name:       "ImageRepositories",
		definition: sdk.ImageRepositoriesDef,
	},
	{
		name:       "ComputePools",
		definition: sdk.ComputePoolsDef,
	},
	{
		name:       "GitRepositories",
		definition: sdk.GitRepositoriesDef,
	},
	{
		name:       "Services",
		definition: sdk.ServicesDef,
	},
	{
		name:       "UserProgrammaticAccessTokens",
		definition: sdk.UserProgrammaticAccessTokensDef,
	},
	{
		name:       "Listings",
		definition: sdk.ListingsDef,
	},
	{
		name:       "OrganizationAccounts",
		definition: sdk.OrganizationAccountsDef,
	},
	{
		name:       "SemanticViews",
		definition: sdk.SemanticViewsDef,
	},
	{
		name:       "Notebooks",
		definition: sdk.NotebooksDef,
	},
}
