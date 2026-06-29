package gen

import (
	"reflect"
	"slices"

	objectassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SdkObjectShowOutputDetails struct {
	dataSourceDef *dataSourceDef
	genhelpers.SdkObjectDetails
}

var dataSourceMappingNormalized = map[string]dataSourceDef{
	// Show output - standard:
	normalized(sdk.Account{}):                {"Accounts"},
	normalized(sdk.ApiIntegration{}):         {"ApiIntegrations"},
	normalized(sdk.AuthenticationPolicy{}):   {"AuthenticationPolicies"},
	normalized(sdk.CatalogIntegration{}):     {"CatalogIntegrations"},
	normalized(sdk.ComputePool{}):            {"ComputePools"},
	normalized(sdk.CortexAgent{}):            {"CortexAgents"},
	normalized(sdk.Database{}):               {"Databases"},
	normalized(sdk.DatabaseRole{}):           {"DatabaseRoles"},
	normalized(sdk.ExternalVolume{}):         {"ExternalVolumes"},
	normalized(sdk.GitRepository{}):          {"GitRepositories"},
	normalized(sdk.ImageRepository{}):        {"ImageRepositories"},
	normalized(sdk.Listing{}):                {"Listings"},
	normalized(sdk.MaskingPolicy{}):          {"MaskingPolicies"},
	normalized(sdk.NetworkPolicy{}):          {"NetworkPolicies"},
	normalized(sdk.NetworkRuleDetails{}):     {"NetworkRules"},
	normalized(sdk.Notebook{}):               {"Notebooks"},
	normalized(sdk.PasswordPolicy{}):         {"PasswordPolicies"},
	normalized(sdk.ResourceMonitor{}):        {"ResourceMonitors"},
	normalized(sdk.RowAccessPolicy{}):        {"RowAccessPolicies"},
	normalized(sdk.Schema{}):                 {"Schemas"},
	normalized(sdk.Secret{}):                 {"Secrets"},
	normalized(sdk.SecurityIntegration{}):    {"SecurityIntegrations"},
	normalized(sdk.SemanticView{}):           {"SemanticViews"},
	normalized(sdk.Service{}):                {"Services"},
	normalized(sdk.Stage{}):                  {"Stages"},
	normalized(sdk.StorageIntegration{}):     {"StorageIntegrations"},
	normalized(sdk.StorageLifecyclePolicy{}): {"StorageLifecyclePolicies"},
	normalized(sdk.Stream{}):                 {"Streams"},
	normalized(sdk.Streamlit{}):              {"Streamlits"},
	normalized(sdk.Tag{}):                    {"Tags"},
	normalized(sdk.Task{}):                   {"Tasks"},
	normalized(sdk.User{}):                   {"Users"},
	normalized(sdk.Warehouse{}):              {"Warehouses"},

	// Describe output:
	normalized(sdk.CortexAgentDetails{}):                   {"CortexAgents"},
	normalized(sdk.ExternalVolumeStorageLocationDetails{}): {"ExternalVolumes"},
	normalized(sdk.PasswordPolicyDetails{}):                {"PasswordPolicies"},
	normalized(sdk.SessionPolicyDetails{}):                 {"SessionPolicies"},
	normalized(sdk.StorageLifecyclePolicyDetails{}):        {"StorageLifecyclePolicies"},
}

type dataSourceDef struct {
	pluralName string
}

// GetFilteredSdkObjectDetails is currently needed to filter out objects that are not resources because the same underlying list of objects is used.
func GetFilteredSdkObjectDetails() []SdkObjectShowOutputDetails {
	allDetails := objectassertgen.GetSdkObjectDetails()
	filtered := collections.Filter(allDetails, func(d genhelpers.SdkObjectDetails) bool {
		return !slices.Contains(objectNamesNotBeingResources, d.Name)
	})
	return collections.Map(filtered, func(d genhelpers.SdkObjectDetails) SdkObjectShowOutputDetails {
		v, _ := dataSourceMappingNormalized[d.Name]
		return SdkObjectShowOutputDetails{&v, d}
	})
}

var (
	objectsNotBeingResources     = []any{sdk.UserWorkloadIdentityAuthenticationMethod{}}
	objectNamesNotBeingResources = collections.Map(objectsNotBeingResources, func(o any) string {
		return reflect.ValueOf(o).Type().String()
	})
)

func normalized(t any) string {
	return reflect.ValueOf(t).Type().String()
}
