FROM golang:1.23.1 AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0
RUN go install github.com/a-h/templ/cmd/templ@v0.2.432
RUN sqlc generate
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./app"]