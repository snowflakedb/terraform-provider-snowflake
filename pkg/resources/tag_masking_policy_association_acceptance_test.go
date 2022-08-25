package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_TagMaskingPolicyAttachment(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAttachmentConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_masking_policy_attachment.test", "masking_policy_database", accName),
					resource.TestCheckResourceAttr("snowflake_tag_masking_policy_attachment.test", "masking_policy_name", accName),
					resource.TestCheckResourceAttr("snowflake_tag_masking_policy_attachment.test", "masking_policy_schema", accName),
					resource.TestCheckResourceAttr("snowflake_tag_masking_policy_attachment.test", "tag_id", fmt.Sprintf("%[1]v|%[1]v|%[1]v", accName)),
				),
			},
		},
	})
}

func tagAttachmentConfig(n string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%[1]v"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
	name = "%[1]v"
	database = snowflake_database.test.name
	comment = "Terraform acceptance test"
}

resource "snowflake_tag" "test" {
	name = "%[1]v"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	allowed_values = [""]
	comment = "Terraform acceptance test"
}

resource "snowflake_masking_policy" "test" {
	name = "%[1]v"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	value_data_type = "VARCHAR"
	masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type = "VARCHAR(16777216)"
	comment = "Terraform acceptance test"
}

resource "snowflake_tag_masking_policy_attachment" "test" {
	tag_id = snowflake_tag.test.id
	masking_policy_database = snowflake_database.test.name
	masking_policy_schema = snowflake_schema.test.name
	masking_policy_name = snowflake_masking_policy.test.name
	
  }
`, n)
}
