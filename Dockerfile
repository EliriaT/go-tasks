# Build stage
FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz


# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY start.sh .
COPY db/migrations ./migrations

ENV DB_SOURCE=mysql://user:password@tcp(sources_db:3306)/sources

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]