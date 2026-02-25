package sdk

import "context"

// TrustCenter provides an interface for managing Trust Center scanner configurations.
// Trust Center uses stored procedures rather than DDL commands, so this implementation
// builds SQL CALL statements directly.
type TrustCenter interface {
	// Scanner Package Configuration (Account Level)
	SetPackageConfiguration(ctx context.Context, req *SetPackageConfigurationRequest) error
	UnsetPackageConfiguration(ctx context.Context, req *UnsetPackageConfigurationRequest) error
	ShowScannerPackages(ctx context.Context, req *ShowScannerPackagesRequest) ([]ScannerPackage, error)
	ShowScannerPackageByID(ctx context.Context, id string) (*ScannerPackage, error)

	// Scanner Configuration (Account Level)
	SetScannerConfiguration(ctx context.Context, req *SetScannerConfigurationRequest) error
	UnsetScannerConfiguration(ctx context.Context, req *UnsetScannerConfigurationRequest) error
	ShowScanners(ctx context.Context, req *ShowScannersRequest) ([]Scanner, error)
	ShowScannerByID(ctx context.Context, packageId, scannerId string) (*Scanner, error)
}

var _ TrustCenter = (*trustCenter)(nil)

type trustCenter struct {
	client *Client
}
