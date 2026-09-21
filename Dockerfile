# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pulse-api main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=builder /app/pulse-api .

EXPOSE 8088

CMD ["./pulse-api"]