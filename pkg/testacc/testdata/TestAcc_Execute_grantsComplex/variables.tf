variable "database_grants" {
  type = list(object({
    database_name = string
    role_id       = string
    privileges    = list(string)
  }))
}
