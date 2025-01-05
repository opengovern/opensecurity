package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
)

const DeleteTasksIndex = "delete_tasks"

type DeletingResource struct {
	Key        []byte
	ResourceID string
	Index      string
}

type DeleteTaskType string

const (
	DeleteTaskTypeResource DeleteTaskType = "resource"
	DeleteTaskTypeQuery    DeleteTaskType = "query"
)

type DeleteTask struct {
	EsID    string `json:"es_id"`
	EsIndex string `json:"es_index"`

	TaskType        DeleteTaskType   `json:"task_type"`
	DiscoveryJobID  uint             `json:"discovery_job_id"`
	IntegrationID   string           `json:"integration_id"`
	ResourceType    string           `json:"resource_type"`
	IntegrationType integration.Type `json:"connector"`

	DeletingResources []DeletingResource `json:"deleting_resources"`
	Query             string             `json:"query"`
	QueryIndex        string             `json:"query_index"`
}

func (d DeleteTask) KeysAndIndex() ([]string, string) {
	var ids []string
	for _, r := range d.DeletingResources {
		ids = append(ids, r.ResourceID)
		ids = append(ids, r.Index)
	}
	ids = append(ids, d.ResourceType)
	ids = append(ids, d.IntegrationID)
	ids = append(ids, string(d.TaskType))
	ids = append(ids, d.Query)
	return ids, DeleteTasksIndex
}

type GetDeleteTasksResponse struct {
	Hits GetDeleteTasksHits `json:"hits"`
}
type GetDeleteTasksHits struct {
	Total opengovernance.SearchTotal `json:"total"`
	Hits  []GetDeleteTasksHit        `json:"hits"`
}
type GetDeleteTasksHit struct {
	ID      string     `json:"_id"`
	Score   float64    `json:"_score"`
	Index   string     `json:"_index"`
	Type    string     `json:"_type"`
	Version int64      `json:"_version,omitempty"`
	Source  DeleteTask `json:"_source"`
	Sort    []any      `json:"sort"`
}

func GetDeleteTasks(ctx context.Context, client opengovernance.Client) (*GetDeleteTasksResponse, error) {
	root := map[string]any{}
	root["size"] = 10000
	root["sort"] = []map[string]any{
		{"_id": "desc"},
	}

	queryBytes, err := json.Marshal(root)
	if err != nil {
		return nil, err
	}

	var response GetDeleteTasksResponse
	err = client.Search(ctx, DeleteTasksIndex,
		string(queryBytes), &response)
	if err != nil {
		fmt.Println("query=", string(queryBytes))
		return nil, err
	}

	return &response, nil
}
