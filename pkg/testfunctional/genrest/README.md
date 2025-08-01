# Generating Framework Code from REST API Definitions

## Purpose
This package demonstrates usage of the HashiCorp generators. See the documentation [overview](https://developer.hashicorp.com/terraform/plugin/code-generation).

This basically works in two steps:
1. Provide the REST API definition to the OpenAPI Provider Spec Generator and generate the intermediary JSON file.
2. Provide the JSON file to the Framework Code Generator and generate the code ready to be used in resource implementation.

## Workflow

### Requirements
Install the generators with the following commands:
```
go install github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@latest
go install github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework@latest
```

### Generate the intermediary JSON file
Run the following:
```bash
tfplugingen-openapi generate \
  --config generator_config.yaml \
  --output provider-code-spec.json \
  ../testdata/openapi/warehouse_modified.yaml
```
Note that the generator does not support REST API definitions from multiple files - the whole API schema must be present in one file. See the following references:
- https://github.com/pb33f/libopenapi/issues/297
- https://github.com/hashicorp/terraform-plugin-codegen-openapi/issues/89
- https://pb33f.io/libopenapi/rolodex/
This means that we use the `warehouse_modified.yaml` file instead. In this file, we appended the `common.yaml` file to the `warehouse.yaml` file and resolved `$ref` references to point to sections in the same file.

Run the following:
```bash
tfplugingen-framework scaffold resource \
  --name warehouse \
  --output-dir . \
  --package genrest
```
This command generates the `warehouse_resource.go` file with the scaffold code. This file is meant to be modified manually by the developers.

Run the following:
```bash
tfplugingen-framework generate resources \
  --input provider-code-spec.json \
  --output .
```
This command generates the `resource_warehouse` package with the model schemas.

Now, we can reference the generated schema from the scaffolded resource in `warehouse_resource.go`:
```go
	resp.Schema = resource_warehouse.WarehouseResourceSchema(ctx)
```
Note that the `tfplugingen-framework` generator outputs the resources in separate packages by default.

## Limitations
- In the `generator_config.yaml` file:
  - The provider name must not contain special characters - it must match `'^[a-z_][a-z0-9_]*$'`.
  - There is no way to specify the path prefix. This means that each path must be prefixed with `/api/v2/` manually.
