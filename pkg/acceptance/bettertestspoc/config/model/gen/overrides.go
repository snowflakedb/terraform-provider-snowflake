package gen

// TODO [SNOW-1501905]: extract all overrides to object definitions
var multilineAttributesOverrides = map[string][]string{
	"User":                             {"rsa_public_key", "rsa_public_key_2"},
	"ServiceUser":                      {"rsa_public_key", "rsa_public_key_2"},
	"LegacyServiceUser":                {"rsa_public_key", "rsa_public_key_2"},
	"FunctionJava":                     {"function_definition"},
	"FunctionJavascript":               {"function_definition"},
	"FunctionPython":                   {"function_definition"},
	"FunctionScala":                    {"function_definition"},
	"FunctionSql":                      {"function_definition"},
	"ProcedureJava":                    {"procedure_definition"},
	"ProcedureJavascript":              {"procedure_definition"},
	"ProcedurePython":                  {"procedure_definition"},
	"ProcedureScala":                   {"procedure_definition"},
	"ProcedureSql":                     {"procedure_definition"},
	"Account":                          {"admin_rsa_public_key"},
	"Saml2SecurityIntegration":         {"saml2_x509_cert"},
	"OauthIntegrationForCustomClients": {"oauth_client_rsa_public_key", "oauth_client_rsa_public_key_2"},
	// TODO(SNOW-1501905): Uncomment once listings are generated
	// "Listing":                          {"manifest.0.from_string"},
}

var complexListAttributesOverrides = map[string]map[string]string{
	"ExternalVolume":  {"storage_location": "sdk.ExternalVolumeStorageLocation"},
	"MaskingPolicy":   {"argument": "sdk.TableColumnSignature"},
	"RowAccessPolicy": {"argument": "sdk.TableColumnSignature"},
	"SemanticView":    {"tables": "sdk.LogicalTable"},
	"TagAssociation":  {"object_identifiers": "sdk.ObjectIdentifier"},
	// TODO [SNOW-1348114]: use better type for override (not null and default are currently not supported)
	"Table":        {"column": "sdk.TableColumnSignature"},
	"SemanticView": {"tables": "sdk.LogicalTable"},
}
