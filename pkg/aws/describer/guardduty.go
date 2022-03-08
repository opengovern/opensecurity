package describer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"gitlab.com/keibiengine/keibi-engine/pkg/aws/model"
)

func GuardDutyFinding(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	var values []Resource

	client := guardduty.NewFromConfig(cfg)

	dpaginator := guardduty.NewListDetectorsPaginator(client, &guardduty.ListDetectorsInput{})
	for dpaginator.HasMorePages() {
		dpage, err := dpaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, detectorId := range dpage.DetectorIds {
			paginator := guardduty.NewListFindingsPaginator(client, &guardduty.ListFindingsInput{
				DetectorId: &detectorId,
			})

			for paginator.HasMorePages() {
				page, err := paginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				findings, err := client.GetFindings(ctx, &guardduty.GetFindingsInput{
					DetectorId: &detectorId,
					FindingIds: page.FindingIds,
				})
				if err != nil {
					return nil, err
				}

				for _, item := range findings.Findings {
					values = append(values, Resource{
						ARN:  *item.Arn,
						Name: *item.Id,
						Description: model.GuardDutyFindingDescription{
							Finding: item,
						},
					})
				}
			}
		}
	}
	return values, nil
}

func GuardDutyDetector(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	var values []Resource

	client := guardduty.NewFromConfig(cfg)

	paginator := guardduty.NewListDetectorsPaginator(client, &guardduty.ListDetectorsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, id := range page.DetectorIds {
			out, err := client.GetDetector(ctx, &guardduty.GetDetectorInput{
				DetectorId: &id,
			})
			if err != nil {
				return nil, err
			}

			values = append(values, Resource{
				ID:   id,
				Name: id,
				Description: model.GuardDutyDetectorDescription{
					DetectorId: id,
					Detector:   out,
				},
			})
		}
	}
	return values, nil
}
