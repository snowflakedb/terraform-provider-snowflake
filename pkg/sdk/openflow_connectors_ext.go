package sdk

import "fmt"

// OpenflowConnectorVersionLocation is the Location of a connector version's configuration stage, which
// Snowflake addresses with a snow:// URI rather than a stage path, so it cannot be a StageLocation.
//
// It takes the connector and version rather than a finished URI because the URI quotes only the identifier
// parts that need it, and an unquoted `my_connector` addresses MY_CONNECTOR instead - a config that looks
// right and resolves to nothing until apply.
type OpenflowConnectorVersionLocation struct {
	connector SchemaObjectIdentifier
	version   string
}

// NewOpenflowConnectorVersionLocation builds the location of one version of a connector's configuration.
// version is the name SHOW VERSIONS reports, or `live` for the uncommitted live version.
func NewOpenflowConnectorVersionLocation(connector SchemaObjectIdentifier, version string) OpenflowConnectorVersionLocation {
	return OpenflowConnectorVersionLocation{connector: connector, version: version}
}

// ToSql renders the URI as Snowflake reports it in location_uri. The version is not an identifier so it is
// not quoted; the whole URI is emitted as a single-quoted literal, which escapes it.
func (l OpenflowConnectorVersionLocation) ToSql() string {
	return fmt.Sprintf(
		"snow://openflow_connector/%s.%s.%s/versions/%s/",
		quoteIdentifierPartIfNeeded(l.connector.DatabaseName()),
		quoteIdentifierPartIfNeeded(l.connector.SchemaName()),
		quoteIdentifierPartIfNeeded(l.connector.Name()),
		l.version,
	)
}

func (r *CreateOpenflowConnectorRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// additionalConvert is required because OpenflowConnectorDetails carries a plain-only Id field. There is
// nothing to derive from the row: Id is set by the caller after Describe returns.
func (r openflowConnectorDetailsRow) additionalConvert(_ *OpenflowConnectorDetails) error {
	return nil
}

// ID returns the identifier threaded in by the caller. DESCRIBE OPENFLOW CONNECTOR does not return
// database_name or schema_name, so it cannot be derived from the response.
func (d *OpenflowConnectorDetails) ID() SchemaObjectIdentifier {
	return d.Id
}
