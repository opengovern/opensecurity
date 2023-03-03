package onboard

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/subscription/mgmt/subscription"
	"gitlab.com/keibiengine/keibi-engine/pkg/azure"
	"gitlab.com/keibiengine/keibi-engine/pkg/onboard/api"
	"gitlab.com/keibiengine/keibi-engine/pkg/source"
)

func discoverAzureSubscriptions(ctx context.Context, req api.DiscoverAzureSubscriptionsRequest) ([]api.DiscoverAzureSubscriptionsResponse, error) {
	authorizer, err := azure.NewAuthorizerFromConfig(azure.AuthConfig{
		TenantID:     req.TenantId,
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
	})
	if err != nil {
		return nil, err
	}

	client := subscription.NewSubscriptionsClient()
	client.Authorizer = authorizer

	authorizer.WithAuthorization()

	it, err := client.ListComplete(ctx)
	if err != nil {
		return nil, err
	}

	subs := []api.DiscoverAzureSubscriptionsResponse{} // don't convert to var so the returned list won't be nil
	for it.NotDone() {
		v := it.Value()
		subs = append(subs, api.DiscoverAzureSubscriptionsResponse{
			ID:             *v.ID,
			SubscriptionID: *v.SubscriptionID,
			Name:           *v.DisplayName,
			Status:         string(v.State),
		})

		if it.NotDone() {
			err := it.NextWithContext(ctx)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return subs, nil
}

func getAzureCredentialsMetadata(ctx context.Context, config api.SourceConfigAzure) (*source.AzureCredentialMetadata, error) {
	//TODO get metadata

	return &source.AzureCredentialMetadata{}, nil
}
