package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudsearch"
	"gitlab.com/keibiengine/keibi-engine/pkg/aws/model"
)

func CloudSearchDomain(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	client := cloudsearch.NewFromConfig(cfg)

	var values []Resource

	output, err := client.ListDomainNames(ctx, &cloudsearch.ListDomainNamesInput{})
	if err != nil {
		return nil, err
	}

	var domainList []string
	for domainName := range output.DomainNames {
		domainList = append(domainList, domainName)
	}

	domains, err := client.DescribeDomains(ctx, &cloudsearch.DescribeDomainsInput{
		DomainNames: domainList,
	})

	for _, domain := range domains.DomainStatusList {
		values = append(values, Resource{
			ARN:  *domain.ARN,
			Name: *domain.DomainName,
			ID:   *domain.DomainId,
			Description: model.CloudSearchDomainDescription{
				DomainStatus: domain,
			},
		})
	}
	return values, nil
}
