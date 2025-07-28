variable "name" {
  type = string
}
variable "external_oauth_type" {
  type = string
}
variable "enabled" {
  type = bool
}
variable "external_oauth_issuer" {
  type = string
}
variable "external_oauth_snowflake_user_mapping_attribute" {
  type = string
}
variable "external_oauth_token_user_mapping_claim" {
  type = set(string)
}
variable "external_oauth_jws_keys_url" {
  type = set(string)
}
variable "external_oauth_scope_mapping_attribute" {
  type = string
}
