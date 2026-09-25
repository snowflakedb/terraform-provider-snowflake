package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStageDescribeProjections(t *testing.T) {
	id := NewSchemaObjectIdentifier("DB", "SCH", "STG")
	formatName := NewSchemaObjectIdentifier("DB", "SCH", "FF")
	d := &StageDetails{
		Id:             id,
		FileFormatName: &formatName,
		DirectoryTable: &StageDirectoryTable{Enable: true, AutoRefresh: false},
		PrivateLink:    &StagePrivateLink{UsePrivatelinkEndpoint: true},
		Location: &StageLocationDetails{
			Url:               []string{"s3://bucket/path"},
			AwsAccessPointArn: "arn:aws:s3:us-east-1:123:accesspoint/ap",
		},
		Credentials: &StageCredentials{AwsKeyId: "should-not-project"},
	}

	common := d.AsCommon()
	require.Equal(t, d.FileFormatName, common.FileFormatName)
	require.Equal(t, d.DirectoryTable, common.DirectoryTable)
	require.Nil(t, (*StageDetails)(nil).AsCommon())

	aws := d.AsAws()
	require.Equal(t, d.FileFormatName, aws.FileFormatName)
	require.Equal(t, d.PrivateLink, aws.PrivateLink)
	require.Equal(t, d.Location, aws.Location)
	require.Nil(t, (*StageDetails)(nil).AsAws())

	compatible := d.AsAwsCompatible()
	require.Equal(t, d.FileFormatName, compatible.FileFormatName)
	require.Equal(t, d.Location, compatible.Location)
	require.Nil(t, (*StageDetails)(nil).AsAwsCompatible())
}
