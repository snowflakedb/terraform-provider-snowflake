package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider"
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

	orderedDataSources := make([]string, 0)
	for key := range provider.Provider().DataSourcesMap {
		orderedDataSources = append(orderedDataSources, key)
	}
	slices.Sort(orderedDataSources)

	deprecatedDataSources := make([]DeprecatedDataSource, 0)
	stableDataSources := make([]FeatureStability, 0)
	previewDataSources := make([]FeatureStability, 0)
	for _, key := range orderedDataSources {
		dataSource := provider.Provider().DataSourcesMap[key]
		nameRelativeLink := docs.RelativeLink(key, filepath.Join("docs", "data-sources", strings.Replace(key, "snowflake_", "", 1)))

		if dataSource.DeprecationMessage != "" {
			deprecatedDataSources = append(deprecatedDataSources, newDeprecatedDataSource(nameRelativeLink, dataSource))
		}

		if slices.Contains(previewfeatures.AllPreviewFeatures, fmt.Sprintf("%s_datasource", key)) {
			previewDataSources = append(previewDataSources, FeatureStability{nameRelativeLink})
		} else {
			stableDataSources = append(stableDataSources, FeatureStability{nameRelativeLink})
		}
	}

	if errs := errors.Join(
		printTo(DeprecatedResourcesTemplate, DeprecatedResourcesContext{deprecatedResources}, filepath.Join(additionalExamplesPath, deprecatedResourcesFilename)),
		printTo(DeprecatedDataSourcesTemplate, DeprecatedDataSourcesContext{deprecatedDataSources}, filepath.Join(additionalExamplesPath, deprecatedDataSourcesFilename)),

		printTo(FeatureStabilityTemplate, FeatureStabilityContext{FeatureTypeResource, FeatureStateStable, stableResources}, filepath.Join(additionalExamplesPath, stableResourcesFilename)),
		printTo(FeatureStabilityTemplate, FeatureStabilityContext{FeatureTypeDataSource, FeatureStateStable, stableDataSources}, filepath.Join(additionalExamplesPath, stableDataSourcesFilename)),

		printTo(FeatureStabilityTemplate, FeatureStabilityContext{FeatureTypeResource, FeatureStatePreview, previewResources}, filepath.Join(additionalExamplesPath, previewResourcesFilename)),
		printTo(FeatureStabilityTemplate, FeatureStabilityContext{FeatureTypeDataSource, FeatureStatePreview, previewDataSources}, filepath.Join(additionalExamplesPath, previewDataSourcesFilename)),
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

func newDeprecatedDataSource(nameRelativeLink string, dataSource *schema.Resource) DeprecatedDataSource {
	replacement, path, _ := docs.GetDeprecatedResourceReplacement(dataSource.DeprecationMessage)
	var replacementRelativeLink string
	if replacement != "" && path != "" {
		replacementRelativeLink = docs.RelativeLink(replacement, filepath.Join("docs", "data-sources", path))
	}

	return DeprecatedDataSource{
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
