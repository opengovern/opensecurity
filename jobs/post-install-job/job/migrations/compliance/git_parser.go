package compliance

import (
	"encoding/json"
	"fmt"
	"github.com/opengovern/opencomply/jobs/post-install-job/config"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/shared"
	"github.com/opengovern/opencomply/jobs/post-install-job/utils"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/model"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/git"
	"github.com/opengovern/opencomply/pkg/types"
	"github.com/opengovern/opencomply/services/compliance/db"
	"github.com/opengovern/opencomply/services/core/db/models"
	"go.uber.org/zap"
)

type GitParser struct {
	logger             *zap.Logger
	benchmarks         []db.Benchmark
	frameworksChildren map[string][]string
	controls           []db.Control
	policies           []db.Policy
	policyParamValues  []models.PolicyParameterValues
	queryViews         []models.QueryView
	coreServiceQueries []models.Query
	controlsPolicies   map[string]db.Policy
	namedPolicies      map[string]NamedQuery
	Comparison         *git.ComparisonResultGrouped

	manualRemediationMap       map[string]string
	cliRemediationMap          map[string]string
	guardrailRemediationMap    map[string]string
	programmaticRemediationMap map[string]string
	noncomplianceCostMap       map[string]string
	usefulnessExampleMap       map[string]string
	explanationMap             map[string]string
}

func populateMdMapFromPath(path string) (map[string]string, error) {
	result := make(map[string]string)
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		id := strings.ToLower(strings.TrimSuffix(filepath.Base(path), ".md"))
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		result[id] = string(content)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (g *GitParser) ExtractNamedQueries() error {
	err := filepath.Walk(config.QueriesGitPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".yaml") {
			id := strings.TrimSuffix(info.Name(), ".yaml")

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var item NamedQuery
			err = yaml.Unmarshal(content, &item)
			if err != nil {
				g.logger.Error("failure in unmarshal", zap.String("path", path), zap.Error(err))
				return nil
			}

			if item.ID != "" {
				id = item.ID
			}

			g.namedPolicies[id] = item
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *GitParser) ExtractControls(complianceControlsPath string, controlEnrichmentBasePath string) error {
	var err error

	g.manualRemediationMap, err = populateMdMapFromPath(path.Join(controlEnrichmentBasePath, "tags", "remediation", "manual"))
	if err != nil {
		g.logger.Warn("failed to load manual remediation", zap.Error(err))
	}

	g.cliRemediationMap, err = populateMdMapFromPath(path.Join(controlEnrichmentBasePath, "tags", "remediation", "cli"))
	if err != nil {
		g.logger.Warn("failed to load cli remediation", zap.Error(err))
	}

	g.guardrailRemediationMap, err = populateMdMapFromPath(path.Join(controlEnrichmentBasePath, "tags", "remediation", "guardrail"))
	if err != nil {
		g.logger.Warn("failed to load cli remediation", zap.Error(err))
	}

	g.programmaticRemediationMap, err = populateMdMapFromPath(path.Join(controlEnrichmentBasePath, "tags", "remediation", "programmatic"))
	if err != nil {
		g.logger.Warn("failed to load cli remediation", zap.Error(err))
	}

	g.noncomplianceCostMap, err = populateMdMapFromPath(path.Join(controlEnrichmentBasePath, "tags", "noncompliance-cost"))
	if err != nil {
		g.logger.Warn("failed to load cli remediation", zap.Error(err))
	}

	g.usefulnessExampleMap, err = populateMdMapFromPath(path.Join(controlEnrichmentBasePath, "tags", "usefulness-example"))
	if err != nil {
		g.logger.Warn("failed to load cli remediation", zap.Error(err))
	}

	g.explanationMap, err = populateMdMapFromPath(path.Join(controlEnrichmentBasePath, "tags", "explanation"))
	if err != nil {
		g.logger.Warn("failed to load cli remediation", zap.Error(err))
	}

	return filepath.WalkDir(complianceControlsPath, func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".yaml") {
			content, err := os.ReadFile(path)
			if err != nil {
				g.logger.Error("failed to read yaml", zap.String("path", path), zap.Error(err))
				return err
			}

			var data map[string]interface{}
			if err := yaml.Unmarshal(content, &data); err != nil {
				g.logger.Error("failed to unmarshal yaml", zap.String("path", path), zap.Error(err))
				return fmt.Errorf("cannot parse YAML as map: %w", err)
			}

			if data["id"] != nil && data["policy"] != nil && data["severity"] != nil {
				if err = g.parseControlFile(content, path); err != nil {
					g.logger.Error("failed to parse control", zap.String("path", path), zap.Error(err))
					return err
				}
			} else if data["id"] != nil && data["definition"] != nil && data["language"] != nil {
				if err = g.parsePolicyFile(content, path); err != nil {
					g.logger.Error("failed to parse control", zap.String("path", path), zap.Error(err))
					return err
				}
			}
		}
		return nil
	})
}

func (g *GitParser) ExtractPolicies(compliancePoliciesPath string) error {
	return filepath.WalkDir(compliancePoliciesPath, func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".yaml") {
			content, err := os.ReadFile(path)
			if err != nil {
				g.logger.Error("failed to read yaml", zap.String("path", path), zap.Error(err))
				return err
			}

			var data map[string]interface{}
			if err := yaml.Unmarshal(content, &data); err != nil {
				g.logger.Error("failed to unmarshal yaml", zap.String("path", path), zap.Error(err))
				return fmt.Errorf("cannot parse YAML as map: %w", err)
			}

			if data["policy"] != nil && data["severity"] != nil {
				if err = g.parseControlFile(content, path); err != nil {
					g.logger.Error("failed to parse control", zap.String("path", path), zap.Error(err))
					return err
				}
			} else if data["definition"] != nil && data["language"] != nil {
				if err = g.parsePolicyFile(content, path); err != nil {
					g.logger.Error("failed to parse control", zap.String("path", path), zap.Error(err))
					return err
				}
			}
		}
		return nil
	})
}

func (g *GitParser) parsePolicyFile(content []byte, path string) error {
	var policy shared.Policy
	err := yaml.Unmarshal(content, &policy)
	if err != nil {
		g.logger.Error("failed to unmarshal policy", zap.String("path", path), zap.Error(err))
		return err
	}

	if policy.ID == nil {
		g.logger.Error("policy id should not be nil", zap.String("path", path))
		return fmt.Errorf("policy id should not be nil")
	}

	listOfTables, err := utils.ExtractTableRefsFromPolicy(policy.Language, policy.Definition)
	if err != nil {
		g.logger.Error("extract control failed: failed to extract table refs from query", zap.String("policy-id", *policy.ID), zap.Error(err))
		return nil
	}

	parameters, err := utils.ExtractParameters(policy.Language, policy.Definition)
	if err != nil {
		g.logger.Error("extract control failed: failed to extract parameters from query", zap.String("policy-id", *policy.ID), zap.Error(err))
		return nil
	}

	q := db.Policy{
		ID:              *policy.ID,
		Definition:      policy.Definition,
		PrimaryResource: policy.PrimaryResource,
		ListOfResources: listOfTables,
		Language:        policy.Language,
		RegoPolicies:    policy.RegoPolicies,
	}
	g.policies = append(g.policies, q)

	for _, parameter := range parameters {
		q.Parameters = append(q.Parameters, db.PolicyParameter{
			PolicyID: *policy.ID,
			Key:      parameter,
		})
	}

	return nil
}

func (g *GitParser) parseControlFile(content []byte, path string) error {
	var control Control
	err := yaml.Unmarshal(content, &control)
	if err != nil {
		g.logger.Error("failed to unmarshal control", zap.String("path", path), zap.Error(err))
		return err
	}
	tags := make([]db.ControlTag, 0, len(control.Tags))
	for tagKey, tagValue := range control.Tags {
		tags = append(tags, db.ControlTag{
			Tag: model.Tag{
				Key:   tagKey,
				Value: tagValue,
			},
			ControlID: control.ID,
		})
	}
	if v, ok := g.manualRemediationMap[strings.ToLower(control.ID)]; ok {
		tags = append(tags, db.ControlTag{
			Tag: model.Tag{
				Key:   "x-opengovernance-manual-remediation",
				Value: []string{v},
			},
			ControlID: control.ID,
		})
	}
	if v, ok := g.cliRemediationMap[strings.ToLower(control.ID)]; ok {
		tags = append(tags, db.ControlTag{
			Tag: model.Tag{
				Key:   "x-opengovernance-cli-remediation",
				Value: []string{v},
			},
			ControlID: control.ID,
		})
	}
	if v, ok := g.guardrailRemediationMap[strings.ToLower(control.ID)]; ok {
		tags = append(tags, db.ControlTag{
			Tag: model.Tag{
				Key:   "x-opengovernance-guardrail-remediation",
				Value: []string{v},
			},
			ControlID: control.ID,
		})
	}
	if v, ok := g.programmaticRemediationMap[strings.ToLower(control.ID)]; ok {
		tags = append(tags, db.ControlTag{
			Tag: model.Tag{
				Key:   "x-opengovernance-programmatic-remediation",
				Value: []string{v},
			},
			ControlID: control.ID,
		})
	}
	if v, ok := g.noncomplianceCostMap[strings.ToLower(control.ID)]; ok {
		tags = append(tags, db.ControlTag{
			Tag: model.Tag{
				Key:   "x-opengovernance-noncompliance-cost",
				Value: []string{v},
			},
			ControlID: control.ID,
		})
	}
	if v, ok := g.explanationMap[strings.ToLower(control.ID)]; ok {
		tags = append(tags, db.ControlTag{
			Tag: model.Tag{
				Key:   "x-opengovernance-explanation",
				Value: []string{v},
			},
			ControlID: control.ID,
		})
	}
	if v, ok := g.usefulnessExampleMap[strings.ToLower(control.ID)]; ok {
		tags = append(tags, db.ControlTag{
			Tag: model.Tag{
				Key:   "x-opengovernance-usefulness-example",
				Value: []string{v},
			},
			ControlID: control.ID,
		})
	}
	if control.Severity == "" {
		control.Severity = "low"
	}

	c := db.Control{
		ID:              control.ID,
		Title:           control.Title,
		Description:     control.Description,
		Tags:            tags,
		IntegrationType: control.IntegrationType,
		Enabled:         true,
		Benchmarks:      nil,
		Severity:        types.ParseComplianceResultSeverity(control.Severity),
	}

	if control.Policy != nil {
		if control.Policy.Ref != nil {
			c.PolicyID = control.Policy.Ref
		} else {
			listOfTables, err := utils.ExtractTableRefsFromPolicy(control.Policy.Language, control.Policy.Definition)
			if err != nil {
				g.logger.Error("extract control failed: failed to extract table refs from query", zap.String("control-id", control.ID), zap.Error(err))
				return nil
			}

			parameters, err := utils.ExtractParameters(control.Policy.Language, control.Policy.Definition)
			if err != nil {
				g.logger.Error("extract control failed: failed to extract parameters from query", zap.String("control-id", control.ID), zap.Error(err))
				return nil
			}

			q := db.Policy{
				ID:              control.ID,
				Definition:      control.Policy.Definition,
				IntegrationType: control.IntegrationType,
				PrimaryResource: control.Policy.PrimaryResource,
				ListOfResources: listOfTables,
				Language:        control.Policy.Language,
				RegoPolicies:    control.Policy.RegoPolicies,
			}
			g.controlsPolicies[control.ID] = q

			controlParameterValues := make(map[string]string)
			for _, parameter := range control.Parameters {
				controlParameterValues[parameter.Key] = parameter.Value
			}

			for _, parameter := range parameters {
				q.Parameters = append(q.Parameters, db.PolicyParameter{
					PolicyID: control.ID,
					Key:      parameter,
				})

				if v, ok := controlParameterValues[parameter]; ok {
					g.policyParamValues = append(g.policyParamValues, models.PolicyParameterValues{
						Key:       parameter,
						Value:     v,
						ControlID: control.ID,
					})
				} else {
					g.logger.Error("extract control failed: control does not contain parameter value", zap.String("control-id", control.ID),
						zap.String("parameter", parameter))
					return nil
				}
			}
			g.policies = append(g.policies, q)
			c.PolicyID = &control.ID
		}
	}
	g.controls = append(g.controls, c)
	return nil
}

func (g *GitParser) ExtractBenchmarks(complianceBenchmarksPath string) error {
	var benchmarks []Benchmark
	var frameworks []Framework
	err := filepath.WalkDir(complianceBenchmarksPath, func(path string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(filepath.Base(path), ".yaml") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			g.logger.Error("failed to read benchmark", zap.String("path", path), zap.Error(err))
			return err
		}

		if len(content) >= 9 && string(content[:10]) == "framework:" {
			var obj FrameworkFile
			err = yaml.Unmarshal(content, &obj)
			if err != nil {
				g.logger.Error("failed to unmarshal benchmark", zap.String("path", path), zap.Error(err))
				return err
			}
			if obj.Framework.ID == "" {
				g.logger.Error("failed to extract benchmark from framework", zap.String("path", path))
			}
			frameworks = append(frameworks, obj.Framework)
		} else if len(content) >= 13 && string(content[:14]) == "control-group:" {
			var obj ControlGroupFile
			err = yaml.Unmarshal(content, &obj)
			if err != nil {
				g.logger.Error("failed to unmarshal benchmark", zap.String("path", path), zap.Error(err))
				return err
			}
			if obj.ControlGroup.ID == "" {
				g.logger.Error("failed to extract benchmark from framework", zap.String("path", path))
			}
			frameworks = append(frameworks, obj.ControlGroup)
		} else {
			var obj Benchmark
			err = yaml.Unmarshal(content, &obj)
			if err != nil {
				g.logger.Error("failed to unmarshal benchmark", zap.String("path", path), zap.Error(err))
				return err
			}
			benchmarks = append(benchmarks, obj)
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = g.HandleBenchmarks(benchmarks)
	if err != nil {
		return err
	}

	err = g.HandleFrameworks(frameworks)
	if err != nil {
		return err
	}

	var newBenchmarks []db.Benchmark
	for _, b := range g.benchmarks {
		if b.ID == "" {
			g.logger.Error("benchmark id should not be empty", zap.String("id", b.ID), zap.String("title", b.Title))
			continue
		}
		newBenchmarks = append(newBenchmarks, b)
	}
	g.benchmarks = newBenchmarks

	g.benchmarks, _ = fillBenchmarksIntegrationTypes(g.benchmarks)

	return nil
}

func (g *GitParser) HandleBenchmarks(benchmarks []Benchmark) error {
	children := map[string][]string{}
	for _, o := range benchmarks {
		tags := make([]db.BenchmarkTag, 0, len(o.Tags))
		for tagKey, tagValue := range o.Tags {
			tags = append(tags, db.BenchmarkTag{
				Tag: model.Tag{
					Key:   tagKey,
					Value: tagValue,
				},
				BenchmarkID: o.ID,
			})
		}

		autoAssign := true
		if o.AutoAssign != nil {
			autoAssign = *o.AutoAssign
		}

		b := db.Benchmark{
			ID:                o.ID,
			Title:             o.Title,
			DisplayCode:       o.SectionCode,
			Description:       o.Description,
			AutoAssign:        autoAssign,
			TracksDriftEvents: o.TracksDriftEvents,
			Tags:              tags,
			Children:          nil,
			Controls:          nil,
		}
		metadataJsonb := pgtype.JSONB{}
		err := metadataJsonb.Set([]byte(""))
		if err != nil {
			return err
		}
		b.Metadata = metadataJsonb

		for _, controls := range g.controls {
			if contains(o.Controls, controls.ID) {
				b.Controls = append(b.Controls, controls)
			}
		}

		integrationTypes := make(map[string]bool)
		for _, c := range b.Controls {
			for _, it := range c.IntegrationType {
				integrationTypes[it] = true
			}
		}
		var integrationTypesList []string
		for k, _ := range integrationTypes {
			integrationTypesList = append(integrationTypesList, k)
		}
		b.IntegrationType = integrationTypesList

		g.benchmarks = append(g.benchmarks, b)
		children[o.ID] = o.Children
	}
	// g.logger.Info("Extracted benchmarks 2", zap.Int("count", len(g.benchmarks)))

	for idx, benchmark := range g.benchmarks {
		for _, childID := range children[benchmark.ID] {
			for _, child := range g.benchmarks {
				if child.ID == childID {
					benchmark.Children = append(benchmark.Children, child)
				}
			}
		}

		if len(children[benchmark.ID]) != len(benchmark.Children) {
			//fmt.Printf("could not find some benchmark children, %d != %d", len(children[benchmark.ID]), len(benchmark.Children))
		}
		g.benchmarks[idx] = benchmark
	}
	// g.logger.Info("Extracted benchmarks 3", zap.Int("count", len(g.benchmarks)))
	return nil
}

func (g *GitParser) HandleFrameworks(frameworks []Framework) error {
	for _, f := range frameworks {
		err := g.HandleSingleFramework(f)
		if err != nil {
			return err
		}
	}
	// g.logger.Info("Extracted benchmarks 2", zap.Int("count", len(g.benchmarks)))

	for idx, benchmark := range g.benchmarks {
		for _, childID := range g.frameworksChildren[benchmark.ID] {
			for _, child := range g.benchmarks {
				if child.ID == childID {
					benchmark.Children = append(benchmark.Children, child)
				}
			}
		}

		if len(g.frameworksChildren[benchmark.ID]) != len(benchmark.Children) {
			//fmt.Printf("could not find some benchmark children, %d != %d", len(children[benchmark.ID]), len(benchmark.Children))
		}
		g.benchmarks[idx] = benchmark
	}
	// g.logger.Info("Extracted benchmarks 3", zap.Int("count", len(g.benchmarks)))
	return nil
}

func (g *GitParser) HandleSingleFramework(framework Framework) error {
	var tags []db.BenchmarkTag
	if framework.Metadata != nil {
		tags = make([]db.BenchmarkTag, 0, len(framework.Metadata.Tags))
		for tagKey, tagValue := range framework.Metadata.Tags {
			tags = append(tags, db.BenchmarkTag{
				Tag: model.Tag{
					Key:   tagKey,
					Value: tagValue,
				},
				BenchmarkID: framework.ID,
			})
		}
	} else {
		tags = make([]db.BenchmarkTag, 0, len(framework.Tags))
		for tagKey, tagValue := range framework.Tags {
			tags = append(tags, db.BenchmarkTag{
				Tag: model.Tag{
					Key:   tagKey,
					Value: tagValue,
				},
				BenchmarkID: framework.ID,
			})
		}
	}

	autoAssign := true
	if framework.Metadata != nil {
		if framework.Metadata.Defaults.AutoAssign != nil {
			autoAssign = *framework.Metadata.Defaults.AutoAssign
		}
	}
	tracksDriftEvents := false
	if framework.Metadata != nil {
		tracksDriftEvents = framework.Metadata.Defaults.TracksDriftEvents
	}

	b := db.Benchmark{
		ID:                framework.ID,
		Title:             framework.Title,
		DisplayCode:       framework.SectionCode,
		Description:       framework.Description,
		AutoAssign:        autoAssign,
		TracksDriftEvents: tracksDriftEvents,
		Tags:              tags,
		Children:          nil,
		Controls:          nil,
	}
	metadataJsonb := pgtype.JSONB{}
	err := metadataJsonb.Set([]byte(""))
	if err != nil {
		return err
	}
	b.Metadata = metadataJsonb

	for _, controls := range g.controls {
		if contains(framework.Controls, controls.ID) {
			b.Controls = append(b.Controls, controls)
		}
	}

	integrationTypes := make(map[string]bool)
	for _, c := range b.Controls {
		for _, it := range c.IntegrationType {
			integrationTypes[it] = true
		}
	}
	var integrationTypesList []string
	for k, _ := range integrationTypes {
		integrationTypesList = append(integrationTypesList, k)
	}
	b.IntegrationType = integrationTypesList

	for _, group := range framework.ControlGroup {
		if len(group.Controls) > 0 || len(group.ControlGroup) > 0 {
			err = g.HandleSingleFramework(group)
			if err != nil {
				return err
			}
		}
		g.frameworksChildren[framework.ID] = append(g.frameworksChildren[framework.ID], group.ID)
	}
	g.benchmarks = append(g.benchmarks, b)
	return nil
}

func fillBenchmarksIntegrationTypes(benchmarks []db.Benchmark) ([]db.Benchmark, []string) {
	var integrationTypes []string
	integrationTypesMap := make(map[string]bool)

	for idx, benchmark := range benchmarks {
		if benchmark.IntegrationType == nil {
			benchmark.Children, benchmark.IntegrationType = fillBenchmarksIntegrationTypes(benchmark.Children)
			benchmarks[idx] = benchmark
		}
		for _, c := range benchmark.IntegrationType {
			if _, ok := integrationTypesMap[c]; !ok {
				integrationTypes = append(integrationTypes, c)
				integrationTypesMap[c] = true
			}
		}
	}
	return benchmarks, integrationTypes
}

func (g *GitParser) CheckForDuplicate() error {
	visited := map[string]bool{}
	for _, b := range g.benchmarks {
		if _, ok := visited[b.ID]; !ok {
			visited[b.ID] = true
		} else {
			return fmt.Errorf("duplicate benchmark id: %s", b.ID)
		}
	}
	return nil
}

func (g GitParser) ExtractBenchmarksMetadata() error {
	for i, b := range g.benchmarks {
		benchmarkControlsCache := make(map[string]BenchmarkControlsCache)
		controlsMap, err := getControlsUnderBenchmark(b, benchmarkControlsCache)
		if err != nil {
			return err
		}
		benchmarkTablesCache := make(map[string]BenchmarkTablesCache)
		primaryTablesMap, listOfTablesMap, err := g.getTablesUnderBenchmark(b, benchmarkTablesCache)
		if err != nil {
			return err
		}
		var listOfTables, primaryTables, controls []string
		for k, _ := range controlsMap {
			controls = append(controls, k)
		}
		for k, _ := range primaryTablesMap {
			primaryTables = append(primaryTables, k)
		}
		for k, _ := range listOfTablesMap {
			listOfTables = append(listOfTables, k)
		}
		metadata := db.BenchmarkMetadata{
			Controls:         controls,
			PrimaryResources: primaryTables,
			ListOfResources:  listOfTables,
		}
		metadataJson, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		metadataJsonb := pgtype.JSONB{}
		err = metadataJsonb.Set(metadataJson)
		if err != nil {
			return err
		}
		g.benchmarks[i].Metadata = metadataJsonb
	}
	return nil
}

func (g *GitParser) ExtractCompliance(compliancePath string, controlEnrichmentBasePath string) error {
	if err := g.ExtractNamedQueries(); err != nil {
		return err
	}
	if err := g.ExtractPolicies(path.Join(compliancePath, "policies")); err != nil {
		return err
	}
	if err := g.ExtractControls(path.Join(compliancePath, "controls"), controlEnrichmentBasePath); err != nil {
		return err
	}
	if err := g.ExtractBenchmarks(path.Join(compliancePath, "frameworks")); err != nil {
		return err
	}
	//if err := g.CheckForDuplicate(); err != nil {
	//	return err
	//}

	if err := g.ExtractBenchmarksMetadata(); err != nil {
		return err
	}
	return nil
}

func (g *GitParser) ExtractQueryViews(viewsPath string) error {
	return filepath.WalkDir(viewsPath, func(path string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			g.logger.Error("failed to read query view", zap.String("path", path), zap.Error(err))
			return err
		}

		var obj QueryView
		err = yaml.Unmarshal(content, &obj)
		if err != nil {
			g.logger.Error("failed to unmarshal query view", zap.String("path", path), zap.Error(err))
			return nil
		}

		qv := models.QueryView{
			ID:          obj.ID,
			Title:       obj.Title,
			Description: obj.Description,
		}

		listOfTables, err := utils.ExtractTableRefsFromPolicy(types.PolicyLanguageSQL, obj.Query)
		if err != nil {
			g.logger.Error("failed to extract table refs from query", zap.String("query-id", obj.ID), zap.Error(err))
		}

		q := models.Query{
			ID:             obj.ID,
			QueryToExecute: obj.Query,
			ListOfTables:   listOfTables,
			Engine:         "sql",
		}

		g.coreServiceQueries = append(g.coreServiceQueries, q)
		qv.QueryID = &obj.ID

		g.queryViews = append(g.queryViews, qv)

		return nil
	})
}

func contains[T uint | int | string](arr []T, ob T) bool {
	for _, o := range arr {
		if o == ob {
			return true
		}
	}
	return false
}
