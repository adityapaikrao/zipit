include .env
export

DB_URL_LOCAL := $(if $(DATABASE_URL),$(DATABASE_URL),postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE))

build_proto:
	mkdir -p gen
	protoc --proto_path=./proto \
	--go_out=gen --go_opt=paths=source_relative \
	--go-grpc_out=gen --go-grpc_opt=paths=source_relative \
	proto/url/url.proto proto/analytics/analytics.proto

migrate_up:
	migrate -path ./migrations -database "$(DB_URL_LOCAL)" up

migrate_down:
	migrate -path ./migrations -database "$(DB_URL_LOCAL)" down

migrate_create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name
