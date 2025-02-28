package opengovernance_client

//
//import (
//	aws "github.com/opengovern/og-aws-describer/aws/model"
//	azure "github.com/opengovern/og-azure-describer/azure/model"
//	opengovernanceAzure "github.com/opengovern/og-azure-describer/pkg/opengovernance-es-sdk"
//	"strings"
//)
//
//// getAwsEc2HostValues get resource values needed for cost estimate from model.EC2HostDescription
//func getAwsEc2HostValues(resource Resource) (map[string]interface{}, error) {
//	var valuesMap map[string]interface{}
//	if v, ok := resource.Description.(opengovernanceAws.EC2Host); ok {
//		valuesMap["availability_zone"] = ""
//		if v.Description.Host.AvailabilityZone != nil {
//			valuesMap["availability_zone"] = *v.Description.Host.AvailabilityZone
//		}
//		valuesMap["instance_type"] = ""
//		valuesMap["instance_family"] = ""
//		if v.Description.Host.HostProperties != nil {
//			if v.Description.Host.HostProperties.InstanceType != nil {
//				valuesMap["instance_type"] = *v.Description.Host.HostProperties.InstanceType
//			}
//			if v.Description.Host.HostProperties.InstanceFamily != nil {
//				valuesMap["instance_family"] = *v.Description.Host.HostProperties.InstanceFamily
//			}
//		}
//	}
//	return valuesMap, nil
//}
//
//// getAwsLambdaFunctionValues get resource values needed for cost estimate from model.LambdaFunctionDescription
//func getAwsLambdaFunctionValues(resource Resource) (map[string]interface{}, error) {
//	var valuesMap map[string]interface{}
//	if v, ok := resource.Description.(opengovernanceAws.LambdaFunction); ok {
//		if v.Description.Function == nil {
//			return nil, nil
//		}
//		if v.Description.Function.Configuration != nil {
//			valuesMap["memory_size"] = ptrInt(v.Description.Function.Configuration.MemorySize)
//			valuesMap["architectures"] = v.Description.Function.Configuration.Architectures
//			if v.Description.Function.Configuration.EphemeralStorage != nil {
//				valuesMap["ephemeral_storage"] = []map[string]interface{}{{"size": ptrInt(v.Description.Function.Configuration.EphemeralStorage.Size)}}
//			}
//		}
//	}
//	return valuesMap, nil
//}
//
//// getAwsEsDomainValues get resource values needed for cost estimate from model.ESDomainDescription
//func getAwsEsDomainValues(resource Resource) (map[string]interface{}, error) {
//	var valuesMap map[string]interface{}
//	if v, ok := resource.Description.(opengovernanceAws.ESDomain); ok {
//		if v.Description.Domain.ElasticsearchClusterConfig != nil {
//			valuesMap["cluster_config"] = []map[string]interface{}{
//				{
//					"instance_type":            v.Description.Domain.ElasticsearchClusterConfig.InstanceType,
//					"instance_count":           ptrInt(v.Description.Domain.ElasticsearchClusterConfig.InstanceCount),
//					"dedicated_master_enabled": ptrBool(v.Description.Domain.ElasticsearchClusterConfig.DedicatedMasterEnabled),
//					"dedicated_master_type":    v.Description.Domain.ElasticsearchClusterConfig.DedicatedMasterType,
//					"dedicated_master_count":   ptrInt(v.Description.Domain.ElasticsearchClusterConfig.DedicatedMasterCount),
//					"warm_enabled":             ptrBool(v.Description.Domain.ElasticsearchClusterConfig.WarmEnabled),
//					"warm_type":                v.Description.Domain.ElasticsearchClusterConfig.WarmType,
//					"warm_count":               ptrInt(v.Description.Domain.ElasticsearchClusterConfig.WarmCount),
//				},
//			}
//		}
//
//		if v.Description.Domain.EBSOptions != nil {
//			valuesMap["ebs_options"] = []map[string]interface{}{
//				{
//					"ebs_enabled": ptrBool(v.Description.Domain.EBSOptions.EBSEnabled),
//					"volume_type": v.Description.Domain.EBSOptions.VolumeType,
//					"volume_size": ptrInt(v.Description.Domain.EBSOptions.VolumeSize),
//					"iops":        ptrInt(v.Description.Domain.EBSOptions.Iops),
//					"throughput":  ptrInt(v.Description.Domain.EBSOptions.Throughput),
//				},
//			}
//		}
//	}
//	return valuesMap, nil
//}
//
//// getAwsOpenSearchDomainValues get resource values needed for cost estimate from model.OpenSearchDomainDescription
//func getAwsOpenSearchDomainValues(resource Resource) (map[string]interface{}, error) {
//	var valuesMap map[string]interface{}
//
//	if v, ok := resource.Description.(opengovernanceAws.OpenSearchDomain); ok {
//
//		if v.Description.Domain.ClusterConfig != nil {
//			valuesMap["cluster_config"] = []map[string]interface{}{
//				{
//					"instance_type":            v.Description.Domain.ClusterConfig.InstanceType,
//					"instance_count":           ptrInt(v.Description.Domain.ClusterConfig.InstanceCount),
//					"dedicated_master_enabled": ptrBool(v.Description.Domain.ClusterConfig.DedicatedMasterEnabled),
//					"dedicated_master_type":    v.Description.Domain.ClusterConfig.DedicatedMasterType,
//					"dedicated_master_count":   ptrInt(v.Description.Domain.ClusterConfig.DedicatedMasterCount),
//					"warm_enabled":             ptrBool(v.Description.Domain.ClusterConfig.WarmEnabled),
//					"warm_type":                v.Description.Domain.ClusterConfig.WarmType,
//					"warm_count":               ptrInt(v.Description.Domain.ClusterConfig.WarmCount),
//				},
//			}
//		}
//
//		if v.Description.Domain.EBSOptions != nil {
//			valuesMap["ebs_options"] = []map[string]interface{}{
//				{
//					"ebs_enabled": ptrBool(v.Description.Domain.EBSOptions.EBSEnabled),
//					"volume_type": v.Description.Domain.EBSOptions.VolumeType,
//					"volume_size": ptrInt(v.Description.Domain.EBSOptions.VolumeSize),
//					"iops":        ptrInt(v.Description.Domain.EBSOptions.Iops),
//					"throughput":  ptrInt(v.Description.Domain.EBSOptions.Throughput),
//				},
//			}
//		}
//	}
//
//	return valuesMap, nil
//}
//
//// getAwsNatGatewayValues get resource values needed for cost estimate
//func getAwsNatGatewayValues(resource Resource) (map[string]interface{}, error) {
//	return map[string]interface{}{}, nil
//}
//
//// getAwsFSXFileSystemValues get resource values needed for cost estimate from model.FSXFileSystemDescription
//func getAwsFSXFileSystemValues(resource Resource) (map[string]interface{}, error) {
//	var valuesMap map[string]interface{}
//	if v, ok := resource.Description.(opengovernanceAws.FSXFileSystem); ok {
//		valuesMap["storage_capacity"] = ptrInt2(v.Description.FileSystem.StorageCapacity)
//		valuesMap["storage_type"] = v.Description.FileSystem.StorageType
//		if v.Description.FileSystem.LustreConfiguration != nil {
//			valuesMap["deployment_type"] = v.Description.FileSystem.LustreConfiguration.DeploymentType
//			valuesMap["data_compression_type"] = v.Description.FileSystem.LustreConfiguration.DataCompressionType
//		}
//		if v.Description.FileSystem.OpenZFSConfiguration != nil {
//			valuesMap["throughput_capacity"] = ptrInt2(v.Description.FileSystem.OpenZFSConfiguration.ThroughputCapacity)
//			valuesMap["automatic_backup_retention_days"] = ptrInt2(v.Description.FileSystem.OpenZFSConfiguration.AutomaticBackupRetentionDays)
//		}
//		if v.Description.FileSystem.OntapConfiguration != nil {
//			if v.Description.FileSystem.OntapConfiguration.DiskIopsConfiguration != nil {
//				valuesMap["disk_iops_configuration"] = []map[string]interface{}{{
//					"iops": ptrInt642(v.Description.FileSystem.OntapConfiguration.DiskIopsConfiguration.Iops),
//					"mode": v.Description.FileSystem.OntapConfiguration.DiskIopsConfiguration.Mode,
//				}}
//			}
//		}
//	}
//	return valuesMap, nil
//}
//
//// getAwsEksNodeGroupValues get resource values needed for cost estimate from model.EKSNodegroupDescription
//func getAwsEksNodeGroupValues(resource Resource) (map[string]interface{}, error) {
//	var valuesMap map[string]interface{}
//	if v, ok := resource.Description.(opengovernanceAws.EKSNodegroup); ok {
//		if v.Description.Nodegroup.ScalingConfig != nil {
//			valuesMap["scaling_config"] = []map[string]interface{}{
//				{
//					"min_size":     ptrInt2(v.Description.Nodegroup.ScalingConfig.MinSize),
//					"desired_size": ptrInt2(v.Description.Nodegroup.ScalingConfig.DesiredSize),
//				},
//			}
//		}
//		valuesMap["instance_types"] = v.Description.Nodegroup.InstanceTypes
//		valuesMap["disk_size"] = ptrInt2(v.Description.Nodegroup.DiskSize)
//		if v.Description.Nodegroup.LaunchTemplate != nil {
//			valuesMap["launch_template"] = []map[string]interface{}{
//				{
//					"id":      ptrStr2(v.Description.Nodegroup.LaunchTemplate.Id),
//					"Name":    ptrStr2(v.Description.Nodegroup.LaunchTemplate.Name),
//					"Version": ptrStr2(v.Description.Nodegroup.LaunchTemplate.Version),
//				},
//			}
//		}
//	}
//	return valuesMap, nil
//}
//
//// getAwsEksNodeGroupValues get resource values needed for cost estimate from model.EKSNodegroupDescription
//func getAwsEksClusterValues(resource Resource) (map[string]interface{}, error) {
//	return map[string]interface{}{}, nil
//}
//
//// getAwsEc2EipValues get resource values needed for cost estimate from model.EC2EIPDescription
//func getAwsEc2EipValues(resource Resource) (map[string]interface{}, error) {
//	if v, ok := resource.Description.(aws.EC2EIPDescription); ok {
//		return map[string]interface{}{
//			"customer_owned_ipv4_pool": ptrStr2(v.Address.CustomerOwnedIpv4Pool),
//			"instance":                 ptrStr2(v.Address.InstanceId),
//			"network_interface":        ptrStr2(v.Address.NetworkInterfaceId),
//		}, nil
//	} else if v, ok := resource.Description.(map[string]interface{}); ok {
//		return map[string]interface{}{
//			"customer_owned_ipv4_pool": v["Address"].(map[string]interface{})["CustomerOwnedIpv4Pool"],
//			"instance":                 v["Address"].(map[string]interface{})["InstanceId"],
//			"network_interface":        v["Address"].(map[string]interface{})["NetworkInterfaceId"],
//		}, nil
//	}
//	return resource.Description.(map[string]interface{}), nil
//}
//
//// getAwsElastiCacheReplicationGroupValues get resource values needed for cost estimate from model.ElastiCacheReplicationGroupDescription
//func getAwsElastiCacheReplicationGroupValues(resource Resource) (map[string]interface{}, error) {
//	var valuesMap map[string]interface{}
//	if v, ok := resource.Description.(opengovernanceAws.ElastiCacheReplicationGroup); ok {
//		valuesMap["node_type"] = ptrStr2(v.Description.ReplicationGroup.CacheNodeType)
//		valuesMap["engine"] = "" // TODO find engine
//		valuesMap["num_node_groups"] = len(v.Description.ReplicationGroup.NodeGroups)
//		var replicas int
//		for _, ng := range v.Description.ReplicationGroup.NodeGroups {
//			replicas += len(ng.NodeGroupMembers)
//		}
//		valuesMap["replicas_per_node_group"] = replicas / len(v.Description.ReplicationGroup.NodeGroups)
//		valuesMap["num_cache_clusters"] = len(v.Description.ReplicationGroup.MemberClusters)
//		valuesMap["snapshot_retention_limit"] = ptrInt2(v.Description.ReplicationGroup.SnapshotRetentionLimit)
//
//		if v.Description.ReplicationGroup.GlobalReplicationGroupInfo != nil {
//			valuesMap["global_replication_group_id"] = ptrStr2(v.Description.ReplicationGroup.GlobalReplicationGroupInfo.GlobalReplicationGroupId)
//		}
//	}
//	return valuesMap, nil
//}
//
//// getAwsElastiCacheClusterValues get resource values needed for cost estimate from model.ElastiCacheClusterDescription
//func getAwsElastiCacheClusterValues(resource Resource) (map[string]interface{}, error) {
//	var valuesMap map[string]interface{}
//	if v, ok := resource.Description.(opengovernanceAws.ElastiCacheCluster); ok {
//		valuesMap["node_type"] = ptrStr2(v.Description.Cluster.CacheNodeType)
//		valuesMap["availability_zone"] = ptrStr2(v.Description.Cluster.PreferredAvailabilityZone)
//		valuesMap["engine"] = ptrStr2(v.Description.Cluster.Engine)
//		valuesMap["replication_group_id"] = ptrStr2(v.Description.Cluster.ReplicationGroupId)
//		valuesMap["num_cache_nodes"] = len(v.Description.Cluster.CacheNodes)
//		valuesMap["snapshot_retention_limit"] = ptrInt2(v.Description.Cluster.SnapshotRetentionLimit)
//	}
//	return valuesMap, nil
//}
//
//// getAwsEfsFileSystemValues get resource values needed for cost estimate from model.EFSFileSystemDescription
//func getAwsEfsFileSystemValues(resource Resource) (map[string]interface{}, error) {
//	var valuesMap map[string]interface{}
//	if v, ok := resource.Description.(opengovernanceAws.EFSFileSystem); ok {
//		valuesMap["availability_zone_name"] = ptrStr2(v.Description.FileSystem.AvailabilityZoneName)
//		//valuesMap["lifecycle_policy"] = nil // TODO
//		valuesMap["throughput_mode"] = v.Description.FileSystem.ThroughputMode
//		valuesMap["provisioned_throughput_in_mibps"] = ptrFloat2(v.Description.FileSystem.ProvisionedThroughputInMibps)
//	}
//	return valuesMap, nil
//}
//
//// getAwsEbsSnapshotValues get resource values needed for cost estimate from model.EC2VolumeSnapshotDescription
//func getAwsEbsSnapshotValues(resource Resource) (map[string]interface{}, error) {
//	if v, ok := resource.Description.(aws.EC2VolumeSnapshotDescription); ok {
//		return map[string]interface{}{
//			"volume_size": v.Snapshot.VolumeSize,
//		}, nil
//	} else if v, ok := resource.Description.(map[string]interface{}); ok {
//		return map[string]interface{}{
//			"volume_size": v["Snapshot"].(map[string]interface{})["VolumeSize"],
//		}, nil
//	}
//	return nil, nil
//}
//
//// getAwsEbsVolumeValues get resource values needed for cost estimate from model.EC2VolumeDescription
//func getAwsEbsVolumeValues(resource Resource) (map[string]interface{}, error) {
//	if v, ok := resource.Description.(aws.EC2VolumeDescription); ok {
//		return map[string]interface{}{
//			"availability_zone": ptrStr2(v.Volume.AvailabilityZone),
//			"type":              v.Volume.VolumeType,
//			"size":              ptrInt2(v.Volume.Size),
//			"iops":              ptrInt2(v.Volume.Iops),
//			"throughput":        ptrInt2(v.Volume.Throughput),
//		}, nil
//	} else if v, ok := resource.Description.(map[string]interface{}); ok {
//		return map[string]interface{}{
//			"availability_zone": v["Volume"].(map[string]interface{})["AvailabilityZone"],
//			"type":              v["Volume"].(map[string]interface{})["VolumeType"],
//			"size":              v["Volume"].(map[string]interface{})["Size"],
//			"iops":              v["Volume"].(map[string]interface{})["Iops"],
//			"throughput":        v["Volume"].(map[string]interface{})["Throughput"],
//		}, nil
//	}
//
//	return resource.Description.(map[string]interface{}), nil
//}
//
//// getAwsEbsVolumeGp3Values get resource values needed for cost estimate from model.EC2VolumeDescription
//func getAwsEbsVolumeGp3Values(resource Resource) (map[string]interface{}, error) {
//	if v, ok := resource.Description.(aws.EC2VolumeDescription); ok {
//		return map[string]interface{}{
//			"availability_zone": ptrStr2(v.Volume.AvailabilityZone),
//			"type":              "gp3",
//			"size":              ptrInt2(v.Volume.Size),
//			"iops":              ptrInt2(v.Volume.Iops),
//			"throughput":        ptrInt2(v.Volume.Throughput),
//		}, nil
//	} else if v, ok := resource.Description.(map[string]interface{}); ok {
//		return map[string]interface{}{
//			"availability_zone": v["Volume"].(map[string]interface{})["AvailabilityZone"],
//			"type":              "gp3",
//			"size":              v["Volume"].(map[string]interface{})["Size"],
//			"iops":              v["Volume"].(map[string]interface{})["Iops"],
//			"throughput":        v["Volume"].(map[string]interface{})["Throughput"],
//		}, nil
//	}
//
//	return resource.Description.(map[string]interface{}), nil
//}
//
//// getAwsRdsDbInstanceValues get resource values needed for cost estimate from model.RDSDBInstanceDescription
//func getAwsRdsDbInstanceValues(resource Resource) (map[string]interface{}, error) {
//	var valuesMap map[string]interface{}
//
//	if v, ok := resource.Description.(opengovernanceAws.RDSDBInstance); ok {
//		valuesMap["instance_class"] = ptrStr2(v.Description.DBInstance.DBInstanceClass)
//		valuesMap["availability_zone"] = ptrStr2(v.Description.DBInstance.AvailabilityZone)
//		valuesMap["engine"] = ptrStr2(v.Description.DBInstance.Engine)
//		valuesMap["license_model"] = ptrStr2(v.Description.DBInstance.LicenseModel)
//		valuesMap["multi_az"] = ptrBool2(v.Description.DBInstance.MultiAZ)
//		valuesMap["allocated_storage"] = ptrInt2(v.Description.DBInstance.AllocatedStorage)
//		valuesMap["storage_type"] = ptrStr2(v.Description.DBInstance.StorageType)
//		valuesMap["iops"] = ptrInt2(v.Description.DBInstance.Iops)
//	}
//
//	return valuesMap, nil
//}
//
//// getAwsAutoscalingGroupValues get resource values needed for cost estimate from model.AutoScalingGroupDescription
//func getAwsAutoscalingGroupValues(resource Resource) (map[string]interface{}, error) {
//	var valuesMap map[string]interface{}
//	if v, ok := resource.Description.(opengovernanceAws.AutoScalingGroup); ok {
//		valuesMap["availability_zones"] = v.Description.AutoScalingGroup.AvailabilityZones
//		valuesMap["launch_configuration"] = v.Description.AutoScalingGroup.LaunchConfigurationName
//	}
//	return valuesMap, nil
//}
//
//// getAwsEc2InstanceValues get resource values needed for cost estimate from model.EC2InstanceDescription
//func getAwsEc2InstanceValues(resource Resource) (map[string]interface{}, error) {
//	var valuesMap map[string]interface{}
//
//	if v, ok := resource.Description.(opengovernanceAws.EC2Instance); ok {
//		valuesMap["instance_type"] = v.Description.Instance.InstanceType
//		if v.Description.Instance.Placement != nil {
//			valuesMap["tenancy"] = v.Description.Instance.Placement.Tenancy
//			valuesMap["availability_zone"] = ptrStr2(v.Description.Instance.Placement.AvailabilityZone)
//		}
//		if v.Description.Instance.Platform == "" {
//			valuesMap["operating_system"] = "Linux"
//		} else {
//			valuesMap["operating_system"] = v.Description.Instance.Platform
//		}
//		valuesMap["ebs_optimized"] = ptrBool2(v.Description.Instance.EbsOptimized)
//		if v.Description.Instance.Monitoring != nil {
//			if v.Description.Instance.Monitoring.State == "disabled" || v.Description.Instance.Monitoring.State == "disabling" {
//				valuesMap["monitoring"] = false
//			} else {
//				valuesMap["monitoring"] = true
//			}
//		}
//		if v.Description.Instance.CpuOptions != nil {
//			valuesMap["credit_specification"] = []map[string]interface{}{{
//				"cpu_credits": ptrInt2(v.Description.Instance.CpuOptions.CoreCount), // TODO not sure
//			}}
//		}
//		// valuesMap["root_block_device"] // TODO must get from another resource
//	}
//	return valuesMap, nil
//}
//
//// getAwsLoadBalancerValues get resource values needed for cost estimate
//func getAwsLoadBalancerValues(resource Resource) (map[string]interface{}, error) {
//	valuesMap := make(map[string]interface{})
//
//	valuesMap["load_balancer_type"] = "classic"
//	valuesMap["region"] = strings.Split(resource.ARN, ":")[3]
//
//	return valuesMap, nil
//}
//
//// getAwsLoadBalancer2Values get resource values needed for cost estimate from model.ElasticLoadBalancingV2LoadBalancerDescription
//func getAwsLoadBalancer2Values(resource Resource) (map[string]interface{}, error) {
//	if v, ok := resource.Description.(aws.ElasticLoadBalancingV2LoadBalancerDescription); ok {
//		return map[string]interface{}{
//			"load_balancer_type": v.LoadBalancer.Type,
//			"region":             strings.Split(*v.LoadBalancer.LoadBalancerArn, ":")[3],
//		}, nil
//	} else if v, ok := resource.Description.(map[string]interface{}); ok {
//		return map[string]interface{}{
//			"load_balancer_type": v["LoadBalancer"].(map[string]interface{})["Type"],
//			"region":             strings.Split(v["LoadBalancer"].(map[string]interface{})["LoadBalancerArn"].(string), ":")[3],
//		}, nil
//	}
//	return nil, nil
//}
//
//// getAwsDynamoDbTableValues get resource values needed for cost estimate from model.ElasticLoadBalancingV2LoadBalancerDescription
//func getAwsDynamoDbTableValues(resource Resource) (map[string]interface{}, error) {
//	var replicas []struct {
//		RegionName string `mapstructure:"region_name"`
//	}
//	if v, ok := resource.Description.(aws.DynamoDbTableDescription); ok {
//		for _, r := range v.Table.Replicas {
//			replicas = append(replicas, struct {
//				RegionName string `mapstructure:"region_name"`
//			}{RegionName: *r.RegionName})
//		}
//		return map[string]interface{}{
//			"billing_mode":   v.Table.BillingModeSummary.BillingMode,
//			"write_capacity": v.Table.ProvisionedThroughput.WriteCapacityUnits,
//			"read_capacity":  v.Table.ProvisionedThroughput.ReadCapacityUnits,
//			"replica":        replicas,
//		}, nil
//	} else if v, ok := resource.Description.(map[string]interface{}); ok {
//		if _, ok2 := v["Table"].(map[string]interface{})["Replicas"].([]interface{}); ok2 {
//			for _, r := range v["Table"].(map[string]interface{})["Replicas"].([]interface{}) {
//				if r == nil {
//					continue
//				}
//				replicas = append(replicas, struct {
//					RegionName string `mapstructure:"region_name"`
//				}{RegionName: r.(map[string]interface{})["RegionName"].(string)})
//			}
//		}
//		output := make(map[string]interface{})
//
//		if v["Table"].(map[string]interface{})["BillingModeSummary"] != nil {
//			output["billing_mode"] = v["Table"].(map[string]interface{})["BillingModeSummary"].(map[string]interface{})["BillingMode"]
//		}
//		if v["Table"].(map[string]interface{})["ProvisionedThroughput"] != nil {
//			output["write_capacity"] = v["Table"].(map[string]interface{})["ProvisionedThroughput"].(map[string]interface{})["WriteCapacityUnits"]
//			output["read_capacity"] = v["Table"].(map[string]interface{})["ProvisionedThroughput"].(map[string]interface{})["ReadCapacityUnits"]
//		}
//
//		output["replica"] = replicas
//		return output, nil
//	}
//	return nil, nil
//}
//
//func getAzureComputeSnapshotValues(resource Resource) (map[string]interface{}, error) {
//	if v, ok := resource.Description.(opengovernanceAzure.ComputeSnapshots); ok {
//		return map[string]interface{}{
//			"disk_size_gb": ptrInt2(v.Description.Snapshot.Properties.DiskSizeGB),
//			"location":     ptrStr2(v.Description.Snapshot.Location),
//		}, nil
//	} else if v, ok := resource.Description.(map[string]interface{}); ok {
//		return map[string]interface{}{
//			"disk_size_gb": v["Snapshot"].(map[string]interface{})["Properties"].(map[string]interface{})["DiskSizeGB"],
//			"location":     v["Snapshot"].(map[string]interface{})["Location"],
//		}, nil
//	}
//	return nil, nil
//}
//
//func getAzureComputeDiskValues(resource Resource) (map[string]interface{}, error) {
//	if v, ok := resource.Description.(opengovernanceAzure.ComputeDisk); ok {
//		return map[string]interface{}{
//			"storage_account_type":       *v.Description.Disk.SKU.Name,
//			"location":                   ptrStr2(v.Description.Disk.Location),
//			"disk_size_gb":               ptrInt2(v.Description.Disk.Properties.DiskSizeGB),
//			"on_demand_bursting_enabled": ptrBool2(v.Description.Disk.Properties.BurstingEnabled),
//			"disk_mbps_read_write":       ptrInt642(v.Description.Disk.Properties.DiskMBpsReadWrite),
//			"disk_iops_read_write":       ptrInt642(v.Description.Disk.Properties.DiskIOPSReadWrite),
//		}, nil
//	} else if v, ok := resource.Description.(map[string]interface{}); ok {
//		return map[string]interface{}{
//			"storage_account_type":       v["Disk"].(map[string]interface{})["SKU"].(map[string]interface{})["Name"],
//			"location":                   v["Disk"].(map[string]interface{})["Location"],
//			"disk_size_gb":               v["Disk"].(map[string]interface{})["Properties"].(map[string]interface{})["DiskSizeGB"],
//			"on_demand_bursting_enabled": v["Disk"].(map[string]interface{})["Properties"].(map[string]interface{})["BurstingEnabled"],
//			"disk_mbps_read_write":       v["Disk"].(map[string]interface{})["Properties"].(map[string]interface{})["DiskMBpsReadWrite"],
//			"disk_iops_read_write":       v["Disk"].(map[string]interface{})["Properties"].(map[string]interface{})["DiskIOPSReadWrite"],
//		}, nil
//	}
//
//	return nil, nil
//}
//
//func getAzureLoadBalancerValues(resource Resource) (map[string]interface{}, error) {
//	if v, ok := resource.Description.(azure.LoadBalancerDescription); ok {
//		return map[string]interface{}{
//			"sku":          *v.LoadBalancer.SKU.Name,
//			"location":     ptrStr2(v.LoadBalancer.Location),
//			"rules_number": len(v.LoadBalancer.Properties.LoadBalancingRules) + len(v.LoadBalancer.Properties.OutboundRules),
//			"sku_tier":     *v.LoadBalancer.SKU.Tier,
//		}, nil
//	} else if v, ok := resource.Description.(map[string]interface{}); ok {
//		rulesNumber := 0
//		if lbRules, ok := v["LoadBalancer"].(map[string]interface{})["Properties"].(map[string]interface{})["LoadBalancingRules"].([]interface{}); ok {
//			rulesNumber += len(lbRules)
//		}
//		if outRules, ok := v["LoadBalancer"].(map[string]interface{})["Properties"].(map[string]interface{})["OutboundRules"].([]interface{}); ok {
//			rulesNumber += len(outRules)
//		}
//		return map[string]interface{}{
//			"sku":          v["LoadBalancer"].(map[string]interface{})["SKU"].(map[string]interface{})["Name"],
//			"location":     v["LoadBalancer"].(map[string]interface{})["Location"],
//			"rules_number": rulesNumber,
//			"sku_tier":     v["LoadBalancer"].(map[string]interface{})["SKU"].(map[string]interface{})["Tier"],
//		}, nil
//	}
//
//	return resource.Description.(map[string]interface{}), nil
//}
//
//func getAzureApplicationGatewayValues(resource Resource) (map[string]interface{}, error) {
//	if v, ok := resource.Description.(azure.ApplicationGatewayDescription); ok {
//		autoscaleConfiguration := struct {
//			MaxCapacity int32 `mapstructure:"max_capacity"`
//			MinCapacity int32 `mapstructure:"min_capacity"`
//		}{
//			MaxCapacity: *v.ApplicationGateway.Properties.AutoscaleConfiguration.MaxCapacity,
//			MinCapacity: *v.ApplicationGateway.Properties.AutoscaleConfiguration.MinCapacity,
//		}
//		return map[string]interface{}{
//			"location":                *v.ApplicationGateway.Location,
//			"sku":                     []interface{}{*v.ApplicationGateway.Properties.SKU},
//			"autoscale_configuration": []interface{}{autoscaleConfiguration},
//		}, nil
//	} else if v, ok := resource.Description.(map[string]interface{}); ok {
//		autoscaleConfiguration := struct {
//			MaxCapacity int64 `mapstructure:"max_capacity"`
//			MinCapacity int64 `mapstructure:"min_capacity"`
//		}{
//			MaxCapacity: v["ApplicationGateway"].(map[string]interface{})["Properties"].(map[string]interface{})["AutoscaleConfiguration"].(map[string]interface{})["MaxCapacity"].(int64),
//			MinCapacity: v["ApplicationGateway"].(map[string]interface{})["Properties"].(map[string]interface{})["AutoscaleConfiguration"].(map[string]interface{})["MinCapacity"].(int64),
//		}
//		return map[string]interface{}{
//			"location":                v["ApplicationGateway"].(map[string]interface{})["Location"],
//			"sku":                     v["ApplicationGateway"].(map[string]interface{})["Properties"].(map[string]interface{})["SKU"],
//			"autoscale_configuration": []interface{}{autoscaleConfiguration},
//		}, nil
//	}
//
//	return resource.Description.(map[string]interface{}), nil
//}
//
////func getAzureVirtualMachineScaleSetValues(resource Resource) (map[string]interface{}, error) {
////	if v, ok := resource.Description.(azure.ComputeVirtualMachineScaleSetDescription); ok {
////		var additionalCapabalities []map[string]bool
////		additionalCapabalities = append(additionalCapabalities, map[string]bool{
////			"ultra_ssd_enabled": *v.VirtualMachineScaleSet.Properties.AdditionalCapabilities.UltraSSDEnabled,
////		})
////		return map[string]interface{}{
////			"size":          *v.VirtualMachineScaleSet.Plan,
////			"location":     *v.VirtualMachineScaleSet.Location,
////			"sku": *v.VirtualMachineScaleSet.SKU,
////			"license_type":     *v.VirtualMachineScaleSet.Properties,
////			"additional_capabilities":     additionalCapabalities,
////			"os_disk":     v.VirtualMachineScaleSet.Properties,
////			"os_profile_windows_config":     *v.LoadBalancer.SKU.Tier,
////			"storage_profile_image_reference":     *v.LoadBalancer.SKU.Tier,
////			"storage_profile_os_disk":     *v.VirtualMachineScaleSetExtensions[0].Properties.VirtualMachineProfile.,
////			"storage_profile_data_disk":     *v.LoadBalancer.SKU.Tier,
////		}, nil
////	} else if v, ok := resource.Description.(map[string]interface{}); ok {
////		return map[string]interface{}{
////			"sku":          v["LoadBalancer"].(map[string]interface{})["SKU"].(map[string]interface{})["Name"],
////			"location":     v["LoadBalancer"].(map[string]interface{})["Location"],
////			"rules_number": len(v["LoadBalancer"].(map[string]interface{})["InboundNatRules"].([]interface{})) + len(v["LoadBalancer"].(map[string]interface{})["LoadBalancingRules"].([]interface{})) + len(v["LoadBalancer"].(map[string]interface{})["OutboundRules"].([]interface{})),
////			"sku_tier":     v["LoadBalancer"].(map[string]interface{})["SKU"].(map[string]interface{})["Tier"],
////		}, nil
////	}
////
////	return resource.Description.(map[string]interface{}), nil
////}
//
//func ptrBool(pointer *bool) interface{} {
//	if pointer == nil {
//		return nil
//	} else {
//		return *pointer
//	}
//}
//
//func ptrInt(pointer *int32) interface{} {
//	if pointer == nil {
//		return nil
//	} else {
//		return *pointer
//	}
//}
//
//func ptrStr(pointer *string) interface{} {
//	if pointer == nil {
//		return nil
//	} else {
//		return *pointer
//	}
//}
//
//func ptrBool2(pointer *bool) bool {
//	if pointer == nil {
//		return false
//	} else {
//		return *pointer
//	}
//}
//
//func ptrInt2(pointer *int32) int32 {
//	if pointer == nil {
//		return 0
//	} else {
//		return *pointer
//	}
//}
//
//func ptrInt642(pointer *int64) int64 {
//	if pointer == nil {
//		return 0
//	} else {
//		return *pointer
//	}
//}
//
//func ptrStr2(pointer *string) string {
//	if pointer == nil {
//		return ""
//	} else {
//		return *pointer
//	}
//}
//
//func ptrFloat2(pointer *float64) float64 {
//	if pointer == nil {
//		return 0
//	} else {
//		return *pointer
//	}
//}
