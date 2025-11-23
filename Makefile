run:
	go run ./cmd/server/main.go

test:
	go test ./... -v

test-coverage:
	go test ./... -cover

build-docs:
	swag init -g /cmd/server/main.go -o api/docs/

build:
	go build -o ./.build/ ./cmd/server/main.go

clean:
	rm -f app coverage.out

fmt:
	go fmt ./...

tidy:
	go mod tidy

deploy:
	./scripts/deploy.sh

lbuild:
	git pull origin clean_arch
	go build -o $HOME/app/back/bin/ ./cmd/server/main.go

# Generate Go code from proto files
proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth/auth.proto proto/core/core.proto