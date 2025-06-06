variable "name" {
  type = string
}

variable "database" {
  type = string
}

variable "schema" {
  type = string
}

variable "min_length" {
  type = number
}

variable "max_length" {
  type = number
}

variable "min_upper_case_chars" {
  type = number
}

variable "min_lower_case_chars" {
  type = number
}

variable "min_numeric_chars" {
  type = number
}

variable "min_special_chars" {
  type = number
}

variable "min_age_days" {
  type = number
}

variable "max_age_days" {
  type = number
}

variable "max_retries" {
  type = number
}

variable "lockout_time_mins" {
  type = number
}

variable "history" {
  type = number
}

variable "comment" {
  type = string
}
