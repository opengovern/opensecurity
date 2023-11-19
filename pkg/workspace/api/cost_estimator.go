package api

import (
	aws "github.com/kaytu-io/kaytu-aws-describer/aws/model"
	azure "github.com/kaytu-io/kaytu-azure-describer/azure/model"
)

type GetEC2InstanceCostRequest struct {
	RegionCode string
	Instance   aws.EC2InstanceDescription
}

type GetEC2VolumeCostRequest struct {
	RegionCode string
	Volume     aws.EC2VolumeDescription
}

type GetLBCostRequest struct {
	RegionCode string
	LBType     string
}

type GetRDSInstanceRequest struct {
	RegionCode string
	DBInstance aws.RDSDBInstanceDescription
}

type GetAzureVmRequest struct {
	RegionCode string
	VM         azure.ComputeVirtualMachineDescription
}

type GetAzureManagedStorageRequest struct {
	RegionCode     string
	ManagedStorage azure.ComputeDiskDescription
}

type GetAzureLoadBalancerRequest struct {
	RegionCode       string
	DailyDataProceed *int64 // (GB)
	LoadBalancer     azure.LoadBalancerDescription
}

type GetAzureVirtualNetworkRequest struct {
	RegionCode            string
	PeeringLocations      []string
	MonthlyDataTransferGB *float64
}

type GetAzureVirtualNetworkPeeringRequest struct {
	SourceLocation        string
	DestinationLocation   string
	MonthlyDataTransferGB *float64
}

type GetAzureSqlServersDatabasesRequest struct {
	RegionCode  string
	SqlServerDB azure.SqlDatabaseDescription
	// MonthlyVCoreHours represents a usage param that allows users to define how many hours of usage a serverless sql database instance uses.
	MonthlyVCoreHours int64
	// ExtraDataStorageGB represents a usage cost of additional backup storage used by the sql database.
	ExtraDataStorageGB float64
	// LongTermRetentionStorageGB defines a usage param that allows users to define how many GB of cold storage the database uses.
	// This is storage that can be kept for up to 10 years.
	LongTermRetentionStorageGB int64
	// BackupStorageGB defines a usage param that allows users to define how many GB Point-In-Time Restore (PITR) backup storage the database uses.
	BackupStorageGB int64
	ResourceId      string
}
