package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	steampipesdk "github.com/opengovern/og-util/pkg/steampipe"
	"go.uber.org/zap"
	"strings"
	"time"
)

type RegoEngine struct {
	logger    *zap.Logger
	steampipe *steampipesdk.Database

	regoFunctions []func(*rego.Rego)
}

var excludedTableSchema = []string{"information_schema", "pg_catalog", "steampipe_internal", "steampipe_command", "public"}

func NewRegoEngine(ctx context.Context, logger *zap.Logger, steampipeDb *steampipesdk.Database) (*RegoEngine, error) {
	engine := RegoEngine{
		logger:    logger,
		steampipe: steampipeDb,
	}

	tries := 5
	for i := 0; i < tries; i++ {
		logger.Info("Getting rego functions", zap.Int("try", i+1))
		functions, err := engine.getRegoFunctionForTables(ctx)
		if err != nil {
			logger.Error("Error getting rego functions", zap.Error(err))
			if i == tries-1 {
				logger.Error("Exhausted all tries to get rego functions")
				return nil, err
			}
			time.Sleep(10 * time.Second)
			continue
		}
		engine.regoFunctions = functions
		break
	}
	logger.Info("Got rego functions")

	return &engine, nil
}

func (r *RegoEngine) getRegoFunctionForTables(ctx context.Context) ([]func(*rego.Rego), error) {
	r.logger.Info("getting rego functions for tables")
	rows, err := r.steampipe.Conn().Query(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema != any ($1)", excludedTableSchema)
	if err != nil {
		r.logger.Error("Unable to query database", zap.Error(err))
		r.logger.Sync()
		return nil, err
	}
	defer rows.Close()
	r.logger.Info("queried database for tables", zap.Int("rows", len(rows.RawValues())))

	results := make([]func(*rego.Rego), 0)
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			r.logger.Error("Unable to scan table name", zap.Error(err))
			r.logger.Sync()
			return nil, err
		}
		r.logger.Info("scanned table name", zap.String("table", tableName))

		f := rego.FunctionDyn(&rego.Function{
			Name:             fmt.Sprintf("opencomply.%s", tableName),
			Description:      "",
			Decl:             types.NewFunction([]types.Type{types.NewObject(nil, &types.DynamicProperty{Key: types.S, Value: types.NewArray(nil, types.A)})}, types.NewArray(nil, types.A)),
			Memoize:          true,
			Nondeterministic: true,
		}, func(bctx rego.BuiltinContext, terms []*ast.Term) (*ast.Term, error) {
			var args []any
			var whereClause string
			if len(terms) > 0 {
				whereObject, err := ast.ValueToInterface(terms[0].Value, nil)
				if err != nil {
					r.logger.Error("Unable to convert to interface", zap.Error(err))
					r.logger.Sync()
					return nil, err
				}
				if whereMap, ok := whereObject.(map[string]any); ok && len(whereMap) > 0 {
					whereClauseBuilder := strings.Builder{}
					whereClauseBuilder.WriteString("WHERE ")
					counter := 1
					for key, value := range whereMap {
						if list, ok := value.([]interface{}); ok {
							whereClauseBuilder.WriteString(fmt.Sprintf("%s IN (", key))
							for i, v := range list {
								whereClauseBuilder.WriteString(fmt.Sprintf("$%d", counter))
								if i < len(list)-1 {
									whereClauseBuilder.WriteString(", ")
								}
								args = append(args, v)
								counter++
							}
							whereClauseBuilder.WriteString(") AND ")
						} else {
							return nil, fmt.Errorf("invalid where clause: %v", whereObject)
						}
					}
					whereClause = strings.TrimSuffix(whereClauseBuilder.String(), "AND ")
				} else {
					return nil, fmt.Errorf("invalid where clause: %v", whereObject)
				}
			}

			rows, err := r.steampipe.Conn().Query(bctx.Context, fmt.Sprintf("SELECT * FROM %s %s", tableName, whereClause), args...)
			if err != nil {
				r.logger.Error("Unable to query database", zap.Error(err), zap.String("table", tableName))
				r.logger.Sync()
				return nil, err
			}
			defer rows.Close()

			results, err := pgx.CollectRows(rows, pgx.RowToMap)
			if err != nil {
				r.logger.Error("Unable to scan row", zap.Error(err), zap.String("table", tableName))
				r.logger.Sync()
				return nil, err
			}

			value, err := ast.InterfaceToValue(results)
			if err != nil {
				r.logger.Error("Unable to convert to value", zap.Error(err), zap.String("table", tableName))
				r.logger.Sync()
				return nil, err
			}

			return ast.NewTerm(value), nil
		})
		results = append(results, f)
		r.logger.Info("created rego function", zap.String("table", tableName))
	}

	// Add raw query function
	f := rego.FunctionDyn(&rego.Function{
		Name:             "opencomply.cloudql",
		Description:      "",
		Decl:             types.NewFunction([]types.Type{types.S}, types.NewArray(nil, types.A)),
		Memoize:          true,
		Nondeterministic: true,
	}, func(bctx rego.BuiltinContext, terms []*ast.Term) (*ast.Term, error) {
		var query string
		if len(terms) == 0 {
			return nil, fmt.Errorf("invalid query")
		}
		query = terms[0].Value.String()

		rows, err := r.steampipe.Conn().Query(bctx.Context, query)
		if err != nil {
			r.logger.Error("Unable to query database", zap.Error(err))
			r.logger.Sync()
			return nil, err
		}
		defer rows.Close()

		results, err := pgx.CollectRows(rows, pgx.RowToMap)
		if err != nil {
			r.logger.Error("Unable to scan row", zap.Error(err))
			r.logger.Sync()
			return nil, err
		}

		value, err := ast.InterfaceToValue(results)
		if err != nil {
			r.logger.Error("Unable to convert to value", zap.Error(err))
			r.logger.Sync()
			return nil, err
		}

		return ast.NewTerm(value), nil
	})

	r.logger.Info("created rego function", zap.String("function", "opencomply.cloudql"))
	results = append(results, f)

	return results, nil
}

func (r *RegoEngine) Evaluate(ctx context.Context, policies []string, query string) (rego.ResultSet, error) {
	params := append(r.regoFunctions, rego.Query(query))
	for i, policy := range policies {
		params = append(params, rego.Module(fmt.Sprintf("policy_%d.rego", i+1), policy))
	}

	regoEngine := rego.New(params...)
	results, err := regoEngine.Eval(ctx)
	if err != nil {
		r.logger.Error("Error evaluating policy", zap.Error(err))
		r.logger.Sync()
		return nil, err
	}

	return results, nil
}
