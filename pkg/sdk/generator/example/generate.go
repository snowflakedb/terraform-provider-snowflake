package example

// TODO [SNOW-3882943]: unit_tests is excluded: it generates against the pkg/sdk unit-test harness (0_sdk_unit_tests_test.go), which isn't mirrored in this example package.
//go:generate go run --tags=sdk_generation_examples ../gen/main/main.go --exclude-generation-part-names=unit_tests $SF_TF_GENERATOR_ARGS
