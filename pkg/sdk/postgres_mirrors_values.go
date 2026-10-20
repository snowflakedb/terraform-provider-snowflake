package sdk

// PostgresMirrorNameMaxLength is enforced by the procedures because the name also becomes the Postgres publication and
// replication slot name. Test name generation stays inside it.
const PostgresMirrorNameMaxLength = 50

// PostgresMirrorDefaultRefreshInterval is what CREATE_MIRROR uses when refresh_interval is omitted.
const PostgresMirrorDefaultRefreshInterval = "10 minutes"
