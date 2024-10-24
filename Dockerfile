FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o jwt-authentication-golang

FROM alpine:3.16
WORKDIR /root/
COPY --from=builder /app/jwt-authentication-golang .

EXPOSE 8080
CMD ["./jwt-authentication-golang"]
