package tasks

import (
	"crypto/rsa"
	"encoding/json"
	api2 "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/opensecurity/pkg/utils"
	"github.com/opengovern/opensecurity/services/tasks/api"
	"github.com/opengovern/opensecurity/services/tasks/db"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type httpRoutes struct {
	logger *zap.Logger

	platformPrivateKey *rsa.PrivateKey
	db                 db.Database
}

func (r *httpRoutes) Register(e *echo.Echo) {
	v1 := e.Group("/api/v1")
	// List all tasks
	v1.GET("/tasks", httpserver.AuthorizeHandler(r.ListTasks, api2.ViewerRole))
	// Get task
	v1.GET("/tasks/:id", httpserver.AuthorizeHandler(r.GetTask, api2.ViewerRole))
	// Create a new task
	v1.POST("/tasks/run", httpserver.AuthorizeHandler(r.RunTask, api2.EditorRole))
	// Get Task Result
	v1.GET("/tasks/run/:id", httpserver.AuthorizeHandler(r.GetTaskRunResult, api2.ViewerRole))
	// List Tasks Result
	v1.GET("/tasks/run", httpserver.AuthorizeHandler(r.ListTaskRunResult, api2.ViewerRole))

}

func bindValidate(ctx echo.Context, i interface{}) error {
	if err := ctx.Bind(i); err != nil {
		return err
	}

	if err := ctx.Validate(i); err != nil {
		return err
	}

	return nil
}

// ListTasks godoc
//
//	@Summary	List tasks
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		cursor		query	int	false	"cursor"
//	@Param		per_page	query	int	false	"per page"
//	@Produce	json
//	@Success	200	{object}	api.ListTaskRunsResponse
//	@Router		/tasks/api/v1/tasks [get]
func (r *httpRoutes) ListTasks(ctx echo.Context) error {
	var cursor, perPage int64
	var err error
	cursorStr := ctx.QueryParam("cursor")
	if cursorStr != "" {
		cursor, err = strconv.ParseInt(cursorStr, 10, 64)
		if err != nil {
			return err
		}
	}
	perPageStr := ctx.QueryParam("per_page")
	if perPageStr != "" {
		perPage, err = strconv.ParseInt(perPageStr, 10, 64)
		if err != nil {
			return err
		}
	}

	items, err := r.db.GetTaskList()
	if err != nil {
		r.logger.Error("failed to get tasks", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, "failed to get tasks")

	}

	totalCount := len(items)
	if perPage != 0 {
		if cursor == 0 {
			items = utils.Paginate(1, perPage, items)
		} else {
			items = utils.Paginate(cursor, perPage, items)
		}
	}
	var taskResponses []api.TaskResponse
	for _, task := range items {
		taskResponses = append(taskResponses, api.TaskResponse{
			ID:          task.ID,
			Name:        task.Name,
			ResultType:  task.ResultType,
			Description: task.Description,
			ImageUrl:    task.ImageUrl,
			Interval:    task.Interval,
			Timeout:     task.Timeout,
		})
	}

	return ctx.JSON(http.StatusOK, api.TaskListResponse{
		TotalCount: totalCount,
		Items:      taskResponses,
	})
}

// GetTask godoc
//
//	@Summary	Get task by id
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		id	path	string	true	"run id"
//	@Produce	json
//	@Success	200	{object}	models.Task
//	@Router		/tasks/api/v1/tasks/:id [get]
func (r *httpRoutes) GetTask(ctx echo.Context) error {
	id := ctx.Param("id")
	task, err := r.db.GetTask(id)
	if err != nil {
		r.logger.Error("failed to get task results", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, "failed to get task results")
	}
	var taskResponse api.TaskResponse
	taskResponse = api.TaskResponse{
		ID:          task.ID,
		Name:        task.Name,
		ResultType:  task.ResultType,
		Description: task.Description,
		ImageUrl:    task.ImageUrl,
		Interval:    task.Interval,
		Timeout:     task.Timeout,
	}

	return ctx.JSON(http.StatusOK, taskResponse)
}

// RunTask godoc
//
//	@Summary	Run a new task
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request	body	api.RunTaskRequest	true	"Run task request"
//	@Produce	json
//	@Success	200	{object}	models.TaskRun
//	@Router		/tasks/api/v1/tasks/run [post]
func (r *httpRoutes) RunTask(ctx echo.Context) error {
	var req api.RunTaskRequest
	if err := bindValidate(ctx, &req); err != nil {
		r.logger.Error("failed to bind task", zap.Error(err))
		return ctx.JSON(http.StatusBadRequest, "failed to bind task")
	}

	task, _ := r.db.GetTask(req.TaskID)
	if task == nil {
		r.logger.Error("failed to find task", zap.String("task", req.TaskID))
		return ctx.JSON(http.StatusInternalServerError, "failed to find task")
	}

	run := models.TaskRun{
		TaskID: req.TaskID,
		Status: models.TaskRunStatusCreated,
	}
	paramsJson, err := json.Marshal(req.Params)
	if err != nil {
		return err
	}
	err = run.Params.Set(paramsJson)
	if err != nil {
		r.logger.Error("failed to set params", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to set params")
	}
	err = run.Result.Set([]byte("{}"))
	if err != nil {
		r.logger.Error("failed to set results", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to set results")
	}

	if err := r.db.CreateTaskRun(&run); err != nil {
		r.logger.Error("failed to create task run", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, "failed to create task run")
	}

	return ctx.JSON(http.StatusCreated, run)
}

// GetTaskRunResult godoc
//
//	@Summary	Get task run
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		id	path	string	true	"run id"
//	@Produce	json
//	@Success	200	{object}	models.TaskRun
//	@Router		/tasks/api/v1/tasks/run/:id [get]
func (r *httpRoutes) GetTaskRunResult(ctx echo.Context) error {
	id := ctx.Param("id")
	var cursor, perPage int64
	taskResults, err := r.db.GetTaskRunResult(id)
	if err != nil {
		r.logger.Error("failed to get task results", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, "failed to get task results")
	}
	cursorStr := ctx.QueryParam("cursor")
	if cursorStr != "" {
		cursor, err = strconv.ParseInt(cursorStr, 10, 64)
		if err != nil {
			return err
		}
	}
	perPageStr := ctx.QueryParam("per_page")
	if perPageStr != "" {
		perPage, err = strconv.ParseInt(perPageStr, 10, 64)
		if err != nil {
			return err
		}
	}
	totalCount := len(taskResults)
	if perPage != 0 {
		if cursor == 0 {
			taskResults = utils.Paginate(1, perPage, taskResults)
		} else {
			taskResults = utils.Paginate(cursor, perPage, taskResults)
		}
	}
	var taskRunResponses []api.TaskRun
	for _, task := range taskResults {
		var params map[string]interface{}
		err := json.Unmarshal(task.Params.Bytes, &params)
		if err != nil {
			r.logger.Error("failed to unmarshal params", zap.Error(err))
			return ctx.JSON(http.StatusInternalServerError, "failed to unmarshal params")
		}
		var result map[string]interface{}
		err = json.Unmarshal(task.Result.Bytes, &result)
		if err != nil {
			r.logger.Error("failed to unmarshal result", zap.Error(err))
			return ctx.JSON(http.StatusInternalServerError, "failed to unmarshal result")
		}
		taskRunResponses = append(taskRunResponses, api.TaskRun{
			ID:             task.ID,
			CreatedAt:      task.CreatedAt,
			UpdatedAt:      task.UpdatedAt,
			TaskID:         task.TaskID,
			Status:         string(task.Status),
			Result:         result,
			Params:         params,
			FailureMessage: task.FailureMessage,
		})
	}

	return ctx.JSON(http.StatusOK, api.ListTaskRunsResponse{
		TotalCount: totalCount,
		Items:      taskRunResponses,
	})

}

// ListTaskRunResult godoc
//
//	@Summary	List task runs
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		cursor		query	int	false	"cursor"
//	@Param		per_page	query	int	false	"per page"
//	@Produce	json
//	@Success	200	{object}	api.ListTaskRunsResponse
//	@Router		/tasks/api/v1/tasks/run [get]
func (r *httpRoutes) ListTaskRunResult(ctx echo.Context) error {
	var cursor, perPage int64
	var err error
	cursorStr := ctx.QueryParam("cursor")
	if cursorStr != "" {
		cursor, err = strconv.ParseInt(cursorStr, 10, 64)
		if err != nil {
			return err
		}
	}
	perPageStr := ctx.QueryParam("per_page")
	if perPageStr != "" {
		perPage, err = strconv.ParseInt(perPageStr, 10, 64)
		if err != nil {
			return err
		}
	}

	items, err := r.db.ListTaskRunResult()
	if err != nil {
		r.logger.Error("failed to get task results", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, "failed to get task results")
	}

	totalCount := len(items)
	if perPage != 0 {
		if cursor == 0 {
			items = utils.Paginate(1, perPage, items)
		} else {
			items = utils.Paginate(cursor, perPage, items)
		}
	}
	var taskRunResponses []api.TaskRun
	for _, task := range items {
		var params map[string]interface{}
		err := json.Unmarshal(task.Params.Bytes, &params)
		if err != nil {
			r.logger.Error("failed to unmarshal params", zap.Error(err))
			return ctx.JSON(http.StatusInternalServerError, "failed to unmarshal params")
		}
		var result map[string]interface{}
		err = json.Unmarshal(task.Result.Bytes, &result)
		if err != nil {
			r.logger.Error("failed to unmarshal result", zap.Error(err))
			return ctx.JSON(http.StatusInternalServerError, "failed to unmarshal result")
		}
		taskRunResponses = append(taskRunResponses, api.TaskRun{
			ID:             task.ID,
			CreatedAt:      task.CreatedAt,
			UpdatedAt:      task.UpdatedAt,
			TaskID:         task.TaskID,
			Status:         string(task.Status),
			Result:         result,
			Params:         params,
			FailureMessage: task.FailureMessage,
		})
	}
	return ctx.JSON(http.StatusOK, api.ListTaskRunsResponse{
		TotalCount: totalCount,
		Items:      taskRunResponses,
	})
}
