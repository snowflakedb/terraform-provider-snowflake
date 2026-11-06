package invokeactionassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type pipeExecutionStateCheck struct {
	id       sdk.SchemaObjectIdentifier
	expected sdk.PipeExecutionState
}

func (c *pipeExecutionStateCheck) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		state, err := testClient.Pipe.ExecutionState(t, c.id)
		if err != nil {
			return err
		}
		if state != c.expected {
			return fmt.Errorf("expected pipe %s execution state to be %s, got %s", c.id.FullyQualifiedName(), c.expected, state)
		}
		return nil
	}
}

// PipeHasExecutionState asserts the pipe execution state returned by SYSTEM$PIPE_STATUS.
func PipeHasExecutionState(t *testing.T, id sdk.SchemaObjectIdentifier, expected sdk.PipeExecutionState) assert.TestCheckFuncProvider {
	t.Helper()
	return &pipeExecutionStateCheck{id: id, expected: expected}
}
