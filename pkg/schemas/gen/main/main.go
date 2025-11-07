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
		getStructDetails,
		gen.ModelFromStructDetails,
		getFilename,
		gen.AllTemplates,
	).
		WithAdditionalObjectsDebugLogs(printAllStructsFields).
		WithAdditionalObjectsDebugLogs(printUniqueTypes).
		RunAndHandleOsReturn()
}

func getStructDetails() []genhelpers.StructDetails {
	allObjects := append(gen.SdkShowResultStructs, gen.AdditionalStructs...)
	allStructsDetails := make([]genhelpers.StructDetails, len(allObjects))
	for idx, s := range allObjects {
		allStructsDetails[idx] = genhelpers.ExtractStructDetails(s)
	}
	return allStructsDetails
}

func getFilename(_ genhelpers.StructDetails, model gen.ShowResultSchemaModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_gen.go"
}

func printAllStructsFields(allStructs []genhelpers.StructDetails) {
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

func printUniqueTypes(allStructs []genhelpers.StructDetails) {
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
