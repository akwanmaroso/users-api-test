run:
	go run ./cmd/api/main.go

build:
	go build ./cmd/api/main.go

test:
	go test -cover ./...

swaggo:
	echo "Starting swagger generating"
	swag init -g **/**/*.go