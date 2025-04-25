package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
)

func main() {
	if len(os.Args) < 2 {
		log.Panic("Requires path as a first arg")
	}

	path := os.Args[1]
	additionalExamplesPath := filepath.Join(path, "examples", "additional")

	orderedResources := make([]string, 0)
	for key := range provider.Provider().ResourcesMap {
		orderedResources = append(orderedResources, key)
	}
	slices.Sort(orderedResources)

	deprecatedResources := make([]DeprecatedResource, 0)
	stableResources := make([]FeatureStability, 0)
	previewResources := make([]FeatureStability, 0)
	for _, key := range orderedResources {
		resource := provider.Provider().ResourcesMap[key]
		nameRelativeLink := docs.RelativeLink(key, filepath.Join("docs", "resources", strings.Replace(key, "snowflake_", "", 1)))

		if resource.DeprecationMessage != "" {
			deprecatedResources = append(deprecatedResources, newDeprecatedResource(nameRelativeLink, resource))
		}

		if slices.Contains(previewfeatures.AllPreviewFeatures, fmt.Sprintf("%s_resource", key)) {
			previewResources = append(previewResources, FeatureStability{nameRelativeLink})
		} else {
			stableResources = append(stableResources, FeatureStability{nameRelativeLink})
		}
	}

	orderedDatasources := make([]string, 0)
	for key := range provider.Provider().DataSourcesMap {
		orderedDatasources = append(orderedDatasources, key)
	}
	slices.Sort(orderedDatasources)

	deprecatedDatasources := make([]DeprecatedDatasource, 0)
	stableDatasources := make([]FeatureStability, 0)
	previewDatasources := make([]FeatureStability, 0)
	for _, key := range orderedDatasources {
		datasource := provider.Provider().DataSourcesMap[key]
		nameRelativeLink := docs.RelativeLink(key, filepath.Join("docs", "data-sources", strings.Replace(key, "snowflake_", "", 1)))

		if datasource.DeprecationMessage != "" {
			deprecatedDatasources = append(deprecatedDatasources, newDeprecatedDatasource(nameRelativeLink, datasource))
		}

		if slices.Contains(previewfeatures.AllPreviewFeatures, fmt.Sprintf("%s_datasource", key)) {
			previewDatasources = append(previewDatasources, FeatureStability{nameRelativeLink})
		} else {
			stableDatasources = append(stableDatasources, FeatureStability{nameRelativeLink})
		}
	}

	if errs := errors.Join(
		printTo(DeprecatedResourcesTemplate, DeprecatedResourcesContext{deprecatedResources}, filepath.Join(additionalExamplesPath, deprecatedResourcesFilename)),
		printTo(DeprecatedDatasourcesTemplate, DeprecatedDatasourcesContext{deprecatedDatasources}, filepath.Join(additionalExamplesPath, deprecatedDatasourcesFilename)),

		printTo(FeatureStabilityTemplate, FeatureStabilityContext{FeatureTypeResource, FeatureStateStable, make([]FeatureStability, 0)}, filepath.Join(additionalExamplesPath, stableResourcesFilename)),
		printTo(FeatureStabilityTemplate, FeatureStabilityContext{FeatureTypeDatasource, FeatureStateStable, make([]FeatureStability, 0)}, filepath.Join(additionalExamplesPath, stableDatasourcesFilename)),

		printTo(FeatureStabilityTemplate, FeatureStabilityContext{FeatureTypeResource, FeatureStatePreview, make([]FeatureStability, 0)}, filepath.Join(additionalExamplesPath, previewResourcesFilename)),
		printTo(FeatureStabilityTemplate, FeatureStabilityContext{FeatureTypeDatasource, FeatureStatePreview, make([]FeatureStability, 0)}, filepath.Join(additionalExamplesPath, previewDatasourcesFilename)),

		//printTo(FeatureStabilityTemplate, FeatureStabilityContext{FeatureTypeResource, FeatureStateStable, stableResources}, filepath.Join(additionalExamplesPath, stableResourcesFilename)),
		//printTo(FeatureStabilityTemplate, FeatureStabilityContext{FeatureTypeDatasource, FeatureStateStable, stableDatasources}, filepath.Join(additionalExamplesPath, stableDatasourcesFilename)),
		//
		//printTo(FeatureStabilityTemplate, FeatureStabilityContext{FeatureTypeResource, FeatureStatePreview, previewResources}, filepath.Join(additionalExamplesPath, previewResourcesFilename)),
		//printTo(FeatureStabilityTemplate, FeatureStabilityContext{FeatureTypeDatasource, FeatureStatePreview, previewDatasources}, filepath.Join(additionalExamplesPath, previewDatasourcesFilename)),
	); errs != nil {
		log.Fatal(errs)
	}
}

func newDeprecatedResource(nameRelativeLink string, resource *schema.Resource) DeprecatedResource {
	replacement, path, _ := docs.GetDeprecatedResourceReplacement(resource.DeprecationMessage)
	var replacementRelativeLink string
	if replacement != "" && path != "" {
		replacementRelativeLink = docs.RelativeLink(replacement, filepath.Join("docs", "resources", path))
	}

	return DeprecatedResource{
		NameRelativeLink:        nameRelativeLink,
		ReplacementRelativeLink: replacementRelativeLink,
	}
}

func newDeprecatedDatasource(nameRelativeLink string, datasource *schema.Resource) DeprecatedDatasource {
	replacement, path, _ := docs.GetDeprecatedResourceReplacement(datasource.DeprecationMessage)
	var replacementRelativeLink string
	if replacement != "" && path != "" {
		replacementRelativeLink = docs.RelativeLink(replacement, filepath.Join("docs", "data-sources", path))
	}

	return DeprecatedDatasource{
		NameRelativeLink:        nameRelativeLink,
		ReplacementRelativeLink: replacementRelativeLink,
	}
}

func printTo(template *template.Template, model any, filepath string) error {
	var writer bytes.Buffer
	err := template.Execute(&writer, model)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, writer.Bytes(), 0o600)
}
