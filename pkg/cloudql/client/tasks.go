package opengovernance_client

import (
	"context"
	"encoding/json"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"runtime"
	"time"

	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/pg"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type TaskRow struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"image_url"`
	LastRun     time.Time `json:"last_run"`
	Params      []string  `json:"params"`
}

func getTaskRowFromTask(ctx context.Context, task models.Task, taskRun *models.TaskRun) (*TaskRow, error) {
	var paramsKeys []string
	var params map[string]interface{}
	err := json.Unmarshal(taskRun.Params.Bytes, &params)
	if err != nil {
		return nil, err
	}

	for k, _ := range params {
		paramsKeys = append(paramsKeys, k)
	}

	row := TaskRow{
		ID:          task.ID,
		Title:       task.Name,
		Description: task.Description,
		ImageUrl:    task.ImageUrl,
		LastRun:     taskRun.CreatedAt,
		Params:      paramsKeys,
	}
	return &row, nil
}

func ListTasks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListTasks")
	runtime.GC()
	cfg := config.GetConfig(d.Connection)
	ke, err := pg.NewTaskClient(cfg, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{PG: ke}

	tasks, err := k.PG.ListTasks(ctx)
	if err != nil {
		return nil, err
	}

	for _, i := range tasks {
		taskRun, err := k.PG.GetLastTaskRun(i.ID)
		if err != nil {
			return nil, err
		}
		row, err := getTaskRowFromTask(ctx, i, taskRun)
		if err != nil {
			plugin.Logger(ctx).Error("ListTasks", "task", i, "error", err)
			continue
		}
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}

func GetTask(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("GetTask")
	runtime.GC()
	cfg := config.GetConfig(d.Connection)
	ke, err := pg.NewTaskClient(cfg, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{PG: ke}

	taskId := d.EqualsQuals["id"].GetStringValue()
	i, err := k.PG.GetTask(ctx, taskId)
	if err != nil {
		return nil, err
	}
	if i == nil {
		return nil, nil
	}

	taskRun, err := k.PG.GetLastTaskRun(i.ID)
	if err != nil {
		return nil, err
	}

	row, err := getTaskRowFromTask(ctx, *i, taskRun)
	if err != nil {
		plugin.Logger(ctx).Error("GetTask", "task", i, "error", err)
		return nil, err
	}
	return row, nil
}
