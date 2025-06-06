variable "name" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "view" {
  type = string
}

variable "copy_grants" {
  type = bool
}

variable "show_initial_rows" {
  type = string
}

variable "append_only" {
  type = string
}

variable "before" {
  type = map(string)
}

variable "comment" {
  type = string
}
