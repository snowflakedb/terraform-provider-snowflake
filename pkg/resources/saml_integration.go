package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"
)

var samlIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the name of the SAML2 integration. This name follows the rules for Object Identifiers. The name should be unique among security integrations in your account.",
	},
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Specifies whether this security integration is enabled or disabled.",
	},
	"saml2_issuer": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The string containing the IdP EntityID / Issuer.",
	},
	"saml2_sso_url": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The string containing the IdP SSO URL, where the user should be redirected by Snowflake (the Service Provider) with a SAML AuthnRequest message.",
	},
	"saml2_provider": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The string describing the IdP. One of the following: OKTA, ADFS, Custom.",
		ValidateFunc: validation.StringInSlice([]string{
			"OKTA", "ADFS", "CUSTOM",
		}, true),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalize := func(s string) string {
				return strings.ToUpper(strings.Replace(s, "-", "", -1))
			}
			return normalize(old) == normalize(new)
		},
	},
	"saml2_x509_cert": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Base64 encoded IdP signing certificate on a single line without the leading -----BEGIN CERTIFICATE----- and ending -----END CERTIFICATE----- markers.",
	},
	"saml2_sp_initiated_login_page_label": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The string containing the label to display after the Log In With button on the login page.",
	},
	"saml2_enable_sp_initiated": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "The Boolean indicating if the Log In With button will be shown on the login page. TRUE: displays the Log in WIth button on the login page.  FALSE: does not display the Log in With button on the login page.",
	},
	// Computed and Optionally Settable. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (SAML2_SNOWFLAKE_METADATA)
	"saml2_snowflake_x509_cert": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The Base64 encoded self-signed certificate generated by Snowflake for use with Encrypting SAML Assertions and Signed SAML Requests. You must have at least one of these features (encrypted SAML assertions or signed SAML responses) enabled in your Snowflake account to access the certificate value.",
	},
	"saml2_sign_request": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "The Boolean indicating whether SAML requests are signed. TRUE: allows SAML requests to be signed. FALSE: does not allow SAML requests to be signed.",
	},
	"saml2_requested_nameid_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The SAML NameID format allows Snowflake to set an expectation of the identifying attribute of the user (i.e. SAML Subject) in the SAML assertion from the IdP to ensure a valid authentication to Snowflake. If a value is not specified, Snowflake sends the urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress value in the authentication request to the IdP. NameID must be one of the following values: urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified, urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress, urn:oasis:names:tc:SAML:1.1:nameid-format:X509SubjectName, urn:oasis:names:tc:SAML:1.1:nameid-format:WindowsDomainQualifiedName, urn:oasis:names:tc:SAML:2.0:nameid-format:kerberos, urn:oasis:names:tc:SAML:2.0:nameid-format:persistent, urn:oasis:names:tc:SAML:2.0:nameid-format:transient .",
		ValidateFunc: validation.StringInSlice([]string{
			"urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
			"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
			"urn:oasis:names:tc:SAML:1.1:nameid-format:X509SubjectName",
			"urn:oasis:names:tc:SAML:1.1:nameid-format:WindowsDomainQualifiedName",
			"urn:oasis:names:tc:SAML:2.0:nameid-format:kerberos",
			"urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
			"urn:oasis:names:tc:SAML:2.0:nameid-format:transient",
		}, true),
	},
	"saml2_post_logout_redirect_url": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The endpoint to which Snowflake redirects users after clicking the Log Out button in the classic Snowflake web interface. Snowflake terminates the Snowflake session upon redirecting to the specified endpoint.",
	},
	"saml2_force_authn": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "The Boolean indicating whether users, during the initial authentication flow, are forced to authenticate again to access Snowflake. When set to TRUE, Snowflake sets the ForceAuthn SAML parameter to TRUE in the outgoing request from Snowflake to the identity provider. TRUE: forces users to authenticate again to access Snowflake, even if a valid session with the identity provider exists. FALSE: does not force users to authenticate again to access Snowflake.",
	},
	// Computed and Optionally Settable. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (SAML2_SNOWFLAKE_METADATA)
	"saml2_snowflake_issuer_url": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The string containing the EntityID / Issuer for the Snowflake service provider. If an incorrect value is specified, Snowflake returns an error message indicating the acceptable values to use.",
	},
	// Computed and Optionally Settable. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (SAML2_SNOWFLAKE_METADATA)
	"saml2_snowflake_acs_url": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The string containing the Snowflake Assertion Consumer Service URL to which the IdP will send its SAML authentication response back to Snowflake. This property will be set in the SAML authentication request generated by Snowflake when initiating a SAML SSO operation with the IdP. If an incorrect value is specified, Snowflake returns an error message indicating the acceptable values to use. Default: https://<account_locator>.<region>.snowflakecomputing.com/fed/login",
	},
	// Computed. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (SAML2_SNOWFLAKE_METADATA)
	"saml2_snowflake_metadata": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Metadata created by Snowflake to provide to SAML2 provider.",
	},
	// Computed. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (SAML2_DIGEST_METHODS_USED)
	"saml2_digest_methods_used": {
		Type:     schema.TypeString,
		Computed: true,
	},
	// Computed. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (SAML2_SIGNATURE_METHODS_USED)
	"saml2_signature_methods_used": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the SAML integration was created.",
	},
}

// SAMLIntegration returns a pointer to the resource representing a SAML2 security integration
func SAMLIntegration() *schema.Resource {
	return &schema.Resource{
		Create: CreateSAMLIntegration,
		Read:   ReadSAMLIntegration,
		Update: UpdateSAMLIntegration,
		Delete: DeleteSAMLIntegration,

		Schema: samlIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateSAMLIntegration implements schema.CreateFunc
func CreateSAMLIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := d.Get("name").(string)

	stmt := snowflake.SamlIntegration(name).Create()

	// Set required fields
	stmt.SetRaw(`TYPE=SAML2`)
	stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))
	stmt.SetString(`SAML2_ISSUER`, d.Get("saml2_issuer").(string))
	stmt.SetString(`SAML2_SSO_URL`, d.Get("saml2_sso_url").(string))
	stmt.SetString(`SAML2_PROVIDER`, d.Get("saml2_provider").(string))

	// Set optional fields
	if _, ok := d.GetOk("saml2_x509_cert"); ok {
		stmt.SetString(`SAML2_X509_CERT`, d.Get("saml2_x509_cert").(string))
	}

	if _, ok := d.GetOk("saml2_sp_initiated_login_page_label"); ok {
		stmt.SetString(`SAML2_SP_INITIATED_LOGIN_PAGE_LABEL`, d.Get("saml2_sp_initiated_login_page_label").(string))
	}

	if _, ok := d.GetOk("saml2_enable_sp_initiated"); ok {
		stmt.SetBool(`SAML2_ENABLE_SP_INITIATED`, d.Get("saml2_enable_sp_initiated").(bool))
	}

	if _, ok := d.GetOk("saml2_snowflake_x509_cert"); ok {
		stmt.SetString(`SAML2_SNOWFLAKE_X509_CERT`, d.Get("saml2_snowflake_x509_cert").(string))
	}

	if _, ok := d.GetOk("saml2_sign_request"); ok {
		stmt.SetString(`SAML2_SIGN_REQUEST`, d.Get("saml2_sign_request").(string))
	}

	if _, ok := d.GetOk("saml2_requested_nameid_format"); ok {
		stmt.SetString(`SAML2_REQUESTED_NAMEID_FORMAT`, d.Get("saml2_requested_nameid_format").(string))
	}

	if _, ok := d.GetOk("saml2_post_logout_redirect_url"); ok {
		stmt.SetString(`SAML2_POST_LOGOUT_REDIRECT_URL`, d.Get("saml2_post_logout_redirect_url").(string))
	}

	if _, ok := d.GetOk("saml2_force_authn"); ok {
		stmt.SetString(`SAML2_FORCE_AUTHN`, d.Get("saml2_force_authn").(string))
	}

	if _, ok := d.GetOk("saml2_snowflake_issuer_url"); ok {
		stmt.SetString(`SAML2_SNOWFLAKE_ISSUER_URL`, d.Get("saml2_snowflake_issuer_url").(string))
	}

	if _, ok := d.GetOk("saml2_snowflake_acs_url"); ok {
		stmt.SetString(`SAML2_SNOWFLAKE_ACS_URL`, d.Get("saml2_snowflake_acs_url").(string))
	}

	err := snowflake.Exec(db, stmt.Statement())
	if err != nil {
		return errors.Wrap(err, "error creating security integration")
	}

	d.SetId(name)

	return ReadSAMLIntegration(d, meta)
}

// ReadSAMLIntegration implements schema.ReadFunc
func ReadSAMLIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.SamlIntegration(id).Show()
	row := snowflake.QueryRow(db, stmt)

	// Some properties can come from the SHOW INTEGRATION call

	s, err := snowflake.ScanSamlIntegration(row)
	if err != nil {
		return errors.Wrap(err, "could not show security integration")
	}

	// Note: category must be Security or something is broken
	if c := s.Category.String; c != "SECURITY" {
		return fmt.Errorf("expected %v to be an Security integration, got %v", id, c)
	}

	// Note: type must be SAML2 or something is broken
	if c := s.IntegrationType.String; c != "SAML2" {
		return fmt.Errorf("expected %v to be a SAML2 integration type, got %v", id, c)
	}

	if err := d.Set("name", s.Name.String); err != nil {
		return err
	}

	if err := d.Set("created_on", s.CreatedOn.String); err != nil {
		return err
	}

	if err := d.Set("enabled", s.Enabled.Bool); err != nil {
		return err
	}

	// Some properties come from the DESCRIBE INTEGRATION call
	// We need to grab them in a loop
	var k, pType string
	var v, unused interface{}
	stmt = snowflake.SamlIntegration(id).Describe()
	rows, err := db.Query(stmt)
	if err != nil {
		return errors.Wrap(err, "could not describe security integration")
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&k, &pType, &v, &unused); err != nil {
			return errors.Wrap(err, "unable to parse security integration rows")
		}
		switch k {
		case "ENABLED":
			// set using the SHOW INTEGRATION, ignoring here
		case "SAML2_ISSUER":
			if err = d.Set("saml2_issuer", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_issuer for security integration")
			}
		case "SAML2_SSO_URL":
			if err = d.Set("saml2_sso_url", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_sso_url for security integration")
			}
		case "SAML2_PROVIDER":
			if err = d.Set("saml2_provider", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_provider for security integration")
			}
		case "SAML2_X509_CERT":
			if err = d.Set("saml2_x509_cert", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_x509_cert for security integration")
			}
		case "SAML2_SP_INITIATED_LOGIN_PAGE_LABEL":
			if err = d.Set("saml2_sp_initiated_login_page_label", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_sp_initiated_login_page_label for security integration")
			}
		case "SAML2_ENABLE_SP_INITIATED":
			b := false
			switch v2 := v.(type) {
			case bool:
				b = v2
			case string:
				b, err = strconv.ParseBool(v.(string))
				if err != nil {
					return errors.Wrap(err, "returned saml2_force_authn that is not boolean")
				}
			default:
				return errors.Wrap(err, "returned saml2_force_authn that is not boolean")
			}
			if err = d.Set("saml2_enable_sp_initiated", b); err != nil {
				return errors.Wrap(err, "unable to set saml2_enable_sp_initiated for security integration")
			}
		case "SAML2_SNOWFLAKE_X509_CERT":
			if err = d.Set("saml2_snowflake_x509_cert", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_snowflake_x509_cert for security integration")
			}
		case "SAML2_SIGN_REQUEST":
			b := false
			switch v2 := v.(type) {
			case bool:
				b = v2
			case string:
				b, err = strconv.ParseBool(v.(string))
				if err != nil {
					return errors.Wrap(err, "returned saml2_force_authn that is not boolean")
				}
			default:
				return errors.Wrap(err, "returned saml2_force_authn that is not boolean")
			}
			if err = d.Set("saml2_sign_request", b); err != nil {
				return errors.Wrap(err, "unable to set saml2_sign_request for security integration")
			}
		case "SAML2_REQUESTED_NAMEID_FORMAT":
			if err = d.Set("saml2_requested_nameid_format", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_requested_nameid_format for security integration")
			}
		case "SAML2_POST_LOGOUT_REDIRECT_URL":
			if err = d.Set("saml2_post_logout_redirect_url", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_post_logout_redirect_url for security integration")
			}
		case "SAML2_FORCE_AUTHN":
			b := false
			switch v2 := v.(type) {
			case bool:
				b = v2
			case string:
				b, err = strconv.ParseBool(v.(string))
				if err != nil {
					return errors.Wrap(err, "returned saml2_force_authn that is not boolean")
				}
			default:
				return errors.Wrap(err, "returned saml2_force_authn that is not boolean")
			}
			if err = d.Set("saml2_force_authn", b); err != nil {
				return errors.Wrap(err, "unable to set saml2_force_authn for security integration")
			}
		case "SAML2_SNOWFLAKE_ISSUER_URL":
			if err = d.Set("saml2_snowflake_issuer_url", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_snowflake_issuer_url for security integration")
			}
		case "SAML2_SNOWFLAKE_ACS_URL":
			if err = d.Set("saml2_snowflake_acs_url", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_snowflake_acs_url for security integration")
			}
		case "SAML2_SNOWFLAKE_METADATA":
			if err = d.Set("saml2_snowflake_metadata", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_snowflake_metadata for security integration")
			}
		case "SAML2_DIGEST_METHODS_USED":
			if err = d.Set("saml2_digest_methods_used", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_digest_methods_used for security integration")
			}
		case "SAML2_SIGNATURE_METHODS_USED":
			if err = d.Set("saml2_signature_methods_used", v.(string)); err != nil {
				return errors.Wrap(err, "unable to set saml2_signature_methods_used for security integration")
			}
		case "COMMENT":
			// COMMENT cannot be set according to snowflake docs, so ignoring
		default:
			log.Printf("[WARN] unexpected security integration property %v returned from Snowflake", k)
		}
	}

	return err
}

// UpdateSAMLIntegration implements schema.UpdateFunc
func UpdateSAMLIntegration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	stmt := snowflake.SamlIntegration(id).Alter()

	var runSetStatement bool

	if d.HasChange("enabled") {
		runSetStatement = true
		stmt.SetBool(`ENABLED`, d.Get("enabled").(bool))
	}

	if d.HasChange("saml2_issuer") {
		runSetStatement = true
		stmt.SetString(`SAML2_ISSUER`, d.Get("saml2_issuer").(string))
	}

	if d.HasChange("saml2_sso_url") {
		runSetStatement = true
		stmt.SetString(`saml2_sso_url`, d.Get("saml2_sso_url").(string))
	}

	if d.HasChange("saml2_provider") {
		runSetStatement = true
		stmt.SetString(`saml2_provider`, d.Get("saml2_provider").(string))
	}

	if d.HasChange("saml2_x509_cert") {
		runSetStatement = true
		stmt.SetString(`saml2_x509_cert`, d.Get("saml2_x509_cert").(string))
	}

	if d.HasChange("saml2_sp_initiated_login_page_label") {
		runSetStatement = true
		stmt.SetString(`saml2_sp_initiated_login_page_label`, d.Get("saml2_sp_initiated_login_page_label").(string))
	}

	if d.HasChange("saml2_enable_sp_initiated") {
		runSetStatement = true
		stmt.SetBool(`saml2_enable_sp_initiated`, d.Get("saml2_enable_sp_initiated").(bool))
	}

	if d.HasChange("saml2_snowflake_x509_cert") {
		runSetStatement = true
		stmt.SetString(`saml2_snowflake_x509_cert`, d.Get("saml2_snowflake_x509_cert").(string))
	}

	if d.HasChange("saml2_sign_request") {
		runSetStatement = true
		stmt.SetBool(`saml2_sign_request`, d.Get("saml2_sign_request").(bool))
	}

	if d.HasChange("saml2_requested_nameid_format") {
		runSetStatement = true
		stmt.SetString(`saml2_requested_nameid_format`, d.Get("saml2_requested_nameid_format").(string))
	}

	if d.HasChange("saml2_post_logout_redirect_url") {
		runSetStatement = true
		stmt.SetString(`saml2_post_logout_redirect_url`, d.Get("saml2_post_logout_redirect_url").(string))
	}

	if d.HasChange("saml2_force_authn") {
		runSetStatement = true
		stmt.SetBool(`saml2_force_authn`, d.Get("saml2_force_authn").(bool))
	}

	if d.HasChange("saml2_snowflake_issuer_url") {
		runSetStatement = true
		stmt.SetString(`saml2_snowflake_issuer_url`, d.Get("saml2_snowflake_issuer_url").(string))
	}

	if d.HasChange("saml2_snowflake_acs_url") {
		runSetStatement = true
		stmt.SetString(`saml2_snowflake_acs_url`, d.Get("saml2_snowflake_acs_url").(string))
	}

	if runSetStatement {
		if err := snowflake.Exec(db, stmt.Statement()); err != nil {
			return errors.Wrap(err, "error updating security integration")
		}
	}

	return ReadSAMLIntegration(d, meta)
}

// DeleteSAMLIntegration implements schema.DeleteFunc
func DeleteSAMLIntegration(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("", snowflake.SamlIntegration)(d, meta)
}
