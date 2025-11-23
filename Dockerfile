FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go -o ./docs/swagger

RUN CGO_ENABLED=0 GOOS=linux go build -o app



FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]