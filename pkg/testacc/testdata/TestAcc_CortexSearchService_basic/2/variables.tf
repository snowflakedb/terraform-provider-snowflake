

variable "name" {
  type = string
}

variable "on" {
  type = string
}

variable "attributes" {
  type = set(string)
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "warehouse" {
  type = string
}

variable "query" {
  type = string
}

variable "comment" {
  type = string
}

variable "table_name" {
  type = string
}

variable "embedding_model" {
  type = string
}
