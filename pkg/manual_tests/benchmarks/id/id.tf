# This module provides a test ID.
terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }

}

resource "random_uuid" "test_id" {
  keeper = {
    first = "${timestamp()}"
  }
}

output "test_id" {
  value = replace(upper(random_uuid.test_id.id), "-", "_")
}
