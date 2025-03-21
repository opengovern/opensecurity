package create_functions

import (
	steampipesdk "github.com/opengovern/og-util/pkg/steampipe"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func CreatePostgresFunctions(ctx context.Context, logger *zap.Logger) {
	selfClient, err := steampipesdk.NewSelfClient(ctx)
	if err != nil {
		logger.Error("Error creating self client for refreshing materialized views", zap.Error(err))
		logger.Sync()
		return
	}

	query := `CREATE OR REPLACE FUNCTION get_columns_for_query(_sql text)
  RETURNS TABLE (ordinal_position int, column_name text, data_type text)
  LANGUAGE plpgsql AS
$$
DECLARE
  temp_table_name text := 'tmpcols_' || md5(random()::text);
BEGIN
  -- 1) Create a temp table from the user query
  EXECUTE format('CREATE TEMP TABLE %I AS (%s) LIMIT 0', temp_table_name, _sql);

  -- 2) Return columns as rows
  RETURN QUERY EXECUTE format(
    'SELECT ordinal_position, column_name, data_type
     FROM information_schema.columns
     WHERE table_name = %L
       AND table_schema = current_schema()
     ORDER BY ordinal_position', temp_table_name
  );

  -- 3) Cleanup
  EXECUTE format('DROP TABLE %I', temp_table_name);
END;
$$;`
	_, err = selfClient.GetConnection().Exec(ctx, query)
	if err != nil {
		logger.Error("Error creating get_columns_for_query function", zap.Error(err))
		logger.Sync()
		return
	}

	return
}
