package opengovernance_client

type ResourceCostEstimate struct {
	ResourceID   string  `json:"resource_id"`
	ResourceType string  `json:"resource_type"`
	Cost         float64 `json:"cost"`
}

func ResourceTypeConversion(resourceType string) string {
	//TODO
	switch resourceType {
	case "aws::elasticloadbalancing::loadbalancer":
		return "aws_lb"
	case "aws::ec2::volumesnapshot":
		return "aws_ebs_snapshot"
	case "aws::elasticloadbalancingv2::loadbalancer":
		return "aws_alb"
	case "aws::ec2::instance":
		return "aws_instance"
	case "aws::autoscaling::autoscalinggroup":
		return "aws_autoscaling_group"
	case "aws::rds::dbinstance":
		return "aws_db_instance"
	case "aws::ec2::volume":
		return "aws_ebs_volume"
	case "aws::ec2::volumegp3":
		return "aws_ebs_volume"
	case "aws::efs::filesystem":
		return "aws_efs_file_system"
	case "aws::elasticache::cluster":
		return "aws_elasticache_cluster"
	case "aws::elasticache::replicationgroup":
		return "aws_elasticache_replication_group"
	case "aws::ec2::eip":
		return "aws_eip"
	case "aws::eks::cluster":
		return "aws_eks_cluster"
	case "aws::eks::nodegroup":
		return "aws_eks_node_group"
	case "aws::fsx::filesystem":
		return "aws_fsx_ontap_file_system"
	case "aws::ec2::natgateway":
		return "aws_nat_gateway"
	case "aws::ec2::host":
		return "aws_ec2_host"
	case "aws::lambda::function":
		return "aws_lambda_function"
	case "aws::elasticsearch::domain":
		return "aws_elasticsearch_domain"
	case "aws::opensearch::domain":
		return "aws_opensearch_domain"
	case "aws::dynamodb::table":
		return "aws_dynamodb_table"

	// Azure
	case "microsoft.compute/virtualmachines":
		return "azurerm_virtual_machine"
	case "microsoft.compute/disks":
		return "azurerm_managed_disk"
	case "microsoft.compute/images":
		return "azurerm_image"
	case "microsoft.compute/snapshots":
		return "azurerm_snapshot"
	case "microsoft.compute/virtualmachinescalesets":
		return "azurerm_virtual_machine_scale_set"
	case "microsoft.network/loadbalancers":
		return "azurerm_lb"
	case "microsoft.network/loadbalancers/loadbalancingeules":
		return "azurerm_lb_rule"
	case "microsoft.network/loadbalancers/outboundrules":
		return "azurerm_lb_outbound_rule"
	case "microsoft.network/applicationgateways":
		return "azurerm_application_gateway"
	case "microsoft.network/natgateways":
		return "azurerm_nat_gateway"
	case "microsoft.network/publicipaddresses":
		return "azurerm_public_ip"
	case "microsoft.network/publicipprefixes":
		return "azurerm_public_ip_prefix"
	case "microsoft.containerregistry/registries":
		return "azurerm_container_registry"
	case "microsoft.network/privateendpoints":
		return "azurerm_private_endpoint"
	case "microsoft.storage/queues":
		return "azurerm_storage_queue"
	case "microsoft.storage/fileshares":
		return "azurerm_storage_share"
	case "microsoft.storage/storageaccounts":
		return "azurerm_storage_account"
	case "microsoft.network/virtualnetworkgateways":
		return "azurerm_virtual_network_gateway"
	case "microsoft.keyvault/vaults/keys":
		return "azurerm_key_vault_key"
	case "microsoft.keyvault/managedhsms":
		return "azurerm_key_vault_managed_hardware_security_module"
	case "microsoft.cdn/profiles/endpoints":
		return "azurerm_cdn_endpoint"
	case "microsoft.network/dnszones":
		return "azurerm_dns_zone"
	case "microsoft.network/privatednszones":
		return "azurerm_private_dns_zone"
	case "microsoft.documentdb/sqldatabases":
		return "azurerm_cosmosdb_sql_database"
	case "microsoft.documentdb/mongodatabases":
		return "azurerm_cosmosdb_mongo_database"
	case "microsoft.documentdb/mongocollection":
		return "azurerm_cosmosdb_mongo_collection"
	case "microsoft.dbformariadb/servers":
		return "azurerm_mariadb_server"
	case "microsoft.sql/servers/databases":
		return "azurerm_sql_database"
	case "microsoft.sql/managedInstances":
		return "azurerm_sql_managed_instance"
	case "microsoft.dbformysql/servers":
		return "azurerm_mysql_server"
	case "microsoft.dbforpostgresql/servers":
		return "azurerm_postgresql_server"
	case "microsoft.dbforpostgresql/flexibleservers":
		return "azurerm_postgresql_flexible_server"
	case "microsoft.dbformysql/flexibleservers":
		return "azurerm_mysql_flexible_server"
	case "microsoft.containerservice/managedclusters":
		return "azurerm_kubernetes_cluster"
	case "microsoft.web/hostingenvironments":
		return "azurerm_app_service_environment"
	case "microsoft.web/plan":
		return "azurerm_app_service_plan"
	case "microsoft.apimanagement/service":
		return "azurerm_api_management"
	case "microsoft.web/sites":
		return "azurerm_function_app"
	case "microsoft.search/searchservices":
		return "azurerm_search_service"
	case "microsoft.automation/automationaccounts":
		return "azurerm_automation_account"
	}
	return resourceType
}
