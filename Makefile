build:
	go build -o ./bin ./cmd/ozon-task/
run:
	go run ./cmd/ozon-task/main.go
test:
	go test ./...
migrateup:
	migrate -path ./migrations -database "postgres://postgres:qwerty@database:5432/postgres?sslmode=disable" up
migratedown:
	migrate -path ./migrations -database "postgres://postgres:qwerty@database:5432/postgres?sslmode=disable" down