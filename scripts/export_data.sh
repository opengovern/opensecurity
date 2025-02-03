mkdir -p /tmp/demo-data

echo "test1" > test.txt
aws s3 cp ./test.txt "s3://test.txt" --endpoint-url="$ENDPOINT_URL" --region "$BUCKET_REGION"

mkdir -p /tmp/demo-data/es-demo
NEW_ELASTICSEARCH_ADDRESS="https://${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}@${ELASTICSEARCH_ADDRESS#https://}"

echo $NEW_ELASTICSEARCH_ADDRESS 
export NODE_TLS_REJECT_UNAUTHORIZED=0
multielasticdump \
  --direction=dump \
  --input="$NEW_ELASTICSEARCH_ADDRESS" \
  --output="/tmp/demo-data/es-demo/" \
  --parallel=20 \
  --match='^(?!\.)(?!.*(security|deleted)).*$' \
  --matchType=alias \
  --limit=10000 \
  --scrollTime=10m \
  --searchBody='{"query": {"bool": {"must_not": {"term": {"deleted": true}}}}}' \
  --ignoreTemplate=true

mkdir -p /tmp/demo-data/postgres
pg_dump --dbname="postgresql://$POSTGRESQL_USERNAME:$POSTGRESQL_PASSWORD@$POSTGRESQL_HOST:$POSTGRESQL_PORT/describe" > /tmp/demo-data/postgres/describe.sql
pg_dump --dbname="postgresql://$POSTGRESQL_USERNAME:$POSTGRESQL_PASSWORD@$POSTGRESQL_HOST:$POSTGRESQL_PORT/compliance" > /tmp/demo-data/postgres/compliance.sql
pg_dump --dbname="postgresql://$POSTGRESQL_USERNAME:$POSTGRESQL_PASSWORD@$POSTGRESQL_HOST:$POSTGRESQL_PORT/integration" > /tmp/demo-data/postgres/integration.sql
pg_dump --dbname="postgresql://$POSTGRESQL_USERNAME:$POSTGRESQL_PASSWORD@$POSTGRESQL_HOST:$POSTGRESQL_PORT/core" --exclude-table=platform_configurations > /tmp/demo-data/postgres/core.sql

cd /tmp
tar -cO demo-data | openssl enc -aes-256-cbc -md md5 -pass pass:"$OPENSSL_PASSWORD" -base64 > demo_data.tar.gz.enc

FILE_SIZE_BYTES=$(stat -c %s /tmp/demo_data.tar.gz.enc)
FILE_SIZE_MB=$(echo "scale=2; $FILE_SIZE_BYTES / 1048576" | bc)
echo "File size: ${FILE_SIZE_MB} MB"

aws s3 cp /tmp/demo_data.tar.gz.enc "$DEMO_DATA_S3_PATH" --endpoint-url="$ENDPOINT_URL" --region "$BUCKET_REGION"
