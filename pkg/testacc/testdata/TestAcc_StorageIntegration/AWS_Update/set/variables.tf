variable "name" {
  type = string
}

variable "comment" {
  type = string
}

variable "allowed_locations" {
  type = set(string)
}

variable "blocked_locations" {
  type = set(string)
}

variable "aws_object_acl" {
  type = string
}

variable "aws_role_arn" {
  type = string
}

variable "external_id" {
  type = string
}
