[example]
accountname = 'account_name'
organizationname = 'organization_name'
user = 'user'
password = 'password'
warehouse = 'SNOWFLAKE'
role = 'ACCOUNTADMIN'
clientip = '1.2.3.4'
protocol = 'https'
port = 443
oktaurl = 'https://example.com'
clienttimeout = 10
jwtclienttimeout = 20
logintimeout = 30
requesttimeout = 40
jwtexpiretimeout = 50
externalbrowsertimeout = 60
maxretrycount = 1
authenticator = 'snowflake'
insecuremode = true
ocspfailopen = true
keepsessionalive = true
disabletelemetry = true
validatedefaultparameters = true
clientrequestmfatoken = true
clientstoretemporarycredential = true
tracing = 'info'
tmpdirpath = '/tmp/terraform-provider/'
disablequerycontextcache = true
includeretryreason = true
disableconsolelogin = true

[example.params]
param_key = 'param_value'
