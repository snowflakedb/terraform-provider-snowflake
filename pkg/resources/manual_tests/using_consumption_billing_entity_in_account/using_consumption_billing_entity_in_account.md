# Using consumption billing entity in the account resource

Important: To run this test, a non-production environment is should be used.

## Snowflake setup

To run this test, the environment it is run in must have more than one billing entity configured.
Consult internal team documentation for more details on how to create a new billing entity.

## Terraform configuration

After preparing the Snowflake environment, use attached Terraform configuration, fill out the random identifiers
and the newly created consumption billing entity name to validate it is working as expected.
