package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/opengovern/og-util/pkg/integration"

	"github.com/opengovern/opensecurity/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	dexApi "github.com/dexidp/dex/api/v2"
	"github.com/jackc/pgtype"
	api3 "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/httpserver"
	model2 "github.com/opengovern/opensecurity/jobs/demo-importer-job/db/model"
	"github.com/opengovern/opensecurity/jobs/post-install-job/db/model"
	complianceapi "github.com/opengovern/opensecurity/services/compliance/api"
	integrationApi "github.com/opengovern/opensecurity/services/integration/api/models"
	integrationClient "github.com/opengovern/opensecurity/services/integration/client"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/labstack/echo/v4"
	"github.com/opengovern/opensecurity/services/core/api"
	"github.com/opengovern/opensecurity/services/core/db/models"
	coreUtils "github.com/opengovern/opensecurity/services/core/utils"
)

func (h *HttpHandler) Register(r *echo.Echo) {
	v1 := r.Group("/api/v1")
	// metadata
	filter := v1.Group("/filter")
	filter.POST("", httpserver.AuthorizeHandler(h.AddFilter, api3.ViewerRole))
	filter.GET("", httpserver.AuthorizeHandler(h.GetFilters, api3.ViewerRole))

	metadata := v1.Group("/metadata")
	metadata.GET("/:key", httpserver.AuthorizeHandler(h.GetConfigMetadata, api3.ViewerRole))
	metadata.POST("", httpserver.AuthorizeHandler(h.SetConfigMetadata, api3.AdminRole))

	queryParameter := v1.Group("/query_parameter")
	queryParameter.POST("/set", httpserver.AuthorizeHandler(h.SetQueryParameter, api3.AdminRole))
	queryParameter.POST("", httpserver.AuthorizeHandler(h.ListQueryParameters, api3.ViewerRole))
	queryParameter.GET("/:key", httpserver.AuthorizeHandler(h.GetQueryParameter, api3.ViewerRole))
	// inventory
	queryV1 := v1.Group("/query")
	queryV1.GET("", httpserver.AuthorizeHandler(h.ListQueries, api3.ViewerRole))
	queryV1.POST("/run", httpserver.AuthorizeHandler(h.RunQuery, api3.ViewerRole))
	queryV1.GET("/run/history", httpserver.AuthorizeHandler(h.GetRecentRanQueries, api3.ViewerRole))
	v2 := r.Group("/api/v2")
	metadatav2 := v2.Group("/metadata")
	metadatav2.GET("/resourcetype", httpserver.AuthorizeHandler(h.ListResourceTypeMetadata, api3.ViewerRole))

	resourceCollectionMetadata := metadata.Group("/resource-collection")
	resourceCollectionMetadata.GET("", httpserver.AuthorizeHandler(h.ListResourceCollectionsMetadata, api3.ViewerRole))
	resourceCollectionMetadata.GET("/:resourceCollectionId", httpserver.AuthorizeHandler(h.GetResourceCollectionMetadata, api3.ViewerRole))

	v3 := r.Group("/api/v3")
	// metadata
	v3.PUT("/sample/purge", httpserver.AuthorizeHandler(h.PurgeSampleData, api3.ViewerRole))
	v3.PUT("/sample/sync", httpserver.AuthorizeHandler(h.SyncDemo, api3.ViewerRole))
	v3.PUT("/sample/loaded", httpserver.AuthorizeHandler(h.WorkspaceLoadedSampleData, api3.ViewerRole))
	v3.GET("/sample/sync/status", httpserver.AuthorizeHandler(h.GetSampleSyncStatus, api3.ViewerRole))
	v3.GET("/migration/status", httpserver.AuthorizeHandler(h.GetMigrationStatus, api3.ViewerRole))
	v3.GET("/configured/status", httpserver.AuthorizeHandler(h.GetConfiguredStatus, api3.ViewerRole))
	v3.PUT("/configured/set", httpserver.AuthorizeHandler(h.SetConfiguredStatus, api3.AdminRole))
	v3.PUT("/configured/unset", httpserver.AuthorizeHandler(h.UnsetConfiguredStatus, api3.ViewerRole))
	v3.GET("/about", httpserver.AuthorizeHandler(h.GetAbout, api3.ViewerRole))

	v3.GET("/vault/configured", httpserver.AuthorizeHandler(h.VaultConfigured, api3.ViewerRole))

	views := v3.Group("/views")
	views.PUT("/reload", httpserver.AuthorizeHandler(h.ReloadViews, api3.AdminRole))
	views.GET("/checkpoint", httpserver.AuthorizeHandler(h.GetViewsCheckpoint, api3.AdminRole))
	views.GET("", httpserver.AuthorizeHandler(h.GetViews, api3.ViewerRole))
	// inventory
	v3.POST("/queries", httpserver.AuthorizeHandler(h.ListQueriesV2, api3.ViewerRole))
	v3.GET("/queries/cache-enabled", httpserver.AuthorizeHandler(h.ListCacheEnabledQueries, api3.ViewerRole))
	v3.GET("/queries/filters", httpserver.AuthorizeHandler(h.ListQueriesFilters, api3.ViewerRole))
	v3.GET("/queries/:query_id", httpserver.AuthorizeHandler(h.GetQuery, api3.ViewerRole))
	v3.GET("/queries/tags", httpserver.AuthorizeHandler(h.ListQueriesTags, api3.ViewerRole))
	v3.POST("/query/run", httpserver.AuthorizeHandler(h.RunQueryByID, api3.ViewerRole))
	v3.GET("/query/async/run/:run_id/result", httpserver.AuthorizeHandler(h.GetAsyncQueryRunResult, api3.ViewerRole))
	v3.GET("/resources/categories", httpserver.AuthorizeHandler(h.GetResourceCategories, api3.ViewerRole))
	v3.GET("/queries/categories", httpserver.AuthorizeHandler(h.GetQueriesResourceCategories, api3.ViewerRole))
	v3.GET("/tables/categories", httpserver.AuthorizeHandler(h.GetTablesResourceCategories, api3.ViewerRole))
	v3.GET("/categories/queries", httpserver.AuthorizeHandler(h.GetCategoriesQueries, api3.ViewerRole))
	v3.GET("/parameters/queries", httpserver.AuthorizeHandler(h.GetParametersQueries, api3.ViewerRole))

	v3.PUT("/plugins/:plugin_id/reload", httpserver.AuthorizeHandler(h.ReloadPluginSteampipeConfig, api3.AdminRole))
	v3.PUT("/plugins/:plugin_id/remove", httpserver.AuthorizeHandler(h.RemovePluginSteampipeConfig, api3.AdminRole))

	v4 := r.Group("/api/v4")
	v4.GET("/about", httpserver.AuthorizeHandler(h.GetAboutShort, api3.ViewerRole))
	v4.GET("/queries/sync", httpserver.AuthorizeHandler(h.SyncQueries, api3.ViewerRole))
	v4.POST("/layout/get", httpserver.AuthorizeHandler(h.GetUserLayouts, api3.ViewerRole))
	v4.POST("/layout/get-default", httpserver.AuthorizeHandler(h.GetUserDefaultLayout, api3.ViewerRole))
	v4.POST("/layout/set", httpserver.AuthorizeHandler(h.SetUserLayout, api3.ViewerRole))
	v4.POST("/layout/change-privacy", httpserver.AuthorizeHandler(h.ChangePrivacy, api3.ViewerRole))
	v4.GET("/layout/public", httpserver.AuthorizeHandler(h.GetPublicLayouts, api3.ViewerRole))

	v4.POST("/chatbot/generate-query", httpserver.AuthorizeHandler(h.GenerateQuery, api3.ViewerRole))
	v4.POST("/chatbot/generate-query/run", httpserver.AuthorizeHandler(h.GenerateQueryAndRun, api3.ViewerRole))
	v4.POST("/chatbot/secret", httpserver.AuthorizeHandler(h.ConfigureChatbotSecret, api3.AdminRole))
}

var tracer = otel.Tracer("core")

func bindValidate(ctx echo.Context, i interface{}) error {
	if err := ctx.Bind(i); err != nil {
		return err
	}

	if err := ctx.Validate(i); err != nil {
		return err
	}

	return nil
}

// GetConfigMetadata godoc
//
//	@Summary		Get key metadata
//	@Description	Returns the config metadata for the given key
//	@Security		BearerToken
//	@Tags			metadata
//	@Produce		json
//	@Param			key	path		string	true	"Key"
//	@Success		200	{object}	models.ConfigMetadata
//	@Router			/metadata/api/v1/metadata/{key} [get]
func (h *HttpHandler) GetConfigMetadata(ctx echo.Context) error {
	key := ctx.Param("key")
	_, span := tracer.Start(ctx.Request().Context(), "new_GetConfigMetadata", trace.WithSpanKind(trace.SpanKindServer))
	span.SetName("new_GetConfigMetadata")

	metadata, err := coreUtils.GetConfigMetadata(h.db, key)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "config not found")
		}
		return err
	}
	span.AddEvent("information", trace.WithAttributes(
		attribute.String("key", key),
	))
	span.End()
	return ctx.JSON(http.StatusOK, metadata.GetCore())
}

// SetConfigMetadata godoc
//
//	@Summary		Set key metadata
//	@Description	Sets the config metadata for the given key
//	@Security		BearerToken
//	@Tags			metadata
//	@Produce		json
//	@Param			req	body	api.SetConfigMetadataRequest	true	"Request Body"
//	@Success		200
//	@Router			/metadata/api/v1/metadata [post]
func (h *HttpHandler) SetConfigMetadata(ctx echo.Context) error {
	var req api.SetConfigMetadataRequest
	if err := bindValidate(ctx, &req); err != nil {
		return err
	}

	key, err := models.ParseMetadataKey(req.Key)
	if err != nil {
		return err
	}

	err = httpserver.RequireMinRole(ctx, key.GetMinAuthRole())
	if err != nil {
		return err
	}
	_, span := tracer.Start(ctx.Request().Context(), "new_SetConfigMetadata", trace.WithSpanKind(trace.SpanKindServer))
	span.SetName("new_SetConfigMetadata")

	err = coreUtils.SetConfigMetadata(h.db, key, req.Value)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.AddEvent("information", trace.WithAttributes(
		attribute.String("key", key.String()),
	))
	span.End()

	return ctx.JSON(http.StatusOK, nil)
}

// AddFilter godoc
//
//	@Summary	add filter
//	@Security	BearerToken
//	@Tags		metadata
//	@Produce	json
//	@Param		req	body	models.Filter	true	"Request Body"
//	@Success	200
//	@Router		/metadata/api/v1/filter [post]
func (h *HttpHandler) AddFilter(ctx echo.Context) error {
	var req models.Filter
	if err := bindValidate(ctx, &req); err != nil {
		return err
	}
	// trace :
	_, span := tracer.Start(ctx.Request().Context(), "new_AddFilter", trace.WithSpanKind(trace.SpanKindServer))
	span.SetName("new_AddFilter")

	err := h.db.AddFilter(models.Filter{Name: req.Name, KeyValue: req.KeyValue})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.AddEvent("information", trace.WithAttributes(
		attribute.String("name", req.Name),
	))
	span.End()
	return ctx.JSON(http.StatusOK, nil)
}

// GetFilters godoc
//
//	@Summary	list filters
//	@Security	BearerToken
//	@Tags		metadata
//	@Produce	json
//	@Success	200	{object}	[]models.Filter
//	@Router		/metadata/api/v1/filter [get]
func (h *HttpHandler) GetFilters(ctx echo.Context) error {
	// trace :
	_, span := tracer.Start(ctx.Request().Context(), "new_ListFilters", trace.WithSpanKind(trace.SpanKindServer))
	span.SetName("new_ListFilters")

	filters, err := h.db.ListFilters()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil
	}
	span.End()
	return ctx.JSON(http.StatusOK, filters)
}

// SetQueryParameter godoc
//
//	@Summary		Set query parameter
//	@Description	Sets the query parameters from the request body
//	@Security		BearerToken
//	@Tags			metadata
//	@Produce		json
//	@Param			req	body	api.SetQueryParameterRequest	true	"Request Body"
//	@Success		200
//	@Router			/metadata/api/v1/query_parameter [post]
func (h *HttpHandler) SetQueryParameter(ctx echo.Context) error {
	var req api.SetQueryParameterRequest
	if err := bindValidate(ctx, &req); err != nil {
		return err
	}

	if len(req.QueryParameters) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "no query parameters provided")
	}

	dbQueryParams := make([]*models.PolicyParameterValues, 0, len(req.QueryParameters))
	for _, apiParam := range req.QueryParameters {
		//key, err := models.ParseQueryParameterKey(apiParam.Key)
		//if err != nil {
		//	return err
		//}
		dbParam := models.QueryParameterFromAPI(apiParam)
		dbParam.Key = apiParam.Key
		dbQueryParams = append(dbQueryParams, &dbParam)
	}

	err := h.db.SetQueryParameters(dbQueryParams)
	if err != nil {
		h.logger.Error("error setting query parameters", zap.Error(err))
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

// ListQueryParameters godoc
//
//	@Summary		List query parameters
//	@Description	Returns the list of query parameters
//	@Security		BearerToken
//	@Tags			metadata
//	@Produce		json
//	@Param			query_id	query		string	false	"Policy ID to filter with"
//	@Param			control_id	query		string	false	"Control ID to filter with"
//	@Param			cursor		query		int		false	"Cursor"
//	@Param			per_page	query		int		false	"Per Page"
//	@Success		200			{object}	api.ListQueryParametersResponse
//	@Router			/metadata/api/v1/query_parameter [post]
func (h *HttpHandler) ListQueryParameters(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: api3.AdminRole}

	var cursor, perPage int64
	var err error
	var request api.ListQueryParametersRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	cursor = request.Cursor
	perPage = request.PerPage

	queryIDs := request.Queries
	controlIDs := request.Controls

	var filteredQueryParams []string
	if controlIDs != nil {
		if !h.complianceEnabled {
			return echo.NewHTTPError(http.StatusBadRequest, "compliance service is not enabled")
		}
		all_control, err := h.complianceClient.ListControl(clientCtx, controlIDs, nil)
		if err != nil {
			h.logger.Error("error getting control", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "error getting control")
		}
		if all_control == nil {
			return echo.NewHTTPError(http.StatusNotFound, "control not found")
		}
		for _, control := range all_control {
			for _, param := range control.Policy.Parameters {
				filteredQueryParams = append(filteredQueryParams, param.Key)
			}
		}
	} else if queryIDs != nil {
		// TODO: Fix this part and write new client on inventory
		queries, err := h.ListQueriesV2Internal(api.ListQueryV2Request{QueryIDs: queryIDs})
		if err != nil {
			h.logger.Error("error getting query", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "error getting query")
		}
		for _, q := range queries.Items {
			for _, param := range q.Query.Parameters {
				filteredQueryParams = append(filteredQueryParams, param.Key)
			}
		}
	}

	var queryParams []models.PolicyParameterValues
	if request.KeyRegex != nil {
		queryParams, err = h.db.GetQueryParametersValues(request.KeyRegex)
		if err != nil {
			h.logger.Error("error getting query parameters", zap.Error(err))
			return err
		}
	} else if controlIDs != nil || queryIDs != nil {
		queryParams, err = h.db.GetQueryParametersByIds(filteredQueryParams)
		if err != nil {
			h.logger.Error("error getting query parameters", zap.Error(err))
			return err
		}
	} else {
		queryParams, err = h.db.GetQueryParametersValues(nil)
		if err != nil {
			h.logger.Error("error getting query parameters", zap.Error(err))
			return err
		}
	}

	parametersMap := make(map[string]*api.QueryParameter)
	for _, dbParam := range queryParams {
		apiParam := dbParam.ToAPI()
		parametersMap[apiParam.Key+apiParam.ControlID] = &apiParam
	}

	var controls []complianceapi.Control

	if h.complianceEnabled {
		controls, err = h.complianceClient.ListControl(clientCtx, nil, nil)
		if err != nil {
			h.logger.Error("error listing controls", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "error listing controls")
		}
	}

	namedQueries, err := h.ListQueriesV2Internal(api.ListQueryV2Request{})
	if err != nil {
		h.logger.Error("error listing queries", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "error listing queries")
	}

	for _, c := range controls {
		for _, p := range c.Policy.Parameters {
			if _, ok := parametersMap[p.Key]; ok {
				parametersMap[p.Key].ControlsCount += 1
			}
			if _, ok := parametersMap[p.Key+c.ID]; ok {
				parametersMap[p.Key+c.ID].ControlsCount += 1
			}
		}
	}
	for _, q := range namedQueries.Items {
		for _, p := range q.Query.Parameters {
			if _, ok := parametersMap[p.Key]; ok {
				parametersMap[p.Key].QueriesCount += 1
			}
		}
	}

	var items []api.QueryParameter
	for _, i := range parametersMap {
		items = append(items, *i)
	}

	totalCount := len(items)

	sortOrder := "desc"
	sortBy := "key"
	if request.SortOrder != nil {
		sortOrder = strings.ToLower(*request.SortOrder)
	}
	if request.SortBy != nil {
		sortBy = strings.ToLower(*request.SortBy)
	}

	switch sortOrder {
	case "asc":
		switch sortBy {
		case "key":
			sort.Slice(items, func(i, j int) bool {
				if items[i].Key == items[j].Key {
					return items[i].ControlsCount < items[j].ControlsCount
				}
				return items[i].Key < items[j].Key
			})
		case "value":
			sort.Slice(items, func(i, j int) bool {
				if items[i].Value == items[j].Value {
					if items[i].Key == items[j].Key {
						return items[i].ControlsCount < items[j].ControlsCount
					}
					return items[i].Key < items[j].Key
				}
				return items[i].Value < items[j].Value
			})
		case "controlid", "control_id":
			sort.Slice(items, func(i, j int) bool {
				if items[i].ControlID == items[j].ControlID {
					return items[i].Key < items[j].Key
				}
				return items[i].ControlID < items[j].ControlID
			})
		case "controlscount", "controls_count":
			sort.Slice(items, func(i, j int) bool {
				if items[i].ControlsCount == items[j].ControlsCount {
					return items[i].Key < items[j].Key
				}
				return items[i].ControlsCount < items[j].ControlsCount
			})
		case "queriescount", "queries_count":
			sort.Slice(items, func(i, j int) bool {
				if items[i].QueriesCount == items[j].QueriesCount {
					return items[i].Key < items[j].Key
				}
				return items[i].QueriesCount < items[j].QueriesCount
			})
		}
	case "desc":
		switch sortBy {
		case "key":
			sort.Slice(items, func(i, j int) bool {
				if items[i].Key == items[j].Key {
					return items[i].ControlsCount > items[j].ControlsCount
				}
				return items[i].Key > items[j].Key
			})
		case "value":
			sort.Slice(items, func(i, j int) bool {
				if items[i].Value == items[j].Value {
					if items[i].Key == items[j].Key {
						return items[i].ControlsCount > items[j].ControlsCount
					}
					return items[i].Key > items[j].Key
				}
				return items[i].Value > items[j].Value
			})
		case "controlid", "control_id":
			sort.Slice(items, func(i, j int) bool {
				if items[i].ControlID == items[j].ControlID {
					return items[i].Key > items[j].Key
				}
				return items[i].ControlID > items[j].ControlID
			})
		case "controlscount", "controls_count":
			sort.Slice(items, func(i, j int) bool {
				if items[i].ControlsCount == items[j].ControlsCount {
					return items[i].Key > items[j].Key
				}
				return items[i].ControlsCount > items[j].ControlsCount
			})
		case "queriescount", "queries_count":
			sort.Slice(items, func(i, j int) bool {
				if items[i].QueriesCount == items[j].QueriesCount {
					return items[i].Key > items[j].Key
				}
				return items[i].QueriesCount > items[j].QueriesCount
			})
		}
	}

	if perPage != 0 {
		if cursor == 0 {
			items = utils.Paginate(1, perPage, items)
		} else {
			items = utils.Paginate(cursor, perPage, items)
		}
	}

	return ctx.JSON(http.StatusOK, api.ListQueryParametersResponse{
		TotalCount: totalCount,
		Items:      items,
	})
}

// GetQueryParameter godoc
//
//	@Summary		Get query parameter
//	@Description	Returns the query parameter for the given key
//	@Security		BearerToken
//	@Tags			metadata
//	@Produce		json
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	models.PolicyParameterValues
//	@Router			/metadata/api/v1/query_parameter/{id} [get]
func (h *HttpHandler) GetQueryParameter(ctx echo.Context) error {
	key := ctx.Param("key")
	clientCtx := &httpclient.Context{UserRole: api3.AdminRole}

	var controls []complianceapi.Control
	var err error
	if h.complianceEnabled {
		controls, err = h.complianceClient.ListControl(clientCtx, nil, nil)
		if err != nil {
			h.logger.Error("error listing controls", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "error listing controls")
		}
	}
	namedQueries, err := h.ListQueriesV2Internal(api.ListQueryV2Request{})
	if err != nil {
		h.logger.Error("error listing queries", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "error listing queries")
	}

	queryParam, err := h.db.GetQueryParameter(key)
	if err != nil {
		h.logger.Error("error getting query parameters", zap.Error(err))
		return err
	}
	var controlsList []complianceapi.Control
	var queriesList []api.NamedQueryItemV2
	for _, c := range controls {
		for _, p := range c.Policy.Parameters {
			if p.Key == key {
				controlsList = append(controlsList, c)
			}

		}
	}
	for _, q := range namedQueries.Items {
		for _, p := range q.Query.Parameters {
			if p.Key == key {
				queriesList = append(queriesList, q)
			}
		}
	}
	var apiControlsList []api.Control
	for _, control := range controlsList {
		var parameters []api.ControlQueryParameter
		for _, parameter := range control.Policy.Parameters {
			parameters = append(parameters, api.ControlQueryParameter{
				Key:      parameter.Key,
				Required: parameter.Required,
			})
		}

		integrationTypes := make([]integration.Type, 0, len(control.Policy.IntegrationType))
		for _, it := range control.Policy.IntegrationType {
			integrationTypes = append(integrationTypes, integration.Type(it))
		}

		query := api.Policy{
			ID:              control.Policy.ID,
			Language:        api.PolicyLanguage(control.Policy.Language),
			Definition:      control.Policy.Definition,
			IntegrationType: integrationTypes,
			PrimaryResource: control.Policy.PrimaryResource,
			ListOfResources: control.Policy.ListOfResources,
			Parameters:      parameters,
			RegoPolicies:    control.Policy.RegoPolicies,
			CreatedAt:       control.Policy.CreatedAt,
			UpdatedAt:       control.Policy.UpdatedAt,
		}
		apiControlsList = append(apiControlsList, api.Control{
			ID:                      control.ID,
			Title:                   control.Title,
			Description:             control.Description,
			Tags:                    control.Tags,
			Explanation:             control.Explanation,
			NonComplianceCost:       control.NonComplianceCost,
			UsefulExample:           control.UsefulExample,
			CliRemediation:          control.CliRemediation,
			ManualRemediation:       control.ManualRemediation,
			GuardrailRemediation:    control.GuardrailRemediation,
			ProgrammaticRemediation: control.ProgrammaticRemediation,
			IntegrationType:         control.IntegrationType,
			Enabled:                 control.Enabled,
			DocumentURI:             control.DocumentURI,
			Policy:                  &query,
			Severity:                api.ComplianceResultSeverity(control.Severity),
			ManualVerification:      control.ManualVerification,
			Managed:                 control.Managed,
			CreatedAt:               control.CreatedAt,
			UpdatedAt:               control.UpdatedAt,
		})
	}
	return ctx.JSON(http.StatusOK, api.GetQueryParamDetailsResponse{
		Key:      key,
		Value:    queryParam.Value,
		Controls: apiControlsList,
		Queries:  queriesList,
	})
}

// PurgeSampleData godoc
//
//	@Summary		Purge all sample data
//	@Description	Purge all sample data
//	@Security		BearerToken
//	@Tags			workspace
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/workspace/api/v3/sample/purge [put]
func (h *HttpHandler) PurgeSampleData(c echo.Context) error {
	ctx := &httpclient.Context{UserRole: api3.AdminRole}

	loaded, err := h.SampleDataLoaded(c)
	if err != nil {
		return err
	}
	if loaded == false {
		return echo.NewHTTPError(http.StatusNotFound, "Workspace does not contain sample data")
	}

	integrations, err := h.integrationClient.PurgeSampleData(ctx)
	if err != nil {
		return err
	}

	err = h.schedulerClient.PurgeSampleData(ctx, integrations)
	if err != nil {
		return err
	}

	if h.complianceEnabled {
		err = h.complianceClient.PurgeSampleData(ctx)
		if err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}

// SyncDemo godoc
//
//	@Summary		Sync demo
//
//	@Description	Syncs demo with the git backend.
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Param			demo_data_s3_url	query	string	false	"Demo Data S3 URL"
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/metadata/api/v3/sample/sync [put]
func (h *HttpHandler) SyncDemo(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()

	var mig *model.Migration
	tx := h.migratorDb.ORM.Model(&model.Migration{}).Where("id = ?", model2.MigrationJobName).Find(&mig)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		h.logger.Error("failed to get migration", zap.Error(tx.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get migration")
	}

	if mig != nil && mig.ID == model2.MigrationJobName {
		h.logger.Info("last migration job", zap.Any("job", *mig))
		if mig.Status != "COMPLETED" && mig.UpdatedAt.After(time.Now().Add(-1*10*time.Minute)) {
			return echo.NewHTTPError(http.StatusBadRequest, "sync sample data already in progress")
		}
	}

	metadata, err := coreUtils.GetConfigMetadata(h.db, string(models.MetadataKeyCustomizationEnabled))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "config not found")
		}
		return err
	}

	cnf := metadata.GetCore()

	var enabled models.IConfigMetadata
	switch cnf.Type {
	case models.ConfigMetadataTypeString:
		enabled = &models.StringConfigMetadata{
			ConfigMetadata: cnf,
		}
	case models.ConfigMetadataTypeInt:
		intValue, err := strconv.ParseInt(cnf.Value, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "failed to parse int value")
		}
		enabled = &models.IntConfigMetadata{
			ConfigMetadata: cnf,
			Value:          int(intValue),
		}
	case models.ConfigMetadataTypeBool:
		boolValue, err := strconv.ParseBool(cnf.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to convert bool to int")
		}
		enabled = &models.BoolConfigMetadata{
			ConfigMetadata: cnf,
			Value:          boolValue,
		}
	case models.ConfigMetadataTypeJSON:
		enabled = &models.JSONConfigMetadata{
			ConfigMetadata: cnf,
			Value:          cnf.Value,
		}
	}

	if !enabled.GetValue().(bool) {
		return echo.NewHTTPError(http.StatusForbidden, "customization is not allowed")
	}

	demoDataS3URL := echoCtx.QueryParam("demo_data_s3_url")
	if demoDataS3URL != "" {
		// validate url
		_, err := url.ParseRequestURI(demoDataS3URL)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid url")
		}
		err = coreUtils.SetConfigMetadata(h.db, models.DemoDataS3URL, demoDataS3URL)
		if err != nil {
			h.logger.Error("set config metadata", zap.Error(err))
			return err
		}
	}

	var importDemoJob batchv1.Job
	err = h.kubeClient.Get(ctx, k8sclient.ObjectKey{
		Namespace: h.cfg.OpengovernanceNamespace,
		Name:      "import-es-demo-data",
	}, &importDemoJob)
	if err != nil {
		return err
	}

	err = h.kubeClient.Delete(ctx, &importDemoJob)
	if err != nil {
		return err
	}

	for {
		err = h.kubeClient.Get(ctx, k8sclient.ObjectKey{
			Namespace: h.cfg.OpengovernanceNamespace,
			Name:      "import-es-demo-data",
		}, &importDemoJob)
		if err != nil {
			if k8sclient.IgnoreNotFound(err) == nil {
				break
			}
			return err
		}

		time.Sleep(1 * time.Second)
	}

	importDemoJob.ObjectMeta = metav1.ObjectMeta{
		Name:      "import-es-demo-data",
		Namespace: h.cfg.OpengovernanceNamespace,
		Annotations: map[string]string{
			"helm.sh/hook":        "post-install,post-upgrade",
			"helm.sh/hook-weight": "0",
		},
	}
	importDemoJob.Spec.Selector = nil
	importDemoJob.Spec.Suspend = aws.Bool(false)
	importDemoJob.Spec.Template.ObjectMeta = metav1.ObjectMeta{}
	importDemoJob.Status = batchv1.JobStatus{}

	err = h.kubeClient.Create(ctx, &importDemoJob)
	if err != nil {
		return err
	}

	var importDemoDbJob batchv1.Job
	err = h.kubeClient.Get(ctx, k8sclient.ObjectKey{
		Namespace: h.cfg.OpengovernanceNamespace,
		Name:      "import-psql-demo-data",
	}, &importDemoDbJob)
	if err != nil {
		return err
	}

	err = h.kubeClient.Delete(ctx, &importDemoDbJob)
	if err != nil {
		return err
	}

	for {
		err = h.kubeClient.Get(ctx, k8sclient.ObjectKey{
			Namespace: h.cfg.OpengovernanceNamespace,
			Name:      "import-psql-demo-data",
		}, &importDemoDbJob)
		if err != nil {
			if k8sclient.IgnoreNotFound(err) == nil {
				break
			}
			return err
		}

		time.Sleep(1 * time.Second)
	}

	importDemoDbJob.ObjectMeta = metav1.ObjectMeta{
		Name:      "import-psql-demo-data",
		Namespace: h.cfg.OpengovernanceNamespace,
		Annotations: map[string]string{
			"helm.sh/hook":        "post-install,post-upgrade",
			"helm.sh/hook-weight": "0",
		},
	}
	importDemoDbJob.Spec.Selector = nil
	importDemoDbJob.Spec.Suspend = aws.Bool(false)
	importDemoDbJob.Spec.Template.ObjectMeta = metav1.ObjectMeta{}
	importDemoDbJob.Status = batchv1.JobStatus{}

	err = h.kubeClient.Create(ctx, &importDemoDbJob)
	if err != nil {
		return err
	}

	jp := pgtype.JSONB{}
	err = jp.Set([]byte(""))
	if err != nil {
		return err
	}
	tx = h.migratorDb.ORM.Model(&model.Migration{}).Where("id = ?", model2.MigrationJobName).Update("status", "Started").Update("jobs_status", jp)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		h.logger.Error("failed to update migration", zap.Error(tx.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update migration")
	}

	return echoCtx.JSON(http.StatusOK, struct{}{})
}

// WorkspaceLoadedSampleData godoc
//
//	@Summary		Sync demo
//
//	@Description	Syncs demo with the git backend.
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Param			demo_data_s3_url	query	string	false	"Demo Data S3 URL"
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/workspace/api/v3/sample/loaded [put]
func (h *HttpHandler) WorkspaceLoadedSampleData(echoCtx echo.Context) error {
	loaded, err := h.SampleDataLoaded(echoCtx)
	if err != nil {
		return err
	}

	if loaded {
		return echoCtx.String(http.StatusOK, "True")
	}
	return echoCtx.String(http.StatusOK, "False")
}

// GetMigrationStatus godoc
//
//	@Summary		Sync demo
//
//	@Description	Syncs demo with the git backend.
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Param			demo_data_s3_url	query	string	false	"Demo Data S3 URL"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	api.GetMigrationStatusResponse
//	@Router			/metadata/api/v3/migration/status [get]
func (h *HttpHandler) GetMigrationStatus(echoCtx echo.Context) error {
	var mig *model.Migration
	tx := h.migratorDb.ORM.Model(&model.Migration{}).Where("id = ?", "main").First(&mig)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		h.logger.Error("failed to get migration", zap.Error(tx.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get migration")
	}
	if mig == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no migration job found")
	}
	jobsStatus := make(map[string]model.JobInfo)

	if len(mig.JobsStatus.Bytes) > 0 {
		err := json.Unmarshal(mig.JobsStatus.Bytes, &jobsStatus)
		if err != nil {
			return err
		}
	}

	var completedJobs int
	for _, status := range jobsStatus {
		if status.Status == model.JobStatusCompleted || status.Status == model.JobStatusFailed {
			completedJobs++
		}
	}

	var jobProgress float64
	if len(jobsStatus) > 0 {
		jobProgress = float64(completedJobs) / float64(len(jobsStatus))
	}
	return echoCtx.JSON(http.StatusOK, api.GetMigrationStatusResponse{
		Status:     mig.Status,
		JobsStatus: jobsStatus,
		Summary: struct {
			TotalJobs          int     `json:"total_jobs"`
			CompletedJobs      int     `json:"completed_jobs"`
			ProgressPercentage float64 `json:"progress_percentage"`
		}{
			TotalJobs:          len(jobsStatus),
			CompletedJobs:      completedJobs,
			ProgressPercentage: jobProgress * 100,
		},
		UpdatedAt: mig.UpdatedAt,
		CreatedAt: mig.CreatedAt,
	})
}

// GetSampleSyncStatus godoc
//
//	@Summary		Sync demo
//
//	@Description	Syncs demo with the git backend.
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Param			demo_data_s3_url	query	string	false	"Demo Data S3 URL"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	api.GetSampleSyncStatusResponse
//	@Router			/workspace/api/v3/sample/sync/status [get]
func (h *HttpHandler) GetSampleSyncStatus(echoCtx echo.Context) error {
	var mig *model.Migration
	tx := h.migratorDb.ORM.Model(&model.Migration{}).Where("id = ?", model2.MigrationJobName).First(&mig)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		h.logger.Error("failed to get migration", zap.Error(tx.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get migration")
	}
	var jobsStatus model2.ESImportProgress

	if len(mig.JobsStatus.Bytes) > 0 {
		err := json.Unmarshal(mig.JobsStatus.Bytes, &jobsStatus)
		if err != nil {
			return err
		}
	}
	return echoCtx.JSON(http.StatusOK, api.GetSampleSyncStatusResponse{
		Status:   mig.Status,
		Progress: jobsStatus.Progress,
	})
}

// GetConfiguredStatus godoc
//
//	@Summary		Sync demo
//
//	@Description	Syncs demo with the git backend.
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Param			demo_data_s3_url	query	string	false	"Demo Data S3 URL"
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/workspace/api/v3/configured/status [get]
func (h *HttpHandler) GetConfiguredStatus(echoCtx echo.Context) error {
	appConfiguration, err := h.db.GetAppConfiguration()
	if err != nil {
		h.logger.Error("failed to get workspace", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get workspace")
	}

	if appConfiguration.Configured {
		return echoCtx.String(http.StatusOK, "True")
	} else {
		return echoCtx.String(http.StatusOK, "False")
	}
}

// SetConfiguredStatus godoc
//
//	@Summary		Sync demo
//
//	@Description	Syncs demo with the git backend.
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Param			demo_data_s3_url	query	string	false	"Demo Data S3 URL"
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/workspace/api/v3/configured/set [put]
func (h *HttpHandler) SetConfiguredStatus(echoCtx echo.Context) error {
	err := h.db.AppConfigured(true)
	if err != nil {
		h.logger.Error("failed to set workspace configured", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to set workspace configured")
	}
	return echoCtx.NoContent(http.StatusOK)
}

// UnsetConfiguredStatus godoc
//
//	@Summary		Sync demo
//
//	@Description	Syncs demo with the git backend.
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Param			demo_data_s3_url	query	string	false	"Demo Data S3 URL"
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/workspace/api/v3/configured/unset [put]
func (h *HttpHandler) UnsetConfiguredStatus(echoCtx echo.Context) error {
	err := h.db.AppConfigured(false)
	if err != nil {
		h.logger.Error("failed to unset workspace configured", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to unset workspace configured")
	}
	return echoCtx.NoContent(http.StatusOK)
}

// GetAbout godoc
//
//	@Summary		Get About info
//
//	@Description	Syncs demo with the git backend.
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	api.About
//	@Router			/workspace/api/v3/about [put]
func (h *HttpHandler) GetAbout(echoCtx echo.Context) error {
	ctx := &httpclient.Context{UserRole: api3.AdminRole}

	version := ""
	var opengovernanceVersionConfig corev1.ConfigMap
	err := h.kubeClient.Get(echoCtx.Request().Context(), k8sclient.ObjectKey{
		Namespace: h.cfg.OpengovernanceNamespace,
		Name:      "platform-version",
	}, &opengovernanceVersionConfig)
	if err == nil {
		version = opengovernanceVersionConfig.Data["version"]
	} else {
		fmt.Printf("failed to load version due to %v\n", err)
	}

	integrationURL := strings.ReplaceAll(h.cfg.Integration.BaseURL, "%NAMESPACE%", h.cfg.OpengovernanceNamespace)
	integrationClient := integrationClient.NewIntegrationServiceClient(integrationURL)
	integrationsResp, err := integrationClient.ListIntegrations(ctx, nil)
	if err != nil {
		h.logger.Error("failed to list integrations", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list integrations")
	}

	integrations := make(map[string][]integrationApi.Integration)
	for _, c := range integrationsResp.Integrations {
		if _, ok := integrations[c.IntegrationType.String()]; !ok {
			integrations[c.IntegrationType.String()] = make([]integrationApi.Integration, 0)
		}
		integrations[c.IntegrationType.String()] = append(integrations[c.IntegrationType.String()], c)
	}

	var engine api.QueryEngine
	engine = api.QueryEngineCloudQL
	query := `SELECT
    (SELECT SUM(cost) FROM azure_costmanagement_costbyresourcetype) +
    (SELECT SUM(amortized_cost_amount) FROM aws_cost_by_service_daily) AS total_cost;`
	var query_req = api.RunQueryRequest{
		Page: api.Page{
			No:   1,
			Size: 1000,
		},
		Engine: &engine,
		Query:  &query,
		Sorts:  nil,
	}

	results, err := h.RunQueryInternal(echoCtx, query_req)
	if err != nil {
		h.logger.Error("failed to run query", zap.Error(err))
	}

	var floatValue float64
	if results != nil {
		h.logger.Info("query result", zap.Any("result", results.Result))
		if len(results.Result) > 0 && len(results.Result[0]) > 0 {
			totalSpent := results.Result[0][0]
			floatValue, _ = totalSpent.(float64)
		}
	}

	var dexConnectors []api.DexConnectorInfo

	if h.dexClient != nil {
		dexRes, err := h.dexClient.ListConnectors(context.Background(), &dexApi.ListConnectorReq{})
		if err != nil {
			h.logger.Error("failed to list dex connectors", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "failed to list dex connectors")
		}
		if dexRes != nil {
			for _, c := range dexRes.Connectors {
				dexConnectors = append(dexConnectors, api.DexConnectorInfo{
					ID:   c.Id,
					Name: c.Name,
					Type: c.Type,
				})
			}
		}
	}

	loaded, err := h.SampleDataLoaded(echoCtx)
	if err != nil {
		h.logger.Error("failed to load data", zap.Error(err))
	}

	appConfiguration, err := h.db.GetAppConfiguration()
	if err != nil {
		h.logger.Error("failed to get workspace", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get workspace")
	}

	creationTime := time.Time{}
	if appConfiguration != nil {
		creationTime = appConfiguration.CreatedAt
	}

	response := api.About{
		InstallID:             appConfiguration.InstallID.String(),
		DexConnectors:         dexConnectors,
		AppVersion:            version,
		WorkspaceCreationTime: creationTime,
		PrimaryDomainURL:      h.cfg.PrimaryDomainURL,
		Integrations:          integrations,
		SampleData:            loaded,
		TotalSpendGoverned:    floatValue,
	}

	return echoCtx.JSON(http.StatusOK, response)
}

// GetAbout godoc
//
//	@Summary		Get About info
//
//	@Description	Syncs demo with the git backend.
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	api.About
//	@Router			/workspace/api/v3/about [put]
func (h *HttpHandler) GetAboutShort(echoCtx echo.Context) error {

	version := ""
	var opengovernanceVersionConfig corev1.ConfigMap
	err := h.kubeClient.Get(echoCtx.Request().Context(), k8sclient.ObjectKey{
		Namespace: h.cfg.OpengovernanceNamespace,
		Name:      "platform-version",
	}, &opengovernanceVersionConfig)
	if err == nil {
		version = opengovernanceVersionConfig.Data["version"]
	} else {
		fmt.Printf("failed to load version due to %v\n", err)
	}

	appConfiguration, err := h.db.GetAppConfiguration()
	if err != nil {
		h.logger.Error("failed to get workspace", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get workspace")
	}
	creationTime := time.Time{}
	if appConfiguration != nil {
		creationTime = appConfiguration.CreatedAt
	}

	response := api.About{
		InstallID:             appConfiguration.InstallID.String(),
		AppVersion:            version,
		WorkspaceCreationTime: creationTime,
		PrimaryDomainURL:      h.cfg.PrimaryDomainURL,
	}

	return echoCtx.JSON(http.StatusOK, response)
}

func newDexClient(hostAndPort string) (dexApi.DexClient, error) {
	conn, err := grpc.NewClient(hostAndPort, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}
	return dexApi.NewDexClient(conn), nil
}

func (h *HttpHandler) SampleDataLoaded(echoCtx echo.Context) (bool, error) {
	ctx := &httpclient.Context{UserRole: api3.AdminRole}

	integrationURL := strings.ReplaceAll(h.cfg.Integration.BaseURL, "%NAMESPACE%", h.cfg.OpengovernanceNamespace)
	integrationClient := integrationClient.NewIntegrationServiceClient(integrationURL)

	integrations, err := integrationClient.ListIntegrations(ctx, nil)
	if err != nil {
		h.logger.Error("failed to list integrations", zap.Error(err))
		return false, echo.NewHTTPError(http.StatusInternalServerError, "failed to list integrations")
	}

	loaded := false
	for _, integration := range integrations.Integrations {
		if integration.State == integrationApi.IntegrationStateSample {
			loaded = true
		}
	}

	return loaded, nil
}

// VaultConfigured godoc
//
//	@Summary		Get About info
//
//	@Description	Syncs demo with the git backend.
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	api.About
//	@Router			/workspace/api/v3/vault/configured [get]
func (h *HttpHandler) VaultConfigured(echoCtx echo.Context) error {

	return echoCtx.String(http.StatusOK, "True")
}

// ReloadViews godoc
//
//	@Summary		Reload views
//
//	@Description	Reloads the views
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/metadata/api/v3/views/reload [put]
func (h *HttpHandler) ReloadViews(echoCtx echo.Context) error {
	h.viewCheckpoint = time.Now()
	return echoCtx.NoContent(http.StatusOK)
}

// GetViewsCheckpoint godoc
//
//	@Summary		Get views checkpoint
//
//	@Description	Returns the views checkpoint
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	api.GetViewsCheckpointResponse
//	@Router			/core/api/v3/views/checkpoint [get]
func (h *HttpHandler) GetViewsCheckpoint(echoCtx echo.Context) error {
	return echoCtx.JSON(http.StatusOK, api.GetViewsCheckpointResponse{
		Checkpoint: h.viewCheckpoint,
	})
}

// GetViews godoc
//
//	@Summary		Get views
//
//	@Description	Returns the views
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Accept			json
//	@Produce		json
//	@Param			cursor		query		int	false	"Cursor"
//	@Param			per_page	query		int	false	"Per Page"
//	@Success		200			{object}	api.GetViewsResponse
//	@Router			/core/api/v3/views [get]
func (h *HttpHandler) GetViews(echoCtx echo.Context) error {
	views, err := h.db.ListQueryViews()
	if err != nil {
		h.logger.Error("failed to list views", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list views")
	}

	var cursor, perPage int64
	cursorStr := echoCtx.QueryParam("cursor")
	if cursorStr != "" {
		cursor, err = strconv.ParseInt(cursorStr, 10, 64)
		if err != nil {
			return err
		}
	}
	perPageStr := echoCtx.QueryParam("per_page")
	if perPageStr != "" {
		perPage, err = strconv.ParseInt(perPageStr, 10, 64)
		if err != nil {
			return err
		}
	}

	apiViews := make([]api.View, 0, len(views))
	for _, view := range views {
		var query api.Query
		if view.Query != nil {
			var parameters []api.QueryParameter
			for _, p := range view.Query.Parameters {
				parameters = append(parameters, api.QueryParameter{
					Key:      p.Key,
					Required: p.Required,
				})
			}
			query = api.Query{
				ID:             view.Query.ID,
				QueryToExecute: view.Query.QueryToExecute,
				PrimaryTable:   view.Query.PrimaryTable,
				ListOfTables:   view.Query.ListOfTables,
				Parameters:     parameters,
				Engine:         view.Query.Engine,
				Global:         view.Query.Global,
			}
		}

		apiView := api.View{
			ID:               view.ID,
			Title:            view.Title,
			Description:      view.Description,
			LastTimeRendered: h.viewCheckpoint,
			Query:            query,
			Dependencies:     view.Dependencies,
			Tags:             make(map[string][]string),
		}
		for _, tag := range view.Tags {
			apiView.Tags[tag.Key] = tag.Value
		}

		apiViews = append(apiViews, apiView)
	}

	totalCount := len(apiViews)
	sort.Slice(apiViews, func(i, j int) bool {
		return apiViews[i].ID < apiViews[j].ID
	})
	if perPage != 0 {
		if cursor == 0 {
			apiViews = utils.Paginate(1, perPage, apiViews)
		} else {
			apiViews = utils.Paginate(cursor, perPage, apiViews)
		}
	}

	return echoCtx.JSON(http.StatusOK, api.GetViewsResponse{
		Views:      apiViews,
		TotalCount: totalCount,
	})
}

// ReloadPluginSteampipeConfig godoc
//
//	@Summary		Update plugin steampipe binary file and config
//	@Description	Update plugin steampipe binary file and config
//	@Security		BearerToken
//	@Tags			compliance
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}
//	@Router			/core/api/v3/plugins/{plugin_id}/reload [put]
func (h *HttpHandler) ReloadPluginSteampipeConfig(echoCtx echo.Context) error {
	pluginId := echoCtx.Param("plugin_id")
	go func() {
		err := h.PluginJob.ReloadSinglePlugin(context.Background(), pluginId)
		if err != nil {
			h.logger.Error("failed to reload plugin", zap.String("plugin_id", pluginId), zap.Error(err))
		}
	}()
	return echoCtx.NoContent(http.StatusOK)
}

// RemovePluginSteampipeConfig godoc
//
//	@Summary		Remove plugin steampipe binary file and config
//	@Description	Remove plugin steampipe binary file and config
//	@Security		BearerToken
//	@Tags			compliance
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}
//	@Router			/core/api/v3/plugins/{plugin_id}/remove [put]
func (h *HttpHandler) RemovePluginSteampipeConfig(echoCtx echo.Context) error {
	pluginId := echoCtx.Param("plugin_id")
	go func() {
		err := h.PluginJob.RemoveSinglePlugin(context.Background(), pluginId)
		if err != nil {
			h.logger.Error("failed to reload plugin", zap.String("plugin_id", pluginId), zap.Error(err))
		}
	}()
	return echoCtx.NoContent(http.StatusOK)
}

// SyncQueries godoc
//
//	@Summary		Sync queries
//
//	@Description	Syncs queries with the git backend.
//
//	@Security		BearerToken
//	@Tags			compliance
//	@Param			configzGitURL	query	string	false	"Git URL"
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/core/api/v4/queries/sync [get]
func (h *HttpHandler) SyncQueries(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()

	var mig *model2.Migration
	tx := h.migratorDb.ORM.Model(&model2.Migration{}).Where("id = ?", "main").First(&mig)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		h.logger.Error("failed to get migration", zap.Error(tx.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get migration")
	}
	if mig != nil {
		if mig.Status == "PENDING" || mig.Status == "IN_PROGRESS" {
			return echo.NewHTTPError(http.StatusBadRequest, "sync sample data already in progress")
		}
	}

	enabled, err := coreUtils.GetConfigMetadata(h.db, string(models.MetadataKeyCustomizationEnabled))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "config not found")
		}
		return err
	}

	if !enabled.GetValue().(bool) {
		return echo.NewHTTPError(http.StatusForbidden, "customization is not allowed")
	}

	configzGitURL := echoCtx.QueryParam("configzGitURL")
	if configzGitURL != "" {
		// validate url
		_, err := url.ParseRequestURI(configzGitURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid url")
		}

		err = coreUtils.SetConfigMetadata(h.db, models.MetadataKeyAnalyticsGitURL, configzGitURL)
		if err != nil {
			return err
		}
	}

	var migratorJob batchv1.Job
	err = h.kubeClient.Get(ctx, k8sclient.ObjectKey{
		Namespace: h.cfg.OpengovernanceNamespace,
		Name:      "post-install-configuration",
	}, &migratorJob)
	if err != nil {
		return err
	}

	err = h.kubeClient.Delete(ctx, &migratorJob)
	if err != nil {
		return err
	}
	envsMap := make(map[string]corev1.EnvVar)
	for _, env := range migratorJob.Spec.Template.Spec.Containers[0].Env {
		envsMap[env.Name] = env
	}
	envsMap["IS_MANUAL"] = corev1.EnvVar{
		Name:  "IS_MANUAL",
		Value: "true",
	}
	var newEnvs []corev1.EnvVar
	for _, v := range envsMap {
		newEnvs = append(newEnvs, v)
	}
	for {
		err = h.kubeClient.Get(ctx, k8sclient.ObjectKey{
			Namespace: h.cfg.OpengovernanceNamespace,
			Name:      "post-install-configuration",
		}, &migratorJob)
		if err != nil {
			if k8sclient.IgnoreNotFound(err) == nil {
				break
			}
			return err
		}

		time.Sleep(1 * time.Second)
	}

	migratorJob.ObjectMeta = metav1.ObjectMeta{
		Name:      "post-install-configuration",
		Namespace: h.cfg.OpengovernanceNamespace,
		Annotations: map[string]string{
			"helm.sh/hook":        "post-install,post-upgrade",
			"helm.sh/hook-weight": "0",
		},
	}
	migratorJob.Spec.Selector = nil
	migratorJob.Spec.Suspend = aws.Bool(false)
	migratorJob.Spec.Template.ObjectMeta = metav1.ObjectMeta{}
	migratorJob.Spec.Template.Spec.Containers[0].Env = newEnvs
	migratorJob.Status = batchv1.JobStatus{}

	err = h.kubeClient.Create(ctx, &migratorJob)
	if err != nil {
		return err
	}

	jp := pgtype.JSONB{}
	err = jp.Set([]byte(""))
	if err != nil {
		return err
	}
	tx = h.migratorDb.ORM.Model(&model2.Migration{}).Where("id = ?", "main").Update("status", "Started").Update("jobs_status", jp)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		h.logger.Error("failed to update migration", zap.Error(tx.Error))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update migration")
	}

	return echoCtx.JSON(http.StatusOK, struct{}{})
}

func (h *HttpHandler) GetUserLayouts(echoCtx echo.Context) error {
	var req api.GetUserLayoutRequest
	if err := bindValidate(echoCtx, &req); err != nil {
		return err
	}
	userId := req.UserID
	layouts, err := h.db.GetUserLayouts(userId)
	if err != nil {
		h.logger.Error("failed to get user layout", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user layout")
	}
	if layouts == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user layout not found")
	}
	var response []api.GetUserLayoutResponse
	for _, layout := range layouts {
		response = append(response, api.GetUserLayoutResponse{
			ID:     layout.ID,
			UserID: userId,
			LayoutConfig: func() []map[string]any {
				var config []map[string]any
				if err := json.Unmarshal(layout.LayoutConfig.Bytes, &config); err != nil {
					h.logger.Error("failed to unmarshal layout config", zap.Error(err))
					return nil
				}
				return config
			}(),
			Name:        layout.Name,
			Description: layout.Description,
			IsDefault:   layout.IsDefault,
			IsPrivate:   layout.IsPrivate,
			UpdatedAt:   layout.UpdatedAt,
		})
	}

	return echoCtx.JSON(http.StatusOK, response)

}

func (h *HttpHandler) GetUserDefaultLayout(echoCtx echo.Context) error {
	var req api.GetUserLayoutRequest
	if err := bindValidate(echoCtx, &req); err != nil {
		return err
	}
	userId := req.UserID
	layout, err := h.db.GetUserDefaultLayout(userId)
	if err != nil {
		h.logger.Error("failed to get user layout", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user layout")
	}
	if layout == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user layout not found")
	}
	return echoCtx.JSON(http.StatusOK, api.GetUserLayoutResponse{
		ID:     layout.ID,
		UserID: userId,
		LayoutConfig: func() []map[string]any {
			var config []map[string]any
			if err := json.Unmarshal(layout.LayoutConfig.Bytes, &config); err != nil {
				h.logger.Error("failed to unmarshal layout config", zap.Error(err))
				return nil
			}
			return config
		}(),
		Name:        layout.Name,
		Description: layout.Description,
		IsPrivate:   layout.IsPrivate,
		IsDefault:   layout.IsDefault,
		UpdatedAt:   layout.UpdatedAt,
	})

}
func (h *HttpHandler) SetUserLayout(echoCtx echo.Context) error {
	var req api.SetUserLayoutRequest
	if err := bindValidate(echoCtx, &req); err != nil {
		return err
	}
	userId := req.UserID
	// convert req.layoutconfig to jsonb
	layout, err := json.Marshal(req.LayoutConfig)
	if err != nil {
		h.logger.Error("failed to marshal layout config", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to marshal layout config")
	}
	layout_config := pgtype.JSONB{}
	layout_config.Set(layout)
	var id string
	if req.ID != "" {
		id = req.ID
	} else {
		id = uuid.New().String()
	}

	user_layout := models.UserLayout{
		ID:           id,
		UserID:       userId,
		LayoutConfig: layout_config,
		IsDefault:    req.IsDefault,
		Name:         req.Name,
		Description:  req.Description,
		IsPrivate:    req.IsPrivate,
		UpdatedAt:    time.Now(),
	}
	err = h.db.SetUserLayout(user_layout)
	if err != nil {
		h.logger.Error("failed to set user layout", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to set user layout")
	}
	return echoCtx.NoContent(http.StatusOK)
}

func (h *HttpHandler) ChangePrivacy(echoCtx echo.Context) error {
	var req api.ChangePrivacyRequest
	if err := bindValidate(echoCtx, &req); err != nil {
		return err
	}
	userId := req.UserID
	privacy := req.IsPrivate
	err := h.db.ChangeLayoutPrivacy(userId, privacy)
	if err != nil {
		h.logger.Error("failed to change privacy", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to change privacy")
	}
	return echoCtx.NoContent(http.StatusOK)
}

// get public layouts
func (h *HttpHandler) GetPublicLayouts(echoCtx echo.Context) error {

	layouts, err := h.db.GetPublicLayouts()
	if err != nil {
		h.logger.Error("failed to get public layouts", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get public layouts")
	}
	if layouts == nil {
		return echo.NewHTTPError(http.StatusNotFound, "public layouts not found")
	}
	// convert layouts to []api.GetUserLayoutResponse
	var layoutsResponse []api.GetUserLayoutResponse
	for _, layout := range layouts {
		layoutResponse := api.GetUserLayoutResponse{
			UserID:      layout.UserID,
			Name:        layout.Name,
			Description: layout.Description,
			UpdatedAt:   layout.UpdatedAt,
			LayoutConfig: func() []map[string]any {
				var config []map[string]any
				if err := json.Unmarshal(layout.LayoutConfig.Bytes, &config); err != nil {
					h.logger.Error("failed to unmarshal layout config", zap.Error(err))
					return nil
				}
				return config
			}(),
		}
		layoutsResponse = append(layoutsResponse, layoutResponse)
	}

	return echoCtx.JSON(http.StatusOK, layoutsResponse)
}
