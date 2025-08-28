variable "name" {
  description = "The name of the storage integration"
  type        = string
}

variable "allowed_locations" {
  description = "Explicitly limits external stages that use the integration to reference one or more storage locations"
  type        = list(string)
}

variable "aws_role_arn" {
  description = "The Amazon Resource Name (ARN) of the role"
  type        = string
}
