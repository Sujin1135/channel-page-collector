# Channel Page Collector

## Run

**Generate pb.go files from protobuf**
```bash
protoc \
  --proto_path=api/grpc/interface/channel-page-collector-interface/protobuf \
  --go_out=paths=source_relative:api/grpc/interface/channel-page-collector-interface/protobuf \
  --go-grpc_out=paths=source_relative:api/grpc/interface/channel-page-collector-interface/protobuf \
  api/grpc/interface/channel-page-collector-interface/protobuf/*.proto
```