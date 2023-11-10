package statemanager

import "github.com/kaytu-io/kaytu-engine/pkg/workspace/db"

func (s *Service) runBootstrapping(workspace *db.Workspace) error {
	creds, err := s.db.ListCredentialsByWorkspace(workspace.Name)
	if err != nil {
		return err
	}

	if !workspace.IsCreated && len(creds) > 0 {
		return s.createWorkspace(workspace)
	}

	if workspace.IsCreated {
		for _, cred := range creds {
			if !cred.IsCreated {
				err := s.addCredentialToWorkspace(workspace.ID, cred.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	if workspace.IsBootstrapInputFinished {
		// run jobs
	}

	// jobs finished
	// change to provisioned

	return nil
}
