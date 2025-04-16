echo "NEEDED ENV Variables are: [PLUGIN_ID, ES_INDEX_PREFIX, DB_PASSWORD, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ENDPOINT_URL, ELASTICSEARCH_PASSWORD, OPENSSL_PASSWORD, DEMO_DATA_S3_PATH]"

mkdir -p /tmp/demo-data


DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD}" # Should be set in the environment
DB_NAME="${DB_NAME:-integration}"
OUTPUT_FILE="/tmp/demo-data/integrations.json"

if [ -z "$DB_PASSWORD" ]; then
    echo "Error: DB_PASSWORD environment variable is not set." >&2
    exit 1
fi

if [ -z "$PLUGIN_ID" ]; then
    echo "Error: PLUGIN_ID environment variable is not set." >&2
    exit 1
fi

if [ -z "$ES_INDEX_PREFIX" ]; then
    echo "Error: ES_INDEX_PREFIX environment variable is not set." >&2
    exit 1
fi

export PGPASSWORD="$DB_PASSWORD"

SQL_QUERY=$(cat <<-EOF
SELECT
    json_build_object(
        'integrationId', integration_id::text, -- Cast UUID to text if needed
        'providerId', provider_id,
        'name', name,
        'integrationType', integration_type, -- Assuming it's stored as text compatible with your Go type
        'annotations', annotations,       -- Assuming annotations is a JSONB column
        'labels', labels              -- Assuming labels is a JSONB column
    )
FROM
    integrations
WHERE
    integration_type = '${PLUGIN_ID}';
EOF
)

echo "Connecting to database '$DB_NAME' on '$DB_HOST:$DB_PORT' as user '$DB_USER'..."
echo "Executing query and exporting to '$OUTPUT_FILE'..."

psql -X -A -t -q -v ON_ERROR_STOP=1 \
    -h "$DB_HOST" \
    -p "$DB_PORT" \
    -U "$DB_USER" \
    -d "$DB_NAME" \
    -c "$SQL_QUERY" | \
jq -s '.' > "$OUTPUT_FILE"

if [ -s "$OUTPUT_FILE" ]; then
    echo "Successfully exported integrations to '$OUTPUT_FILE'."
else
    echo "Error: Failed to export data or the table is empty. Output file is empty or not created." >&2
    exit 1
fi

unset PGPASSWORD

echo "Starting dumping elasticsearch data"

mkdir -p /tmp/demo-data/es-demo
ELASTICSEARCH_USERNAME="${ELASTICSEARCH_USERNAME:-admin}"
ELASTICSEARCH_ADDRESS="${ELASTICSEARCH_ADDRESS:-https://localhost:9200}"
if [ -z "$ELASTICSEARCH_PASSWORD" ]; then
    echo "Error: ELASTICSEARCH_PASSWORD environment variable is not set." >&2
    exit 1
fi
NEW_ELASTICSEARCH_ADDRESS="https://${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD}@${ELASTICSEARCH_ADDRESS#https://}"
mkdir -p ~/.config/rclone

cat <<EOF > ~/.config/rclone/rclone.conf
[r2]
type = s3
provider = Cloudflare
access_key_id = $AWS_ACCESS_KEY_ID
secret_access_key = $AWS_SECRET_ACCESS_KEY
region = auto
endpoint = $ENDPOINT_URL
acl = private
EOF

echo $NEW_ELASTICSEARCH_ADDRESS
export NODE_TLS_REJECT_UNAUTHORIZED=0
multielasticdump \
  --direction=dump \
  --input="$NEW_ELASTICSEARCH_ADDRESS" \
  --output="/tmp/demo-data/es-demo/" \
  --parallel=20 \
  --match='^'${ES_INDEX_PREFIX}'.*$' \
  --matchType=alias \
  --limit=10000 \
  --scrollTime=10m \
  --searchBody='{"query": {"bool": {"must_not": {"term": {"deleted": true}}}}}' \
  --ignoreTemplate=true

cd /tmp
tar -cO demo-data | openssl enc -aes-256-cbc -md md5 -pass pass:"$OPENSSL_PASSWORD" -base64 > demo_data.tar.gz.enc

FILE_SIZE_BYTES=$(stat -c %s /tmp/demo_data.tar.gz.enc)
FILE_SIZE_MB=$(echo "scale=2; $FILE_SIZE_BYTES / 1048576" | bc)
echo "File size: ${FILE_SIZE_MB} MB"

rclone copy /tmp/demo_data.tar.gz.enc "$DEMO_DATA_S3_PATH"
