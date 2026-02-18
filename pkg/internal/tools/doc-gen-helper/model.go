package main

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"

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

type ExperimentalFeatures struct {
	ActiveExperiments       []Experiment
	DiscontinuedExperiments []Experiment
}

type Experiment struct {
	Name        string
	Description string
}

func toExperimentModel(e experimentalfeatures.Experiment) Experiment {
	return Experiment{string(e.Name()), e.Description()}
}
