package resources

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// v2_23_0_ShareStateUpgrader migrates the state of the share:
//   - secure_objects_only did not exist before, so it is set to its default. Without it the default would be planned as
//     a change of the ForceNew field and replace the share.
//   - the id was the plain name of the share. It is now the fully qualified name, like for the other resources.
func v2_23_0_ShareStateUpgrader(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	if rawState == nil {
		return rawState, nil
	}
	if _, ok := rawState["secure_objects_only"]; !ok {
		rawState["secure_objects_only"] = BooleanDefault
	}
	if id, ok := rawState["id"].(string); ok && id != "" && !strings.HasPrefix(id, `"`) {
		rawState["id"] = sdk.NewAccountObjectIdentifier(id).FullyQualifiedName()
	}
	return rawState, nil
}
