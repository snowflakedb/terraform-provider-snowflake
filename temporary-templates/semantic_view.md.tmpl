---
page_title: "{{.Name}} {{.Type}} - {{.ProviderName}}"
subcategory: "Preview"
description: |-
{{ if gt (len (split .Description "<deprecation>")) 1 -}}
{{ index (split .Description "<deprecation>") 1 | plainmarkdown | trimspace | prefixlines "  " }}
{{- else -}}
{{ .Description | plainmarkdown | trimspace | prefixlines "  " }}
{{- end }}
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

<!-- TODO(SNOW-2344309): address follow ups in semantic view SDK and resource -->
-> **Note** The `snowflake_semantic_view` resource currently does not handle external changes. It will be added during the resource stabilization.

-> **Note** Object renaming is currently not supported. It will be added during resource stabilization.

-> **Note** Copy Grants is currently not supported. It will be added during resource stabilization.

-> **Note** Lowercase IDs are not currently supported. It will be added during resource stabilization.

-> **Note** PRIVATE/ PUBLIC qualifiers for semantic expressions are not currently supported. They are treated as PUBLIC by default. It will be added during resource stabilization.

-> **Note** Excluding dimensions in `window_function:over_clause:partition_by` clause is not currently supported. It will be added during resource stabilization.

-> **Note** The `window_function:over_clause:order_by` clause must be a complete SQL expression, including any `[ ASC | DESC ] [ NULLS { FIRST | LAST } ]` modifiers. Support to break it down will be added during resource stabilization.

# {{.Name}} ({{.Type}})

{{ .Description | trimspace }}

{{ if .HasExample -}}
## Example Usage

-> **Note** Instead of using fully_qualified_name, you can reference objects managed outside Terraform by constructing a correct ID, consult [identifiers guide](../guides/identifiers_rework_design_decisions#new-computed-fully-qualified-name-field-in-resources).
<!-- TODO(SNOW-1634854): include an example showing both methods-->

{{ tffile .ExampleFile }}

{{- end }}

-> **Note** If a field has a default value, it is shown next to the type in the schema.

{{ .SchemaMarkdown | trimspace }}
{{- if .HasImport }}

## Import

Import is supported using the following syntax:

{{ codefile "shell" (printf "examples/resources/%s/import.sh" .Name)}}
{{- end }}
