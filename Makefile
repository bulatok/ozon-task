postgres_dst := ${POSTGRES_URL}

build:
	go build -o don ./cmd/ozon-task/main.go

run:
	go run ./cmd/ozon-task/main.go

migrateup:
	migrate -path ./migrations -database ${postgres_dst} up

migratedown:
	migrate -path ./migrations -database ${postgres_dst} down

proto:
	protoc --proto_path=pkg/pb/proto/ --go_out=pkg/pb \
           --go_opt=paths=source_relative --go-grpc_out=pkg/pb \
           --go-grpc_opt=paths=source_relative pkg/pb/proto/links.proto

grpc-client:
	go build -o client/cmd/grpc/don client/cmd/grpc/main.go

mock:
	mockgen -source=internal/ozon-task/store/store.go \
		-destination=internal/ozon-task/store/mock.go -package store

start:
	docker-compose up --build