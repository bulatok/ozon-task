build:
	go build -o ./bin ./cmd/ozon-task/
run:
	go run ./cmd/ozon-task/main.go
test:
	go test ./...