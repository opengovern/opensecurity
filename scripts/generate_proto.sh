protoc --go_out=services/describe/proto/src/golang --go_opt=paths=source_relative \
    --go-grpc_out=services/describe/proto/src/golang --go-grpc_opt=paths=source_relative \
    services/describe/proto/*.proto
mv services/describe/proto/src/golang/services/describe/proto/* services/describe/proto/src/golang/
rm -rf services/describe/proto/src/golang/pkg