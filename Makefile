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

# Запуск всех тестов
test:
	@echo "Running all tests..."
	@go test ./... -v

# Покрытие всего кода с детальной статистикой
test-coverage:
	@echo "Running all tests with coverage..."
	@echo ""
	@go test ./... -cover 2>&1 | grep -E "(ok|FAIL)" | grep -v "compile: version"
	@echo ""
	@echo "Generating coverage report..."
	@go test ./... -coverprofile=coverage.out 2>&1 | grep -v "compile: version" | grep -v "no test files" > /dev/null
	@echo ""
	@echo "=== TOTAL COVERAGE ==="
	@go tool cover -func=coverage.out | grep total | awk '{printf "Total: %s\n", $$3}'
	@echo ""
	@echo "=== COVERAGE BY COMPONENT ==="
	@echo "Repositories:"
	@go tool cover -func=coverage.out | grep "repository/implementations" | awk '{sum+=$$NF; count++} END {if(count>0) printf "  %.1f%% (%d files)\n", sum/count, count; else print "  No coverage"}'
	@echo "Services:"
	@go tool cover -func=coverage.out | grep "service/implementations" | awk '{sum+=$$NF; count++} END {if(count>0) printf "  %.1f%% (%d files)\n", sum/count, count; else print "  No coverage"}'  
	@echo "Middleware:"
	@go tool cover -func=coverage.out | grep "handler/middleware" | awk '{sum+=$$NF; count++} END {if(count>0) printf "  %.1f%% (%d files)\n", sum/count, count; else print "  No coverage"}'
	@echo "Converters:"
	@go tool cover -func=coverage.out | grep "converters" | awk '{sum+=$$NF; count++} END {if(count>0) printf "  %.1f%% (%d files)\n", sum/count, count; else print "  No coverage"}'
	@echo "Handlers:"
	@go tool cover -func=coverage.out | grep "handler/http" | awk '{sum+=$$NF; count++} END {if(count>0) printf "  %.1f%% (%d files)\n", sum/count, count; else print "  No coverage"}'
	@echo "Adapters:"
	@go tool cover -func=coverage.out | grep "adapters" | awk '{sum+=$$NF; count++} END {if(count>0) printf "  %.1f%% (%d files)\n", sum/count, count; else print "  No coverage"}'
	@echo ""
	@echo "Full report: coverage.out"
	@echo "HTML report: go tool cover -html=coverage.out"

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

# Test только Repository (БД)
test-repository:
	@echo "Running repository tests..."
	@go test -v ./internal/repository/implementations/postgres/...
	@go test -v ./internal/repository/implementations/redis/...
	@go test -v ./internal/repository/implementations/minio/...

# Test с покрытием для Repository
test-repository-coverage:
	@echo "Running repository tests with coverage..."
	@go test -cover ./internal/repository/implementations/postgres/...
	@go test -cover ./internal/repository/implementations/redis/...
	@go test -cover ./internal/repository/implementations/minio/...
	@echo "\nDetailed coverage:"
	@go test -coverprofile=coverage-repository.out ./internal/repository/implementations/...
	@go tool cover -func=coverage-repository.out | grep total

# Test только Service (бизнес-логика)
test-service:
	@echo "Running service tests..."
	@go test -v ./internal/service/implementations/...

# Test с покрытием для Service
test-service-coverage:
	@echo "Running service tests with coverage..."
	@go test -cover ./internal/service/implementations/...
	@echo "\nDetailed coverage:"
	@go test -coverprofile=coverage-service.out ./internal/service/implementations/...
	@go tool cover -func=coverage-service.out | grep total

# Test только Middleware
test-middleware:
	@echo "Running middleware tests..."
	@go test ./internal/handler/middleware/... -v

# Test Middleware с покрытием
test-middleware-coverage:
	@echo "Running middleware tests with coverage..."
	@go test ./internal/handler/middleware/... -cover

# Test только Handlers
test-handler:
	@echo "Running handler tests..."
	@go test ./internal/handler/http/... -v

# Test Handlers с покрытием
test-handler-coverage:
	@echo "Running handler tests with coverage..."
	@go test ./internal/handler/http/... -cover

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

.PHONY: run run-auth run-core run-chat test test-coverage \
	test-converters test-converters-coverage \
	test-repository test-repository-coverage \
	test-service test-service-coverage \
	test-middleware test-middleware-coverage \
	test-handler test-handler-coverage \
	build-docs build build-bin clean fmt tidy deploy proto-gen
down:
	docker stop $(docker ps -q)
