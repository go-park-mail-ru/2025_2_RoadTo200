run:
	go run ./cmd/server/main.go

test:
	go test ./... -v

test-coverage:
	go test ./... -cover

build-docs:
	swag init -g /cmd/server/main.go -o docs/

build:
	go build -o ./../build/ ./cmd/server/main.go

clean:
	rm -f app coverage.out

fmt:
	go fmt ./...

tidy:
	go mod tidy
