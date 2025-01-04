package utils

import (
	"fmt"
	"github.com/opengovern/opencomply/pkg/types"
	pg_query "github.com/pganalyze/pg_query_go/v4"
	"regexp"
)

// extractSQLTableRefs walks the pg_query_go AST to find all referenced tables.
func extractSQLTableRefs(node *pg_query.Node) []string {
	if node == nil {
		return nil
	}
	var tables []string
	switch n := node.Node.(type) {

	// SELECT
	case *pg_query.Node_SelectStmt:
		sel := n.SelectStmt
		// FROM
		for _, fromItem := range sel.FromClause {
			tables = append(tables, extractSQLTableRefs(fromItem)...)
		}
		// WHERE
		if sel.WhereClause != nil {
			tables = append(tables, extractSQLTableRefs(sel.WhereClause)...)
		}
		// GROUP BY
		for _, g := range sel.GroupClause {
			tables = append(tables, extractSQLTableRefs(g)...)
		}
		// HAVING
		if sel.HavingClause != nil {
			tables = append(tables, extractSQLTableRefs(sel.HavingClause)...)
		}
		// TARGET LIST
		for _, t := range sel.TargetList {
			tables = append(tables, extractSQLTableRefs(t)...)
		}
		// UNION, etc.
		if sel.Larg != nil {
			leftNode := &pg_query.Node{Node: &pg_query.Node_SelectStmt{SelectStmt: sel.Larg}}
			tables = append(tables, extractSQLTableRefs(leftNode)...)
		}
		if sel.Rarg != nil {
			rightNode := &pg_query.Node{Node: &pg_query.Node_SelectStmt{SelectStmt: sel.Rarg}}
			tables = append(tables, extractSQLTableRefs(rightNode)...)
		}

	// JOIN
	case *pg_query.Node_JoinExpr:
		j := n.JoinExpr
		if j.Larg != nil {
			tables = append(tables, extractSQLTableRefs(j.Larg)...)
		}
		if j.Rarg != nil {
			tables = append(tables, extractSQLTableRefs(j.Rarg)...)
		}
		if j.Quals != nil {
			tables = append(tables, extractSQLTableRefs(j.Quals)...)
		}

	// RANGE SUBSELECT
	case *pg_query.Node_RangeSubselect:
		rs := n.RangeSubselect
		if rs.Subquery != nil {
			tables = append(tables, extractSQLTableRefs(rs.Subquery)...)
		}

	// SUBLINK
	case *pg_query.Node_SubLink:
		s := n.SubLink
		if s.Subselect != nil {
			tables = append(tables, extractSQLTableRefs(s.Subselect)...)
		}

	// TABLE REFERENCE
	case *pg_query.Node_RangeVar:
		rv := n.RangeVar
		if rv.Relname != "" {
			tables = append(tables, rv.Relname)
		}

	// INSERT
	case *pg_query.Node_InsertStmt:
		ins := n.InsertStmt
		if ins.Relation != nil && ins.Relation.Relname != "" {
			tables = append(tables, ins.Relation.Relname)
		}
		if ins.SelectStmt != nil {
			tables = append(tables, extractSQLTableRefs(ins.SelectStmt)...)
		}

	// UPDATE
	case *pg_query.Node_UpdateStmt:
		upd := n.UpdateStmt
		if upd.Relation != nil && upd.Relation.Relname != "" {
			tables = append(tables, upd.Relation.Relname)
		}
		for _, f := range upd.FromClause {
			tables = append(tables, extractSQLTableRefs(f)...)
		}
		if upd.WhereClause != nil {
			tables = append(tables, extractSQLTableRefs(upd.WhereClause)...)
		}

	// DELETE
	case *pg_query.Node_DeleteStmt:
		del := n.DeleteStmt
		if del.Relation != nil && del.Relation.Relname != "" {
			tables = append(tables, del.Relation.Relname)
		}
		for _, u := range del.UsingClause {
			tables = append(tables, extractSQLTableRefs(u)...)
		}
		if del.WhereClause != nil {
			tables = append(tables, extractSQLTableRefs(del.WhereClause)...)
		}
	}

	return tables
}

func extractSQLQueryParameters(query string) ([]string, error) {
	// Define the regex pattern to match parameters in the format {{.parameterName}}
	pattern := `{{\.\s*([a-zA-Z0-9_]+)\s*}}`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}

	// Find all matches
	matches := re.FindAllStringSubmatch(query, -1)

	// Use a map to ensure uniqueness
	paramMap := make(map[string]struct{})
	for _, match := range matches {
		if len(match) > 1 {
			paramMap[match[1]] = struct{}{}
		}
	}

	// Convert the map keys to a slice
	var parameters []string
	for param := range paramMap {
		parameters = append(parameters, param)
	}

	return parameters, nil
}

func extractRegoParameters(modules []string) []string {
	paramRefRegex := regexp.MustCompile(`input\.params\.([A-Za-z0-9_]+)`)
	parameters := make(map[string]bool)
	for _, mod := range modules {
		matches := paramRefRegex.FindAllStringSubmatch(mod, -1)
		for _, m := range matches {
			parameters[m[1]] = true
		}
	}

	var paramsList []string
	for param := range parameters {
		paramsList = append(paramsList, param)
	}

	return paramsList
}

func ExtractParameters(language types.PolicyLanguage, definition string) ([]string, error) {
	var parameters []string
	var err error

	switch language {
	case types.PolicyLanguageSQL:
		parameters, err = extractSQLQueryParameters(definition)
		if err != nil {
			return nil, err
		}
	}

	return parameters, nil
}

func ExtractTableRefsFromPolicy(language types.PolicyLanguage, definition string) ([]string, error) {
	var tables []string

	switch language {
	case types.PolicyLanguageSQL:
		// Parse the SQL using pg_query_go
		parseResult, err := pg_query.Parse(definition)
		if err != nil {
			return nil, err
		}

		// Collect table references from each statement
		for _, rawStmt := range parseResult.Stmts {
			stmtTables := extractSQLTableRefs(rawStmt.Stmt)
			tables = append(tables, stmtTables...)
		}
	default:
		return nil, fmt.Errorf("unsupported policy language: %s", language)
	}

	return tables, nil
}
