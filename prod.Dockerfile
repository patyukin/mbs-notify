FROM golang:1.23.2-alpine3.20 AS builder

COPY . /app
WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go build -o ./bin/notify cmd/notify/main.go

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/bin/notify .
ENV YAML_CONFIG_FILE_PATH=config.yaml

CMD ["./notify"]
