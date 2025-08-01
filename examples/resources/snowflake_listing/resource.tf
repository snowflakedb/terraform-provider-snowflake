# basic resource with inlined manifest
resource "snowflake_listing" "basic_inlined" {
  name = "LISTING"
  manifest {
    from_string = <<-EOT
title: title
subtitle: subtitle
description: description
listing_terms:
  type: OFFLINE
EOT
  }
}

# basic resource with manifest in a stage
resource "snowflake_listing" "basic_staged" {
  name = "LISTING"
  manifest {
    from_stage = {
      stage = snowflake_stage.test_stage.fully_qualified_name
    }
  }
}

# complete resource with inlined manifest
resource "snowflake_listing" "basic_inlined" {
  name = "LISTING"
  manifest {
    from_string = <<-EOT
title: title
subtitle: subtitle
description: description
listing_terms:
  type: OFFLINE
EOT
  }

  share = snowflake_share.test_share.fully_qualified_name
  # or
  application_package = "test_application_package"

  publish = true
  comment = "This is a comment for the listing"
}

# complete resource with manifest in a stage
resource "snowflake_listing" "basic_staged" {
  name = "LISTING"
  manifest {
    from_stage = {
      stage = snowflake_stage.test_stage.fully_qualified_name
    }
  }

  share = snowflake_share.test_share.fully_qualified_name
  # or
  application_package = "test_application_package"

  publish = true
  comment = "This is a comment for the listing"
}
