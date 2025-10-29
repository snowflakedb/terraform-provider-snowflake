# Test output processor

Is a simple script that takes in (from the STDIN) the output of the `go test` command with the `-json` flag enabled, 
and processes it to:
- Log any output from the tests to the console
- Write to the STDOUT the CSV formatted output

## CSV Format

The CSV contains the following headers:
- PACKAGE (package in which the test is located)
- TEST (test name)
- ACTION (action performed; one of: pass, fail)
- ELAPSED (time taken to run the test in seconds; floating point number format)

## Usage

Example usage:
```shell
TF_ACC=1 TF_LOG=DEBUG SNOWFLAKE_DRIVER_TRACING=debug SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE=true SF_TF_ACC_TEST_ENABLE_ALL_PREVIEW_FEATURES=true \
go test --tags=non_account_level_tests,account_level_tests -run TestAcc_GrantPrivilegesToAccountRole_OnSchema_ExactlyOneOf -v -timeout=20m ./pkg/testacc -json \
| go run ./pkg/scripts/test_output_processor/test_output_processor.go \
1> ./pkg/scripts/test_output_processor/output.csv
```

The script parts were grouped by new liens for clarity:
1. We set the necessary environment variables for the test run.
2. We run the `go test` command with the required tags, test name, verbosity, timeout, and JSON output.
3. We pipe (notice the `|`) the JSON output to the `test_output_processor.go` script.
4. We redirect the standard output of the script (notice the `1>`; `1` denotes STDOUT) to a file named `output.csv`. 
   
> Note: We only redirect STDOUT as logger within the script uses STDERR to print the test command output to the console
