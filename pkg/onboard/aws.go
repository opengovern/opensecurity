package onboard

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/smithy-go"
	keibiaws "gitlab.com/keibiengine/keibi-engine/pkg/aws"
	"gitlab.com/keibiengine/keibi-engine/pkg/aws/describer"
	"gitlab.com/keibiengine/keibi-engine/pkg/onboard/api"
	"gitlab.com/keibiengine/keibi-engine/pkg/source"
)

var PermissionError = errors.New("PermissionError")

func discoverAwsAccounts(ctx context.Context, req api.DiscoverAWSAccountsRequest) ([]api.DiscoverAWSAccountsResponse, error) {
	err := keibiaws.CheckDescribeRegionsPermission(req.AccessKey, req.SecretKey)
	if err != nil {
		return nil, PermissionError
	}

	isAttached, err := keibiaws.CheckAttachedPolicy(req.AccessKey, req.SecretKey, keibiaws.SecurityAuditPolicyARN)
	if err != nil {
		return nil, PermissionError
	}
	if !isAttached {
		return nil, PermissionError
	}

	cfg, err := keibiaws.GetConfig(ctx, req.AccessKey, req.SecretKey, "", "")
	if err != nil {
		return nil, err
	}

	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	acc, err := currentAwsAccount(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if acc.Name == "" {
		acc.Name = acc.AccountID
	}
	return []api.DiscoverAWSAccountsResponse{acc}, nil
	//
	//accounts, err := describer.OrganizationAccounts(ctx, cfg)
	//if err != nil {
	//	if !ignoreAwsOrgError(err) {
	//		return nil, err
	//	}
	//	return []api.DiscoverAWSAccountsResponse{acc}, nil
	//}
	//if len(accounts) == 0 {
	//	return []api.DiscoverAWSAccountsResponse{acc}, nil
	//}
	//
	//discovered := make([]api.DiscoverAWSAccountsResponse, 0, len(accounts))
	//for _, item := range accounts {
	//	if *item.Name == "" {
	//		*item.Name = *item.Id
	//	}
	//	discovered = append(discovered, api.DiscoverAWSAccountsResponse{
	//		AccountID:      *item.Id,
	//		Status:         string(item.Status),
	//		Name:           *item.Name,
	//		Email:          *item.Email,
	//		OrganizationID: acc.OrganizationID,
	//	})
	//}
	//
	//return discovered, nil
}

func currentAwsAccount(ctx context.Context, cfg aws.Config) (api.DiscoverAWSAccountsResponse, error) {
	accID, err := describer.STSAccount(ctx, cfg)
	if err != nil {
		return api.DiscoverAWSAccountsResponse{}, err
	}

	var (
		orgId    string
		accName  string
		accEmail string
	)
	orgs, err := describer.OrganizationOrganization(ctx, cfg)
	if err != nil {
		if !ignoreAwsOrgError(err) {
			return api.DiscoverAWSAccountsResponse{}, err
		}
	} else {
		orgId = *orgs.Id
	}

	acc, err := describer.OrganizationAccount(ctx, cfg, accID)
	if err != nil {
		if !ignoreAwsOrgError(err) {
			return api.DiscoverAWSAccountsResponse{}, err
		}
	} else {
		accName = *acc.Name
		accEmail = *acc.Email
	}

	return api.DiscoverAWSAccountsResponse{
		AccountID:      accID,
		Status:         string(types.AccountStatusActive),
		OrganizationID: orgId,
		Name:           accName,
		Email:          accEmail,
	}, nil
}

func getAWSCredentialsMetadata(ctx context.Context, config api.SourceConfigAWS) (*source.AWSCredentialMetadata, error) {
	creds, err := keibiaws.GetConfig(ctx, config.AccessKey, config.SecretKey, "", "")
	if err != nil {
		return nil, err
	}
	if creds.Region == "" {
		creds.Region = "us-east-1"
	}

	iamClient := iam.NewFromConfig(creds)
	user, err := iamClient.GetUser(ctx, &iam.GetUserInput{})
	if err != nil {
		fmt.Printf("failed to get user: %v", err)
		return nil, err
	}
	paginator := iam.NewListAttachedUserPoliciesPaginator(iamClient, &iam.ListAttachedUserPoliciesInput{
		UserName: user.User.UserName,
	})

	policyARNs := make([]string, 0)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			fmt.Printf("failed to get policy page: %v", err)
			return nil, err
		}
		for _, policy := range page.AttachedPolicies {
			policyARNs = append(policyARNs, *policy.PolicyArn)
		}
	}

	//TODO get metadata from aws

	return &source.AWSCredentialMetadata{
		AccountID:        config.AccountId,
		IamUserName:      user.User.UserName,
		AttachedPolicies: policyARNs,
	}, nil

}

func ignoreAwsOrgError(err error) bool {
	var ae smithy.APIError
	return errors.As(err, &ae) &&
		(ae.ErrorCode() == (&types.AWSOrganizationsNotInUseException{}).ErrorCode() ||
			ae.ErrorCode() == (&types.AccessDeniedException{}).ErrorCode())
}
