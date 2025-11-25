#GOOS=linux
#GOARCH=amd64

run:
	go run ./cmd/server/main.go

run-auth:
	CONFIG_PATH=config/auth-config.yaml go run ./cmd/auth/main.go

run-core:
	CONFIG_PATH=config/core-config.yaml go run ./cmd/core/main.go

run-chat:
	CONFIG_PATH=config/chat-config.yaml go run ./cmd/chat/main.go

test:
	go test ./... -v

test-coverage:
	go test ./... -cover

# Test только Converters
test-converters:
	@echo "Running converter tests..."
	@go test -v ./internal/core-service/converters/...
	@go test -v ./internal/chat-service/converters/...
	@go test -v ./internal/gateway/adapters/... -run TestCoreProto

# Test с покрытием для Converters
test-converters-coverage:
	@echo "Running converter tests with coverage..."
	@go test -cover ./internal/core-service/converters/...
	@go test -cover ./internal/chat-service/converters/...
	@go test -cover ./internal/gateway/adapters/... -run TestCoreProto
	@echo "\nDetailed coverage:"
	@go test -coverprofile=coverage-converters.out ./internal/core-service/converters/... ./internal/chat-service/converters/... ./internal/gateway/adapters/...
	@go tool cover -func=coverage-converters.out | grep total

build-docs:
	swag init -g /cmd/auth/main.go -o api/auth/
	swag init -g /cmd/core/main.go -o api/core/
	swag init -g /cmd/chat/main.go -o api/chat/
	swag init -g /cmd/server/main.go -o api/server/

build:
	go build -o ./.build/auth ./cmd/auth/main.go
	go build -o ./.build/core ./cmd/core/main.go

build-bin:
	GOOS=linux GOARCH=amd64 go build -o ./.build/auth ./cmd/auth/main.go
	GOOS=linux GOARCH=amd64 go build -o ./.build/core ./cmd/core/main.go
	GOOS=linux GOARCH=amd64 go build -o ./.build/chat ./cmd/chat/main.go
	GOOS=linux GOARCH=amd64 go build -o ./.build/server ./cmd/server/main.go

clean:
	rm -f app coverage.out

fmt:
	go fmt ./...

tidy:
	go mod tidy

deploy:
	./scripts/deploy.sh

# Generate Go code from proto files
proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth/auth.proto proto/core/core.proto proto/chat/chat.proto