# User default database and role

This test shows how Snowflake behaves whenever you log into a user that has default namespace (database) and default role set up.
The cases are covering the situations where:
- Either database or role is not present.
- Database and role are present, but not granted to the user.
- Database and role are granted to the user, but the casing is not matched.
- Database and role are granted to the user and the casing is matching.
- Only role is granted to the user and the casing is matching.
- Database and role are created with lowercase, granted to the user.

## Setup

As a testing environment, I chose VSCode with the Snowflake extension for being able to quickly go between ACCOUNTADMIN account
and the tested one. For testing account, I added this configuration for the tested user:

```toml
[connections.Test]
account_name = '***'
username = '***'
password= '***'
host = '***'
```

The only thing that won't change between the tests is role and database, so we can create it beforehand:

```sql
CREATE ROLE TEST_ROLE;
CREATE DATABASE TEST_DATABASE;
CREATE ROLE "test_role";
CREATE DATABASE "test_database";
```

> Note: We have to use `TYPE = LEGACY_SERVICE` to be able to quickly log in using `login_name` + `password`.

> Note: It doesn't matter if you use single or double quotes for user properties (they result in the same behavior).

### 1. Either database or role is not present.

1. Create user with the non-existing default database and role.

```sql
CREATE OR REPLACE USER TEST_USER
    TYPE = LEGACY_SERVICE
    LOGIN_NAME = 'login_name'
    PASSWORD = 'password'
    DEFAULT_NAMESPACE = 'NON_EXISTING_TEST_DATABASE'
    DEFAULT_ROLE = 'NON_EXISTING_TEST_ROLE';
```

2. Log into the user (logs in successfully, but no database and role is selected in the context).

### 2. Database and role are present, but not granted to the user.

1. Replace the user with existing default database and role.

```sql
CREATE OR REPLACE USER TEST_USER
    TYPE = LEGACY_SERVICE
    LOGIN_NAME = 'login_name'
    PASSWORD = 'password'
    DEFAULT_NAMESPACE = 'TEST_DATABASE'
    DEFAULT_ROLE = 'TEST_ROLE';
```

2. Log into the user (logs in successfully, but no database and role is selected in the context).

### 3. Database and role are granted to the user, but the casing is not matched.

1. Replace the user with existing default database and role, but with lowercase names.

```sql
CREATE OR REPLACE USER TEST_USER
    TYPE = LEGACY_SERVICE
    LOGIN_NAME = 'login_name'
    PASSWORD = 'password'
    DEFAULT_NAMESPACE = 'test_database'
    DEFAULT_ROLE = 'test_role';
```

2. Grant the role to the user and grant usage on the database to the role.

```sql
GRANT ROLE TEST_ROLE TO USER TEST_USER;
GRANT USAGE ON DATABASE TEST_DATABASE TO ROLE TEST_ROLE;
```

3. Log into the user (logs in successfully, but no database and role is selected in the context).

### 4. Database and role are granted to the user and the casing is matching.

1. Replace the user with existing default database and role and exact casing.

```sql
CREATE OR REPLACE USER TEST_USER
    TYPE = LEGACY_SERVICE
    LOGIN_NAME = 'login_name'
    PASSWORD = 'password'
    DEFAULT_NAMESPACE = 'TEST_DATABASE'
    DEFAULT_ROLE = 'TEST_ROLE';
```

2. Grant the role to the user and grant usage on the database to the role.

```sql
GRANT ROLE TEST_ROLE TO USER TEST_USER;
GRANT USAGE ON DATABASE TEST_DATABASE TO ROLE TEST_ROLE;
```

3. Log into the user (logs in successfully and the database and role are selected in the context).

### 5. Only role is granted to the user and the casing is matching.

1. Replace the user with existing default database and role and exact casing.

```sql
CREATE OR REPLACE USER TEST_USER
    TYPE = LEGACY_SERVICE
    LOGIN_NAME = 'login_name'
    PASSWORD = 'password'
    DEFAULT_NAMESPACE = 'TEST_DATABASE'
    DEFAULT_ROLE = 'TEST_ROLE';
```

2. Grant the role to the user and revoke usage on the database to the role.

```sql
GRANT ROLE TEST_ROLE TO USER TEST_USER;
REVOKE USAGE ON DATABASE TEST_DATABASE FROM ROLE TEST_ROLE;
```

3. Log into the user (logs in successfully and role is selected in the context, but the database is not).

### 6. Database and role are created with lowercase, granted to the user.

#### 1. Matching casing

1. Replace the user with existing default database and role and exact casing.

```sql
CREATE OR REPLACE USER TEST_USER
    TYPE = LEGACY_SERVICE
    LOGIN_NAME = 'login_name'
    PASSWORD = 'password'
    DEFAULT_NAMESPACE = 'test_database'
    DEFAULT_ROLE = 'test_role';
```

2. Grant the role to the user and grant usage on the database to the role.

```sql
GRANT ROLE "test_role" TO USER TEST_USER;
GRANT USAGE ON DATABASE "test_database" TO ROLE "test_role";
```

3. Log into the user (logs in successfully and role is selected in the context, but the database is not).

#### 2. Additionally quoting the database

1. Replace the user with existing default database and role and exact casing.

```sql
CREATE OR REPLACE USER TEST_USER
    TYPE = LEGACY_SERVICE
    LOGIN_NAME = 'login_name'
    PASSWORD = 'password'
    DEFAULT_NAMESPACE = '"test_database"'
    DEFAULT_ROLE = 'test_role';
```

2. Grant the role to the user and grant usage on the database to the role.

```sql
GRANT ROLE "test_role" TO USER TEST_USER;
GRANT USAGE ON DATABASE "test_database" TO ROLE "test_role";
```

3. Log into the user (logs in successfully and the database and role are selected in the context).

#### 3. Additionally quoting the role

1. Replace the user with existing default database and role and exact casing.

```sql
CREATE OR REPLACE USER TEST_USER
    TYPE = LEGACY_SERVICE
    LOGIN_NAME = 'login_name'
    PASSWORD = 'password'
    DEFAULT_NAMESPACE = 'test_database'
    DEFAULT_ROLE = '"test_role"';
```

2. Grant the role to the user and grant usage on the database to the role.

```sql
GRANT ROLE "test_role" TO USER TEST_USER;
GRANT USAGE ON DATABASE "test_database" TO ROLE "test_role";
```

3. Log into the user (logs in successfully and the database and role are not selected in the context, because the value of `"test_role"` is saved as role name which is **not** perceived as a valid role name).

## Clean up

To clean up all the objects used in the tests, run the following commands.

```sql
DROP DATABASE TEST_DATABASE;
DROP ROLE TEST_ROLE;
DROP DATABASE "test_database";
DROP ROLE "test_role";
DROP USER TEST_USER;
```

## Summary

When specifying `DEFAULT_NAMESPACE` and `DEFAULT_ROLE` we have to take into the account that:
- `DEFAULT_NAMESPACE` is always being uppercased in Snowflake unless it's wrapped into double quotes, e.g. `DEFAULT_NAMESPACE = '"test_database"'`.
- `DEFAULT_ROLE` is always saving the input you are passing as is. This may cause issues when double-quoted id is passed, e.g. `DEFAULT_ROLE = '"test_role"'` will be saved as `"test_role"` (which won't work when logging into that user), and not `test_role` like in the case of `DEFAULT_NAMESPACE`.
