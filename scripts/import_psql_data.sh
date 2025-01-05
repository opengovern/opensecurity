# https://github.com/elasticsearch-dump/elasticsearch-dump

curl -O "$DEMO_DATA_S3_URL"

openssl enc -d -aes-256-cbc -md md5 -pass pass:"$OPENSSL_PASSWORD" -base64 -in demo_data.tar.gz.enc -out demo_data.tar.gz
tar -xvf demo_data.tar.gz

echo "$POSTGRESQL_HOST"
echo "$POSTGRESQL_PORT"
echo "$POSTGRESQL_USERNAME"
echo "$POSTGRESQL_PASSWORD"

# drop and recreate the databases
PGPASSWORD="$POSTGRESQL_PASSWORD" psql --host="$POSTGRESQL_HOST" --port="$POSTGRESQL_PORT" --username "$POSTGRESQL_USERNAME" -c "DROP DATABASE describe;"
PGPASSWORD="$POSTGRESQL_PASSWORD" psql --host="$POSTGRESQL_HOST" --port="$POSTGRESQL_PORT" --username "$POSTGRESQL_USERNAME" -c "CREATE DATABASE describe;"
PGPASSWORD="$POSTGRESQL_PASSWORD" psql --host="$POSTGRESQL_HOST" --port="$POSTGRESQL_PORT" --username "$POSTGRESQL_USERNAME" -c "DROP DATABASE integration;"
PGPASSWORD="$POSTGRESQL_PASSWORD" psql --host="$POSTGRESQL_HOST" --port="$POSTGRESQL_PORT" --username "$POSTGRESQL_USERNAME" -c "CREATE DATABASE integration;"
PGPASSWORD="$POSTGRESQL_PASSWORD" psql --host="$POSTGRESQL_HOST" --port="$POSTGRESQL_PORT" --username "$POSTGRESQL_USERNAME" -c "DROP DATABASE core;"
PGPASSWORD="$POSTGRESQL_PASSWORD" psql --host="$POSTGRESQL_HOST" --port="$POSTGRESQL_PORT" --username "$POSTGRESQL_USERNAME" -c "CREATE DATABASE core;"

PGPASSWORD="$POSTGRESQL_PASSWORD" psql --host="$POSTGRESQL_HOST" --port="$POSTGRESQL_PORT" --username "$POSTGRESQL_USERNAME" --dbname "describe" < /demo-data/postgres/describe.sql
PGPASSWORD="$POSTGRESQL_PASSWORD" psql --host="$POSTGRESQL_HOST" --port="$POSTGRESQL_PORT" --username "$POSTGRESQL_USERNAME" --dbname "integration" < /demo-data/postgres/integration.sql
PGPASSWORD="$POSTGRESQL_PASSWORD" psql --host="$POSTGRESQL_HOST" --port="$POSTGRESQL_PORT" --username "$POSTGRESQL_USERNAME" --dbname "core" < /demo-data/postgres/core.sql
PGPASSWORD="$POSTGRESQL_PASSWORD" psql --host="$POSTGRESQL_HOST" --port="$POSTGRESQL_PORT" --username "$POSTGRESQL_USERNAME" --dbname "integration" -c "DELETE FROM credentials;"

rm -rf /demo-data/postgres