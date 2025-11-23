#GOOS=linux
#GOARCH=amd64

run:
	go run ./cmd/server/main.go

test:
	go test ./... -v

test-coverage:
	go test ./... -cover

build-docs:
	swag init -g /cmd/auth/main.go -o api/auth/
	swag init -g /cmd/core/main.go -o api/core/
	swag init -g /cmd/server/main.go -o api/server/

build:
	go build -o ./.build/auth ./cmd/auth/main.go
	go build -o ./.build/core ./cmd/core/main.go

build-bin:
	GOOS=linux GOARCH=amd64 go build -o ./.build/auth ./cmd/auth/main.go
	GOOS=linux GOARCH=amd64 go build -o ./.build/core ./cmd/core/main.go
	#go build -o ./.build/server ./cmd/server/main.go

clean:
	rm -f app coverage.out

fmt:
	go fmt ./...

tidy:
	go mod tidy

deploy:
	./scripts/deploy.sh

sbuild:
	git pull origin clean_arch
	go build -o $HOME/app/back/bin/ ./cmd/server/main.go

# Generate Go code from proto files
proto-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth/auth.proto proto/core/core.proto