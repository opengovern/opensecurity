package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/opengovern/og-util/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	Orm *gorm.DB
}

func (db Database) Initialize(ctx context.Context) error {
	err := db.Orm.WithContext(ctx).AutoMigrate(
		&Policy{},
		&PolicyParameter{},
		&Control{},
		&ControlTag{},
		&Benchmark{},
		&BenchmarkTag{},
		&BenchmarkAssignment{},
	)
	if err != nil {
		return err
	}

	return nil
}

// =========== Benchmarks ===========

func (db Database) ListBenchmarks(ctx context.Context) ([]Benchmark, error) {
	var s []Benchmark
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).Preload(clause.Associations).
		Find(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) ListBenchmarksBare(ctx context.Context) ([]Benchmark, error) {
	var s []Benchmark
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).Preload("Tags").
		Find(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return s, nil
}

// ListRootBenchmarks returns all benchmarks that are not children of any other benchmark
// is it important to note that this function does not return the children of the root benchmarks neither the controls
func (db Database) ListRootBenchmarks(ctx context.Context, tags map[string][]string) ([]Benchmark, error) {
	var benchmarks []Benchmark
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).Preload(clause.Associations).
		Where("NOT EXISTS (SELECT 1 FROM benchmark_children WHERE benchmark_children.child_id = benchmarks.id)")
	if len(tags) > 0 {
		tx = tx.Joins("JOIN benchmark_tags AS tags ON tags.benchmark_id = benchmarks.id")
		for key, values := range tags {
			if len(values) != 0 {
				tx = tx.Where("tags.key = ? AND (tags.value && ?)", key, pq.StringArray(values))
			} else {
				tx = tx.Where("tags.key = ?", key)
			}
		}
	}
	err := tx.Find(&benchmarks).Error
	if err != nil {
		return nil, err
	}

	return benchmarks, nil
}

// ListBenchmarksFiltered returns all benchmarks with the associated filters
func (db Database) ListBenchmarksFiltered(ctx context.Context, titleRegex *string, root bool, tags map[string][]string, parentBenchmarkId []string,
	assigned *bool, isBaseline *bool, integrationIds []string) ([]Benchmark, error) {
	var benchmarks []Benchmark
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).Preload(clause.Associations)

	if isBaseline != nil {
		tx = tx.Where("benchmarks.is_baseline = ?", *isBaseline)
	}
	if titleRegex != nil {
		tx = tx.Where("title ~* ?", *titleRegex)
	}
	if root {
		tx = tx.Where("NOT EXISTS (SELECT 1 FROM benchmark_children WHERE benchmark_children.child_id = benchmarks.id)")
	}
	if len(parentBenchmarkId) > 0 {
		tx = tx.Where("EXISTS (SELECT 1 FROM benchmark_children WHERE benchmark_children.child_id = benchmarks.id AND benchmark_children.benchmark_id IN ?)", parentBenchmarkId)
	}
	if assigned != nil {
		if *assigned {
			tx = tx.Where("(EXISTS (SELECT 1 FROM benchmark_assignments WHERE benchmark_assignments.benchmark_id = benchmarks.id) OR benchmarks.is_baseline = true)")
		} else {
			tx = tx.Where("(NOT EXISTS (SELECT 1 FROM benchmark_assignments WHERE benchmark_assignments.benchmark_id = benchmarks.id) OR benchmarks.is_baseline = true)")
		}
	}
	if len(integrationIds) > 0 {
		tx = tx.Where("(EXISTS (SELECT 1 FROM benchmark_assignments WHERE benchmark_assignments.integration_id IN ? AND benchmark_assignments.benchmark_id = benchmarks.id) OR benchmarks.is_baseline = true)",
			integrationIds)
	}

	if len(tags) > 0 {
		tx = tx.Joins("JOIN benchmark_tags AS tags ON tags.benchmark_id = benchmarks.id")
		for key, values := range tags {
			if len(values) != 0 {
				tx = tx.Where("tags.key = ? AND (tags.value && ?)", key, pq.StringArray(values))
			} else {
				tx = tx.Where("tags.key = ?", key)
			}
		}
	}
	err := tx.Find(&benchmarks).Error
	if err != nil {
		return nil, err
	}

	return benchmarks, nil
}

func (db Database) ListRootBenchmarksWithSubtreeControls(ctx context.Context, tags map[string][]string) ([]Benchmark, error) {
	var benchmarks []Benchmark

	allBenchmarks, err := db.ListBenchmarks(ctx)
	if err != nil {
		return nil, err
	}
	allBenchmarksMap := make(map[string]Benchmark)
	for _, b := range allBenchmarks {
		allBenchmarksMap[b.ID] = b
	}

	var populateControls func(ctx context.Context, benchmark *Benchmark) error
	populateControls = func(ctx context.Context, benchmark *Benchmark) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if benchmark == nil {
			return nil
		}
		if len(benchmark.Children) > 0 {
			for _, child := range benchmark.Children {
				child := allBenchmarksMap[child.ID]
				err := populateControls(ctx, &child)
				if err != nil {
					return err
				}
				for _, control := range child.Controls {
					found := false
					for _, c := range benchmark.Controls {
						if c.ID == control.ID {
							found = true
							break
						}
					}
					if !found {
						benchmark.Controls = append(benchmark.Controls, control)
					}
				}
			}
		}
		return nil
	}

	rootBenchmarks, err := db.ListRootBenchmarks(ctx, tags)
	if err != nil {
		return nil, err
	}

	for _, b := range rootBenchmarks {
		err := populateControls(ctx, &b)
		if err != nil {
			return nil, err
		}
		benchmarks = append(benchmarks, b)
	}

	return benchmarks, nil
}

func (db Database) GetFramework(ctx context.Context, frameworkId string) (*Benchmark, error) {
	var s Benchmark
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).Preload(clause.Associations).
		Where("id = ?", frameworkId).
		First(&s)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &s, nil
}

func (db Database) GetFrameworks(ctx context.Context, frameworkIds []string) ([]Benchmark, error) {
	var s []Benchmark
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).Preload(clause.Associations).
		Where("id IN ?", frameworkIds).
		First(&s)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return s, nil
}

func (db Database) GetFrameworkWithControlQueries(ctx context.Context, frameworkId string) (*Benchmark, error) {
	var s Benchmark
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).Preload(clause.Associations).Preload("Controls").Preload("Controls.Policy").
		Where("id = ?", frameworkId).
		First(&s)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &s, nil
}

func (db Database) GetFrameworkBare(ctx context.Context, benchmarkId string) (*Benchmark, error) {
	var s Benchmark
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).Preload("Tags").
		Where("id = ?", benchmarkId).
		First(&s)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &s, nil
}

func (db Database) GetFrameworksBare(ctx context.Context, benchmarkIds []string) ([]Benchmark, error) {
	var s []Benchmark
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).Preload("Tags").
		Where("id in ?", benchmarkIds).
		Find(&s)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return s, nil
}

func (db Database) SetFrameworkAutoAssign(ctx context.Context, benchmarkId string, isBaseline bool) error {
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).Where("id = ?", benchmarkId).Update("is_baseline", isBaseline)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) SetFrameworkEnabled(ctx context.Context, benchmarkId string, enabled bool) error {
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).Where("id = ?", benchmarkId).Update("enabled", enabled)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) ListDistinctRootFrameworksFromControlIds(ctx context.Context, controlIds []string) ([]Benchmark, error) {
	s := make(map[string]Benchmark)

	findControls := make(map[string]struct{})
	for _, controlId := range controlIds {
		findControls[controlId] = struct{}{}
	}

	rootBenchmarksWithControls, err := db.ListRootBenchmarksWithSubtreeControls(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, b := range rootBenchmarksWithControls {
		for _, c := range b.Controls {
			if _, ok := findControls[c.ID]; ok {
				s[b.ID] = b
				break
			}
		}
	}

	var res []Benchmark
	for _, b := range s {
		res = append(res, b)
	}

	return res, nil
}

func (db Database) GetPolicy(ctx context.Context, policyID string) (*Policy, error) {
	var s Policy
	tx := db.Orm.WithContext(ctx).Model(&Policy{}).Preload(clause.Associations).
		Where("id = ?", policyID).
		First(&s)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &s, nil
}

func (db Database) GetFrameworksTitle(ctx context.Context, ds []string) (map[string]string, error) {
	var bs []Benchmark
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).
		Where("id in ?", ds).
		Select("id, title").
		Find(&bs)

	if tx.Error != nil {
		return nil, tx.Error
	}

	res := map[string]string{}
	for _, b := range bs {
		res[b.ID] = b.Title
	}
	return res, nil
}

func (db Database) GetControlsTitle(ctx context.Context, ds []string) (map[string]string, error) {
	var bs []Control
	tx := db.Orm.WithContext(ctx).Model(&Control{}).
		Where("id in ?", ds).
		Select("id, title").
		Find(&bs)

	if tx.Error != nil {
		return nil, tx.Error
	}

	res := map[string]string{}
	for _, b := range bs {
		res[b.ID] = b.Title
	}
	return res, nil
}

// =========== Control ===========

func (db Database) GetControl(ctx context.Context, id string) (*Control, error) {
	var s Control
	tx := db.Orm.WithContext(ctx).Model(&Control{}).
		Preload(clause.Associations).
		Where("id = ?", id).
		First(&s)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	if s.PolicyID != nil {
		var policy Policy
		tx := db.Orm.WithContext(ctx).Model(&Policy{}).Preload(clause.Associations).Where("id = ?", *s.PolicyID).First(&policy)
		if tx.Error != nil {
			return nil, tx.Error
		}
		s.Policy = &policy
	}

	return &s, nil
}

func (db Database) GetFrameworkParent(ctx context.Context, benchmarkID string) (string, error) {
	var benchmarkIDs string

	tx := db.Orm.WithContext(ctx).
		Model(&BenchmarkChild{}).
		Select("benchmark_id").
		Where("child_id = ?", benchmarkID).
		Find(&benchmarkIDs)

	if tx.Error != nil {
		return "", tx.Error
	}

	return benchmarkIDs, nil
}

func (db Database) GetFrameworkIdsByControlID(ctx context.Context, controlID string) ([]string, error) {
	var benchmarkIDs []string

	tx := db.Orm.WithContext(ctx).
		Model(&BenchmarkControls{}).
		Select("benchmark_id").
		Where("control_id = ?", controlID).
		Find(&benchmarkIDs)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return benchmarkIDs, nil
}

func (db Database) ListControlsByFrameworkID(ctx context.Context, benchmarkID string) ([]Control, error) {
	var s []Control
	tx := db.Orm.WithContext(ctx).Model(&Control{}).
		Preload("Tags").
		Preload("Benchmarks").
		Where(Control{Benchmarks: []Benchmark{{ID: benchmarkID}}}).Find(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}

	policyIDs := make([]string, 0, len(s))
	for _, control := range s {
		if control.PolicyID != nil {
			policyIDs = append(policyIDs, *control.PolicyID)
		}
	}
	var queriesMap map[string]Policy
	if len(policyIDs) > 0 {
		var policies []Policy
		qtx := db.Orm.WithContext(ctx).Model(&Policy{}).Preload(clause.Associations).Where("id IN ?", policyIDs).Find(&policies)
		if qtx.Error != nil {
			return nil, qtx.Error
		}
		queriesMap = make(map[string]Policy)
		for _, policy := range policies {
			queriesMap[policy.ID] = policy
		}
	}

	for i, c := range s {
		if c.PolicyID != nil {
			v := queriesMap[*c.PolicyID]
			s[i].Policy = &v
		}
	}

	return s, nil
}
func (db Database) ListControlsByFilter(ctx context.Context, controlIDs, integrationTypes []string, severity []string, benchmarkIDs []string,
	tagFilters map[string][]string, hasParameters *bool, primaryResource []string, listOfResources []string, params []string) ([]Control, error) {
	var s []Control

	m := db.Orm.WithContext(ctx).Model(&Control{}).Distinct("controls.*").
		Preload("Tags").
		Preload("Benchmarks")

	// Add filtering by tag keys and values if any filters are provided
	if len(tagFilters) > 0 {
		i := 0
		for key, values := range tagFilters {
			// Generate unique alias for each join to avoid alias collision
			alias := fmt.Sprintf("t%d", i)
			joinCondition := fmt.Sprintf("JOIN control_tags %s ON %s.control_id = controls.id", alias, alias)

			// Use PostgreSQL array operator @> to filter by tag values (if array comparison is required)
			m = m.Joins(joinCondition).Where(fmt.Sprintf("%s.key = ? AND %s.value::text[] @> ?", alias, alias), key, pq.Array(values))

			i++ // Increment the alias index
		}
	}

	for i, c := range integrationTypes {
		integrationTypes[i] = strings.ToLower(c)
	}

	if len(integrationTypes) > 0 {
		m = m.Where("controls.integration_type::text[] && ?", pq.Array(integrationTypes))
	}

	if len(severity) > 0 {
		m = m.Where("controls.severity IN ?", severity)
	}

	if len(controlIDs) > 0 {
		m = m.Where("controls.id IN ?", controlIDs)
	}

	if len(benchmarkIDs) > 0 {
		m = m.Joins("JOIN benchmark_controls bc ON bc.control_id = controls.id").
			Where("bc.benchmark_id IN ?", benchmarkIDs)
	}

	if hasParameters != nil || len(params) > 0 || len(primaryResource) > 0 || len(listOfResources) > 0 {
		m = m.Joins("JOIN policies p ON p.id = controls.policy_id")
	}

	if hasParameters != nil {
		if *hasParameters {
			m = m.Joins("LEFT JOIN policy_parameters pp ON pp.policy_id = p.id").
				Group("controls.id").
				Having("COUNT(pp.policy_id) > 0")
		} else {
			m = m.Joins("LEFT JOIN policy_parameters pp ON pp.policy_id = p.id").
				Group("controls.id").
				Having("COUNT(pp.policy_id) = 0")
		}
	}

	if len(params) > 0 {
		m = m.Joins("LEFT JOIN policy_parameters pp ON pp.policy_id = p.id").
			Where("pp.key IN ?", params).
			Group("controls.id").
			Having("COUNT(pp.policy_id) > 0")
	}

	if len(primaryResource) > 0 {
		m = m.Where("p.primary_resource IN ?", primaryResource)
	}

	if len(listOfResources) > 0 {
		m = m.Where("p.list_of_resources::text[] && ?", pq.Array(listOfResources))
	}

	// Execute the policy
	tx := m.Find(&s)

	if tx.Error != nil {
		return nil, tx.Error
	}

	v := map[string]Control{}
	for _, item := range s {
		if _, ok := v[item.ID]; !ok {
			v[item.ID] = item
		}
	}

	var res []Control
	for _, val := range v {
		res = append(res, val)
	}

	policyIDs := make([]string, 0, len(res))
	for _, control := range res {
		if control.PolicyID != nil {
			policyIDs = append(policyIDs, *control.PolicyID)
		}
	}
	var queriesMap map[string]Policy
	if len(policyIDs) > 0 {
		var policies []Policy
		qtx := db.Orm.WithContext(ctx).Model(&Policy{}).Preload(clause.Associations).Where("id IN ?", policyIDs).Find(&policies)
		if qtx.Error != nil {
			return nil, qtx.Error
		}
		queriesMap = make(map[string]Policy)
		for _, policy := range policies {
			queriesMap[policy.ID] = policy
		}
	}

	for i, c := range res {
		if c.PolicyID != nil {
			v := queriesMap[*c.PolicyID]
			res[i].Policy = &v
		}
	}

	return res, nil
}

func (db Database) GetControlsTags() ([]ControlTagsResult, error) {
	var results []ControlTagsResult

	// Execute the raw SQL query
	query := `SELECT 
    key, 
    ARRAY_AGG(DISTINCT value) AS unique_values
FROM (
    SELECT 
        key, 
        UNNEST(value) AS value
    FROM control_tags
) AS expanded_values
GROUP BY key;
`
	err := db.Orm.Raw(query).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (db Database) GetFrameworksTags() ([]BenchmarkTagsResult, error) {
	var results []BenchmarkTagsResult

	// Execute the raw SQL query
	query := `SELECT 
    key, 
    ARRAY_AGG(DISTINCT value) AS unique_values
FROM (
    SELECT 
        key, 
        UNNEST(value) AS value
    FROM benchmark_tags
) AS expanded_values
GROUP BY key;
`
	err := db.Orm.Raw(query).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (db Database) GetControls(ctx context.Context, controlIDs []string, tags map[string][]string) ([]Control, error) {
	var s []Control
	tx := db.Orm.WithContext(ctx).Model(&Control{}).
		Preload(clause.Associations)
	if len(controlIDs) > 0 {
		tx = tx.Where("id IN ?", controlIDs)
	}
	if len(tags) > 0 {
		tx = tx.Joins("JOIN control_tags AS tags ON tags.control_id = controls.id")
		for key, values := range tags {
			if len(values) != 0 {
				tx = tx.Where("tags.key = ? AND (tags.value && ?)", key, pq.StringArray(values))
			} else {
				tx = tx.Where("tags.key = ?", key)
			}
		}
	}
	if tx.Find(&s).Error != nil {
		return nil, tx.Error
	}

	policyIDs := make([]string, 0, len(s))
	for _, control := range s {
		if control.PolicyID != nil {
			policyIDs = append(policyIDs, *control.PolicyID)
		}
	}
	var queriesMap map[string]Policy
	if len(policyIDs) > 0 {
		var policies []Policy
		qtx := db.Orm.WithContext(ctx).Model(&Policy{}).Preload(clause.Associations).Where("id IN ?", policyIDs).Find(&policies)
		if qtx.Error != nil {
			return nil, qtx.Error
		}
		queriesMap = make(map[string]Policy)
		for _, policy := range policies {
			queriesMap[policy.ID] = policy
		}
	}

	for i, c := range s {
		if c.PolicyID != nil {
			v := queriesMap[*c.PolicyID]
			s[i].Policy = &v
		}
	}

	return s, nil
}

func (db Database) GetPolicies(ctx context.Context, policyIDs []string) ([]Policy, error) {
	var s []Policy
	tx := db.Orm.WithContext(ctx).Model(&Policy{}).
		Where("id IN ?", policyIDs).
		Find(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) GetPoliciesIdAndIntegrationType(ctx context.Context, policyIDs []string) ([]Policy, error) {
	var s []Policy
	tx := db.Orm.WithContext(ctx).Model(&Policy{}).
		Select("id, integration_type").
		Where("id IN ?", policyIDs).
		Find(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return s, nil
}

// =========== BenchmarkAssignment ===========

func (db Database) CleanupAllBenchmarkAssignments() error {
	tx := db.Orm.Where("1 = 1").Unscoped().Delete(&BenchmarkAssignment{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) AddBenchmarkAssignment(ctx context.Context, assignment *BenchmarkAssignment) error {
	tx := db.Orm.WithContext(ctx).Where(BenchmarkAssignment{
		BenchmarkId:   assignment.BenchmarkId,
		IntegrationID: assignment.IntegrationID,
	}).FirstOrCreate(assignment)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) GetBenchmarkAssignmentsByIntegrationId(ctx context.Context, integrationId string) ([]BenchmarkAssignment, error) {
	var s []BenchmarkAssignment
	tx := db.Orm.WithContext(ctx).Model(&BenchmarkAssignment{}).
		Where(BenchmarkAssignment{IntegrationID: &integrationId}).
		Where("resource_collection IS NULL").Scan(&s)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return s, nil
}

func (db Database) GetBenchmarkAssignmentsByBenchmarkId(ctx context.Context, benchmarkId string) ([]BenchmarkAssignment, error) {
	var s []BenchmarkAssignment
	tx := db.Orm.WithContext(ctx).Model(&BenchmarkAssignment{}).Where(BenchmarkAssignment{BenchmarkId: benchmarkId}).Scan(&s)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return s, nil
}

func (db Database) ListBenchmarkAssignments(ctx context.Context) ([]BenchmarkAssignment, error) {
	var s []BenchmarkAssignment
	tx := db.Orm.WithContext(ctx).Model(&BenchmarkAssignment{}).Scan(&s)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return s, nil
}

func (db Database) GetBenchmarkAssignmentByIds(ctx context.Context, benchmarkId string, integrationId *string) (*BenchmarkAssignment, error) {
	var s BenchmarkAssignment
	tx := db.Orm.WithContext(ctx).Model(&BenchmarkAssignment{}).Where(BenchmarkAssignment{
		BenchmarkId:   benchmarkId,
		IntegrationID: integrationId,
	}).Find(&s)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &s, nil
}

func (db Database) GetBenchmarkAssignmentsCount() ([]BenchmarkAssignmentsCount, error) {
	var results []BenchmarkAssignmentsCount
	tx := db.Orm.Table("benchmark_assignments").
		Select("benchmark_id, COUNT(*) as count").
		Group("benchmark_id").
		Scan(&results)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return results, nil
}

func (db Database) DeleteBenchmarkAssignmentByIds(ctx context.Context, benchmarkId string, integrationId *string) error {
	tx := db.Orm.WithContext(ctx).Unscoped().Where(BenchmarkAssignment{
		BenchmarkId:   benchmarkId,
		IntegrationID: integrationId,
	}).Delete(&BenchmarkAssignment{})

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) DeleteBenchmarkAssignmentByBenchmarkId(ctx context.Context, benchmarkId string) error {
	tx := db.Orm.WithContext(ctx).Where("benchmark_id = ?", benchmarkId).Delete(&BenchmarkAssignment{})

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) ListComplianceTagKeysWithPossibleValues(ctx context.Context) (map[string][]string, error) {
	var tags []BenchmarkTag
	tx := db.Orm.WithContext(ctx).Model(BenchmarkTag{}).Find(&tags)
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

func (db Database) ListControls(controlIDs []string, tags map[string][]string) ([]Control, error) {
	var s []Control
	tx := db.Orm.Model(&Control{}).
		Preload(clause.Associations).
		Preload("Policy.Parameters")

	if len(controlIDs) > 0 {
		tx = tx.Where("id IN ?", controlIDs)
	}
	if len(tags) > 0 {
		tx = tx.Joins("JOIN control_tags AS tags ON tags.control_id = controls.id")
		for key, values := range tags {
			if len(values) > 0 {
				tx = tx.Where("tags.key = ? AND (tags.value && ?)", key, pq.StringArray(values))
			} else {
				tx = tx.Where("tags.key = ?", key)
			}
		}
	}
	tx = tx.Find(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) ListPolicies(ctx context.Context) ([]Policy, error) {
	var s []Policy
	tx := db.Orm.WithContext(ctx).Model(&Policy{}).Preload(clause.Associations).
		Find(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) ListControlsBare(ctx context.Context) ([]Control, error) {
	var s []Control
	tx := db.Orm.WithContext(ctx).Model(&Control{}).Preload("Tags").
		Find(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) UpdateBenchmarkTrackDriftEvents(ctx context.Context, benchmarkId string, tracksDriftEvents bool) error {
	tx := db.Orm.WithContext(ctx).Model(&Benchmark{}).Where("id = ?", benchmarkId).Update("tracks_drift_events", tracksDriftEvents)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) ListControlsUniqueIntegrationTypes(ctx context.Context) ([]string, error) {
	var integrationTypes []string

	tx := db.Orm.WithContext(ctx).
		Model(&Control{}).
		Select("DISTINCT UNNEST(integration_type) AS unique_integration_types").
		Scan(&integrationTypes)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return integrationTypes, nil
}

func (db Database) ListControlsUniqueSeverity(ctx context.Context) ([]string, error) {
	var severities []string

	tx := db.Orm.WithContext(ctx).
		Model(&Control{}).
		Select("DISTINCT severity").
		Scan(&severities)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return severities, nil
}

func (db Database) ListControlsUniqueParentBenchmarks(ctx context.Context) ([]string, error) {
	var parentBenchmarks []string

	tx := db.Orm.WithContext(ctx).
		Model(&BenchmarkControls{}).
		Select("DISTINCT benchmark_id").
		Scan(&parentBenchmarks)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return parentBenchmarks, nil
}

func (db Database) ListPoliciesUniquePrimaryResources(ctx context.Context) ([]string, error) {
	var primaryTables []string

	tx := db.Orm.WithContext(ctx).
		Model(&Policy{}).
		Select("DISTINCT primary_resource").
		Where("primary_resource is not NULL").
		Scan(&primaryTables)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return primaryTables, nil
}

func (db Database) ListPolicyUniqueResources(ctx context.Context) ([]string, error) {
	var tables []string

	tx := db.Orm.WithContext(ctx).
		Model(&Policy{}).
		Select("DISTINCT UNNEST(list_of_resources)").
		Scan(&tables)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tables, nil
}

func (db Database) GetPolicyParameters(ctx context.Context) ([]string, error) {
	var parameters []string

	tx := db.Orm.WithContext(ctx).
		Select("DISTINCT key").
		Model(&PolicyParameter{}).
		Scan(&parameters)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return parameters, nil
}
