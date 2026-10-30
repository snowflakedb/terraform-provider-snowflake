variable "name" {
  type = string
}
variable "comment" {
  type = string
}
variable "allow_writes" {
  type = string
}
variable "s3_storage_locations" {
  type = list(map(string))
}
variable "gcs_storage_locations" {
  type = list(map(string))
}
variable "azure_storage_locations" {
  type = list(map(string))
}
