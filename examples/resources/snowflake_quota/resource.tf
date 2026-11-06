## Minimal
resource "snowflake_quota" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "quota_name"
}

## AI domains: per-user Cortex Code and AI function limits for the finance team
resource "snowflake_quota" "ai" {
  database = "database_name"
  schema   = "schema_name"
  name     = "ai_quota"
  comment  = "per-user AI quota"

  shared_resource {
    domain = "CORTEX CODE"
  }
  shared_resource {
    domain   = "AI FUNCTION"
    instance = "AI_CLASSIFY"
  }
  shared_resource {
    domain   = "CORTEX AGENT"
    instance = snowflake_cortex_agent.example.fully_qualified_name
  }

  user_tags {
    operator = "UNION"
    tag {
      tag   = snowflake_tag.cost_center.fully_qualified_name
      value = "finance"
    }
  }
  included_users = [snowflake_user.alice.name]
  excluded_users = [snowflake_user.carol.name]

  per_user_limit_monthly = 1000
  per_user_limit_weekly  = 300
  per_user_limit_daily   = 60
}

## All users except the ones with a tag
resource "snowflake_quota" "all_but_high_tier" {
  database = "database_name"
  schema   = "schema_name"
  name     = "default_ai_quota"

  shared_resource {
    domain = "AI FUNCTION"
  }

  excluded_user_tag {
    tag   = snowflake_tag.usage_tier.fully_qualified_name
    value = "high"
  }

  per_user_limit_monthly = 200
}

## Warehouse domain: per-user warehouse credit limits
resource "snowflake_quota" "warehouse" {
  database = "database_name"
  schema   = "schema_name"
  name     = "warehouse_quota"

  shared_resource {
    domain   = "WAREHOUSE"
    instance = snowflake_warehouse.example.name
  }

  per_user_limit_monthly = 100
}
