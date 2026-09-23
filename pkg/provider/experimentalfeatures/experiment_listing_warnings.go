package experimentalfeatures

import "fmt"

type ExperimentListingWarning struct {
	Summary string
	Detail  string
}

func ExperimentListingWarnings(userEnabled, userDisabled []string) []ExperimentListingWarning {
	return experimentListingWarnings(allExperiments, userEnabled, userDisabled)
}

func experimentListingWarnings(experiments []Experiment, userEnabled, userDisabled []string) []ExperimentListingWarning {
	warnings := make([]ExperimentListingWarning, 0)
	for _, experiment := range experiments {
		inEnabled := isListed(experiment.name, userEnabled)
		inDisabled := isListed(experiment.name, userDisabled)
		if !inEnabled && !inDisabled {
			continue
		}
		name := string(experiment.name)
		if inEnabled && inDisabled {
			warnings = append(warnings, ExperimentListingWarning{
				Summary: "Experiment listed in both enabled and disabled lists.",
				Detail:  fmt.Sprintf("Experimental feature %s is listed in both experimental_features_enabled and experimental_features_disabled. The disabled list wins; the enabled entry is ignored.", name),
			})
		}
		switch experiment.state {
		case ExperimentalFeatureStateEnabledByDefault:
			if inEnabled {
				warnings = append(warnings, ExperimentListingWarning{
					Summary: "Experiment already enabled by default.",
					Detail:  fmt.Sprintf("Experimental feature %s is enabled by default. Listing it in experimental_features_enabled is redundant; you can remove it.", name),
				})
			}
		case ExperimentalFeatureStatePromoted:
			warnings = append(warnings, ExperimentListingWarning{
				Summary: "Promoted experiment listed in configuration.",
				Detail:  fmt.Sprintf("Experimental feature %s is now default provider behavior. Listing it in experimental_features_enabled or experimental_features_disabled has no effect. Please remove it from your configuration.", name),
			})
		case ExperimentalFeatureStateDiscontinued:
			warnings = append(warnings, ExperimentListingWarning{
				Summary: "Discontinued experiment listed in configuration.",
				Detail:  fmt.Sprintf("Experimental feature %s was discontinued and has no effect. It will be rejected in the next major version. Please remove it from your configuration.", name),
			})
		}
	}
	return warnings
}
