package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Database(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_DATABASE_TESTS"); ok {
		t.Skip("Skipping TestAccDatabase")
	}

	prefix := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	prefix2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: dbConfig(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", prefix),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment"),
					resource.TestCheckResourceAttrSet("snowflake_database.db", "data_retention_time_in_days"),
				),
			},
			// RENAME
			{
				Config: dbConfig(prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment"),
					resource.TestCheckResourceAttrSet("snowflake_database.db", "data_retention_time_in_days"),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: dbConfig2(prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.db", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_database.db", "comment", "test comment 2"),
					resource.TestCheckResourceAttr("snowflake_database.db", "data_retention_time_in_days", "3"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_database.db",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func dbConfig(prefix string) string {
	s := `
resource "snowflake_database" "db" {
	name = "%s"
	comment = "test comment"
}
`
	return fmt.Sprintf(s, prefix)
}

func dbConfig2(prefix string) string {
	s := `
resource "snowflake_database" "db" {
	name = "%s"
	comment = "test comment 2"
	data_retention_time_in_days = 3
}
`
	return fmt.Sprintf(s, prefix)
}
