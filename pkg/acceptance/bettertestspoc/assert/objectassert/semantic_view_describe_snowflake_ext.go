package objectassert

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SemanticViewDetailsCollection struct {
	Details *sdk.SemanticViewDescribeDetails
}
type SemanticViewDetailsAssert struct {
	*assert.SnowflakeObjectAssert[SemanticViewDetailsCollection, sdk.SchemaObjectIdentifier]
}

func SemanticViewDetails(t *testing.T, id sdk.SchemaObjectIdentifier) *SemanticViewDetailsAssert {
	t.Helper()
	return &SemanticViewDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeSemanticView, id, func(testClient *helpers.TestClient) assert.ObjectProvider[SemanticViewDetailsCollection, sdk.SchemaObjectIdentifier] {
			return func(t *testing.T, id sdk.SchemaObjectIdentifier) (*SemanticViewDetailsCollection, error) {
				t.Helper()
				details, err := testClient.SemanticView.Describe(t, id)
				if err != nil {
					return nil, err
				}
				return &SemanticViewDetailsCollection{Details: details}, nil
			}
		}),
	}
}

func (s *SemanticViewDetailsAssert) HasDetailsCount(expected int) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *SemanticViewDetailsCollection) error {
		t.Helper()
		if o.Details.DescribeRowCount != expected {
			return fmt.Errorf("expected %d semantic view details; got: %d", expected, o.Details.DescribeRowCount)
		}
		return nil
	})
	return s
}

func (s *SemanticViewDetailsAssert) HasComment(expected string) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *SemanticViewDetailsCollection) error {
		t.Helper()
		if o.Details.Comment != expected {
			return fmt.Errorf("expected semantic view comment %q; got: %q", expected, o.Details.Comment)
		}
		return nil
	})
	return s
}

func (s *SemanticViewDetailsAssert) ContainsTable(expected sdk.SemanticViewTableDetails) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *SemanticViewDetailsCollection) error {
		t.Helper()
		if !slices.ContainsFunc(o.Details.Tables, func(d sdk.SemanticViewTableDetails) bool {
			return reflect.DeepEqual(d, expected)
		}) {
			return fmt.Errorf("expected semantic view to contain table %+v", expected)
		}
		return nil
	})
	return s
}

func (s *SemanticViewDetailsAssert) ContainsRelationship(expected sdk.SemanticViewRelationshipDetails) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *SemanticViewDetailsCollection) error {
		t.Helper()
		if !slices.ContainsFunc(o.Details.Relationships, func(d sdk.SemanticViewRelationshipDetails) bool {
			return reflect.DeepEqual(d, expected)
		}) {
			return fmt.Errorf("expected semantic view to contain relationship %+v", expected)
		}
		return nil
	})
	return s
}

func (s *SemanticViewDetailsAssert) ContainsDimension(expected sdk.SemanticViewDimensionDetails) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *SemanticViewDetailsCollection) error {
		t.Helper()
		if !slices.ContainsFunc(o.Details.Dimensions, func(d sdk.SemanticViewDimensionDetails) bool {
			return reflect.DeepEqual(d, expected)
		}) {
			return fmt.Errorf("expected semantic view to contain dimension %+v", expected)
		}
		return nil
	})
	return s
}

func (s *SemanticViewDetailsAssert) ContainsFact(expected sdk.SemanticViewFactDetails) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *SemanticViewDetailsCollection) error {
		t.Helper()
		if !slices.ContainsFunc(o.Details.Facts, func(d sdk.SemanticViewFactDetails) bool {
			return reflect.DeepEqual(d, expected)
		}) {
			return fmt.Errorf("expected semantic view to contain fact %+v", expected)
		}
		return nil
	})
	return s
}

func (s *SemanticViewDetailsAssert) ContainsMetric(expected sdk.SemanticViewMetricDetails) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *SemanticViewDetailsCollection) error {
		t.Helper()
		if !slices.ContainsFunc(o.Details.Metrics, func(d sdk.SemanticViewMetricDetails) bool {
			fmt.Printf("o.Details.Metrics: %+v\n", d)
			fmt.Printf("expected: %+v\n", expected)

			return reflect.DeepEqual(d, expected)
		}) {
			return fmt.Errorf("expected semantic view to contain metric %+v", expected)
		}
		return nil
	})
	return s
}
