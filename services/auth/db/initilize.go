package db

import (
	// Needed for ColumnType
	"errors"
	"fmt"
	"reflect" // Needed for schema reflection
	"strings" // Needed for joining error messages

	// Potentially needed if Initialize does more time-based setup
	"go.uber.org/zap" // For logging
	"gorm.io/gorm"    // GORM core
	// GORM schema utilities
)

// Initialize runs AutoMigrate and then verifies the schema strictly.
// This method is associated with the Database struct defined in db.go
func (db Database) Initialize() error {
	// Define models in one place for consistency
	modelsToMigrate := []interface{}{
		&ApiKey{},
		&User{},
		&Configuration{},
		&Connector{},
	}

	db.Logger.Info("Running AutoMigrate...")
	err := db.Orm.AutoMigrate(modelsToMigrate...) // Pass the slice of models
	if err != nil {
		db.Logger.Error("AutoMigrate failed", zap.Error(err))
		return fmt.Errorf("automigrate failed: %w", err) // Wrap error for context
	}
	db.Logger.Info("AutoMigrate finished successfully.")

	// --- Schema Verification Step ---
	db.Logger.Info("Verifying database schema alignment (strict existence check)...")
	err = db.verifySchemaStrictExistence(modelsToMigrate...) // Call the verification method
	if err != nil {
		db.Logger.Error("Database schema verification failed", zap.Error(err))
		// Return the verification error to stop service startup if schema doesn't match
		return fmt.Errorf("schema verification failed: %w", err)
	}
	db.Logger.Info("Database schema verification successful.")
	// --- End Schema Verification Step ---

	return nil
}

// verifySchemaStrictExistence checks:
// 1. Tables defined in models exist.
// 2. Columns defined in models exist in the corresponding tables.
// 3. No extra columns exist in the tables that are not defined in the models.
// It focuses purely on existence, making it DB-type agnostic.
// This method is associated with the Database struct defined in db.go
func (db Database) verifySchemaStrictExistence(models ...interface{}) error {
	migrator := db.Orm.Migrator()
	var allErrors []string

	for _, model := range models {
		// build a new Statement; Parse will use db.Orm.Config.NamingStrategy under the hood
		stmt := &gorm.Statement{
			DB: db.Orm,
		}

		// ensure model is a non-nil pointer
		modelValue := reflect.ValueOf(model)
		if modelValue.Kind() != reflect.Ptr {
			ptr := reflect.New(modelValue.Type())
			ptr.Elem().Set(modelValue)
			model = ptr.Interface()
		} else if modelValue.IsNil() {
			db.Logger.Debug("verifySchema encountered nil model pointer, creating instance",
				zap.String("type", modelValue.Type().Elem().Name()))
			model = reflect.New(modelValue.Type().Elem()).Interface()
		}

		// parse the schema
		if err := stmt.Parse(model); err != nil {
			detail := fmt.Sprintf("failed to parse schema for %T: %v", model, err)
			db.Logger.Warn("Schema verification issue: parsing failed", zap.String("detail", detail))
			allErrors = append(allErrors, detail)
			continue
		}

		tableName := stmt.Schema.Table
		modelName := stmt.Schema.Name
		db.Logger.Debug("Verifying schema", zap.String("model", modelName), zap.String("table", tableName))

		// 1. table exists?
		if !migrator.HasTable(tableName) {
			detail := fmt.Sprintf("table '%s' for model '%s' does not exist", tableName, modelName)
			db.Logger.Warn("Schema verification issue: table missing", zap.String("detail", detail))
			allErrors = append(allErrors, detail)
			continue
		}

		// 2. model→column exists?
		expectedCols := map[string]bool{}
		for _, f := range stmt.Schema.Fields {
			// Only consider fields that map to a database column
			if f.DBName == "" {
				continue
			}
			expectedCols[f.DBName] = true // Record expected column name (resolved by Parse)
			db.Logger.Debug("  - Checking model column existence in DB", zap.String("column", f.DBName), zap.String("field", f.Name), zap.String("model", modelName))
			if !migrator.HasColumn(model, f.DBName) { // Check using the resolved DBName
				detail := fmt.Sprintf("column '%s' (field '%s') defined in model '%s' missing from table '%s'",
					f.DBName, f.Name, modelName, tableName)
				db.Logger.Warn("Schema verification issue: model column missing in DB", zap.String("detail", detail))
				allErrors = append(allErrors, detail)
			}
		}

		// 3. no extra columns?
		actualCols, err := migrator.ColumnTypes(model)
		if err != nil {
			detail := fmt.Sprintf("could not retrieve actual columns for table '%s': %v", tableName, err)
			db.Logger.Warn("Schema verification issue: failed to get actual columns", zap.String("detail", detail), zap.String("model", modelName))
			allErrors = append(allErrors, detail)
		} else {
			db.Logger.Debug("  - Checking for extra columns in DB table", zap.String("model", modelName), zap.String("table", tableName))
			for _, col := range actualCols {
				name := col.Name()
				// Check if the column found in the DB was expected based on the model
				if !expectedCols[name] {
					detail := fmt.Sprintf("table '%s' has an EXTRA column '%s' which is not defined in model '%s'",
						tableName, name, modelName)
					db.Logger.Warn("Schema verification issue: extra column found in DB", zap.String("detail", detail))
					allErrors = append(allErrors, detail)
				}
			}
		}
	} // End loop through models

	if len(allErrors) > 0 {
		return errors.New("schema verification failed:\n - " + strings.Join(allErrors, "\n - "))
	}

	db.Logger.Info("Strict schema existence verification passed for all checked models.")
	return nil
}
