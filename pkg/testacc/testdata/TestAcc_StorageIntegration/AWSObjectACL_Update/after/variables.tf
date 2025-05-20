variable "name" {
  type = string
}

variable "allowed_locations" {
  type = set(string)
}

variable "aws_object_acl" {
  type = string
}

