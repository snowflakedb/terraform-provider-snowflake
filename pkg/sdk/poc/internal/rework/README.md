> ⚠️ **Disclaimer**: This file will be deleted, it's just temporary to discuss moving the SDK generator to our common generator builder.

## SDK generator rework

The following objects are generated as a rework example:
- Sequences
- Streamlits

There are following generation parts:
- default
- dto
- dto_builders
- impl 
- unit_tests 
- validations

The reworked generator offers filtering by object name and by generation part.

Experiment with the following commands:

```shell
# generate all objects and all files
make clean-sdk generate-sdk
```
```shell
# generate all objects and chosen files only
make clean-sdk generate-sdk SF_TF_GENERATOR_ARGS='--filter-generation-part-names=default,dto,validations'
```
```shell
# generate chosen objects only and all files
make clean-sdk generate-sdk SF_TF_GENERATOR_ARGS='--filter-object-names=Sequences'
```
```shell
# generate chosen objects and chosen files only
make clean-sdk generate-sdk SF_TF_GENERATOR_ARGS='--filter-generation-part-names=default,impl --filter-object-names=Streamlits'
```
```shell
# show usage
make generate-sdk SF_TF_GENERATOR_ARGS='-h'
```
