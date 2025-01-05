package db


import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-util/pkg/integration"
	"strings"
	"time"
	"github.com/lib/pq"
	"github.com/opengovern/og-util/pkg/model"
	"gorm.io/gorm/clause"
	"github.com/opengovern/opencomply/services/core/db/models"
)



func (db Database) GetQueriesWithFilters(search *string) ([]models.NamedQuery, error) {
	var s []models.NamedQuery

	m := db.orm.Model(&models.NamedQuery{})

	if search != nil {
		m = m.Where("title like ?", "%"+*search+"%")
	}
	tx := m.Find(&s)

	if tx.Error != nil {
		return nil, tx.Error
	}

	v := map[string]models.NamedQuery{}
	for _, item := range s {
		if _, ok := v[item.ID]; !ok {
			v[item.ID] = item
		}
	}
	var res []models.NamedQuery
	for _, val := range v {
		res = append(res, val)
	}

	for i, sq := range res {
		if sq.QueryID != nil {
			var query models.Query
			tx := db.orm.Model(&models.Query{}).Preload(clause.Associations).Where("id = ?", *sq.QueryID).First(&query)
			if tx.Error != nil {
				return nil, tx.Error
			}
			res[i].Query = &query
		}
	}

	return res, nil
}

func (db Database) ListQueries(queryIdsFilter []string, primaryTable []string, listOfTables []string, params []string, isBookmarked *bool) ([]models.NamedQuery, error) {
	var s []models.NamedQuery

	m := db.orm.Model(&models.NamedQuery{}).Distinct("named_queries.*")

	if len(queryIdsFilter) > 0 {
		m = m.Where("id in ?", queryIdsFilter)
	}

	if isBookmarked != nil {
		m = m.Where("is_bookmarked = ?", *isBookmarked)
	}

	if len(params) > 0 || len(primaryTable) > 0 || len(listOfTables) > 0 {
		m = m.Joins("JOIN queries q ON q.id = named_queries.query_id")
	}

	if len(params) > 0 {
		m = m.Joins("LEFT JOIN query_parameters qp ON qp.query_id = q.id").
			Where("qp.key IN ?", params).
			Group("named_queries.id").
			Having("COUNT(qp.query_id) > 0")
	}

	if len(primaryTable) > 0 {
		m = m.Where("q.primary_table IN ?", primaryTable)
	}

	if len(listOfTables) > 0 {
		m = m.Where("q.list_of_tables::text[] && ?", pq.Array(listOfTables))
	}

	tx := m.Find(&s)

	if tx.Error != nil {
		return nil, tx.Error
	}

	v := map[string]models.NamedQuery{}
	for _, item := range s {
		if _, ok := v[item.ID]; !ok {
			v[item.ID] = item
		}
	}

	var res []models.NamedQuery
	for _, val := range v {
		res = append(res, val)
	}

	queryIds := make([]string, 0, len(res))
	for _, control := range res {
		if control.QueryID != nil {
			queryIds = append(queryIds, *control.QueryID)
		}
	}
	var queriesMap map[string]models.Query
	if len(queryIds) > 0 {
		var queries []models.Query
		qtx := db.orm.Model(&models.Query{}).Preload(clause.Associations).Where("id IN ?", queryIds).Find(&queries)
		if qtx.Error != nil {
			return nil, qtx.Error
		}
		queriesMap = make(map[string]models.Query)
		for _, query := range queries {
			queriesMap[query.ID] = query
		}
	}

	for i, c := range res {
		if c.QueryID != nil {
			v := queriesMap[*c.QueryID]
			res[i].Query = &v
		}
	}
	return res, nil
}

func (db Database) GetQuery(id string) (*models.NamedQuery, error) {
	var s models.NamedQuery
	tx := db.orm.Model(models.NamedQuery{}).Preload(clause.Associations).Preload("Tags").Where("id = ?", id).First(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if s.QueryID != nil {
		var query models.Query
		tx := db.orm.Model(&models.Query{}).Preload(clause.Associations).Where("id = ?", *s.QueryID).First(&query)
		if tx.Error != nil {
			return nil, tx.Error
		}
		s.Query = &query
	}

	return &s, nil
}

func (db Database) ListQueriesByFilters(queryIds []string, search *string, tagFilters map[string][]string, integrationTypes []string,
	hasParameters *bool, primaryTable []string, listOfTables []string, params []string) ([]models.NamedQuery, error) {
	var s []models.NamedQuery

	m := db.orm.Model(&models.NamedQuery{}).Distinct("named_queries.*").Preload(clause.Associations).Preload("Policy.Parameters").Preload("Tags")

	if len(queryIds) > 0 {
		m = m.Where("id IN ?", queryIds)
	}

	if search != nil {
		m = m.Where("title LIKE ?", "%"+*search+"%")
	}

	for i, c := range integrationTypes {
		integrationTypes[i] = strings.ToLower(c)
	}

	if len(integrationTypes) > 0 {
		m = m.Where("named_queries.integration_types::text[] @> ?", pq.Array(integrationTypes))
	}

	if len(tagFilters) > 0 {
		i := 0
		for key, values := range tagFilters {
			alias := fmt.Sprintf("t%d", i)
			joinCondition := fmt.Sprintf("JOIN named_query_tags %s ON %s.named_query_id = named_queries.id", alias, alias)

			m = m.Joins(joinCondition).Where(fmt.Sprintf("%s.key = ? AND %s.value::text[] @> ?", alias, alias), key, pq.Array(values))

			i++
		}
	}

	if hasParameters != nil || len(params) > 0 || len(primaryTable) > 0 || len(listOfTables) > 0 {
		m = m.Joins("JOIN queries q ON q.id = named_queries.query_id")
	}

	if hasParameters != nil {
		if *hasParameters {
			m = m.Joins("LEFT JOIN query_parameters qp ON qp.query_id = q.id").
				Group("named_queries.id").
				Having("COUNT(qp.query_id) > 0")
		} else {
			m = m.Joins("LEFT JOIN query_parameters qp ON qp.query_id = q.id").
				Group("named_queries.id").
				Having("COUNT(qp.query_id) = 0")
		}
	}

	if len(params) > 0 {
		m = m.Joins("LEFT JOIN query_parameters qp ON qp.query_id = q.id").
			Where("qp.key IN ?", params).
			Group("named_queries.id").
			Having("COUNT(qp.query_id) > 0")
	}

	if len(primaryTable) > 0 {
		m = m.Where("q.primary_table IN ?", primaryTable)
	}

	if len(listOfTables) > 0 {
		m = m.Where("q.list_of_tables::text[] && ?", pq.Array(listOfTables))
	}

	tx := m.Find(&s)

	if tx.Error != nil {
		return nil, tx.Error
	}

	v := map[string]models.NamedQuery{}
	for _, item := range s {
		if _, ok := v[item.ID]; !ok {
			v[item.ID] = item
		}
	}

	var res []models.NamedQuery
	for _, val := range v {
		res = append(res, val)
	}

	for i, sq := range res {
		if sq.QueryID != nil {
			var query models.Query
			tx := db.orm.Model(&models.Query{}).Preload(clause.Associations).Where("id = ?", *sq.QueryID).First(&query)
			if tx.Error != nil {
				return nil, tx.Error
			}
			res[i].Query = &query
		}
	}

	return res, nil
}

func (db Database) GetQueriesTags() ([]models.NamedQueryTagsResult, error) {
	var results []models.NamedQueryTagsResult

	// Execute the raw SQL query
	query := `SELECT 
    key, 
    ARRAY_AGG(DISTINCT value) AS unique_values
FROM (
    SELECT 
        key, 
        UNNEST(value) AS value
    FROM named_query_tags
) AS expanded_values
GROUP BY key;
`
	err := db.orm.Raw(query).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (db Database) GetQueryHistory() ([]models.NamedQueryHistory, error) {
	var history []models.NamedQueryHistory
	tx := db.orm.Order("executed_at desc").Limit(3).Find(&history)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return history, nil
}

func (db Database) UpdateQueryHistory(query string) error {
	history := models.NamedQueryHistory{
		Query:      query,
		ExecutedAt: time.Now(),
	}
	// Upsert query history
	err := db.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "query"}},
		DoUpdates: clause.AssignmentColumns([]string{"executed_at"}),
	}).Create(&history).Error
	if err != nil {
		return err
	}

	// Only keep latest 100 queries in history
	const keepNumber = 100
	var count int64
	err = db.orm.Model(&models.NamedQueryHistory{}).Count(&count).Error
	if err != nil {
		return err
	}
	if count > keepNumber {
		var oldest models.NamedQueryHistory
		err = db.orm.Model(&models.NamedQueryHistory{}).Order("executed_at desc").Offset(keepNumber - 1).Limit(1).Find(&oldest).Error
		if err != nil {
			return err
		}

		err = db.orm.Model(&models.NamedQueryHistory{}).Where("executed_at < ?", oldest.ExecutedAt).Delete(&models.NamedQueryHistory{}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (db Database) ListResourceTypeTagsKeysWithPossibleValues(integrationTypes []integration.Type, doSummarize *bool) (map[string][]string, error) {
	var tags []models.ResourceTypeTag
	tx := db.orm.Model(models.ResourceTypeTag{}).Joins("JOIN resource_types ON resource_type_tags.resource_type = resource_types.resource_type")
	if doSummarize != nil {
		tx = tx.Where("resource_types.do_summarize = ?", true)
	}
	if len(integrationTypes) > 0 {
		tx = tx.Where("resource_types.integration_type in ?", integrationTypes)
	}
	tx.Find(&tags)
	if tx.Error != nil {
		return nil, tx.Error
	}
	tagLikes := make([]model.TagLike, 0, len(tags))
	for _, tag := range tags {
		tagLikes = append(tagLikes, tag)
	}
	result := model.GetTagsMap(tagLikes)
	return result, nil
}

func (db Database) ListFilteredResourceTypes(tags map[string][]string, resourceTypeNames []string, serviceNames []string, integrationTypes []integration.Type, doSummarize bool) ([]models.ResourceType, error) {
	var resourceTypes []models.ResourceType
	query := db.orm.Model(models.ResourceType{}).Preload(clause.Associations)
	if doSummarize {
		query = query.Where("resource_types.do_summarize = ?", doSummarize)
	}
	if len(tags) != 0 {
		query = query.Joins("JOIN resource_type_tags AS tags ON tags.resource_type = resource_types.resource_type")
		for key, values := range tags {
			if len(values) != 0 {
				query = query.Where("tags.key = ? AND (tags.value && ?)", key, pq.StringArray(values))
			} else {
				query = query.Where("tags.key = ?", key)
			}
		}
	}
	if len(serviceNames) != 0 {
		query = query.Where("service_name IN ?", serviceNames)
	}
	if len(integrationTypes) != 0 {
		query = query.Where("integration_type IN ?", integrationTypes)
	}
	if len(resourceTypeNames) != 0 {
		query = query.Where("resource_types.resource_type IN ?", resourceTypeNames)
	}
	tx := query.Find(&resourceTypes)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return resourceTypes, nil
}

func (db Database) GetResourceType(resourceType string) (*models.ResourceType, error) {
	var rtObj models.ResourceType
	tx := db.orm.Model(models.ResourceType{}).Preload(clause.Associations).Where("resource_type = ?", resourceType).First(&rtObj)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &rtObj, nil
}

func (db Database) ListResourceCollections(ids []string, statuses []models.ResourceCollectionStatus) ([]models.ResourceCollection, error) {
	var resourceCollections []models.ResourceCollection
	tx := db.orm.Model(models.ResourceCollection{}).Preload(clause.Associations)
	if len(ids) > 0 {
		tx = tx.Where("id IN ?", ids)
	}
	if len(statuses) > 0 {
		tx = tx.Where("status IN ?", statuses)
	}
	tx.Find(&resourceCollections)
	if tx.Error != nil {
		return nil, tx.Error
	}
	for i := range resourceCollections {
		if len(resourceCollections[i].FiltersJson.Bytes) > 0 {
			err := json.Unmarshal(resourceCollections[i].FiltersJson.Bytes, &resourceCollections[i].Filters)
			if err != nil {
				return nil, err
			}
		}
	}

	return resourceCollections, nil
}

func (db Database) GetResourceCollection(collectionID string) (*models.ResourceCollection, error) {
	var collection models.ResourceCollection
	tx := db.orm.Model(models.ResourceCollection{}).Preload(clause.Associations).Where("id = ?", collectionID).First(&collection)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if len(collection.FiltersJson.Bytes) > 0 {
		err := json.Unmarshal(collection.FiltersJson.Bytes, &collection.Filters)
		if err != nil {
			return nil, err
		}
	}

	return &collection, nil
}

func (db Database) ListNamedQueriesUniqueProviders() ([]string, error) {
	var integrationTypes []string

	tx := db.orm.
		Model(&models.NamedQuery{}).
		Select("DISTINCT UNNEST(integration_types)").
		Scan(&integrationTypes)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return integrationTypes, nil
}

func (db Database) ListResourceTypesUniqueCategories() ([]string, error) {
	var integrationTypes []string

	tx := db.orm.
		Model(&models.ResourceTypeV2{}).
		Select("DISTINCT category").
		Scan(&integrationTypes)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return integrationTypes, nil
}

func (db Database) ListCategoryResourceTypes(category string) ([]models.ResourceTypeV2, error) {
	var resourceTypes []models.ResourceTypeV2

	tx := db.orm.
		Model(&models.ResourceTypeV2{}).
		Where("category = ?", category)

	tx = tx.Find(&resourceTypes)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return resourceTypes, nil
}

func (db Database) ListResourceTypes(tables []string, categories []string) ([]models.ResourceTypeV2, error) {
	var resourceTypes []models.ResourceTypeV2

	tx := db.orm.
		Model(&models.ResourceTypeV2{})

	if len(tables) > 0 {
		tx = tx.Where("steampipe_table IN ?", tables)
	}

	if len(categories) > 0 {
		tx = tx.Where("category IN ?", categories)
	}

	tx = tx.Find(&resourceTypes)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return resourceTypes, nil
}

func (db Database) ListUniqueCategoriesAndTablesForTables(tables []string) ([]models.CategoriesTables, error) {
	var results []models.CategoriesTables

	query := `
        SELECT 
            json_build_object(
                'category', category,
                'tables', ARRAY_AGG(steampipe_table)
            ) AS category_tables
        FROM 
            resource_type_v2`

	var rows *sql.Rows
	var err error
	if len(tables) > 0 {
		query = query + ` WHERE steampipe_table IN ? GROUP BY category`
		rows, err = db.orm.Raw(query, tables).Rows()

	} else {
		query = query + ` GROUP BY category`
		rows, err = db.orm.Raw(query).Rows()

	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var jsonData []byte
		if err := rows.Scan(&jsonData); err != nil {
			return nil, err
		}

		var categoryTable models.CategoriesTables
		if err := json.Unmarshal(jsonData, &categoryTable); err != nil {
			return nil, err
		}

		results = append(results, categoryTable)
	}

	return results, nil
}

func (db Database) GetQueryParameters() ([]string, error) {
	var parameters []string

	tx := db.orm.Select("DISTINCT key").
		Model(&models.QueryParameter{}).
		Scan(&parameters)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return parameters, nil
}
