FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -mod=vendor -o ./build/api-gateway ./services/api-gateway

FROM alpine

WORKDIR /app
COPY --from=builder /app/build/api-gateway .

EXPOSE 8081
CMD ["/app/api-gateway"]