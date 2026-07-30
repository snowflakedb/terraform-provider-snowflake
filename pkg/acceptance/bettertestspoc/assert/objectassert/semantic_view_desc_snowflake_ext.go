package objectassert

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SemanticViewDetailsAssert) ContainsTable(expected sdk.SemanticViewTableDetails) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.SemanticViewDetails) error {
		t.Helper()
		if !slices.ContainsFunc(o.Tables, func(d sdk.SemanticViewTableDetails) bool {
			return reflect.DeepEqual(d, expected)
		}) {
			return fmt.Errorf("expected semantic view to contain table %+v", expected)
		}
		return nil
	})
	return s
}

func (s *SemanticViewDetailsAssert) ContainsRelationship(expected sdk.SemanticViewRelationshipDetails) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.SemanticViewDetails) error {
		t.Helper()
		if !slices.ContainsFunc(o.Relationships, func(d sdk.SemanticViewRelationshipDetails) bool {
			return reflect.DeepEqual(d, expected)
		}) {
			return fmt.Errorf("expected semantic view to contain relationship %+v", expected)
		}
		return nil
	})
	return s
}

func (s *SemanticViewDetailsAssert) ContainsDimension(expected sdk.SemanticViewDimensionDetails) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.SemanticViewDetails) error {
		t.Helper()
		if !slices.ContainsFunc(o.Dimensions, func(d sdk.SemanticViewDimensionDetails) bool {
			return reflect.DeepEqual(d, expected)
		}) {
			return fmt.Errorf("expected semantic view to contain dimension %+v", expected)
		}
		return nil
	})
	return s
}

func (s *SemanticViewDetailsAssert) ContainsFact(expected sdk.SemanticViewFactDetails) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.SemanticViewDetails) error {
		t.Helper()
		if !slices.ContainsFunc(o.Facts, func(d sdk.SemanticViewFactDetails) bool {
			return reflect.DeepEqual(d, expected)
		}) {
			return fmt.Errorf("expected semantic view to contain fact %+v", expected)
		}
		return nil
	})
	return s
}

func (s *SemanticViewDetailsAssert) ContainsMetric(expected sdk.SemanticViewMetricDetails) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.SemanticViewDetails) error {
		t.Helper()
		if !slices.ContainsFunc(o.Metrics, func(d sdk.SemanticViewMetricDetails) bool {
			return reflect.DeepEqual(d, expected)
		}) {
			return fmt.Errorf("expected semantic view to contain metric %+v", expected)
		}
		return nil
	})
	return s
}
