//go:build exclude

package main

import (
	"fmt"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas/gen"
	"golang.org/x/exp/maps"
)

const (
	name    = "SDK to schema"
	version = "0.1.0"
)

func main() {
	genhelpers.NewGenerator(
		genhelpers.NewPreambleModel(name, version).
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk").
			WithImport("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"),
		gen.GetShowResultSchemaDetails,
		gen.ModelFromStructDetails,
		getFilename,
		gen.AllTemplates,
	).
		WithDescription("Generate SHOW/DESCRIBE output schemas and ToSchema mappers.").
		WithMakefileCommandPart("show-output-schemas").
		WithAdditionalObjectsDebugLogs(printAllStructsFields).
		WithAdditionalObjectsDebugLogs(printUniqueTypes).
		RunAndHandleOsReturn()
}

func getFilename(_ gen.ShowResultSchemaDetails, model gen.ShowResultSchemaModel) string {
	return model.Filename()
}

func printAllStructsFields(allStructs []gen.ShowResultSchemaDetails) {
	for _, s := range allStructs {
		fmt.Println("===========================")
		fmt.Printf("%s\n", s.Name)
		fmt.Println("===========================")
		for _, field := range s.Fields {
			fmt.Println(genhelpers.ColumnOutput(40, field.Name, field.ConcreteType, field.UnderlyingType))
		}
		fmt.Println()
	}
}

func printUniqueTypes(allStructs []gen.ShowResultSchemaDetails) {
	uniqueTypes := make(map[string]bool)
	for _, s := range allStructs {
		for _, f := range s.Fields {
			uniqueTypes[f.ConcreteType] = true
		}
	}
	fmt.Println("===========================")
	fmt.Println("Unique types")
	fmt.Println("===========================")
	keys := maps.Keys(uniqueTypes)
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Println(k)
	}
}
