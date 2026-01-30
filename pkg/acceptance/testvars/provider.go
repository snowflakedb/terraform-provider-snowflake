package testvars

func V097CompatibleProvider() string {
	// The variable names are all uppercase because GitHub forces all env variables to be uppercase.
	return `
provider "snowflake" {
	authenticator = "JWT"
	private_key = var.V097_COMPATIBLE_PRIVATE_KEY
	private_key_passphrase = var.V097_COMPATIBLE_PRIVATE_KEY_PASSPHRASE
}

variable "V097_COMPATIBLE_PRIVATE_KEY" {
	type      = string
	sensitive = true
}

variable "V097_COMPATIBLE_PRIVATE_KEY_PASSPHRASE" {
	type      = string
	sensitive = true
}
`
}
