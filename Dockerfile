FROM golang:alpine as builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./don ./cmd/ozon-task/*

FROM alpine

WORKDIR /app

COPY --from=builder /src/don ./
COPY --from=builder /src/config.yml ./

ENTRYPOINT ["/app/don", "-db", "redis"]