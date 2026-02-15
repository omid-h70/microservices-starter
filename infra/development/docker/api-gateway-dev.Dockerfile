FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk update && \
    apk add git && \
    go install github.com/go-delve/delve/cmd/dlv@latest

# Debug build: disable optimizations so Delve works properly
RUN go build -gcflags "all=-N -l" -mod=vendor -o ./build/api-gateway ./services/api-gateway

EXPOSE 40000 8081
CMD ["dlv", "exec", "/app/build/api-gateway", "--headless", "--listen=:40000", "--api-version=2", "--accept-multiclient"]