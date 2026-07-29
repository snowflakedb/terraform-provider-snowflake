package customassert

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func BetweenFunc(min, max int) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		got, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("expected value to be an integer; got %s: %w", value, err)
		}
		if got < min || max < got {
			return fmt.Errorf("expected value to be between %d and %d; got %d", min, max, got)
		}
		return nil
	}
}
