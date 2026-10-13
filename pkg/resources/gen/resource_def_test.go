package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/stretchr/testify/require"
)

func Test_ResourceObject_implementsHasObjectGenerationSettings(t *testing.T) {
	o := ResourceObject{ObjectGenerationSettings: &genhelpers.ObjectGenerationSettings{}}
	_, ok := any(o).(genhelpers.HasObjectGenerationSettings)
	require.True(t, ok)
}
