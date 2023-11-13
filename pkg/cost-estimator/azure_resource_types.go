package cost_estimator

import (
	"fmt"
	azureModel "github.com/kaytu-io/kaytu-azure-describer/azure/model"
	apiAuth "github.com/kaytu-io/kaytu-engine/pkg/auth/api"
	"github.com/kaytu-io/kaytu-engine/pkg/cost-estimator/es"
	"github.com/kaytu-io/kaytu-engine/pkg/internal/httpclient"
	"github.com/kaytu-io/kaytu-engine/pkg/workspace/api"
	"go.uber.org/zap"
)

func GetComputeVirtualMachineCost(h *HttpHandler, _ string, resourceId string) (float64, error) {
	response, err := es.GetElasticsearch(h.logger, h.client, resourceId, "Microsoft.Compute/virtualMachines")
	if err != nil {
		h.logger.Error("failed to get resource", zap.Error(err))
		return 0, fmt.Errorf("failed to get resource")
	}
	if len(response.Hits.Hits) == 0 {
		return 0, fmt.Errorf("no resource found")
	}
	var request api.GetAzureVmRequest
	if vm, ok := response.Hits.Hits[0].Source.Description.(azureModel.ComputeVirtualMachineDescription); ok {
		request = api.GetAzureVmRequest{
			RegionCode: response.Hits.Hits[0].Source.Location,
			VM:         vm,
		}
	} else {
		return 0, fmt.Errorf("cannot parse resource")
	}
	cost, err := h.workspaceClient.GetAzureVm(&httpclient.Context{UserRole: apiAuth.InternalRole}, request)
	if err != nil {
		h.logger.Error("failed in calculating cost", zap.Error(err))
		return 0, err
	}

	return cost, nil
}

func GetManagedStorageCost(h *HttpHandler, _ string, resourceId string) (float64, error) {
	response, err := es.GetElasticsearch(h.logger, h.client, resourceId, "Microsoft.Compute/disks")
	if err != nil {
		h.logger.Error("failed to get resource", zap.Error(err))
		return 0, fmt.Errorf("failed to get resource")
	}
	if len(response.Hits.Hits) == 0 {
		return 0, fmt.Errorf("no resource found")
	}
	var request api.GetAzureManagedStorageRequest
	if storage, ok := response.Hits.Hits[0].Source.Description.(azureModel.ComputeDiskDescription); ok {
		request = api.GetAzureManagedStorageRequest{
			RegionCode:     response.Hits.Hits[0].Source.Location,
			ManagedStorage: storage,
		}
	} else {
		return 0, fmt.Errorf("cannot parse resource")
	}
	cost, err := h.workspaceClient.GetAzureManagedStorage(&httpclient.Context{UserRole: apiAuth.InternalRole}, request)
	if err != nil {
		h.logger.Error("failed in calculating cost", zap.Error(err))
		return 0, err
	}

	return cost, nil
}

func GetLoadBalancerCost(h *HttpHandler, _ string, resourceId string) (float64, error) {
	response, err := es.GetElasticsearch(h.logger, h.client, resourceId, "Microsoft.Network/loadBalancers")
	if err != nil {
		h.logger.Error("failed to get resource", zap.Error(err))
		return 0, fmt.Errorf("failed to get resource")
	}
	if len(response.Hits.Hits) == 0 {
		return 0, fmt.Errorf("no resource found")
	}
	var request api.GetAzureLoadBalancerRequest
	if lb, ok := response.Hits.Hits[0].Source.Description.(azureModel.LoadBalancerDescription); ok {
		request = api.GetAzureLoadBalancerRequest{
			RegionCode:   response.Hits.Hits[0].Source.Location,
			LoadBalancer: lb,
		}
	} else {
		return 0, fmt.Errorf("cannot parse resource")
	}
	cost, err := h.workspaceClient.GetAzureLoadBalancer(&httpclient.Context{UserRole: apiAuth.InternalRole}, request)
	if err != nil {
		h.logger.Error("failed in calculating cost", zap.Error(err))
		return 0, err
	}

	return cost, nil
}

func GetVirtualNetworkCost(h *HttpHandler, _ string, resourceId string) (float64, error) {
	//var resource azureCompute.VirtualNetwork
	//err := h.GetResource("Microsoft.Network/virtualNetworks", resourceId, &resource)
	//if err != nil {
	//	return 0, err
	//}

	return 0, nil
}
