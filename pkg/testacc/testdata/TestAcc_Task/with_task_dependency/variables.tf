variable "tasks" {
  type = list(object({
    database      = string
    schema        = string
    name          = string
    started       = bool
    sql_statement = string

    # Optionals
    comment  = optional(string)
    schedule = optional(map(string))
    after    = optional(set(string))
    finalize = optional(string)

    # Parameters
    suspend_task_after_num_failures = optional(number)
  }))
}
