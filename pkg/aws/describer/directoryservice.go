package describer

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directoryservice"
	"gitlab.com/keibiengine/keibi-engine/pkg/aws/model"
)

func DirectoryServiceDirectory(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := directoryservice.NewFromConfig(cfg)
	paginator := directoryservice.NewDescribeDirectoriesPaginator(client, &directoryservice.DescribeDirectoriesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			if !isErr(err, "InvalidParameterValueException") && !isErr(err, "ResourceNotFoundFault") && !isErr(err, "EntityDoesNotExistException") {
				return nil, err
			}
			continue
		}

		for _, v := range page.DirectoryDescriptions {
			arn := fmt.Sprintf("arn:%s:ds:%s:%s:directory/%s", describeCtx.Partition, describeCtx.Region, describeCtx.AccountID, *v.DirectoryId)

			tags, err := client.ListTagsForResource(ctx, &directoryservice.ListTagsForResourceInput{
				ResourceId: v.DirectoryId,
			})
			if err != nil {
				if !isErr(err, "InvalidParameterValueException") && !isErr(err, "ResourceNotFoundFault") && !isErr(err, "EntityDoesNotExistException") {
					return nil, err
				}
				tags = &directoryservice.ListTagsForResourceOutput{}
			}

			values = append(values, Resource{
				ARN:  arn,
				Name: *v.Name,
				Description: model.DirectoryServiceDirectoryDescription{
					Directory: v,
					Tags:      tags.Tags,
				},
			})
		}
	}

	return values, nil
}
