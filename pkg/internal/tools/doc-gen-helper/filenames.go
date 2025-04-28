package main

import "fmt"

func markdown(name string) string {
	return fmt.Sprintf(`%s.MD`, name)
}

var (
	deprecatedResourcesFilename   = markdown("deprecated_resources")
	deprecatedDataSourcesFilename = markdown("deprecated_data_sources")

	stableResourcesFilename   = markdown("stable_resources")
	stableDataSourcesFilename = markdown("stable_data_sources")

	previewResourcesFilename   = markdown("preview_resources")
	previewDataSourcesFilename = markdown("preview_data_sources")
)
