variable "name" {
  type = string
}

variable "google_audience" {
  type = string
}

variable "api_allowed_prefixes" {
  type = list(string)
}

variable "api_blocked_prefixes" {
  type = list(string)
}

variable "comment" {
  type = string
}

variable "enabled" {
  type = bool
}

