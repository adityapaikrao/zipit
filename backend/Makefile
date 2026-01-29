build_proto:
	mkdir -p gen
	protoc --proto_path=./proto \
	--go_out=gen --go_opt=paths=source_relative \
	--go-grpc_out=gen --go-grpc_opt=paths=source_relative \
	proto/url/url.proto proto/analytics/analytics.proto