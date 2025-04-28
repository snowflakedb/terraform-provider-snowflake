package main

type DeprecatedResourcesContext struct {
	Resources []DeprecatedResource
}

type DeprecatedResource struct {
	NameRelativeLink        string
	ReplacementRelativeLink string
}

type DeprecatedDataSourcesContext struct {
	DataSources []DeprecatedDataSource
}

type DeprecatedDataSource struct {
	NameRelativeLink        string
	ReplacementRelativeLink string
}

type FeatureType string

const (
	FeatureTypeResource   FeatureType = "resource"
	FeatureTypeDataSource FeatureType = "data source"
)

type FeatureState string

const (
	FeatureStateStable  FeatureState = "stable"
	FeatureStatePreview FeatureState = "preview"
)

type FeatureStabilityContext struct {
	FeatureType  FeatureType
	FeatureState FeatureState
	Features     []FeatureStability
}

type FeatureStability struct {
	NameRelativeLink string
}
