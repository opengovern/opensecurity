package statemanager

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kaytu-io/kaytu-engine/pkg/workspace/api"
	"github.com/kaytu-io/kaytu-engine/pkg/workspace/db"
	"github.com/kaytu-io/kaytu-engine/pkg/workspace/state"
	"github.com/kaytu-io/kaytu-engine/pkg/workspace/transactions"
	"github.com/sony/sonyflake"
)

func (s *Service) UseReservationIfPossible(workspace db.Workspace) error {
	rs, err := s.db.GetReservedWorkspace()
	if err != nil {
		return err
	}

	if rs == nil {
		return nil
	}

	err = s.db.DeleteWorkspace(workspace.ID)
	if err != nil {
		return err
	}

	err = s.db.UpdateCredentialWSID(workspace.ID, rs.ID)
	if err != nil {
		return err
	}

	workspace.ID = rs.ID
	err = s.db.UpdateWorkspace(&workspace)
	if err != nil {
		return err
	}

	err = s.db.DeleteWorkspaceTransaction(workspace.ID, string(transactions.Transaction_CreateHelmRelease))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) handleReservation() error {
	rs, err := s.db.GetReservedWorkspace()
	if err != nil {
		return err
	}

	if rs != nil {
		return nil
	}

	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := sf.NextID()
	if err != nil {
		return err
	}

	awsUID, err := sf.NextID()
	if err != nil {
		return err
	}

	workspace := &db.Workspace{
		ID:             fmt.Sprintf("ws-%d", id),
		Name:           "",
		AWSUniqueId:    aws.String(fmt.Sprintf("aws-uid-%d", awsUID)),
		OwnerId:        nil,
		Status:         string(state.StateID_Reserving),
		Size:           api.SizeXS,
		Tier:           api.Tier_Teams,
		OrganizationID: nil,
	}

	if err := s.db.CreateWorkspace(workspace); err != nil {
		return err
	}
	return nil
}
