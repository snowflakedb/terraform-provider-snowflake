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
