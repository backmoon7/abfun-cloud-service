# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

ARG GO_PROXY=https://goproxy.cn,direct
ENV GOPROXY=${GO_PROXY}

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

RUN go mod tidy

# Build the application
ARG SERVICE_NAME
RUN go build -o main cmd/${SERVICE_NAME}/main.go

# Run Stage
FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
