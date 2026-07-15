# basic resource with inlined manifest
resource "snowflake_organization_listing" "basic_inlined" {
  name = "ORG_LISTING"
  manifest {
    from_string = <<-EOT
title: "My Organization Listing"
subtitle: "subtitle"
description: "Sharing data within our Snowflake organization"
organization_targets:
  access:
  - all_internal_accounts: true
locations:
  access_regions:
  - name: "ALL"
EOT
  }
}

# basic resource with manifest in a stage
resource "snowflake_organization_listing" "basic_staged" {
  name = "ORG_LISTING"
  manifest {
    from_stage {
      stage = snowflake_stage.test_stage.fully_qualified_name
    }
  }
}

# complete resource with inlined manifest and share
resource "snowflake_organization_listing" "complete_inlined" {
  name = "ORG_LISTING"
  manifest {
    from_string = <<-EOT
title: "My Organization Listing"
subtitle: "subtitle"
description: "Sharing data within our Snowflake organization"
organization_targets:
  access:
  - all_internal_accounts: true
locations:
  access_regions:
  - name: "ALL"
EOT
  }

  share = snowflake_share.test_share.fully_qualified_name
  # or
  # application_package = "test_application_package"

  publish = "true"
  comment = "Organization listing for internal data sharing"
}

# complete resource with manifest in a stage and version management
resource "snowflake_organization_listing" "complete_staged" {
  name = "ORG_LISTING"
  manifest {
    from_stage {
      stage           = snowflake_stage.test_stage.fully_qualified_name
      location        = "path/to/manifest"
      version_name    = "v1.0.0"
      version_comment = "Initial version of the organization listing manifest"
    }
  }

  share = snowflake_share.test_share.fully_qualified_name

  publish = "false"
  comment = "Organization listing with staged manifest"
}
