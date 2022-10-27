package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	"gitlab.com/keibiengine/keibi-engine/pkg/aws/model"
)

func DynamoDbTable(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	client := dynamodb.NewFromConfig(cfg)
	paginator := dynamodb.NewListTablesPaginator(client, &dynamodb.ListTablesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, table := range page.TableNames {
			// This prevents Implicit memory aliasing in for loop
			table := table
			v, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: &table,
			})
			if err != nil {
				return nil, err
			}

			continuousBackup, err := client.DescribeContinuousBackups(ctx, &dynamodb.DescribeContinuousBackupsInput{
				TableName: &table,
			})
			if err != nil {
				return nil, err
			}

			tags, err := client.ListTagsOfResource(ctx, &dynamodb.ListTagsOfResourceInput{
				ResourceArn: v.Table.TableArn,
			})
			if err != nil {
				return nil, err
			}

			values = append(values, Resource{
				ARN:  *v.Table.TableArn,
				Name: *v.Table.TableName,
				Description: model.DynamoDbTableDescription{
					Table:            v.Table,
					ContinuousBackup: continuousBackup.ContinuousBackupsDescription,
					Tags:             tags.Tags,
				},
			})
		}
	}

	return values, nil
}

func DynamoDbGlobalSecondaryIndex(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	client := dynamodb.NewFromConfig(cfg)
	paginator := dynamodb.NewListTablesPaginator(client, &dynamodb.ListTablesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, table := range page.TableNames {
			// This prevents Implicit memory aliasing in for loop
			table := table
			tableOutput, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: &table,
			})
			if err != nil {
				return nil, err
			}

			for _, v := range tableOutput.Table.GlobalSecondaryIndexes {
				values = append(values, Resource{
					ARN:  *v.IndexArn,
					Name: *v.IndexName,
					Description: model.DynamoDbGlobalSecondaryIndexDescription{
						GlobalSecondaryIndex: v,
					},
				})
			}
		}
	}

	return values, nil
}

func DynamoDbLocalSecondaryIndex(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	client := dynamodb.NewFromConfig(cfg)
	paginator := dynamodb.NewListTablesPaginator(client, &dynamodb.ListTablesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, table := range page.TableNames {
			// This prevents Implicit memory aliasing in for loop
			table := table
			tableOutput, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: &table,
			})
			if err != nil {
				return nil, err
			}

			for _, v := range tableOutput.Table.LocalSecondaryIndexes {
				values = append(values, Resource{
					ARN:  *v.IndexArn,
					Name: *v.IndexName,
					Description: model.DynamoDbLocalSecondaryIndexDescription{
						LocalSecondaryIndex: v,
					},
				})
			}
		}
	}

	return values, nil
}

func DynamoDbStream(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	client := dynamodbstreams.NewFromConfig(cfg)

	streams, err := client.ListStreams(ctx, &dynamodbstreams.ListStreamsInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range streams.Streams {
		values = append(values, Resource{
			ARN:  *v.StreamArn,
			Name: *v.StreamLabel,
			Description: model.DynamoDbStreamDescription{
				Stream: v,
			},
		})
	}

	return values, nil
}

func DynamoDbBackUp(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	client := dynamodb.NewFromConfig(cfg)
	backups, err := client.ListBackups(ctx, &dynamodb.ListBackupsInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, v := range backups.BackupSummaries {
		values = append(values, Resource{
			ARN:  *v.BackupArn,
			Name: *v.BackupName,
			Description: model.DynamoDbBackupDescription{
				Backup: v,
			},
		})
	}

	return values, nil
}

func DynamoDbGlobalTable(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	client := dynamodb.NewFromConfig(cfg)
	globalTables, err := client.ListGlobalTables(ctx, &dynamodb.ListGlobalTablesInput{})
	if err != nil {
		return nil, err
	}

	var values []Resource
	for _, table := range globalTables.GlobalTables {
		globalTable, err := client.DescribeGlobalTable(ctx, &dynamodb.DescribeGlobalTableInput{
			GlobalTableName: table.GlobalTableName,
		})
		if err != nil {
			return nil, err
		}
		values = append(values, Resource{
			ARN:  *globalTable.GlobalTableDescription.GlobalTableArn,
			Name: *globalTable.GlobalTableDescription.GlobalTableName,
			Description: model.DynamoDbGlobalTableDescription{
				GlobalTable: *globalTable.GlobalTableDescription,
			},
		})
	}

	return values, nil
}
