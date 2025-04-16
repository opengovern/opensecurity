mkdir -p /tmp/demo-data

mkdir -p /tmp/demo-data/es-demo
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
  --match='^(?!\.)(?!.*(security|deleted|logs|metrics)).*$' \
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
