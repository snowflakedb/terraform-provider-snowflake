variable "comment" {
  type = string
}
variable "enabled" {
  type = bool
}
variable "name" {
  type = string
}
variable "oauth_access_token_validity" {
  type = number
}
variable "oauth_authorization_endpoint" {
  type = string
}
variable "oauth_client_auth_method" {
  type = string
}
variable "oauth_client_id" {
  type = string
}
variable "oauth_client_secret" {
  type = string
}
variable "oauth_refresh_token_validity" {
  type = number
}
variable "oauth_token_endpoint" {
  type = string
}
variable "oauth_allowed_scopes" {
  type = set(string)
}
