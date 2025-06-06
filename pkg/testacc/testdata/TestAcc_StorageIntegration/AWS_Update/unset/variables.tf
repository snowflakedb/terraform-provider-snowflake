variable "name" {
  type = string
}

variable "allowed_locations" {
  type = set(string)
}

variable "aws_role_arn" {
  type = string
}
