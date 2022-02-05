build:
	go build -o ./bin/ozon-task ./cmd/ozon-task/
run:
	go run ./cmd/ozon-task/main.go
test:
	go test ./...
migrateup:
	migrate -path ./migrations -database "postgres://postgres:qwerty@database:5432/ozon_task?sslmode=disable" up
migratedown:
	migrate -path ./migrations -database "postgres://postgres:qwerty@database:5432/ozon_task?sslmode=disable" down
start:
	docker-compose up --build