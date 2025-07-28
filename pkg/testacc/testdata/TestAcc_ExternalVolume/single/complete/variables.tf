variable "name" {
  type = string
}
variable "comment" {
  type = string
}
variable "allow_writes" {
  type = string
}
variable "storage_location" {
  type = list(map(string))
}
