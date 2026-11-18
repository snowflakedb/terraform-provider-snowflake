//go:build account_level_tests

package testacc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The failover group and replication group data sources have no LIKE, and IN ACCOUNT returns every group of the account,
// so the created group is not at a known index. Other data source tests filter with LIKE until the row is at index 0.
// groupsDataSource finds the group by name instead. resultName is the attribute that holds the list, such as failover_groups.
type groupsDataSource struct {
	resultName string
}

func groupsInDataSource(resultName string) groupsDataSource {
	return groupsDataSource{resultName: resultName}
}

func (g groupsDataSource) resource(s *terraform.State, ref string) (*terraform.ResourceState, error) {
	rs, ok := s.RootModule().Resources[ref]
	if !ok {
		return nil, fmt.Errorf("data source %s not found", ref)
	}
	return rs, nil
}

func (g groupsDataSource) indexOf(rs *terraform.ResourceState, name string) (int, error) {
	count, err := strconv.Atoi(rs.Primary.Attributes[g.resultName+".#"])
	if err != nil {
		return -1, err
	}
	for i := range count {
		if rs.Primary.Attributes[fmt.Sprintf("%s.%d.show_output.0.name", g.resultName, i)] == name {
			return i, nil
		}
	}
	return -1, nil
}

// Listed checks that the group is in the data source.
func (g groupsDataSource) Listed(ref, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, err := g.resource(s, ref)
		if err != nil {
			return err
		}
		idx, err := g.indexOf(rs, name)
		if err != nil {
			return err
		}
		if idx < 0 {
			return fmt.Errorf("group %s is not listed in %s\n%s", name, ref, g.listing(rs))
		}
		return nil
	}
}

// HasShape checks the role of the group and the account locator it reports.
func (g groupsDataSource) HasShape(ref, name string, isPrimary bool, locator string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, err := g.resource(s, ref)
		if err != nil {
			return err
		}
		idx, err := g.indexOf(rs, name)
		if err != nil {
			return err
		}
		if idx < 0 {
			return fmt.Errorf("group %s not found in %s\n%s", name, ref, g.listing(rs))
		}
		prefix := fmt.Sprintf("%s.%d.show_output.0.", g.resultName, idx)
		gotPrimary := rs.Primary.Attributes[prefix+"is_primary"]
		gotLocator := rs.Primary.Attributes[prefix+"account_locator"]
		if gotPrimary != strconv.FormatBool(isPrimary) || gotLocator != locator {
			return fmt.Errorf("group %s is_primary=%s locator=%s, want is_primary=%t locator=%s\n%s",
				name, gotPrimary, gotLocator, isPrimary, locator, g.listing(rs))
		}
		return nil
	}
}

// AtIndex runs the check that the show output builder returns for the index of the group. The builders need an index.
func (g groupsDataSource) AtIndex(ref, name string, check func(idx int) resource.TestCheckFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, err := g.resource(s, ref)
		if err != nil {
			return err
		}
		idx, err := g.indexOf(rs, name)
		if err != nil {
			return err
		}
		if idx < 0 {
			return fmt.Errorf("group %s not found in %s\n%s", name, ref, g.listing(rs))
		}
		return check(idx)(s)
	}
}

func (g groupsDataSource) listing(rs *terraform.ResourceState) string {
	count, err := strconv.Atoi(rs.Primary.Attributes[g.resultName+".#"])
	if err != nil {
		return rs.Primary.ID
	}
	var b strings.Builder
	for i := range count {
		prefix := fmt.Sprintf("%s.%d.show_output.0.", g.resultName, i)
		fmt.Fprintf(&b, "  name=%s is_primary=%s locator=%s account=%s.%s\n",
			rs.Primary.Attributes[prefix+"name"],
			rs.Primary.Attributes[prefix+"is_primary"],
			rs.Primary.Attributes[prefix+"account_locator"],
			rs.Primary.Attributes[prefix+"organization_name"],
			rs.Primary.Attributes[prefix+"account_name"],
		)
	}
	return b.String()
}
