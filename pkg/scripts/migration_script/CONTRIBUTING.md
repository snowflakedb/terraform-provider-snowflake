# Migration script - contributing guide

<!-- TOC -->
* [Migration script - contributing guide](#migration-script---contributing-guide)
  * [1. Defining a new object type](#1-defining-a-new-object-type)
  * [2. Defining an input schema](#2-defining-an-input-schema)
  * [3. Providing an object migration function](#3-providing-an-object-migration-function)
<!-- TOC -->

The following steps will guide you through contributing to the migration script.
You will learn the project structure, expectations, and how to create or modify migrations for object types.
While the steps focus on adding a new object type, you can apply the same principles to update or improve existing migrations,
ensuring consistency and maintainability across the project.

It's important to establish the scope of the change before the actual implementation.
We would appreciate that any planned changes are discussed first (e.g., in an issue).

## 1. Defining a new object type

If you want to add a new object type, you need to define it in the migration script.
The [`program.go`](./program.go) file contains the main logic for the user interactions with the migration script through the terminal.
At the top of the file, you need to add a new constant for the object type. Remember to reflect this change 
in the help text (in `parseInputArguments` method for the Program struct) and the readme file ([syntax section](./README.md#syntax)).

Next, you can use the newly defined object type in the `generateOutput` Program method.
As providing the object migration function is the last step, you can handle the new case by returning an empty string without error.

As a last step, we need to ensure our tests cover the new object type.
Add a new test cases for the newly added object type in the [`program_test.go`](./program_test.go) file
and update any failing tests (you will need to update the help text test).

## 2. Defining an input schema

To define an input schema, you need to first understand the output of the Snowflake command that corresponds to the object type you are working with.
For almost every object type, you want to focus on SHOW commands. Sometimes the object may need to be constructed from multiple commands (e.g., SHOW + DESC + SHOW PARAMETERS) to be completely generated.

If you want to add a support for database, you would look at the output of [`SHOW DATABASES`](https://docs.snowflake.com/en/sql-reference/sql/show-databases) command.
Since all objects supported by the provider have already representation in the underlying SDK, you can look through the SDK code to find the object definition.
It's usually named in the format of `<object>Row`, so in this case we would be `databaseRow` struct.

Once you have identified the relevant Snowflake command and the corresponding SDK object, you can define the CSV schema.
It should be in the `<object_type>_converter.go` file, where `<object_type>` is the name of the object type you are working with.
It should be pretty much exactly the same as `databaseRow` struct, but with "csv" tags instead of "db" ones and different types.
Currently, only strings and booleans are supported. Given that, the CSV schema struct should look similar to the following:

```go
// Note: it can have the same name as the SDK object, because they are in different packages.
type databaseRow struct {
  CreatedOn     string `csv:"created_on"`
  Name          string `csv:"name"`
  IsDefault     string `csv:"is_default"` // Note: that databaseRow is also treating it as string, becasuse it is a string in Snowflake output. For "true" booleans, we would use bool type here.
  IsCurrent     string `csv:"is_current"`
  Origin        string `csv:"origin"`
  Owner         string `csv:"owner"`
  Comment       string `csv:"comment"`
  Options       string `csv:"options"`
  RetentionTime string `csv:"retention_time"`
  ResourceGroup string `csv:"resource_group"`
  DroppedOn     string `csv:"dropped_on"`
  Kind          string `csv:"kind"`
  OwnerRoleType string `csv:"owner_role_type"`
}
```

Under schema definition, you need to provide a conversion function that takes the CSV schema struct as input and returns the corresponding SDK object.
In our case, it would be `func (row databaseRow) convert() (*sdk.Database, error)`.

> The SDK's databaseRow has the same convert method. 
> It's worth checking as it may contain the necessary parts for converting Snowflake output for certain cases.

## 3. Providing an object migration function

Now, you need to provide a function that would take the CSV input and return the generated resources and imports in the form of string (see `HandleGrants` function).
The file with mapping function should be placed in the file named `<object_type>_migration.go`, where `<object_type>` is the name of the object type you are working with.

At the top of the function, you need to parse the CSV input into the CSV schema struct you have defined in the previous step.
This is done by the predefined `ConvertCsvInput` function (see `HandleGrants` for example usage).
Now, you can iterate over the parsed rows and generate the resources and imports.

To do this, you should use the model package in our project that contains the logic for generating resource definitions.
It should contain the resource model struct and functions for transforming the model (if not you should add them, look at [generators documentation](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/pkg/acceptance/bettertestspoc/README.md)).
They in combination with `TransformResourceModel` function produce the final resource definitions output.

To generate the import statements or blocks, you should use the `TransformImportModel` function that expects you 
to provide the import model with resource address and identifier used for importing a given object.
You should look at given resource documentation to understand how to construct the resource import (e.g., https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/database#import).

At the end, you should provide tests that cover the migration function (see [`mappings_grants_test.go`](./mappings_grants_test.go)).
In case there are any limitations of the implementation, they should be documented both in the help text and in the readme file ([syntax section](./README.md#syntax)).
