# Authentication methods manual tests

This directory is dedicated to hold steps for manual authentication methods tests in the provider that are not possible to re-recreate in automated tests (or very hard to set up). These tests are disabled by default and require `TEST_SF_TF_ENABLE_MANUAL_TESTS` environmental variable to be set.

## Okta authenticator test
This test checks `Okta` authenticator option. It requires manual steps because of additional setup on Okta side. It assumes that `default` profile uses a standard values of account name, user, password, etc.
1. Set up a developer Okta account [here](https://developer.okta.com/signup/).
2. Go to admin panel and select Applications -> Create App Integration.
3. Create a new application with SAML 2.0 type and give it a unique name
4. Fill SAML settings - paste the URLs for the testing accounts, like `https://example.snowflakecomputing.com/fed/login` for Single sign on URL, Recipient URL, Destination URL and Audience URI (SP Entity ID)
5. Click Next and Finish
6. After the app gets created, click View SAML setup instructions
7. Save the values provided: IDP SSO URL, IDP Issuer, and X509 certificate
8. Create a new security integration in Snowflake:
```
CREATE SECURITY INTEGRATION MyIDP
TYPE=SAML2
ENABLED=true
SAML2_ISSUER='http://www.okta.com/example'
SAML2_SSO_URL='https://dev-123456.oktapreview.com/app/dev-123456_test_1/example/sso/saml'
SAML2_PROVIDER='OKTA'
SAML2_SP_INITIATED_LOGIN_PAGE_LABEL='myidp - okta'
SAML2_ENABLE_SP_INITIATED=false
SAML2_X509_CERT='<x509 cert, without headers>';
```
9. Note that Snowflake and Okta login name must match, otherwise create a temporary user with a login name matching the one in Okta.
10. Prepare a TOML config like:
```
[okta]
organization_name='ORGANIZATION_NAME'
account_name='ACCOUNT_NAME'
user='LOGIN_NAME' # This is a value used to login in Okta
password='PASSWORD' # This is a password in Okta
okta_url='https://dev-123456.okta.com' # URL of your Okta environment
```
11. Run the tests - you should be able to authenticate with Okta.


## UsernamePasswordMFA authenticator test
This test checks `UsernamePasswordMFA` authenticator option. It requires manual steps because of additional verification via MFA device. It assumes that `default` profile uses a standard values of account name, user, password, etc.
1. Make sure the user you're testing with has enabled MFA (see [docs](https://docs.snowflake.com/en/user-guide/ui-snowsight-profile#enroll-in-multi-factor-authentication-mfa)) and an MFA bypass is not set (check `mins_to_bypass_mfa` in `SHOW USERS` output for the given user).
1. After running the test, you should get pinged 3 times in MFA app:
    - The first two notifications are just test setups, also present in other acceptance tests.
    - The third notification verifies that MFA is used for the first test step.
    - For the second test step we are caching MFA token, so there is not any notification.

## UsernamePasswordMFA authenticator with passcode test
This test checks `UsernamePasswordMFA` authenticator option with using `passcode`. It requires manual steps because of additional verification via MFA device. It assumes that `default_with_passcode` profile uses a standard values of account name, user, password, etc. with `passcode` set to a value in your MFA app.
1. Make sure the user you're testing with has enabled MFA (see [docs](https://docs.snowflake.com/en/user-guide/ui-snowsight-profile#enroll-in-multi-factor-authentication-mfa)) and an MFA bypass is not set (check `mins_to_bypass_mfa` in `SHOW USERS` output for the given user).
1. After running the test, you should get pinged 2 times in MFA app:
    - The first two notifications are just test setups, also present in other acceptance tests.
    - The first step asks for permission to access your device keychain.
    - For the second test step we are caching MFA token, so there is not any notification.

## OAUTH_AUTHORIZATION_CODE with Snowflake and External IdPs test
This test checks `OAUTH_AUTHORIZATION_CODE` authenticator option for both IdPs providers. They require manual steps (signing in Snowflake or giving access to secret store).

The `oauth_authorization_code_external_idp` directory contains setup and test required for the Okta flow. Run it as following:
1. Fill the following values in `main.tfvars` file:
- `issuer` - the value from the Okta Authorization Server API.
- `audience` - the account url.
- `login_name` - the user login name. Must be the same as the login name in Okta.
- `password` - a password for the new user.
- `oauth_client_id` - the value from the Okta application.
- `oauth_client_secret` - the value from the Okta application.
- `organization_name` - the organization name of the account.
- `account_name` - the account name.
2. Run terraform commands like `terraform apply -var-file="main.tfvars"` to include these variables.
3. Run the test steps.
4. Remember to destroy created resources.

The `oauth_authorization_code_snowflake_idp` directory contains setup and test required for the Snowflake flow. Run it as following:
1. Run the Step 1 and 2 to get the credentials for the security integration
2. Fill the following values in `main.tfvars` file:
- `login_name` - the login name for the user in Snowflake.
- `oauth_client_id` - the oauth client id from the `snowflake_execute` resource.
- `oauth_client_secret` - the oauth client secret from the `snowflake_execute` resource.
- `issuer` - the account url.
- `organization_name` - the organization name of the account.
- `account_name` - the account name.
2. Run terraform commands like `terraform apply -var-file="main.tfvars"` to include these variables.
3. Run the test steps.
4. Remember to destroy created resources.


## WIF + EKS OIDC authenticator test

To test the [WIF + EKS OIDC setup]
(https://docs.snowflake.com/en/user-guide/workload-identity-federation#authenticate-to-snowflake-using-openid-connect-oidc-issuer-from-aws-kubernetes),
some additional manual steps are required to create the necessary resources in a k8s cluster and Snowflake.

Pre-requisites:
- An existing AWS EKS cluster with OIDC provider enabled.

1. Create a Kubernetes service account.
2. Create a `main.tf` file with the following provider configuration:
   ```hcl
   provider "snowflake" {
     organization_name          = "ORGANIZATION_NAME"
     account_name               = "ACCOUNT_NAME"
     user                       = "USER_NAME"
     authenticator              = "WORKLOAD_IDENTITY"
     token                      = file("<token_file_path>")
     workload_identity_provider = "OIDC"
   }
   ```
3. Build a container with terraform installed. Mount/copy the terraform files incl. the above provider configuration into the image.
4. Create a user in Snowflake with
   ```sql
   CREATE OR REPLACE USER USER_NAME
     WORKLOAD_IDENTITY = (
       TYPE = OIDC
       ISSUER = 'https://oidc.eks.<region>.amazonaws.com/id/<oidc-provider-id>'
       SUBJECT = 'system:serviceaccount:<namespace>:<service-account-name>'
     )
     TYPE = SERVICE;
   ```
5. Start a pod/job/deployment in the EKS cluster. The pod/job/deployment must link the above created service account.
