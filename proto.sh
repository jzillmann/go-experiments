# Server
rm -rf ./server/proto
mkdir ./server/proto

protoc --go_out=./server/proto --go_opt=paths=source_relative \
    --go-grpc_out=./server/proto --go-grpc_opt=paths=source_relative \
    --proto_path=./proto \
    ./proto/*.proto