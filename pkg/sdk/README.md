# Snowflake Go SDK

[secondary_test_account] credentials are required in the Snowflake profile if running integration tests that provision resources in multiple accounts:

Required environment variable to run sweepers (cleanup up resources created by integration tests):

```
TEST_SF_TF_ENABLE_SWEEP=1
```
Required environment variable to test creating an account. Note that this cannot be cleaned up by sweepers:

```
TEST_SF_TF_TEST_ACCOUNT_CREATE=1
```

## SQL clause types

| ddl tag            | function              | output format                                                                 |
| ------------------ | --------------------- | ----------------------------------------------------------------------------- |
| `ddl:"static"`     | `sqlStaticClause`     | `WORD`                                                                        |
| `ddl:"keyword"`    | `sqlKeywordClause`    | `"WORD"` (quotes configurable)                                                |
| `ddl:"identifier"` | `sqlIdentifierClause` | `"a.b.c"` or `OBJ_TYPE "a.b.c"`                                               |
| `ddl:"parameter"`  | `sqlParameterClause`  | `PARAM = "value"` (quotes configurable) or `PARAM = 2`                        |
| `ddl:"list"`       | `sqlListClause`       | `WORD (<subclause>, <subclause>)` (WORD, parentheses, separator configurable) |
