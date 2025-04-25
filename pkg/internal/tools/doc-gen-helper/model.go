package main

type DeprecatedResourcesContext struct {
	Resources []DeprecatedResource
}

type DeprecatedResource struct {
	NameRelativeLink        string
	ReplacementRelativeLink string
}

type DeprecatedDatasourcesContext struct {
	Datasources []DeprecatedDatasource
}

type DeprecatedDatasource struct {
	NameRelativeLink        string
	ReplacementRelativeLink string
}

type FeatureType string

const (
	FeatureTypeResource   FeatureType = "resource"
	FeatureTypeDatasource FeatureType = "datasource"
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
